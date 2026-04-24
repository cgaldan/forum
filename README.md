# Real-Time Forum

A modern, production-ready real-time forum application with WebSocket support, built with Go and vanilla JavaScript.

## Features

- 🔐 **User Authentication** - Secure registration and login with bcrypt password hashing
- 💬 **Real-time Messaging** - Private messaging between users using WebSockets
- 📝 **Forum Posts** - Create and view posts with categories
- 💭 **Comments** - Comment on posts with real-time updates
- 👥 **Online Users** - See who's currently online
- 🚀 **WebSocket Support** - Real-time updates for messages and user status
- 🔒 **Security** - Rate limiting, CORS, security headers, and more
- 🐳 **Docker Support** - Easy deployment with Docker and Docker Compose
- 📊 **Health Checks** - Built-in health check endpoint
- 📝 **Structured Logging** - Comprehensive logging with log levels

## Deployment Process

Deployment scenarios are being analyzed here
[DEPLOYMENT.md](/DEPLOYMENT.md)

## Project Structure

You can find detailed project structure here
[PROJECT_STRUCTURE.md](/PROJECT_STRUCTURE.md)

## Configuration

Configuration is managed through environment variables. See `.env.example` for all available options.

### Key Configuration Options

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `PORT` | Server port | `8000` |
| `SERVER_READ_TIMEOUT` | Server read timeout | `5s` |
| `SERVER_WRITE_TIMEOUT` | Server write timeout | `15s` |
| `SERVER_IDLE_TIMEOUT` | Server idle timeout | `60s` |
| `DATABASE_PATH` | SQLite database file path | `./data/database/forum.db` |
| `SESSION_DURATION` | Session expiration time | `24h` |
| `RATE_LIMIT_ENABLED` | Enable rate limiting | `true` |
| `RATE_LIMIT_RPM` | Requests per minute | `100` |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | `*` |
| `WS_READ_BUFFER_SIZE` | WebSocket read buffer size | `1024` |
| `WS_WRITE_BUFFER_SIZE` | WebSocket write buffer size | `1024` |
| `WS_PING_PERIOD` | WebSocket ping period | `54s` |
| `WS_PONG_WAIT` | WebSocket pong wait time | `60s` |
| `WS_WRITE_WAIT` | WebSocket write wait time | `10s` |
| `FRONTEND_PATH` | Path to frontend files | `../frontend` |

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user

### Posts

- `GET /api/posts` - List posts (with optional category, limit, offset filters)
- `POST /api/posts` - Create a new post
- `GET /api/posts/{id}` - Get a single post with comments
- `POST /api/posts/{id}/comments` - Create a comment

### Messages

- `GET /api/messages/conversations` - Get all conversations
- `GET /api/messages/{id}` - Get messages with a user (with optional limit, offset)
- `POST /api/messages/{id}` - Send a message to a user

### WebSocket

- `GET /ws?token={token}` - WebSocket connection for real-time updates

### Health Check

- `GET /health` - Health check endpoint

## Development

### Available Make Commands

```bash
make help                 # Show available commands
make build                # Build the application
make run                  # Run the application
make test                 # Run tests
make clean                # Clean build artifacts
make docker-build         # Build Docker image
make docker-run           # Run Docker container
make docker-compose-up    # Run Docker Compose
make docker-compose-down  # Remove Docker Container
make deps                 # Install dependencies
```

### Running Tests

```bash
make test
```

This will run all tests with race detection and generate a coverage report.

### Code Formatting

```bash
make fmt
```

### Linting

```bash
make lint
```

## Testing

The project includes comprehensive tests for:

- Repository layer
- Service layer
- Handlers
- Middleware
- WebSocket functionality

Run tests with:

```bash
go test -v -race -coverprofile=coverage.out ./...
```

## Database Schema

### Users
- `id` - Primary key
- `nickname` - Unique username
- `email` - Unique email
- `password_hash` - Hashed password
- `first_name`, `last_name` - User's name
- `age` - User's age
- `gender` - User's gender
- `created_at` - Registration timestamp
- `last_seen` - Last activity timestamp

### Posts
- `id` - Primary key
- `user_id` - Foreign key to users
- `title` - Post title
- `content` - Post content
- `category` - Post category
- `created_at`, `updated_at` - Timestamps

### Comments
- `id` - Primary key
- `post_id` - Foreign key to posts
- `user_id` - Foreign key to users
- `content` - Comment content
- `created_at` - Timestamp

### Messages
- `id` - Primary key
- `sender_id` - Foreign key to users
- `receiver_id` - Foreign key to users
- `content` - Message content
- `created_at` - Timestamp
- `read_at` - Read timestamp (nullable)

### Sessions
- `id` - Session ID (primary key)
- `user_id` - Foreign key to users
- `created_at` - Creation timestamp
- `expires_at` - Expiration timestamp

## WebSocket Protocol

### Message Types

#### From Server

**Online Users**
```json
{
  "type": "online_users",
  "payload": {
    "users": [
      {"user_id": 1, "nickname": "john", "online": true}
    ]
  }
}
```

**User Status**
```json
{
  "type": "user_status",
  "payload": {
    "user_id": 1,
    "nickname": "john",
    "online": true
  }
}
```

**New Message**
```json
{
  "type": "new_message",
  "payload": {
    "id": 1,
    "sender_id": 1,
    "receiver_id": 2,
    "content": "Hello!",
    "created_at": "2024-01-01T00:00:00Z",
    "sender_name": "john"
  }
}
```

#### From Client

**Ping**
```json
{
  "type": "ping"
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, please open an issue in the GitHub repository.

## Acknowledgments

- Built with [Gorilla Mux](https://github.com/gorilla/mux) for routing
- [Gorilla WebSocket](https://github.com/gorilla/websocket) for WebSocket support
- [SQLite](https://www.sqlite.org/) for the database
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) for password hashing


## Authors

- Christos Gkaldanidis
- Christos Markos

Creators and primary Developers