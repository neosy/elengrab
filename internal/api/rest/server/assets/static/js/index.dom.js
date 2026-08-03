export const DOM_IDS = {
    rowTopPlaceholder: "row-top-placeholder",
    
    row: (id) => `row-${id}`,
    progressRow: (id) => `progress-${id}`,
};

export const DOM_CLASSES = {
    rowRefreshing: "row--refreshing",
    rowRemoving: "row--removing",
    mediaResultRows: "media-result__rows",
    mediaResultRow: "media-result__row",
    mediaResultRowThumbnail: "media-result__row-thumbnail",
    mediaResultRowThumbnailImageWrapper: "media-result__thumbnail-image__wrapper",
};

export const DOM_ELEMENTS = {
    resultRows: null,
    rowNoItems: null,

    grabForm: null,
    mediaURLInput: null,
    inputActionBtn: null,
    inputActionSettingsBtn: null,
    grabOptionsCollapse: null,
    grabOptions: null,

    historySearchInputWrapper: null,
    historySearchClearButton: null,

    resultInfo: null,
    resultInfoRow: null,
    resultInfoFailed: null,

    sysInfoServerStatusDot: null,
    sysInfoDiskFree: null,
    sysInfoDiskUsed: null,
};

export function initDomElements() {    
    DOM_ELEMENTS.resultRows = document.getElementById("media-result-rows");
    DOM_ELEMENTS.rowNoItems = document.getElementById("row-no-items");

    DOM_ELEMENTS.grabForm = document.getElementById("grabForm");
    DOM_ELEMENTS.mediaURLInput = document.getElementById("mediaURLInput");
    DOM_ELEMENTS.inputActionBtn = document.getElementById("inputActionBtn");
    DOM_ELEMENTS.inputActionSettingsBtn = document.getElementById("inputActionSettingsBtn");
    DOM_ELEMENTS.grabOptionsCollapse = document.getElementById("grabOptionsCollapse");
    DOM_ELEMENTS.grabOptions = document.getElementById("grabOptions");

    DOM_ELEMENTS.historySearchInputWrapper = document.getElementById("historySearchInputWrapper");
    DOM_ELEMENTS.historySearchClearButton = document.getElementById("historySearchClearButton");

    DOM_ELEMENTS.resultInfo = document.getElementById("result-info");
    DOM_ELEMENTS.resultInfoRow = document.getElementById("result-info-row");
    DOM_ELEMENTS.resultInfoFailed = document.getElementById("result-info-failed");

    DOM_ELEMENTS.sysInfoServerStatusDot = document.getElementById("server-status-dot");
    DOM_ELEMENTS.sysInfoDiskFree = document.getElementById("disk-free");
    DOM_ELEMENTS.sysInfoDiskUsed = document.getElementById("disk-used");
}