import * as browser from './browser.js';
import * as common from "./common.js";
import * as sseClient from "./sse.js";
import * as dialog from "./dialog.js";
import * as tooltip from './tooltip.js';
import { initPlayer } from './player.js';
import * as videoPreview from './video-preview.js';

import * as rowEventHandlers from './pages-list.sse.events.js';
import { initChannelMenus as initMenu } from './pages-list.menu-configs.js';
import * as view from './pages-list.view.js';
import * as search from './pages-list.search.js';
import { CLASS_NAMES, CLASS_SELECTORS, DOM_ELEMENTS, initDomElements } from "./channel-page.dom.js";

// Global variables
let searchInputClearButton = null;

// -------------------------------------------------------------
// Main Init
// -------------------------------------------------------------

// Disable browser scroll position restoration
if ('scrollRestoration' in history) {
    history.scrollRestoration = 'manual';
}

document.addEventListener('DOMContentLoaded', () => {
    initDomElements();

    // Apply persisted grid/list layout state on initial page load 
    view.initGridView();
    document.body.classList.add('layout-ready');

    // Initialize common page functionality
    common.initHTMX();

    // Initialize viewport height sync (fixes mobile PWA viewport issues)
    browser.initViewportHeightVar();

    // Initialize header auto-hide on scroll
    common.initHeaderAutoHide();

    // Init dialogs
    dialog.initDialogs();

    // Init tooltips
    tooltip.initTooltips();

    // Init menu
    initMenu();
    
    // Init inline media player
    initPlayer(DOM_ELEMENTS.mediaPlayer);

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
        "row-update": rowEventHandlers.handleRowUpdate,
        "row-patch": rowEventHandlers.handleRowPatch,
        "row-delete": rowEventHandlers.handleRowDelete,
        "row-patch-field": rowEventHandlers.handleRowPatchField,
        "row-start-refreshing": rowEventHandlers.handleRowStartRefreshing,
        "notification": rowEventHandlers.handleNotification,
    };    
    sseClient.initSSE(sseEventHandlers);
});

// Force scroll to top after full page load
window.addEventListener('load', () => {
    window.scrollTo(0, 0);
});
