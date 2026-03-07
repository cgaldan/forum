import { CONFIG } from '../config.js';
import store from '../state/store.js';
import { getElement } from '../utils/helpers.js';
import { handleNewMessage } from './messages.js';
import { renderOnlineUsers } from '../ui/ui.js';


export function connectWebSocket() {
    const token = store.get('token');
    if (!token) return;

    updateConnectionStatus('connecting', 'Connecting...');

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?token=${token}`;

    try {
        const ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            console.log('WebSocket connected');
            updateConnectionStatus('online', 'Connected');
            store.set('wsReconnectAttempts', 0);
        };

        ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                handleWebSocketMessage(message);
            } catch (error) {
                console.error('WebSocket message parse error:', error);
            }
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            updateConnectionStatus('offline', 'Connection error');
        };

        ws.onclose = () => {
            console.log('WebSocket disconnected');
            updateConnectionStatus('offline', 'Disconnected');

            const attempts = store.get('wsReconnectAttempts');
            if (attempts < CONFIG.WS_RECONNECT_ATTEMPTS) {
                store.set('wsReconnectAttempts', attempts + 1);
                console.log(`Reconnecting... Attempt ${attempts + 1}`);
                setTimeout(connectWebSocket, CONFIG.WS_RECONNECT_DELAY);
            } else {
                updateConnectionStatus('offline', 'Connection failed');
            }
        };

        store.set('ws', ws);
    } catch (error) {
        console.error('WebSocket connection error:', error);
        updateConnectionStatus('offline', 'Connection failed');
    }
}

async function handleWebSocketMessage(message) {
    console.log('WebSocket message:', message);

    switch (message.type) {
        case 'user_status':
            handleUserStatus(message.payload);
            break;
        case 'online_users':
            handleOnlineUsers(message.payload);
            break;
        case 'new_message': {
            handleNewMessage(message.payload);
            break;
        }
        case 'pong':
            break;
        default:
            console.log('Unknown message type:', message.type);
    }
}

async function handleOnlineUsers(payload) {
    store.set('onlineUsers', payload.users || []);
    updateOnlineCount();
    renderOnlineUsers();
    
    console.log(`${store.get('onlineUsers').length} users online`);
}

async function handleUserStatus(payload) {
    const { user_id, nickname, online } = payload;
    const onlineUsers = store.get('onlineUsers');

    if (online) {
        if (!onlineUsers.find(u => u.user_id === user_id)) {
            onlineUsers.push({ user_id, nickname, online: true });
            store.set('onlineUsers', onlineUsers);
        }
        console.log(`${nickname} is now online`);
    } else {
        const filtered = onlineUsers.filter(u => u.user_id !== user_id);
        store.set('onlineUsers', filtered);
        console.log(`${nickname} is now offline`);
    }

    updateOnlineCount();
    renderOnlineUsers();
}

function updateConnectionStatus(status, text) {
    const indicator = getElement('ws-status');
    const statusText = getElement('ws-status-text');

    store.set('connectionStatus', status);
    
    if (indicator) {
        indicator.className = `status-indicator ${status}`;
    }
    if (statusText) {
        statusText.textContent = text;
    }
}

function updateOnlineCount() {
    const countElement = getElement('online-count');
    if (countElement) {
        countElement.textContent = store.get('onlineUsers').length;
    }
}

export function disconnectWebSocket() {
    const ws = store.get('ws');
    if (ws) {
        ws.close();
        store.set('ws', null);
    }
}
