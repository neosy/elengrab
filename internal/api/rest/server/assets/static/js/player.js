// -------------------------------------------------------------
// Player Overlay Logic
// Handles video/audio playback in modal overlay or bottom bar
// -------------------------------------------------------------

const WATCH_TRACKING_URL_TEMPLATE = "/downloader/items/{itemId}/watch-tracking";
const WATCH_POSITION_URL_TEMPLATE = "/downloader/items/{itemId}/watch-position";

export function initPlayer() {
    const overlay           = document.getElementById("media-player-overlay");
    const videoWrapper      = document.getElementById("media-player__wrapper");
    let   audioBarContainer = document.getElementById("audio-bar");
    const playerHash        = "#player"
    let isOpenVideoPlayer   = false;

    if (!overlay || !videoWrapper) return;

    // Create audio container if missing
    if (!audioBarContainer) {
        audioBarContainer = document.createElement("div");
        audioBarContainer.id = "audio-bar";
        document.body.appendChild(audioBarContainer);
    }

    // Force initial hidden state
    overlay.style.display = "none";

    document.addEventListener("click", async (event) => {
        const playBtn = event.target.closest(".media-result__play-button");
        if (!playBtn) return;

        const row = playBtn.closest(".media-result__row");
        if (!row) return;

        const itemId = row.dataset.itemId;

        const mediaURL = row.dataset.media;
        if (!mediaURL) return;

        const isAudio = row.dataset.isAudio === "true";
        const shouldLoop = row.dataset.loop === "true";

        // Clean previous players
        videoWrapper.innerHTML = "";
        audioBarContainer.innerHTML = "";

        let positionMs = 0;

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

            positionMs = data.position;
        } catch (err) {
            console.error("Failed to get watch position", err);
        }        

        let element;
        if (isAudio) {
            element = document.createElement("audio");
        } else {
            element = document.createElement("video");
            element.style.background = "black";
            // Disable Picture-in-Picture
            element.disablePictureInPicture = true;
        }

        element.controls = true;
        element.autoplay = true;
        element.loop = shouldLoop;

        element.addEventListener("loadedmetadata", () => {
            if (positionMs > 0) {
                element.currentTime = positionMs / 1000;
            }
        }, { once: true });        

        element.src = mediaURL;

        if (isAudio) {
            // Audio → bottom fixed bar
            const bar = document.createElement("div");
            bar.className = "audio-player-bar";

            const closeBtn = document.createElement("button");
            closeBtn.className = "audio-close-btn";
            closeBtn.innerHTML = "×";
            closeBtn.setAttribute("aria-label", "Close audio player");
            closeBtn.onclick = closePlayer;

            bar.appendChild(closeBtn);
            bar.appendChild(element);
            audioBarContainer.appendChild(bar);

            overlay.style.display = "none !important";   // forceful hide
            document.body.style.overflow = "";
            document.body.classList.add("audio-playing");
        } else {
            isOpenVideoPlayer = true
            location.hash = playerHash

            // Video → centered overlay
            const wrapper = document.createElement("div");
            wrapper.className = "player-wrapper-modern";

            const closeBtn = document.createElement("button");
            closeBtn.className = "player-close-modern";
            closeBtn.innerHTML = "×";
            closeBtn.setAttribute("aria-label", "Close player");
            closeBtn.onclick = closePlayer;

            wrapper.appendChild(closeBtn);
            wrapper.appendChild(element);
            videoWrapper.appendChild(wrapper);

            overlay.style.display = "flex";
            document.body.style.overflow = "hidden";
        }

        if (itemId) {
            const watchTracker = new MediaWatchTracker(element, itemId);
            watchTracker.init();        
        }
    });

    // Handle middle-click on play button to open in new tab (for videos only)
    document.addEventListener('pointerdown', (e) => {
        if (e.button !== 1) return;

        const el = e.target;
        if (!(el instanceof Element)) return;

        const playBtn = el.closest(".media-result__play-button");
        if (!playBtn) return;

        const row = playBtn.closest(".media-result__row");
        const isAudio = row.dataset.isAudio === "true";

        if (isAudio) return; // To open in new tab only applies to videos
        
        if (e.button === 1) {
            e.preventDefault();
            window.open(playBtn.dataset.watchUrl, '_blank', 'noopener,noreferrer');
            return;
        }
    });

    // Close overlay on background click (video only)
    overlay.addEventListener("click", (event) => {
        if (event.target === overlay) closePlayer();
    });

    // Global ESC handler
    document.addEventListener("keydown", (event) => {
        if (event.key === "Escape" && location.hash === playerHash) closePlayer();
    });

    if (location.hash === playerHash) {
        history.replaceState(null, "", location.pathname + location.search);
    }

    window.addEventListener('hashchange', syncPlayerWithHash);

    function syncPlayerWithHash() {
        if (location.hash === playerHash) {
            return;
        }

        if (isOpenVideoPlayer) {
            closePlayer();
        }
    }    

    function closePlayer() {
        isOpenVideoPlayer = false
        history.replaceState(null, "", location.pathname + location.search);

        videoWrapper.innerHTML = "";
        if (audioBarContainer) audioBarContainer.innerHTML = "";
        overlay.style.display = "none";
        document.body.style.overflow = "";
        document.body.classList.remove("audio-playing");
    }

    function getWatchPositionUrl(itemId) {
        return WATCH_POSITION_URL_TEMPLATE.replace(
            "{itemId}",
            itemId
        );
    }

}

export const MEDIA_WATCH_EVENT_PAUSE = "pause";
export const MEDIA_WATCH_EVENT_ENDED = "ended";
export const MEDIA_WATCH_EVENT_SEEK = "seek";
export const MEDIA_WATCH_EVENT_HEARTBEAT = "heartbeat";

export class MediaWatchTracker {
    constructor(video, itemId) {
        this.video = video;
        this.itemId = itemId;

        this.heartbeatInterval = 5000;
        this.heartbeatTimer = null;
        this.lastSentAt = null;
        this.lastPosition = null;

        this.stopPromise = null;
    }

    getWatchEventUrl() {
        return WATCH_TRACKING_URL_TEMPLATE.replace(
            "{itemId}",
            this.itemId
        );
    }

    async sendWatchEvent(eventType = null) {
        if (!this.video || this.lastSentAt === null) {
            return;
        }

        if (this.lastPosition === null) {
            return;
        }

        const now = Date.now();
        const intervalMs = now - this.lastSentAt;

        // Ignore empty or invalid intervals.
        if (intervalMs <= 0) {
            return;
        }

        let positionMs = Math.floor(this.lastPosition * 1000);
        if (positionMs < intervalMs) {
            positionMs = intervalMs;
        }

        if (eventType === MEDIA_WATCH_EVENT_ENDED) {
            positionMs = 0;
        }

        const event = {
            positionMs: positionMs,
            intervalMs,
            eventType,
        };

        try {
            await fetch(this.getWatchEventUrl(), {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(event),
            });

            this.lastSentAt = now;
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

        // Start measuring playback time.
        this.lastSentAt = Date.now();

        // Store the current position to track playback progress.
        this.lastPosition = this.video.currentTime;

        this.heartbeatTimer = setInterval(() => {
            if (!this.video.paused && !this.video.ended) {
                this.sendWatchEvent(MEDIA_WATCH_EVENT_HEARTBEAT);
            }
        }, this.heartbeatInterval);
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
        await this.sendWatchEvent(type);

        // Reset the heartbeat timer to avoid sending duplicate events.
        this.heartbeatTimer = null;

        // Reset the last sent timestamp to avoid sending duplicate events.
        this.lastSentAt = null;

        // Reset the last position to avoid sending duplicate events.
        this.lastPosition = null;
    }

    init() {
        this.video.addEventListener("play", () => {
            this.startHeartbeat();
        });

        this.video.addEventListener("pause", () => {
            if (this.video.seeking) {
                this.stopHeartbeat(MEDIA_WATCH_EVENT_SEEK);
            } else {
                this.stopHeartbeat(MEDIA_WATCH_EVENT_PAUSE);
            }
        });

        this.video.addEventListener("ended", () => {
            this.stopHeartbeat(MEDIA_WATCH_EVENT_ENDED);
        });

        this.video.addEventListener("timeupdate", () => {
            if (!this.video.seeking) {
                this.lastPosition = this.video.currentTime;
            }
        });

        this.video.addEventListener("seeked", async () => {
            if (!this.video.paused && !this.video.ended && this.video.currentTime === 0) {
                await this.stopHeartbeat(MEDIA_WATCH_EVENT_ENDED);
                this.startHeartbeat();
            }
        });
    }
}