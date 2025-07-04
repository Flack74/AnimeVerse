// Enhanced UI Components and Animations
(function() {
    // Add loading skeleton animation
    function createSkeleton(count = 6) {
        return Array(count).fill(0).map(() => `
            <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden animate-pulse">
                <div class="bg-gray-300 dark:bg-gray-600 h-64 w-full"></div>
                <div class="p-4 space-y-2">
                    <div class="bg-gray-300 dark:bg-gray-600 h-4 rounded w-3/4"></div>
                    <div class="bg-gray-300 dark:bg-gray-600 h-3 rounded w-1/2"></div>
                </div>
            </div>
        `).join('');
    }

    // Enhanced toast notifications
    function showToast(message, type = 'success') {
        const toast = document.createElement('div');
        toast.className = `fixed top-4 right-4 z-50 px-6 py-3 rounded-lg shadow-lg transform translate-x-full transition-transform duration-300 ${
            type === 'success' ? 'bg-green-500 text-white' : 
            type === 'error' ? 'bg-red-500 text-white' : 
            'bg-blue-500 text-white'
        }`;
        toast.textContent = message;
        
        document.body.appendChild(toast);
        setTimeout(() => toast.classList.remove('translate-x-full'), 100);
        setTimeout(() => {
            toast.classList.add('translate-x-full');
            setTimeout(() => document.body.removeChild(toast), 300);
        }, 3000);
    }

    // Smooth scroll to top
    function addScrollToTop() {
        const button = document.createElement('button');
        button.innerHTML = '↑';
        button.className = 'fixed bottom-20 right-6 w-12 h-12 bg-primary text-white rounded-full shadow-lg opacity-0 transition-all duration-300 hover:bg-indigo-600 z-40';
        button.onclick = () => window.scrollTo({ top: 0, behavior: 'smooth' });
        
        document.body.appendChild(button);
        
        window.addEventListener('scroll', () => {
            if (window.scrollY > 300) {
                button.classList.remove('opacity-0');
                button.classList.add('opacity-100');
            } else {
                button.classList.add('opacity-0');
                button.classList.remove('opacity-100');
            }
        });
    }

    // Enhanced search with debounce
    function enhanceSearch() {
        const searchInputs = document.querySelectorAll('input[type="text"]');
        searchInputs.forEach(input => {
            let timeout;
            input.addEventListener('input', (e) => {
                clearTimeout(timeout);
                timeout = setTimeout(() => {
                    if (e.target.value.length > 2) {
                        // Trigger search
                        const event = new CustomEvent('enhancedSearch', { detail: e.target.value });
                        document.dispatchEvent(event);
                    }
                }, 300);
            });
        });
    }

    // Add intersection observer for animations
    function addScrollAnimations() {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.classList.add('animate-fadeInUp');
                }
            });
        }, { threshold: 0.1 });

        document.querySelectorAll('.anime-card, .content-section').forEach(el => {
            observer.observe(el);
        });
    }

    // Global functions
    window.showToast = showToast;
    window.createSkeleton = createSkeleton;

    // Initialize enhancements
    document.addEventListener('DOMContentLoaded', () => {
        addScrollToTop();
        enhanceSearch();
        setTimeout(addScrollAnimations, 100);
    });

    // Add CSS animations
    const style = document.createElement('style');
    style.textContent = `
        @keyframes fadeInUp {
            from { opacity: 0; transform: translateY(30px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .animate-fadeInUp { animation: fadeInUp 0.6s ease-out forwards; }
        .glass-effect { backdrop-filter: blur(10px); background: rgba(255, 255, 255, 0.1); }
        .dark .glass-effect { background: rgba(0, 0, 0, 0.2); }
    `;
    document.head.appendChild(style);
})();