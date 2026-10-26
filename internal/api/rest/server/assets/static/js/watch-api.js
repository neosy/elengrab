import { MEDIA_WATCH, API_PATH_TEMPLATES } from './constants.js';

/**
 * Gets the saved watch position for a media item.
 *
 * @param {string} itemId - Media item ID.
 */
export async function getWatchPosition(itemId) {
    if (!itemId) return 0;

    try {
        const response = await fetch(getWatchPositionUrl(itemId), {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (!response.ok) {
            throw new Error(`Failed to get watch position: ${response.status}`);
        }

        const data = await response.json();

        return data.position;
    } catch (err) {
        console.error("Failed to get watch position", err);
        return 0;
    }

    function getWatchPositionUrl(itemId) {
        return API_PATH_TEMPLATES.downloaderWatchPosition.replace(
            "{itemId}",
            itemId
        );
    }
}

const LOOP_THRESHOLD = 0.2;
const LOOP_END_THRESHOLD = 1;

// Extra time on top of minIntervalMs for the early send,
// because currentPosition is updated only on timeupdate (every ~250 ms).
const EARLY_SEND_MARGIN_MS = 500;

// Media watch event types
export const MEDIA_WATCH_EVENT = {
    pause: "pause",
    ended: "ended",
    seek: "seek",
    heartbeat: "heartbeat",
};

export class MediaWatchTracker {
    constructor(video, itemId) {
        this.video = video;

        this.itemId = null;
        if (itemId) {
            this.itemId = itemId;
        }

        // Single timer of the send chain (early send, then regular heartbeats).
        this.heartbeatTimer = null;

        this.currentPosition = null;
        this.lastSentPosition = null;

        this.stopPromise = null;

        // Incremented whenever the playback baseline is reset (seek).
        this.baselineEpoch = 0;

        // Queue that serializes sending of watch events.
        this.sendChain = Promise.resolve();
    }

    async setItemId(itemId) {
        // Send the remaining time to the previous item first.
        await this.stopHeartbeat(MEDIA_WATCH_EVENT.pause);

        // The new item starts counting from scratch.
        this.lastSentPosition = null;
        this.itemId = itemId || null;
    }

    getWatchEventUrl() {
        if (!this.itemId) {
            return null;
        }

        return API_PATH_TEMPLATES.downloaderWatchTracking.replace(
            "{itemId}",
            this.itemId
        );
    }

    sendWatchEvent(eventType = null) {
        this.sendChain = this.sendChain
            .then(() => this.sendWatchEventInternal(eventType))
            .catch((err) => console.error("Media watch send failed", err));
        return this.sendChain;
    }

    async sendWatchEventInternal(eventType = null) {
        if (!this.video) {
            return;
        }

        if (this.currentPosition === null || this.lastSentPosition === null) {
            return;
        }

        const intervalMs = Math.floor((this.currentPosition - this.lastSentPosition) * 1000);

        // Ignore empty or invalid intervals.
        if (intervalMs <= 0) {
            return;
        }

        // Ignore very short watch intervals for intermediate events.
        // The "ended" event is always sent to ensure the final playback state is recorded.
        if (eventType !== MEDIA_WATCH_EVENT.ended && intervalMs < MEDIA_WATCH.minIntervalMs) {
            return;
        }

        // Ignore invalid large intervals and restart counting from the current position.
        if (intervalMs > MEDIA_WATCH.maxIntervalMs) {
            this.lastSentPosition = this.currentPosition;
            return;
        }

        let positionMs = Math.floor(this.currentPosition * 1000);

        if (eventType === MEDIA_WATCH_EVENT.ended) {
            positionMs = 0;
        }

        const event = {
            positionMs,
            intervalMs,
            eventType,
        };

        const url = this.getWatchEventUrl();

        if (!url) {
            return;
        }

        // Remember the position that is actually being sent.
        const sentPosition = this.currentPosition;
        const epoch = this.baselineEpoch;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(event),
            });

            if (!response.ok) {
                throw new Error(`Failed to send media watch event: ${response.status}`);
            }

            // Skip the update if a seek re-baselined the positions while the request was in flight.
            if (epoch === this.baselineEpoch) {
                this.lastSentPosition = sentPosition;
            }
        } catch (err) {
            console.error("Failed to send media watch event", err);
        }
    }

    // Schedules the next send. After it fires, the chain continues
    // with a full heartbeat interval until the heartbeat is stopped.
    scheduleNext(delayMs) {
        clearTimeout(this.heartbeatTimer);

        this.heartbeatTimer = setTimeout(() => {
            if (this.video && !this.video.paused && !this.video.ended) {
                this.sendWatchEvent(MEDIA_WATCH_EVENT.heartbeat);
            }

            // heartbeatTimer is null only after stopHeartbeatInternal has started.
            if (this.heartbeatTimer !== null) {
                this.scheduleNext(MEDIA_WATCH.heartbeatIntervalMs);
            }
        }, delayMs);
    }

    // Early send: the new position is saved sooner than a regular heartbeat.
    scheduleEarlySend() {
        this.scheduleNext(MEDIA_WATCH.minIntervalMs + EARLY_SEND_MARGIN_MS);
    }

    async startHeartbeat() {
        if (this.stopPromise !== null) {
            await this.stopPromise;
        }

        if (this.heartbeatTimer !== null) {
            return;
        }

        // Initialize the current position and last sent position to the current playback time.
        this.currentPosition = this.video.currentTime;

        if (this.lastSentPosition === null || this.lastSentPosition > this.currentPosition) {
            this.lastSentPosition = this.currentPosition;
        }

        this.scheduleEarlySend();
    }

    async stopHeartbeat(type = null) {
        if (this.stopPromise !== null) {
            await this.stopPromise;
            return;
        }

        if (this.heartbeatTimer === null) {
            return;
        }

        this.stopPromise = this.stopHeartbeatInternal(type);

        try {
            await this.stopPromise;
        } finally {
            this.stopPromise = null;
        }
    }

    async stopHeartbeatInternal(type = null) {
        // Stop the chain right away, so no new send is scheduled.
        clearTimeout(this.heartbeatTimer);
        this.heartbeatTimer = null;

        // Send the remaining playback time.
        if (type) {
            await this.sendWatchEvent(type);
        }

        // Reset the current position to null to avoid sending duplicate events.
        this.currentPosition = null;

        if (type !== MEDIA_WATCH_EVENT.pause) {
            this.lastSentPosition = null;
        }
    }

    init() {
        this.onPlay = () => {
            this.startHeartbeat();
        };

        this.onPause = async () => {
            if (!this.video.seeking && this.currentPosition !== null) {
                this.currentPosition = this.video.currentTime;
            }

            if (this.video.ended) {
                this.stopHeartbeat(MEDIA_WATCH_EVENT.ended);
            }
            else if (this.video.seeking) {
                this.stopHeartbeat(MEDIA_WATCH_EVENT.seek);
            } else {
                this.stopHeartbeat(MEDIA_WATCH_EVENT.pause);
            }
        };

        this.onEnded = () => {
            if (this.currentPosition !== null) {
                this.currentPosition = this.video.currentTime;
            }

            this.stopHeartbeat(MEDIA_WATCH_EVENT.ended);
        };

        this.onTimeUpdate = () => {
            if (
                !this.video.seeking &&
                this.currentPosition !== null &&
                this.video.currentTime >= this.currentPosition
            ) {
                this.currentPosition = this.video.currentTime;
            }
        };

        this.onSeeked = async () => {
            if (this.video.paused) {
                this.baselineEpoch++;
                this.lastSentPosition = null;
                return;
            }

            if (this.video.ended) {
                return;
            }

            // Looped video: playback jumped from the end back to the start.
            const wasNearEnd =
                this.currentPosition !== null &&
                this.video.duration - this.currentPosition <= LOOP_END_THRESHOLD;

            if (this.video.loop && wasNearEnd && this.video.currentTime <= LOOP_THRESHOLD) {
                await this.stopHeartbeat(MEDIA_WATCH_EVENT.ended);
                this.startHeartbeat();
                return;
            }

            // Playback continues after a seek: re-baseline both positions,
            // so the jumped range is not counted as watched time.
            if (this.heartbeatTimer !== null) {
                this.baselineEpoch++;
                this.currentPosition = this.video.currentTime;
                this.lastSentPosition = this.video.currentTime;

                // Save the new position soon, without waiting for the regular heartbeat.
                this.scheduleEarlySend();
            }
        };

        this.video.addEventListener("play", this.onPlay);
        this.video.addEventListener("pause", this.onPause);
        this.video.addEventListener("ended", this.onEnded);
        this.video.addEventListener("timeupdate", this.onTimeUpdate);
        this.video.addEventListener("seeked", this.onSeeked);
    }

    async destroy() {
        await this.stopHeartbeat(MEDIA_WATCH_EVENT.pause);

        if (!this.video) {
            return;
        }

        this.video.removeEventListener("play", this.onPlay);
        this.video.removeEventListener("pause", this.onPause);
        this.video.removeEventListener("ended", this.onEnded);
        this.video.removeEventListener("timeupdate", this.onTimeUpdate);
        this.video.removeEventListener("seeked", this.onSeeked);

        this.video = null;
        this.itemId = null;
    }
}