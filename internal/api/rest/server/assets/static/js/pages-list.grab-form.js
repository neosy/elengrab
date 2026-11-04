import * as downloadOptions from './download-options.js';
import * as actionButton from './action-buttons.js';

// -------------------------------------------------------------
// Element selectors and cookie names
// -------------------------------------------------------------
const CLASS_NAMES = {
    actionPasteButton: "grab-area__action-button--paste",
    actionSettingsButton: "grab-area__action-button--settings",
}

const STORAGE_KEYS = new downloadOptions.StorageKeys({
    qualityCodec: "grab_form_quality_codec",
    qualityResolution: "grab_form_quality_resolution",
    format: "grab_form_format"
});

export const DOM_IDS = {
    mediaURLInput: "mediaURLInput",
}

export const DOM_ELEMENTS = {
    grabForm: null,
    mediaURLInput: null,
    inputActionBtn: null,
    inputActionSettingsBtn: null,
};

function $(root, className) {
    return root.querySelector(`.${className}`);
}

function initDomElements() {    
    DOM_ELEMENTS.grabForm = document.getElementById("grabForm");

    DOM_ELEMENTS.mediaURLInput = document.getElementById(DOM_IDS.mediaURLInput);
    DOM_ELEMENTS.inputActionBtn = $(DOM_ELEMENTS.grabForm, CLASS_NAMES.actionPasteButton);
    DOM_ELEMENTS.inputActionSettingsBtn = $(DOM_ELEMENTS.grabForm, CLASS_NAMES.actionSettingsButton);
}

export function initMediaGrabForm() {
    initDomElements();

    if (!DOM_ELEMENTS.grabForm) {
        console.error("Grab form not found. Initialization aborted.");
        return;
    }

    const grabForm = DOM_ELEMENTS.grabForm;
    const buttonGrab = document.querySelector('.url-input-button');
    const grabURLInput = DOM_ELEMENTS.mediaURLInput;

    // Init quality/format sync
    const ctx = {
        root: grabForm,
        storageKeys: STORAGE_KEYS,
    };
    const options = downloadOptions.initQualityFormatLogic(ctx);
    if (!options || !options.elements) {
        console.error("Failed to initialize quality/format logic.");
        return;
    }

    const urlActions = {
        pasteClear: null,
    };

    // Submit on Enter
    if (grabURLInput) {
        grabURLInput.addEventListener('keydown', (event) => {
            if (event.key === 'Enter') {
                event.preventDefault();
                buttonGrab.click();
            }
        });
    }

    // Clear before HTMX request
    if (grabForm) {
        htmx.on(grabForm, 'htmx:beforeRequest', () => {
            if (grabURLInput) {
                grabURLInput.value = '';
                // update action button after clearing
                urlActions.pasteClear?.update();
            }
            if (DOM_ELEMENTS.resultInfo) DOM_ELEMENTS.resultInfo.classList.remove("show");
        });

        // Handle HTMX response for grab form
        htmx.on(grabForm, 'htmx:afterOnLoad', (event) => {
            const xhr = event.detail.xhr;

            // --- Error handling (HTTP >= 400, except 503) ---
            if (xhr.status >= 400 && xhr.status !== 503) {
                if (grabURLInput) grabURLInput.value = '';

                if (DOM_ELEMENTS.resultInfo && DOM_ELEMENTS.resultInfoFailed) {
                    let text = xhr.responseText;

                    try {
                        const data = JSON.parse(text);
                        if (data && typeof data === "object" && "message" in data) {
                            text = data.message;
                        }
                    } catch (e) {
                        // ignore non-JSON
                    }
                }

                return;
            }

            // --- Success: guest session created ---
            let data;
            try {
                data = JSON.parse(xhr.responseText);
            } catch (e) {
                console.error('Invalid JSON response');
                return;
            }

            if (data.guestCreated === true) {
                // Reload to apply new session (cookie)
                window.location.reload();
            }
        });
    }

    // Init action button for input field
    urlActions.pasteClear =  actionButton.initInputPasteClearButton(DOM_ELEMENTS.mediaURLInput, DOM_ELEMENTS.inputActionBtn);

    if (options.elements) {
        // Init settiongs action button
        actionButton.initInputSettingsButton(DOM_ELEMENTS.inputActionSettingsBtn, options.elements.downloadOptions, options.elements.downloadOptionsList);
    }
}