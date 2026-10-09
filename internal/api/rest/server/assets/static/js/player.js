// -------------------------------------------------------------
// Player Overlay Logic
// Handles video/audio playback in modal overlay or bottom bar
// -------------------------------------------------------------

import * as watchAPI from './watch-api.js';
import { CLASS_NAMES as CONST_CLASS_NAMES, MEDIA_WATCH, EVENT_NAMES } from './constants.js';

export const CLASS_NAMES = {
    ...CONST_CLASS_NAMES,

    mediaResultPlayButton: "media-play-button",
    mediaResultRow: "media-result__row",

    mediaPlayerWrapper: "media-player__wrapper",
    mediaPlayerClose: "media-player__close",
    showControls: "show-controls",

    mediaPlayerVideo: "media-player__video",
    mediaPlayerAudio: "media-player__audio",

    audioPlaying: "audio-playing",
    isAudioPaused: 'is-audio-paused',

    audioPlayingOverlay: 'media-player__audio-playing-overlay',
};

export const CLASS_SELECTORS = Object.fromEntries(
    Object.entries(CLASS_NAMES)
        .filter(([, value]) => typeof value === "string")
        .map(([key, value]) => [key, `.${value}`])
);

const DOM_IDS = {
    audioPlayingOverlay: "media-player-audio-playing-overlay",
};

const DOM_ELEMENTS = {
    audioPlayingOverlay: null,
}  

function initDomElements(elements = DOM_ELEMENTS) {    
    elements.audioPlayingOverlay = document.getElementById(DOM_IDS.audioPlayingOverlay);
}

let watchTracker = null;

const mediaPlayer = {
    isAudio: false,
    isOpen: false,

    player: null,

    playingRow: null,
    rowItemId: "",

    clear() {
        this.isAudio = false;
        this.isOpen = false;

        if (this.player) {
            const unloadVideo = () => {
                this.player.removeAttribute("src");
                this.player.load();
            };

            if (this.player.paused) {
                unloadVideo();
            } else {
                this.player.addEventListener("pause", unloadVideo, { once: true });
                this.player.pause();
            }    

            this.player.remove();
            this.player = null;
        }

        this.playingRow = null;
        this.rowItemId = "";
    },
};

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
 * @param {string} overlayTargetClassName Class name of the element to move the audio overlay into.
 */
export function initPlayer(playerContainer, overlayTargetClassName) {
    if (!playerContainer) return;

    initDomElements();

    const videoContainer = playerContainer.querySelector(CLASS_SELECTORS.mediaPlayerVideo);
    const audioContainer = playerContainer.querySelector(CLASS_SELECTORS.mediaPlayerAudio);

    if (!videoContainer || !audioContainer) return;
    
    const videoPlayerHash        = "#player"
    let cleanupControls = null;

    // Create audio container if missing
    if (!audioContainer) {
        audioContainer = document.createElement("div");
        audioContainer.className = CLASS_NAMES.mediaPlayerAudio;
        playerContainer.appendChild(audioContainer);
    }

    // Force initial hidden state
    videoContainer.style.display = "none";

    initWatchTracker();

    document.addEventListener("click", async (event) => {
        const playBtn = event.target.closest(CLASS_SELECTORS.mediaResultPlayButton);
        if (!playBtn) return;

        const row = playBtn.closest(CLASS_SELECTORS.mediaResultRow);
        if (!row) return;

        const itemId = row.dataset.itemId;

        if (mediaPlayer.isOpen) {
            if (mediaPlayer.isAudio && itemId === mediaPlayer.rowItemId) {
                mediaPlayer.player.paused ? mediaPlayer.player.play() : mediaPlayer.player.pause();
                return;
            }

            await closePlayer();
        }

        mediaPlayer.rowItemId = itemId;
        mediaPlayer.playingRow = row;

        const mediaURL = mediaPlayer.playingRow.dataset.media;
        if (!mediaURL) return;

        mediaPlayer.isAudio = mediaPlayer.playingRow.dataset.isAudio === "true";

        if (!mediaPlayer.isAudio) {
            document.dispatchEvent(new Event(EVENT_NAMES.videoPlayerOpened));
        }

        const shouldLoop = mediaPlayer.playingRow.dataset.loop === "true";

        let positionMs = await watchAPI.getWatchPosition(itemId);

        if (positionMs < MEDIA_WATCH.startThresholdMs) {
            positionMs = 0;
        }

        // Clean previous players
        videoContainer.innerHTML = "";
        audioContainer.innerHTML = "";

        if (mediaPlayer.isAudio) {
            mediaPlayer.player = document.createElement("audio");
        } else {
            mediaPlayer.player = document.createElement("video");
            mediaPlayer.player.style.background = "black";
            // Disable Picture-in-Picture
            mediaPlayer.player.disablePictureInPicture = true;
        }

        mediaPlayer.player.controls = true;
        mediaPlayer.player.autoplay = true;
        mediaPlayer.player.loop = shouldLoop;

        mediaPlayer.player.addEventListener("loadedmetadata", () => {
            if (positionMs > 0) {
                mediaPlayer.player.currentTime = positionMs / 1000;
            }
        }, { once: true });        

        mediaPlayer.isOpen = true
        mediaPlayer.player.src = mediaURL;

        if (mediaPlayer.isAudio) {
            // Audio → bottom fixed bar
            const wrapper = document.createElement("div");
            wrapper.className = CLASS_NAMES.mediaPlayerWrapper;

            const closeBtn = document.createElement("button");
            closeBtn.className = CLASS_NAMES.mediaPlayerClose;
            closeBtn.innerHTML = "×";
            closeBtn.setAttribute("aria-label", "Close audio player");
            closeBtn.onclick = closePlayer;

            wrapper.appendChild(mediaPlayer.player);
            wrapper.appendChild(closeBtn);

            audioContainer.appendChild(wrapper);

            cleanupControls = initAudioControls();

            mediaPlayer.player.focus({ preventScroll: true });

            videoContainer.style.display = "none !important";   // forceful hide
            document.body.style.overflow = "";

            const overlayTarget = mediaPlayer.playingRow.querySelector(`.${overlayTargetClassName}`);
            if (overlayTarget) {
                showAudioPlayingOverlay(overlayTarget);
            }

            requestAnimationFrame(() => {
                document.body.classList.add(CLASS_NAMES.audioPlaying);
            });
        } else {
            document.documentElement.classList.add(CLASS_NAMES.ui.blockingActive);
            
            location.hash = videoPlayerHash

            // Video → centered overlay
            const wrapper = document.createElement("div");
            wrapper.className = CLASS_NAMES.mediaPlayerWrapper;

            wrapper.appendChild(mediaPlayer.player);
            videoContainer.appendChild(wrapper);

            cleanupControls = initVideoControls();

            videoContainer.style.display = "flex";

            mediaPlayer.player.focus({ preventScroll: true });
        }

        if (itemId) {
            initWatchTracker(mediaPlayer.player, itemId)
        }
    });

    // Handle middle-click on play button to open in new tab (for videos only)
    document.addEventListener('pointerdown', (e) => {
        if (e.button !== 1) return;

        const el = e.target;
        if (!(el instanceof Element)) return;

        const playBtn = el.closest(CLASS_SELECTORS.mediaResultPlayButton);
        if (!playBtn) return;

        const row = playBtn.closest(CLASS_SELECTORS.mediaResultRow);
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
        if (event.key === "Escape" && location.hash === videoPlayerHash) closePlayer();
    });

    if (location.hash === videoPlayerHash) {
        history.replaceState(null, "", location.pathname + location.search);
    }

    window.addEventListener('hashchange', syncPlayerWithHash);

    function syncPlayerWithHash() {
        if (mediaPlayer.isAudio) {
            return;
        }

        if (location.hash === videoPlayerHash) {
            return;
        }

        if (mediaPlayer.isOpen) {
            closePlayer();
        }
    }    

    async function closePlayer() {
        if (mediaPlayer.player !== null) {
            if (!mediaPlayer.player.paused) {
                await mediaPlayer.player.pause();
            }
        }

        history.replaceState(null, "", location.pathname + location.search);

        cleanupControls?.();
        cleanupControls = null;

        mediaPlayer.clear();

        if (videoContainer) videoContainer.innerHTML = "";
        if (audioContainer) audioContainer.innerHTML = "";

        videoContainer.style.display = "none";
        document.body.style.overflow = "";

        document.body.classList.remove(CLASS_NAMES.audioPlaying);
        document.body.classList.remove(CLASS_NAMES.isAudioPaused);

        document.documentElement.classList.remove(CLASS_NAMES.ui.blockingActive);

        destroyWatchTracker();
    }

    function initAudioControls() {
        const handlePlay = () => {
            document.body.classList.remove(CLASS_NAMES.isAudioPaused);
        };

        const handlePause = () => {
            document.body.classList.add(CLASS_NAMES.isAudioPaused);
        };

        mediaPlayer.player.addEventListener("play", handlePlay);
        mediaPlayer.player.addEventListener("pause", handlePause);

        return () => {
            mediaPlayer.player.removeEventListener("play", handlePlay);
            mediaPlayer.player.removeEventListener("pause", handlePause);
        };
    }
    
    function initVideoControls() {
        const wrapper = videoContainer.querySelector(CLASS_SELECTORS.mediaPlayerWrapper);

        if (!wrapper) return;

        let controlsTimeout;

        const closeBtn = document.createElement("button");
        closeBtn.className = CLASS_NAMES.mediaPlayerClose;
        closeBtn.innerHTML = "×";
        closeBtn.setAttribute("aria-label", "Close player");
        closeBtn.onclick = closePlayer;

        wrapper.appendChild(closeBtn);

        const showControls = () => {
            if (!mediaPlayer.isOpen) {
                return;
            }

            videoContainer.classList.add(CLASS_NAMES.showControls);

            clearTimeout(controlsTimeout);

            controlsTimeout = setTimeout(() => {
                videoContainer.classList.remove(CLASS_NAMES.showControls);
            }, 3000);
        };

        videoContainer.addEventListener("mouseenter", showControls);
        mediaPlayer.player.addEventListener("mousemove", showControls);
        mediaPlayer.player.addEventListener("play", showControls);
        mediaPlayer.player.addEventListener("pause", showControls);
        mediaPlayer.player.addEventListener("click", showControls);
        mediaPlayer.player.addEventListener("touchstart", showControls);

        showControls();
        
        return () => {
            videoContainer.removeEventListener("mouseenter", showControls);
            mediaPlayer.player.removeEventListener("mousemove", showControls);
            mediaPlayer.player.removeEventListener("play", showControls);
            mediaPlayer.player.removeEventListener("pause", showControls);
            mediaPlayer.player.removeEventListener("click", showControls);
            mediaPlayer.player.removeEventListener("touchstart", showControls);

            clearTimeout(controlsTimeout);

            videoContainer.classList.remove(CLASS_NAMES.showControls);
        };
    }

    function showAudioPlayingOverlay(target) {
        if (!target || !DOM_ELEMENTS.audioPlayingOverlay) return;

        target.append(DOM_ELEMENTS.audioPlayingOverlay);
    }
}
