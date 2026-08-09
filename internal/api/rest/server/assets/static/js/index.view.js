import * as constants from './constants.js';
import storageState from './storage-state.js';
import * as notify from './notifications.js';

export function applyGridView(isGridView) {
    document.body.classList.toggle(constants.CLASS_NAMES.gridView, isGridView);
    document.body.classList.toggle(constants.CLASS_NAMES.listView, !isGridView);
}

export function toggleGridView() {
    const current = storageState.get(constants.STORAGE_KEYS.settingsGridView, true);
    const next = !current;

    storageState.set(constants.STORAGE_KEYS.settingsGridView, next);

    applyGridView(next);

    return next;
}

export function initGridView() {
    const isGridView = storageState.get(constants.STORAGE_KEYS.settingsGridView, true);

    applyGridView(isGridView);

    return isGridView;
}

export function getGridView() {
    return storageState.get(constants.STORAGE_KEYS.settingsGridView, true);
}

export function initSearching(clear) {
    const searchBtn = document.getElementById("userMenuSearchButton");
    const backBtn = document.getElementById("historySearchBackButton");
    const header = document.getElementById("header");
    const searchInput = document.getElementById("historySearchInput");

    if (!searchBtn || !header || !backBtn) return;

    searchBtn.addEventListener('click', () => {
        openSearching(header, searchInput);
    });    

    backBtn.addEventListener('click', () => {
        closeSearching(header, clear);
    });
}

export function initHeaderUserMenu() {
    const btn = document.getElementById("userMenudownloadButton");
    const grabInput = document.getElementById("mediaURLInput");

    if (!btn) return;

    btn.addEventListener('click', () => {
        selectGrabInput(grabInput);
    });    
}

function openSearching(header, input) {
    header.classList.toggle(constants.CLASS_NAMES.isSearch, true);
    input.focus();
}

function closeSearching(header, clear) {
    header.classList.toggle(constants.CLASS_NAMES.isSearch, false);
    clear();
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
}