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

/**
 * Update button state for paste/clear behavior.
 * @param {HTMLInputElement} input
 * @param {HTMLElement} btn
 */
export async function updateInputPasteClearButton(input, btn) {
    if (!input || !btn) return;

    const showPaste = true; //!input.value.trim();

    if (showPaste) {
        // Paste button only if Clipboard API is available
        // if (!navigator.clipboard || !navigator.clipboard.readText) {
        //     btn.style.display = 'none';
        //     btn.dataset.state = 'none';
        //     return;
        // } else {
        //     btn.style.display = '';
        // }

        setIcon(btn, ICON_PASTE, 'Paste from clipboard');
        btn.dataset.state = 'paste';

        btn.onclick = async () => {
            try {
                const text = await navigator.clipboard.readText();
                input.value = text;
                input.focus();
                updateInputPasteClearButton(input, btn);
            } catch (e) {
                console.error('Clipboard read failed', e);
                btn.style.display = 'none';
                btn.dataset.state = 'none';
            }
        };
    } else {
        // Clear button
        btn.style.display = '';
        setIcon(btn, ICON_CLEAR, 'Clear input');
        btn.dataset.state = 'clear';

        btn.onclick = () => {
            input.value = '';
            input.focus();
            updateInputPasteClearButton(input, btn);
        };
    }
}

/**
 * Initialize action button for the specified input.
 * @param {HTMLInputElement} input
 * @param {HTMLElement} btn
 */
export function initInputPasteClearButton(input, btn) {
    if (!input || !btn) return;

    // Initial setup
    updateInputPasteClearButton(input, btn);

    // Attach input listener
    input.addEventListener('input', () => {
        updateInputPasteClearButton(input, btn);
    });
}

// -------------------------------------------------------------
// Function: initInputClearButton
// Adds a clear (×) button to any input field inside a wrapper
// wrapperSelector — CSS selector of the container that contains the input + button
// -------------------------------------------------------------
export function initInputClearButton(wrapperSelector) {
    const wrappers = document.querySelectorAll(wrapperSelector);
    if (!wrappers.length) return;

    wrappers.forEach(wrapper => {
        const input = wrapper.querySelector('input');
        const clearBtn = wrapper.querySelector('button.history-search__clear-button');

        if (!input || !clearBtn) return;

        const updateState = () => {
            const hasValue = input.value.trim() !== '';
            wrapper.classList.toggle('has-value', hasValue);
            input.classList.toggle('has-value', hasValue);
        };

        input.addEventListener('input', updateState);

        clearBtn.addEventListener('click', () => {
            input.value = '';
            input.focus();
            updateState();
            input.dispatchEvent(new Event('input', { bubbles: true }));
        });

        // Initialization during loading (for example, auto-completion)
        updateState();
    });
}