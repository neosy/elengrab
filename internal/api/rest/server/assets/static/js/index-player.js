// -------------------------------------------------------------
// Player Overlay Logic
// Handles video/audio playback in modal overlay or bottom bar
// -------------------------------------------------------------

export function initPlayer() {
    const overlay           = document.getElementById("media-player-overlay");
    const videoWrapper      = document.getElementById("media-player__wrapper");
    let   audioBarContainer = document.getElementById("audio-bar");

    if (!overlay || !videoWrapper) return;

    // Create audio container if missing
    if (!audioBarContainer) {
        audioBarContainer = document.createElement("div");
        audioBarContainer.id = "audio-bar";
        document.body.appendChild(audioBarContainer);
    }

    // Force initial hidden state
    overlay.style.display = "none";

    document.addEventListener("click", (event) => {
        const playBtn = event.target.closest(".avatar-play__button");
        if (!playBtn) return;

        const row = playBtn.closest(".grab-result__row");
        if (!row) return;

        const mediaURL = row.dataset.media;
        if (!mediaURL) return;

        const isAudio = row.dataset.isAudio === "true";

        // Clean previous players
        videoWrapper.innerHTML = "";
        audioBarContainer.innerHTML = "";

        let element;
        if (isAudio) {
            element = document.createElement("audio");
            element.controls = true;
             element.autoplay = true;
        } else {
            element = document.createElement("video");
            element.style.background = "black";
            element.controls = true;
            element.autoplay = true;
            // Disable Picture-in-Picture
            element.disablePictureInPicture = true;
        }

        element.src = mediaURL;

        if (isAudio) {
            // Audio → bottom fixed bar
            const bar = document.createElement("div");
            bar.className = "audio-player-bar";

            const closeBtn = document.createElement("button");
            closeBtn.className = "audio-close-btn";
            closeBtn.innerHTML = "×";
            closeBtn.setAttribute("aria-label", "Close audio player");
            closeBtn.onclick = closePlayer;

            bar.appendChild(closeBtn);
            bar.appendChild(element);
            audioBarContainer.appendChild(bar);

            overlay.style.display = "none !important";   // forceful hide
            document.body.style.overflow = "";
            document.body.classList.add("audio-playing");
        } else {
            // Video → centered overlay
            const wrapper = document.createElement("div");
            wrapper.className = "player-wrapper-modern";

            const closeBtn = document.createElement("button");
            closeBtn.className = "player-close-modern";
            closeBtn.innerHTML = "×";
            closeBtn.setAttribute("aria-label", "Close player");
            closeBtn.onclick = closePlayer;

            wrapper.appendChild(closeBtn);
            wrapper.appendChild(element);
            videoWrapper.appendChild(wrapper);

            overlay.style.display = "flex";
            document.body.style.overflow = "hidden";
        }
    });

    // Close overlay on background click (video only)
    overlay.addEventListener("click", (event) => {
        if (event.target === overlay) closePlayer();
    });

    // Global ESC handler
    document.addEventListener("keydown", (event) => {
        if (event.key === "Escape") closePlayer();
    });

    function closePlayer() {
        videoWrapper.innerHTML = "";
        if (audioBarContainer) audioBarContainer.innerHTML = "";
        overlay.style.display = "none";
        document.body.style.overflow = "";
        document.body.classList.remove("audio-playing");
    }
}