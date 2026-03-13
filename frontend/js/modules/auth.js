import api from '../api/client.js';
import store from '../state/store.js';
import { showError, clearError, showToast } from '../ui/ui.js';
import { getElement, setText } from '../utils/helpers.js';

export async function initAuth() {
    const savedToken = localStorage.getItem('token');
    if (savedToken) {
        store.set('token', savedToken);
        api.setToken(savedToken);
        await verifySession();
    }
}

export async function handleLogin(e) {
    e.preventDefault();
    clearError('login-error');

    const identifier = getElement('login-identifier').value;
    const password = getElement('login-password').value;

    const submitBtn = e.target.querySelector('button[type="submit"]');
    if (submitBtn) { submitBtn.disabled = true; submitBtn.textContent = 'Logging in...'; }

    try {
        const data = await api.login(identifier, password);

        if (data.success) {
            store.set('token', data.token);
            store.set('currentUser', data.user);
            api.setToken(data.token);
            localStorage.setItem('token', data.token);
            showMainView();
        } else {
            showError('login-error', data.message || 'Login failed');
        }
    } catch (error) {
        showError('login-error', 'Network error. Please check if the server is running.');
        console.error('Login error:', error);
    } finally {
        if (submitBtn) { submitBtn.disabled = false; submitBtn.textContent = 'Login'; }
    }
}

export async function handleRegister(e) {
    e.preventDefault();
    clearError('register-error');

    const formData = {
        nickname: getElement('register-nickname').value,
        email: getElement('register-email').value,
        password: getElement('register-password').value,
        first_name: getElement('register-firstname').value,
        last_name: getElement('register-lastname').value,
        age: parseInt(getElement('register-age').value),
        gender: getElement('register-gender').value
    };

    const submitBtn = e.target.querySelector('button[type="submit"]');
    if (submitBtn) { submitBtn.disabled = true; submitBtn.textContent = 'Registering...'; }

    try {
        const data = await api.register(formData);

        if (data.success) {
            store.set('token', data.token);
            store.set('currentUser', data.user);
            api.setToken(data.token);
            localStorage.setItem('token', data.token);
            showMainView();
        } else {
            showError('register-error', data.message || 'Registration failed');
        }
    } catch (error) {
        showError('register-error', 'Network error. Please check if the server is running.');
        console.error('Register error:', error);
    } finally {
        if (submitBtn) { submitBtn.disabled = false; submitBtn.textContent = 'Register'; }
    }
}


export async function handleLogout() {
    const { disconnectWebSocket } = await import('./websocket.js');
    const { closeChatPanel } = await import('./messages.js');

    disconnectWebSocket();

    closeChatPanel();

    try {
        await api.logout();
    } catch (error) {
        console.error('Logout error:', error);
    }

    clearAuthState();
    showAuthView();
}

export async function verifySession() {
    console.log('Verifying session...');
    try {
        const data = await api.verifySession();

        if (data.success && data.user) {
            store.set('currentUser', data.user);
            console.log('Session valid, showing main view');
            showMainView();
        } else {
            clearAuthState();
            showAuthView();
        }
    } catch (error) {
        console.error('Session verification error:', error);
        clearAuthState();
        showAuthView();
    }
}

export function clearAuthState() {
    store.set('token', null);
    store.set('currentUser', null);
    api.setToken(null);
    localStorage.removeItem('token');
}

function showAuthView() {
    const authView = getElement('auth-view');
    const mainView = getElement('main-view');
    
    authView.classList.add('active');
    authView.classList.remove('hidden');
    mainView.classList.remove('active');
    mainView.classList.add('hidden');
    
    store.set('currentView', 'auth');

    const messagePanel = getElement('message-panel');
    if (messagePanel) messagePanel.classList.add('hidden');

    const conversationsContainer = getElement('conversations-container');
    if (conversationsContainer) {
        conversationsContainer.innerHTML = '<div class="no-conversations">No conversations yet</div>';
    }
    
    const onlineUsersContainer = getElement('online-users-container');
    if (onlineUsersContainer) {
        onlineUsersContainer.innerHTML = '<div class="loading">Loading...</div>';
    }
    
    const postsContainer = getElement('posts-container');
    if (postsContainer) {
        postsContainer.innerHTML = '<div class="loading">Loading posts...</div>';
    }
}

async function showMainView() {
    const authView = getElement('auth-view');
    const mainView = getElement('main-view');
    
    authView.classList.remove('active');
    authView.classList.add('hidden');
    mainView.classList.add('active');
    mainView.classList.remove('hidden');
    
    store.set('currentView', 'main');

    const currentUser = store.get('currentUser');
    if (currentUser) {
        setText('user-nickname', currentUser.nickname);
    }

    const { loadPosts } = await import('./posts.js');
    const { loadConversations } = await import('./messages.js');
    const { connectWebSocket } = await import('./websocket.js');

    loadPosts();
    loadConversations();
    connectWebSocket();
}
