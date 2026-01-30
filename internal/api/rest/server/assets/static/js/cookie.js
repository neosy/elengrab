import * as helper from './helper.js';
import { SELECT_NAMES, COOKIE_NAMES } from './constants.js';

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
export function saveAllSelectsToCookie() {
    Object.entries(SELECT_NAMES).forEach(([key, name]) => {
        const el = helper.getSelectByName(name);
        if (el) cookie.set(COOKIE_NAMES[key], el.value);
    });
}

// -------------------------------------------------------------
// Restore select value from cookie
// -------------------------------------------------------------
export function setupCookieSelectSync(selectName, cookieName) {
    const selectElement = helper.getSelectByName(selectName);
    if (!selectElement) return;

    const savedValue = cookie.get(cookieName);
    if (savedValue) {
        const option = selectElement.querySelector(`option[value="${savedValue}"]`);
        if (option) option.selected = true;
    }

    // For resolution, save only its value on change
    if (selectElement.name === SELECT_NAMES.qualityResolution) {
        selectElement.addEventListener("change", () => {
            cookie.set(cookieName, selectElement.value);
        });
    }
}
