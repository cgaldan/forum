// Configuration
const API_URL = '/api';

// State management
const state = {
    currentUser: null,
    token: null,
    currentView: 'auth'
};

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    initializeApp();
});

function initializeApp() {
    // Check for existing session
    const savedToken = localStorage.getItem('token');
    if (savedToken) {
        state.token = savedToken;
        verifySession();
    } else {
        showAuthView();
    }

    // Set up event listeners
    setupAuthListeners();
    setupMainViewListeners();
}

// Authentication Listeners
function setupAuthListeners() {
    const loginTab = document.getElementById('login-tab');
    const registerTab = document.getElementById('register-tab');
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');

    loginTab.addEventListener('click', () => switchAuthTab('login'));
    registerTab.addEventListener('click', () => switchAuthTab('register'));
    loginForm.addEventListener('submit', handleLogin);
    registerForm.addEventListener('submit', handleRegister);
}

function setupMainViewListeners() {
    const logoutBtn = document.getElementById('logout-btn');
    logoutBtn.addEventListener('click', handleLogout);
}

// Tab switching
function switchAuthTab(tab) {
    const loginTab = document.getElementById('login-tab');
    const registerTab = document.getElementById('register-tab');
    const loginFormDiv = document.getElementById('login-form');
    const registerFormDiv = document.getElementById('register-form');

    if (tab === 'login') {
        loginTab.classList.add('active');
        registerTab.classList.remove('active');
        loginFormDiv.classList.add('active');
        registerFormDiv.classList.remove('active');
        clearError('login-error');
    } else {
        registerTab.classList.add('active');
        loginTab.classList.remove('active');
        registerFormDiv.classList.add('active');
        loginFormDiv.classList.remove('active');
        clearError('register-error');
    }
}

// Handle Login
async function handleLogin(e) {
    e.preventDefault();
    clearError('login-error');

    const identifier = document.getElementById('login-identifier').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ identifier, password })
        });

        const data = await response.json();

        if (data.success) {
            state.token = data.token;
            state.currentUser = data.user;
            localStorage.setItem('token', data.token);
            showMainView();
        } else {
            showError('login-error', data.message || 'Login failed');
        }
    } catch (error) {
        showError('login-error', 'Network error. Please check if the server is running.');
        console.error('Login error:', error);
    }
}

// Handle Register
async function handleRegister(e) {
    e.preventDefault();
    clearError('register-error');

    const formData = {
        nickname: document.getElementById('register-nickname').value,
        email: document.getElementById('register-email').value,
        password: document.getElementById('register-password').value,
        first_name: document.getElementById('register-firstname').value,
        last_name: document.getElementById('register-lastname').value,
        age: parseInt(document.getElementById('register-age').value),
        gender: document.getElementById('register-gender').value
    };

    try {
        const response = await fetch(`${API_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        const data = await response.json();

        if (data.success) {
            state.token = data.token;
            state.currentUser = data.user;
            localStorage.setItem('token', data.token);
            showMainView();
        } else {
            showError('register-error', data.message || 'Registration failed');
        }
    } catch (error) {
        showError('register-error', 'Network error. Please check if the server is running.');
        console.error('Register error:', error);
    }
}

// Handle Logout
async function handleLogout() {
    try {
        await fetch(`${API_URL}/auth/logout`, {
            method: 'POST',
            headers: {
                'Authorization': state.token
            }
        });
    } catch (error) {
        console.error('Logout error:', error);
    }

    // Clear state regardless of API response
    state.token = null;
    state.currentUser = null;
    localStorage.removeItem('token');
    showAuthView();
}

// Verify existing session
async function verifySession() {
    try {
        const response = await fetch(`${API_URL}/auth/me`, {
            method: 'GET',
            headers: {
                'Authorization': state.token
            }
        });

        const data = await response.json();

        if (data.success && data.user) {
            state.currentUser = data.user;
            showMainView();
        } else {
            // Invalid session
            localStorage.removeItem('token');
            state.token = null;
            showAuthView();
        }
    } catch (error) {
        console.error('Session verification error:', error);
        localStorage.removeItem('token');
        state.token = null;
        showAuthView();
    }
}

// View management
function showAuthView() {
    document.getElementById('auth-view').classList.add('active');
    document.getElementById('auth-view').classList.remove('hidden');
    document.getElementById('main-view').classList.remove('active');
    document.getElementById('main-view').classList.add('hidden');
    state.currentView = 'auth';
}

function showMainView() {
    document.getElementById('auth-view').classList.remove('active');
    document.getElementById('auth-view').classList.add('hidden');
    document.getElementById('main-view').classList.add('active');
    document.getElementById('main-view').classList.remove('hidden');
    state.currentView = 'main';

    // Update user info in the UI
    if (state.currentUser) {
        document.getElementById('user-nickname').textContent = state.currentUser.nickname;
        document.getElementById('profile-nickname').textContent = state.currentUser.nickname;
        document.getElementById('profile-email').textContent = state.currentUser.email;
        document.getElementById('profile-name').textContent = 
            `${state.currentUser.first_name} ${state.currentUser.last_name}`;
    }
}

// Error handling
function showError(elementId, message) {
    const errorElement = document.getElementById(elementId);
    errorElement.textContent = message;
    errorElement.classList.add('show');
}

function clearError(elementId) {
    const errorElement = document.getElementById(elementId);
    errorElement.textContent = '';
    errorElement.classList.remove('show');
}

