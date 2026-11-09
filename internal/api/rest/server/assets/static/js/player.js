// -------------------------------------------------------------
// Player Overlay Logic
// Handles video/audio playback in modal overlay or bottom bar
// -------------------------------------------------------------

import { CLASS_NAMES as CONST_CLASS_NAMES, MEDIA_WATCH, EVENT_NAMES } from './constants.js';
import * as watchAPI from './watch-api.js';
import * as playerUtils from './player-utils.js';

const PLAYER_HASH = "#player"

const CLASS_NAMES = {
    ...CONST_CLASS_NAMES,

    mediaResultRow: "media-result__row",
    rowPlayButton: "media-play-button",

    isPlayerAttached: "is-player-attached",

    mediaPlayerWrapper: "media-player__wrapper",
    mediaPlayerCloseButton: "media-player__control--close",
    mediaPlayerMinimizeButton: "media-player__control--minimize",

    showControls: "media-player__video--show-controls",
    videoPlayerMinimized: "media-player__video--minimized",

    videoPlayer: "media-player__video-player",
    audioPlayer: "media-player__audio-player",

    playingDuration: "playing-overlay__duration-value",
    playingProgress: "playing-overlay__progress-value",
};

const CLASS_SELECTORS = Object.fromEntries(
    Object.entries(CLASS_NAMES)
        .filter(([, value]) => typeof value === "string")
        .map(([key, value]) => [key, `.${value}`])
);

const CSS_VARIABLE_NAMES = {
    watchBuffer: "--watch-buffer",
    watchProgress: "--watch-progress"
};

const DOM_IDS = {
    videoPlayerContainer: "media-player-video",
    audioPlayerContainer: "media-player-audio",

    videoPlayingOverlay: "media-player-video-playing-overlay",
    audioPlayingOverlay: "media-player-audio-playing-overlay",
};

const DOM_ELEMENTS = {
    videoPlayerContainer: null,
    audioPlayerContainer: null,

    videoPlayer: null,
    audioPlayer: null,

    videoPlayerMinimizeButton: null,
    videoPlayerCloseButton: null,
    audioPlayerCloseButton: null,

    videoPlayingOverlay: null,
    videoPlayingDuration: null,
    videoProgressValue: null,

    audioPlayingOverlay: null,
    audioPlayingDuration: null,
    audioProgressValue: null,
}  

function initDomElements(elements = DOM_ELEMENTS) {    
    elements.videoPlayerContainer = document.getElementById(DOM_IDS.videoPlayerContainer);
    elements.audioPlayerContainer = document.getElementById(DOM_IDS.audioPlayerContainer);

    elements.videoPlayingOverlay = document.getElementById(DOM_IDS.videoPlayingOverlay);
    elements.audioPlayingOverlay = document.getElementById(DOM_IDS.audioPlayingOverlay);

    if (elements.videoPlayerContainer) {
        elements.videoPlayer = elements.videoPlayerContainer.querySelector(CLASS_SELECTORS.videoPlayer);
        elements.videoPlayerMinimizeButton = elements.videoPlayerContainer.querySelector(CLASS_SELECTORS.mediaPlayerMinimizeButton)
        elements.videoPlayerCloseButton = elements.videoPlayerContainer.querySelector(CLASS_SELECTORS.mediaPlayerCloseButton)
    }

    if (elements.audioPlayerContainer) {
        elements.audioPlayer = elements.audioPlayerContainer.querySelector(CLASS_SELECTORS.audioPlayer);
        elements.audioPlayerCloseButton = elements.audioPlayerContainer.querySelector(CLASS_SELECTORS.mediaPlayerCloseButton)
    }

    if (elements.videoPlayingOverlay) {
        elements.videoPlayingDuration = elements.videoPlayingOverlay.querySelector(CLASS_SELECTORS.playingDuration)
        elements.videoProgressValue = elements.videoPlayingOverlay.querySelector(CLASS_SELECTORS.playingProgress)
    }

    if (elements.audioPlayingOverlay) {
        elements.audioPlayingDuration = elements.audioPlayingOverlay.querySelector(CLASS_SELECTORS.playingDuration)
        elements.audioProgressValue = elements.audioPlayingOverlay.querySelector(CLASS_SELECTORS.playingProgress)
    }
}

let watchTracker = null;

const mediaPlayer = {
    isAudio: false,
    isOpen: false,

    player: null,

    playingRow: null,
    rowItemId: "",

    unload() {
        this.player?.removeAttribute("src");
        this.player?.load();
    },

    clear() {
        this.isAudio = false;
        this.isOpen = false;

        if (this.player) {
            if (this.player.paused) {
                this.unload();
            } else {
                this.player.addEventListener("pause", () => this.unload(), { once: true });
                this.player.pause();
            }    

            this.player = null;
        }

        this.playingRow = null;
        this.rowItemId = "";

        this.syncWithHash();
    },

    isVideoPlayer() {
        return this.isOpen && !this.isAudio;
    },

    isAudioPlayer() {
        return this.isOpen && this.isAudio;
    },

    isMinimizedVideoPlayer() {
        return this.isOpen && !this.isAudio && DOM_ELEMENTS.videoPlayerContainer.classList.contains(CLASS_NAMES.videoPlayerMinimized);
    },

    isFullVideoPlayer() {
        return this.isOpen && !this.isAudio && !this.isMinimizedVideoPlayer();
    },

    updateBodyState() {
        document.body.classList.toggle(CLASS_NAMES.body.playerPaused, this.player.paused);
        document.body.classList.toggle(CLASS_NAMES.body.playerSoundOn, !this.player.paused && !this.player.muted);
    },

    updateBlockingState() {
        document.documentElement.classList.toggle(CLASS_NAMES.ui.blockingActive, !this.isMinimizedVideoPlayer());
    },

    toggleMinimizeVideoPlayer() {
        DOM_ELEMENTS.videoPlayerContainer.classList.toggle(CLASS_NAMES.videoPlayerMinimized);

        const label = this.isMinimizedVideoPlayer() ? "Restore player" : "Minimize player";

        DOM_ELEMENTS.videoPlayerMinimizeButton.title = label;
        DOM_ELEMENTS.videoPlayerMinimizeButton.setAttribute("aria-label", label);

        this.updateBlockingState();
        this.syncWithHash();
    },

    // Returns true if the current URL hash matches the player hash.
    isHash() {
        return window.location.hash === PLAYER_HASH;
    },

    // Removes the hash from the URL if present, without adding a history entry.
    clearHash() {
        if (!this.isHash()) return;

        history.replaceState(
            history.state,
            "",
            window.location.pathname + window.location.search
        );
    },

    // Player Overlay Functions
    syncWithHash() {
        if (!this.isOpen) {
            this.clearHash();
            return;
        }

        if (this.isAudio) {
            this.clearHash();
            return;
        }

        if (this.isMinimizedVideoPlayer()) {
            this.clearHash();
            return;
        }
        
        history.pushState({ player: true }, "", PLAYER_HASH);
    },
};

function initWatchTracker(player, itemId) {
    if (!player) return;
    if (watchTracker !== null) return;

    watchTracker = new watchAPI.MediaWatchTracker(player, itemId);
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

    mediaPlayer.syncWithHash();

    let cleanupHandlers = null;

    DOM_ELEMENTS.videoPlayerCloseButton.addEventListener("click", closePlayer);
    DOM_ELEMENTS.videoPlayerMinimizeButton.addEventListener("click", () => mediaPlayer.toggleMinimizeVideoPlayer());
    DOM_ELEMENTS.audioPlayerCloseButton.addEventListener("click", closePlayer);

    document.addEventListener("click", async (event) => {
        const playBtn = event.target.closest(CLASS_SELECTORS.rowPlayButton);
        if (!playBtn) return;

        const row = playBtn.closest(CLASS_SELECTORS.mediaResultRow);
        if (!row) return;

        openPlayer(row);
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
        if (event.key === "Escape" &&  mediaPlayer.isFullVideoPlayer()) closePlayer();
    });

    // Close the full video player when its hash is removed or changed.
    window.addEventListener("hashchange", () => {
        if (mediaPlayer.isFullVideoPlayer() && !mediaPlayer.isHash()) {
            closePlayer();
        }
    });

    async function openPlayer(row) {
        const itemId = row.dataset.itemId;
        const itemIsAudio = row.dataset.isAudio === "true";

        if (mediaPlayer.isOpen) {
            if (itemId === mediaPlayer.rowItemId) {
                mediaPlayer.player.paused ? mediaPlayer.player.play() : mediaPlayer.player.pause();
                return;
            }

            if (mediaPlayer.isAudio != itemIsAudio) {
                await closePlayer(false);
            } else {
                mediaPlayer.unload();
                destroyWatchTracker();
            }
        }

        if (mediaPlayer.playingRow) {
            mediaPlayer.playingRow.classList.remove(CLASS_NAMES.player.isAttached);
        }

        mediaPlayer.rowItemId = itemId;
        mediaPlayer.playingRow = row;

        mediaPlayer.playingRow.classList.add(CLASS_NAMES.player.isAttached);

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

        mediaPlayer.syncWithHash();

        if (mediaPlayer.isAudio) {
            cleanupHandlers = initAudioHandlers();

            mediaPlayer.player.focus({ preventScroll: true });

            document.body.style.overflow = "";

            const overlayTarget = mediaPlayer.playingRow.querySelector(`.${overlayTargetClassName}`);
            if (overlayTarget) {
                showAudioPlayingOverlay(overlayTarget);
            }

            requestAnimationFrame(() => {
                document.body.classList.add(CLASS_NAMES.body.mediaPlayer);
                document.body.classList.add(CLASS_NAMES.body.mediaPlayerAudio);
            });
        } else {
            mediaPlayer.updateBlockingState();
            
            cleanupHandlers = initVideoHandlers();

            const overlayTarget = mediaPlayer.playingRow.querySelector(`.${overlayTargetClassName}`);
            if (overlayTarget) {
                showVideoPlayingOverlay(overlayTarget);
            }

            requestAnimationFrame(() => {
                document.body.classList.add(CLASS_NAMES.body.mediaPlayer);
                document.body.classList.add(CLASS_NAMES.body.mediaPlayerVideo);
            });

            mediaPlayer.player.focus({ preventScroll: true });
        }

        if (itemId) {
            initWatchTracker(mediaPlayer.player, itemId)
        }
    }

    async function closePlayer(fullReset = false) {
        if (mediaPlayer.player !== null) {
            if (!mediaPlayer.player.paused) {
                await mediaPlayer.player.pause();
            }
        }

        history.replaceState(null, "", location.pathname + location.search);

        if (mediaPlayer.playingRow) {
            mediaPlayer.playingRow.classList.remove(CLASS_NAMES.player.isAttached);
        }

        cleanupHandlers?.();
        cleanupHandlers = null;

        mediaPlayer.clear();

        document.body.style.overflow = "";

        document.body.classList.remove(CLASS_NAMES.body.mediaPlayer);
        document.body.classList.remove(CLASS_NAMES.body.mediaPlayerVideo);
        document.body.classList.remove(CLASS_NAMES.body.mediaPlayerAudio);
        document.body.classList.remove(CLASS_NAMES.body.playerPaused);
        document.body.classList.remove(CLASS_NAMES.body.playerSoundOn);

        document.documentElement.classList.remove(CLASS_NAMES.ui.blockingActive);

        if (fullReset) {
            DOM_ELEMENTS.videoPlayerContainer.classList.remove(CLASS_NAMES.videoPlayerMinimized);
        }

        destroyWatchTracker();
    }
}

function initVideoHandlers() {
    if (!DOM_ELEMENTS.videoPlayerContainer) return;

    const videoContainer = DOM_ELEMENTS.videoPlayerContainer;

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

    const handlePlay = () => {
        showControls();
        document.body.classList.remove(CLASS_NAMES.body.playerPaused);
        mediaPlayer.updateBodyState();
    };

    const handlePause = () => {
        showControls();
        mediaPlayer.updateBodyState();
    };

    const handleVolumeChange = () => {
        showControls();
        mediaPlayer.updateBodyState();
    };

    const handleTimeUpdate = () => {
        playerUtils.updateRemainingDuration(DOM_ELEMENTS.videoPlayingDuration, mediaPlayer.player)

        playerUtils.updateProgress(
            null,
            DOM_ELEMENTS.videoProgressValue,
            mediaPlayer.player,
            {
                watchBuffer: CSS_VARIABLE_NAMES.watchBuffer,
                watchProgress: CSS_VARIABLE_NAMES.watchProgress,
            }
        )
    }

    videoContainer.addEventListener("mouseenter", showControls);
    mediaPlayer.player.addEventListener("mousemove", showControls);
    mediaPlayer.player.addEventListener("play", handlePlay);
    mediaPlayer.player.addEventListener("pause", handlePause);
    mediaPlayer.player.addEventListener("click", showControls);
    mediaPlayer.player.addEventListener("touchstart", showControls);
    mediaPlayer.player.addEventListener("volumechange", handleVolumeChange);
    mediaPlayer.player.addEventListener("timeupdate", handleTimeUpdate);

    showControls();
    
    return () => {
        videoContainer.removeEventListener("mouseenter", showControls);
        mediaPlayer.player.removeEventListener("mousemove", showControls);
        mediaPlayer.player.removeEventListener("play", handlePlay);
        mediaPlayer.player.removeEventListener("pause", handlePause);
        mediaPlayer.player.removeEventListener("click", showControls);
        mediaPlayer.player.removeEventListener("touchstart", showControls);
        mediaPlayer.player.removeEventListener("volumechange", handleVolumeChange);
        mediaPlayer.player.removeEventListener("timeupdate", handleTimeUpdate);

        clearTimeout(controlsTimeout);

        videoContainer.classList.remove(CLASS_NAMES.showControls);
    };
}

function initAudioHandlers() {
    const handlePlay = () => {
        mediaPlayer.updateBodyState();
    };

    const handlePause = () => {
        mediaPlayer.updateBodyState();
    };

    const handleVolumeChange = () => {
        mediaPlayer.updateBodyState();
    };

    const handleTimeUpdate = () => {
        playerUtils.updateRemainingDuration(DOM_ELEMENTS.audioPlayingDuration, mediaPlayer.player)

        playerUtils.updateProgress(
            null,
            DOM_ELEMENTS.audioProgressValue,
            mediaPlayer.player,
            {
                watchBuffer: CSS_VARIABLE_NAMES.watchBuffer,
                watchProgress: CSS_VARIABLE_NAMES.watchProgress,
            }
        )
    }

    mediaPlayer.player.addEventListener("play", handlePlay);
    mediaPlayer.player.addEventListener("pause", handlePause);
    mediaPlayer.player.addEventListener("volumechange", handleVolumeChange);
    mediaPlayer.player.addEventListener("timeupdate", handleTimeUpdate);

    return () => {
        mediaPlayer.player.removeEventListener("play", handlePlay);
        mediaPlayer.player.removeEventListener("pause", handlePause);
        mediaPlayer.player.removeEventListener("volumechange", handleVolumeChange);
        mediaPlayer.player.removeEventListener("timeupdate", handleTimeUpdate);
    };
}

function showAudioPlayingOverlay(target) {
    if (!target || !DOM_ELEMENTS.audioPlayingOverlay) return;

    target.append(DOM_ELEMENTS.audioPlayingOverlay);
}

function showVideoPlayingOverlay(target) {
    if (!target || !DOM_ELEMENTS.videoPlayingOverlay) return;

    target.append(DOM_ELEMENTS.videoPlayingOverlay);
}
