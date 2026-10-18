const DIALOG_SELECTORS = {
    open: '[data-dialog-open]',
    close: '[data-dialog-close]',
    static: '[data-dialog-static]',
};

export function initDialogs() {
    document.addEventListener('click', (event) => {
        const openButton = event.target.closest(DIALOG_SELECTORS.open);

        if (openButton) {
            event.preventDefault();

            const dialogId = openButton.dataset.dialogOpen;
            const dialog = document.getElementById(dialogId);

            if (dialog instanceof HTMLDialogElement) {
                dialog.showModal();
                dialog.focus();
            }

            return;
        }

        const closeButton = event.target.closest(DIALOG_SELECTORS.close);

        if (closeButton) {
            const dialog = closeButton.closest('dialog');

            if (dialog instanceof HTMLDialogElement) {
                dialog.close();
            }

            return;
        }

        if (event.target instanceof HTMLDialogElement) {
            // Ignore clicks on the backdrop of static dialogs
            if (event.target.matches(DIALOG_SELECTORS.static)) {
                return;
            }

            event.target.close();
        }
    });
}