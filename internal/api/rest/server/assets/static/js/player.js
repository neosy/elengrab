// -------------------------------------------------------------
// Player Overlay Logic
// Handles video/audio playback in modal overlay or bottom bar
// -------------------------------------------------------------

import * as watchAPI from './watch-api.js';
import { MEDIA_WATCH } from './constants.js';

let watchTracker = null;
let player = null

function initWatchTracker(video, itemId) {
    if (!video) return;
    if (watchTracker !== null) return;

    watchTracker = new watchAPI.MediaWatchTracker(video, itemId);
    watchTracker.init();        
}

async function destroyWatchTracker() {
    if (watchTracker === null) return;

    await watchTracker.destroy();
    watchTracker = null;
}

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

    initWatchTracker();

    document.addEventListener("click", async (event) => {
        const playBtn = event.target.closest(".media-result__play-button");
        if (!playBtn) return;

        const row = playBtn.closest(".media-result__row");
        if (!row) return;

        const itemId = row.dataset.itemId;

        const mediaURL = row.dataset.media;
        if (!mediaURL) return;

        document.dispatchEvent(new Event("preview-player-opened"));

        const isAudio = row.dataset.isAudio === "true";
        const shouldLoop = row.dataset.loop === "true";

        let positionMs = await watchAPI.getWatchPosition(itemId);

        if (positionMs < MEDIA_WATCH.startThresholdMs) {
            positionMs = 0;
        }

        // Clean previous players
        videoWrapper.innerHTML = "";
        audioBarContainer.innerHTML = "";

        let element;
        if (isAudio) {
            element = document.createElement("audio");
        } else {
            element = document.createElement("video");
            element.style.background = "black";
            // Disable Picture-in-Picture
            element.disablePictureInPicture = true;
        }

        player = element;

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
            initWatchTracker(element, itemId)
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

    async function closePlayer() {
        if (player !== null) {
            if (!player.paused) {
                await player.pause();
            }
            player = null;
        }

        isOpenVideoPlayer = false
        history.replaceState(null, "", location.pathname + location.search);

        videoWrapper.innerHTML = "";
        if (audioBarContainer) audioBarContainer.innerHTML = "";
        overlay.style.display = "none";
        document.body.style.overflow = "";
        document.body.classList.remove("audio-playing");

        destroyWatchTracker();
    }
}
