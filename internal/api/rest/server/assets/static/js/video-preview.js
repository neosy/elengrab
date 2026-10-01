import * as watchAPI from './watch-api.js';
import { CLASS_NAMES as CONST_CLASS_NAMES, MEDIA_WATCH, VIDEO_PREVIEW } from './constants.js';
import { isMobileScreen } from './browser.js';

const CLASS_NAMES = {
    ...CONST_CLASS_NAMES,
    soundOff: "video-preview__sound-off",
    soundOn: "video-preview__sound-on",
    previewPlaying: VIDEO_PREVIEW.previewPlayingClassName,

    rowRefreshing: "row--refreshing",
};

export const CLASS_SELECTORS = Object.fromEntries(
    Object.entries(CLASS_NAMES)
        .filter(([, value]) => typeof value === "string")
        .map(([key, value]) => [key, `.${value}`])
);

const CSS_VARIABLE_NAMES = {
    watchBuffer: "--video-preview-watch-buffer",
    watchProgress: "--video-preview-watch-progress"
};

const DOM_IDS = {
    row: (id) => `row-${id}`,
}

const DOM_ELEMENTS = {
    sound: {
        button: null,
        iconOff: null,
        iconOn: null,
    },
    preview: {
        container: null,
        player: null,
        durationRemaining: null,
        progressValue: null,
        progressBuffer: null,

        isAvailable() {
            return Boolean(this.container && this.player);
        },
    },
}

const PREVIEW_STATE = {
    hoverTimer: null,
    scrollTimer: null,

    ended: false,

    currentItemId: null,
    currentVideoUrl: null,

    updateCenteredRequestId: 0,
    showVideoPreviewRequestId: 0,

    watchTracker: null,

    clear() {
        clearTimeout(this.hoverTimer);
        this.currentItemId = null;
        this.currentVideoUrl = null;
    },
};

function initDomElements() {
    DOM_ELEMENTS.preview.container = document.getElementById("video-preview-container");
    DOM_ELEMENTS.preview.player = document.getElementById("video-preview-player");
    DOM_ELEMENTS.preview.durationRemaining = document.getElementById("video-preview-duration-remaining");
    DOM_ELEMENTS.preview.progressValue = document.getElementById("video-preview-watch-progress-value");
    DOM_ELEMENTS.preview.progressBuffer = document.getElementById("video-preview-watch-progress-buffer");

    const soundButton = document.getElementById("video-preview-sound");
    if (soundButton != null) {
        DOM_ELEMENTS.sound.button = soundButton;
        DOM_ELEMENTS.sound.iconOff = soundButton.querySelector(CLASS_SELECTORS.soundOff);
        DOM_ELEMENTS.sound.iconOn = soundButton.querySelector(CLASS_SELECTORS.soundOn);
    }
}

export function initVideoPreview() {
    initDomElements();

    if (!DOM_ELEMENTS.preview.isAvailable()) {
        return;
    }

    DOM_ELEMENTS.preview.container.hidden = true;

    DOM_ELEMENTS.preview.player.muted = true;
    DOM_ELEMENTS.preview.player.playsInline = true;

    initWatchTracker(DOM_ELEMENTS.preview.player);

    if (DOM_ELEMENTS.sound.button !== null) {
        DOM_ELEMENTS.sound.button.addEventListener("click", (event) => {
            event.stopPropagation();
            event.preventDefault();

            DOM_ELEMENTS.preview.player.muted = !DOM_ELEMENTS.preview.player.muted;

            DOM_ELEMENTS.sound.iconOff.hidden = !DOM_ELEMENTS.preview.player.muted;
            DOM_ELEMENTS.sound.iconOn.hidden = DOM_ELEMENTS.preview.player.muted;

            const title = DOM_ELEMENTS.preview.player.muted ? "Turn on the sound" : "Turn off the sound";
            DOM_ELEMENTS.sound.button.setAttribute("title", title);
            DOM_ELEMENTS.sound.button.setAttribute("aria-label", title);
        });    
    }
}

function toggleVideoPreviewSound() {
    DOM_ELEMENTS.preview.player.muted = !DOM_ELEMENTS.preview.player.muted;

    DOM_ELEMENTS.sound.iconOff.hidden = !DOM_ELEMENTS.preview.player.muted;
    DOM_ELEMENTS.sound.iconOn.hidden = DOM_ELEMENTS.preview.player.muted;

    const title = DOM_ELEMENTS.preview.player.muted
        ? "Turn on the sound"
        : "Turn off the sound";

    DOM_ELEMENTS.sound.button.setAttribute("title", title);
    DOM_ELEMENTS.sound.button.setAttribute("aria-label", title);
}

function initWatchTracker(video) {
    if (!video) {
        return;
    }

    PREVIEW_STATE.watchTracker = new watchAPI.MediaWatchTracker(video);
    PREVIEW_STATE.watchTracker.init();        
}

function setWatchTrackerItemId(itemId) {
    if (!PREVIEW_STATE.watchTracker) {
        return;
    }

    PREVIEW_STATE.watchTracker.setItemId(itemId);
}

export function initVideoPreviewHover(previewAreaContainer, previewContainerClassName, thumbClassName) {
    if (!previewAreaContainer) {
        return;
    }

    previewAreaContainer.addEventListener("mouseover", async (event) => {
        const previewContainer = event.target.closest(`.${previewContainerClassName}`);
        const thumbnailElement = previewContainer?.querySelector(`.${thumbClassName}`);

        if (!thumbnailElement) {
            return;
        }

        if (!shouldShowVideoPreview(previewContainer)) {
            return;
        }

        if (!previewAreaContainer.contains(previewContainer)) {
            return;
        }

        if (event.relatedTarget && previewContainer.contains(event.relatedTarget)) {
            return;
        }

        clearTimeout(PREVIEW_STATE.hoverTimer);

        PREVIEW_STATE.hoverTimer = setTimeout(async () => {
            if (PREVIEW_STATE.ended) {
                return;
            }

            showVideoPreview(
                thumbnailElement,
                previewContainer,
            );
        }, 300);
    });

    previewAreaContainer.addEventListener("mouseout", (event) => {
        if (isMobileScreen()) return;

        const previewContainer = event.target.closest(`.${previewContainerClassName}`);

        if (!previewContainer || !previewAreaContainer.contains(previewContainer)) {
            return;
        }

        if (event.relatedTarget && previewContainer.contains(event.relatedTarget)) {
            return;
        }

        hideVideoPreview();
    });

    document.addEventListener(VIDEO_PREVIEW.playerOpenedEventName, () => {
        hideVideoPreview();
    });

    DOM_ELEMENTS.preview.player.addEventListener("ended", () => {
        PREVIEW_STATE.ended = true;
        hideVideoPreview();
    });

    DOM_ELEMENTS.preview.player.addEventListener("timeupdate", updateVideoPreviewDuration);
}

export function shouldShowVideoPreview(previewContainer) {
    if (isMobileScreen()) {
        return false;
    }

    if (!previewContainer) {
        return false;
    }

    if (document.body.classList.contains(CLASS_NAMES.listView)) {
        return false;
    }

    if (!previewContainer.classList.contains(CLASS_NAMES.rowStatus.success)) {
        return false;
    }

    if (previewContainer.dataset.isAudio === "true") {
        return false;
    }

    return Boolean(previewContainer.dataset.itemId && previewContainer.dataset.media);
}

export async function showVideoPreview(thumbnailElement, previewContainer) {
    const videoUrl = String(previewContainer.dataset.media)
    const itemId = String(previewContainer.dataset.itemId)
    const requestId = ++PREVIEW_STATE.showVideoPreviewRequestId;

    if (!DOM_ELEMENTS.preview.isAvailable()) {
        return;
    }

    const positionMs = await watchAPI.getWatchPosition(itemId);

    if (requestId !== PREVIEW_STATE.showVideoPreviewRequestId) {
        return
    }

    PREVIEW_STATE.ended = false;

    const itemEl = document.getElementById(DOM_IDS.row(itemId));
    if (itemEl) {
        const isPreviewBlocked = itemEl.classList.contains(CLASS_NAMES.rowRefreshing);
        if (isPreviewBlocked) return;

        itemEl.classList.add(CLASS_NAMES.previewPlaying);
    }

    thumbnailElement.appendChild(DOM_ELEMENTS.preview.container);

    setWatchTrackerItemId(itemId);

    if (PREVIEW_STATE.currentVideoUrl !== videoUrl) {
        PREVIEW_STATE.currentVideoUrl = videoUrl;
        DOM_ELEMENTS.preview.player.src = videoUrl;

        await new Promise(resolve => {
            DOM_ELEMENTS.preview.player.onloadedmetadata = resolve;
        });
    }

    if (requestId !== PREVIEW_STATE.showVideoPreviewRequestId) {
        return
    }

    if (positionMs < MEDIA_WATCH.startThresholdMs) {
        DOM_ELEMENTS.preview.player.currentTime = 0;
    } else {
        DOM_ELEMENTS.preview.player.currentTime = positionMs / 1000;
    }

    DOM_ELEMENTS.preview.player.playsInline = true;
    DOM_ELEMENTS.preview.player.loop = false;

    PREVIEW_STATE.currentItemId = itemId;

    try {
        await DOM_ELEMENTS.preview.player.play();

        if (!PREVIEW_STATE.currentItemId) {
            hideVideoPreview();
            return
        }

        DOM_ELEMENTS.preview.container.hidden = false;
    } catch (error) {
        hideVideoPreview();
    }
}

export function hideVideoPreview() {
    if (!DOM_ELEMENTS.preview.isAvailable()) {
        return;
    }

    ++PREVIEW_STATE.showVideoPreviewRequestId;

    const itemEl = document.getElementById(DOM_IDS.row(PREVIEW_STATE.currentItemId));
    if (itemEl) {
        itemEl.classList.remove(CLASS_NAMES.previewPlaying);
    }

    PREVIEW_STATE.clear();

    DOM_ELEMENTS.preview.container.hidden = true;
    DOM_ELEMENTS.preview.player.pause();
}

function updateVideoPreviewDuration() {
    if (!DOM_ELEMENTS.preview.player || !DOM_ELEMENTS.preview.durationRemaining) {
        return;
    }

    const remainingSeconds = Math.max(
        0,
        Math.floor(DOM_ELEMENTS.preview.player.duration - DOM_ELEMENTS.preview.player.currentTime)
    );

    DOM_ELEMENTS.preview.durationRemaining.textContent = formatDuration(remainingSeconds);

    if (DOM_ELEMENTS.preview.progressBuffer !== null) {
        const bufferPercent = (DOM_ELEMENTS.preview.player.duration && DOM_ELEMENTS.preview.player.buffered.length > 0)
        ? Math.floor((DOM_ELEMENTS.preview.player.buffered.end(DOM_ELEMENTS.preview.player.buffered.length - 1) / DOM_ELEMENTS.preview.player.duration) * 100)
        : 0;

        DOM_ELEMENTS.preview.progressBuffer.style.setProperty(
            CSS_VARIABLE_NAMES.watchBuffer,
            `${bufferPercent}%`
        );
    }

    if (DOM_ELEMENTS.preview.progressValue !== null) {
        const progressPercent = DOM_ELEMENTS.preview.player.duration > 0
        ? Math.floor((DOM_ELEMENTS.preview.player.currentTime / DOM_ELEMENTS.preview.player.duration) * 1000) / 10
        : 0;

        DOM_ELEMENTS.preview.progressValue.style.setProperty(
            CSS_VARIABLE_NAMES.watchProgress,
            `${progressPercent}%`
        );
    }
}

function formatDuration(seconds) {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const sec = (seconds % 60).toString().padStart(2, "0");

    return days > 0 
        ? `${days}:${hours.toString().padStart(2, "0")}:${minutes.toString().padStart(2, "0")}:${sec}` 
        : hours > 0 
            ? `${hours}:${minutes.toString().padStart(2, "0")}:${sec}`
            : `${minutes}:${sec}`;
}

export function initVideoPreviewScroll(container, previewContainerClassName, thumbClassName) {
    if (!container) {
        return;
    }

    const refreshPreviewWithForce = (force=false) => {
        updateVideoPreview(container, previewContainerClassName, thumbClassName, force);
    };    

    const refreshPreview = () => {
        refreshPreviewWithForce(false);
    };    

    window.addEventListener("scroll", refreshPreview, { passive: true });
    window.addEventListener("resize", refreshPreview);

    // We launch it immediately after opening the page.
    refreshPreview();

    return refreshPreviewWithForce;
}

function updateVideoPreview(container, previewContainerClassName, thumbClassName, force=false) {
    if (!isMobileScreen()) {
        return;
    }

    if (PREVIEW_STATE.currentItemId) {
        const element = document.getElementById(DOM_IDS.row(PREVIEW_STATE.currentItemId));

        // Stop the current preview when less than 60% is visible.
        if (force || !element || !isElementVisibleEnough(element)) {
            hideVideoPreview();
        }
    }

    clearTimeout(PREVIEW_STATE.scrollTimer);

    PREVIEW_STATE.scrollTimer = setTimeout(() => {
        updateCenteredPreview(container, previewContainerClassName, thumbClassName);
    }, 120);
}

function isElementInViewport(element) {
    const rect = element.getBoundingClientRect();

    return rect.bottom > 0 && rect.top < window.innerHeight;
}

async function updateCenteredPreview(container, previewContainerClassName, thumbClassName) {
    const previewContainer = findCenteredElement(container, previewContainerClassName);

    if (!previewContainer) {
        hideVideoPreview();
        return;
    }

    if (!previewContainer.classList.contains(CLASS_NAMES.rowStatus.success)) {
        return;
    }

    const itemId = previewContainer.dataset.itemId;

    if (itemId === PREVIEW_STATE.currentItemId) {
        return;
    }

    const requestId = ++PREVIEW_STATE.updateCenteredRequestId;

    hideVideoPreview();

    // While waiting for a response, the user has already scrolled through the list.
    if (requestId !== PREVIEW_STATE.updateCenteredRequestId) {
        return;
    }

    const thumbnail = previewContainer.querySelector(`.${thumbClassName}`);
    if (!thumbnail) {
        return;
    }

    await showVideoPreview(
        thumbnail,
        previewContainer,
    );
}

function isElementVisibleEnough(element, threshold = 0.6) {
    const rect = element.getBoundingClientRect();

    // Completely off-screen.
    if (rect.bottom <= 0 || rect.top >= window.innerHeight) {
        return false;
    }

    const visibleHeight =
        Math.min(rect.bottom, window.innerHeight) -
        Math.max(rect.top, 0);

    // More than 60% visible.
    return visibleHeight >= rect.height * threshold;
}

function findCenteredElement(container, previewContainerClassName) {
    const viewportCenter = window.innerHeight / 2;

    let bestElement = null;
    let bestDistance = Number.MAX_VALUE;

    const items = container.querySelectorAll(`.${previewContainerClassName}`);

    for (const item of items) {
        if (item.dataset.isAudio === "true") {
            continue;
        }

        const rect = item.getBoundingClientRect();

        // Less than 60% visible.
        if (!isElementVisibleEnough(item, 0.6)) {
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