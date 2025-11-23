# API Documentation

## Base URL

```
http://localhost:8000/api
```

## Authentication

Most endpoints require authentication via a Bearer token in the `Authorization` header.

```
Authorization: {session_token}
```

---

## Authentication Endpoints

### Register

Create a new user account.

**Endpoint:** `POST /api/auth/register`

**Request Body:**
```json
{
  "nickname": "john_doe",
  "email": "john@example.com",
  "password": "securePassword123",
  "first_name": "John",
  "last_name": "Doe",
  "age": 25,
  "gender": "male"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "nickname": "john_doe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "age": 25,
    "gender": "male",
    "created_at": "2024-01-01T00:00:00Z",
    "last_seen": "2024-01-01T00:00:00Z"
  },
  "token": "session_token_here"
}
```

**Validation Rules:**
- `nickname`: Minimum 3 characters, unique
- `email`: Valid email format, unique
- `password`: Minimum 6 characters
- `age`: Between 13 and 120
- All fields are required

---

### Login

Authenticate an existing user.

**Endpoint:** `POST /api/auth/login`

**Request Body:**
```json
{
  "identifier": "john_doe",
  "password": "securePassword123"
}
```

The `identifier` can be either the nickname or email.

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "user": {
    "id": 1,
    "nickname": "john_doe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "age": 25,
    "gender": "male",
    "created_at": "2024-01-01T00:00:00Z",
    "last_seen": "2024-01-01T00:00:00Z"
  },
  "token": "session_token_here"
}
```

---

### Logout

End the current user session.

**Endpoint:** `POST /api/auth/logout`

**Headers:**
```
Authorization: {session_token}
```

**Response:**
```json
{
  "success": true,
  "message": "Logout successful"
}
```

---

### Get Current User

Get information about the currently authenticated user.

**Endpoint:** `GET /api/auth/me`

**Headers:**
```
Authorization: {session_token}
```

**Response:**
```json
{
  "success": true,
  "user": {
    "id": 1,
    "nickname": "john_doe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "age": 25,
    "gender": "male",
    "created_at": "2024-01-01T00:00:00Z",
    "last_seen": "2024-01-01T00:00:00Z"
  }
}
```

---

## Post Endpoints

### List Posts

Get a list of posts with optional filtering and pagination.

**Endpoint:** `GET /api/posts`

**Query Parameters:**
- `category` (optional): Filter by category
- `limit` (optional): Number of posts per page (default: 20, max: 100)
- `offset` (optional): Pagination offset (default: 0)

**Example:**
```
GET /api/posts?category=technology&limit=10&offset=0
```

**Response:**
```json
{
  "success": true,
  "posts": [
    {
      "id": 1,
      "user_id": 1,
      "title": "My First Post",
      "content": "This is the content of my post...",
      "category": "technology",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "author": "john_doe"
    }
  ]
}
```

---

### Get Single Post

Get a single post with its comments.

**Endpoint:** `GET /api/posts/{id}`

**Response:**
```json
{
  "success": true,
  "post": {
    "id": 1,
    "user_id": 1,
    "title": "My First Post",
    "content": "This is the content of my post...",
    "category": "technology",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "author": "john_doe"
  },
  "comments": [
    {
      "id": 1,
      "post_id": 1,
      "user_id": 2,
      "content": "Great post!",
      "created_at": "2024-01-01T00:00:00Z",
      "author": "jane_doe"
    }
  ]
}
```

---

### Create Post

Create a new post.

**Endpoint:** `POST /api/posts`

**Headers:**
```
Authorization: {session_token}
```

**Request Body:**
```json
{
  "title": "My First Post",
  "content": "This is the content of my post...",
  "category": "technology"
}
```

**Validation Rules:**
- `title`: Minimum 3 characters
- `content`: Minimum 10 characters
- `category`: Required (e.g., "general", "technology", "gaming", "music", "sports", "other")

**Response:**
```json
{
  "success": true,
  "message": "Post created successfully",
  "post": {
    "id": 1,
    "user_id": 1,
    "title": "My First Post",
    "content": "This is the content of my post...",
    "category": "technology",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "author": "john_doe"
  }
}
```

---

### Create Comment

Add a comment to a post.

**Endpoint:** `POST /api/posts/{id}/comments`

**Headers:**
```
Authorization: {session_token}
```

**Request Body:**
```json
{
  "content": "This is my comment on the post."
}
```

**Validation Rules:**
- `content`: Must not be empty

**Response:**
```json
{
  "success": true,
  "message": "Comment created successfully",
  "comment": {
    "id": 1,
    "post_id": 1,
    "user_id": 2,
    "content": "This is my comment on the post.",
    "created_at": "2024-01-01T00:00:00Z",
    "author": "jane_doe"
  }
}
```

---

## Message Endpoints

### Get Conversations

Get a list of all conversations for the current user.

**Endpoint:** `GET /api/messages/conversations`

**Headers:**
```
Authorization: {session_token}
```

**Response:**
```json
{
  "success": true,
  "conversations": [
    {
      "user_id": 2,
      "nickname": "jane_doe",
      "last_message": "Hey, how are you?",
      "last_time": "2024-01-01T00:00:00Z",
      "unread_count": 3
    }
  ]
}
```

---

### Get Messages

Get message history with a specific user.

**Endpoint:** `GET /api/messages/{id}`

**Headers:**
```
Authorization: {session_token}
```

**Query Parameters:**
- `limit` (optional): Number of messages per page (default: 50, max: 100)
- `offset` (optional): Pagination offset (default: 0)

**Example:**
```
GET /api/messages/2?limit=20&offset=0
```

**Response:**
```json
{
  "success": true,
  "messages": [
    {
      "id": 1,
      "sender_id": 1,
      "receiver_id": 2,
      "content": "Hello!",
      "created_at": "2024-01-01T00:00:00Z",
      "read_at": "2024-01-01T00:01:00Z",
      "sender_name": "john_doe"
    },
    {
      "id": 2,
      "sender_id": 2,
      "receiver_id": 1,
      "content": "Hi there!",
      "created_at": "2024-01-01T00:02:00Z",
      "read_at": null,
      "sender_name": "jane_doe"
    }
  ]
}
```

---

### Send Message

Send a private message to a user.

**Endpoint:** `POST /api/messages/{id}`

**Headers:**
```
Authorization: {session_token}
```

**Request Body:**
```json
{
  "content": "Hello! How are you?"
}
```

**Validation Rules:**
- `content`: Must be 1-1000 characters
- Cannot send messages to yourself

**Response:**
```json
{
  "success": true,
  "message": "Message sent",
  "msg": {
    "id": 1,
    "sender_id": 1,
    "receiver_id": 2,
    "content": "Hello! How are you?",
    "created_at": "2024-01-01T00:00:00Z",
    "sender_name": "john_doe"
  }
}
```

---

## User Endpoints

### Get Online Users

Get a list of currently online users.

**Endpoint:** `GET /api/users/online`

**Headers:**
```
Authorization: {session_token}
```

**Response:**
```json
{
  "success": true,
  "users": [
    {
      "user_id": 1,
      "nickname": "john_doe",
      "online": true
    },
    {
      "user_id": 2,
      "nickname": "jane_doe",
      "online": true
    }
  ]
}
```

---

## WebSocket Connection

### Connect to WebSocket

Establish a WebSocket connection for real-time updates.

**Endpoint:** `GET /ws?token={session_token}`

**Protocol:** WebSocket

**Connection Example:**
```javascript
const token = 'your_session_token';
const ws = new WebSocket(`ws://localhost:8000/ws?token=${token}`);

ws.onopen = () => {
  console.log('Connected to WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received message:', message);
};
```

### WebSocket Message Types

#### Online Users (Server → Client)
```json
{
  "type": "online_users",
  "payload": {
    "users": [
      {"user_id": 1, "nickname": "john_doe", "online": true}
    ]
  }
}
```

#### User Status (Server → Client)
```json
{
  "type": "user_status",
  "payload": {
    "user_id": 1,
    "nickname": "john_doe",
    "online": true
  }
}
```

#### New Message (Server → Client)
```json
{
  "type": "new_message",
  "payload": {
    "id": 1,
    "sender_id": 1,
    "receiver_id": 2,
    "content": "Hello!",
    "created_at": "2024-01-01T00:00:00Z",
    "sender_name": "john_doe"
  }
}
```

#### Ping (Client → Server)
```json
{
  "type": "ping"
}
```

#### Pong (Server → Client)
```json
{
  "type": "pong"
}
```

---

## Health Check

### Health Status

Check the health status of the application.

**Endpoint:** `GET /health`

**No Authentication Required**

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0"
}
```

---

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "success": false,
  "message": "Error description here"
}
```

### Common HTTP Status Codes

- `200 OK`: Request successful
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required or invalid
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Default Limit:** 100 requests per minute per IP address
- **Response:** `429 Too Many Requests` when limit exceeded

---

## CORS

The API supports Cross-Origin Resource Sharing (CORS) for browser-based applications.

**Allowed Methods:** GET, POST, PUT, DELETE, OPTIONS
**Allowed Headers:** Content-Type, Authorization

---

## Categories

Available post categories:
- `general`
- `technology`
- `gaming`
- `music`
- `sports`
- `other`

