import { DOM_IDS } from "./index-dom-ids.js";

export const DOM_ELEMENTS = {
    mediaURL: null,
    inputActionBtn: null,
    inputActionSettingsBtn: null,
    grabOptionsCollapse: null,
    grabOptions: null,

    resultInfo: null,
    resultInfoRow: null,
    resultInfoFailed: null,

    sysInfoServerStatusDot: null,
    sysInfoDiskFree: null,
    sysInfoDiskUsed: null,
};

document.addEventListener("DOMContentLoaded", () => {
    DOM_ELEMENTS.mediaURL = document.getElementById(DOM_IDS.mediaURL);
    DOM_ELEMENTS.inputActionBtn = document.getElementById(DOM_IDS.inputActionBtn);
    DOM_ELEMENTS.inputActionSettingsBtn = document.getElementById(DOM_IDS.inputActionSettingsBtn);
    DOM_ELEMENTS.grabOptionsCollapse = document.getElementById(DOM_IDS.grabOptionsCollapse);
    DOM_ELEMENTS.grabOptions = document.getElementById(DOM_IDS.grabOptions);

    DOM_ELEMENTS.resultInfo = document.getElementById(DOM_IDS.resultInfo);
    DOM_ELEMENTS.resultInfoRow = document.getElementById(DOM_IDS.resultInfoRow);
    DOM_ELEMENTS.resultInfoFailed = document.getElementById(DOM_IDS.resultInfoFailed);

    DOM_ELEMENTS.sysInfoServerStatusDot = document.getElementById(DOM_IDS.sysInfoServerStatusDot);
    DOM_ELEMENTS.sysInfoDiskFree = document.getElementById(DOM_IDS.sysInfoDiskFree);
    DOM_ELEMENTS.sysInfoDiskUsed = document.getElementById(DOM_IDS.sysInfoDiskUsed);

    if (!DOM_ELEMENTS.sysInfoServerStatusDot || !DOM_ELEMENTS.sysInfoDiskFree || !DOM_ELEMENTS.sysInfoDiskUsed) {
        console.warn("System info elements not found in DOM");
    }
});