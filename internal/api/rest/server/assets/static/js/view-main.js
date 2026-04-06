const DOM_ELEMENTS = {
    player: null,

    customControls: null,
    backButton: null,
    progressBar: null,
    volumeControl: null,
    playButton: null,
    skipBackwardButton: null,
    skipForwardButton: null,
    toggleFullscreenButton: null,

    mediaTitle: null,
    mediaDescription: null,
};

function initDomElements() {
    DOM_ELEMENTS.backButton = document.getElementById("backButton");
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

document.addEventListener('DOMContentLoaded', () => {
    initDomElements();

    document.addEventListener('keydown', function(e) {
        if (e.code === "Space") {
            e.preventDefault();
            togglePlay();
        }
        if (e.code === "ArrowRight") skipForward();
        if (e.code === "ArrowLeft") skipBackward();
    });

    initPlayer();
});

function initPlayer() {
    DOM_ELEMENTS.backButton && (DOM_ELEMENTS.backButton.addEventListener('click', goBack));

    if (DOM_ELEMENTS.customControls) {
        DOM_ELEMENTS.playButton.addEventListener('click', togglePlay);
        DOM_ELEMENTS.skipBackwardButton.addEventListener('click', skipBackward);
        DOM_ELEMENTS.skipForwardButton.addEventListener('click', skipForward);
        DOM_ELEMENTS.toggleFullscreenButton.addEventListener('click', toggleFullscreen);

        DOM_ELEMENTS.progressBar.addEventListener('click', function(e) {
            const rect = this.getBoundingClientRect();
            const clickPosition = (e.clientX - rect.left) / rect.width;
            DOM_ELEMENTS.player.currentTime = clickPosition * DOM_ELEMENTS.player.duration;
        });

        DOM_ELEMENTS.volumeControl.addEventListener('input', function() {
            DOM_ELEMENTS.player.volume = this.value;
        });
    }

    const params = new URLSearchParams(window.location.search);
    if (params.has('title')) currentMedia.title = params.get('title');
    if (params.has('url')) currentMedia.url = params.get('url');
    if (params.has('type')) currentMedia.type = params.get('type');

    // Auto play
    DOM_ELEMENTS.player.addEventListener('loadedmetadata', () => {
        DOM_ELEMENTS.player.autoplay = true;
        DOM_ELEMENTS.player.play().catch(err => {
            console.warn('Autorun failed', err);
        });
    });    
}

function formatTime(seconds) {
    if (!seconds || isNaN(seconds)) return "0:00";
    const min = Math.floor(seconds / 60);
    const sec = Math.floor(seconds % 60);
    return `${min}:${sec < 10 ? '0' : ''}${sec}`;
}

function togglePlay() {
    if (DOM_ELEMENTS.player.paused) {
        DOM_ELEMENTS.player.play();
        if (DOM_ELEMENTS.playButton) {
            DOM_ELEMENTS.playButton.innerHTML = '❚❚';
        }
    } else {
        DOM_ELEMENTS.player.pause();
        if (DOM_ELEMENTS.playButton) {
            DOM_ELEMENTS.playButton.innerHTML = '▶';
        }
    }
}

function updateProgressBar() {
    const progress = document.getElementById('progress');
    const currentTimeEl = document.getElementById('currentTime');
    
    const percentage = (DOM_ELEMENTS.player.currentTime / DOM_ELEMENTS.player.duration) * 100;
    progress && (progress.style.width = percentage + '%');
    currentTimeEl && (currentTimeEl.textContent = formatTime(DOM_ELEMENTS.player.currentTime));
}

function skipForward() {
    DOM_ELEMENTS.player.currentTime += 10;
}

function skipBackward() {
    DOM_ELEMENTS.player.currentTime -= 10;
}

function toggleFullscreen() {
    const container = document.getElementById('mediaContainer');
    if (!document.fullscreenElement) {
        container.requestFullscreen();
    } else {
        document.exitFullscreen();
    }
}

function goBack() {
    // window.history.back();
    console.log("1111");
    window.location.href = "/";
}
