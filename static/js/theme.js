// Dark Mode Theme Manager
class ThemeManager {
    constructor() {
        this.theme = localStorage.getItem('theme') || 'light';
        this.init();
    }

    init() {
        this.applyTheme();
        this.createToggleButton();
        // Update header button multiple times to ensure it works
        this.updateHeaderButton();
        setTimeout(() => this.updateHeaderButton(), 50);
        setTimeout(() => this.updateHeaderButton(), 200);
    }

    applyTheme() {
        if (this.theme === 'dark') {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
        localStorage.setItem('theme', this.theme);
    }

    toggle() {
        this.theme = this.theme === 'light' ? 'dark' : 'light';
        this.applyTheme();
        this.updateToggleButton();
        this.updateHeaderButton();
    }

    createToggleButton() {
        const button = document.createElement('button');
        button.id = 'theme-toggle';
        button.className = 'fixed bottom-6 right-6 w-12 h-12 bg-primary hover:bg-indigo-600 text-white rounded-full shadow-lg transition-all z-50 flex items-center justify-center text-lg';
        button.innerHTML = this.theme === 'dark' ? '☀️' : '🌙';
        button.onclick = () => this.toggle();
        document.body.appendChild(button);
    }

    updateToggleButton() {
        const button = document.getElementById('theme-toggle');
        if (button) {
            button.innerHTML = this.theme === 'dark' ? '☀️' : '🌙';
        }
    }
    
    updateHeaderButton() {
        const icon = document.getElementById('theme-icon');
        if (icon) {
            icon.textContent = this.theme === 'dark' ? '☀️' : '🌙';
        }
    }
    
    // Global toggle function
    toggleTheme() {
        this.toggle();
    }
}

// Initialize theme manager
const themeManager = new ThemeManager();

// Global function for manual toggle
function toggleTheme() {
    themeManager.toggle();
}

// Ensure theme is applied immediately and on page load
themeManager.updateHeaderButton();
document.addEventListener('DOMContentLoaded', () => {
    themeManager.updateHeaderButton();
    // Force update after a short delay to ensure DOM is ready
    setTimeout(() => themeManager.updateHeaderButton(), 100);
});