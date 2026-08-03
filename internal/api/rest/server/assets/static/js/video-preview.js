import * as watchAPI from './watch-api.js';
import { MEDIA_WATCH } from './constants.js';
import { isMobileScreen } from './browser.js';

let hoverTimer;
let watchTracker = null;

let videoPreviewEnded = false;

let soundElements = null;

let previewElements = null;

const cssClassNames = {
    soundOff: "video-preview__sound-off",
    soundOn: "video-preview__sound-on",
};

const cssVarNames = {
    watchBuffer: "--video-preview-watch-buffer",
    watchProgress: "--video-preview-watch-progress"
};

export function initVideoPreview() {
    previewElements = {
        container: document.getElementById("video-preview-container"),
        player: document.getElementById("video-preview-player"),
        durationRemaining: document.getElementById("video-preview-duration-remaining"),
        progressValue: document.getElementById("video-preview-watch-progress-value"),
        progressBuffer: document.getElementById("video-preview-watch-progress-buffer")
    };

    const soundButton = document.getElementById("video-preview-sound");
    if (soundButton != null) {
        soundElements = {
            button: soundButton,
            iconOff: soundButton.querySelector(`.${cssClassNames.soundOff}`),
            iconOn: soundButton.querySelector(`.${cssClassNames.soundOn}`)
        }
    }

    if (!previewElements.container || !previewElements.player) {
        return;
    }

    previewElements.container.hidden = true;

    previewElements.player.muted = true;
    previewElements.player.playsInline = true;

    initWatchTracker(previewElements.player);

    if (soundButton !== null) {
        soundButton.addEventListener("click", (event) => {
            event.stopPropagation();
            event.preventDefault();

            previewElements.player.muted = !previewElements.player.muted;

            soundElements.iconOff.hidden = !previewElements.player.muted;
            soundElements.iconOn.hidden = previewElements.player.muted;

            const title = previewElements.player.muted ? "Turn on the sound" : "Turn off the sound";
            soundElements.button.setAttribute("title", title);
            soundElements.button.setAttribute("aria-label", title);
        });    
    }
}

function toggleVideoPreviewSound() {
    previewElements.player.muted = !previewElements.player.muted;

    soundElements.iconOff.hidden = !previewElements.player.muted;
    soundElements.iconOn.hidden = previewElements.player.muted;

    const title = previewElements.player.muted
        ? "Turn on the sound"
        : "Turn off the sound";

    soundElements.button.setAttribute("title", title);
    soundElements.button.setAttribute("aria-label", title);
}

function initWatchTracker(video) {
    if (!video) {
        return;
    }

    watchTracker = new watchAPI.MediaWatchTracker(video);
    watchTracker.init();        
}

function setWatchTrackerItemId(itemId) {
    if (!watchTracker) {
        return;
    }

    watchTracker.setItemId(itemId);
}

export function initVideoPreviewHover(container, elementClassName, thumbClassName) {
    if (!container) {
        return;
    }

    container.addEventListener("mouseover", async (event) => {
        if (isMobileScreen()) return;

        const el = event.target.closest(`.${elementClassName}`);

        if (!el || !container.contains(el)) {
            return;
        }

        if (event.relatedTarget && el.contains(event.relatedTarget)) {
            return;
        }

        const thumbnail = el.querySelector(`.${thumbClassName}`);

        if (!thumbnail) {
            return;
        }

        const isAudio = el.dataset.isAudio === "true";
        if (isAudio) {
            return;
        }

        const itemId = el.dataset.itemId;

        clearTimeout(hoverTimer);

        hoverTimer = setTimeout(async () => {
            if (videoPreviewEnded) {
                return;
            }

            const positionMs = await watchAPI.getWatchPosition(itemId);

            showVideoPreview(
                thumbnail,
                el.dataset.media,
                itemId,
                positionMs
            );
        }, 300);
    });

    container.addEventListener("mouseout", (event) => {
        if (isMobileScreen()) return;

        const el = event.target.closest(`.${elementClassName}`);

        if (!el || !container.contains(el)) {
            return;
        }

        if (event.relatedTarget && el.contains(event.relatedTarget)) {
            return;
        }

        videoPreviewEnded = false;

        clearTimeout(hoverTimer);
        hideVideoPreview();
    });

    document.addEventListener("preview-player-opened", () => {
        clearTimeout(hoverTimer);
        hideVideoPreview();
    });

    previewElements.player.addEventListener("ended", () => {
        videoPreviewEnded = true;
        hideVideoPreview();
    });

    previewElements.player.addEventListener("timeupdate", updateVideoPreviewDuration);
}

export async function showVideoPreview(thumbnail, videoUrl, itemId, positionMs) {
    if (!previewElements.container || !previewElements.player) {
        return;
    }

    videoPreviewEnded = false;

    thumbnail.appendChild(previewElements.container);

    setWatchTrackerItemId(itemId);

    previewElements.player.src = videoUrl;

    await new Promise((resolve) => {
        previewElements.player.onloadedmetadata = resolve;
    });

    if (positionMs < MEDIA_WATCH.startThresholdMs) {
        previewElements.player.currentTime = 0;
    } else {
        previewElements.player.currentTime = positionMs / 1000;
    }

    previewElements.player.playsInline = true;
    previewElements.player.loop = false;

    try {
        await previewElements.player.play();
        previewElements.container.hidden = false;
    } catch (error) {
        hideVideoPreview();
        console.debug("Video preview play failed", error);
    }
}

export function hideVideoPreview() {
    if (!previewElements.container || !previewElements.player) {
        return;
    }

    previewElements.player.pause();

    previewElements.container.hidden = true;
}

function setVideoPreviewPosition(element) {
    const rect = element.getBoundingClientRect();

    previewElements.container.style.position = "fixed";

    previewElements.container.style.left = `${rect.left}px`;
    previewElements.container.style.top = `${rect.top}px`;

    previewElements.container.style.width = `${rect.width}px`;
    previewElements.container.style.height = `${rect.height}px`;
}

function updateVideoPreviewDuration() {
    if (!previewElements.player || !previewElements.durationRemaining) {
        return;
    }

    const remainingSeconds = Math.max(
        0,
        Math.floor(previewElements.player.duration - previewElements.player.currentTime)
    );

    previewElements.durationRemaining.textContent = formatDuration(remainingSeconds);

    if (previewElements.progressBuffer !== null) {
        const bufferPercent = (previewElements.player.duration && previewElements.player.buffered.length > 0)
        ? Math.floor((previewElements.player.buffered.end(previewElements.player.buffered.length - 1) / previewElements.player.duration) * 100)
        : 0;

        previewElements.progressBuffer.style.setProperty(
            cssVarNames.watchBuffer,
            `${bufferPercent}%`
        );
    }

    if (previewElements.progressValue !== null) {
        const progressPercent = previewElements.player.duration > 0
        ? Math.floor((previewElements.player.currentTime / previewElements.player.duration) * 1000) / 10
        : 0;

        previewElements.progressValue.style.setProperty(
            cssVarNames.watchProgress,
            `${progressPercent}%`
        );
    }
}

function formatDuration(seconds) {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;

    return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
}