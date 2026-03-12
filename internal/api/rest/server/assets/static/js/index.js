import * as helper from './helper.js';
import * as cookie from './cookie.js';
import * as actionButton from './action-button.js';
import * as player from './player.js';
import { SELECT_NAMES, COOKIE_NAMES } from './constants.js';

// -------------------------------------------------------------
// Function: setupQualityFormatLogic
// Handles enabling/disabling format options based on quality
// -------------------------------------------------------------
function setupQualityFormatLogic() {
    const qualityCodecSelect = helper.getSelectByName(SELECT_NAMES.qualityCodec);
    const qualityResolutionSelect = helper.getSelectByName(SELECT_NAMES.qualityResolution);
    const formatSelect = helper.getSelectByName(SELECT_NAMES.format);

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
        cookie.saveAllSelectsToCookie();
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
    const inputURL = document.querySelector('#mediaURL');
    const resultDivInfo = document.querySelector('#grab-result-info');

    // Sync selects with cookies
    cookie.setupCookieSelectSync(SELECT_NAMES.qualityCodec, COOKIE_NAMES.qualityCodec);
    cookie.setupCookieSelectSync(SELECT_NAMES.qualityResolution, COOKIE_NAMES.qualityResolution);
    cookie.setupCookieSelectSync(SELECT_NAMES.format, COOKIE_NAMES.format);

    // Submit on Enter
    inputURL.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            event.preventDefault();
            buttonGrab.click();
        }
    });

    // Clear before HTMX request
    htmx.on('#form-grab', 'htmx:beforeRequest', () => {
        if (inputURL) {
            inputURL.value = '';
            // update action button after clearing
            actionButton.updateActionButton();
        }
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
                        <div class="div-grab-result-info-row">
                            <span class="result-failed">Error: ${event.detail.xhr.responseText}</span>
                        </div>
                    `;
                }
            }
        }
    });

    // Init quality/format sync
    setupQualityFormatLogic();

    // Init action button for input field
    actionButton.updateActionButton();
    inputURL.addEventListener('input', actionButton.updateActionButton);

    // Subscribe to SSE row-delete event
    const es = new EventSource("/ui/downloader/files/events");
    es.addEventListener("row-add", actionButton.handleRowAdd);
    es.addEventListener("row-update", actionButton.handleRowUpdate);
    es.addEventListener("row-delete", actionButton.handleRowDelete);

    // Init inline media player
    player.initPlayer();
});
