// constants.js

export const STORAGE_KEYS = {
    grabOptionsCollapsed: "grabOptionsCollapsed",
    settingsGridView: "settingsGridView",
};

// Action button icon URLs
export const ICON_PASTE = 'static/icons/action-paste-v2-icon.svg';
export const ICON_CLEAR = 'static/icons/action-clear-icon.svg';

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

// Media watch constants
export const MEDIA_WATCH = {
    // Minimum watched time from the beginning to restore playback position
    startThresholdMs: 8000,

    // Minimum watch interval duration to send a watch event
    minIntervalMs: 2000,

    // Maximum allowed watch interval duration (with playback speed tolerance)
    maxIntervalMs: 15500,
};

export const VIDEO_PREVIEW = {
    // Event dispatched when the full video player is opened
    playerOpenedEventName: "video-preview-player-opened",

    previewPlayingClassName: "video-preview-playing",
};
