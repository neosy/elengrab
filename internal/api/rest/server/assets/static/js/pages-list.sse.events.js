import { CLASS_NAMES, CLASS_SELECTORS, CSS_VARIABLE_NAMES, DOM_IDS, VIEW_MODES, CONTENT_TYPES } from "./pages-list.dom.js";
import { CLASS_PREFIXES, EVENT_NAMES } from './constants.js';

function shouldAddRow(data) {
    const currentContentType = CONTENT_TYPES.getCurrent();

    if (currentContentType) {
        if (!data.isReady && !CONTENT_TYPES.isInProgress()) return false;

        if (data.isReady) {
            switch (currentContentType) {
                case CONTENT_TYPES.all:
                    break;

                case CONTENT_TYPES.video:
                    if (data.isShort || data.isAudio) return false;
                    break;
                    
                case CONTENT_TYPES.shorts:
                    if (!data.isShort) return false;
                    break;

                case CONTENT_TYPES.audio:
                    if (!data.isAudio) return false;
                    break;

                case CONTENT_TYPES.uploads:
                    if (data.sourceType !== "upload") return false;
                    break;

                default:
                    return false;
            }
        }
    }

    if (VIEW_MODES.getCurrent() && !VIEW_MODES.isNew()) return false;

    return true;
}

// -------------------------------------------------------------
// Handle row-add SSE event with multiple rows in one payload
// and insert them after the top placeholder div
// -------------------------------------------------------------
const ADD_MS = 350;
const FADE_IN_MS = 700;

// First: snapshot of the current row positions (call before inserting the new row)
function snapshotRowPositions(itemsEl) {
    const rows = Array.from(itemsEl.querySelectorAll(CLASS_SELECTORS.mediaResultRow));
    return new Map(rows.map((r) => [r, r.getBoundingClientRect()]));
}

// Last + Invert + Play: existing rows glide to new positions, the new row fades in
function animateRowAdd(newEl, first) {
    // No animation for users who prefer reduced motion
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return;

    // Read phase: measure all new positions first
    const moves = [];
    first.forEach((firstRect, r) => {
        if (!r.isConnected) return;
        const last = r.getBoundingClientRect();
        const dx = firstRect.left - last.left;
        const dy = firstRect.top - last.top;
        if (dx || dy) moves.push({ r, dx, dy });
    });

    // Write phase: start all animations
    moves.forEach(({ r, dx, dy }) => {
        r.animate(
            [
                { transform: `translate(${dx}px, ${dy}px)` },
                { transform: 'translate(0, 0)' },
            ],
            { duration: ADD_MS, easing: EASING }
        );
    });

    // The new row fades in slower than the siblings move, after the slot starts opening
    newEl.animate(
        [
            { opacity: 0, transform: 'translateY(-12px)' },
            { opacity: 1, transform: 'translateY(0)' },
        ],
        { duration: FADE_IN_MS, delay: 120, easing: 'ease-out', fill: 'backwards' }
    );
}

export function handleRowAdd(event, usernameFilter = "") {
    try {
        // Payload is raw HTML with multiple rows
        const data = JSON.parse(event.data);
        if (!data.itemId || !data.html) return

        if (!shouldAddRow(data)) return;

        const owner = data.owner;

        if (usernameFilter && owner !== usernameFilter) {
            return;
        }

        // Find the top placeholder div
        const itemsEl = document.getElementById(DOM_IDS.mediaResultItems);
        if (!itemsEl) return;

        const rowId = DOM_IDS.row(data.itemId);

        // Check if the row already exists
        let newEl = document.getElementById(rowId);
        if (newEl) return;

        // Check whether the incoming row is older than the first existing row
        const firstRow = itemsEl.querySelector(CLASS_SELECTORS.mediaResultRow);
        if (firstRow) {
            const firstCreatedAt = Number(firstRow.dataset.createdAtUnixMs);
            const incomingCreatedAt = Number(data.createdAtUnixMs);

            if (
                Number.isFinite(firstCreatedAt) &&
                Number.isFinite(incomingCreatedAt) &&
                firstCreatedAt > incomingCreatedAt
            ) {
                return;
            }
        }        

        // Snapshot of row positions before insertion
        const first = snapshotRowPositions(itemsEl);

        // Insert all rows at the beginning
        itemsEl.insertAdjacentHTML("afterbegin", data.html);

        newEl = document.getElementById(DOM_IDS.row(data.itemId));
        if (!newEl) return;

        // If the "no items" element exists, remove it from the DOM
        const elNoItems = document.getElementById(DOM_IDS.rowNoItems);
        if (elNoItems) {
            elNoItems.remove();
        }

        // Let HTMX process this element and activate any hx-* attributes
        htmx.process(newEl);
        newEl.classList.add('animated');

        animateRowAdd(newEl, first);        
    } catch (err) {
        console.error(`SSE "${event.type}" handler error:`, err);
    }
}

export function handleRowUpdate(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.itemId || !data.html) return

        const el = document.getElementById(DOM_IDS.row(data.itemId));
        if (!el) return;

        const temp = document.createElement("div");
        temp.innerHTML = data.html.trim();

        const newEl = temp.firstElementChild;
        if (!newEl) return;

        // Stop the active video preview before replacing the media row.
        document.dispatchEvent(new CustomEvent(EVENT_NAMES.videoPreviewStop));

        el.replaceWith(newEl);

        //HTMX will scan this element and activate all hx-* attributes
        htmx.process(newEl);
    } catch (err) {
        console.error(`SSE "${event.type}" handler error:`, err);
    }
}

// Title string `json:"title,omitempty"`
// WatchPercent  string  `json:"watchPercent,omitempty"`
// Watched       *bool   `json:"watched,omitempty"`
// ViewCount     *uint32 `json:"viewCount,omitempty"`
// ViewCountText string  `json:"viewCountText,omitempty"`
export function handleRowPatch(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.itemId) return

        const el = document.getElementById(DOM_IDS.row(data.itemId));
        if (!el) return;

        if (data.title !== undefined && data.title !== "") {
            const titleEl = el.querySelector(CLASS_SELECTORS.mediaResultTitleLabel);
            if (titleEl) {
                titleEl.textContent = data.title;
                titleEl.title = data.title;
            }
        }

        if (data.visibility !== undefined) {
            const visibilityEl = el.querySelector(CLASS_SELECTORS.mediaResultVisibility);

            if (visibilityEl) {
                const oldClass = [...visibilityEl.classList]
                    .find(className => className.startsWith(CLASS_PREFIXES.visibility));

                if (oldClass) {
                    visibilityEl.classList.remove(oldClass);
                }

                visibilityEl.classList.toggle("hidden", !data.visibility.visible);
                visibilityEl.classList.add(`${CLASS_PREFIXES.visibility}${data.visibility.value}`);
                visibilityEl.dataset.tooltip = data.visibility.label;
                visibilityEl.title = data.visibility.label;
                visibilityEl.ariaLabel = data.visibility.label;
                if (data.visibility.visible) {
                    visibilityEl.innerHTML = data.visibility.icon;
                } else {
                    visibilityEl.innerHTML = "";
                }
            }
        }

        if (data.hasShareLink !== undefined) {
            const linkEl = el.querySelector(CLASS_SELECTORS.mediaResultShareLink);
            if (linkEl) {
                linkEl.classList.toggle("hidden", !data.hasShareLink);
            }
        }

        if (data.watched !== undefined && data.watched) {
            const watchedEl = el.querySelector(CLASS_SELECTORS.mediaResultRowThumbnailWatched);
            if (watchedEl) {
                watchedEl.classList.remove("hidden");
            }
        }

        if (data.watchPercent !== undefined) {
            const progressEl = el.querySelector(CLASS_SELECTORS.thumbnailWatchProgress);
            const progressValueEl = el.querySelector(CLASS_SELECTORS.thumbnailWatchProgressValue);
            if (progressEl && progressValueEl) {
                const watchPercent = Number(data.watchPercent);
                progressEl.classList.toggle(CLASS_NAMES.thumbnailWatchProgressHidden, watchPercent === 0);
                progressValueEl.style.setProperty(
                    CSS_VARIABLE_NAMES.watchProgress,
                    `${watchPercent}%`,
                );
            }
        }

        if (data.viewCount !== undefined) {
            const viewCountEl = el.querySelector(CLASS_SELECTORS.mediaResultViewCount);
            if (viewCountEl) {
                viewCountEl.textContent = data.viewCount.text;
            }
        }
    } catch (err) {
        console.error(`SSE "${event.type}" handler error:`, err);
    }
}

export function handleRowStartRefreshing(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.itemId) return

        const el = document.getElementById(DOM_IDS.row(data.itemId));
        if (!el) return;

        el.classList.add(CLASS_NAMES.row.refreshing);
    } catch (err) {
        console.error(`SSE "${event.type}" handler error:`, err);
    }
}

// -------------------------------------------------------------
// Handle row-delete SSE event and remove element if it exists
// -------------------------------------------------------------
const REMOVE_MS = 350;
const EASING = 'cubic-bezier(0.4, 0, 0.2, 1)';

function animateRemoval(el) {
    const parent = el.parentElement;
    const siblings = Array.from(parent.children).filter((c) => c !== el);

    // First: positions of all elements before removal
    const first = new Map(siblings.map((c) => [c, c.getBoundingClientRect()]));
    const rect = el.getBoundingClientRect();

    // Visual copy of the row, pinned to the viewport at the exact same place
    const ghost = el.cloneNode(true);
    ghost.removeAttribute('id');
    ghost.querySelectorAll('[id]').forEach((n) => n.removeAttribute('id'));
    ghost.classList.add(CLASS_NAMES.rowRemoving);
    Object.assign(ghost.style, {
        top: `${rect.top}px`,
        left: `${rect.left}px`,
        width: `${rect.width}px`,
        height: `${rect.height}px`,
    });

    // The original leaves the layout immediately, so siblings reflow right away
    el.remove();
    parent.appendChild(ghost);

    // Last + Invert + Play: siblings glide from old positions to new ones
    siblings.forEach((c) => {
        const last = c.getBoundingClientRect();
        const dx = first.get(c).left - last.left;
        const dy = first.get(c).top - last.top;
        if (!dx && !dy) return;

        c.animate(
            [
                { transform: `translate(${dx}px, ${dy}px)` },
                { transform: 'translate(0, 0)' },
            ],
            { duration: REMOVE_MS, easing: EASING }
        );
    });

    // The copy fades and shrinks in place while siblings are moving
    ghost.animate(
        [
            { opacity: 1, transform: 'scale(1)' },
            { opacity: 0, transform: 'scale(0.9)' },
        ],
        { duration: REMOVE_MS, easing: 'ease-in', fill: 'forwards' }
    ).finished
        .then(() => ghost.remove())
        .catch(() => ghost.remove());
}

export function handleRowDelete(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.itemId) return;

        const el = document.getElementById(DOM_IDS.row(data.itemId));
        if (!el) return;

        // No animation for users who prefer reduced motion
        if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
            el.remove();
            return;
        }

        animateRemoval(el);

    } catch (err) {
        console.error(`SSE "${event.type}" handler error:`, err);
    }
}

export function handleRowPatchField(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.field || !data.itemId || !data.value) return

        switch (data.field) {
            case "progress": handleRowPatchProgress(data.itemId, data.value);
        }
    } catch (err) {
        console.error(`SSE ${event.type} handler error:`, err);
    }
}

function handleRowPatchProgress(itemId, value) {
    if (!itemId || !value) return

    const el = document.getElementById(DOM_IDS.progress(itemId));
    if (!el) return;

    el.textContent = value
}
