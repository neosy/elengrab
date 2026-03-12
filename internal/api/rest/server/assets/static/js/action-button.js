import { ICON_PASTE, ICON_CLEAR } from './constants.js';

// -------------------------------------------------------------
// Set SVG icon for the action button
// -------------------------------------------------------------
const svgCache = {}; //cache for loaded SVGs

/**
 * Set SVG icon on the given button
 * @param {HTMLElement} btn - the button element
 * @param {string} url - path to the SVG file
 * @param {string} alt - text for aria-label
 */
function setIcon(btn, url, alt) {
    if (!btn) return;

    // use cached SVG if available
    if (svgCache[url]) {
        btn.innerHTML = svgCache[url];
        btn.setAttribute('aria-label', alt);
        btn.setAttribute('title', alt);
        return;
    }

    // fetch SVG from the server
    fetch(url)
        .then(res => {
            if (!res.ok) throw new Error(`Failed to load SVG: ${url}`);
            return res.text();
        })
        .then(svgText => {
            svgCache[url] = svgText; // cache the SVG
            btn.innerHTML = svgText;
            btn.setAttribute('aria-label', alt);
            btn.setAttribute('title', alt);
        })
        .catch(err => console.error(err));
}

// -------------------------------------------------------------
// Update action button based on input field state
// -------------------------------------------------------------
export function updateActionButton() {
    const input = document.getElementById('mediaURL');
    const btn = document.getElementById('inputActionBtn');

    if (!input || !btn) return;

    const showPaste = !input.value.trim();

    // Skip if the button is already in correct state
    if (showPaste && btn.dataset.state === 'paste') return;
    if (!showPaste && btn.dataset.state === 'clear') return;

    if (showPaste) {
        // Paste button only if Clipboard API is available
        if (!navigator.clipboard || !navigator.clipboard.readText) {
            btn.style.display = 'none'; // hide paste button if clipboard not available
            btn.dataset.state = 'none';
            return;
        } else {
            btn.style.display = ''; // show paste button
        }

        setIcon(btn, ICON_PASTE, 'Paste from clipboard');
        btn.dataset.state = 'paste';

        btn.onclick = async () => {
            try {
                const text = await navigator.clipboard.readText();
                input.value = text;
                updateActionButton();
            } catch (e) {
                console.error('Clipboard read failed', e);
                btn.style.display = 'none'; // hide button if reading fails
                btn.dataset.state = 'none';
            }
        };
    } else {
        // Clear button always visible
        btn.style.display = '';
        setIcon(btn, ICON_CLEAR, 'Clear input');
        btn.dataset.state = 'clear';

        btn.onclick = () => {
            input.value = '';
            updateActionButton();
        };
    }
}

// -------------------------------------------------------------
// Handle row-add SSE event with multiple rows in one payload
// and insert them after the top placeholder div
// -------------------------------------------------------------
export function handleRowAdd(event) {
    try {
        // Payload is raw HTML with multiple rows
        const data = JSON.parse(event.data);
        if (!data.id || !data.html) return

        // Find the container that holds all rows
        const container = document.getElementById("grab-result-rows");
        if (!container) return;

        // Find the top placeholder div
        const placeholder = document.getElementById("row-top-placeholder");
        if (!placeholder) return;

        // Insert all rows right after the placeholder
        placeholder.insertAdjacentHTML("afterend", data.html);

        const newEl = document.getElementById(data.id);
        if (newEl) {
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
        console.error("SSE row-add handler error:", err);
    }
}

export function handleRowUpdate(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.id || !data.html) return

        const el = document.getElementById(data.id);
        if (!el) return;

        const temp = document.createElement("div");
        temp.innerHTML = data.html.trim();

        const newEl = temp.firstElementChild;
        if (!newEl) return;

        el.replaceWith(newEl);
        //HTMX will scan this element and activate all hx-* attributes
        htmx.process(newEl);
    } catch (err) {
        console.error("SSE row-update handler error:", err);
    }
}

// -------------------------------------------------------------
// Handle row-delete SSE event and remove element if it exists
// -------------------------------------------------------------
export function handleRowDelete(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.id) return;

        const el = document.getElementById(data.id);
        if (!el) return;

        // initial styles for animation
        el.style.transition = "opacity 0.4s ease, transform 0.4s ease, height 0.4s ease, margin 0.4s ease";
        el.style.opacity = "1";
        el.style.transform = "translateY(0)";
        el.style.height = el.offsetHeight + "px";

        // forcing layout so that the browser applies height
        el.offsetHeight;

        // final state
        el.style.opacity = "0";
        el.style.transform = "translateY(-10px)";
        el.style.height = "0";
        el.style.marginTop = "0";
        el.style.marginBottom = "0";
        el.style.paddingTop = "0";
        el.style.paddingBottom = "0";

        el.addEventListener("transitionend", () => {
            el.remove();
        }, { once: true });
    } catch (err) {
        console.error("SSE row-delete handler error:", err);
    }
}