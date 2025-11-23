# Frontend Documentation

## Overview

The frontend is built with vanilla JavaScript, HTML, and CSS. It provides a modern, responsive user interface for the real-time forum application.

## Architecture

### Core Components

- **index.html** - Main HTML structure
- **app.js** - Application logic and state management
- **styles.css** - Styling and responsive design

### State Management

The application uses a simple state object to manage global state:

```javascript
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
    conversations: [],
    messageOffset: 0,
    hasMoreMessages: true,
    isLoadingMessages: false
};
```

## Features

### Authentication
- User registration with validation
- Login with nickname or email
- Session management with localStorage
- Automatic session restoration

### Forum Posts
- Create posts with categories
- View posts with pagination
- Filter posts by category
- View post details with comments
- Add comments to posts

### Real-time Messaging
- Private messaging between users
- Real-time message delivery via WebSocket
- Message pagination (load older messages)
- Conversation list with unread counts
- Online/offline user status

### WebSocket Integration
- Automatic connection on login
- Reconnection on disconnect
- Real-time updates for:
  - New messages
  - User online/offline status
  - Online users list

## Code Structure

### Initialization
```javascript
document.addEventListener('DOMContentLoaded', () => {
    initializeApp();
    setupKeyboardNavigation();
});
```

### Event Listeners
- Authentication form handlers
- Post and comment handlers
- Message handlers
- WebSocket message handlers
- Keyboard navigation

### API Communication
All API calls use the Fetch API with error handling:

```javascript
try {
    const response = await fetch(url, options);
    const data = await response.json();
    
    if (data.success) {
        // Handle success
    } else {
        // Handle error
    }
} catch (error) {
    console.error('Error:', error);
    showToast('Error message', 'error');
}
```

### WebSocket Protocol
```javascript
// Connection
const ws = new WebSocket(`${protocol}//${host}/ws?token=${token}`);

// Message handling
ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    handleWebSocketMessage(message);
};
```

## UI/UX Features

### Responsive Design
- Mobile-first approach
- Flexible layout for all screen sizes
- Touch-friendly interface

### Accessibility
- ARIA labels and roles
- Keyboard navigation
- Screen reader support
- Focus management

### User Feedback
- Toast notifications for actions
- Loading indicators
- Error messages
- Real-time status updates

### Keyboard Shortcuts
- **ESC** - Close modals/panels
- **Enter** - Submit forms
- **Shift+Enter** - New line in message input

## Customization

### Styling
Edit `styles.css` to customize:
- Colors and themes
- Layout and spacing
- Typography
- Component styles

### Configuration
Edit `app.js` to customize:
- API URL
- WebSocket reconnection attempts
- Pagination limits
- Message limits

## Browser Support

Tested and supported on:
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Performance Optimizations

### Implemented
- Debouncing for scroll events
- Throttling for frequent operations
- Efficient DOM updates
- Minimal reflows and repaints
- WebSocket message batching

### Best Practices
- Avoid unnecessary API calls
- Cache data when appropriate
- Use event delegation
- Minimize DOM manipulation
- Lazy load content

## Security

### Client-side Security
- XSS prevention via escapeHtml()
- Token storage in localStorage
- HTTPS enforcement (production)
- Input validation

### Token Management
```javascript
// Store token
localStorage.setItem('token', token);

// Retrieve token
const token = localStorage.getItem('token');

// Clear token
localStorage.removeItem('token');
```

## Development

### Local Development
Simply open `index.html` in a browser or use a local server:

```bash
# Python
python3 -m http.server 8080

# Node.js
npx http-server -p 8080
```

### Testing
- Manual testing in multiple browsers
- Test WebSocket reconnection
- Test offline behavior
- Test responsive design

## Future Enhancements

Potential improvements:
- Service Worker for offline support
- Push notifications
- File upload support
- Emoji picker
- Rich text editor
- Image preview
- User profiles
- Search functionality
- Advanced filtering
- Theme switcher (dark mode)

## Troubleshooting

### WebSocket connection fails
- Check if backend is running
- Verify token is valid
- Check browser console for errors
- Ensure CORS is configured

### Messages not updating
- Check WebSocket connection status
- Verify token hasn't expired
- Check browser console for errors
- Try refreshing the page

### UI not responsive
- Clear browser cache
- Hard refresh (Ctrl+Shift+R)
- Check for JavaScript errors
- Verify CSS is loading

## Contributing

When contributing to the frontend:
1. Follow existing code style
2. Add comments for complex logic
3. Test in multiple browsers
4. Ensure responsive design
5. Update documentation

## License

Same as the main project license.

