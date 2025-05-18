// Authentication related functionality

// Constants
const TOKEN_KEY = 'lumo_token';
const REFRESH_TOKEN_KEY = 'lumo_refresh_token';
const USERNAME_KEY = 'lumo_username';
const TOKEN_EXPIRY_KEY = 'lumo_token_expiry';

// Check if user is authenticated
function isAuthenticated() {
    const token = localStorage.getItem(TOKEN_KEY);
    const expiry = localStorage.getItem(TOKEN_EXPIRY_KEY);
    
    if (!token || !expiry) {
        return false;
    }
    
    // Check if token is expired
    const expiryDate = new Date(parseInt(expiry));
    const now = new Date();
    
    // If token expires in less than 5 minutes, try to refresh it
    if ((expiryDate - now) < 5 * 60 * 1000) {
        refreshToken();
    }
    
    return expiryDate > now;
}

// Login function
async function login(username, password) {
    try {
        const response = await fetch('/api/v1/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password })
        });
        
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
            throw new Error(errorData.error || response.statusText);
        }
        
        const data = await response.json();
        
        // Save auth data to localStorage
        localStorage.setItem(TOKEN_KEY, data.token);
        localStorage.setItem(REFRESH_TOKEN_KEY, data.refresh_token);
        localStorage.setItem(USERNAME_KEY, data.username);
        
        // Calculate and store expiry time
        const expiryTime = new Date().getTime() + (data.expires_in * 1000);
        localStorage.setItem(TOKEN_EXPIRY_KEY, expiryTime.toString());
        
        return data;
    } catch (error) {
        console.error('Login error:', error);
        throw error;
    }
}

// Logout function
function logout() {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem(USERNAME_KEY);
    localStorage.removeItem(TOKEN_EXPIRY_KEY);
    
    // Redirect to login page
    showLoginPage();
}

// Refresh token function
async function refreshToken() {
    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
    
    if (!refreshToken) {
        logout();
        return;
    }
    
    try {
        const response = await fetch('/api/v1/auth/refresh', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ refresh_token: refreshToken })
        });
        
        if (!response.ok) {
            // If refresh fails, logout
            logout();
            return;
        }
        
        const data = await response.json();
        
        // Update auth data in localStorage
        localStorage.setItem(TOKEN_KEY, data.token);
        localStorage.setItem(REFRESH_TOKEN_KEY, data.refresh_token);
        localStorage.setItem(USERNAME_KEY, data.username);
        
        // Calculate and store expiry time
        const expiryTime = new Date().getTime() + (data.expires_in * 1000);
        localStorage.setItem(TOKEN_EXPIRY_KEY, expiryTime.toString());
        
        return data;
    } catch (error) {
        console.error('Token refresh error:', error);
        logout();
    }
}

// Change password function
async function changePassword(currentPassword, newPassword) {
    try {
        const token = localStorage.getItem(TOKEN_KEY);
        
        if (!token) {
            throw new Error('Not authenticated');
        }
        
        const response = await fetch('/api/v1/auth/change-password', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ 
                current_password: currentPassword,
                new_password: newPassword
            })
        });
        
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
            throw new Error(errorData.error || response.statusText);
        }
        
        return await response.json();
    } catch (error) {
        console.error('Change password error:', error);
        throw error;
    }
}

// Get auth token for API requests
function getAuthToken() {
    return localStorage.getItem(TOKEN_KEY);
}

// Get current username
function getUsername() {
    return localStorage.getItem(USERNAME_KEY);
}

// Show login page
function showLoginPage() {
    document.getElementById('login-page').classList.remove('hidden');
    document.getElementById('app-page').classList.add('hidden');
}

// Show app page
function showAppPage() {
    document.getElementById('login-page').classList.add('hidden');
    document.getElementById('app-page').classList.remove('hidden');
    
    // Update username display
    const username = getUsername();
    document.getElementById('username-display').textContent = username;
}
