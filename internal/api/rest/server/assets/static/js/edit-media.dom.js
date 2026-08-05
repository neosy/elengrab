export const DOM_CLASSES = {
    mediaEditorFieldInput: "media-editor__field--input",
    mediaEditorFieldError: "media-editor__field--error",
    mediaEditorTextArea: "media-editor__textarea",
    mediaEditorCharCounter: "media-editor__char-counter",
};

export const DOM_ELEMENTS = {
    saveButton: null,
    cancelButton: null,
    mediaTitleInput: null,
    mediaDescriptionInput: null,
    mediaVisibilitySelect: null,
};

export function initDomElements() {
    DOM_ELEMENTS.saveButton = document.getElementById("saveButton");
    DOM_ELEMENTS.cancelButton = document.getElementById("cancelButton");
    DOM_ELEMENTS.mediaTitleInput = document.getElementById("mediaTitleInput");
    DOM_ELEMENTS.mediaDescriptionInput = document.getElementById("mediaDescriptionInput");
    DOM_ELEMENTS.mediaVisibilitySelect = document.getElementById("mediaVisibilitySelect");
}
