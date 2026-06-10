const API_BASE = 'http://localhost:8080/api';

let authToken = localStorage.getItem('auth_token');
let currentUser = null;

async function apiRequest(endpoint, options = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };
    
    if (authToken) {
        headers['Authorization'] = `Bearer ${authToken}`;
    }
    
    const response = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers
    });
    
    if (response.status === 401) {
        logout();
        throw new Error('Unauthorized');
    }
    
    return response;
}

async function register(username, email, password) {
    const response = await apiRequest('/register', {
        method: 'POST',
        body: JSON.stringify({ username, email, password })
    });
    
    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Registration failed');
    }
    
    const data = await response.json();
    authToken = data.token;
    currentUser = data.user;
    localStorage.setItem('auth_token', authToken);
    localStorage.setItem('current_user', JSON.stringify(currentUser));
    return data;
}

async function login(username, password) {
    const response = await apiRequest('/login', {
        method: 'POST',
        body: JSON.stringify({ username, password })
    });
    
    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Login failed');
    }
    
    const data = await response.json();
    authToken = data.token;
    currentUser = data.user;
    localStorage.setItem('auth_token', authToken);
    localStorage.setItem('current_user', JSON.stringify(currentUser));
    return data;
}

function logout() {
    authToken = null;
    currentUser = null;
    localStorage.removeItem('auth_token');
    localStorage.removeItem('current_user');
    window.location.reload();
}

function checkAuth() {
    const token = localStorage.getItem('auth_token');
    if (token) {
        authToken = token;
        const user = localStorage.getItem('current_user');
        if (user) {
            currentUser = JSON.parse(user);
        }
        return true;
    }
    return false;
}

async function searchUsers(query) {
    const response = await apiRequest(`/users/search?q=${encodeURIComponent(query)}`);
    if (!response.ok) throw new Error('Search failed');
    return response.json();
}