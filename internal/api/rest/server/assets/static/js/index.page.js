import { DOM_ELEMENTS, DOM_CLASSES, initDomElements, DOM_SELECTORS } from "./index.dom.js";
import * as utils from './utils.js';
import * as cookie from './cookie.js';
import * as browser from './browser.js';
import * as actionButton from './action-buttons.js';
import * as rowEventHandlers from './index.sse.events.js';
import { initPlayer } from './player.js';
import * as tooltip from './tooltip.js';
import { initIndexMenus as initMenu } from './index.menu-configs.js';
import { SELECT_NAMES, COOKIE_NAMES } from './constants.js';
import * as notify from './notifications.js';
import * as view from './index.view.js';
import * as videoPreview from './video-preview.js';

// Global variables
let globalEventSource = null;

// -------------------------------------------------------------
// Function: setupQualityFormatLogic
// Handles enabling/disabling format options based on quality
// -------------------------------------------------------------
function setupQualityFormatLogic() {
    const qualityCodecSelect = utils.getSelectByName(SELECT_NAMES.qualityCodec);
    const qualityResolutionSelect = utils.getSelectByName(SELECT_NAMES.qualityResolution);
    const formatSelect = utils.getSelectByName(SELECT_NAMES.format);

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
        globalEventSource.addEventListener("row-patch", rowEventHandlers.handleRowPatch);
        globalEventSource.addEventListener("row-delete", rowEventHandlers.handleRowDelete);
        globalEventSource.addEventListener("row-patch-field", rowEventHandlers.handleRowPatchField);
        globalEventSource.addEventListener("row-start-refreshing", rowEventHandlers.handleRowStartRefreshing);
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

function initHeaderAutoHide() {
    let lastScrollTop = 0;
    let ticking = false;
    const header = document.getElementById('header');
    const hiddenClass = "header--hidden";

    function updateHeader() {
        const scrollTop = window.scrollY || document.documentElement.scrollTop;

        if (scrollTop > lastScrollTop && scrollTop > 250) {
            header.classList.add(hiddenClass);
        } else {
            header.classList.remove(hiddenClass);
        }

        lastScrollTop = Math.max(scrollTop, 0);
        ticking = false;
    }

    function handleScroll() {
        if (!ticking) {
            requestAnimationFrame(updateHeader);
            ticking = true;
        }
    }

    window.addEventListener('scroll', handleScroll, { passive: true });
}

// -------------------------------------------------------------
// Main Init
// -------------------------------------------------------------

// Disable browser scroll position restoration
if ('scrollRestoration' in history) {
    history.scrollRestoration = 'manual';
}

// Apply persisted grid/list layout state on initial page load 
view.initGridView();
document.body.classList.add('layout-ready');

document.addEventListener('DOMContentLoaded', () => {
    initDomElements();

    const grabForm = DOM_ELEMENTS.grabForm;
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
    if (grabURLInput) {
        grabURLInput.addEventListener('keydown', (event) => {
            if (event.key === 'Enter') {
                event.preventDefault();
                buttonGrab.click();
            }
        });
    }

    // Clear before HTMX request
    if (grabForm) {
        htmx.on(grabForm, 'htmx:beforeRequest', () => {
            if (grabURLInput) {
                grabURLInput.value = '';
                // update action button after clearing
                actionButton.updateInputPasteClearButton(grabURLInput, grabInputActionBtn);
            }
            if (DOM_ELEMENTS.resultInfo) DOM_ELEMENTS.resultInfo.classList.remove("show");
        });

        // Handle HTMX response for grab form
        htmx.on(grabForm, 'htmx:afterOnLoad', (event) => {
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
                        }
                    } catch (e) {
                        // ignore non-JSON
                    }

                    /// TODO Duplicate message output has been disabled.
                    // notify.showErrorMessage(
                    //     text,
                    //     DOM_ELEMENTS.resultInfo,
                    //     DOM_ELEMENTS.resultInfoFailed
                    // );
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
    }
    

    // Guest session created
    htmx.on("guestCreated", (event) => {
        if (event.detail.value === true) {
            // Reload to apply new session (cookie)
            window.location.reload();
        }
    });

    document.body.addEventListener('htmx:responseError', function (event) {
        const xhr = event.detail.xhr;

        // --- Error handling (HTTP >= 400, except 503) ---
        if (xhr.status >= 400 && xhr.status !== 503) {
            try {
                const data = JSON.parse(xhr.responseText);
                if (data && typeof data === "object" && "message" in data) {
                    notify.show(data.message, notify.notifyType.ERROR);
                }
            } catch (e) {
                // ignore non-JSON
            }
        }
    });

    // Init quality/format sync
    setupQualityFormatLogic();

    // Initialize viewport height sync (fixes mobile PWA viewport issues)
    browser.initViewportHeightVar();

    // Initialize header auto-hide on scroll
    initHeaderAutoHide();

    // Init tooltips
    tooltip.initTooltips();

    // Init menu
    initMenu();
    
    // Init inline media player
    initPlayer();

    // Init settiongs action button
    actionButton.initInputSettingsButton(DOM_ELEMENTS.inputActionSettingsBtn, DOM_ELEMENTS.grabOptionsCollapse, DOM_ELEMENTS.grabOptions);

    // Init action button for input field
    actionButton.initInputPasteClearButton(grabURLInput, grabInputActionBtn);

    // Init action button for search input
    const searchInputClearButton = actionButton.initInputClearButton(DOM_ELEMENTS.historySearchInputWrapper, DOM_ELEMENTS.historySearchClearButton);

    // Init search elements
    view.initSearching(searchInputClearButton.clear);

    // Init header user menu elements
    view.initHeaderUserMenu();

    // Initialize video preview player
    videoPreview.initVideoPreview();
    videoPreview.initVideoPreviewHover(
        DOM_ELEMENTS.result,
        DOM_CLASSES.mediaResultRow, DOM_CLASSES.mediaResultRowThumbnailImageWrapper
    );
    videoPreview.initVideoPreviewScroll(
        DOM_ELEMENTS.result,
        DOM_CLASSES.mediaResultRow, DOM_CLASSES.mediaResultRowThumbnailImageWrapper
    );

    // Lazy-load video thumbnails.
    view.initLazyImages({
        containerSelector: DOM_SELECTORS.mediaResultThumbnailPlayButton,
        placeholderSelector: DOM_SELECTORS.mediaResultThumbnailPlaceholder,
    });

    // Lazy-load channel avatars.
    view.initLazyImages({
        containerSelector: DOM_SELECTORS.mediaResultAvatar,
    });

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
