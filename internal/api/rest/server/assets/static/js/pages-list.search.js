import { CLASS_NAMES as PAGES_CLASS_NAMES, CLASS_SELECTORS, DOM_IDS } from "./pages-list.dom.js";
import { initInputClearButton } from './action-buttons.js';
import { isMobileScreen } from "./browser.js";

const URL_KEYS = {
    query: 'query',
    searchQuery: 'searchQuery',
    searchParams: 'sp',
};

const CLASS_NAMES = {
    ...PAGES_CLASS_NAMES,
    isSearch: "is-search",
}

const DOM_ELEMENTS = {
    searchInput: null,
}

function initDomElements() {
    DOM_ELEMENTS.searchInput = document.getElementById("historySearchInput");
}

export function initSearching() {
    const searchBtn = document.getElementById("headerActionsSearchButton");
    const backBtn = document.getElementById("historySearchBackButton");
    const header = document.getElementById("header");
    const searchInputWrapper = document.getElementById("historySearchInputWrapper");
    const searchClearButton = document.getElementById("historySearchClearButton");

    initDomElements();

    const clearButton = initInputClearButton(searchInputWrapper, searchClearButton);

    if (!searchBtn || !header || !backBtn) return;

    searchBtn.addEventListener('click', () => {
        openSearching(header, DOM_ELEMENTS.searchInput);
    });    

    if (clearButton) {
        backBtn.addEventListener('click', () => {
            closeSearching(header, clearButton.clear);
        });
    }

    if (isMobileScreen() && DOM_ELEMENTS.searchInput.value !== "") {
        header.classList.toggle(CLASS_NAMES.isSearch, true);
    }

    window.configureSearch = function (event) {
        configureSearch(event)
    }
    window.isSearchQueryValid = isSearchQueryValid
    window.handleSearchSuccess = handleSearchSuccess

    return {
        open: openSearching,
        close: closeSearching,
        clearButton: clearButton,
    }
}

export function getSearchQueryValues() {
    const viewModeTabs = document.querySelector(CLASS_SELECTORS.viewModeTabs)

    let queryValues = {
        query: DOM_ELEMENTS.searchInput?.value ?? "",
        viewMode: viewModeTabs.dataset.viewMode,
    }

    const rows = document.getElementById(DOM_IDS.mediaResultRows);
    if (!rows) {
        return queryValues;
    }

    if (rows.dataset.searchFiltersJson) {
        const searchFilters = JSON.parse(rows.dataset.searchFiltersJson);
        const channelId = searchFilters.channelId;

        if (channelId) {
            queryValues.channelId = channelId;
        }
    }

    return queryValues;
}

function getSearchValues() {
    let values = {
        query: DOM_ELEMENTS.searchInput.value,
        hasSearchFilters: false,
        encodeParams: null,
    }

    const rows = document.getElementById(DOM_IDS.mediaResultRows);
    if (rows) {
        values.hasSearchFilters = rows.dataset.hasSearchFilters === "true";
        if (rows.dataset.searchParameters) {
            values.encodeParams = rows.dataset.searchParameters; 
        }
    }

    return values;
}

function configureSearch(event) {
    Object.assign(event.detail.parameters, getSearchQueryValues())
}

function isSearchQueryValid(input) {
    return input.value.length === 0 || input.value.length >= 2
}

function hasSearchParametersForUrl() {
    const values = getSearchValues();

    if (values.query || values.hasSearchFilters) {
        return true
    }

    return false
}

export function setSearchParametersInUrl(event) {
    const values = getSearchValues();
    const url = new URL(window.location.href);

    url.searchParams.delete(URL_KEYS.searchQuery);
    url.searchParams.delete(URL_KEYS.searchParams);

    if (!hasSearchParametersForUrl()) {
        window.history.replaceState(null, '', url);
        return;
    }

    if (values.query) {
        url.searchParams.set(URL_KEYS.searchQuery, values.query);
    }

    if (values.encodeParams) {
        url.searchParams.set(URL_KEYS.searchParams, values.encodeParams);
    }

    window.history.replaceState(null, '', url);
}

function openSearching(header, input) {
    header.classList.toggle(CLASS_NAMES.isSearch, true);
    input.focus();
}

function closeSearching(header, clear) {
    header.classList.toggle(CLASS_NAMES.isSearch, false);
    clear();
}

function handleSearchSuccess(event) {
    setSearchParametersInUrl(event);
}
