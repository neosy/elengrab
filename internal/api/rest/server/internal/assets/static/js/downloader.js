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
// Element selectors and cookie names
// -------------------------------------------------------------
const SELECT_NAMES = {
    qualityCodec: "quality-codec",
    qualityResolution: "quality-resolution",
    format: "format"
};

const COOKIE_NAMES = {
    qualityCodec: "selectQualityCodec",
    qualityResolution: "selectQualityResolution",
    format: "selectFormat"
};

// -------------------------------------------------------------
// Helper: get select element by name
// -------------------------------------------------------------
function getSelectByName(name) {
    return document.querySelector(`select[name="${name}"]`);
}

// -------------------------------------------------------------
// Save all select values to cookies
// -------------------------------------------------------------
function saveAllSelectsToCookie() {
    Object.entries(SELECT_NAMES).forEach(([key, name]) => {
        const el = getSelectByName(name);
        if (el) cookie.set(COOKIE_NAMES[key], el.value);
    });
}

// -------------------------------------------------------------
// Restore select value from cookie
// -------------------------------------------------------------
function setupCookieSelectSync(selectName, cookieName) {
    const selectElement = getSelectByName(selectName);
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

// -------------------------------------------------------------
// Function: setupQualityFormatLogic
// Handles enabling/disabling format options based on quality
// -------------------------------------------------------------
function setupQualityFormatLogic() {
    const qualityCodecSelect = getSelectByName(SELECT_NAMES.qualityCodec);
    const qualityResolutionSelect = getSelectByName(SELECT_NAMES.qualityResolution);
    const formatSelect = getSelectByName(SELECT_NAMES.format);

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
        saveAllSelectsToCookie();
    };

    updateFormatOptions();

    qualityCodecSelect.addEventListener("change", updateFormatOptions);
    formatSelect.addEventListener("change", updateFormatOptions);
}

// -------------------------------------------------------------
// Main Init
// -------------------------------------------------------------
document.addEventListener('DOMContentLoaded', () => {

    const formGrab = document.querySelector('#form-grab');
    const buttonGrab = document.querySelector('.button-grab-get');
    const inputURL = document.querySelector('#youtubeURL');
    const resultDivInfo = document.querySelector('#grab-result-info');

    // Sync selects with cookies
    setupCookieSelectSync(SELECT_NAMES.qualityCodec, COOKIE_NAMES.qualityCodec);
    setupCookieSelectSync(SELECT_NAMES.qualityResolution, COOKIE_NAMES.qualityResolution);
    setupCookieSelectSync(SELECT_NAMES.format, COOKIE_NAMES.format);

    // Submit on Enter
    inputURL.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            event.preventDefault();
            buttonGrab.click();
        }
    });

    // Clear before HTMX request
    htmx.on('#form-grab', 'htmx:beforeRequest', () => {
        if (inputURL) inputURL.value = '';
        if (resultDivInfo) resultDivInfo.innerHTML = '';
    });

    // Display error on non-200 + non-503
    document.body.addEventListener('htmx:afterOnLoad', (event) => {
        if (event.detail.elt === formGrab) {
            if (inputURL) inputURL.value = '';
            if (event.detail.xhr.status !== 200 &&
                event.detail.xhr.status !== 503) {

                if (resultDivInfo) {
                    resultDivInfo.innerHTML = `
                        <div class="div-grab-result-info-item">
                            <span class="result-failed">Error: ${event.detail.xhr.responseText}</span>
                        </div>
                    `;
                }
            }
        }
    });

    // Save resultItemsOnlyOne before HTMX request
    document.body.addEventListener('htmx:beforeRequest', function (evt) {
        const div = document.getElementById("grab-result-item-replaceme");
        if (div) {
            cookie.set('resultItemsOnlyOne', div.dataset.onlyOne, 7);
        }
    });

    // Init quality/format sync
    setupQualityFormatLogic();
});
