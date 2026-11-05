import * as notify from "./notifications.js";
import {formatDuration} from "player-utils.js";

const CLASS_NAMES = {
    urlInput: "download-list-dialog__input",
    fetchButton: "download-list-dialog__button--fetch",
    loading: "download-list-dialog__loading",
    results: "download-list-dialog__results",
    list: "download-list-dialog__list",
    itemTemplate: "download-list-dialog__item-template",
    selectAll: "download-list-dialog__checkbox--all",
    counter: "download-list-dialog__counter",
    submitButton: "download-list-dialog__button--submit",
    submitSpinner: "download-list-dialog__spinner--small",
    submitText: "download-list-dialog__submit-text",
    itemCheckbox: "download-list-dialog__checkbox--item",
    itemLabel: "download-list-dialog__item-label",
    itemTitle: "download-list-dialog__item-title",
    itemDuration: "download-list-dialog__item-duration",
    itemUrl: "download-list-dialog__item-url",
};

const DOM_IDS = {
    dialog: "download-list-dialog",
};

const TEXTS = {
    submitDefault: "Download Selected",
    submitBusy: "Sending…",
};

/* ------------------------------------------------------------------ */
/* Helpers                                                             */
/* ------------------------------------------------------------------ */

function $(root, className) {
    return root.querySelector(`.${className}`);
}

function getDomElements(dialog) {
    const elements = {
        dialog,
        urlInput: $(dialog, CLASS_NAMES.urlInput),
        fetchButton: $(dialog, CLASS_NAMES.fetchButton),
        loading: $(dialog, CLASS_NAMES.loading),
        results: $(dialog, CLASS_NAMES.results),
        list: $(dialog, CLASS_NAMES.list),
        itemTemplate: $(dialog, CLASS_NAMES.itemTemplate),
        selectAll: $(dialog, CLASS_NAMES.selectAll),
        counter: $(dialog, CLASS_NAMES.counter),
        submitButton: $(dialog, CLASS_NAMES.submitButton),
        submitSpinner: $(dialog, CLASS_NAMES.submitSpinner),
        submitText: $(dialog, CLASS_NAMES.submitText),
    };

    return Object.values(elements).some((el) => el === null) ? null : elements;
}

function getItemCheckboxes(elements) {
    return elements.list.querySelectorAll(`.${CLASS_NAMES.itemCheckbox}`);
}

function getSelectedItems({ elements, state }) {
    const selected = [];

    getItemCheckboxes(elements).forEach((checkbox, index) => {
        if (checkbox.checked) {
            selected.push(state.items[index]);
        }
    });

    return selected;
}

/* ------------------------------------------------------------------ */
/* UI state                                                            */
/* ------------------------------------------------------------------ */

function updateControls(ctx) {
    const { elements, state } = ctx;
    const total = getItemCheckboxes(elements).length;
    const selectedCount = getSelectedItems(ctx).length;

    elements.counter.textContent = `Selected: ${selectedCount} of ${total}`;
    elements.selectAll.checked = total > 0 && selectedCount === total;
    elements.selectAll.indeterminate = selectedCount > 0 && selectedCount < total;

    elements.fetchButton.disabled = state.busy || !elements.urlInput.value.trim();
    elements.submitButton.disabled = state.busy || selectedCount === 0;
}

function setBusy(ctx, value) {
    const { elements, state } = ctx;

    state.busy = value;
    elements.urlInput.disabled = value;
    elements.results.setAttribute("aria-busy", String(value));
    updateControls(ctx);
}

function setSubmitting(elements, value) {
    elements.submitSpinner.hidden = !value;
    elements.submitText.textContent = value ? TEXTS.submitBusy : TEXTS.submitDefault;
}

function renderItems(ctx) {
    const { elements, state } = ctx;
    const fragment = document.createDocumentFragment();

    state.items.forEach((item, index) => {
        const node = elements.itemTemplate.content.cloneNode(true);
        const checkbox = $(node, CLASS_NAMES.itemCheckbox);

        checkbox.id = `downloadListItem${index}`;
        checkbox.checked = true;
        $(node, CLASS_NAMES.itemTitle).textContent = item.title || item.url;
        $(node, CLASS_NAMES.itemUrl).textContent = item.url;
        $(node, CLASS_NAMES.itemDuration).textContent = formatDuration(item.durationSec);
        $(node, CLASS_NAMES.itemLabel).title = item.url;
        

        fragment.appendChild(node);
    });

    elements.list.replaceChildren(fragment);
    elements.results.hidden = state.items.length === 0;
    updateControls(ctx);
}

function resetDialog(ctx) {
    const { elements, state } = ctx;

    state.controller?.abort();
    state.controller = null;
    state.items = [];

    elements.list.replaceChildren();
    elements.results.hidden = true;
    elements.loading.hidden = true;
    elements.urlInput.value = "";
    setSubmitting(elements, false);
    setBusy(ctx, false);
}

/* ------------------------------------------------------------------ */
/* Server requests                                                     */
/* ------------------------------------------------------------------ */

async function fetchList(ctx) {
    const { elements, state, options } = ctx;
    const pageUrl = elements.urlInput.value.trim();

    if (!pageUrl || state.busy) {
        return;
    }

    state.controller = new AbortController();
    elements.results.hidden = true;
    elements.loading.hidden = false;
    setBusy(ctx, true);

    try {
        const response = await fetch(options.fetchUrl, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ url: pageUrl }),
            signal: state.controller.signal,
        });

        if (!response.ok) {
            throw new Error(`Server responded with ${response.status}`);
        }

        const data = await response.json();
        state.items = Array.isArray(data.items) ? data.items : [];

        if (state.items.length === 0) {
            console.error("No media found for URL:", pageUrl);
            notify.show("No media found on this page", notify.notifyType.ERROR);
        }

        renderItems(ctx);
    } catch (error) {
        if (error.name !== "AbortError") {
            console.error("Failed to fetch media list:", error);
            notify.show(error.message, notify.notifyType.ERROR);
        }
    } finally {
        elements.loading.hidden = true;
        setBusy(ctx, false);
    }
}

async function submitSelected(ctx) {
    const { dialog, elements, state, options } = ctx;
    const selected = getSelectedItems(ctx);

    if (selected.length === 0 || state.busy) {
        return;
    }

    setBusy(ctx, true);
    setSubmitting(elements, true);

    try {
        const response = await fetch(options.submitUrl, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ items: selected }),
        });
        const data = await response.json().catch(() => ({}));

        if (!response.ok || data.success === false) {
            throw new Error(data.error || `Server responded with ${response.status}`);
        }

        resetDialog(ctx);

        if (typeof options.onSuccess === "function") {
            options.onSuccess();
        } else if (typeof dialog.close === "function") {
            dialog.close();
        }
    } catch (error) {
        // Keep the dialog open so the user can retry
        console.error("Failed to submit selected media:", error);
        notify.show(error.message, notify.notifyType.ERROR);
    } finally {
        setSubmitting(elements, false);
        setBusy(ctx, false);
    }
}

/* ------------------------------------------------------------------ */
/* Events                                                              */
/* ------------------------------------------------------------------ */

function onUrlKeydown(ctx, event) {
    if (event.key === "Enter") {
        event.preventDefault();
        fetchList(ctx);
    }
}

function onSelectAllChange(ctx) {
    const { elements } = ctx;

    getItemCheckboxes(elements).forEach((checkbox) => {
        checkbox.checked = elements.selectAll.checked;
    });
    updateControls(ctx);
}

function onListChange(ctx, event) {
    if (event.target.matches(`.${CLASS_NAMES.itemCheckbox}`)) {
        updateControls(ctx);
    }
}

function bindEvents(ctx) {
    const { dialog, elements } = ctx;

    elements.urlInput.addEventListener("input", () => updateControls(ctx));
    elements.urlInput.addEventListener("keydown", (event) => onUrlKeydown(ctx, event));
    elements.fetchButton.addEventListener("click", () => fetchList(ctx));
    elements.submitButton.addEventListener("click", () => submitSelected(ctx));
    elements.selectAll.addEventListener("change", () => onSelectAllChange(ctx));
    elements.list.addEventListener("change", (event) => onListChange(ctx, event));

    // Native <dialog> close event: abort pending request and clear the state
    dialog.addEventListener("close", () => resetDialog(ctx));
}

//* ------------------------------------------------------------------ */
/* Dialog                                                              */
/* ------------------------------------------------------------------ */
/**
 * Download list dialog.
 * @param {Object} options
 * @param {string} options.fetchUrl - endpoint that returns the media list
 * @param {string} options.submitUrl - endpoint that receives the selected items
 * @param {Function} [options.onSuccess] - called after a positive submit status
 */
export function initDownloadListDialog(options) {
    const dialog = document.getElementById(DOM_IDS.dialog);

    if (!dialog) {
        return;
    }

    const elements = getDomElements(dialog);

    if (!elements) {
        return;
    }

    const ctx = {
        dialog,
        elements,
        options,
        state: {
            items: [],
            busy: false,
            controller: null,
        },
    };

    bindEvents(ctx);
    updateControls(ctx);
}