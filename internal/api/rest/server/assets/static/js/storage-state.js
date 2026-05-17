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

  return {
    get,
    set,
    remove,
    toggle,
  };
})();

export default StorageState;