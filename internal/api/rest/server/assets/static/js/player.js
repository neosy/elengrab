// -------------------------------------------------------------
// Player Overlay Logic
// Handles video/audio playback in modal overlay or bottom bar
// -------------------------------------------------------------

import * as watchAPI from './watch-api.js';
import { CLASS_NAMES, MEDIA_WATCH, VIDEO_PREVIEW } from './constants.js';

export const LOCAL_CLASS_NAMES = {
    mediaResultPlayButton: "media-play-button",
    mediaResultRow: "media-result__row",

    mediaPlayerWrapper: "media-player__wrapper",
    mediaPlayerClose: "media-player__close",
    showControls: "show-controls",

    mediaPlayerVideo: "media-player__video",
    mediaPlayerAudio: "media-player__audio",

    audioPlaying: "audio-playing",
};

export const LOCAL_CLASS_SELECTORS = Object.fromEntries(
    Object.entries(LOCAL_CLASS_NAMES)
        .filter(([, value]) => typeof value === "string")
        .map(([key, value]) => [key, `.${value}`])
);

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

/**
 * @param {HTMLElement} playerContainer Media player container.
 */
export function initPlayer(playerContainer) {
    if (!playerContainer) return;

    const videoContainer = playerContainer.querySelector(LOCAL_CLASS_SELECTORS.mediaPlayerVideo);
    const audioContainer = playerContainer.querySelector(LOCAL_CLASS_SELECTORS.mediaPlayerAudio);

    if (!videoContainer || !audioContainer) return;
    
    const playerHash        = "#player"
    let isOpenVideoPlayer   = false;
    let cleanupControls = null;

    // Create audio container if missing
    if (!audioContainer) {
        audioContainer = document.createElement("div");
        audioContainer.className = LOCAL_CLASS_NAMES.mediaPlayerAudio;
        playerContainer.appendChild(audioContainer);
    }

    // Force initial hidden state
    videoContainer.style.display = "none";

    initWatchTracker();

    document.addEventListener("click", async (event) => {
        const playBtn = event.target.closest(LOCAL_CLASS_SELECTORS.mediaResultPlayButton);
        if (!playBtn) return;

        if (isOpenVideoPlayer) await closePlayer();

        const row = playBtn.closest(LOCAL_CLASS_SELECTORS.mediaResultRow);
        if (!row) return;

        const itemId = row.dataset.itemId;

        const mediaURL = row.dataset.media;
        if (!mediaURL) return;

        document.dispatchEvent(new Event(VIDEO_PREVIEW.playerOpenedEventName));

        const isAudio = row.dataset.isAudio === "true";
        const shouldLoop = row.dataset.loop === "true";

        let positionMs = await watchAPI.getWatchPosition(itemId);

        if (positionMs < MEDIA_WATCH.startThresholdMs) {
            positionMs = 0;
        }

        // Clean previous players
        videoContainer.innerHTML = "";
        audioContainer.innerHTML = "";

        if (isAudio) {
            player = document.createElement("audio");
        } else {
            player = document.createElement("video");
            player.style.background = "black";
            // Disable Picture-in-Picture
            player.disablePictureInPicture = true;
        }

        player.controls = true;
        player.autoplay = true;
        player.loop = shouldLoop;

        player.addEventListener("loadedmetadata", () => {
            if (positionMs > 0) {
                player.currentTime = positionMs / 1000;
            }
        }, { once: true });        

        player.src = mediaURL;

        isOpenVideoPlayer = true

        if (isAudio) {
            // Audio → bottom fixed bar
            const wrapper = document.createElement("div");
            wrapper.className = LOCAL_CLASS_NAMES.mediaPlayerWrapper;

            const closeBtn = document.createElement("button");
            closeBtn.className = LOCAL_CLASS_NAMES.mediaPlayerClose;
            closeBtn.innerHTML = "×";
            closeBtn.setAttribute("aria-label", "Close audio player");
            closeBtn.onclick = closePlayer;

            wrapper.appendChild(player);
            wrapper.appendChild(closeBtn);

            audioContainer.appendChild(wrapper);

            player.focus({ preventScroll: true });

            videoContainer.style.display = "none !important";   // forceful hide
            document.body.style.overflow = "";
            document.body.classList.add(LOCAL_CLASS_NAMES.audioPlaying);
        } else {
            document.documentElement.classList.add(CLASS_NAMES.ui.blockingActive);
            
            location.hash = playerHash

            // Video → centered overlay
            const wrapper = document.createElement("div");
            wrapper.className = LOCAL_CLASS_NAMES.mediaPlayerWrapper;

            wrapper.appendChild(player);
            videoContainer.appendChild(wrapper);

            cleanupControls = initVideoControls();

            videoContainer.style.display = "flex";

            player.focus({ preventScroll: true });
        }

        if (itemId) {
            initWatchTracker(player, itemId)
        }
    });

    // Handle middle-click on play button to open in new tab (for videos only)
    document.addEventListener('pointerdown', (e) => {
        if (e.button !== 1) return;

        const el = e.target;
        if (!(el instanceof Element)) return;

        const playBtn = el.closest(LOCAL_CLASS_SELECTORS.mediaResultPlayButton);
        if (!playBtn) return;

        const row = playBtn.closest(LOCAL_CLASS_SELECTORS.mediaResultRow);
        const isAudio = row.dataset.isAudio === "true";

        if (isAudio) return; // To open in new tab only applies to videos
        
        if (e.button === 1) {
            e.preventDefault();
            window.open(playBtn.dataset.watchUrl, '_blank', 'noopener,noreferrer');
            return;
        }
    });

    // Close overlay on background click (video only)
    videoContainer.addEventListener("click", (event) => {
        if (event.target === videoContainer) closePlayer();
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
        }

        isOpenVideoPlayer = false;

        history.replaceState(null, "", location.pathname + location.search);

        cleanupControls?.();
        cleanupControls = null;

        player = null;

        if (videoContainer) videoContainer.innerHTML = "";
        if (audioContainer) audioContainer.innerHTML = "";
        videoContainer.style.display = "none";
        document.body.style.overflow = "";
        document.body.classList.remove(LOCAL_CLASS_NAMES.audioPlaying);

        document.documentElement.classList.remove(CLASS_NAMES.ui.blockingActive);

        destroyWatchTracker();
    }
    
    function initVideoControls() {
        const wrapper = videoContainer.querySelector(LOCAL_CLASS_SELECTORS.mediaPlayerWrapper);

        if (!wrapper) return;

        let controlsTimeout;

        const closeBtn = document.createElement("button");
        closeBtn.className = LOCAL_CLASS_NAMES.mediaPlayerClose;
        closeBtn.innerHTML = "×";
        closeBtn.setAttribute("aria-label", "Close player");
        closeBtn.onclick = closePlayer;

        wrapper.appendChild(closeBtn);

        const showControls = () => {
            if (!isOpenVideoPlayer) {
                return;
            }

            videoContainer.classList.add(LOCAL_CLASS_NAMES.showControls);

            clearTimeout(controlsTimeout);

            controlsTimeout = setTimeout(() => {
                videoContainer.classList.remove(LOCAL_CLASS_NAMES.showControls);
            }, 3000);
        };

        videoContainer.addEventListener("mouseenter", showControls);
        player.addEventListener("mousemove", showControls);
        player.addEventListener("play", showControls);
        player.addEventListener("pause", showControls);
        player.addEventListener("click", showControls);
        player.addEventListener("touchstart", showControls);

        showControls();
        
        return () => {
            videoContainer.removeEventListener("mouseenter", showControls);
            player.removeEventListener("mousemove", showControls);
            player.removeEventListener("play", showControls);
            player.removeEventListener("pause", showControls);
            player.removeEventListener("click", showControls);
            player.removeEventListener("touchstart", showControls);

            clearTimeout(controlsTimeout);

            videoContainer.classList.remove(LOCAL_CLASS_NAMES.showControls);
        };
    }
}
