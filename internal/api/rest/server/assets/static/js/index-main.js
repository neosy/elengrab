import { DOM_IDS } from "./index-dom-ids.js";
import { DOM_ELEMENTS } from "./index-dom-elements.js";
import * as helper from './helper.js';
import * as cookie from './cookie.js';
import * as actionButton from './action-buttons.js';
import * as rowEventHandlers from './index-row-event-handlers.js';
import { initPlayer } from './index-player.js';
import { initTooltips } from './tooltip.js';
import { initIndexMenus as initMenu } from './menu-configs.js';
import { SELECT_NAMES, COOKIE_NAMES } from './constants.js';
import * as notify from './notifications.js';

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
            console.info("SSE connection opened");
            el.classList.add("online");
        } else {
            console.warn("SSE connection closed");
            el.classList.remove("online");
        }
    }

    // Internal function to (re)connect
    function connect() {
        globalEventSource?.close();

        globalEventSource = new EventSource("/downloader/events");

        // Server is considered online when these events arrive
        globalEventSource.addEventListener("connected", () => setServerStatus(true));
        // globalEventSource.addEventListener("ping", () => setServerStatus(true));

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
            if (globalEventSource && globalEventSource.readyState !== EventSource.CLOSED) {
                console.error("SSE connection lost:", err);
                setServerStatus(false);
                globalEventSource?.close();
            }

            // Reconnect after delay
            setTimeout(connect, 5000);
        };
    }

    // Initial connection
    connect();

    // Return only API to close connection from outside
    return {
        close: () => {
            setServerStatus(false);
            globalEventSource?.close();
        }};
}

function initSearching(clear) {
    const searchBtn = document.getElementById("userMenuSearchButton");
    const backBtn = document.getElementById("historySearchBackButton");
    const mainHeader = document.getElementById("mainHeader");
    const searchInput = document.getElementById("historySearchInput");

    if (!searchBtn || !mainHeader || !backBtn) return;

    searchBtn.addEventListener('click', () => {
        openSearching();
    });    

    backBtn.addEventListener('click', () => {
        closeSearching();
    });

    function openSearching() {
        mainHeader.classList.toggle('is-search', true);
        searchInput.focus();
    }

    function closeSearching() {
        mainHeader.classList.toggle('is-search', false);
        clear();
    }
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

    const grabURLInput = DOM_ELEMENTS.mediaURLInput;
    const grabInputActionBtn = DOM_ELEMENTS.inputActionBtn;

    // Sync selects with cookies
    cookie.setupCookieSelectSync(SELECT_NAMES.qualityCodec, COOKIE_NAMES.qualityCodec);
    cookie.setupCookieSelectSync(SELECT_NAMES.qualityResolution, COOKIE_NAMES.qualityResolution);
    cookie.setupCookieSelectSync(SELECT_NAMES.format, COOKIE_NAMES.format);

    // TODO Disabled. Problems resetting the position when pressing Back
    // Reload page if restored from bfcache (back/forward navigation)
    // window.addEventListener('pageshow', (event) => {
    //     if (event.persisted) {
    //         window.location.reload();
    //     }
    // });

    // Submit on Enter
    grabURLInput.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            event.preventDefault();
            buttonGrab.click();
        }
    });

    // Clear before HTMX request
    htmx.on(`#${DOM_IDS.grabForm}`, 'htmx:beforeRequest', () => {
        if (grabURLInput) {
            grabURLInput.value = '';
            // update action button after clearing
            actionButton.updateInputPasteClearButton(grabURLInput, grabInputActionBtn);
        }
        if (DOM_ELEMENTS.resultInfo) DOM_ELEMENTS.resultInfo.classList.remove("show");
    });

    // Handle HTMX response for grab form
    htmx.on(`#${DOM_IDS.grabForm}`, 'htmx:afterOnLoad', (event) => {
        const xhr = event.detail.xhr;

        // --- Error handling (HTTP >= 400, except 503) ---
        if (xhr.status >= 400 && xhr.status !== 503) {
            if (grabURLInput) grabURLInput.value = '';

            if (DOM_ELEMENTS.resultInfo && DOM_ELEMENTS.resultInfoFailed) {
                let text = xhr.responseText;

                try {
                    const data = JSON.parse(text);
                    if (data && typeof data === "object" && "message" in data) {
                        text = data.message;
                        notify.show(data.message, notify.notifyType.ERROR);
                    }
                } catch (e) {
                    // ignore non-JSON
                }

                helper.showErrorMessage(
                    text,
                    DOM_ELEMENTS.resultInfo,
                    DOM_ELEMENTS.resultInfoFailed
                );
            }

            return;
        }

        // --- Success: guest session created ---
        let data;
        try {
            data = JSON.parse(xhr.responseText);
        } catch (e) {
            console.error('Invalid JSON response');
            return;
        }

        if (data.guestCreated === true) {
            // Reload to apply new session (cookie)
            window.location.reload();
        }
    });
    

    // Guest session created
    htmx.on("guestCreated", (event) => {
        if (event.detail.value === true) {
        // Reload to apply new session (cookie)
            window.location.reload();
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

    // Init settiongs action button
    actionButton.initInputSettingsButton(DOM_ELEMENTS.inputActionSettingsBtn, DOM_ELEMENTS.grabOptionsCollapse, DOM_ELEMENTS.grabOptions)

    // Init action button for input field
    actionButton.initInputPasteClearButton(grabURLInput, grabInputActionBtn)

    // Init action button for search input
    const searchInputClearButton = actionButton.initInputClearButton(DOM_ELEMENTS.historySearchInputWrapper, DOM_ELEMENTS.historySearchClearButton);

    // Init search elements
    initSearching(searchInputClearButton.clear);

    // Create SSE connection
    var sse = null;
    window.addEventListener('pageshow', () => {
        if (!globalEventSource || globalEventSource.readyState === EventSource.CLOSED) {
            sse = createSSEConnection();
        }
    });    

    // Close SSE on page unload
    window.addEventListener("beforeunload", () => {
        if (!globalEventSource || globalEventSource.readyState === EventSource.CLOSED) {
            sse = createSSEConnection();
        }
    });

    // ------------------------------------------------------------
    // Reconnect SSE when tab becomes visible again
    // (fixes lost Server-Sent Events after browser tab sleep,
    // background throttling, or mobile suspension)
    // ------------------------------------------------------------
    document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "visible") {
        if (!globalEventSource || globalEventSource.readyState === EventSource.CLOSED) {
            sse = createSSEConnection();
        }
    }
});

});

// Force scroll to top after full page load
window.addEventListener('load', () => {
    window.scrollTo(0, 0);
});
