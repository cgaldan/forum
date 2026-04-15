import { initAuth } from './modules/auth.js';
import { setupEventListeners } from './ui/events.js';

async function initializeApp() {
    try {
        await initAuth();

        setupEventListeners();

        console.log('Application initialized');
    } catch (error) {
        console.error('Failed to initialize app:', error);
    }
}

document.addEventListener('DOMContentLoaded', initializeApp);
