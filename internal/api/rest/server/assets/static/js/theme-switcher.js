;(() => {
    const html = document.documentElement;
    const btn = document.getElementById('theme-toggle-btn');

    // Apply theme and save to localStorage
    const setTheme = (theme) => {
        html.classList.toggle('dark', theme === 'dark');
        localStorage.setItem('theme', theme);
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

    // Follow system preference changes only if user hasn't made a manual choice
    window.matchMedia('(prefers-color-scheme: dark)')
        .addEventListener('change', (e) => {
            if (!localStorage.getItem('theme')) {
                setTheme(e.matches ? 'dark' : 'light');
            }
        });
})();