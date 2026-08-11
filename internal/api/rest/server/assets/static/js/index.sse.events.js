import { DOM_CLASS_PREFIXES, DOM_CLASSES, DOM_SELECTORS, DOM_CSS_VAR_NAMES, DOM_ELEMENTS } from "./index.dom.js";
import * as notify from './notifications.js';
import { DOM_IDS, VIDEO_PREVIEW } from './constants.js';

// -------------------------------------------------------------
// Handle row-add SSE event with multiple rows in one payload
// and insert them after the top placeholder div
// -------------------------------------------------------------
export function handleRowAdd(event) {
    try {
        // Payload is raw HTML with multiple rows
        const data = JSON.parse(event.data);
        if (!data.itemId || !data.html) return

        // Find the top placeholder div
        const placeholder = document.getElementById(DOM_IDS.rowTopPlaceholder);
        if (!placeholder) return;

        // Insert all rows right after the placeholder
        placeholder.insertAdjacentHTML("afterend", data.html);

        const newEl = document.getElementById(DOM_IDS.row(data.itemId));
        if (newEl) {
            // If the element exists, remove it from the DOM
            const elNoItems = document.getElementById(DOM_IDS.rowNoItems);;
            if (elNoItems) {
                elNoItems.remove();
            }

            // Let HTMX process this element and activate any hx-* attributes
            htmx.process(newEl);
            
            // Add animation class and initial styles
            newEl.classList.add('animated');
            newEl.style.opacity = 0;
            newEl.style.transform = 'translateY(-20px)';
            newEl.style.transition = 'opacity 0.75s ease, transform 0.75s ease';

            // Trigger animation on the next browser frame
            requestAnimationFrame(() => {
                newEl.style.opacity = 1;
                newEl.style.transform = 'translateY(0)';
            });
        }
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

        if (el.classList.contains(VIDEO_PREVIEW.previewPlayingClassName)) {
            return;
        }

        const temp = document.createElement("div");
        temp.innerHTML = data.html.trim();

        const newEl = temp.firstElementChild;
        if (!newEl) return;

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
            const titleEl = el.querySelector(DOM_SELECTORS.mediaResultTitleLink);
            if (titleEl) {
                titleEl.textContent = data.title;
                titleEl.title = data.title;
            }
        }

        if (data.visibility !== undefined) {
            const visibilityEl = el.querySelector(DOM_SELECTORS.mediaResultVisibility);

            if (visibilityEl) {
                const oldClass = [...visibilityEl.classList]
                    .find(className => className.startsWith(DOM_CLASS_PREFIXES.visibility));

                if (oldClass) {
                    visibilityEl.classList.remove(oldClass);
                }

                visibilityEl.classList.toggle("hidden", !data.visibility.visible);
                visibilityEl.classList.add(`${DOM_CLASS_PREFIXES.visibility}${data.visibility.value}`);
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
            const linkEl = el.querySelector(DOM_SELECTORS.mediaResultShareLink);
            if (linkEl) {
                linkEl.classList.toggle("hidden", !data.hasShareLink);
            }
        }

        if (data.watched !== undefined && data.watched) {
            const watchedEl = el.querySelector(DOM_SELECTORS.mediaResultRowThumbnailWatched);
            if (watchedEl) {
                watchedEl.classList.remove("hidden");
            }
        }

        if (data.watchPercent !== undefined) {
            const progressEl = el.querySelector(DOM_SELECTORS.mediaResultRowThumbnailWatchProgress);
            const progressValueEl = el.querySelector(DOM_SELECTORS.mediaResultRowThumbnailWatchProgressValue);
            if (progressEl && progressValueEl) {
                const watchPercent = Number(data.watchPercent);
                progressEl.classList.toggle("hidden", watchPercent === 0);
                progressValueEl.style.setProperty(
                    DOM_CSS_VAR_NAMES.watchProgress,
                    `${watchPercent}%`,
                );
            }
        }

        if (data.viewCount !== undefined) {
            const viewCountEl = el.querySelector(DOM_SELECTORS.mediaResultViewCount);
            if (viewCountEl) {
                viewCountEl.textContent = data.viewCount.text;
            }
        }

        if (el.classList.contains(VIDEO_PREVIEW.previewPlayingClassName)) {
            return;
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

        el.classList.add(DOM_CLASSES.rowRefreshing);
    } catch (err) {
        console.error(`SSE "${event.type}" handler error:`, err);
    }
}

// -------------------------------------------------------------
// Handle row-delete SSE event and remove element if it exists
// -------------------------------------------------------------
export function handleRowDelete(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.itemId) return;

        const el = document.getElementById(DOM_IDS.row(data.itemId));
        if (!el) return;

        // Fix current height to enable smooth animation
        el.style.height = `${el.offsetHeight}px`;

        // Force layout recalculation
        void el.offsetHeight;

        // Start animation by adding CSS class
        el.classList.add(DOM_CLASSES.rowRemoving);

        // Remove element after transition completes
        el.addEventListener('transitionend', () => {
            el.remove();
        }, { once: true });

    } catch (err) {
        console.error(`SSE "${event.type}" handler error:`, err);
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


export function handleNotification(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.module || !data.type || !data.message) return

        notify.show(data.message, data.type)
    } catch (err) {
        console.error(`SSE ${event.type} handler error:`, err);
    }
}