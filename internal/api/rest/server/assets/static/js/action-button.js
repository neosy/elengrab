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
        })
        .catch(err => console.error(err));
}

// -------------------------------------------------------------
// Update action button based on input field state
// -------------------------------------------------------------
export function updateActionButton() {
    const input = document.getElementById('youtubeURL');
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

