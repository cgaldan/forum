class Store {
    constructor() {
        this.state = {
            currentUser: null,
            token: null,
            currentView: 'auth',
            posts: [],
            currentPost: null,
            currentCategory: '',
            ws: null,
            wsReconnectAttempts: 0,
            connectionStatus: 'offline',
            onlineUsers: [],
            currentChatUser: null,
            messages: [],
            conversations: [],
            messageOffset: 0,
            hasMoreMessages: true,
            isLoadingMessages: false
        };

        this.listeners = [];
    }

    getState() {
        return this.state;
    }

    setState(updates) {
        this.state = { ...this.state, ...updates };
        this.notifyListeners();
    }

    get(key) {
        return this.state[key];
    }

    set(key, value) {
        this.setState({ [key]: value });
    }

    subscribe(listener) {
        this.listeners.push(listener);
        return () => {
            this.listeners = this.listeners.filter(l => l !== listener);
        };
    }

    notifyListeners() {
        this.listeners.forEach(listener => listener(this.state));
    }

    clear() {
        this.setState({
            currentUser: null,
            token: null,
            onlineUsers: [],
            currentChatUser: null,
            messages: [],
            conversations: [],
            posts: [],
            currentPost: null,
            messageOffset: 0,
            hasMoreMessages: true,
            isLoadingMessages: false
        });
    }
}

export default new Store();
