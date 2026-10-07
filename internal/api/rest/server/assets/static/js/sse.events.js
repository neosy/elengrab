import * as notify from './notifications.js';

const DOM_ELEMENTS = {
    sysInfoDiskFree: null,
    sysInfoDiskUsed: null,
};

function initDomElements(elements = DOM_ELEMENTS) {    
    elements.sysInfoDiskFree = document.getElementById("disk-free");
    elements.sysInfoDiskUsed = document.getElementById("disk-used");
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