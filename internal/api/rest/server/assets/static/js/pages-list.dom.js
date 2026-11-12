import { CLASS_NAMES as CONST_CLASS_NAMES } from './constants.js';

export const CLASS_NAMES = {
    ...CONST_CLASS_NAMES,
    
    rowRemoving: "row--removing",

    mediaResultRow: "media-result__row",

    viewModeTabs: "view-mode-tabs",
    viewModeTab: "view-mode-tab",

    mediaResultRowThumbnail: "media-result__row-thumbnail",
    mediaResultThumbnailPlayButton: "media-result__thumbnail-play-button",
    mediaResultRowThumbnailImageWrapper: "media-result__thumbnail-image__wrapper",
    mediaResultThumbnailPlaceholder: "media-result__thumbnail-placeholder",
    mediaResultRowThumbnailWatched: "media-result__thumbnail-watched",

    thumbnailWatchProgress: "thumbnail-watch-progress",
    thumbnailWatchProgressHidden: "thumbnail-watch-progress--hidden",
    thumbnailWatchProgressValue: "thumbnail-watch-progress-value",

    mediaResultAvatar: "media-result__channel",

    mediaResultTitleLabel: "media-result__title-label",
    mediaResultViewCount: "media-result__content-view-count",
    mediaResultVisibility: "media-result__content-visibility",
    mediaResultShareLink: "media-result__content-share-link",
};

export const CLASS_SELECTORS = Object.fromEntries(
    Object.entries(CLASS_NAMES)
        .filter(([, value]) => typeof value === "string")
        .map(([key, value]) => [key, `.${value}`])
);

export const CSS_VARIABLE_NAMES = {
    watchProgress: '--watch-progress',
};

export const DOM_IDS = {
    mediaResultRows: 'media-result-rows',
    mediaResultItems: 'media-result-items',
    rowNoItems: "row-no-items",

    row: (id) => `row-${id}`,
    progress: (id) => `progress-${id}`,
}

export const DOM_ELEMENTS = {
    result: null,
    mediaResultItems: null,
    
    resultInfo: null,
    resultInfoRow: null,
    resultInfoFailed: null,

    mediaPlayer: null,
};

function $(root, className) {
    return root?.querySelector(`.${className}`);
}

export function initDomElements(elements = DOM_ELEMENTS) {
    elements.result = document.getElementById("media-result");

    elements.resultInfo = document.getElementById("result-info");
    elements.resultInfoRow = document.getElementById("result-info-row");
    elements.resultInfoFailed = document.getElementById("result-info-failed");

    elements.mediaPlayer = document.getElementById("media-player");
}

export const VIEW_MODES = {
    New: "new",
    Popular: "popular",
    Old: "old",

    getCurrent: getViewMode,
    isNew: isNewViewMode,
    isPopular: isPopularViewMode,
    isOld: isOldViewMode,
};

/**
 * @returns {string}
 */
export function getViewMode() {
    const viewModeTabs = $(document, CLASS_NAMES.viewModeTabs);
    return viewModeTabs?.dataset.viewMode || "";
}

/**
 * @returns {boolean}
 */
export function isNewViewMode() {
    return getViewMode() === VIEW_MODES.New;
}

/**
 * @returns {boolean}
 */
export function isPopularViewMode() {
    return getViewMode() === VIEW_MODES.Popular;
}

/**
 * @returns {boolean}
 */
export function isOldViewMode() {
    return getViewMode() === VIEW_MODES.Old;
}
