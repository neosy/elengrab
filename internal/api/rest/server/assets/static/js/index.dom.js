export const DOM_IDS = {
    row: (id) => `row-${id}`,
    progressRow: (id) => `progress-${id}`,
};

export const DOM_ELEMENTS = {
    resultRows: null,
    rowTopPlaceholder: null,
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
    DOM_ELEMENTS.rowTopPlaceholder = document.getElementById("row-top-placeholder");
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

    if (!DOM_ELEMENTS.sysInfoServerStatusDot || !DOM_ELEMENTS.sysInfoDiskFree || !DOM_ELEMENTS.sysInfoDiskUsed) {
        console.warn("System info elements not found in DOM");
    }
}