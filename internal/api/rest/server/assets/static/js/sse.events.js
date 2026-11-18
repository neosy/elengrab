import * as notify from './notifications.js';
import { CSS_VARIABLE_NAMES } from './constants.js';
import { store } from './store.js';

const DOM_ELEMENTS = {
    sysInfoDiskFree: null,
    sysInfoDiskUsed: null,
};

function initDomElements(elements = DOM_ELEMENTS) {    
    elements.sysInfoDiskFree = document.getElementById("disk-free");
    elements.sysInfoDiskUsed = document.getElementById("disk-used");
}

export function handleDownloadsInProgress(event) {
    try {
        const data = JSON.parse(event.data);

        store.downloadsCount.userInProgress = Number.isFinite(data.userInProgressCount) ? data.userInProgressCount : 0;
        store.downloadsCount.userFailed = Number.isFinite(data.userFailedCount) ? data.userFailedCount : 0;
        store.downloadsCount.userTotal = Number.isFinite(data.userCount) ? data.userCount : 0;
        store.downloadsCount.totalInProgress = Number.isFinite(data.totalInProgressCount) ? data.totalInProgressCount : 0;
        store.downloadsCount.totalFailed = Number.isFinite(data.totalFailedCount) ? data.totalFailedCount : nu0l;
        store.downloadsCount.total = Number.isFinite(data.totalCount) ? data.totalCount : 0;

        document.body.classList.toggle(
            "has-user-downloads-in-progress",
            store.downloadsCount.userInProgress > 0
        );

        document.body.classList.toggle(
            "has-user-downloads-failed",
            store.downloadsCount.userFailed > 0
        );

        document.body.classList.toggle(
            "has-downloads-in-progress",
            store.downloadsCount.userInProgress > 0 || store.downloadsCount.totalInProgress > 0
        );
       
        document.body.classList.toggle(
            "has-downloads-failed",
            store.downloadsCount.userFailed > 0 || store.downloadsCount.totalFailed > 0
        );

        document.body.classList.toggle(
            "has-total-downloads-in-progress",
            store.downloadsCount.totalInProgress > 0
        );

        document.body.classList.toggle(
            "has-total-downloads-failed",
            store.downloadsCount.totalFailed > 0
        );
        
        // Counts for CSS (e.g. a badge on the tab): set when positive, removed when zero
        setCssVariable(CSS_VARIABLE_NAMES.downloadsCount.userInProgress, store.downloadsCount.userInProgress > 0 ? store.downloadsCount.userInProgress : null);
        setCssVariable(CSS_VARIABLE_NAMES.downloadsCount.userFailed, store.downloadsCount.userFailed > 0 ? store.downloadsCount.userFailed : null);
        setCssVariable(CSS_VARIABLE_NAMES.downloadsCount.totalInProgress, store.downloadsCount.totalInProgress > 0 ? store.downloadsCount.totalInProgress : null);
        setCssVariable(CSS_VARIABLE_NAMES.downloadsCount.totalFailed, store.downloadsCount.totalFailed > 0 ? store.downloadsCount.totalFailed : null);
    } catch (err) {
        console.error(`SSE ${event.type} handler error:`, err);
    }
}

function setCssVariable(name, value, element = document.body) {
    if (value === null || value === undefined || value === "") {
        element.style.removeProperty(name);
    } else {
        element.style.setProperty(name, value);
    }
}

export function handleSystemInfoUpdate(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.diskFree || !data.diskUsed) return

        if (!DOM_ELEMENTS.sysInfoDiskFree || !DOM_ELEMENTS.sysInfoDiskUsed) return;

        DOM_ELEMENTS.sysInfoDiskFree.textContent = data.diskFree
        DOM_ELEMENTS.sysInfoDiskUsed.textContent = data.diskUsed
    } catch (err) {
        console.error(`SSE ${event.type} handler error:`, err);
    }
}

export function handleNotification(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.module || !data.type || !data.message) return

        notify.show(data.message, data.type)
    } catch (err) {
        console.error(`SSE ${event.type} handler error:`, err);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    initDomElements();
})