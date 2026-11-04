/* ------------------------------------------------------------------ */
/* Helpers                                                             */
/* ------------------------------------------------------------------ */

function $select(root, name) {
    return root.querySelector(`select[name="${name}"]`);
}

// -------------------------------------------------------------
// Cookie Helper Module
// -------------------------------------------------------------
const cookie = (() => {

    // Get cookie by name
    const get = (name) => {
        return document.cookie
            .split("; ")
            .find(row => row.startsWith(name + "="))
            ?.split("=")[1];
    };

    // Set cookie with expiration (days)
    const set = (name, value, days = 365) => {
        const d = new Date();
        d.setTime(d.getTime() + (days * 24 * 60 * 60 * 1000));
        document.cookie =
            `${name}=${encodeURIComponent(value)};expires=${d.toUTCString()};path=/`;
    };

    return { get, set };
})();

// -------------------------------------------------------------
// Save all select values to cookies
// -------------------------------------------------------------
export function saveAllSelectsToCookie(root, names, cookieNames) {
    Object.entries(names).forEach(([key, name]) => {
        const el = $select(root, name);
        if (el) cookie.set(cookieNames[key], el.value);
    });
}

// -------------------------------------------------------------
// Restore select value from cookie
// -------------------------------------------------------------
export function setupCookieSelectSync(root, selectName, cookieName, saveOnChange = false) {
    const selectElement = $select(root, selectName);
    if (!selectElement) return;

    const savedValue = cookie.get(cookieName);
    if (savedValue) {
        const option = selectElement.querySelector(`option[value="${savedValue}"]`);
        if (option) option.selected = true;
    }

    // Save the selected value on change when enabled
    if (saveOnChange) {
        selectElement.addEventListener("change", () => {
            cookie.set(cookieName, selectElement.value);
        });
    }
}
