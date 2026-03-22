import { DOM_IDS } from "./dom-ids.js";

export const DOM_ELEMENTS = {
    mediaURL: null,
    inputActionBtn: null,
    resultInfo: null,

    sysInfoServerStatusDot: null,
    sysInfoDiskFree: null,
    sysInfoDiskUsed: null,
};

document.addEventListener("DOMContentLoaded", () => {
    DOM_ELEMENTS.mediaURL = document.getElementById(DOM_IDS.mediaURL);
    DOM_ELEMENTS.inputActionBtn = document.getElementById(DOM_IDS.inputActionBtn);
    DOM_ELEMENTS.resultInfo = document.getElementById(DOM_IDS.resultInfo);

    DOM_ELEMENTS.sysInfoServerStatusDot = document.getElementById(DOM_IDS.sysInfoServerStatusDot);
    DOM_ELEMENTS.sysInfoDiskFree = document.getElementById(DOM_IDS.sysInfoDiskFree);
    DOM_ELEMENTS.sysInfoDiskUsed = document.getElementById(DOM_IDS.sysInfoDiskUsed);

    if (!DOM_ELEMENTS.sysInfoServerStatusDot || !DOM_ELEMENTS.sysInfoDiskFree || !DOM_ELEMENTS.sysInfoDiskUsed) {
        console.warn("System info elements not found in DOM");
    }
});