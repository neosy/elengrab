import * as notify from './notifications.js';

export function initHeaderAutoHide() {
    let lastScrollTop = 0;
    let ticking = false;
    const header = document.getElementById('header');
    const hiddenClass = "header--hidden";

    function updateHeader() {
        const scrollTop = window.scrollY || document.documentElement.scrollTop;

        if (scrollTop > lastScrollTop && scrollTop > 250) {
            header.classList.add(hiddenClass);
        } else {
            header.classList.remove(hiddenClass);
        }

        lastScrollTop = Math.max(scrollTop, 0);
        ticking = false;
    }

    function handleScroll() {
        if (!ticking) {
            requestAnimationFrame(updateHeader);
            ticking = true;
        }
    }

    window.addEventListener('scroll', handleScroll, { passive: true });
}

export function initHTMX() {
    // Guest session created
    htmx.on("guestCreated", (event) => {
        if (event.detail.value === true) {
            // Reload to apply new session (cookie)
            window.location.reload();
        }
    });

    document.body.addEventListener("htmx:responseError", (event) => {
        const xhr = event.detail.xhr;

        // Error handling (HTTP >= 400, except 503)
        if (xhr.status >= 400 && xhr.status !== 503) {
            try {
                const data = JSON.parse(xhr.responseText);

                if (data && typeof data === "object" && "message" in data) {
                    notify.show(data.message, notify.notifyType.ERROR);
                }
            } catch (e) {
                // Ignore non-JSON responses.
            }
        }
    });
}