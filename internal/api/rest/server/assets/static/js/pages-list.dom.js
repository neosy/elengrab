export const CLASS_NAMES = {
    rowRemoving: "row--removing",

    mediaResultRow: "media-result__row",

    viewModeTabs: "view-mode-tabs",
    viewModeTab: "view-mode-tab",

    mediaExtChannelLinkButton: "channel-button--link",

    mediaResultRowThumbnail: "media-result__row-thumbnail",
    mediaResultThumbnailPlayButton: "media-result__thumbnail-play-button",
    mediaResultRowThumbnailImageWrapper: "media-result__thumbnail-image__wrapper",
    mediaResultThumbnailPlaceholder: "media-result__thumbnail-placeholder",
    mediaResultRowThumbnailWatched: "media-result__thumbnail-watched",
    mediaResultRowThumbnailWatchProgress: "media-result__thumbnail-watch-progress",
    mediaResultRowThumbnailWatchProgressValue: "media-result__thumbnail-watch-progress-value",

    mediaResultAvatar: "media-result__channel",

    mediaResultTitleLink: "media-result__title-link",
    mediaResultViewCount: "media-result__content-view-count",
    mediaResultVisibility: "media-result__content-visibility",
    mediaResultShareLink: "media-result__content-share-link",

    row: {
        rowRefreshing: "row--refreshing",
    },

    watchProgress: "watch-progress",
};

export const CLASS_SELECTORS = Object.fromEntries(
    Object.entries(CLASS_NAMES)
        .filter(([, value]) => typeof value === "string")
        .map(([key, value]) => [key, `.${value}`])
);

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

    sysInfoDiskFree: null,
    sysInfoDiskUsed: null,
};

export function initDomElements(elements = DOM_ELEMENTS) {    
    elements.result = document.getElementById("media-result");

    elements.resultInfo = document.getElementById("result-info");
    elements.resultInfoRow = document.getElementById("result-info-row");
    elements.resultInfoFailed = document.getElementById("result-info-failed");

    elements.mediaPlayer = document.getElementById("media-player");

    elements.sysInfoDiskFree = document.getElementById("disk-free");
    elements.sysInfoDiskUsed = document.getElementById("disk-used");
}