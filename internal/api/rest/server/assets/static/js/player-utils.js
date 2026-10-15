function formatDuration(seconds) {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const sec = (seconds % 60).toString().padStart(2, "0");

    return days > 0 
        ? `${days}:${hours.toString().padStart(2, "0")}:${minutes.toString().padStart(2, "0")}:${sec}` 
        : hours > 0 
            ? `${hours}:${minutes.toString().padStart(2, "0")}:${sec}`
            : `${minutes}:${sec}`;
}

/**
Updates the remaining playback duration of the specified player.
@param {HTMLElement} durationElement - Element to display the remaining duration.
@param {HTMLMediaElement} player - Audio or video player to get the playback time from.
*/
export function updateRemainingDuration(durationElement, player) {
    if (!durationElement || !player) {
        return;
    }

    let remainingSeconds = null;

    if (Number.isFinite(player.duration)) {
        remainingSeconds = Math.max(
            0,
            Math.floor(player.duration - player.currentTime)
        );
    }

    if (remainingSeconds !== null) {
        durationElement.textContent = formatDuration(remainingSeconds);
    } else {
        durationElement.textContent = "-:-";
    }
}

/**
 * Updates the playback progress and buffer indicators.
 *
 * @param {HTMLElement} bufferElement - Element displaying the buffered progress.
 * @param {HTMLElement} progressElement - Element displaying the current playback progress.
 * @param {HTMLMediaElement} player - Audio or video player to get playback data from.
 * @param {Object} cssVariables - CSS custom property names used to update the indicators.
 * @param {string} cssVariables.watchBuffer - CSS custom property for the buffered progress.
 * @param {string} cssVariables.watchProgress - CSS custom property for the playback progress.
 */
export function updateProgress(bufferElement, progressElement, player, cssVariables) {
    if (!player || !cssVariables) {
        return;
    }

    if (bufferElement && cssVariables.watchBuffer) {
        const bufferPercent = (player.duration > 0 && player.buffered.length > 0)
        ? Math.floor((player.buffered.end(player.buffered.length - 1) / player.duration) * 100)
        : 0;

        bufferElement.style.setProperty(
            cssVariables.watchBuffer,
            `${bufferPercent}%`
        );
    }

    if (progressElement && cssVariables.watchProgress) {
        const progressPercent = player.duration > 0
        ? Math.floor((player.currentTime / player.duration) * 1000) / 10
        : 0;

        progressElement.style.setProperty(
            cssVariables.watchProgress,
            `${progressPercent}%`
        );
    }
}
