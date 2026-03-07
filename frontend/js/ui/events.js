import { getElement } from '../utils/helpers.js';
import { 
    handleLogin, 
    handleLogout, 
    handleRegister 
} from '../modules/auth.js';
import { 
    handleCreatePost, 
    handleCategoryFilter, 
    handleCreateComment, 
    closePostModal 
} from '../modules/posts.js';
import { 
    handleSendMessage, 
    closeChatPanel 
} from '../modules/messages.js';
import { switchAuthTab, sidebar } from './ui.js';


export function setupEventListeners() {
    setupAuthListeners();
    setupMainViewListeners();
    setupKeyboardNavigation();
}

function setupAuthListeners() {
    const loginTab = getElement('login-tab');
    const registerTab = getElement('register-tab');
    const loginForm = getElement('loginForm');
    const registerForm = getElement('registerForm');

    if (loginTab) loginTab.addEventListener('click', () => switchAuthTab('login'));
    if (registerTab) registerTab.addEventListener('click', () => switchAuthTab('register'));
    if (loginForm) loginForm.addEventListener('submit', handleLogin);
    if (registerForm) registerForm.addEventListener('submit', handleRegister);
}

function setupMainViewListeners() {
    const logoutBtn = getElement('logout-btn');
    if (logoutBtn) logoutBtn.addEventListener('click', handleLogout);

    const createPostForm = getElement('createPostForm');
    if (createPostForm) createPostForm.addEventListener('submit', handleCreatePost);

    const categoryFilter = getElement('category-filter');
    if (categoryFilter) categoryFilter.addEventListener('change', handleCategoryFilter);

    const createCommentForm = getElement('createCommentForm');
    if (createCommentForm) createCommentForm.addEventListener('submit', handleCreateComment);

    const closeModal = document.querySelector('.close-modal');
    if (closeModal) closeModal.addEventListener('click', closePostModal);

    const modal = getElement('post-modal');
    if (modal) {
        modal.addEventListener('click', (e) => {
            if (e.target === modal) closePostModal();
        });
    }

    const sendMessageForm = getElement('sendMessageForm');
    if (sendMessageForm) sendMessageForm.addEventListener('submit', handleSendMessage);

    const closeChat = getElement('close-chat');
    if (closeChat) closeChat.addEventListener('click', closeChatPanel);

    const toggleSidebar = getElement('toggle-sidebar');
    if (toggleSidebar) toggleSidebar.addEventListener('click', sidebar);
}


function setupKeyboardNavigation() {
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            const postModal = getElement('post-modal');
            if (postModal && !postModal.classList.contains('hidden')) {
                closePostModal();
            }

            const messagePanel = getElement('message-panel');
            if (messagePanel && !messagePanel.classList.contains('hidden')) {
                closeChatPanel();
            }
        }
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && !e.shiftKey && e.target.tagName !== 'TEXTAREA') {
            const form = e.target.closest('form');
            if (form) {
                e.preventDefault();
                form.dispatchEvent(new Event('submit'));
            }
        }
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && !e.shiftKey && e.target.id === 'message-input') {
            e.preventDefault();
            const form = e.target.closest('form');
            if (form) {
                form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
            }
        }
    });
}
