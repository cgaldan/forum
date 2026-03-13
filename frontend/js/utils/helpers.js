export function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

export function formatDate(dateString) {
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

export function getElement(id) {
    return document.getElementById(id);
}

export function setText(id, text) {
    const el = getElement(id);
    if (el) el.textContent = text;
}

export function toggleClass(id, className) {
    const el = getElement(id);
    if (el) el.classList.toggle(className);
}

export function noMoreMessages(hasMoreMessages, fragment) {
    if (!hasMoreMessages) {
        const noMoreMessages = document.createElement('div');
        noMoreMessages.classList.add('no-more-messages');
        noMoreMessages.textContent = 'No more messages';
        fragment.appendChild(noMoreMessages);
    }
}