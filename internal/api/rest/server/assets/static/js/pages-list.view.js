import {CLASS_NAMES, STORAGE_KEYS, STORAGE_VALUES, API_PATHS} from './constants.js';
import storageState from './storage-state.js';
import * as notify from './notifications.js';
import { CLASS_SELECTORS, DOM_IDS, CONTENT_TYPES } from "./pages-list.dom.js";

export function applyGridView(isGridView) {
    document.body.classList.toggle(CLASS_NAMES.body.gridView, isGridView);
    document.body.classList.toggle(CLASS_NAMES.body.listView, !isGridView);
}

function getCurrentCardLayout() {
    const contentType = CONTENT_TYPES.getCurrent();

    const settings = {
        [CONTENT_TYPES.inProgress]: {
            key: STORAGE_KEYS.settingsCardLayout.inProgress,
            defaultValue: STORAGE_VALUES.cardLayout.defaultInProgress,
        },
    }[contentType] ?? {
        key: STORAGE_KEYS.settingsCardLayout.default,
        defaultValue: STORAGE_VALUES.cardLayout.default,
    };

    return {
        key: settings.key,
        value: storageState.get(settings.key, settings.defaultValue),
    };
}

export function toggleGridView() {
    const { key: cardLayoutKey, value: currentCardLayout } = getCurrentCardLayout();

    const nextIsGrid = currentCardLayout !== STORAGE_VALUES.cardLayout.grid;
    const nextCardLayout = nextIsGrid
        ? STORAGE_VALUES.cardLayout.grid
        : STORAGE_VALUES.cardLayout.list;

    storageState.set(cardLayoutKey, nextCardLayout);

    applyGridView(nextIsGrid);

    return nextIsGrid;
}

export function initGridView() {
    const { value: cardLayout } = getCurrentCardLayout();
    const isGridView = cardLayout === STORAGE_VALUES.cardLayout.grid;

    applyGridView(isGridView);

    return isGridView;
}

export function getGridView() {
    return getCurrentCardLayout().value === STORAGE_VALUES.cardLayout.grid;
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
 * @returns {{ setActiveTab: (value: string, options?: { reload?: boolean }) => Promise<void> } | undefined}
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

        applyTab(tab);

        setViewMode();
    });

    function applyTab(tab) {
        tabs.dataset.viewMode = tab.dataset.viewMode;

        tabs.querySelectorAll(CLASS_SELECTORS.viewModeTab).forEach((item) => {
            const selected = item === tab;

            item.classList.toggle('active', selected);
            item.setAttribute('aria-selected', String(selected));
        });
    }

    /**
     * Switch to the tab with the given value.
     * By default it also reloads the rows; pass { reload: false } to only update the UI.
     */
    async function setActiveTab(value, { reload = true } = {}) {
        const tab = Array.from(tabs.querySelectorAll(CLASS_SELECTORS.viewModeTab))
            .find((item) => item.dataset.viewMode === value);

        if (!tab) {
            return;
        }

        applyTab(tab);

        if (reload) {
            await setViewMode();
        }
    }

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

    return {
        setActiveTab,
    };
}

/**
 * @param {Object} options
 * @param {() => Object} options.getSearchQueryValues
 * @param {() => void} options.onSuccess
 */
export function initContentTypeBar({
    getSearchQueryValues,
    onSuccess,
}) {
    const tabs = document.querySelector(CLASS_SELECTORS.contentTypeTabs);

    if (!tabs) {
        return;
    }

    tabs.addEventListener('click', (event) => {
        const tab = event.target.closest(CLASS_SELECTORS.contentTypeTab);

        if (!tab || tab.classList.contains('active')) {
            return;
        }

        tabs.dataset.contentType = tab.dataset.contentType;

        tabs.querySelectorAll(CLASS_SELECTORS.contentTypeTab).forEach((item) => {
            const selected = item === tab;

            item.classList.toggle('active', selected);
            item.setAttribute('aria-selected', String(selected));
        });

        setContentType();
    });

    async function setContentType() {
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

    return;
}
