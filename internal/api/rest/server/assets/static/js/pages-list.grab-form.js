import * as cookie from './cookie.js';
import * as utils from './utils.js';
import * as actionButton from './action-buttons.js';

// -------------------------------------------------------------
// Element selectors and cookie names
// -------------------------------------------------------------
const SELECT_NAMES = {
    qualityCodec: "quality-codec",
    qualityResolution: "quality-resolution",
    format: "format"
};

const COOKIE_NAMES = {
    qualityCodec: "select_quality_codec",
    qualityResolution: "select_quality_resolution",
    format: "select_format"
};

export const DOM_IDS = {
    mediaURLInput: "mediaURLInput",
}

export const DOM_ELEMENTS = {
    grabForm: null,
    mediaURLInput: null,
    inputActionBtn: null,
    inputActionSettingsBtn: null,
    grabOptionsCollapse: null,
    grabOptions: null,
};

function initDomElements() {    
    DOM_ELEMENTS.grabForm = document.getElementById("grabForm");
    DOM_ELEMENTS.mediaURLInput = document.getElementById(DOM_IDS.mediaURLInput);
    DOM_ELEMENTS.inputActionBtn = document.getElementById("inputActionBtn");
    DOM_ELEMENTS.inputActionSettingsBtn = document.getElementById("inputActionSettingsBtn");
    DOM_ELEMENTS.grabOptionsCollapse = document.getElementById("grabOptionsCollapse");
    DOM_ELEMENTS.grabOptions = document.getElementById("grabOptions");
}

// -------------------------------------------------------------
// Function: setupQualityFormatLogic
// Handles enabling/disabling format options based on quality
// -------------------------------------------------------------
function setupQualityFormatLogic() {
    const qualityCodecSelect = utils.getSelectByName(SELECT_NAMES.qualityCodec);
    const qualityResolutionSelect = utils.getSelectByName(SELECT_NAMES.qualityResolution);
    const formatSelect = utils.getSelectByName(SELECT_NAMES.format);

    if (!qualityCodecSelect || !qualityResolutionSelect || !formatSelect) return;

    const videoFormats = ["auto", "mp4", "webm"];
    const audioFormats = ["auto", "mp3", "m4a", "flac", "opus"];
    const onlyAudioFormats = ["mp3", "m4a", "flac", "opus"];
    const webmCodecs = ["best", "av1"];
    const videoFormatWebmValue = "webm";
    const bestValue = "best"
    const maxValue = "max"
    const autoValue = "auto"
    const emptyValue = "empty"
    const onlyAudioValue = "only_audio";
    const videoFormatDefault = autoValue;
    const audioFormatDefault = autoValue;

    // Update format options based on quality
    const updateFormatOptions = () => {
        const qualityCodec = qualityCodecSelect.value;
        const isQualityCodecBest = (qualityCodec === bestValue);
        const isCodecOnlyAudio = (qualityCodec === onlyAudioValue);
        const isCodecFormatOnlyAudio = isCodecOnlyAudio || (isQualityCodecBest && onlyAudioFormats.includes(formatSelect.value));

        if (isCodecFormatOnlyAudio && !qualityResolutionSelect.disabled) {
            qualityResolutionSelect.disabled = true;
            qualityResolutionSelect.value = emptyValue;
        }

        if (!isCodecFormatOnlyAudio && qualityResolutionSelect.disabled) {
            qualityResolutionSelect.disabled = false;
            qualityResolutionSelect.value = maxValue;
        }

        formatSelect.querySelectorAll("option").forEach(option => {
            const value = option.value;

            if (isCodecOnlyAudio) {
                const allowed = audioFormats.includes(value);
                option.disabled = !allowed;
            } else {
                let allowed = isQualityCodecBest || videoFormats.includes(value);
                if (value == videoFormatWebmValue) {
                    allowed = webmCodecs.includes(qualityCodec);
                }
                option.disabled = !allowed;
            }
        });

        if (!isQualityCodecBest) {
            if (isCodecOnlyAudio && !audioFormats.includes(formatSelect.value)) {
                formatSelect.value = audioFormatDefault;
            }

            if (!isCodecOnlyAudio && !videoFormats.includes(formatSelect.value)) {
                formatSelect.value = videoFormatDefault;
            }

            if (!isCodecOnlyAudio && formatSelect.value == videoFormatWebmValue && !webmCodecs.includes(qualityCodec)) {
                formatSelect.value = videoFormatDefault;
            }
        }
        cookie.saveAllSelectsToCookie(SELECT_NAMES, COOKIE_NAMES);
    };

    updateFormatOptions();

    qualityCodecSelect.addEventListener("change", updateFormatOptions);
    formatSelect.addEventListener("change", updateFormatOptions);
}

export function initMediaGrabForm() {
    initDomElements();

    if (!DOM_ELEMENTS.grabForm) {
        return;
    }

    const grabForm = DOM_ELEMENTS.grabForm;
    const buttonGrab = document.querySelector('.grab-area__submit-button');

    const grabURLInput = DOM_ELEMENTS.mediaURLInput;
    const grabInputActionBtn = DOM_ELEMENTS.inputActionBtn;

    cookie.setupCookieSelectSync(SELECT_NAMES.qualityCodec, COOKIE_NAMES.qualityCodec);
    cookie.setupCookieSelectSync(SELECT_NAMES.qualityResolution, COOKIE_NAMES.qualityResolution, true);
    cookie.setupCookieSelectSync(SELECT_NAMES.format, COOKIE_NAMES.format);

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
                actionButton.updateInputPasteClearButton(grabURLInput, grabInputActionBtn);
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

    // Init quality/format sync
    setupQualityFormatLogic();

    // Init settiongs action button
    actionButton.initInputSettingsButton(DOM_ELEMENTS.inputActionSettingsBtn, DOM_ELEMENTS.grabOptionsCollapse, DOM_ELEMENTS.grabOptions);

    // Init action button for input field
    actionButton.initInputPasteClearButton(DOM_ELEMENTS.mediaURLInput, DOM_ELEMENTS.inputActionBtn);
}