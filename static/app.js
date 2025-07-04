// Animeverse Frontend JavaScript
// Handles image caching, modals, and dynamic interactions

class AnimeImageCache {
    constructor() {
        this.cache = new Map();
        this.pendingRequests = new Map();
    }

    async fetchAndCacheImages(malId, anilistId, animeElement) {
        const cacheKey = `${malId}-${anilistId}`;
        
        // Check local cache first
        if (this.cache.has(cacheKey)) {
            this.updateAnimeImages(animeElement, this.cache.get(cacheKey));
            return;
        }

        // Check if request is already pending
        if (this.pendingRequests.has(cacheKey)) {
            return this.pendingRequests.get(cacheKey);
        }

        // Create new request
        const request = this.performImageFetch(malId, anilistId, animeElement, cacheKey);
        this.pendingRequests.set(cacheKey, request);
        
        try {
            await request;
        } finally {
            this.pendingRequests.delete(cacheKey);
        }
    }

    async performImageFetch(malId, anilistId, animeElement, cacheKey) {
        try {
            // Check backend cache first
            const checkResponse = await fetch(`/api/images/check?mal_id=${malId}&anilist_id=${anilistId}`);
            const checkData = await checkResponse.json();
            
            if (checkData.success && checkData.data) {
                const images = checkData.data;
                this.cache.set(cacheKey, images);
                this.updateAnimeImages(animeElement, images);
                return;
            }

            // Fetch from external APIs
            const [jikanResponse, anilistResponse] = await Promise.allSettled([
                malId ? this.fetchJikanImage(malId) : Promise.reject('No MAL ID'),
                anilistId ? this.fetchAniListImage(anilistId) : Promise.reject('No AniList ID')
            ]);

            let imageUrl = '';
            let bannerUrl = '';

            if (jikanResponse.status === 'fulfilled') {
                imageUrl = jikanResponse.value;
            }

            if (anilistResponse.status === 'fulfilled') {
                bannerUrl = anilistResponse.value;
            }

            // Save to backend if we got any images
            if (imageUrl || bannerUrl) {
                const images = { image_url: imageUrl, banner_url: bannerUrl };
                
                await fetch('/api/images/save', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        mal_id: malId,
                        anilist_id: anilistId,
                        image_url: imageUrl,
                        banner_url: bannerUrl
                    })
                });

                this.cache.set(cacheKey, images);
                this.updateAnimeImages(animeElement, images);
            }
        } catch (error) {
            console.error('Error fetching images:', error);
        }
    }

    async fetchJikanImage(malId) {
        const response = await fetch(`https://api.jikan.moe/v4/anime/${malId}`);
        if (!response.ok) throw new Error('Jikan API error');
        
        const data = await response.json();
        return data.data?.images?.jpg?.large_image_url || '';
    }

    async fetchAniListImage(anilistId) {
        const query = `
            query ($id: Int) {
                Media (id: $id, type: ANIME) {
                    bannerImage
                    coverImage {
                        large
                    }
                }
            }
        `;

        const response = await fetch('https://graphql.anilist.co', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                query: query,
                variables: { id: anilistId }
            })
        });

        if (!response.ok) throw new Error('AniList API error');
        
        const data = await response.json();
        return data.data?.Media?.bannerImage || data.data?.Media?.coverImage?.large || '';
    }

    updateAnimeImages(element, images) {
        if (!element) return;
        
        const img = element.querySelector('img');
        if (img && images.image_url) {
            img.src = images.image_url;
            img.onerror = null; // Remove error handler after successful load
        }
    }
}

// Modal Management
class ModalManager {
    constructor() {
        this.currentModal = null;
        this.setupEventListeners();
    }

    setupEventListeners() {
        // Close modals on escape key
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                this.closeAllModals();
            }
        });

        // Close modals on backdrop click
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('modal-backdrop')) {
                this.closeAllModals();
            }
        });
    }

    async showAnimeModal(animeName) {
        try {
            const response = await fetch(`/api/anime/${animeName}`, {
                headers: { 'HX-Request': 'true' }
            });
            
            if (!response.ok) throw new Error('Failed to load anime details');
            
            const html = await response.text();
            const modalContent = document.getElementById('modal-content');
            const modal = document.getElementById('anime-modal');
            
            if (modalContent && modal) {
                modalContent.innerHTML = html;
                modal.classList.remove('hidden');
                modal.classList.add('flex');
                this.currentModal = 'anime-modal';
            }
        } catch (error) {
            console.error('Error loading anime modal:', error);
        }
    }

    showPlayer() {
        const modal = document.getElementById('player-modal');
        if (modal) {
            modal.classList.remove('hidden');
            modal.classList.add('flex');
            this.currentModal = 'player-modal';
        }
    }

    closeModal() {
        const modal = document.getElementById('anime-modal');
        if (modal) {
            modal.classList.add('hidden');
            modal.classList.remove('flex');
        }
        this.currentModal = null;
    }

    closePlayer() {
        const modal = document.getElementById('player-modal');
        const videoPlayer = document.getElementById('video-player');
        
        if (modal) {
            modal.classList.add('hidden');
            modal.classList.remove('flex');
        }
        
        if (videoPlayer) {
            videoPlayer.innerHTML = '';
        }
        
        this.currentModal = null;
    }

    closeAllModals() {
        this.closeModal();
        this.closePlayer();
    }
}

// Search and Filter Management
class SearchManager {
    constructor() {
        this.searchTimeout = null;
        this.setupSearchHandlers();
    }

    setupSearchHandlers() {
        const searchInput = document.getElementById('search-input');
        if (searchInput) {
            searchInput.addEventListener('input', () => {
                if (searchInput.value.trim() === '') {
                    // Hide search results when input is empty
                    const searchResultsSection = document.getElementById('search-results-section');
                    if (searchResultsSection) {
                        searchResultsSection.classList.add('hidden');
                    }
                } else {
                    this.debounceSearch();
                }
            });
            
            // Handle Enter key
            searchInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    e.preventDefault();
                    this.performSearch();
                }
            });
        }
    }

    debounceSearch() {
        clearTimeout(this.searchTimeout);
        this.searchTimeout = setTimeout(() => {
            this.performSearch();
        }, 500);
    }

    async performSearch() {
        const searchInput = document.getElementById('search-input');
        const genreFilter = document.getElementById('genre-filter');
        const yearFilter = document.getElementById('year-filter');
        const seasonFilter = document.getElementById('season-filter');
        
        const params = new URLSearchParams();
        
        if (searchInput?.value) params.append('search', searchInput.value);
        if (genreFilter?.value) params.append('genre', genreFilter.value);
        if (yearFilter?.value) params.append('year', yearFilter.value);
        if (seasonFilter?.value) params.append('season', seasonFilter.value);

        try {
            const response = await fetch(`/api/animes/filter?${params.toString()}`, {
                headers: { 'HX-Request': 'true' }
            });
            
            if (response.ok) {
                const html = await response.text();
                
                // If there's a search query, show results at the top
                if (searchInput?.value) {
                    const searchResultsSection = document.getElementById('search-results-section');
                    const searchResults = document.getElementById('search-results');
                    
                    if (searchResultsSection && searchResults) {
                        searchResults.innerHTML = html;
                        searchResultsSection.classList.remove('hidden');
                        searchResultsSection.scrollIntoView({ behavior: 'smooth' });
                    }
                } else {
                    // Hide search results section if no search query
                    const searchResultsSection = document.getElementById('search-results-section');
                    if (searchResultsSection) {
                        searchResultsSection.classList.add('hidden');
                    }
                    
                    // Show in main content for filters
                    const mainContent = document.getElementById('main-content');
                    if (mainContent) {
                        mainContent.innerHTML = html;
                    }
                }
            }
        } catch (error) {
            console.error('Search error:', error);
        }
    }
}

// Utility Functions
class Utils {
    static async randomAnime() {
        try {
            const response = await fetch('/api/animes/random');
            const data = await response.json();
            
            if (data.success && data.data) {
                const animeName = data.data.name.toLowerCase().replace(/ /g, '-');
                window.modalManager.showAnimeModal(animeName);
            }
        } catch (error) {
            console.error('Error fetching random anime:', error);
        }
    }

    static setCurrentDay() {
        const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
        const currentDayElement = document.getElementById('current-day');
        if (currentDayElement) {
            currentDayElement.textContent = days[new Date().getDay()];
        }
    }

    static async addAnimeToList(animeName, status) {
        // This would require authentication
        console.log(`Adding ${animeName} to ${status} list`);
        // Implementation would depend on your auth system
    }
}

// Global Functions (for inline event handlers)
function showAnimeModal(animeName) {
    window.modalManager.showAnimeModal(animeName);
}

function closeModal() {
    window.modalManager.closeModal();
}

function showPlayer() {
    window.modalManager.showPlayer();
}

function closePlayer() {
    window.modalManager.closePlayer();
}

function randomAnime() {
    Utils.randomAnime();
}

function addAnimeToList(animeName, status) {
    Utils.addAnimeToList(animeName, status);
}

function clearSearch() {
    const searchInput = document.getElementById('search-input');
    const searchResultsSection = document.getElementById('search-results-section');
    
    if (searchInput) {
        searchInput.value = '';
    }
    
    if (searchResultsSection) {
        searchResultsSection.classList.add('hidden');
    }
}

function performSearchFromButton() {
    if (window.searchManager) {
        window.searchManager.performSearch();
    }
}

function toggleMobileMenu() {
    const mobileMenu = document.getElementById('mobile-menu');
    const menuBtn = document.getElementById('mobile-menu-btn');
    
    if (mobileMenu && menuBtn) {
        mobileMenu.classList.toggle('hidden');
        
        // Update button icon
        const svg = menuBtn.querySelector('svg');
        const path = svg.querySelector('path');
        
        if (mobileMenu.classList.contains('hidden')) {
            path.setAttribute('d', 'M4 6h16M4 12h16M4 18h16');
        } else {
            path.setAttribute('d', 'M6 18L18 6M6 6l12 12');
        }
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    // Initialize managers
    window.imageCache = new AnimeImageCache();
    window.modalManager = new ModalManager();
    window.searchManager = new SearchManager();
    
    // Set current day
    Utils.setCurrentDay();
    
    console.log('🌸 Animeverse frontend initialized');
});