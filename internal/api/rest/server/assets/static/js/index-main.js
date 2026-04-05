import * as helper from './helper.js';
import * as cookie from './cookie.js';
import * as actionButton from './action-buttons.js';
import * as rowEventHandlers from './index-row-event-handlers.js';
import { initPlayer } from './index-player.js';
import { initTooltips } from './tooltip.js';
import { initIndexMenus as initMenu } from './menu-configs.js';
import { DOM_IDS } from "./index-dom-ids.js";
import { DOM_ELEMENTS } from "./index-dom-elements.js";
import { SELECT_NAMES, COOKIE_NAMES } from './constants.js';

// Global variables
let globalEventSource = null;

// -------------------------------------------------------------
// Function: setupQualityFormatLogic
// Handles enabling/disabling format options based on quality
// -------------------------------------------------------------
function setupQualityFormatLogic() {
    const qualityCodecSelect = helper.getSelectByName(SELECT_NAMES.qualityCodec);
    const qualityResolutionSelect = helper.getSelectByName(SELECT_NAMES.qualityResolution);
    const formatSelect = helper.getSelectByName(SELECT_NAMES.format);

    if (!qualityCodecSelect || !qualityResolutionSelect || !formatSelect) return;

    const videoFormats = ["auto", "mp4", "webm"];
    const audioFormats = ["auto", "mp3", "m4a", "flac", "opus"];
    const onlyAudioFormats = ["mp3", "m4a", "flac", "opus"];
    const webmCodecs = ["best", "av1"];
    const videoFormatWebmValue = "webm";
    const bestValue = "best"
    const maxValue = "max"
    const autoValue = "auto"
    const emptyValue = "empty"
    const onlyAudioValue = "only_audio";
    const videoFormatDefault = autoValue;
    const audioFormatDefault = autoValue;

    // Update format options based on quality
    const updateFormatOptions = () => {
        const qualityCodec = qualityCodecSelect.value;
        const isQualityCodecBest = (qualityCodec === bestValue);
        const isCodecOnlyAudio = (qualityCodec === onlyAudioValue);
        const isCodecFormatOnlyAudio = isCodecOnlyAudio || (isQualityCodecBest && onlyAudioFormats.includes(formatSelect.value));

        if (isCodecFormatOnlyAudio && !qualityResolutionSelect.disabled) {
            qualityResolutionSelect.disabled = true;
            qualityResolutionSelect.value = emptyValue;
        }

        if (!isCodecFormatOnlyAudio && qualityResolutionSelect.disabled) {
            qualityResolutionSelect.disabled = false;
            qualityResolutionSelect.value = maxValue;
        }

        formatSelect.querySelectorAll("option").forEach(option => {
            const value = option.value;

            if (isCodecOnlyAudio) {
                const allowed = audioFormats.includes(value);
                option.disabled = !allowed;
            } else {
                let allowed = isQualityCodecBest || videoFormats.includes(value);
                if (value == videoFormatWebmValue) {
                    allowed = webmCodecs.includes(qualityCodec);
                }
                option.disabled = !allowed;
            }
        });

        if (!isQualityCodecBest) {
            if (isCodecOnlyAudio && !audioFormats.includes(formatSelect.value)) {
                formatSelect.value = audioFormatDefault;
            }

            if (!isCodecOnlyAudio && !videoFormats.includes(formatSelect.value)) {
                formatSelect.value = videoFormatDefault;
            }

            if (!isCodecOnlyAudio && formatSelect.value == videoFormatWebmValue && !webmCodecs.includes(qualityCodec)) {
                formatSelect.value = videoFormatDefault;
            }
        }
        cookie.saveAllSelectsToCookie();
    };

    updateFormatOptions();

    qualityCodecSelect.addEventListener("change", updateFormatOptions);
    formatSelect.addEventListener("change", updateFormatOptions);
}

function createSSEConnection() {
    function setServerStatus(online) {
        const el = DOM_ELEMENTS.sysInfoServerStatusDot
        if (!el) return;

        if (online) {
            el.classList.add("online");
        } else {
            el.classList.remove("online");
        }
    }

    // Internal function to (re)connect
    function connect() {
        if (globalEventSource) {
            globalEventSource.close();
        }        

        globalEventSource = new EventSource("/downloader/files/events");

        // Server is considered online when these events arrive
        globalEventSource.addEventListener("connected", () => setServerStatus(true));
        globalEventSource.addEventListener("ping", () => setServerStatus(true));

        // Business events
        globalEventSource.addEventListener("row-add", rowEventHandlers.handleRowAdd);
        globalEventSource.addEventListener("row-update", rowEventHandlers.handleRowUpdate);
        globalEventSource.addEventListener("row-delete", rowEventHandlers.handleRowDelete);
        globalEventSource.addEventListener("row-patch-field", rowEventHandlers.handleRowPatchField);
        globalEventSource.addEventListener("system-info-update", rowEventHandlers.handleSystemInfoUpdate);
        globalEventSource.addEventListener("notification", rowEventHandlers.handleNotification);

        // Fallback: any default message marks server as online
        globalEventSource.onmessage = () => setServerStatus(true);

        // On error: mark offline and reconnect
        globalEventSource.onerror = function(err) {
            console.error("SSE connection lost:", err);
            setServerStatus(false);

            globalEventSource.close();

            // Reconnect after delay
            setTimeout(connect, 5000);
        };
    }

    // Initial connection
    connect();

    // Return only API to close connection from outside
    return {
        close: () => globalEventSource?.close()
    };
}

// -------------------------------------------------------------
// Main Init
// -------------------------------------------------------------

// Disable browser scroll position restoration
if ('scrollRestoration' in history) {
    history.scrollRestoration = 'manual';
}

document.addEventListener('DOMContentLoaded', () => {
    const grabForm = document.querySelector(`#${DOM_IDS.grabForm}`);
    const buttonGrab = document.querySelector('.grab-area__submit-button');

    const grabInputURL = DOM_ELEMENTS.mediaURL;
    const grabInputActionBtn = DOM_ELEMENTS.inputActionBtn;

    let sessionToken = cookie.getSessionToken();

    // Sync selects with cookies
    cookie.setupCookieSelectSync(SELECT_NAMES.qualityCodec, COOKIE_NAMES.qualityCodec);
    cookie.setupCookieSelectSync(SELECT_NAMES.qualityResolution, COOKIE_NAMES.qualityResolution);
    cookie.setupCookieSelectSync(SELECT_NAMES.format, COOKIE_NAMES.format);

    // Submit on Enter
    grabInputURL.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            event.preventDefault();
            buttonGrab.click();
        }
    });

    // Clear before HTMX request
    htmx.on(`#${DOM_IDS.grabForm}`, 'htmx:beforeRequest', () => {
        if (grabInputURL) {
            grabInputURL.value = '';
            // update action button after clearing
            actionButton.updateInputPasteClearButton(grabInputURL, grabInputActionBtn);
        }
        if (DOM_ELEMENTS.resultInfo) DOM_ELEMENTS.resultInfo.classList.remove("show");
    });

    // Display error on non-200 + non-503
    document.body.addEventListener('htmx:afterOnLoad', (event) => {
        if (event.detail.elt === grabForm) {
            if (grabInputURL) grabInputURL.value = '';
            if (event.detail.xhr.status !== 200 &&
                event.detail.xhr.status !== 503) {
                    
                if (DOM_ELEMENTS.resultInfo && DOM_ELEMENTS.resultInfoFailed) {
                    let text = event.detail.xhr.responseText;
                    try {
                        const data = JSON.parse(text);
                        if (data && typeof data === "object" && "message" in data) {
                            text = data.message;
                        }
                    } catch (e) {}

                    helper.showErrorMessage(text, DOM_ELEMENTS.resultInfo, DOM_ELEMENTS.resultInfoFailed)
                }
            }

            // Reload the page after the user has been authorized (e.g., guest user created and session token set)
            if (!sessionToken) {
                const curToken = cookie.getSessionToken();
                if (curToken) {
                    window.location.reload();
                }
            }
        }
    });

    // Init quality/format sync
    setupQualityFormatLogic();

    // Init tooltips
    initTooltips();

    // Init menu
    initMenu();
    
    // Init inline media player
    initPlayer();

    // Init action button for input field
    actionButton.initInputPasteClearButton(grabInputURL, grabInputActionBtn)

    // Init action button for input
    actionButton.initInputClearButton('.history-search__input_wrapper');

    // Create SSE connection
    const sse = createSSEConnection();

    // Close SSE on page unload
    window.addEventListener("beforeunload", () => {
        console.warn("SSE connection closed");
        sse?.close();
    });
});

// Force scroll to top after full page load
window.addEventListener('load', () => {
    window.scrollTo(0, 0);
});
