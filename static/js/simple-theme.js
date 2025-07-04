// Simple, working dark mode implementation
(function() {
    // Apply theme immediately
    function applyTheme() {
        const theme = localStorage.getItem('theme') || 'light';
        if (theme === 'dark') {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
        updateIcon();
    }

    // Update theme icon
    function updateIcon() {
        const theme = localStorage.getItem('theme') || 'light';
        const icon = document.getElementById('theme-icon');
        if (icon) {
            icon.textContent = theme === 'dark' ? '☀️' : '🌙';
        }
    }

    // Toggle theme function
    window.toggleTheme = function() {
        const currentTheme = localStorage.getItem('theme') || 'light';
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        
        localStorage.setItem('theme', newTheme);
        applyTheme();
        
        console.log('Theme switched to:', newTheme);
    };

    // Apply theme immediately when script loads
    applyTheme();

    // Apply theme when DOM is ready
    document.addEventListener('DOMContentLoaded', function() {
        applyTheme();
        // Update icon after a short delay to ensure DOM is ready
        setTimeout(updateIcon, 100);
    });

    // Update icon when page is fully loaded
    window.addEventListener('load', updateIcon);
})();