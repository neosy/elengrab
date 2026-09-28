import * as watchAPI from './watch-api.js';
import { CLASS_NAMES as CONST_CLASS_NAMES, MEDIA_WATCH, VIDEO_PREVIEW } from './constants.js';
import { isMobileScreen } from './browser.js';

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
    },
}

const PREVIEW_STATE = {
    hoverTimer: null,
    scrollTimer: null,

    ended: false,

    currentItemId: null,
    currentVideoUrl: null,

    requestId: 0,

    watchTracker: null,
};

const CLASS_NAMES = {
    ...CONST_CLASS_NAMES,
    soundOff: "video-preview__sound-off",
    soundOn: "video-preview__sound-on",
    previewPlaying: VIDEO_PREVIEW.previewPlayingClassName,

    rowRefreshing: "row--refreshing",
};

const cssVarNames = {
    watchBuffer: "--video-preview-watch-buffer",
    watchProgress: "--video-preview-watch-progress"
};

const DOM_IDS = {
    row: (id) => `row-${id}`,
}

function initDomElements() {
    DOM_ELEMENTS.preview.container = document.getElementById("video-preview-container");
    DOM_ELEMENTS.preview.player = document.getElementById("video-preview-player");
    DOM_ELEMENTS.preview.durationRemaining = document.getElementById("video-preview-duration-remaining");
    DOM_ELEMENTS.preview.progressValue = document.getElementById("video-preview-watch-progress-value");
    DOM_ELEMENTS.preview.progressBuffer = document.getElementById("video-preview-watch-progress-buffer");

    const soundButton = document.getElementById("video-preview-sound");
    if (soundButton != null) {
        DOM_ELEMENTS.sound.button = soundButton;
        DOM_ELEMENTS.sound.iconOff = soundButton.querySelector(`.${CLASS_NAMES.soundOff}`);
        DOM_ELEMENTS.sound.iconOn = soundButton.querySelector(`.${CLASS_NAMES.soundOn}`);
    }
}

export function initVideoPreview() {
    initDomElements();

    if (!DOM_ELEMENTS.preview.container || !DOM_ELEMENTS.preview.player) {
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

export function initVideoPreviewHover(container, previewElementClassName, thumbClassName) {
    if (!container) {
        return;
    }

    container.addEventListener("mouseover", async (event) => {
        const el = event.target.closest(`.${previewElementClassName}`);

        if (!shouldShowVideoPreview(el)) return;

        if (!container.contains(el)) {
            return;
        }

        if (event.relatedTarget && el.contains(event.relatedTarget)) {
            return;
        }

        const itemId = el.dataset.itemId;

        const thumbnail = el.querySelector(`.${thumbClassName}`);

        if (!thumbnail) {
            return;
        }

        clearTimeout(PREVIEW_STATE.hoverTimer);

        PREVIEW_STATE.hoverTimer = setTimeout(async () => {
            if (PREVIEW_STATE.ended) {
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

        const el = event.target.closest(`.${previewElementClassName}`);

        if (!el || !container.contains(el)) {
            return;
        }

        if (event.relatedTarget && el.contains(event.relatedTarget)) {
            return;
        }

        PREVIEW_STATE.ended = false;

        clearTimeout(PREVIEW_STATE.hoverTimer);
        hideVideoPreview();
    });

    document.addEventListener(VIDEO_PREVIEW.playerOpenedEventName, () => {
        clearTimeout(PREVIEW_STATE.hoverTimer);
        hideVideoPreview();
    });

    DOM_ELEMENTS.preview.player.addEventListener("ended", () => {
        PREVIEW_STATE.ended = true;
        hideVideoPreview();
    });

    DOM_ELEMENTS.preview.player.addEventListener("timeupdate", updateVideoPreviewDuration);
}

export function shouldShowVideoPreview(previewElement) {
    if (isMobileScreen()) {
        return false;
    }

    if (!previewElement) {
        return false;
    }

    if (document.body.classList.contains(CLASS_NAMES.listView)) {
        return false;
    }

    if (!previewElement.classList.contains(CLASS_NAMES.rowStatus.success)) {
        return false;
    }

    if (previewElement.dataset.isAudio === "true") {
        return false;
    }

    return Boolean(previewElement.dataset.itemId && previewElement.dataset.media);
}

export async function showVideoPreview(thumbnail, videoUrl, itemId) {
    if (!DOM_ELEMENTS.preview.container || !DOM_ELEMENTS.preview.player) {
        return;
    }

    const positionMs = await watchAPI.getWatchPosition(itemId);

    PREVIEW_STATE.ended = false;

    const itemEl = document.getElementById(DOM_IDS.row(itemId));
    if (itemEl) {
        const isPreviewBlocked = itemEl.classList.contains(CLASS_NAMES.rowRefreshing);
        if (isPreviewBlocked) return;

        itemEl.classList.add(CLASS_NAMES.previewPlaying);
    }

    thumbnail.appendChild(DOM_ELEMENTS.preview.container);

    setWatchTrackerItemId(itemId);

    if (PREVIEW_STATE.currentVideoUrl !== videoUrl) {
        PREVIEW_STATE.currentVideoUrl = videoUrl;
        DOM_ELEMENTS.preview.player.src = videoUrl;

        await new Promise(resolve => {
            DOM_ELEMENTS.preview.player.onloadedmetadata = resolve;
        });
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
        DOM_ELEMENTS.preview.container.hidden = false;
    } catch (error) {
        hideVideoPreview();
        console.debug("Video preview play failed", error);
    }
}

export function hideVideoPreview() {
    if (!DOM_ELEMENTS.preview.container || !DOM_ELEMENTS.preview.player) {
        return;
    }

    const itemEl = document.getElementById(DOM_IDS.row(PREVIEW_STATE.currentItemId));
    if (itemEl) {
        itemEl.classList.remove(CLASS_NAMES.previewPlaying);
    }

    PREVIEW_STATE.currentItemId = null;

    DOM_ELEMENTS.preview.player.pause();

    DOM_ELEMENTS.preview.container.hidden = true;
}

function setVideoPreviewPosition(element) {
    const rect = element.getBoundingClientRect();

    DOM_ELEMENTS.preview.container.style.position = "fixed";

    DOM_ELEMENTS.preview.container.style.left = `${rect.left}px`;
    DOM_ELEMENTS.preview.container.style.top = `${rect.top}px`;

    DOM_ELEMENTS.preview.container.style.width = `${rect.width}px`;
    DOM_ELEMENTS.preview.container.style.height = `${rect.height}px`;
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
            cssVarNames.watchBuffer,
            `${bufferPercent}%`
        );
    }

    if (DOM_ELEMENTS.preview.progressValue !== null) {
        const progressPercent = DOM_ELEMENTS.preview.player.duration > 0
        ? Math.floor((DOM_ELEMENTS.preview.player.currentTime / DOM_ELEMENTS.preview.player.duration) * 1000) / 10
        : 0;

        DOM_ELEMENTS.preview.progressValue.style.setProperty(
            cssVarNames.watchProgress,
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

export function initVideoPreviewScroll(container, previewElementClassName, thumbClassName) {
    if (!container) {
        return;
    }

    const refreshPreview = (force=false) => {
        updateVideoPreview(container, previewElementClassName, thumbClassName, force);
    };    

    window.addEventListener("scroll", refreshPreview, { passive: true });
    window.addEventListener("resize", refreshPreview);

    // We launch it immediately after opening the page.
    refreshPreview();

    return refreshPreview;
}

function updateVideoPreview(container, previewElementClassName, thumbClassName, force=false) {
    if (!isMobileScreen()) {
        return;
    }

    if (PREVIEW_STATE.currentItemId) {
        const element = document.getElementById(
            ids.row(PREVIEW_STATE.currentItemId)
        );

        if (force || !element || !isElementInViewport(element)) {
            hideVideoPreview();
        }
    }

    clearTimeout(PREVIEW_STATE.scrollTimer);

    PREVIEW_STATE.scrollTimer = setTimeout(() => {
        updateCenteredPreview(container, previewElementClassName, thumbClassName);
    }, 120);
}

function isElementInViewport(element) {
    const rect = element.getBoundingClientRect();

    return rect.bottom > 0 && rect.top < window.innerHeight;
}

async function updateCenteredPreview(container, previewElementClassName, thumbClassName) {
    const element = findCenteredElement(container, previewElementClassName);

    if (!element) {
        hideVideoPreview();
        return;
    }

    if (!element.classList.contains(CLASS_NAMES.rowStatus.success)) {
        return;
    }

    const itemId = element.dataset.itemId;

    if (itemId === PREVIEW_STATE.currentItemId) {
        return;
    }

    const requestId = ++PREVIEW_STATE.requestId;

    hideVideoPreview();

    // While waiting for a response, the user has already scrolled through the list.
    if (requestId !== PREVIEW_STATE.requestId) {
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

function findCenteredElement(container, previewElementClassName) {
    const viewportCenter = window.innerHeight / 2;

    let bestElement = null;
    let bestDistance = Number.MAX_VALUE;

    const items = container.querySelectorAll(`.${previewElementClassName}`);

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