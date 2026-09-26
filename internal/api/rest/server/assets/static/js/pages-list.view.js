import {CLASS_NAMES, STORAGE_KEYS, API_PATHS} from './constants.js';
import storageState from './storage-state.js';
import * as notify from './notifications.js';
import { DOM_ELEMENTS, CLASS_SELECTORS, DOM_IDS } from "./pages-list.dom.js";

export function applyGridView(isGridView) {
    document.body.classList.toggle(CLASS_NAMES.gridView, isGridView);
    document.body.classList.toggle(CLASS_NAMES.listView, !isGridView);
}

export function toggleGridView() {
    const current = storageState.get(STORAGE_KEYS.settingsGridView, true);
    const next = !current;

    storageState.set(STORAGE_KEYS.settingsGridView, next);

    applyGridView(next);

    return next;
}

export function initGridView() {
    const isGridView = storageState.get(STORAGE_KEYS.settingsGridView, true);

    applyGridView(isGridView);

    return isGridView;
}

export function getGridView() {
    return storageState.get(STORAGE_KEYS.settingsGridView, true);
}

export function initHeaderUserMenu() {
    const btn = document.getElementById("headerActionsDownloadButton");
    const grabInput = document.getElementById("mediaURLInput");

    if (!btn) return;

    btn.addEventListener('click', () => {
        selectGrabInput(grabInput);
    });    
}

function selectGrabInput(grabInput) {
    if (!grabInput) {
        notify.show("You do not have permission to add downloads", notify.notifyType.ERROR);
        return;
    }

    grabInput.scrollIntoView({
        behavior: "smooth",
        block: "center",
    });

    grabInput.focus();
}

export function initLazyImages({
    containerSelector,
    imageSelector = "img[data-src]",
    placeholderSelector,
    rootMargin = "300px 0px",
}) {
    const observer = new IntersectionObserver(
        (entries) => {
            for (const entry of entries) {
                if (!entry.isIntersecting) {
                    continue;
                }

                loadImages(entry.target);
                observer.unobserve(entry.target);
            }
        },
        {
            rootMargin,
        },
    );

    function observe(root) {
        if (root.matches?.(containerSelector)) {
            observer.observe(root);
        }

        for (const container of root.querySelectorAll(containerSelector)) {
            observer.observe(container);
        }
    }

    function loadImages(container) {
        for (const image of container.querySelectorAll(imageSelector)) {
            const placeholder = placeholderSelector
                ? image.parentElement.querySelector(placeholderSelector)
                : null;

            const removePlaceholder = () => {
                requestAnimationFrame(() => {
                    placeholder?.remove();
                    image.removeAttribute("data-src");
                });
            };

            image.addEventListener("load", removePlaceholder, { once: true });

            image.src = image.dataset.src;

            if (image.complete && image.naturalWidth > 0) {
                removePlaceholder();
            }
        }
    }

    observe(document);

    document.addEventListener("htmx:afterSwap", (event) => {
        observe(event.target);
    });

    return { observe };
}

/**
 * @param {Object} options
 * @param {() => Object} options.getSearchQueryValues
 * @param {() => void} options.onSuccess
 */
export function initViewModeBar({
    getSearchQueryValues,
    onSuccess,
}) {
    const tabs = document.querySelector(CLASS_SELECTORS.viewModeTabs);

    if (!tabs) {
        return;
    }

    tabs.addEventListener('click', (event) => {
        const tab = event.target.closest(CLASS_SELECTORS.viewModeTab);

        if (!tab) {
            return;
        }

        tabs.dataset.viewMode = tab.dataset.viewMode;

        tabs.querySelectorAll(CLASS_SELECTORS.viewModeTab).forEach((item) => {
            const selected = item === tab;

            item.classList.toggle('active', selected);
            item.setAttribute('aria-selected', String(selected));
        });

        setViewMode();
    });

    async function setViewMode() {
        const rows = document.getElementById(DOM_IDS.mediaResultRows);

        if (!rows) {
            return;
        }

        const searchQuery = {
            ...getSearchQueryValues(),
        };

        const response = await fetch(API_PATHS.downloaderSearch, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(searchQuery),
        });

        if (!response.ok) {
            throw new Error(`Search failed: ${response.status}`);
        }

        rows.outerHTML = await response.text();

        const newRows = document.getElementById(DOM_IDS.mediaResultRows);

        htmx.process(newRows);

        onSuccess(newRows);
    }
}

export function initExtChannelLink({
    getSearchQueryValues,
    onSuccess,
}) {
    document.addEventListener("click", async (event) => {
        const linkBtn = event.target.closest(CLASS_SELECTORS.mediaExtChannelLinkButton);
        if (!linkBtn || linkBtn.dataset.channelId === '') return;

        const channelId = linkBtn.dataset.channelId;

        click(channelId);
    })

    async function click(channelId) {
        const rows = document.getElementById(DOM_IDS.mediaResultRows);

        if (!rows) {
            return;
        }

        const searchQuery = {
            ...getSearchQueryValues(),
            
            channelId: channelId,
        };

        const response = await fetch(API_PATHS.downloaderSearch, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(searchQuery),
        });

        if (!response.ok) {
            throw new Error(`Search failed: ${response.status}`);
        }

        rows.outerHTML = await response.text();

        const newRows = document.getElementById(DOM_IDS.mediaResultRows);

        htmx.process(newRows);

        window.scrollTo(0, 0);

        onSuccess(newRows);
    }
}