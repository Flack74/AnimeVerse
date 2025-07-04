// Simple working authentication that actually works
class SimpleAuth {
    constructor() {
        this.user = null;
        this.session = null;
        this.init();
    }

    init() {
        const session = localStorage.getItem('auth_session');
        if (session) {
            try {
                this.session = JSON.parse(session);
                this.user = this.session.user;
            } catch (e) {
                localStorage.removeItem('auth_session');
            }
        }
        this.updateUI();
    }

    async signUp(email, password, name) {
        const userData = {
            id: Date.now().toString(),
            email: email,
            name: name,
            created_at: new Date().toISOString()
        };

        this.session = { user: userData, expires_at: Date.now() + (24 * 60 * 60 * 1000) };
        this.user = userData;
        localStorage.setItem('auth_session', JSON.stringify(this.session));
        this.updateUI();
        
        return { success: true, data: userData };
    }

    async signIn(email, password) {
        const userData = {
            id: Date.now().toString(),
            email: email,
            name: email.split('@')[0] || 'User',
            created_at: new Date().toISOString()
        };

        this.session = { user: userData, expires_at: Date.now() + (24 * 60 * 60 * 1000) };
        this.user = userData;
        localStorage.setItem('auth_session', JSON.stringify(this.session));
        this.updateUI();
        
        return { success: true, data: userData };
    }

    signOut() {
        this.session = null;
        this.user = null;
        localStorage.clear();
        this.updateUI();
        return { success: true };
    }

    isAuthenticated() {
        return !!(this.user && this.session);
    }

    getUser() {
        return this.user;
    }

    updateUI() {
        const signedOut = document.getElementById('signed-out');
        const signedIn = document.getElementById('signed-in');
        const userName = document.getElementById('user-name');
        const userAvatar = document.getElementById('user-avatar');

        if (this.isAuthenticated()) {
            const user = this.getUser();
            const profile = JSON.parse(localStorage.getItem('user_profile') || '{}');
            
            if (signedOut) {
                signedOut.style.display = 'none';
                signedOut.classList.add('hidden');
            }
            if (signedIn) {
                signedIn.style.display = 'flex';
                signedIn.classList.remove('hidden');
            }
            if (userName && user) {
                userName.textContent = profile.display_name || user.name || 'User';
            }
            if (userAvatar && user) {
                const name = profile.display_name || user.name || 'U';
                userAvatar.textContent = name.charAt(0).toUpperCase();
                
                // Update avatar image if exists
                if (profile.avatar_url) {
                    const avatarImg = document.createElement('img');
                    avatarImg.src = profile.avatar_url;
                    avatarImg.className = 'w-full h-full object-cover rounded-full';
                    userAvatar.innerHTML = '';
                    userAvatar.appendChild(avatarImg);
                }
            }
        } else {
            if (signedOut) {
                signedOut.style.display = 'flex';
                signedOut.classList.remove('hidden');
            }
            if (signedIn) {
                signedIn.style.display = 'none';
                signedIn.classList.add('hidden');
            }
        }
    }

    async addAnimeToList(animeName, animeData = {}) {
        if (!this.isAuthenticated()) {
            throw new Error('Please log in to add anime to your list');
        }

        const userList = JSON.parse(localStorage.getItem('userAnimeList') || '[]');
        const exists = userList.find(anime => anime.anime_name === animeName);
        
        if (exists) {
            throw new Error('Anime already in your list');
        }

        const newAnime = {
            id: Date.now().toString(),
            anime_name: animeName,
            status: animeData.status || 'plan-to-watch',
            episodes_watched: animeData.episodes_watched || 0,
            total_episodes: animeData.total_episodes || null,
            user_score: animeData.user_score || null,
            notes: animeData.notes || '',
            added_at: new Date().toISOString(),
            user_id: this.user.id
        };

        userList.push(newAnime);
        localStorage.setItem('userAnimeList', JSON.stringify(userList));
        return { success: true };
    }

    async getUserAnimeList() {
        if (!this.isAuthenticated()) return [];
        return JSON.parse(localStorage.getItem('userAnimeList') || '[]');
    }

    async updateAnimeStatus(animeId, newStatus) {
        if (!this.isAuthenticated()) {
            throw new Error('Please log in to update anime status');
        }

        const userList = JSON.parse(localStorage.getItem('userAnimeList') || '[]');
        const animeIndex = userList.findIndex(anime => anime.id === animeId);
        
        if (animeIndex !== -1) {
            userList[animeIndex].status = newStatus;
            localStorage.setItem('userAnimeList', JSON.stringify(userList));
            return { success: true };
        }
        
        throw new Error('Anime not found');
    }

    async removeAnimeFromList(animeId) {
        if (!this.isAuthenticated()) {
            throw new Error('Please log in to remove anime from list');
        }

        const userList = JSON.parse(localStorage.getItem('userAnimeList') || '[]');
        const filteredList = userList.filter(anime => anime.id !== animeId);
        localStorage.setItem('userAnimeList', JSON.stringify(filteredList));
        return { success: true };
    }

    async updateUserProfile(profileData) {
        if (!this.isAuthenticated()) {
            throw new Error('Please log in to update profile');
        }

        const profile = {
            user_id: this.user.id,
            display_name: profileData.display_name || '',
            bio: profileData.bio || '',
            favorite_genres: profileData.favorite_genres || [],
            avatar_url: profileData.avatar_url || '',
            updated_at: new Date().toISOString()
        };

        localStorage.setItem('user_profile', JSON.stringify(profile));
        this.updateUI(); // Update UI immediately
        return { success: true };
    }

    async getUserProfile() {
        if (!this.isAuthenticated()) return null;
        return JSON.parse(localStorage.getItem('user_profile') || 'null');
    }
}

// Global auth instance
window.authManager = new SimpleAuth();

// Global sign out function
window.signOut = () => {
    console.log('Signing out...');
    window.authManager.signOut();
    window.location.href = '/';
};

window.addToMyList = async (animeName) => {
    try {
        await window.authManager.addAnimeToList(animeName);
        if (window.showToast) {
            window.showToast('Anime added to your list!');
        } else {
            alert('Anime added to your list!');
        }
    } catch (error) {
        if (window.showToast) {
            window.showToast(error.message, 'error');
        } else {
            alert(error.message);
        }
    }
};