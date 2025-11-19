// Configuration
const API_URL = '/api';

// State management
const state = {
    currentUser: null,
    token: null,
    currentView: 'auth',
    posts: [],
    currentPost: null,
    currentCategory: '',
    ws: null,
    wsReconnectAttempts: 0,
    wsMaxReconnectAttempts: 5,
    onlineUsers: [],
    currentChatUser: null,
    messages: [],
    conversations: []
};

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    initializeApp();
    setupKeyboardNavigation();
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

    const createPostForm = document.getElementById('createPostForm');
    createPostForm.addEventListener('submit', handleCreatePost);

    const categoryFilter = document.getElementById('category-filter');
    categoryFilter.addEventListener('change', handleCategoryFilter);

    const createCommentForm = document.getElementById('createCommentForm');
    createCommentForm.addEventListener('submit', handleCreateComment);

    const closeModal = document.querySelector('.close-modal');
    closeModal.addEventListener('click', closePostModal);

    const modal = document.getElementById('post-modal');
    modal.addEventListener('click', (e) => {
        if (e.target === modal) closePostModal();
    });

    const sendMessageForm = document.getElementById('sendMessageForm');
    sendMessageForm.addEventListener('submit', handleSendMessage);

    const closeChat = document.getElementById('close-chat');
    closeChat.addEventListener('click', closeChatPanel);
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
    // Disconnect WebSocket
    disconnectWebSocket();

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
    state.onlineUsers = [];
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
    }

    // Load posts
    loadPosts();

    // Load conversations
    loadConversations();

    // Connect WebSocket
    connectWebSocket();
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

// Toast Notifications
function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    toast.setAttribute('role', 'alert');
    toast.setAttribute('aria-live', 'polite');
    
    document.body.appendChild(toast);
    
    setTimeout(() => {
        toast.style.animation = 'slideIn 0.3s ease-out reverse';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

// Loading States
function showLoading(containerId) {
    const container = document.getElementById(containerId);
    if (container) {
        container.innerHTML = '<div class="spinner"></div>';
    }
}

function showSkeletonPosts() {
    const container = document.getElementById('posts-container');
    container.innerHTML = `
        <div class="skeleton skeleton-post"></div>
        <div class="skeleton skeleton-post"></div>
        <div class="skeleton skeleton-post"></div>
    `;
}

// Post Management
async function loadPosts(category = '') {
    showSkeletonPosts();
    
    try {
        let url = `${API_URL}/posts`;
        if (category) {
            url += `?category=${category}`;
        }

        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error('Failed to fetch posts');
        }
        
        const data = await response.json();

        if (data.success) {
            state.posts = data.posts || [];
            renderPosts();
        } else {
            throw new Error(data.message || 'Failed to load posts');
        }
    } catch (error) {
        console.error('Load posts error:', error);
        document.getElementById('posts-container').innerHTML = 
            '<div class="no-posts">Failed to load posts. Please try again.</div>';
        showToast('Failed to load posts', 'error');
    }
}

function renderPosts() {
    const container = document.getElementById('posts-container');
    
    if (state.posts.length === 0) {
        container.innerHTML = '<div class="no-posts">No posts yet. Be the first to post!</div>';
        return;
    }

    container.innerHTML = state.posts.map(post => `
        <div class="post-card" onclick="openPost(${post.id})">
            <div class="post-card-header">
                <h3 class="post-title">${escapeHtml(post.title)}</h3>
                <span class="post-category">${escapeHtml(post.category)}</span>
            </div>
            <div class="post-meta">
                by ${escapeHtml(post.author)} • ${formatDate(post.created_at)}
            </div>
            <div class="post-content-preview">
                ${escapeHtml(post.content.substring(0, 200))}${post.content.length > 200 ? '...' : ''}
            </div>
        </div>
    `).join('');
}

async function handleCreatePost(e) {
    e.preventDefault();
    clearError('post-error');

    const title = document.getElementById('post-title').value.trim();
    const category = document.getElementById('post-category').value;
    const content = document.getElementById('post-content').value.trim();

    // Client-side validation
    if (!title || title.length < 3) {
        showError('post-error', 'Title must be at least 3 characters');
        return;
    }
    if (!content || content.length < 10) {
        showError('post-error', 'Content must be at least 10 characters');
        return;
    }
    if (!category) {
        showError('post-error', 'Please select a category');
        return;
    }

    const submitBtn = e.target.querySelector('button[type="submit"]');
    submitBtn.disabled = true;
    submitBtn.textContent = 'Posting...';

    try {
        const response = await fetch(`${API_URL}/posts`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': state.token
            },
            body: JSON.stringify({ title, category, content })
        });

        const data = await response.json();

        if (data.success) {
            document.getElementById('createPostForm').reset();
            showToast('Post created successfully!', 'success');
            loadPosts(state.currentCategory);
        } else {
            showError('post-error', data.message || 'Failed to create post');
        }
    } catch (error) {
        showError('post-error', 'Network error. Please try again.');
        showToast('Failed to create post', 'error');
        console.error('Create post error:', error);
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Post';
    }
}

function handleCategoryFilter(e) {
    state.currentCategory = e.target.value;
    loadPosts(state.currentCategory);
}

async function openPost(postId) {
    try {
        const response = await fetch(`${API_URL}/posts/${postId}`);
        const data = await response.json();

        if (data.success) {
            state.currentPost = data.post;
            renderPostDetail(data.post, data.comments || []);
            document.getElementById('post-modal').classList.add('active');
            document.getElementById('post-modal').classList.remove('hidden');
        }
    } catch (error) {
        console.error('Open post error:', error);
    }
}

function renderPostDetail(post, comments) {
    document.getElementById('post-detail').innerHTML = `
        <h2 class="post-detail-title">${escapeHtml(post.title)}</h2>
        <div class="post-detail-meta">
            <span class="post-category">${escapeHtml(post.category)}</span> • 
            by ${escapeHtml(post.author)} • ${formatDate(post.created_at)}
        </div>
        <div class="post-detail-content">${escapeHtml(post.content)}</div>
    `;

    const commentsContainer = document.getElementById('comments-container');
    if (comments.length === 0) {
        commentsContainer.innerHTML = '<div class="no-posts">No comments yet. Be the first to comment!</div>';
    } else {
        commentsContainer.innerHTML = comments.map(comment => `
            <div class="comment-card">
                <div class="comment-author">${escapeHtml(comment.author)}</div>
                <div class="comment-date">${formatDate(comment.created_at)}</div>
                <div class="comment-content">${escapeHtml(comment.content)}</div>
            </div>
        `).join('');
    }
}

async function handleCreateComment(e) {
    e.preventDefault();

    if (!state.currentPost) return;

    const content = document.getElementById('comment-content').value;

    try {
        const response = await fetch(`${API_URL}/posts/${state.currentPost.id}/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': state.token
            },
            body: JSON.stringify({ content })
        });

        const data = await response.json();

        if (data.success) {
            document.getElementById('comment-content').value = '';
            // Reload post to get updated comments
            openPost(state.currentPost.id);
        }
    } catch (error) {
        console.error('Create comment error:', error);
    }
}

function closePostModal() {
    document.getElementById('post-modal').classList.remove('active');
    document.getElementById('post-modal').classList.add('hidden');
    state.currentPost = null;
}

// Utility functions
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'just now';
    if (diffMins < 60) return `${diffMins} minute${diffMins > 1 ? 's' : ''} ago`;
    if (diffHours < 24) return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`;
    if (diffDays < 7) return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`;
    
    return date.toLocaleDateString();
}

// WebSocket Management
function connectWebSocket() {
    if (!state.token) return;

    updateConnectionStatus('connecting', 'Connecting...');

    // Use ws:// for localhost, wss:// for production
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?token=${state.token}`;

    try {
        state.ws = new WebSocket(wsUrl);

        state.ws.onopen = () => {
            console.log('WebSocket connected');
            updateConnectionStatus('online', 'Connected');
            state.wsReconnectAttempts = 0;
        };

        state.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                handleWebSocketMessage(message);
            } catch (error) {
                console.error('WebSocket message parse error:', error);
            }
        };

        state.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            updateConnectionStatus('offline', 'Connection error');
        };

        state.ws.onclose = () => {
            console.log('WebSocket disconnected');
            updateConnectionStatus('offline', 'Disconnected');
            
            // Attempt to reconnect
            if (state.wsReconnectAttempts < state.wsMaxReconnectAttempts) {
                state.wsReconnectAttempts++;
                console.log(`Reconnecting... Attempt ${state.wsReconnectAttempts}`);
                setTimeout(connectWebSocket, 3000);
            } else {
                updateConnectionStatus('offline', 'Connection failed');
            }
        };
    } catch (error) {
        console.error('WebSocket connection error:', error);
        updateConnectionStatus('offline', 'Connection failed');
    }
}

function handleWebSocketMessage(message) {
    console.log('WebSocket message:', message);

    switch (message.type) {
        case 'user_status':
            handleUserStatus(message.payload);
            break;
        case 'online_users':
            handleOnlineUsers(message.payload);
            break;
        case 'new_message':
            handleNewMessage(message.payload);
            break;
        case 'pong':
            // Heartbeat response
            break;
        default:
            console.log('Unknown message type:', message.type);
    }
}

function handleOnlineUsers(payload) {
    state.onlineUsers = payload.users || [];
    updateOnlineCount();
    renderOnlineUsers();
    console.log(`${state.onlineUsers.length} users online`);
}

function updateConnectionStatus(status, text) {
    const indicator = document.getElementById('ws-status');
    const statusText = document.getElementById('ws-status-text');

    indicator.className = `status-indicator ${status}`;
    statusText.textContent = text;
}

function updateOnlineCount() {
    const countElement = document.getElementById('online-count');
    if (countElement) {
        countElement.textContent = state.onlineUsers.length;
    }
}

function disconnectWebSocket() {
    if (state.ws) {
        state.ws.close();
        state.ws = null;
    }
}

// Messaging Functions
function handleUserStatus(payload) {
    const { user_id, nickname, online } = payload;
    
    if (online) {
        // User came online
        if (!state.onlineUsers.find(u => u.user_id === user_id)) {
            state.onlineUsers.push({ user_id, nickname, online: true });
        }
        console.log(`${nickname} is now online`);
    } else {
        // User went offline
        state.onlineUsers = state.onlineUsers.filter(u => u.user_id !== user_id);
        console.log(`${nickname} is now offline`);
    }

    updateOnlineCount();
    renderOnlineUsers();
}

function renderOnlineUsers() {
    const container = document.getElementById('online-users-container');
    
    if (state.onlineUsers.length === 0) {
        container.innerHTML = '<div class="no-conversations">No users online</div>';
        return;
    }

    // Filter out current user
    const otherUsers = state.onlineUsers.filter(u => u.user_id !== state.currentUser.id);

    container.innerHTML = otherUsers.map(user => `
        <div class="user-item" onclick="openChat(${user.user_id}, '${escapeHtml(user.nickname)}')">
            <div class="user-item-name">${escapeHtml(user.nickname)}</div>
            <div class="user-item-status">● Online</div>
        </div>
    `).join('');
}

async function loadConversations() {
    try {
        const response = await fetch(`${API_URL}/messages/conversations`, {
            headers: {
                'Authorization': state.token
            }
        });

        const data = await response.json();
        if (data.success) {
            state.conversations = data.conversations || [];
            renderConversations();
        }
    } catch (error) {
        console.error('Load conversations error:', error);
    }
}

function renderConversations() {
    const container = document.getElementById('conversations-container');
    
    if (state.conversations.length === 0) {
        container.innerHTML = '<div class="no-conversations">No conversations yet</div>';
        return;
    }

    container.innerHTML = state.conversations.map(conv => `
        <div class="conversation-item" onclick="openChat(${conv.user_id}, '${escapeHtml(conv.nickname)}')">
            <div class="conversation-name">
                ${escapeHtml(conv.nickname)}
                ${conv.unread_count > 0 ? `<span class="conversation-unread">${conv.unread_count}</span>` : ''}
            </div>
            <div class="conversation-preview">${escapeHtml(conv.last_message)}</div>
        </div>
    `).join('');
}

async function openChat(userId, nickname) {
    state.currentChatUser = { id: userId, nickname };
    
    document.getElementById('chat-user-name').textContent = nickname;
    document.getElementById('message-panel').classList.remove('hidden');
    
    // Load message history
    await loadMessages(userId);
}

async function loadMessages(userId) {
    try {
        const response = await fetch(`${API_URL}/messages/${userId}`, {
            headers: {
                'Authorization': state.token
            }
        });

        const data = await response.json();
        if (data.success) {
            state.messages = data.messages || [];
            renderMessages();
        }
    } catch (error) {
        console.error('Load messages error:', error);
    }
}

function renderMessages() {
    const container = document.getElementById('messages-container');
    
    if (state.messages.length === 0) {
        container.innerHTML = '<div class="no-messages">No messages yet. Start the conversation!</div>';
        return;
    }

    container.innerHTML = state.messages.map(msg => {
        const isSent = msg.sender_id === state.currentUser.id;
        return `
            <div class="message-bubble ${isSent ? 'sent' : 'received'}">
                ${!isSent ? `<div class="message-sender">${escapeHtml(msg.sender_name)}</div>` : ''}
                <div>${escapeHtml(msg.content)}</div>
                <div class="message-time">${formatDate(msg.created_at)}</div>
            </div>
        `;
    }).join('');

    // Scroll to bottom
    container.scrollTop = container.scrollHeight;
}

async function handleSendMessage(e) {
    e.preventDefault();

    if (!state.currentChatUser) return;

    const content = document.getElementById('message-input').value.trim();
    if (!content) return;

    try {
        const response = await fetch(`${API_URL}/messages/${state.currentChatUser.id}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': state.token
            },
            body: JSON.stringify({ content })
        });

        const data = await response.json();
        if (data.success) {
            document.getElementById('message-input').value = '';
            // Message will be added via WebSocket
        }
    } catch (error) {
        console.error('Send message error:', error);
    }
}

function handleNewMessage(message) {
    console.log('New message received:', message);

    // If chat is open with this user, add message to view
    if (state.currentChatUser && 
        (message.sender_id === state.currentChatUser.id || 
         message.receiver_id === state.currentChatUser.id)) {
        state.messages.push(message);
        renderMessages();
    }

    // Reload conversations to update unread count
    loadConversations();
}

function closeChatPanel() {
    document.getElementById('message-panel').classList.add('hidden');
    state.currentChatUser = null;
    state.messages = [];
}

// Keyboard Navigation
function setupKeyboardNavigation() {
    // ESC key to close modals
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            // Close post modal
            const postModal = document.getElementById('post-modal');
            if (postModal && !postModal.classList.contains('hidden')) {
                closePostModal();
            }
            
            // Close chat panel
            const messagePanel = document.getElementById('message-panel');
            if (messagePanel && !messagePanel.classList.contains('hidden')) {
                closeChatPanel();
            }
        }
    });

    // Enter key to submit forms (when not in textarea)
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && !e.shiftKey && e.target.tagName !== 'TEXTAREA') {
            const form = e.target.closest('form');
            if (form) {
                e.preventDefault();
                form.dispatchEvent(new Event('submit'));
            }
        }
    });
}

// Performance: Debounce function
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Performance: Throttle function
function throttle(func, limit) {
    let inThrottle;
    return function(...args) {
        if (!inThrottle) {
            func.apply(this, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    };
}

