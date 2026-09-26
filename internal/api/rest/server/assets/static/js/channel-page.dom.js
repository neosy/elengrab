import * as CONST from './constants.js';
import * as PAGES_DOM from './pages-list.dom.js';

export const CLASS_PREFIXES = {
    ...CONST.CLASS_PREFIXES,
}

export const CLASS_NAMES = {
    ...CONST.CLASS_NAMES,
    ...PAGES_DOM.CLASS_NAMES,
};

export const CLASS_SELECTORS = Object.fromEntries(
    Object.entries(CLASS_NAMES)
        .filter(([, value]) => typeof value === "string")
        .map(([key, value]) => [key, `.${value}`])
);

export const DOM_IDS = {
    ...PAGES_DOM.DOM_IDS,
}

export const DOM_ELEMENTS = {
    ...PAGES_DOM.DOM_ELEMENTS,
};

export function initDomElements() {    
    PAGES_DOM.initDomElements(DOM_ELEMENTS);
}