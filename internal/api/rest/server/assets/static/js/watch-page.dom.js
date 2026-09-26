export const DOM_ELEMENTS = {
    mainContentPlayer: null,
    player: null,
    video: null,
    videoWrapper: null,
    mediaInfo: null,

    customControls: null,
    backButton: null,
    progressBar: null,
    volumeControl: null,
    playButton: null,
    skipBackwardButton: null,
    skipForwardButton: null,
    toggleFullscreenButton: null,

    main: null,
    mediaTitle: null,
    mediaDescription: null,
};

export function initDomElements() {
    DOM_ELEMENTS.main = document.querySelector("main");
    DOM_ELEMENTS.backButton = document.getElementById("backButton");

    DOM_ELEMENTS.mainContentPlayer = document.getElementById("mainContentPlayer");
    DOM_ELEMENTS.videoWrapper = document.getElementById("videoWrapper");
    DOM_ELEMENTS.video = document.getElementById("videoElement");
    DOM_ELEMENTS.mediaInfo = document.getElementById("mainContentMediaInfo")
    DOM_ELEMENTS.player = document.getElementById("videoElement");
    if (!DOM_ELEMENTS.player) {
        DOM_ELEMENTS.player = document.getElementById("audioElement");
    }

    DOM_ELEMENTS.customControls = document.getElementById("customControls");
    if (DOM_ELEMENTS.customControls) {
        DOM_ELEMENTS.progressBar = document.getElementById("progressBar");
        DOM_ELEMENTS.volumeControl = document.getElementById("volumeControl");
        DOM_ELEMENTS.playButton = document.getElementById("playButton");
        DOM_ELEMENTS.skipBackwardButton = document.getElementById("skipBackwardButton");
        DOM_ELEMENTS.skipForwardButton = document.getElementById("skipForwardButton");
        DOM_ELEMENTS.toggleFullscreenButton = document.getElementById("toggleFullscreenButton");
    }

    DOM_ELEMENTS.mediaTitle = document.getElementById("mediaTitle");
    DOM_ELEMENTS.mediaDescription = document.getElementById("mediaDescription");
}
