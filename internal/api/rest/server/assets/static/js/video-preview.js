import * as watchAPI from './watch-api.js';
import { CLASS_NAMES, MEDIA_WATCH, VIDEO_PREVIEW, DOM_IDS } from './constants.js';
import { isMobileScreen } from './browser.js';

const soundElements = {
    button: null,
    iconOff: null,
    iconOn: null
};

const previewElements = {
    container: null,
    player: null,
    durationRemaining: null,
    progressValue: null,
    progressBuffer: null
};

const previewState = {
    hoverTimer: null,
    scrollTimer: null,

    ended: false,

    currentItemId: null,
    currentVideoUrl: null,

    requestId: 0,
};

let watchTracker = null;

const cssClassNames = {
    soundOff: "video-preview__sound-off",
    soundOn: "video-preview__sound-on",
    previewPlaying: VIDEO_PREVIEW.previewPlayingClassName,
};

const cssVarNames = {
    watchBuffer: "--video-preview-watch-buffer",
    watchProgress: "--video-preview-watch-progress"
};

export function initVideoPreview() {
    previewElements.container = document.getElementById("video-preview-container");
    previewElements.player = document.getElementById("video-preview-player");
    previewElements.durationRemaining = document.getElementById("video-preview-duration-remaining");
    previewElements.progressValue = document.getElementById("video-preview-watch-progress-value");
    previewElements.progressBuffer = document.getElementById("video-preview-watch-progress-buffer");

    const soundButton = document.getElementById("video-preview-sound");
    if (soundButton != null) {
        soundElements.button = soundButton;
        soundElements.iconOff = soundButton.querySelector(`.${cssClassNames.soundOff}`);
        soundElements.iconOn = soundButton.querySelector(`.${cssClassNames.soundOn}`);
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

        if (!el.classList.contains(CLASS_NAMES.rowStatus.success)) {
            return;
        }

        const isAudio = el.dataset.isAudio === "true";
        if (isAudio) {
            return;
        }

        const itemId = el.dataset.itemId;

        const thumbnail = el.querySelector(`.${thumbClassName}`);

        if (!thumbnail) {
            return;
        }

        clearTimeout(previewState.hoverTimer);

        previewState.hoverTimer = setTimeout(async () => {
            if (previewState.ended) {
                return;
            }

            showVideoPreview(
                thumbnail,
                el.dataset.media,
                itemId
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

        previewState.ended = false;

        clearTimeout(previewState.hoverTimer);
        hideVideoPreview();
    });

    document.addEventListener(VIDEO_PREVIEW.playerOpenedEventName, () => {
        clearTimeout(previewState.hoverTimer);
        hideVideoPreview();
    });

    previewElements.player.addEventListener("ended", () => {
        previewState.ended = true;
        hideVideoPreview();
    });

    previewElements.player.addEventListener("timeupdate", updateVideoPreviewDuration);
}

export async function showVideoPreview(thumbnail, videoUrl, itemId) {
    if (!previewElements.container || !previewElements.player) {
        return;
    }

    const positionMs = await watchAPI.getWatchPosition(itemId);

    previewState.ended = false;

    const itemEl = document.getElementById(DOM_IDS.row(itemId));
    if (itemEl) {
        itemEl.classList.add(cssClassNames.previewPlaying);
    }

    thumbnail.appendChild(previewElements.container);

    setWatchTrackerItemId(itemId);

    if (previewState.currentVideoUrl !== videoUrl) {
        previewState.currentVideoUrl = videoUrl;
        previewElements.player.src = videoUrl;

        await new Promise(resolve => {
            previewElements.player.onloadedmetadata = resolve;
        });
    }

    if (positionMs < MEDIA_WATCH.startThresholdMs) {
        previewElements.player.currentTime = 0;
    } else {
        previewElements.player.currentTime = positionMs / 1000;
    }

    previewElements.player.playsInline = true;
    previewElements.player.loop = false;

    previewState.currentItemId = itemId;

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

    const itemEl = document.getElementById(DOM_IDS.row(previewState.currentItemId));
    if (itemEl) {
        itemEl.classList.remove(cssClassNames.previewPlaying);
    }

    previewState.currentItemId = null;

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

export function initVideoPreviewScroll(container, elementClassName, thumbClassName) {
    if (!container) {
        return;
    }

    const onScroll = () => {
        if (!isMobileScreen()) {
            return;
        }

        clearTimeout(previewState.scrollTimer);

        previewState.scrollTimer = setTimeout(() => {
            updateCenteredPreview(container, elementClassName, thumbClassName);
        }, 120);
    };

    window.addEventListener("scroll", onScroll, { passive: true });
    window.addEventListener("resize", onScroll);

    // We launch it immediately after opening the page.
    onScroll();
}

async function updateCenteredPreview(container, elementClassName, thumbClassName) {
    const element = findCenteredElement(container, elementClassName);

    if (!element) {
        hideVideoPreview();
        return;
    }

    if (!element.classList.contains(CLASS_NAMES.rowStatus.success)) {
        return;
    }

    const itemId = element.dataset.itemId;

    if (itemId === previewState.currentItemId) {
        return;
    }

    const requestId = ++previewState.requestId;

    hideVideoPreview();

    // While waiting for a response, the user has already scrolled through the list.
    if (requestId !== previewState.requestId) {
        return;
    }

    const thumbnail = element.querySelector(`.${thumbClassName}`);
    if (!thumbnail) {
        return;
    }

    await showVideoPreview(
        thumbnail,
        element.dataset.media,
        itemId
    );
}

function findCenteredElement(container, elementClassName) {
    const viewportCenter = window.innerHeight / 2;

    let bestElement = null;
    let bestDistance = Number.MAX_VALUE;

    const items = container.querySelectorAll(`.${elementClassName}`);

    for (const item of items) {
        if (item.dataset.isAudio === "true") {
            continue;
        }

        const rect = item.getBoundingClientRect();

        // Completely off-screen.
        if (rect.bottom <= 0 || rect.top >= window.innerHeight) {
            continue;
        }

        // Less than 40% visible.
        const visibleHeight =
            Math.min(rect.bottom, window.innerHeight) -
            Math.max(rect.top, 0);

        if (visibleHeight < rect.height * 0.4) {
            continue;
        }

        const center = rect.top + rect.height / 2;
        const distance = Math.abs(center - viewportCenter);

        if (distance < bestDistance) {
            bestDistance = distance;
            bestElement = item;
        }
    }

    return bestElement;
}