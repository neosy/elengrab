import * as constants from './constants.js';
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