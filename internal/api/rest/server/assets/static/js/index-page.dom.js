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

    mediaURLInput: "mediaURLInput",
}

export const DOM_ELEMENTS = {
    ...PAGES_DOM.DOM_ELEMENTS,

    grabForm: null,
    mediaURLInput: null,
    inputActionBtn: null,
    inputActionSettingsBtn: null,
    grabOptionsCollapse: null,
    grabOptions: null,
};

export function initDomElements() {    
    PAGES_DOM.initDomElements(DOM_ELEMENTS);

    DOM_ELEMENTS.grabForm = document.getElementById("grabForm");
    DOM_ELEMENTS.mediaURLInput = document.getElementById(DOM_IDS.mediaURLInput);
    DOM_ELEMENTS.inputActionBtn = document.getElementById("inputActionBtn");
    DOM_ELEMENTS.inputActionSettingsBtn = document.getElementById("inputActionSettingsBtn");
    DOM_ELEMENTS.grabOptionsCollapse = document.getElementById("grabOptionsCollapse");
    DOM_ELEMENTS.grabOptions = document.getElementById("grabOptions");
}