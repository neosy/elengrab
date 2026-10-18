import * as notify from "./notifications.js";

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

// Containers that can hold browser-playable codecs
// (H.264/VP9/AV1 video; AAC/Opus/MP3/FLAC/Vorbis/PCM audio).
// Codecs themselves are validated on the server.
const MEDIA_EXTENSIONS =
    /\.(mp4|m4v|m4a|m4b|mov|3gp|3gpp|3g2|webm|weba|mkv|ogg|oga|ogv|opus|mp3|mpga|aac|flac|wav)$/i;

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

function fitFileName(name, maxLength) {
    if (name.length <= maxLength) {
        return name;
    }

    const dot = name.lastIndexOf(".");
    const extension = dot > 0 ? name.slice(dot) : "";

    return name.slice(0, maxLength - extension.length) + extension;
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
/* Interface: start(file, callbacks) -> { abort }, queue(fileId, data) -> Promise<{ itemId }>, remove(fileId) -> Promise */
/* callbacks: onProgress(percent), onSent(), onSuccess({ fileId }), onError(msg) */
/* ------------------------------------------------------------------ */

function createXhrUploader(url) {
    const fileUrl = (fileId) => `${url.replace(/\/+$/, "")}/${encodeURIComponent(fileId)}`;

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
                    let fileId = null;

                    try {
                        const response = JSON.parse(xhr.responseText);

                        if (typeof response.fileId === "string" && response.fileId) {
                            fileId = response.fileId;
                        }
                    } catch {
                        // Handled below: the response is not valid JSON.
                    }

                    if (fileId === null) {
                        onError("The server returned an unexpected response.");
                        return;
                    }

                    onSuccess({ fileId });
                    return;
                }

                let message = "Upload failed. Please try again.";

                try {
                    const response = JSON.parse(xhr.responseText);

                    if (typeof response.message === "string" && response.message) {
                        message = response.message;
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

        async queue(fileId, { title, description, visibility }) {
            let response;

            try {
                response = await fetch(fileUrl(fileId), {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        Accept: "application/json",
                    },
                    body: JSON.stringify({ title, description, visibility }),
                });
            } catch {
                throw new Error("A network error occurred. Please try again.");
            }

            const body = await response.json().catch(() => null);

            if (!response.ok) {
                const message = body?.message;

                throw new Error(
                    typeof message === "string" && message
                        ? message
                        : "Could not save the file. Please try again."
                );
            }

            if (typeof body?.itemId !== "string" || !body.itemId) {
                throw new Error("The server returned an unexpected response.");
            }

            return { itemId: body.itemId };
        },

        async remove(fileId) {
            try {
                // keepalive lets the request finish even if the page is closing.
                await fetch(fileUrl(fileId), { method: "DELETE", keepalive: true });
            } catch {
                // Best effort: the user can't do anything about a failed cleanup.
            }
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
                        onSuccess({ fileId: `mock-${Date.now()}` });
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

        queue() {
            return new Promise((resolve, reject) => {
                setTimeout(() => {
                    if (Math.random() < failRate) {
                        reject(new Error("Mock save failed."));
                    } else {
                        resolve({ itemId: `mock-item-${Date.now()}` });
                    }
                }, 600);
            });
        },

        remove() {
            return Promise.resolve();
        },        
    };
}

/* ------------------------------------------------------------------ */
/* Dialog                                                              */
/* ------------------------------------------------------------------ */

/**
 * @param {string} [uploadUrl] Upload URL. If omitted (or mock: true), the mock uploader is used.
 * @param {object} [options]
 * @param {(result: {itemId: string, title: string, description: string, visibility: string}) => void} [options.onDone]
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
    const closeButtons = dialog.querySelectorAll("[data-dialog-close]");

    const uploader =
        mock || !uploadUrl
            ? createMockUploader({ failRate: mockFailRate })
            : createXhrUploader(uploadUrl);

    let uploading = false;
    let submitting = false;
    let uploaded = false;
    let fileId = null;
    let job = null;

    const resetUploadResult = () => {
        uploaded = false;
        fileId = null;
    };

    const setUploadResult = (id) => {
        uploaded = true;
        fileId = id;
    };

    const discardUploadedFile = () => {
        const id = fileId;

        resetUploadResult();

        if (id !== null) {
            uploader.remove(id);
        }
    };

    /* ---------- UI helpers ---------- */

    const updateDone = () => {
        doneButton.disabled = !(
            uploaded &&
            !uploading &&
            !submitting &&
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

        if (isError) {
            notify.show(message, notify.notifyType.ERROR);
        }
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

    const setSubmitting = (value) => {
        submitting = value;

        input.disabled = value;
        chooseButton.disabled = value;
        removeButton.disabled = value;

        closeButtons.forEach((button) => {
            button.disabled = value;
        });

        updateDone();
    };

    const resetProgress = () => {
        progress.classList.remove(CLASS_NAMES.progressVisible);
        setProgress(0);
        progressText.textContent = "Uploading...";
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
        resetUploadResult();

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

        discardUploadedFile();
        setStatus("");
        setProgress(0);
        setUploading(true);
        updateDone();

        job = uploader.start(file, {
            onProgress: setProgress,

            onSent: () => {
                progressText.textContent = "Processing...";
            },

            onSuccess: ({ fileId: id }) => {
                job = null;
                setUploadResult(id);

                setProgress(100);
                setUploading(false);

                form.classList.remove(CLASS_NAMES.formUploadFailed);
                form.classList.add(CLASS_NAMES.formUploaded);

                // setStatus("File uploaded successfully.");
                updateDone();
            },

            onError: handleFailure,
        });
    };

    const selectFile = (file) => {
        if (uploading || submitting || !file) {
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
            titleField.value = fitFileName(file.name, getMaxLength(titleField));
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
        discardUploadedFile();

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
            resetAll();
        }
    });

    dropzone.addEventListener("dragenter", (event) => {
        event.preventDefault();

        if (!uploading && !submitting) {
            dropzone.classList.add(CLASS_NAMES.dropzoneDragover);
        }
    });

    dropzone.addEventListener("dragover", (event) => {
        event.preventDefault();

        if (!uploading && !submitting) {
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

    doneButton.addEventListener("click", async (event) => {
        if (doneButton.disabled || submitting) {
            return;
        }

        // Capture the values now: the user could edit the fields while the request is in flight.
        const data = {
            title: titleField.value.trim(),
            description: descriptionField.value.trim(),
            visibility: visibilityField.value,
        };

        setSubmitting(true);
        setStatus("");

        let itemId;

        try {
            ({ itemId } = await uploader.queue(fileId, data));
        } catch (error) {
            console.error("Failed to queue uploaded file", error);
            notify.show(error.message, notify.notifyType.ERROR);
            return;
        } finally {
            setSubmitting(false);
        }

        // The file is now owned by the queued item: closing must not delete it.
        resetUploadResult();
        dialog.close();
        onDone?.({ itemId, ...data });        
    });

    dialog.addEventListener("cancel", (event) => {
        if (submitting) {
            event.preventDefault();
        }
    });

    // Fires on Cancel, the close button, Esc, and after Done (dialog.close()).
    // An uploaded file that wasn't queued is deleted in clearFile.
    dialog.addEventListener("close", resetAll);

    window.addEventListener("pagehide", () => {
    if (fileId !== null && !submitting) {
        uploader.remove(fileId);
    }
    });

    // Initial state.
    resetState();
}
