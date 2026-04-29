const notificationArea = document.getElementById("notification-area");
const durationDefault = 4000;

export const notifyType = {
    SUCCESS: 'success',
    ERROR:   'error',
    WARNING: 'warning',
    INFO:    'info'
};

export function show(message, type = NotificationType.INFO, duration = durationDefault) {
    const notification = document.createElement('div');
    notification.className = `notification notification--${type}`;
    notification.textContent = message;

    notificationArea.appendChild(notification);

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