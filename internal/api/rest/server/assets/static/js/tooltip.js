let tooltip;

function createTooltip() {
    tooltip = document.createElement("div");
    tooltip.className = "tooltip";
    document.body.appendChild(tooltip);
}

function showTooltip(el) {
    const text = el.dataset.tooltip || el.title;
    if (!text) return;

    tooltip.textContent = text;

    const rect = el.getBoundingClientRect();

    tooltip.style.left = rect.left + rect.width / 2 + "px";
    tooltip.style.top = rect.top - 30 + window.scrollY + "px";

    tooltip.classList.add("show");

    setTimeout(() => {
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