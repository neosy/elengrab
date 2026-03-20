// -------------------------------------------------------------
// Handle row-add SSE event with multiple rows in one payload
// and insert them after the top placeholder div
// -------------------------------------------------------------
export function handleRowAdd(event) {
    try {
        // Payload is raw HTML with multiple rows
        const data = JSON.parse(event.data);
        if (!data.id || !data.html) return

        // Find the container that holds all rows
        const container = document.getElementById("result-rows");
        if (!container) return;

        // Find the top placeholder div
        const placeholder = document.getElementById("row-top-placeholder");
        if (!placeholder) return;

        // Insert all rows right after the placeholder
        placeholder.insertAdjacentHTML("afterend", data.html);

        const newEl = document.getElementById(data.id);
        if (newEl) {
            // If the element exists, remove it from the DOM
            const elNoItems = document.getElementById("row-no-items");
            if (elNoItems) {
                elNoItems.remove();
            }

            // Let HTMX process this element and activate any hx-* attributes
            htmx.process(newEl);
            
            // Add animation class and initial styles
            newEl.classList.add('animated');
            newEl.style.opacity = 0;
            newEl.style.transform = 'translateY(-20px)';
            newEl.style.transition = 'opacity 0.75s ease, transform 0.75s ease';

            // Trigger animation on the next browser frame
            requestAnimationFrame(() => {
                newEl.style.opacity = 1;
                newEl.style.transform = 'translateY(0)';
            });
        }
    } catch (err) {
        console.error("SSE row-add handler error:", err);
    }
}

export function handleRowUpdate(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.id || !data.html) return

        const el = document.getElementById(data.id);
        if (!el) return;

        const temp = document.createElement("div");
        temp.innerHTML = data.html.trim();

        const newEl = temp.firstElementChild;
        if (!newEl) return;

        el.replaceWith(newEl);
        //HTMX will scan this element and activate all hx-* attributes
        htmx.process(newEl);
    } catch (err) {
        console.error("SSE row-update handler error:", err);
    }
}

// -------------------------------------------------------------
// Handle row-delete SSE event and remove element if it exists
// -------------------------------------------------------------
export function handleRowDelete(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.id) return;

        const el = document.getElementById(data.id);
        if (!el) return;

        // initial styles for animation
        el.style.transition = "opacity 0.4s ease, transform 0.4s ease, height 0.4s ease, margin 0.4s ease";
        el.style.opacity = "1";
        el.style.transform = "translateY(0)";
        el.style.height = el.offsetHeight + "px";

        // forcing layout so that the browser applies height
        el.offsetHeight;

        // final state
        el.style.opacity = "0";
        el.style.transform = "translateY(-10px)";
        el.style.height = "0";
        el.style.marginTop = "0";
        el.style.marginBottom = "0";
        el.style.paddingTop = "0";
        el.style.paddingBottom = "0";

        el.addEventListener("transitionend", () => {
            el.remove();
        }, { once: true });
    } catch (err) {
        console.error("SSE row-delete handler error:", err);
    }
}

export function handleSystemInfoUpdate(event) {
    try {
        const data = JSON.parse(event.data);
        if (!data.diskFree || !data.diskUsed) return

        const elDiskFree = document.getElementById("disk-free");
        const elDiskUsed = document.getElementById("disk-used");
        if (!elDiskFree || !elDiskUsed) return;

        elDiskFree.textContent = data.diskFree
        elDiskUsed.textContent = data.diskUsed
    } catch (err) {
        console.error("SSE row-update handler error:", err);
    }
}
