import { DOM_IDS, DOM_CLASSES, DOM_ELEMENTS, initDomElements } from "./edit-media.dom.js";
import * as notify from './notifications.js';

document.addEventListener('DOMContentLoaded', () => {
    initDomElements();
    initInputElements();
    initSaveButton();
});

function initInputElements() {
    const fields = document.querySelectorAll(`.${DOM_CLASSES.mediaEditorFieldInput}`);
    const maxWidth = 160;

    function updateCounterAndStyle(field, textArea, charCounter, maxLength) {
        if (!field || !textArea || !charCounter) {
            console.error("Input elements was not found.");
            return;
        }

        const currentLength = textArea.value.length;
        charCounter.textContent = `${currentLength}/${maxLength}`;

        if (currentLength > maxLength) {
            field.classList.add(DOM_CLASSES.mediaEditorFieldError);
        } else {
            field.classList.remove(DOM_CLASSES.mediaEditorFieldError);
        }
    }

    function autoResize(textArea) {
        textArea.style.height = 'auto';
        const newHeight = Math.min(textArea.scrollHeight, maxWidth);
        textArea.style.height = newHeight + 'px';
    }    

    fields.forEach(field => {
        const textArea = field.querySelector(`.${DOM_CLASSES.mediaEditorTextArea}`);
        const charCounter = field.querySelector(`.${DOM_CLASSES.mediaEditorCharCounter}`);

        if (!textArea || !charCounter) return;

        const maxLength = textArea.dataset.maxLength;

        field.addEventListener('click', () => {
            textArea.focus();
        });

        autoResize(textArea);

        textArea.addEventListener('input', () => {
            updateCounterAndStyle(field, textArea, charCounter, maxLength);
            autoResize(textArea);
        });

        textArea.addEventListener('focus', () => {
            updateCounterAndStyle(field, textArea, charCounter, maxLength);
            autoResize(textArea);
        });
    });
}

function initSaveButton() {
    const saveButton = DOM_ELEMENTS.saveButton;
    const cancelButton = DOM_ELEMENTS.cancelButton;
    const titleEl = DOM_ELEMENTS.mediaTitleInput;
    const descriptionEl = DOM_ELEMENTS.mediaDescriptionInput;
    const visibilityEl = DOM_ELEMENTS.mediaVisibilitySelect;

    if (!titleEl || !descriptionEl || !visibilityEl) return;

    async function save() {
        const response = await fetch(saveButton.dataset.patchUrl, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                title: titleEl.value,
                description: descriptionEl.value,
                visibility: visibilityEl.value
            })
        })

        if (!response.ok) {
            let message = "";

            if (response.status >= 400 && response.status < 500) {
                try {
                    const data = await response.json();
                    message = data.message;
                } catch {
                    message = `Request error: ${response.status}`;
                }
            } else {
                message = `Server error: ${response.status}`;
            }

            notify.show(message, notify.notifyType.ERROR);
            return;
        }        

        notify.show("Changes saved successfully", notify.notifyType.INFO);
    }

    if (saveButton) {
        saveButton.addEventListener('click', save);
    }

    if (cancelButton) {
        cancelButton.addEventListener('click', () => {
            history.back();
        });
    }
}