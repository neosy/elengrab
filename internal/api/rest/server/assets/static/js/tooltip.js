let tooltip, tooltipTimer;

function createTooltip() {
    tooltip = document.createElement("div");
    tooltip.className = "tooltip";
    document.body.appendChild(tooltip);
}

function showTooltip(el) {
    const text = el.dataset.tooltip || el.title;
    if (!text) return;

    // clear previous timer
    if (tooltipTimer) {
        clearTimeout(tooltipTimer);
        tooltipTimer = null;
    }

    tooltip.innerHTML = text.replace(/\n\s*/g, "<br>");

    const rect = el.getBoundingClientRect();

    // remove class first to restart animation
    tooltip.classList.remove("show");

    // temporarily set top to 0 to measure height correctly
    tooltip.style.top = "0px";
    tooltip.style.left = rect.left + rect.width / 2 + "px";

    // get tooltip height dynamically
    const tooltipHeight = tooltip.offsetHeight;

    // position tooltip above element with margin
    tooltip.style.top = rect.top - tooltipHeight - 8 + window.scrollY + "px"; // 8px gap above

    // wait next frame to add 'show' class for smooth animation
    requestAnimationFrame(() => {
        requestAnimationFrame(() => {
            tooltip.classList.add("show");
        });
    });

    // hide tooltip after 2s
    tooltipTimer = setTimeout(() => {
        tooltip.classList.remove("show");
    }, 2000);
}

export function initTooltips(selector = "[data-tooltip]") {
    createTooltip();

    document.body.addEventListener("click", (event) => {
        const el = event.target.closest(selector);
        if (!el) return;

        showTooltip(el);
    });
}