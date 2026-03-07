import api from '../api/client.js';
import store from '../state/store.js';
import { escapeHtml, formatDate, getElement } from '../utils/helpers.js';
import { showError, clearError, showToast } from '../ui/ui.js';
import { CONFIG } from '../config.js';

export async function loadPosts(category = '') {
    const container = getElement('posts-container');
    if (container) {
        container.innerHTML = '<div class="loading">Loading posts...</div>';
    }

    try {
        const data = await api.getPosts(category);

        if (data.success) {
            store.set('posts', data.posts || []);
            renderPosts();
        } else {
            throw new Error(data.message || 'Failed to load posts');
        }
    } catch (error) {
        console.error('Load posts error:', error);
        if (container) {
            container.innerHTML = '<div class="no-posts">Failed to load posts. Please try again.</div>';
        }
        showToast('Failed to load posts', 'error');
    }
}

export function renderPosts() {
    const container = getElement('posts-container');
    const posts = store.get('posts');

    if (posts.length === 0) {
        container.innerHTML = '<div class="no-posts">No posts yet. Be the first to post!</div>';
        return;
    }

    container.innerHTML = posts.map(post => `
        <div class="post-card" onclick="window.postsModule.openPost(${post.id})">
            <div class="post-card-header">
                <h3 class="post-title">${escapeHtml(post.title)}</h3>
                <span class="post-category">${escapeHtml(post.category)}</span>
            </div>
            <div class="post-meta">
                by ${escapeHtml(post.author)} • ${formatDate(post.created_at)}
            </div>
            <div class="post-content-preview">
                ${escapeHtml(post.content.substring(0, CONFIG.POST_PREVIEW_LENGTH))}${post.content.length > CONFIG.POST_PREVIEW_LENGTH ? '...' : ''}
            </div>
        </div>
    `).join('');
}

export async function handleCreatePost(e) {
    e.preventDefault();
    clearError('post-error');

    const title = getElement('post-title').value.trim();
    const category = getElement('post-category').value;
    const content = getElement('post-content').value.trim();

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
        const data = await api.createPost(title, category, content);

        if (data.success) {
            getElement('createPostForm').reset();
            showToast('Post created successfully!', 'success');
            loadPosts(store.get('currentCategory'));
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

export function handleCategoryFilter(e) {
    store.set('currentCategory', e.target.value);
    loadPosts(store.get('currentCategory'));
}

export async function openPost(postId) {
    try {
        const data = await api.getPost(postId);

        if (data.success) {
            store.set('currentPost', data.post);
            renderPostDetail(data.post, data.comments || []);
            const modal = getElement('post-modal');
            modal.classList.add('active');
            modal.classList.remove('hidden');
        }
    } catch (error) {
        console.error('Open post error:', error);
    }
}

function renderPostDetail(post, comments) {
    getElement('post-detail').innerHTML = `
        <h2 class="post-detail-title">${escapeHtml(post.title)}</h2>
        <div class="post-detail-meta">
            <span class="post-category">${escapeHtml(post.category)}</span> • 
            by ${escapeHtml(post.author)} • ${formatDate(post.created_at)}
        </div>
        <div class="post-detail-content">${escapeHtml(post.content)}</div>
    `;

    const commentsContainer = getElement('comments-container');
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

export async function handleCreateComment(e) {
    e.preventDefault();

    const currentPost = store.get('currentPost');
    if (!currentPost) return;

    const content = getElement('comment-content').value;

    try {
        const data = await api.createComment(currentPost.id, content);

        if (data.success) {
            getElement('comment-content').value = '';
            openPost(currentPost.id);
        }
    } catch (error) {
        console.error('Create comment error:', error);
    }
}

export function closePostModal() {
    const modal = getElement('post-modal');
    modal.classList.remove('active');
    modal.classList.add('hidden');
    store.set('currentPost', null);
}

window.postsModule = { openPost };
