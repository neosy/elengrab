/* ------------------------------------------------------------------ */
/* Helpers                                                            */
/* ------------------------------------------------------------------ */

function $select(root, name) {
    return root.querySelector(`select[name="${name}"]`);
}

/* ------------------------------------------------------------------ */
/* Storage state                                                      */
/* ------------------------------------------------------------------ */

const StorageState = (() => {
  /**
   * Save value
   * @param {string} key
   * @param {any} value
   */
  function set(key, value) {
    localStorage.setItem(key, JSON.stringify(value));
  }

  /**
   * Read value
   * @param {string} key
   * @param {any} defaultValue
   */
  function get(key, defaultValue = null) {
    const raw = localStorage.getItem(key);

    if (raw === null) return defaultValue;

    try {
      return JSON.parse(raw);
    } catch (e) {
      return defaultValue;
    }
  }

  /**
   * Remove value
   * @param {string} key
   */
  function remove(key) {
    localStorage.removeItem(key);
  }

  /**
   * Toggle boolean state
   * @param {string} key
   * @param {boolean} fallback
   */
  function toggle(key, fallback = false) {
    const current = get(key, fallback);
    const next = !current;
    set(key, next);
    return next;
  }

  
  // Function: saveAllSelects
  /**
   * Save the current value of every select to storage.
   * Property names in `names` and `keys` must match; entries missing from `keys` are skipped with a warning.
   * @param {Object.<string, string>} names - Map of logical keys to select `name` attributes (e.g. `{ format: 'format' }`).
   * @param {Object.<string, string>} keys - Map of logical keys to storage keys (e.g. `{ format: 'cookie_format' }`). Must contain every property present in `names`.
   * @returns {void}
   */
  function saveAllSelects(root, names, keys) {
      Object.entries(names).forEach(([key, name]) => {
          if (!(key in keys)) {
              console.warn(`saveAllSelects: missing storage key for "${key}"`);
              return;
          }
          const el = $select(root, name);
          if (el) set(keys[key], el.value);
      });
  }  

  // Function: syncSelect
  /**
   * Sync a select element with storage: restore its saved value and optionally persist changes.
   * @param {string} selectName - The value of the select's `name` attribute.
   * @param {string} storageKey - The storage key to read the saved value from and, when `saveOnChange` is true, to write it to.
   * @param {boolean} [saveOnChange=false] - Whether to save the select's value on change.
   * @returns {void}
   */
  function syncSelect(root, selectName, storageKey, saveOnChange = false) {
      const selectElement = $select(root, selectName);
      if (!selectElement) return;
  
      const savedValue = get(storageKey);
      if (savedValue) {
          const option = selectElement.querySelector(`option[value="${savedValue}"]`);
          if (option) option.selected = true;
      }
  
      // Save the selected value on change when enabled
      if (saveOnChange) {
          selectElement.addEventListener("change", () => {
              set(storageKey, selectElement.value);
          });
      }
  }
  
  return {
    get,
    set,
    remove,
    toggle,
    saveAllSelects,
    syncSelect,
  };
})();

export default StorageState;