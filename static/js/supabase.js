// FINAL WORKING SUPABASE AUTH
const SUPABASE_URL = 'https://nnjjmlvhbejtrkmrajoz.supabase.co';
const SUPABASE_ANON_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im5uamptbHZoYmVqdHJrbXJham96Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTE2MTMzMDMsImV4cCI6MjA2NzE4OTMwM30.ACROHWi6_1_w2a3vGLp4JR50D3IvPVzWzrSnv-lXCvE';

class SupabaseAuth {
    constructor() {
        this.user = null;
        this.session = null;
        this.initialized = false;
        console.log('🚀 SupabaseAuth constructor called');
        this.init();
    }

    async init() {
        console.log('🔄 Initializing auth...');
        
        // Small delay to ensure DOM is ready
        await new Promise(resolve => setTimeout(resolve, 200));
        
        const session = localStorage.getItem('sb_session');
        if (session) {
            try {
                const sessionData = JSON.parse(session);
                if (sessionData.expires_at && sessionData.expires_at > Date.now()) {
                    this.session = sessionData;
                    this.user = sessionData.user;
                    console.log('✅ Session restored for:', this.user?.email);
                } else {
                    console.log('⏰ Session expired');
                    localStorage.removeItem('sb_session');
                }
            } catch (e) {
                console.log('❌ Invalid session data');
                localStorage.removeItem('sb_session');
            }
        }
        
        this.initialized = true;
        this.updateUI();
        
        // Force update after delays
        setTimeout(() => this.updateUI(), 500);
        setTimeout(() => this.updateUI(), 1000);
    }

    async signUp(email, password, name) {
        console.log('📝 Signing up:', email);
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
            console.log('📝 Signup response:', data);
            
            if (response.ok && data.user) {
                this.session = {
                    ...data,
                    expires_at: Date.now() + (3600 * 1000)
                };
                this.user = data.user;
                localStorage.setItem('sb_session', JSON.stringify(this.session));
                this.forceUIUpdate();
                return { success: true, data };
            } else {
                return { success: false, error: data.msg || data.error_description || 'Signup failed' };
            }
        } catch (error) {
            console.error('❌ Signup error:', error);
            return { success: false, error: error.message };
        }
    }

    async signIn(email, password) {
        console.log('🔐 Attempting sign in for:', email);
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
            console.log('🔐 Login response status:', response.status);
            console.log('🔐 Login response data:', data);
            
            if (response.ok && data.user) {
                this.session = {
                    ...data,
                    expires_at: Date.now() + (3600 * 1000)
                };
                this.user = data.user;
                localStorage.setItem('sb_session', JSON.stringify(this.session));
                
                console.log('✅ Login successful! User:', this.user.email);
                this.forceUIUpdate();
                
                return { success: true, data };
            } else {
                console.log('❌ Login failed:', data);
                return { 
                    success: false, 
                    error: data.error_description || data.msg || 'Invalid email or password' 
                };
            }
        } catch (error) {
            console.error('❌ Login error:', error);
            return { success: false, error: 'Login failed. Please check your connection.' };
        }
    }

    signOut() {
        console.log('🚪 Signing out...');
        this.session = null;
        this.user = null;
        localStorage.removeItem('sb_session');
        this.forceUIUpdate();
        console.log('✅ Signed out successfully');
        return { success: true };
    }

    isAuthenticated() {
        const result = !!(this.user && this.session);
        return result;
    }

    getUser() {
        return this.user;
    }

    forceUIUpdate() {
        console.log('🎨 Force updating UI...');
        this.updateUI();
        setTimeout(() => this.updateUI(), 100);
        setTimeout(() => this.updateUI(), 300);
        setTimeout(() => this.updateUI(), 500);
        setTimeout(() => this.updateUI(), 1000);
    }

    updateUI() {
        if (!this.initialized) return;
        
        const isAuth = this.isAuthenticated();
        console.log('🎨 Updating UI - Authenticated:', isAuth);
        
        const signedOut = document.getElementById('signed-out');
        const signedIn = document.getElementById('signed-in');
        const userName = document.getElementById('user-name');
        const userAvatar = document.getElementById('user-avatar');

        if (isAuth) {
            const user = this.getUser();
            console.log('👤 Showing user info for:', user.email);
            
            // Hide login buttons
            if (signedOut) {
                signedOut.style.display = 'none';
                signedOut.classList.add('hidden');
            }
            
            // Show user info
            if (signedIn) {
                signedIn.style.display = 'flex';
                signedIn.classList.remove('hidden');
            }
            
            // Update user name
            if (userName) {
                const displayName = user.user_metadata?.full_name || user.email?.split('@')[0] || 'User';
                userName.textContent = displayName;
            }
            
            // Update avatar
            if (userAvatar) {
                const name = user.user_metadata?.full_name || user.email || 'U';
                userAvatar.textContent = name.charAt(0).toUpperCase();
            }
        } else {
            console.log('🔓 Showing login buttons');
            
            // Show login buttons
            if (signedOut) {
                signedOut.style.display = 'flex';
                signedOut.classList.remove('hidden');
            }
            
            // Hide user info
            if (signedIn) {
                signedIn.style.display = 'none';
                signedIn.classList.add('hidden');
            }
        }
    }

    async addAnimeToList(animeName, animeData = {}) {
        if (!this.isAuthenticated()) {
            if (confirm(`Please log in to add "${animeName}" to your list. Go to login page?`)) {
                window.location.href = '/static/login.html';
            }
            return { success: false, error: 'Please log in first' };
        }

        try {
            const response = await fetch(`${SUPABASE_URL}/rest/v1/user_anime_list`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'apikey': SUPABASE_ANON_KEY,
                    'Authorization': `Bearer ${this.session.access_token}`,
                    'Prefer': 'return=minimal'
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

            if (response.ok || response.status === 201) {
                return { success: true };
            } else {
                const error = await response.json();
                throw new Error(error.message || 'Failed to add anime');
            }
        } catch (error) {
            throw new Error(error.message);
        }
    }

    async getUserAnimeList() {
        if (!this.isAuthenticated()) {
            return [];
        }

        try {
            const response = await fetch(`${SUPABASE_URL}/rest/v1/user_anime_list?user_id=eq.${this.user.id}&select=*`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'apikey': SUPABASE_ANON_KEY,
                    'Authorization': `Bearer ${this.session.access_token}`
                }
            });

            if (response.ok) {
                const data = await response.json();
                return data || [];
            } else {
                console.error('Failed to fetch user anime list');
                return [];
            }
        } catch (error) {
            console.error('Error fetching user anime list:', error);
            return [];
        }
    }

    async updateAnimeStatus(animeId, newStatus) {
        if (!this.isAuthenticated()) {
            throw new Error('Please log in first');
        }

        try {
            const response = await fetch(`${SUPABASE_URL}/rest/v1/user_anime_list?id=eq.${animeId}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                    'apikey': SUPABASE_ANON_KEY,
                    'Authorization': `Bearer ${this.session.access_token}`
                },
                body: JSON.stringify({
                    status: newStatus,
                    updated_at: new Date().toISOString()
                })
            });

            if (!response.ok) {
                throw new Error('Failed to update anime status');
            }
        } catch (error) {
            throw new Error(error.message);
        }
    }

    async removeAnimeFromList(animeId) {
        if (!this.isAuthenticated()) {
            throw new Error('Please log in first');
        }

        try {
            const response = await fetch(`${SUPABASE_URL}/rest/v1/user_anime_list?id=eq.${animeId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    'apikey': SUPABASE_ANON_KEY,
                    'Authorization': `Bearer ${this.session.access_token}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to remove anime from list');
            }
        } catch (error) {
            throw new Error(error.message);
        }
    }

    async updateUserProfile(profileData) {
        if (!this.isAuthenticated()) {
            throw new Error('Please log in first');
        }

        try {
            const response = await fetch(`${SUPABASE_URL}/rest/v1/user_profiles?user_id=eq.${this.user.id}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                    'apikey': SUPABASE_ANON_KEY,
                    'Authorization': `Bearer ${this.session.access_token}`
                },
                body: JSON.stringify({
                    ...profileData,
                    updated_at: new Date().toISOString()
                })
            });

            if (!response.ok) {
                // Try to create profile if it doesn't exist
                const createResponse = await fetch(`${SUPABASE_URL}/rest/v1/user_profiles`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'apikey': SUPABASE_ANON_KEY,
                        'Authorization': `Bearer ${this.session.access_token}`
                    },
                    body: JSON.stringify({
                        user_id: this.user.id,
                        ...profileData,
                        created_at: new Date().toISOString(),
                        updated_at: new Date().toISOString()
                    })
                });
                
                if (!createResponse.ok) {
                    throw new Error('Failed to update profile');
                }
            }
        } catch (error) {
            throw new Error(error.message);
        }
    }
}

// Create global auth instance
console.log('🚀 Creating global auth manager...');
window.authManager = new SupabaseAuth();

// Global sign out function
window.signOut = function() {
    console.log('🚪 Global signOut called');
    if (window.authManager) {
        window.authManager.signOut();
        // Small delay then redirect
        setTimeout(() => {
            window.location.href = '/';
        }, 500);
    }
};

// Add to list function
window.addToMyList = async function(animeName) {
    console.log('📝 Adding to list:', animeName);
    
    if (!window.authManager || !window.authManager.isAuthenticated()) {
        if (confirm(`Please log in to add "${animeName}" to your list. Go to login page?`)) {
            window.location.href = '/static/login.html';
        }
        return;
    }
    
    try {
        const result = await window.authManager.addAnimeToList(animeName);
        if (result.success) {
            showToast(`✅ "${animeName}" added to your list!`);
        } else {
            showToast(`❌ Failed to add anime: ${result.error}`, 'error');
        }
    } catch (error) {
        showToast(`❌ Error: ${error.message}`, 'error');
    }
};

// Toast notification function
window.showToast = function(message, type = 'success') {
    const toast = document.createElement('div');
    toast.className = `fixed top-4 right-4 z-50 px-6 py-3 rounded-lg shadow-lg text-white font-medium transition-all transform translate-x-full ${
        type === 'error' ? 'bg-red-500' : 'bg-green-500'
    }`;
    toast.textContent = message;
    
    document.body.appendChild(toast);
    
    // Animate in
    setTimeout(() => {
        toast.classList.remove('translate-x-full');
    }, 100);
    
    // Animate out and remove
    setTimeout(() => {
        toast.classList.add('translate-x-full');
        setTimeout(() => {
            if (toast.parentNode) {
                toast.parentNode.removeChild(toast);
            }
        }, 300);
    }, 3000);
};

console.log('✅ Supabase auth script loaded');

// Create skeleton loading function
window.createSkeleton = function(count) {
    return Array(count).fill(0).map(() => 
        '<div class="bg-gray-200 dark:bg-gray-700 rounded-xl h-80 animate-pulse"></div>'
    ).join('');
};