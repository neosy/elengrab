const notificationArea = document.getElementById("notification-area");
const durationDefault = 4000;

export const notifyType = {
    SUCCESS: 'success',
    ERROR:   'error',
    WARNING: 'warning',
    INFO:    'info'
};

// Move the notification area to the top of the top layer,
// so it is rendered above any open modal dialog.
function raiseNotificationArea() {
    if (notificationArea.matches(':popover-open')) {
        notificationArea.hidePopover();
    }
    notificationArea.showPopover();
}

export function show(message, type = notifyType.INFO, duration = durationDefault) {
    const notification = document.createElement('div');
    notification.className = `notification notification--${type}`;
    notification.textContent = message;

    notificationArea.appendChild(notification);
    raiseNotificationArea();

    setTimeout(() => {
        notification.classList.add('notification--show');
    }, 10);

    setTimeout(() => {
        notification.classList.remove('notification--show');

        setTimeout(() => {
            notification.remove();
        }, 400);
    }, duration);
}

let errorMessageTimer = null;

/**
 * Show an error message in a result info element with fade-in/out effect.
 * @param {string} text - Message to display
 * @param {HTMLElement} container - Element that will get the "show" class
 * @param {HTMLElement} messageElem - Element where the message text will be inserted
 * @param {number} duration - Duration in ms to show the message (default 5000)
 */
export function showErrorMessage(text, container, messageElem, duration = 5000) {
    if (!container || !messageElem) return;

    // Reset the previous timer
    if (errorMessageTimer) {
        clearTimeout(errorMessageTimer);
        errorMessageTimer = null;
    }

    // Update the message text
    messageElem.textContent = `Error: ${text}`;

    // Show the element
    container.classList.add("show");

    // Hide after a set time
    errorMessageTimer = setTimeout(() => {
        container.classList.remove("show");
        errorMessageTimer = null;
    }, duration);
}