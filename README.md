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

## Architecture

### Backend Structure

```
backend/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/    # HTTP request handlers
│   │   ├── middleware/  # HTTP middleware
│   │   └── router/      # Route definitions
│   ├── config/          # Configuration management
│   ├── domain/          # Domain models and DTOs
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic layer
│   └── websocket/       # WebSocket hub and client
├── pkg/
│   └── logger/          # Logging package
├── .env.example         # Environment variables template
├── Dockerfile           # Docker image definition
├── Makefile             # Build and development commands
└── go.mod               # Go dependencies
```

### Frontend Structure

```
frontend/
├── app.js               # Main application logic
├── index.html           # HTML template
└── styles.css           # Styles
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (optional)
- Make (optional, for using Makefile commands)

### Local Development

1. **Clone the repository**

```bash
git clone <repository-url>
cd real-time-forum
```

2. **Set up environment variables**

```bash
cd backend
cp .env.example .env
# Edit .env with your configuration
```

3. **Install dependencies**

```bash
make deps
```

4. **Run the application**

```bash
make run
```

The server will start on `http://localhost:8000`

### Using Docker

1. **Build and run with Docker Compose**

```bash
docker-compose up -d
```

2. **View logs**

```bash
docker-compose logs -f
```

3. **Stop the application**

```bash
docker-compose down
```

## Configuration

Configuration is managed through environment variables. See `.env.example` for all available options.

### Key Configuration Options

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `PORT` | Server port | `8000` |
| `DATABASE_PATH` | SQLite database file path | `./data/forum.db` |
| `SESSION_DURATION` | Session expiration time | `24h` |
| `RATE_LIMIT_ENABLED` | Enable rate limiting | `true` |
| `RATE_LIMIT_RPM` | Requests per minute | `100` |

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user

### Posts

- `GET /api/posts` - List posts (with optional category filter)
- `POST /api/posts` - Create a new post
- `GET /api/posts/{id}` - Get a single post with comments
- `POST /api/posts/{id}/comments` - Create a comment

### Messages

- `GET /api/messages/conversations` - Get all conversations
- `GET /api/messages/{id}` - Get messages with a user
- `POST /api/messages/{id}` - Send a message to a user

### WebSocket

- `GET /ws?token={token}` - WebSocket connection for real-time updates

### Health Check

- `GET /health` - Health check endpoint

## Development

### Available Make Commands

```bash
make help          # Show available commands
make build         # Build the application
make run           # Run the application
make test          # Run tests
make clean         # Clean build artifacts
make docker-build  # Build Docker image
make docker-run    # Run Docker container
make lint          # Run linter
make fmt           # Format code
make deps          # Install dependencies
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

## Deployment

### Docker Deployment

1. **Build the Docker image**

```bash
docker build -t forum-backend:latest ./backend
```

2. **Run the container**

```bash
docker run -p 8000:8000 \
  -e ENVIRONMENT=production \
  -e DATABASE_PATH=/app/data/forum.db \
  -v forum-data:/app/data \
  forum-backend:latest
```

### Docker Compose Deployment

```bash
docker-compose up -d
```

### Production Considerations

1. **Security**
   - Use strong session secrets
   - Configure proper CORS origins
   - Enable rate limiting
   - Use HTTPS in production

2. **Database**
   - Regular backups of SQLite database
   - Consider using PostgreSQL for high-traffic sites

3. **Monitoring**
   - Set up health check monitoring
   - Configure log aggregation
   - Monitor WebSocket connections

4. **Performance**
   - Adjust rate limiting based on traffic
   - Configure connection pool sizes
   - Use reverse proxy (nginx) for static files

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

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, please open an issue in the GitHub repository.

## Acknowledgments

- Built with [Gorilla Mux](https://github.com/gorilla/mux) for routing
- [Gorilla WebSocket](https://github.com/gorilla/websocket) for WebSocket support
- [SQLite](https://www.sqlite.org/) for the database
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) for password hashing

