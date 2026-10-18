const CLASS_NAMES = {
    form: "upload-file-dialog__form",
    input: "upload-file-dialog__input",

    dropzone: "upload-file-dialog__dropzone",
    dropzoneDragover: "upload-file-dialog__dropzone--dragover",

    fileName: "upload-file-dialog__file-name",
    fileSize: "upload-file-dialog__file-size",

    progress: "upload-file-dialog__progress",
    progressVisible: "upload-file-dialog__progress--visible",
    progressTrack: "upload-file-dialog__progress-track",
    progressBar: "upload-file-dialog__progress-bar",
    progressText: "upload-file-dialog__progress-text",

    status: "upload-file-dialog__status",
    statusVisible: "upload-file-dialog__status--visible",
    statusError: "upload-file-dialog__status--error",
    fieldError: "upload-file-dialog__field--error",

    chooseButton: "upload-file-dialog__button--choose",
    doneButton: "upload-file-dialog__button--done",
    removeButton: "upload-file-dialog__remove",

    charCounter: "upload-file-dialog__char-counter",

    formFileSelected: "upload-file-dialog__form--file-selected",
    formUploading: "upload-file-dialog__form--uploading",
    formUploadFailed: "upload-file-dialog__form--upload-failed",
    formUploaded: "upload-file-dialog__form--uploaded",
};

const DOM_IDS = {
    dialog: "upload-file-dialog",
};

const MEDIA_EXTENSIONS =
    /\.(mp3|wav|flac|m4a|aac|ogg|oga|opus|wma|mp4|m4v|mov|mkv|webm|avi|wmv|flv|mpeg|mpg)$/i;

/* ------------------------------------------------------------------ */
/* Helpers                                                             */
/* ------------------------------------------------------------------ */

function $(root, className) {
    return root.querySelector(`.${className}`);
}

function getDomElements(dialog) {
    const elements = {
        dialog,
        form: $(dialog, CLASS_NAMES.form),
        input: $(dialog, CLASS_NAMES.input),
        dropzone: $(dialog, CLASS_NAMES.dropzone),
        chooseButton: $(dialog, CLASS_NAMES.chooseButton),
        fileName: $(dialog, CLASS_NAMES.fileName),
        fileSize: $(dialog, CLASS_NAMES.fileSize),
        removeButton: $(dialog, CLASS_NAMES.removeButton),
        progress: $(dialog, CLASS_NAMES.progress),
        progressTrack: $(dialog, CLASS_NAMES.progressTrack),
        progressBar: $(dialog, CLASS_NAMES.progressBar),
        progressText: $(dialog, CLASS_NAMES.progressText),
        status: $(dialog, CLASS_NAMES.status),
        doneButton: $(dialog, CLASS_NAMES.doneButton),
    };

    return Object.values(elements).some((el) => el === null) ? null : elements;
}

function formatFileSize(size) {
    if (size < 1024) {
        return `${size} B`;
    }
    if (size < 1024 * 1024) {
        return `${(size / 1024).toFixed(1)} KB`;
    }
    if (size < 1024 * 1024 * 1024) {
        return `${(size / (1024 * 1024)).toFixed(1)} MB`;
    }
    return `${(size / (1024 * 1024 * 1024)).toFixed(1)} GB`;
}

function isMediaFile(file) {
    // For some formats (mkv, flac, m4a...) the browser reports an empty type,
    // so fall back to checking the extension. Server-side validation is still required.
    return /^(audio|video)\//.test(file.type) || MEDIA_EXTENSIONS.test(file.name);
}

function stripExtension(name) {
    return name.replace(/\.[^./\\]+$/, "");
}

function getMaxLength(field) {
    // The limit lives in data-max-length (not the maxlength attribute),
    // so typing and pasting are never blocked.
    const value = Number(field.dataset.maxLength);

    return value > 0 ? value : Infinity;
}

function releaseNativeLimit(field) {
    // A native maxlength attribute blocks typing and pasting.
    // If the markup still has one, move its value to data-max-length and remove it.
    const nativeLimit = field.getAttribute("maxlength");

    if (nativeLimit === null) {
        return;
    }

    if (!field.dataset.maxLength) {
        field.dataset.maxLength = nativeLimit;
    }

    field.removeAttribute("maxlength");
}

function isOverLimit(field) {
    return field.value.length > getMaxLength(field);
}

/* ------------------------------------------------------------------ */
/* Uploaders (transport). Interface: start(file, callbacks) -> { abort }  */
/* callbacks: onProgress(percent), onSent(), onSuccess(), onError(msg)   */
/* ------------------------------------------------------------------ */

function createXhrUploader(url) {
    return {
        start(file, { onProgress, onSent, onSuccess, onError }) {
            const xhr = new XMLHttpRequest();
            const formData = new FormData();
            formData.append("file", file);

            xhr.open("POST", url);
            xhr.setRequestHeader("Accept", "application/json");

            xhr.upload.addEventListener("progress", (event) => {
                if (event.lengthComputable) {
                    onProgress((event.loaded / event.total) * 100);
                }
            });

            // The whole file is sent, but the server may still be processing it.
            xhr.upload.addEventListener("load", () => onSent());

            xhr.addEventListener("load", () => {
                if (xhr.status >= 200 && xhr.status < 300) {
                    onSuccess();
                    return;
                }

                let message = "Upload failed. Please try again.";

                try {
                    const response = JSON.parse(xhr.responseText);

                    if (typeof response.error === "string" && response.error) {
                        message = response.error;
                    }
                } catch {
                    // Keep the default message when the response is not valid JSON.
                }

                onError(message);
            });

            xhr.addEventListener("error", () => {
                onError("A network error occurred. Please try again.");
            });

            xhr.addEventListener("timeout", () => {
                onError("The upload timed out. Please try again.");
            });

            xhr.send(formData);

            return { abort: () => xhr.abort() };
        },
    };
}

function createMockUploader({ failRate = 0.25, stepMs = 150 } = {}) {
    return {
        start(file, { onProgress, onSent, onSuccess, onError }) {
            let percent = 0;
            let finishTimer = null;

            const timer = setInterval(() => {
                percent += 5 + Math.random() * 10;

                if (percent < 100) {
                    onProgress(percent);
                    return;
                }

                clearInterval(timer);
                onProgress(100);
                onSent();

                finishTimer = setTimeout(() => {
                    if (Math.random() < failRate) {
                        onError("Mock upload failed.");
                    } else {
                        onSuccess();
                    }
                }, 600);
            }, stepMs);

            return {
                abort: () => {
                    clearInterval(timer);
                    clearTimeout(finishTimer);
                },
            };
        },
    };
}

/* ------------------------------------------------------------------ */
/* Dialog                                                              */
/* ------------------------------------------------------------------ */

/**
 * @param {string} [uploadUrl] Upload URL. If omitted (or mock: true), the mock uploader is used.
 * @param {object} [options]
 * @param {(result: {file: File|null, title: string, description: string, visibility: string}) => void} [options.onDone]
 * @param {boolean} [options.mock]
 * @param {number} [options.mockFailRate] Mock failure rate (0..1), defaults to 0.25
 */
export function initUploadFileDialog(uploadUrl, { onDone, mock = false, mockFailRate = 0.25 } = {}) {
    const dialog = document.getElementById(DOM_IDS.dialog);

    if (!dialog) {
        return;
    }

    const elements = getDomElements(dialog);

    if (!elements) {
        return;
    }

    const {
        form,
        input,
        dropzone,
        chooseButton,
        fileName,
        fileSize,
        removeButton,
        progress,
        progressTrack,
        progressBar,
        progressText,
        status,
        doneButton,
    } = elements;

    if (form.dataset.initialized === "true") {
        return;
    }

    form.dataset.initialized = "true";

    const titleField = form.elements.title;
    const descriptionField = form.elements.description;
    const visibilityField = form.elements.visibility;

    const uploader =
        mock || !uploadUrl
            ? createMockUploader({ failRate: mockFailRate })
            : createXhrUploader(uploadUrl);

    let uploading = false;
    let uploaded = false;
    let job = null;

    /* ---------- UI helpers ---------- */

    const updateDone = () => {
        doneButton.disabled = !(
            uploaded &&
            !uploading &&
            titleField.value.trim() !== "" &&
            !isOverLimit(titleField) &&
            !isOverLimit(descriptionField) &&
            form.checkValidity()
        );
    };

    const setStatus = (message, isError = false) => {
        status.textContent = message;
        status.classList.toggle(CLASS_NAMES.statusVisible, Boolean(message));
        status.classList.toggle(CLASS_NAMES.statusError, Boolean(message) && isError);
    };

    const setProgress = (percent) => {
        const value = Math.max(0, Math.min(100, percent));

        progressBar.style.width = `${value}%`;
        progressTrack.setAttribute("aria-valuenow", String(Math.round(value)));
        progressText.textContent = `${Math.round(value)}%`;
    };

    const setUploading = (value) => {
        uploading = value;

        form.classList.toggle(CLASS_NAMES.formUploading, value);

        input.disabled = value;
        chooseButton.disabled = value;
        removeButton.disabled = value;
    };

    const resetProgress = () => {
        progress.classList.remove(CLASS_NAMES.progressVisible);
        setProgress(0);
        progressText.textContent = "Uploading&";
    };

    const resetState = () => {
        form.classList.remove(
            CLASS_NAMES.formFileSelected,
            CLASS_NAMES.formUploading,
            CLASS_NAMES.formUploadFailed,
            CLASS_NAMES.formUploaded
        );

        setUploading(false);
        resetProgress();
        setStatus("");
        updateDone();
    };

    const handleFailure = (message) => {
        job = null;
        uploaded = false;

        // Allows picking the same file again; otherwise "change" won't fire.
        input.value = "";

        setUploading(false);

        form.classList.remove(CLASS_NAMES.formUploaded);
        form.classList.add(CLASS_NAMES.formUploadFailed);

        progress.classList.remove(CLASS_NAMES.progressVisible);
        setStatus(message, true);
        updateDone();
    };

    /* ---------- Text fields ---------- */

    const bindTextField = (field) => {
        releaseNativeLimit(field);

        const wrapper = field.parentElement;
        const counter = $(wrapper, CLASS_NAMES.charCounter);
        const maxLength = getMaxLength(field);

        const sync = () => {
            const overLimit = isOverLimit(field);

            if (counter) {
                counter.textContent = Number.isFinite(maxLength)
                    ? `${field.value.length}/${maxLength}`
                    : String(field.value.length);
            }

            // Highlight the field in red while the text is too long.
            wrapper.classList.toggle(CLASS_NAMES.fieldError, overLimit);
            field.setAttribute("aria-invalid", String(overLimit));

            // Auto-resize the textarea (the CSS max-height caps the growth).
            // Resetting the height to "auto" resets the inner scroll position,
            // so remember it and restore it after measuring.
            const scrollTop = field.scrollTop;

            field.style.height = "auto";
            field.style.height = `${field.scrollHeight}px`;
            field.scrollTop = scrollTop;

            updateDone();
        };

        field.addEventListener("input", sync);
        sync();

        return sync;
    };

    const syncTitle = bindTextField(titleField);
    const syncDescription = bindTextField(descriptionField);

    // Name is a single-line field: don't allow line breaks (typing or pasting).
    titleField.addEventListener("keydown", (event) => {
        if (event.key === "Enter") {
            event.preventDefault();
        }
    });

    titleField.addEventListener("input", () => {
        if (titleField.value.includes("\n")) {
            titleField.value = titleField.value.replace(/\s*\n+\s*/g, " ");
            syncTitle();
        }
    });

    // The form must never be submitted the regular way.
    form.addEventListener("submit", (event) => event.preventDefault());

    /* ---------- Upload flow ---------- */

    const uploadFile = (file) => {
        if (uploading || !file) {
            return;
        }

        form.classList.remove(CLASS_NAMES.formUploadFailed, CLASS_NAMES.formUploaded);
        form.classList.add(CLASS_NAMES.formFileSelected);
        progress.classList.add(CLASS_NAMES.progressVisible);

        uploaded = false;
        setStatus("");
        setProgress(0);
        setUploading(true);
        updateDone();

        job = uploader.start(file, {
            onProgress: setProgress,

            onSent: () => {
                progressText.textContent = "Processing&";
            },

            onSuccess: () => {
                job = null;
                uploaded = true;

                setProgress(100);
                setUploading(false);

                form.classList.remove(CLASS_NAMES.formUploadFailed);
                form.classList.add(CLASS_NAMES.formUploaded);

                setStatus("File uploaded successfully.");
                updateDone();
            },

            onError: handleFailure,
        });
    };

    const selectFile = (file) => {
        if (uploading || !file) {
            return;
        }

        if (!isMediaFile(file)) {
            input.value = "";
            setStatus("Please select an audio or video file.", true);
            return;
        }

        fileName.textContent = file.name;
        fileSize.textContent = formatFileSize(file.size);

        // Prefill Name from the file name if the user hasn't typed anything yet.
        const shouldPrefillTitle = titleField.value.trim() === "";

        if (shouldPrefillTitle) {
            titleField.value = stripExtension(file.name).slice(0, getMaxLength(titleField));
        }

        // Start the upload first: it reveals the metadata block,
        // so the textarea can be measured correctly (hidden elements have scrollHeight 0).
        uploadFile(file);

        if (shouldPrefillTitle) {
            syncTitle();
        }
    };

    // Resets only the file and upload state; metadata is kept.
    const clearFile = () => {
        job?.abort();
        job = null;
        uploaded = false;

        input.value = "";
        fileName.textContent = "";
        fileSize.textContent = "";

        dropzone.classList.remove(CLASS_NAMES.dropzoneDragover);
        resetState();
    };

    // Full dialog reset (used on close).
    const resetAll = () => {
        clearFile();

        form.reset(); // title, description, visibility (back to the template's selected option)
        syncTitle();
        syncDescription();
        updateDone();
    };

    /* ---------- Events ---------- */

    chooseButton.addEventListener("click", () => {
        if (!uploading) {
            input.click();
        }
    });

    input.addEventListener("change", () => {
        selectFile(input.files?.[0]);
    });

    removeButton.addEventListener("click", () => {
        if (!uploading) {
            clearFile();
        }
    });

    dropzone.addEventListener("dragenter", (event) => {
        event.preventDefault();

        if (!uploading) {
            dropzone.classList.add(CLASS_NAMES.dropzoneDragover);
        }
    });

    dropzone.addEventListener("dragover", (event) => {
        event.preventDefault();

        if (!uploading) {
            dropzone.classList.add(CLASS_NAMES.dropzoneDragover);
        }
    });

    dropzone.addEventListener("dragleave", (event) => {
        // dragleave also fires when moving onto child elements; that's not leaving the dropzone.
        if (!dropzone.contains(event.relatedTarget)) {
            dropzone.classList.remove(CLASS_NAMES.dropzoneDragover);
        }
    });

    dropzone.addEventListener("drop", (event) => {
        event.preventDefault();
        dropzone.classList.remove(CLASS_NAMES.dropzoneDragover);

        if (!uploading) {
            selectFile(event.dataTransfer?.files?.[0]);
        }
    });

    doneButton.addEventListener("click", () => {
        if (doneButton.disabled) {
            return;
        }

        // The dialog's "close" event is dispatched asynchronously, after this handler,
        // so the form data hasn't been reset yet.
        onDone?.({
            file: input.files?.[0] ?? null,
            title: titleField.value.trim(),
            description: descriptionField.value.trim(),
            visibility: visibilityField.value,
        });
    });

    // Fires on Cancel, Done, the close button and Esc.
    dialog.addEventListener("close", resetAll);

    // Initial state.
    resetState();
}
