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