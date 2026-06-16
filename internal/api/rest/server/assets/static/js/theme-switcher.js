;(() => {
    const html = document.documentElement;
    const btn = document.getElementById('theme-toggle-btn');
    const themeColorMeta = document.getElementById('theme-color-meta');

    if (!btn || !themeColorMeta) return;

    // Apply theme, persist to localStorage and update mobile status bar color
    const setTheme = (theme) => {
        const isDarkTheme = theme === 'dark';
        
        html.classList.toggle('dark', isDarkTheme);
        localStorage.setItem('theme', theme);

        // Update <meta name="theme-color"> for mobile browsers (removes white bar)
        if (themeColorMeta) {
            const color = isDarkTheme ? '#171717' : '#f0f0f0'; // adjust to your exact colors
            themeColorMeta.setAttribute('content', color);
        }

        if (btn) {
            btn.title = isDarkTheme ? "Switch to light mode" : "Switch to dark mode";
        }
    };

    // Determine initial theme on page load
    const savedTheme = localStorage.getItem('theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const initialTheme = savedTheme || (prefersDark ? 'dark' : 'light');

    setTheme(initialTheme);

    // Toggle theme on button click
    btn.addEventListener('click', () => {
        const isDark = html.classList.contains('dark');
        setTheme(isDark ? 'light' : 'dark');
    });

    // Follow system theme changes only if user has not made a manual choice
    window.matchMedia('(prefers-color-scheme: dark)')
        .addEventListener('change', (e) => {
            if (!localStorage.getItem('theme')) {
                setTheme(e.matches ? 'dark' : 'light');
            }
        });
})();