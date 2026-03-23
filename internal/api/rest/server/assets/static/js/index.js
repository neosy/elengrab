import * as helper from './helper.js';
import * as cookie from './cookie.js';
import * as actionButton from './action-button.js';
import * as rowEventHandlers from './row-event-handlers.js';
import * as player from './player.js';
import * as tooltip from './tooltip.js';
import { DOM_IDS } from "./dom-ids.js";
import { DOM_ELEMENTS } from "./dom-elements.js";
import { SELECT_NAMES, COOKIE_NAMES } from './constants.js';

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
    let eventSource;

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
        eventSource = new EventSource("/ui/downloader/files/events");

        // Server is considered online when these events arrive
        eventSource.addEventListener("connected", () => setServerStatus(true));
        eventSource.addEventListener("ping", () => setServerStatus(true));

        // Business events
        eventSource.addEventListener("row-add", rowEventHandlers.handleRowAdd);
        eventSource.addEventListener("row-update", rowEventHandlers.handleRowUpdate);
        eventSource.addEventListener("row-delete", rowEventHandlers.handleRowDelete);
        eventSource.addEventListener("row-patch-field", rowEventHandlers.handleRowPatchField);
        eventSource.addEventListener("system-info-update", rowEventHandlers.handleSystemInfoUpdate);
        eventSource.addEventListener("notification", rowEventHandlers.handleNotification);

        // Fallback: any default message marks server as online
        eventSource.onmessage = () => setServerStatus(true);

        // On error: mark offline and reconnect
        eventSource.onerror = function(err) {
            console.error("SSE connection lost:", err);
            setServerStatus(false);

            eventSource.close();

            // Reconnect after delay
            setTimeout(connect, 5000);
        };
    }

    // Initial connection
    connect();

    // Return only API to close connection from outside
    return {
        close: () => eventSource?.close()
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
    const formGrab = document.querySelector(`#${DOM_IDS.grabForm}`);
    const buttonGrab = document.querySelector('.button-grab-get');

    const grabInputURL = DOM_ELEMENTS.mediaURL;
    const grabInputActionBtn = DOM_ELEMENTS.inputActionBtn;

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
        if (event.detail.elt === formGrab) {
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
        }
    });

    // Init quality/format sync
    setupQualityFormatLogic();

    // Init tooltips
    tooltip.initTooltips();
    
    // Init inline media player
    player.initPlayer();

    // Init action button for input field
    actionButton.initInputPasteClearButton(grabInputURL, grabInputActionBtn)

    // Init action button for input
    actionButton.initInputClearButton('.history-search__wrapper');

    // Create SSE connection
    const sse = createSSEConnection();

    // Close SSE on page unload
    window.addEventListener("beforeunload", () => {
        // Disabled: fires on file downloads too, causing unwanted SSE disconnects
        // console.warn("SSE connection closed");
        // sse.close();
    });
});

// Force scroll to top after full page load
window.addEventListener('load', () => {
    window.scrollTo(0, 0);
});
