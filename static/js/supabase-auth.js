// Supabase Authentication with new credentials
const SUPABASE_URL = 'https://nnjjmlvhbejtrkmrajoz.supabase.co';
const SUPABASE_ANON_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im5uamptbHZoYmVqdHJrbXJham96Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTE2MTMzMDMsImV4cCI6MjA2NzE4OTMwM30.ACROHWi6_1_w2a3vGLp4JR50D3IvPVzWzrSnv-lXCvE';

class SupabaseAuth {
    constructor() {
        this.user = null;
        this.session = null;
        this.init();
    }

    async init() {
        const session = localStorage.getItem('sb_session');
        if (session) {
            try {
                const sessionData = JSON.parse(session);
                if (sessionData.expires_at > Date.now()) {
                    this.session = sessionData;
                    this.user = sessionData.user;
                } else {
                    localStorage.removeItem('sb_session');
                }
            } catch (e) {
                localStorage.removeItem('sb_session');
            }
        }
        
        this.updateUI();
        // Update UI periodically to catch auth changes
        setInterval(() => this.updateUI(), 1000);
    }

    async signUp(email, password, name) {
        try {
            const response = await fetch(`${SUPABASE_URL}/auth/v1/signup`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'apikey': SUPABASE_ANON_KEY
                },
                body: JSON.stringify({
                    email: email,
                    password: password,
                    data: { full_name: name }
                })
            });

            const data = await response.json();
            
            if (response.ok && data.user) {
                this.session = {
                    ...data,
                    expires_at: Date.now() + (3600 * 1000)
                };
                this.user = data.user;
                localStorage.setItem('sb_session', JSON.stringify(this.session));
                this.updateUI();
                return { success: true, data };
            } else {
                return { success: false, error: data.msg || 'Signup failed' };
            }
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    async signIn(email, password) {
        try {
            const response = await fetch(`${SUPABASE_URL}/auth/v1/token?grant_type=password`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'apikey': SUPABASE_ANON_KEY
                },
                body: JSON.stringify({
                    email: email,
                    password: password
                })
            });

            const data = await response.json();
            
            if (response.ok && data.user) {
                this.session = {
                    ...data,
                    expires_at: Date.now() + (3600 * 1000)
                };
                this.user = data.user;
                localStorage.setItem('sb_session', JSON.stringify(this.session));
                
                // Force multiple UI updates to ensure navbar changes
                this.updateUI();
                setTimeout(() => this.updateUI(), 100);
                setTimeout(() => this.updateUI(), 500);
                
                return { success: true, data };
            } else {
                return { success: false, error: data.msg || 'Invalid credentials' };
            }
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    async signOut() {
        try {
            if (this.session?.access_token) {
                await fetch(`${SUPABASE_URL}/auth/v1/logout`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'apikey': SUPABASE_ANON_KEY,
                        'Authorization': `Bearer ${this.session.access_token}`
                    }
                });
            }
        } catch (error) {
            console.log('Logout API call failed, continuing with local cleanup');
        }
        
        this.session = null;
        this.user = null;
        localStorage.removeItem('sb_session');
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

        console.log('Updating UI, authenticated:', this.isAuthenticated());

        if (this.isAuthenticated()) {
            const user = this.getUser();
            console.log('User authenticated, updating UI for:', user?.email);
            
            if (signedOut) {
                signedOut.style.display = 'none';
                signedOut.classList.add('hidden');
            }
            if (signedIn) {
                signedIn.style.display = 'flex';
                signedIn.classList.remove('hidden');
            }
            if (userName && user) {
                // Try to get display name from profile first
                this.getUserProfile().then(profile => {
                    if (profile && profile.display_name) {
                        userName.textContent = profile.display_name;
                    } else {
                        userName.textContent = user.user_metadata?.full_name || user.email?.split('@')[0] || 'User';
                    }
                });
            }
            if (userAvatar && user) {
                // Check for profile avatar first
                this.getUserProfile().then(profile => {
                    if (profile?.avatar_url) {
                        userAvatar.innerHTML = `<img src="${profile.avatar_url}" class="w-full h-full object-cover rounded-full" alt="Avatar">`;
                    } else {
                        const name = user.user_metadata?.full_name || user.email || 'U';
                        userAvatar.textContent = name.charAt(0).toUpperCase();
                    }
                }).catch(() => {
                    const name = user.user_metadata?.full_name || user.email || 'U';
                    userAvatar.textContent = name.charAt(0).toUpperCase();
                });
            }
        } else {
            console.log('User not authenticated, showing login/signup');
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

        try {
            const response = await fetch(`${SUPABASE_URL}/rest/v1/user_anime_list`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'apikey': SUPABASE_ANON_KEY,
                    'Authorization': `Bearer ${this.session.access_token}`
                },
                body: JSON.stringify({
                    user_id: this.user.id,
                    anime_name: animeName,
                    status: animeData.status || 'plan-to-watch',
                    episodes_watched: animeData.episodes_watched || 0,
                    total_episodes: animeData.total_episodes || null,
                    user_score: animeData.user_score || null,
                    notes: animeData.notes || '',
                    added_at: new Date().toISOString()
                })
            });

            if (response.ok) {
                return { success: true };
            } else {
                const error = await response.json();
                throw new Error(error.message || 'Failed to add anime');
            }
        } catch (error) {
            throw new Error(error.message);
        }
    }

    async updateUserProfile(profileData) {
        if (!this.isAuthenticated()) {
            throw new Error('Please log in to update profile');
        }

        try {
            // Use UPSERT to handle both insert and update
            const response = await fetch(`${SUPABASE_URL}/rest/v1/user_profiles`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'apikey': SUPABASE_ANON_KEY,
                    'Authorization': `Bearer ${this.session.access_token}`,
                    'Prefer': 'resolution=merge-duplicates'
                },
                body: JSON.stringify({
                    user_id: this.user.id,
                    display_name: profileData.display_name || '',
                    bio: profileData.bio || '',
                    favorite_genres: profileData.favorite_genres || [],
                    avatar_url: profileData.avatar_url || ''
                })
            });

            if (response.ok || response.status === 409) {
                // If conflict (409), try update instead
                if (response.status === 409) {
                    const updateResponse = await fetch(`${SUPABASE_URL}/rest/v1/user_profiles?user_id=eq.${this.user.id}`, {
                        method: 'PATCH',
                        headers: {
                            'Content-Type': 'application/json',
                            'apikey': SUPABASE_ANON_KEY,
                            'Authorization': `Bearer ${this.session.access_token}`
                        },
                        body: JSON.stringify({
                            display_name: profileData.display_name || '',
                            bio: profileData.bio || '',
                            favorite_genres: profileData.favorite_genres || [],
                            avatar_url: profileData.avatar_url || ''
                        })
                    });
                    return { success: updateResponse.ok };
                }
                return { success: true };
            } else {
                const error = await response.json();
                throw new Error(error.message || 'Failed to update profile');
            }
        } catch (error) {
            throw new Error(error.message);
        }
    }

    async getUserProfile() {
        if (!this.isAuthenticated()) return null;

        try {
            const response = await fetch(`${SUPABASE_URL}/rest/v1/user_profiles?user_id=eq.${this.user.id}`, {
                headers: {
                    'apikey': SUPABASE_ANON_KEY,
                    'Authorization': `Bearer ${this.session.access_token}`
                }
            });

            if (response.ok) {
                const profiles = await response.json();
                return profiles[0] || null;
            }
            return null;
        } catch (error) {
            return null;
        }
    }

    async getUserAnimeList() {
        if (!this.isAuthenticated()) return [];

        try {
            const response = await fetch(`${SUPABASE_URL}/rest/v1/user_anime_list?user_id=eq.${this.user.id}&order=added_at.desc`, {
                headers: {
                    'apikey': SUPABASE_ANON_KEY,
                    'Authorization': `Bearer ${this.session.access_token}`
                }
            });

            if (response.ok) {
                return await response.json();
            } else {
                return [];
            }
        } catch (error) {
            return [];
        }
    }

    async updateAnimeStatus(animeId, newStatus) {
        if (!this.isAuthenticated()) {
            throw new Error('Please log in to update anime status');
        }

        try {
            const response = await fetch(`${SUPABASE_URL}/rest/v1/user_anime_list?id=eq.${animeId}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                    'apikey': SUPABASE_ANON_KEY,
                    'Authorization': `Bearer ${this.session.access_token}`
                },
                body: JSON.stringify({ status: newStatus })
            });

            if (response.ok) {
                return { success: true };
            } else {
                throw new Error('Failed to update status');
            }
        } catch (error) {
            throw new Error(error.message);
        }
    }

    async removeAnimeFromList(animeId) {
        if (!this.isAuthenticated()) {
            throw new Error('Please log in to remove anime from list');
        }

        try {
            const response = await fetch(`${SUPABASE_URL}/rest/v1/user_anime_list?id=eq.${animeId}`, {
                method: 'DELETE',
                headers: {
                    'apikey': SUPABASE_ANON_KEY,
                    'Authorization': `Bearer ${this.session.access_token}`
                }
            });

            if (response.ok) {
                return { success: true };
            } else {
                throw new Error('Failed to remove anime');
            }
        } catch (error) {
            throw new Error(error.message);
        }
    }
}

// Global auth instance
window.authManager = new SupabaseAuth();

// Global functions - WORKING SIGN OUT
window.signOut = () => {
    console.log('Signing out NOW...');
    
    // Clear all data immediately
    localStorage.clear();
    
    // Reset auth manager
    if (window.authManager) {
        window.authManager.session = null;
        window.authManager.user = null;
    }
    
    // Force redirect immediately
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