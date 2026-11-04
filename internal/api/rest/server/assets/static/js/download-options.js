import storage from './storage-state.js';

// -------------------------------------------------------------
// CSS class names and selectors
// -------------------------------------------------------------
const CLASS_NAMES = {
    downloadOptions: "download-options",
    downloadOptionsList: "download-options__list",
}

const OPTION_NAMES = {
    qualityCodec: "quality-codec",
    qualityResolution: "quality-resolution",
    format: "format"
};

export class StorageKeys {
    qualityCodec = null;
    qualityResolution = null;
    format = null;

    constructor({ qualityCodec, qualityResolution, format } = {}) {
        if (qualityCodec !== undefined) this.qualityCodec = qualityCodec;
        if (qualityResolution !== undefined) this.qualityResolution = qualityResolution;
        if (format !== undefined) this.format = format;
    }    

    isNull() {
        return Object.values(this).some((el) => el === null);
    }
}

//-------------------------------------------------------------
// Function: $
/**
 * Helper function to select an element by class name within a root element.
 * @param {HTMLElement} root - The root element to search within.
 * @param {string} className - The class name of the element to select.
 * @returns {HTMLElement|null} - The selected element or null if not found.
 */
function $(root, className) {
    return root.querySelector(`.${className}`);
}

function $select(root, name) {
    return root.querySelector(`select[name="${name}"]`);
}

//-------------------------------------------------------------
// Function: initDomElements
/**
 * Initialize DOM elements for download options.
 * @param {HTMLElement} root - The root element containing the download options.
 */
// Handles enabling/disabling format options based on quality
function initDomElements(root) {
    const elements = {
        downloadOptions: $(root, CLASS_NAMES.downloadOptions),
        downloadOptionsList: $(root, CLASS_NAMES.downloadOptionsList),
    };

    return Object.values(elements).some((el) => el === null) ? null : elements;
}

// -------------------------------------------------------------
// Function: setupQualityFormatLogic
// Handles enabling/disabling format options based on quality
/**
 * Initialize the logic for enabling/disabling format options based on the selected quality.
 * @param {Object}
 * @param {HTMLElement|Document} ctx.root - Root element to query within.
 * @param {StorageKeys} ctx.storageKeys - Instance of StorageKeys holding storage key mapping.
 * @returns {Object|null} - An object containing the initialized elements or null if initialization failed.
*/
// -------------------------------------------------------------
export function initQualityFormatLogic(ctx) {
    const {root, storageKeys} = ctx;

    const elements = initDomElements(root);

    if (!elements) {
        return;
    }

    const qualityCodecSelect = $select(root, OPTION_NAMES.qualityCodec);
    const qualityResolutionSelect = $select(root, OPTION_NAMES.qualityResolution);
    const formatSelect = $select(root, OPTION_NAMES.format);

    if (!qualityCodecSelect || !qualityResolutionSelect || !formatSelect) {
        return;
    }

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

    if (storageKeys) {
        storage.syncSelect(elements.downloadOptions, OPTION_NAMES.qualityCodec, storageKeys.qualityCodec);
        storage.syncSelect(elements.downloadOptions, OPTION_NAMES.qualityResolution, storageKeys.qualityResolution, true);
        storage.syncSelect(elements.downloadOptions, OPTION_NAMES.format, storageKeys.format);
    }

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

        if (storageKeys) {
            storage.saveAllSelects(elements.downloadOptions, OPTION_NAMES, storageKeys);
        }
    };

    updateFormatOptions();

    qualityCodecSelect.addEventListener("change", updateFormatOptions);
    formatSelect.addEventListener("change", updateFormatOptions);

    return {
        elements: elements,
        optionNames: OPTION_NAMES,
    };
}
