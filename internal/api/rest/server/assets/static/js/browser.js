export function initViewportHeightVar() {
    syncViewportHeightVar();

    const onChange = () => syncViewportHeightVar();

    window.addEventListener('resize', onChange);
    window.visualViewport?.addEventListener('resize', onChange);
    window.addEventListener('orientationchange', onChange);
    document.addEventListener('visibilitychange', handleVisibility);
    window.addEventListener('pageshow', onChange);

    return function destroy() {
        window.removeEventListener('resize', onChange);
        window.visualViewport?.removeEventListener('resize', onChange);
        window.removeEventListener('orientationchange', onChange);
        document.removeEventListener('visibilitychange', handleVisibility);
        window.removeEventListener('pageshow', onChange);
    };
}

function syncViewportHeightVar() {
    const height = window.visualViewport?.height || window.innerHeight;
    document.documentElement.style.setProperty('--vh', `${height * 0.01}px`);
}

function handleVisibility() {
    if (!document.hidden) {
    requestAnimationFrame(onChange);
    }
}
