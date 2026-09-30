const sysInfoServerStatusDot = document.getElementById("server-status-dot");

let globalEventSource = null;

/**
 * @param {Record<string, EventListener>} eventHandlers
 */
export function initSSE(eventHandlers) {
    let sse = null;

    const connect = () => {
        if (!globalEventSource || globalEventSource.readyState === EventSource.CLOSED) {
            sse = createSSEConnection(eventHandlers);
        }
    };

    window.addEventListener("pageshow", connect);

    window.addEventListener("beforeunload", () => {
        if (!globalEventSource || globalEventSource.readyState === EventSource.CLOSED) {
            sse = createSSEConnection(eventHandlers);
        }
    });

    document.addEventListener("visibilitychange", () => {
        if (document.visibilityState === "visible") {
            connect();
        }
    });
}

/**
 * @param {Record<string, EventListener>} eventHandlers
 */
function createSSEConnection(eventHandlers) {
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
        globalEventSource?.close();

        globalEventSource = new EventSource("/downloader/events");

        // Server is considered online when these events arrive
        globalEventSource.addEventListener("connected", () => setServerStatus(true));
        // globalEventSource.addEventListener("ping", () => setServerStatus(true));

        // Business events
        Object.entries(eventHandlers).forEach(([eventName, handler]) => {
            globalEventSource.addEventListener(eventName, handler);
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

            // Reconnect after delay
            setTimeout(connect, 5000);
        };
    }

    // Initial connection
    connect();

    // Return only API to close connection from outside
    return {
        close: () => {
            setServerStatus(false);
            globalEventSource?.close();
        }};
}
