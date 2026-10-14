// -------------------------------------------------------------
// Player Overlay Logic
// Handles video/audio playback in modal overlay or bottom bar
// -------------------------------------------------------------

import * as watchAPI from './watch-api.js';
import { CLASS_NAMES as CONST_CLASS_NAMES, MEDIA_WATCH, EVENT_NAMES } from './constants.js';

const CLASS_NAMES = {
    ...CONST_CLASS_NAMES,

    mediaResultRow: "media-result__row",
    rowPlayButton: "media-play-button",

    mediaPlayerWrapper: "media-player__wrapper",
    mediaPlayerClose: "media-player__close",

    showControls: "show-controls",

    videoPlayer: "media-player__video-player",
    audioPlayer: "media-player__audio-player",

    audioPlayingOverlay: 'media-player__audio-playing-overlay',
};

const CLASS_SELECTORS = Object.fromEntries(
    Object.entries(CLASS_NAMES)
        .filter(([, value]) => typeof value === "string")
        .map(([key, value]) => [key, `.${value}`])
);

const DOM_IDS = {
    audioPlayingOverlay: "media-player-audio-playing-overlay",
    videoPlayerContainer: "media-player-video",
    audioPlayerContainer: "media-player-audio",
};

const DOM_ELEMENTS = {
    audioPlayingOverlay: null,
    videoPlayerContainer: null,
    audioPlayerContainer: null,
    videoPlayer: null,
    audioPlayer: null,
}  

function initDomElements(elements = DOM_ELEMENTS) {    
    elements.audioPlayingOverlay = document.getElementById(DOM_IDS.audioPlayingOverlay);
    elements.videoPlayerContainer = document.getElementById(DOM_IDS.videoPlayerContainer);
    elements.audioPlayerContainer = document.getElementById(DOM_IDS.audioPlayerContainer);

    if (elements.videoPlayerContainer) {
        elements.videoPlayer = elements.videoPlayerContainer.querySelector(CLASS_SELECTORS.videoPlayer);
    }

    if (elements.audioPlayerContainer) {
        elements.audioPlayer = elements.audioPlayerContainer.querySelector(CLASS_SELECTORS.audioPlayer);
    }
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
            if (this.player.paused) {
                unloadMediaPlayer(this.player);
            } else {
                this.player.addEventListener("pause", () => unloadMediaPlayer(this.player), { once: true });
                this.player.pause();
            }    

            this.player = null;
        }

        this.playingRow = null;
        this.rowItemId = "";
    },
};

function unloadMediaPlayer(player) {
    player.removeAttribute("src");
    player.load();
}

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

    const videoContainer = DOM_ELEMENTS.videoPlayerContainer;
    const audioContainer = DOM_ELEMENTS.audioPlayerContainer;
    const videoPlayer = DOM_ELEMENTS.videoPlayer;
    const audioPlayer = DOM_ELEMENTS.audioPlayer;

    if (!videoContainer || !audioContainer || !videoPlayer || !audioPlayer) return;

    const videoPlayerHash = "#player"
    let cleanupControls = null;

    initWatchTracker();

    videoContainer.querySelector(CLASS_SELECTORS.mediaPlayerClose).addEventListener("click", closePlayer);
    audioContainer.querySelector(CLASS_SELECTORS.mediaPlayerClose).addEventListener("click", closePlayer);

    document.addEventListener("click", async (event) => {
        const playBtn = event.target.closest(CLASS_SELECTORS.rowPlayButton);
        if (!playBtn) return;

        const row = playBtn.closest(CLASS_SELECTORS.mediaResultRow);
        if (!row) return;

        const itemId = row.dataset.itemId;
        const itemIsAudio = row.dataset.isAudio === "true";

        if (mediaPlayer.isOpen) {
            if (mediaPlayer.isAudio && itemId === mediaPlayer.rowItemId) {
                mediaPlayer.player.paused ? mediaPlayer.player.play() : mediaPlayer.player.pause();
                return;
            }

            if (mediaPlayer.isAudio != itemIsAudio) {
                await closePlayer();
            } else {
                unloadMediaPlayer(mediaPlayer.player);
            }
        }

        mediaPlayer.rowItemId = itemId;
        mediaPlayer.playingRow = row;

        const mediaURL = mediaPlayer.playingRow.dataset.media;
        if (!mediaURL) return;

        mediaPlayer.isAudio = itemIsAudio;

        if (!mediaPlayer.isAudio) {
            mediaPlayer.player = videoPlayer;
            document.dispatchEvent(new Event(EVENT_NAMES.videoPlayerOpened));
        } else {
            mediaPlayer.player = audioPlayer;
        }

        const shouldLoop = mediaPlayer.playingRow.dataset.loop === "true";

        let positionMs = await watchAPI.getWatchPosition(itemId);

        if (positionMs < MEDIA_WATCH.startThresholdMs) {
            positionMs = 0;
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
            cleanupControls = initAudioControls();

            mediaPlayer.player.focus({ preventScroll: true });

            document.body.style.overflow = "";

            const overlayTarget = mediaPlayer.playingRow.querySelector(`.${overlayTargetClassName}`);
            if (overlayTarget) {
                showAudioPlayingOverlay(overlayTarget);
            }

            requestAnimationFrame(() => {
                document.body.classList.add(CLASS_NAMES.body.audioPlayer);
            });
        } else {
            document.documentElement.classList.add(CLASS_NAMES.ui.blockingActive);
            
            location.hash = videoPlayerHash

            cleanupControls = initVideoControls();

            requestAnimationFrame(() => {
                document.body.classList.add(CLASS_NAMES.body.videoPlayer);
            });

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

        const playBtn = el.closest(CLASS_SELECTORS.rowPlayButton);
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

        document.body.style.overflow = "";

        document.body.classList.remove(CLASS_NAMES.body.videoPlayer);
        document.body.classList.remove(CLASS_NAMES.body.audioPlayer);
        document.body.classList.remove(CLASS_NAMES.body.playerPaused);

        document.documentElement.classList.remove(CLASS_NAMES.ui.blockingActive);

        destroyWatchTracker();
    }

    function initAudioControls() {
        const handlePlay = () => {
            document.body.classList.remove(CLASS_NAMES.body.playerPaused);
        };

        const handlePause = () => {
            document.body.classList.add(CLASS_NAMES.body.playerPaused);
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
