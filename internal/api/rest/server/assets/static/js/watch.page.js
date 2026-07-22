import { DOM_ELEMENTS, initDomElements } from "./watch.dom.js";
import * as browser from './browser.js';
import * as actionButton from './action-buttons.js';
import * as notify from './notifications.js';

const isPWA =
      window.matchMedia('(display-mode: standalone)').matches ||
      window.matchMedia('(display-mode: fullscreen)').matches ||
      window.navigator.standalone === true;

document.addEventListener('DOMContentLoaded', () => {
    initDomElements();

    document.addEventListener('keydown', function(e) {
        if (e.code === "Space") {
            e.preventDefault();
            togglePlay();
        }
        if (e.code === "ArrowRight") skipForward();
        if (e.code === "ArrowLeft") skipBackward();
    });

    // Initialize viewport height sync (fixes mobile PWA viewport issues)
    browser.initViewportHeightVar();

    actionButton.initCopyUrlButtons(
        (url) => {
            notify.show(`Link copied: ${url}`, notify.notifyType.SUCCESS);
        },
        (url) => {
            notify.show(`Failed to copy link: ${url}`, notify.notifyType.ERROR);
        }
    );

    initPlayer();
});

function initPlayer() {
    if (DOM_ELEMENTS.backButton) {
        const display = DOM_ELEMENTS.backButton.style.display;
        DOM_ELEMENTS.backButton.style.display = isPWA ? display : 'none';
        isPWA && (DOM_ELEMENTS.backButton.addEventListener('click', goBack));
    }

    // Update progress bar as the media plays
    if (DOM_ELEMENTS.customControls) {
        DOM_ELEMENTS.playButton.addEventListener('click', togglePlay);
        DOM_ELEMENTS.skipBackwardButton.addEventListener('click', skipBackward);
        DOM_ELEMENTS.skipForwardButton.addEventListener('click', skipForward);
        DOM_ELEMENTS.toggleFullscreenButton.addEventListener('click', toggleFullscreen);

        DOM_ELEMENTS.progressBar.addEventListener('click', function(e) {
            const rect = this.getBoundingClientRect();
            const clickPosition = (e.clientX - rect.left) / rect.width;
            DOM_ELEMENTS.player.currentTime = clickPosition * DOM_ELEMENTS.player.duration;
        });

        DOM_ELEMENTS.volumeControl.addEventListener('input', function() {
            DOM_ELEMENTS.player.volume = this.value;
        });
    }

    // Set video resize on metadata load and window resize
    if (DOM_ELEMENTS.video) {
        const video = DOM_ELEMENTS.video;

        video.addEventListener("loadedmetadata", fitVideo);
        window.addEventListener("resize", fitVideo);

        if (video.readyState >= 1) {
            fitVideo();
        }
    }

    // Auto play
    DOM_ELEMENTS.player.addEventListener('loadedmetadata', () => {
        DOM_ELEMENTS.player.autoplay = true;
        DOM_ELEMENTS.player.play().catch(err => {
            console.warn('Autorun failed', err);
        });
    });    
}

// Adjust video size to fit within max dimensions while maintaining aspect ratio
// Also adds "is-full" class to wrapper if video is already full width (e.g. on mobile)
function fitVideo() {
    const video = DOM_ELEMENTS.video;
    const mediaInfo = DOM_ELEMENTS.mediaInfo;

    // Check if video wrapper is already full width (e.g. on mobile)
    const rect = DOM_ELEMENTS.videoWrapper.getBoundingClientRect();
    const isFullWidth = Math.abs(rect.width - document.documentElement.clientWidth) < 2;
    DOM_ELEMENTS.videoWrapper.classList.toggle("is-full", isFullWidth); 

    // If video has no intrinsic size or max dimensions are not set,
    // remove explicit sizing to let it fill the wrapper
    const styles = getComputedStyle(video);
    if (!video.videoWidth || !video.videoHeight || styles.maxWidth === "none"  || styles.maxHeight === "none" ) {
        video.style.removeProperty("width");
        video.style.removeProperty("height");
        mediaInfo.style.removeProperty("width");
        return;
    }

    const maxW = parseFloat(styles.maxWidth);
    const maxH = parseFloat(styles.maxHeight);

    const vw = video.videoWidth;
    const vh = video.videoHeight;

    const scale = Math.min(maxW / vw, maxH / vh);
    const width = vw * scale;
    const maxWidth = 1000;

    video.style.width = `${width}px`;
    video.style.height = `${vh * scale}px`;

    if (width > maxWidth) {
        mediaInfo.style.removeProperty("max-width");
        mediaInfo.style.width = `${width}px`;
    } else {
        mediaInfo.style.removeProperty("width");
        mediaInfo.style.maxWidth = `${maxWidth}px`;
    }
    
}

function formatTime(seconds) {
    if (!seconds || isNaN(seconds)) return "0:00";
    const min = Math.floor(seconds / 60);
    const sec = Math.floor(seconds % 60);
    return `${min}:${sec < 10 ? '0' : ''}${sec}`;
}

function togglePlay() {
    if (DOM_ELEMENTS.player.paused) {
        DOM_ELEMENTS.player.play();
        if (DOM_ELEMENTS.playButton) {
            DOM_ELEMENTS.playButton.innerHTML = '❚❚';
        }
    } else {
        DOM_ELEMENTS.player.pause();
        if (DOM_ELEMENTS.playButton) {
            DOM_ELEMENTS.playButton.innerHTML = '▶';
        }
    }
}

function updateProgressBar() {
    const progress = document.getElementById('progress');
    const currentTimeEl = document.getElementById('currentTime');
    
    const percentage = (DOM_ELEMENTS.player.currentTime / DOM_ELEMENTS.player.duration) * 100;
    progress && (progress.style.width = percentage + '%');
    currentTimeEl && (currentTimeEl.textContent = formatTime(DOM_ELEMENTS.player.currentTime));
}

function skipForward() {
    DOM_ELEMENTS.player.currentTime += 10;
}

function skipBackward() {
    DOM_ELEMENTS.player.currentTime -= 10;
}

function toggleFullscreen() {
    const container = document.getElementById('mediaContainer');
    if (!document.fullscreenElement) {
        container.requestFullscreen();
    } else {
        document.exitFullscreen();
    }
}

function goBack() {
    window.history.back();
}
