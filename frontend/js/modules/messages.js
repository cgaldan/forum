import api from '../api/client.js';
import store from '../state/store.js';
import { escapeHtml, formatDate, getElement, noMoreMessages } from '../utils/helpers.js';
import { CONFIG } from '../config.js';

export async function loadConversations() {
    try {
        const data = await api.getConversations();
        if (data.success) {
            store.set('conversations', data.conversations || []);
            renderConversations();
        }
    } catch (error) {
        console.error('Load conversations error:', error);
    }
}

function renderConversations() {
    const container = getElement('conversations-container');
    const conversations = store.get('conversations');

    if (conversations.length === 0) {
        container.innerHTML = '<div class="no-conversations">No conversations yet</div>';
        return;
    }

    container.innerHTML = conversations.map(conv => `
        <div class="conversation-item" onclick="window.messagesModule.openChat(${conv.user_id}, '${escapeHtml(conv.nickname)}')">
            <div class="conversation-name">
                ${escapeHtml(conv.nickname)}
                ${conv.unread_count > 0 ? `<span class="conversation-unread">${conv.unread_count}</span>` : ''}
            </div>
            <div class="conversation-preview">${escapeHtml(conv.last_message)}</div>
        </div>
    `).join('');
}

let observer;
export async function openChat(userId, nickname) {
    store.setState({
        currentChatUser: { id: userId, nickname },
        messages: [],
        messageOffset: 0,
        hasMoreMessages: true
    });

    const list = getElement('messages-list');
    if (list) list.innerHTML = '';

    getElement('chat-user-name').textContent = nickname;
    getElement('message-panel').classList.remove('hidden');

    await loadMessages(userId, true);
    await loadConversations();

    if (!observer) {
        setupMessageObserver();
    }
}

export async function loadMessages(userId, isInitial = false) {
    if (store.get('isLoadingMessages')) {
        return;
    }

    const state = store.getState();
    if (!state.hasMoreMessages && !isInitial) {
        return;
    }

    store.set('isLoadingMessages', true);
    const loadingIndicator = getElement('loading-more');
    if (loadingIndicator) loadingIndicator.style.display = 'block';

    const container = getElement('messages-container');

    const previousScrollHeight = container.scrollHeight;
    let didLoadNewMessages = false;

    try {
        const data = await api.getMessages(userId, CONFIG.MESSAGE_LOAD_LIMIT, state.messageOffset);
        if (store.get('currentChatUser')?.id !== userId) return;
    
        if (data.success) {
            const newMessages = data.messages || [];

            if (newMessages.length < CONFIG.MESSAGE_LOAD_LIMIT) {
                store.set('hasMoreMessages', false);
            }

            if (newMessages.length > 0) {
                didLoadNewMessages = true;
                const currentMessages = store.get('messages');
                const messages = [...newMessages, ...currentMessages];
                store.setState({
                    messages,
                    messageOffset: state.messageOffset + newMessages.length
                });

                if (isInitial) {
                    renderMessages();
                } else {
                    prependMessages(newMessages);
                }
            }
        }
    } catch (error) {
        console.error('Load messages error:', error);
    } finally {
        store.set('isLoadingMessages', false);
        const loadingIndicator = getElement('loading-more');
        if (loadingIndicator) loadingIndicator.style.display = 'none';
    }

    if (didLoadNewMessages) {
        requestAnimationFrame(() => {
            if (isInitial) {
                container.scrollTop = container.scrollHeight;
            } else {
                container.scrollTop += container.scrollHeight - previousScrollHeight;
            }
        });
    }
}

function setupMessageObserver() {
    const container = getElement('messages-container');
    const trigger = getElement('load-trigger');

    observer = new IntersectionObserver(entries => {
        const entry = entries[0];

        if (!entry.isIntersecting) return;

        const state = store.getState();
        if (state.currentChatUser && state.hasMoreMessages && !state.isLoadingMessages) {
            loadMessages(state.currentChatUser.id, false);
        }

    }, {
        root: container,
        threshold: 0.1
    });

    observer.observe(trigger);
}

function renderMessages() {
    const list = getElement('messages-list');
    const state = store.getState();
    const messages = store.get('messages');
    const currentUser = store.get('currentUser');
    const hasMoreMessages = store.get('hasMoreMessages');
    const loadingIndicator = getElement('loading-more');

    if (!loadingIndicator){
        console.log('Loading indicator not found'); 
        return;
    }

    if (messages.length === 0) {
        list.innerHTML = '<div class="no-messages">No messages yet. Start the conversation!</div>';
        return;
    }

    let html = '';

    noMoreMessages(messages.length, hasMoreMessages, html, CONFIG.MESSAGE_LOAD_LIMIT);

    html += messages.map(msg => {
        const isSent = msg.sender_id === currentUser.id;
        return `
            <div class="message-bubble ${isSent ? 'sent' : 'received'}">
                ${isSent ? '' : `<div class="message-sender">${escapeHtml(msg.sender_name)}</div>`}
                <div>${escapeHtml(msg.content)}</div>
                <div class="message-time">${formatDate(msg.created_at)}</div>
            </div>
        `;
    }).join('');

    list.innerHTML = html;
}

function prependMessages(newMessages) {
    const list = getElement('messages-list');
    const currentUser = store.get('currentUser');
    const hasMoreMessages = store.get('hasMoreMessages');
    const messages = store.get('messages');

    const empty = list.querySelector('.no-messages');
    if (empty) {
        empty.remove();
    }

    const fragment = document.createDocumentFragment();

    noMoreMessages(messages.length, hasMoreMessages, fragment, CONFIG.MESSAGE_LOAD_LIMIT);

    newMessages.forEach(msg => {
        const isSent = msg.sender_id === currentUser.id;

        const bubble = document.createElement('div');
        bubble.classList.add('message-bubble', isSent ? 'sent' : 'received');
        
        bubble.innerHTML = `
            ${isSent ? '' : `<div class="message-sender">${escapeHtml(msg.sender_name)}</div>`}
            <div>${escapeHtml(msg.content)}</div>
            <div class="message-time">${formatDate(msg.created_at)}</div>
        `;
        fragment.appendChild(bubble);
    });
    
    list.prepend(fragment);
}

export async function handleSendMessage(e) {
    e.preventDefault();

    const currentChatUser = store.get('currentChatUser');
    if (!currentChatUser) return;

    const content = getElement('message-input').value.trim();
    if (!content) return;

    try {
        const data = await api.sendMessage(currentChatUser.id, content);
        if (data.success) {
            getElement('message-input').value = '';
        }
    } catch (error) {
        console.error('Send message error:', error);
    }
}

export function handleNewMessage(message) {
    const currentChatUser = store.get('currentChatUser');

    if (currentChatUser &&
        (message.sender_id === currentChatUser.id ||
         message.receiver_id === currentChatUser.id)) {
        const messages = [...store.get('messages'), message];
        store.set('messages', messages);
        renderMessages();

        const container = getElement('messages-container');
        container.scrollTop = container.scrollHeight;
    }

    loadConversations();
}

export function closeChatPanel() {
    const panel = getElement('message-panel');
    panel.classList.add('hidden');

    store.setState({
        currentChatUser: null,
        messages: [],
        messageOffset: 0,
        hasMoreMessages: true
    });
}

window.messagesModule = { openChat };
