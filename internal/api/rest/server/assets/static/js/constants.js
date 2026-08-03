// constants.js

// -------------------------------------------------------------
// Element selectors and cookie names
// -------------------------------------------------------------
export const SELECT_NAMES = {
    qualityCodec: "quality-codec",
    qualityResolution: "quality-resolution",
    format: "format"
};

export const COOKIE_NAMES = {
    qualityCodec: "select_quality_codec",
    qualityResolution: "select_quality_resolution",
    format: "select_format"
};

export const STORAGE_KEYS = {
    grabOptionsCollapsed: "grabOptionsCollapsed",
    settingsGridView: "settingsGridView",
}

// Action button icon URLs
export const ICON_PASTE = 'static/icons/action-paste-v2-icon.svg';
export const ICON_CLEAR = 'static/icons/action-clear-icon.svg';

// Class names
export const CLASS_NAMES = {
    gridView: "grid-view",
    listView: "list-view",
    isSearch: "is-search",
}
// Media watch constants
export const MEDIA_WATCH = {
    // Minimum watched time from the beginning to restore playback position
    startThresholdMs: 8000,

    // Minimum watch interval duration to send a watch event
    minIntervalMs: 2000,

    // Maximum allowed watch interval duration (with playback speed tolerance)
    maxIntervalMs: 15500,
};
