const sysInfoServerStatusDot = document.getElementById("server-status-dot");

let globalEventSource = null;
let reconnectTimer = null;

/**
 * @param {Record<string, EventListener>} eventHandlers
 * @param {string} usernameFilter
 */
export function initSSE(eventHandlers, usernameFilter = "") {
    let sse = null;

    const connect = () => {
        if (!globalEventSource || globalEventSource.readyState === EventSource.CLOSED) {
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

/**
 * @param {Record<string, EventListener>} eventHandlers
 */
function createSSEConnection(eventHandlers, usernameFilter = "") {
    function setServerStatus(online) {
        const el = sysInfoServerStatusDot;
        if (!el) return;

        if (online) {
            console.info("SSE connection opened");
            el.classList.add("online");
        } else {
            console.warn("SSE connection closed");
            el.classList.remove("online");
        }
    }

    // Internal function to (re)connect
    function connect() {
        clearTimeout(reconnectTimer);
        globalEventSource?.close();

        globalEventSource = new EventSource("/downloader/events");

        // Server is considered online when these events arrive
        globalEventSource.addEventListener("connected", () => setServerStatus(true));
        // globalEventSource.addEventListener("ping", () => setServerStatus(true));

        // Business events
        Object.entries(eventHandlers).forEach(([eventName, handler]) => {
            globalEventSource.addEventListener(eventName, (event) => {
                handler(event, usernameFilter);
            });
        });

        // Fallback: any default message marks server as online
        globalEventSource.onmessage = () => setServerStatus(true);

        // On error: mark offline and reconnect
        globalEventSource.onerror = function(err) {
            if (globalEventSource && globalEventSource.readyState !== EventSource.CLOSED) {
                console.error("SSE connection lost:", err);
                setServerStatus(false);
                globalEventSource?.close();
            }

            // Reconnect after delay (cancelable, single timer for the module)
            clearTimeout(reconnectTimer);
            reconnectTimer = setTimeout(connect, 5000);
        };
    }

    // Initial connection
    connect();

    // Return only API to close connection from outside
    return {
        close: () => {
            clearTimeout(reconnectTimer);
            setServerStatus(false);
            globalEventSource?.close();
            globalEventSource = null;
        }};
}