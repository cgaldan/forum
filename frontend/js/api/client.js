import { CONFIG } from '../config.js';

class ApiClient {
    constructor(baseUrl = CONFIG.API_URL) {
        this.baseUrl = baseUrl;
        this.token = null;
    }

    setToken(token) {
        this.token = token;
    }

    getHeaders(contentType = 'application/json') {
        const headers = {
            'Content-Type': contentType
        };
        if (this.token) {
            headers['Authorization'] = this.token;
        }
        return headers;
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseUrl}${endpoint}`;
        const config = {
            ...options,
            headers: {
                ...this.getHeaders(options.contentType),
                ...options.headers
            }
        };

        try {
            const response = await fetch(url, config);
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error(`API request failed: ${endpoint}`, error);
            throw error;
        }
    }

    get(endpoint) {
        return this.request(endpoint, { method: 'GET' });
    }

    post(endpoint, data) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    /**
     * Authentication Endpoints
     */

    login(identifier, password) {
        return this.post('/auth/login', { identifier, password });
    }

    register(userData) {
        return this.post('/auth/register', userData);
    }

    logout() {
        return this.post('/auth/logout', {});
    }

    verifySession() {
        return this.get('/auth/me');
    }

    /**
     * Posts Endpoints
     */

    getPosts(category = '') {
        const endpoint = category ? `/posts?category=${category}` : '/posts';
        return this.get(endpoint);
    }

    getPost(postId) {
        return this.get(`/posts/${postId}`);
    }

    createPost(title, category, content) {
        return this.post('/posts', { title, category, content });
    }

    /**
     * Comments Endpoints
     */

    createComment(postId, content) {
        return this.post(`/posts/${postId}/comments`, { content });
    }

    /**
     * Messages Endpoints
     */

    getConversations() {
        return this.get('/messages/conversations');
    }

    getMessages(userId, limit = 10, offset = 0) {
        return this.get(`/messages/${userId}?limit=${limit}&offset=${offset}`);
    }

    sendMessage(userId, content) {
        return this.post(`/messages/${userId}`, { content });
    }
}

export default new ApiClient();
