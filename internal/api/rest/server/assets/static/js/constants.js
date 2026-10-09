// constants.js

export const STORAGE_KEYS = {
    grabOptionsCollapsed: "grabOptionsCollapsed",
    settingsGridView: "settingsGridView",
};

export const ICONS = {
    // Action button icon URLs
    paste: 'static/icons/action-paste-v2-icon.svg',
    clear: 'static/icons/action-clear-icon.svg',
};

// Internal API request paths.
export const API_PATHS = {
    downloaderItems: "/downloader/items",
    downloaderSearch: "/downloader/search",
}

// Internal API request path templates.
export const API_PATH_TEMPLATES = {
    downloaderWatchTracking: "/downloader/items/{itemId}/watch-tracking",
    downloaderWatchPosition: "/downloader/items/{itemId}/watch-position",
}

export const CLASS_PREFIXES = {
    visibility: "visibility--",
}

// Class names
export const CLASS_NAMES = {
    gridView: "grid-view",
    listView: "list-view",

    rowStatus: {
        success: "success",
        inProgress: "inprogress",
    },

    ui: {
        blockingActive: "ui-blocking-active",
    },
};

export const EVENT_NAMES = {
    videoPreviewStop: 'video-preview-stop',

    // Event dispatched when the full video player is opened
    videoPlayerOpened: "video-player-opened",
};

// Media watch constants
export const MEDIA_WATCH = {
    // Minimum watched time from the beginning to restore playback position
    startThresholdMs: 5000,

    // Minimum watch interval duration to send a watch event
    minIntervalMs: 1200,

    // Maximum allowed watch interval duration (with playback speed tolerance)
    maxIntervalMs: 15500,

    // Interval between periodic watch-tracking events
    heartbeatIntervalMs: 5000,
};
