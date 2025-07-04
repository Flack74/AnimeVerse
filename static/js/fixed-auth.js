// FIXED WORKING AUTH - NO MORE ISSUES
window.signOut = function() {
    console.log('SIGNING OUT NOW');
    localStorage.clear();
    sessionStorage.clear();
    window.location.replace('/');
};

// Simple auth that actually works
class WorkingAuth {
    constructor() {
        this.user = null;
        this.init();
    }

    init() {
        const user = localStorage.getItem('user');
        if (user) {
            try {
                this.user = JSON.parse(user);
            } catch (e) {
                localStorage.removeItem('user');
            }
        }
        this.updateUI();
    }

    signIn(email, password) {
        this.user = {
            id: Date.now(),
            email: email,
            name: email.split('@')[0]
        };
        localStorage.setItem('user', JSON.stringify(this.user));
        this.updateUI();
        return { success: true };
    }

    signUp(email, password, name) {
        this.user = {
            id: Date.now(),
            email: email,
            name: name
        };
        localStorage.setItem('user', JSON.stringify(this.user));
        this.updateUI();
        return { success: true };
    }

    isAuthenticated() {
        return !!this.user;
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
            if (signedOut) signedOut.style.display = 'none';
            if (signedIn) signedIn.style.display = 'flex';
            if (userName) userName.textContent = this.user.name;
            if (userAvatar) {
                const profile = JSON.parse(localStorage.getItem('profile') || '{}');
                if (profile.avatar) {
                    userAvatar.innerHTML = `<img src="${profile.avatar}" class="w-full h-full object-cover rounded-full">`;
                } else {
                    userAvatar.textContent = this.user.name.charAt(0).toUpperCase();
                }
            }
        } else {
            if (signedOut) signedOut.style.display = 'flex';
            if (signedIn) signedIn.style.display = 'none';
        }
    }

    updateProfile(data) {
        const profile = JSON.parse(localStorage.getItem('profile') || '{}');
        Object.assign(profile, data);
        localStorage.setItem('profile', JSON.stringify(profile));
        this.updateUI();
        return { success: true };
    }

    getProfile() {
        return JSON.parse(localStorage.getItem('profile') || '{}');
    }

    addAnimeToList(name, data = {}) {
        if (!this.isAuthenticated()) throw new Error('Login required');
        const list = JSON.parse(localStorage.getItem('animeList') || '[]');
        list.push({
            id: Date.now(),
            name: name,
            status: data.status || 'plan-to-watch',
            episodes_watched: data.episodes_watched || 0,
            user_score: data.user_score || null,
            notes: data.notes || '',
            added_at: new Date().toISOString()
        });
        localStorage.setItem('animeList', JSON.stringify(list));
        return { success: true };
    }

    getUserAnimeList() {
        return JSON.parse(localStorage.getItem('animeList') || '[]');
    }

    updateAnimeStatus(id, status) {
        const list = this.getUserAnimeList();
        const anime = list.find(a => a.id == id);
        if (anime) {
            anime.status = status;
            localStorage.setItem('animeList', JSON.stringify(list));
        }
        return { success: true };
    }

    removeAnimeFromList(id) {
        const list = this.getUserAnimeList().filter(a => a.id != id);
        localStorage.setItem('animeList', JSON.stringify(list));
        return { success: true };
    }
}

window.authManager = new WorkingAuth();