import * as constants from './index-constants.js';
import storageState from './storage-state.js';

export function applyGridView(isGridView) {
  document.body.classList.toggle(constants.CLASS_NAMES.gridView, isGridView);
  document.body.classList.toggle(constants.CLASS_NAMES.listView, !isGridView);
}

export function toggleGridView() {
  const current = storageState.get(constants.STORAGE_KEYS.settingsGridView, true);
  const next = !current;

  storageState.set(constants.STORAGE_KEYS.settingsGridView, next);

  applyGridView(next);

  return next;
}

export function initGridView() {
  const isGridView = storageState.get(constants.STORAGE_KEYS.settingsGridView, true);

  applyGridView(isGridView);

  return isGridView;
}

export function getGridView() {
  return storageState.get(constants.STORAGE_KEYS.settingsGridView, true);
}

export function initSearching(clear) {
    const searchBtn = document.getElementById("userMenuSearchButton");
    const backBtn = document.getElementById("historySearchBackButton");
    const header = document.getElementById("header");
    const searchInput = document.getElementById("historySearchInput");

    if (!searchBtn || !header || !backBtn) return;

    searchBtn.addEventListener('click', () => {
        openSearching(header, searchInput);
    });    

    backBtn.addEventListener('click', () => {
        closeSearching(header);
    });
}

function openSearching(header, input) {
    header.classList.toggle(constants.CLASS_NAMES.isSearch, true);
    input.focus();
}

function closeSearching(header) {
    header.classList.toggle(constants.CLASS_NAMES.isSearch, false);
    clear();
}
