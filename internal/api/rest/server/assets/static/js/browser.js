let stableHeight = 0;

export function initViewportHeightVar() {
    syncViewportHeightVar();

    let rafId = null;

    const scheduleSync = () => {
        if (rafId !== null) cancelAnimationFrame(rafId);

        rafId = requestAnimationFrame(() => {
            syncViewportHeightVar();
            rafId = null;
        });
    };

    function handleVisibility() {
        if (!document.hidden) {
            scheduleSync();
        }
    }

    function handleOrientationChange() {
        stableHeight = 0;
        scheduleSync();
    }

    function destroy() {
        window.removeEventListener('resize', scheduleSync);
        window.visualViewport?.removeEventListener('resize', scheduleSync);
        window.removeEventListener('orientationchange', handleOrientationChange);
        document.removeEventListener('visibilitychange', handleVisibility);
        window.removeEventListener('pageshow', scheduleSync);
    }

    window.addEventListener('resize', scheduleSync);
    window.visualViewport?.addEventListener('resize', scheduleSync);
    window.addEventListener('orientationchange', handleOrientationChange);
    document.addEventListener('visibilitychange', handleVisibility);
    window.addEventListener('pageshow', scheduleSync);

    return destroy;
}

function syncViewportHeightVar() {
    const layoutHeight = window.innerHeight;
    const visualHeight = window.visualViewport?.height || layoutHeight;

    stableHeight = Math.max(
        stableHeight,
        layoutHeight,
        visualHeight
    );

    document.documentElement.style.setProperty(
        '--vh-stable',
        `${stableHeight * 0.01}px`
    );

    document.documentElement.style.setProperty(
        '--vh-layout',
        `${layoutHeight * 0.01}px`
    );

    document.documentElement.style.setProperty(
        '--vh-visual',
        `${visualHeight * 0.01}px`
    );
}