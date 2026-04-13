// -------------------------------------------------------------
// Helper: get select element by name
// -------------------------------------------------------------
export function getSelectByName(name) {
    return document.querySelector(`select[name="${name}"]`);
}


let errorMessageTimer = null;

/**
 * Show an error message in a result info element with fade-in/out effect.
 * @param {string} text - Message to display
 * @param {HTMLElement} container - Element that will get the "show" class
 * @param {HTMLElement} messageElem - Element where the message text will be inserted
 * @param {number} duration - Duration in ms to show the message (default 2000)
 */
export function showErrorMessage(text, container, messageElem, duration = 5000) {
    if (!container || !messageElem) return;

    // Reset the previous timer
    if (errorMessageTimer) {
        clearTimeout(errorMessageTimer);
        errorMessageTimer = null;
    }

    // Updating the message text
    messageElem.textContent = `Error: ${text}`;

    // Показать элемент
    container.classList.add("show");

    // Hide after a set time
    errorMessageTimer = setTimeout(() => {
        container.classList.remove("show");
        errorMessageTimer = null;
    }, duration);
}

/**
 * Universal share helper
 * @param {Object} options
 * @param {string} options.endpoint - API endpoint to get short URL
 * @param {string} [options.method='GET'] - HTTP method (GET, POST, etc.)
 * @param {Object} [options.fetchOptions] - additional fetch options (headers, body, etc.)
 * @param {string} [options.title]
 * @param {string} [options.text]
 * @returns {Promise<{status: string, url?: string, error?: any}>}
 */
export async function shareLink({
  endpoint = '/api/share',
  method = 'GET',
  fetchOptions = {},
  title = 'My App',
  text = 'Check this out'
} = {}) {
  try {
    // 1. Request short URL
    const response = await fetch(endpoint, {
      method,
      ...fetchOptions
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const { url } = await response.json();

    if (!url) {
      throw new Error('URL missing in response');
    }

    // 2. Native share
    if (navigator.share) {
      await navigator.share({ title, text, url });
      return { status: 'shared', url };
    }

    // 3. Clipboard fallback
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(url);
      return { status: 'copied', url };
    }

    // 4. Last fallback
    return { status: 'manual', url };

  } catch (error) {
    if (error.name === 'AbortError') {
      return { status: 'cancelled', url };
    }

    return { status: 'error', error };
  }
}