import * as utils from './utils.js';
import * as cookie from './cookie.js';
import * as browser from './browser.js';
import * as actionButton from './action-buttons.js';
import { initPlayer } from './player.js';
import * as tooltip from './tooltip.js';
import * as videoPreview from './video-preview.js';
import * as common from "./common.js";
import * as sseClient from "./sse.js";

import * as rowEventHandlers from './pages-list.sse.events.js';
import { initIndexMenus as initMenu } from './pages-list.menu-configs.js';
import * as view from './pages-list.view.js';
import * as search from './pages-list.search.js';
import {  CLASS_NAMES, CLASS_SELECTORS, DOM_IDS, DOM_ELEMENTS, initDomElements } from "./index-page.dom.js";

// -------------------------------------------------------------
// Element selectors and cookie names
// -------------------------------------------------------------
export const SELECT_NAMES = {
    qualityCodec: "quality-codec",
    qualityResolution: "quality-resolution",
    format: "format"
};
export const COOKIE_NAMES = {
    qualityCodec: "select_quality_codec",
    qualityResolution: "select_quality_resolution",
    format: "select_format"
};

// Global variables
let searchInputClearButton = null;

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
        cookie.saveAllSelectsToCookie(SELECT_NAMES, COOKIE_NAMES);
    };

    updateFormatOptions();

    qualityCodecSelect.addEventListener("change", updateFormatOptions);
    formatSelect.addEventListener("change", updateFormatOptions);
}

// -------------------------------------------------------------
// Main Init
// -------------------------------------------------------------

// Disable browser scroll position restoration
if ('scrollRestoration' in history) {
    history.scrollRestoration = 'manual';
}

document.addEventListener('DOMContentLoaded', () => {
    initDomElements();

    const grabForm = DOM_ELEMENTS.grabForm;
    const buttonGrab = document.querySelector('.grab-area__submit-button');

    const grabURLInput = DOM_ELEMENTS.mediaURLInput;
    const grabInputActionBtn = DOM_ELEMENTS.inputActionBtn;

    cookie.setupCookieSelectSync(SELECT_NAMES.qualityCodec, COOKIE_NAMES.qualityCodec);
    cookie.setupCookieSelectSync(SELECT_NAMES.qualityResolution, COOKIE_NAMES.qualityResolution, true);
    cookie.setupCookieSelectSync(SELECT_NAMES.format, COOKIE_NAMES.format);

    // Apply persisted grid/list layout state on initial page load 
    view.initGridView();
    document.body.classList.add('layout-ready');

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
    
    // Initialize common page functionality
    common.initHTMX();

    // Initialize viewport height sync (fixes mobile PWA viewport issues)
    browser.initViewportHeightVar();

    // Initialize header auto-hide on scroll
    common.initHeaderAutoHide();

    // Init quality/format sync
    setupQualityFormatLogic();

    // Init tooltips
    tooltip.initTooltips();

    // Init menu
    initMenu();
    
    // Init inline media player
    initPlayer(DOM_ELEMENTS.mediaPlayer);

    // Init settiongs action button
    actionButton.initInputSettingsButton(DOM_ELEMENTS.inputActionSettingsBtn, DOM_ELEMENTS.grabOptionsCollapse, DOM_ELEMENTS.grabOptions);

    // Init action button for input field
    actionButton.initInputPasteClearButton(grabURLInput, grabInputActionBtn);

    // Init search elements
    const searching = search.initSearching();
    searchInputClearButton = searching.clearButton;

    // Init header user menu elements
    view.initHeaderUserMenu();

    // Initialize video preview player
    videoPreview.initVideoPreview();
    videoPreview.initVideoPreviewHover(
        DOM_ELEMENTS.result,
        CLASS_NAMES.mediaResultRow, CLASS_NAMES.mediaResultRowThumbnailImageWrapper
    );
    const refreshVideoPreview = videoPreview.initVideoPreviewScroll(
        DOM_ELEMENTS.result,
        CLASS_NAMES.mediaResultRow, CLASS_NAMES.mediaResultRowThumbnailImageWrapper
    );

    // Lazy-load video thumbnails.
    const thumbnailLazyImages = view.initLazyImages({
        containerSelector: CLASS_SELECTORS.mediaResultThumbnailPlayButton,
        placeholderSelector: CLASS_SELECTORS.mediaResultThumbnailPlaceholder,
    });

    // Lazy-load channel avatars.
    const avatarLazyImages = view.initLazyImages({
        containerSelector: CLASS_SELECTORS.mediaResultAvatar,
    });

    // Initialize view mode bar
    view.initViewModeBar({
        getSearchQueryValues: search.getSearchQueryValues,
        onSuccess: (rows) => {
            refreshVideoPreview(true);
            search.setSearchParametersInUrl();

            const lazyObservers= [
                thumbnailLazyImages,
                avatarLazyImages,
            ];

            for (const observer of lazyObservers) {
                observer.observe(rows);
            }
        },
    });

    // Initialize the external channel link button
    view.initExtChannelLink({
        getSearchQueryValues: search.getSearchQueryValues,
        onSuccess: (rows) => {
            refreshVideoPreview(true);

            searchInputClearButton.clearInputOnly();
            search.setSearchParametersInUrl();

            const lazyObservers= [
                thumbnailLazyImages,
                avatarLazyImages,
            ];

            for (const observer of lazyObservers) {
                observer.observe(rows);
            }
        },
    });

    // Create SSE connection
    const sseEventHandlers = {
        "row-add": rowEventHandlers.handleRowAdd,
        "row-update": rowEventHandlers.handleRowUpdate,
        "row-patch": rowEventHandlers.handleRowPatch,
        "row-delete": rowEventHandlers.handleRowDelete,
        "row-patch-field": rowEventHandlers.handleRowPatchField,
        "row-start-refreshing": rowEventHandlers.handleRowStartRefreshing,
        "system-info-update": rowEventHandlers.handleSystemInfoUpdate,
        "notification": rowEventHandlers.handleNotification,
    };    
    sseClient.initSSE(sseEventHandlers);
});

// Force scroll to top after full page load
window.addEventListener('load', () => {
    window.scrollTo(0, 0);
});
