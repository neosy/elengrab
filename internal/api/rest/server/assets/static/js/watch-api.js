import { MEDIA_WATCH } from './constants.js';

const WATCH_TRACKING_URL_TEMPLATE = "/downloader/items/{itemId}/watch-tracking";
const WATCH_POSITION_URL_TEMPLATE = "/downloader/items/{itemId}/watch-position";

function getWatchPositionUrl(itemId) {
    return WATCH_POSITION_URL_TEMPLATE.replace(
        "{itemId}",
        itemId
    );
}

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
}

const HEARTBEAT_INTERVAL_MS = 5000;
const LOOP_THRESHOLD = 0.2;

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

        this.heartbeatTimer = null;

        this.currentPosition = null;
        this.lastSentPosition = null;

        this.stopPromise = null;
    }

    setItemId(itemId) {
        this.stopHeartbeatInternal();
        this.itemId = itemId || null;
    }
        
    getWatchEventUrl() {
        if (!this.itemId) {
            return null;
        }

        return WATCH_TRACKING_URL_TEMPLATE.replace(
            "{itemId}",
            this.itemId
        );
    }

    async sendWatchEvent(eventType = null) {
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

        // Ignore invalid large intervals.
        if (intervalMs > MEDIA_WATCH.maxIntervalMs) {
            return;
        }

        let positionMs = Math.floor(this.currentPosition * 1000);

        if (eventType === MEDIA_WATCH_EVENT.ended) {
            positionMs = 0;
        }

        const event = {
            positionMs: positionMs,
            intervalMs,
            eventType,
        };

        const url = this.getWatchEventUrl();

        if (!url) {
            return;
        }        

        try {
            await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(event),
            });

            this.lastSentPosition = this.currentPosition;
        } catch (err) {
            console.error("Failed to send media watch event", err);
        }
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
            this.lastSentPosition = this.currentPosition
        }

        this.heartbeatTimer = setInterval(() => {
            if (!this.video.paused && !this.video.ended) {
                this.sendWatchEvent(MEDIA_WATCH_EVENT.heartbeat);
            }
        }, HEARTBEAT_INTERVAL_MS);
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
        clearInterval(this.heartbeatTimer);

        // Send the remaining playback time.
        if (type) {
            await this.sendWatchEvent(type);
        }

        // Reset the heartbeat timer to avoid sending duplicate events.
        this.heartbeatTimer = null;

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
                this.lastSentPosition = null;
            }

            if (!this.video.paused &&
                !this.video.ended &&
                this.video.currentTime <= LOOP_THRESHOLD) {

                await this.stopHeartbeat(MEDIA_WATCH_EVENT.ended);
                this.startHeartbeat();
            }
        };

        this.video.addEventListener("play", this.onPlay);
        this.video.addEventListener("pause", this.onPause);
        this.video.addEventListener("ended", this.onEnded);
        this.video.addEventListener("timeupdate", this.onTimeUpdate);
        this.video.addEventListener("seeked", this.onSeeked);

        // [
        //     "play",
        //     "playing",
        //     "pause",
        //     "ended",
        //     "seeking",
        //     "seeked",
        //     "timeupdate"
        // ].forEach(name => {
        //     this.video.addEventListener(name, () => {
        //         console.log(name, this.video.currentTime);
        //     });
        // });        
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