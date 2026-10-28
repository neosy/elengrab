const DOWNLOADER_EVENTS_URL = "/downloader/events";
const RECONNECT_DELAY_MS = 5000;

let lastEventId = null;
let isReconnect = false;

const CLASS_NAMES = {
    serverStatusOnline: "online",
}

export const DOM_IDS = {
    serverStatusDot: "server-status-dot",
}

const DOM_ELEMENTS = {
    serverStatusDot: null,
}

export function initDomElements(elements = DOM_ELEMENTS) {
    elements.serverStatusDot = document.getElementById(DOM_IDS.serverStatusDot);
}

/**
 * @param {Record<string, EventListener>} eventHandlers
 * @param {string} usernameFilter
 */
export function initSSE(eventHandlers, usernameFilter = "") {
    let sse = null;

    initDomElements();

    const connect = () => {
        if (!sse) {
            sse = createSSEConnection(eventHandlers, usernameFilter);
        }
    };

    window.addEventListener("pageshow", connect);

    // Close the connection when leaving the page so the server slot is released
    window.addEventListener("pagehide", () => {
        sse?.close();
        sse = null;
    });

    document.addEventListener("visibilitychange", () => {
        if (document.visibilityState === "visible") {
            connect();
        }
    });

    // Connect immediately in case "pageshow" has already fired
    connect();
}

function setServerStatus(online) {
    const el = DOM_ELEMENTS.serverStatusDot;
    if (!el) return;

    if (online) {
        console.info("SSE connection opened");
        el.classList.add(CLASS_NAMES.serverStatusOnline);
    } else {
        console.warn("SSE connection closed");
        el.classList.remove(CLASS_NAMES.serverStatusOnline);
    }
}

// Remember the last received event id (events without id keep the previous value).
// With reset, an event without id clears the stored value (used for "connected").
function trackEventId(event, reset = false) {
    if (event.lastEventId) {
        lastEventId = event.lastEventId;
    } else if (reset) {
        lastEventId = null;
    }
}

/**
 * @param {Record<string, EventListener>} eventHandlers
 */
function createSSEConnection(eventHandlers, usernameFilter = "") {
    const state = {
        eventSource: null,
        reconnectTimer: null,

        // Stop the timer and close the connection
        clear() {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;

            this.eventSource?.close();
            this.eventSource = null;
        },

        // Mark the connection as lost: log, set offline status and close the source
        markLost(err) {
            if (this.eventSource && this.eventSource.readyState !== EventSource.CLOSED) {
                console.error("SSE connection lost:", err);
                setServerStatus(false);
                this.eventSource.close();
            }
        },

        // Reconnect after delay (cancelable, single timer per connection)
        scheduleReconnect(connect) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = setTimeout(connect, RECONNECT_DELAY_MS);
        },
    };

    // Internal function to (re)connect
    function connect() {
        state.clear();

        const params = new URLSearchParams();
        if (lastEventId) params.set("lastEventId", lastEventId);
        if (isReconnect) params.set("reconnect", "true");

        const query = params.toString();
        const url = query ? `${DOWNLOADER_EVENTS_URL}?${query}` : DOWNLOADER_EVENTS_URL;

        state.eventSource = new EventSource(url);

        // Server is considered online when these events arrive
        state.eventSource.addEventListener("connected", (event) => {
            trackEventId(event, true);
            isReconnect = true;
            setServerStatus(true);
        });
        // state.eventSource.addEventListener("ping", () => setServerStatus(true));

        // Business events
        Object.entries(eventHandlers).forEach(([eventName, handler]) => {
            state.eventSource.addEventListener(eventName, (event) => {
                trackEventId(event);
                handler(event, usernameFilter);
            });
        });

        // Fallback: any default message marks server as online
        state.eventSource.onmessage = () => setServerStatus(true);

        // On error: mark offline and reconnect
        state.eventSource.onerror = function(err) {
            state.markLost(err);
            state.scheduleReconnect(connect);
        };
    }

    // Initial connection
    connect();

    // Return only API to close connection from outside
    return {
        close: () => {
            state.clear();
            setServerStatus(false);
        }};
}