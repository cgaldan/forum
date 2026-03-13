import { getElement, toggleClass } from '../utils/helpers.js';
import { CONFIG } from '../config.js';
import { escapeHtml, formatDate } from '../utils/helpers.js';
import store from '../state/store.js';

export function showError(elementId, message) {
    const errorElement = getElement(elementId);
    if (errorElement) {
        errorElement.textContent = message;
        errorElement.classList.add('show');
    }
}

export function clearError(elementId) {
    const errorElement = getElement(elementId);
    if (errorElement) {
        errorElement.textContent = '';
        errorElement.classList.remove('show');
    }
}

export function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    toast.setAttribute('role', 'alert');
    toast.setAttribute('aria-live', 'polite');

    document.body.appendChild(toast);

    setTimeout(() => {
        toast.style.animation = 'slideIn 0.3s ease-out reverse';
        setTimeout(() => toast.remove(), 300);
    }, CONFIG.TOAST_DURATION);
}

export function renderOnlineUsers() {
    const container = getElement('online-users-container');
    const onlineUsers = store.get('onlineUsers');
    const currentUser = store.get('currentUser');

    if (!onlineUsers || onlineUsers.length === 0) {
        if (container) {
            container.innerHTML = '<div class="no-conversations">No other users online</div>';
        }
        return;
    }

    const otherUsers = onlineUsers
        .filter(u => u.user_id !== currentUser.id)
        .sort((a, b) => a.nickname.localeCompare(b.nickname));

    if (otherUsers.length === 0) {
        if (container) {
            container.innerHTML = '<div class="no-conversations">No other users online</div>';
        }
        return;
    }

    if (container) {
        container.innerHTML = otherUsers.map(user => `
            <div class="user-item" onclick="window.messagesModule.openChat(${user.user_id}, '${escapeHtml(user.nickname)}')">
                <div class="sidebar-avatar" aria-hidden="true">${escapeHtml(user.nickname.charAt(0))}</div>
                <div class="user-item-info">
                    <div class="user-item-name">${escapeHtml(user.nickname)}</div>
                    <div class="user-item-status">Online</div>
                </div>
            </div>
        `).join('');
    }
}

export function switchAuthTab(tab) {
    const loginTab = getElement('login-tab');
    const registerTab = getElement('register-tab');
    const loginFormDiv = getElement('login-form');
    const registerFormDiv = getElement('register-form');

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

export function sidebar() {
    const sidebarEl = document.querySelector('.messaging-sidebar');
    if (!sidebarEl) return;
    const isMobile = window.innerWidth <= 768;

    if (isMobile) {
        if (sidebarEl.classList.contains('open')) {
            closeSidebar();
        } else {
            sidebarEl.classList.add('open');
            const backdrop = getElement('sidebar-backdrop');
            if (backdrop) backdrop.classList.add('visible');
            document.body.style.overflow = 'hidden';
        }
    } else {
        const isCollapsed = sidebarEl.classList.toggle('collapsed');
        const btn = getElement('toggle-sidebar');
        if (btn) btn.setAttribute('aria-expanded', isCollapsed ? 'false' : 'true');
    }
}

export function closeSidebar() {
    const sidebarEl = document.querySelector('.messaging-sidebar');
    if (!sidebarEl) return;
    const isMobile = window.innerWidth <= 768;
    if (isMobile) {
        sidebarEl.classList.remove('open');
        const backdrop = getElement('sidebar-backdrop');
        if (backdrop) backdrop.classList.remove('visible');
        document.body.style.overflow = '';
    } else {
        sidebarEl.classList.add('collapsed');
        const btn = getElement('toggle-sidebar');
        if (btn) btn.setAttribute('aria-expanded', 'false');
    }
}

export function createMessageHTML(message, isSent = false) {
    return `
        ${isSent ? '' : `<div class="message-sender">${escapeHtml(message.sender_name)}</div>`}
        <div class="message-bubble-body">${escapeHtml(message.content)}</div>
        <div class="message-time">${formatDate(message.created_at)}</div>
    `;
}