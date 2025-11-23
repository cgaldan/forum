# Project Structure

This document provides a detailed overview of the project's file and directory structure.

## Root Directory

```
real-time-forum/
├── .github/                    # GitHub-specific files
│   └── workflows/             # CI/CD workflows
│       └── ci.yml            # Continuous Integration pipeline
├── backend/                   # Backend Go application
├── frontend/                  # Frontend application
├── .gitignore                # Git ignore rules
├── API.md                    # API documentation
├── CHANGELOG.md              # Version history
├── CONTRIBUTING.md           # Contribution guidelines
├── DEPLOYMENT.md             # Deployment guide
├── docker-compose.yml        # Docker Compose configuration
├── LICENSE                   # MIT License
├── PROJECT_STRUCTURE.md      # This file
└── README.md                 # Main documentation
```

## Backend Structure

```
backend/
├── cmd/                      # Application entry points
│   └── server/              # Main server application
│       └── main.go          # Entry point with initialization
│
├── internal/                # Private application code
│   ├── api/                # API layer
│   │   ├── handlers/       # HTTP request handlers
│   │   │   ├── auth_handler.go        # Authentication endpoints
│   │   │   ├── comment_handler.go     # Comment endpoints
│   │   │   ├── health_handler.go      # Health check endpoint
│   │   │   ├── message_handler.go     # Messaging endpoints
│   │   │   ├── post_handler.go        # Post endpoints
│   │   │   └── websocket_handler.go   # WebSocket endpoint
│   │   │
│   │   ├── middleware/     # HTTP middleware
│   │   │   └── middleware.go          # All middleware functions
│   │   │
│   │   └── router/         # Routing configuration
│   │       └── router.go              # Route definitions
│   │
│   ├── config/             # Configuration management
│   │   └── config.go                  # Configuration loading and validation
│   │
│   ├── domain/             # Domain models and DTOs
│   │   ├── models.go                  # Core domain models
│   │   ├── requests.go                # Request DTOs
│   │   └── responses.go               # Response DTOs
│   │
│   ├── repository/         # Data access layer
│   │   ├── comment_repository.go      # Comment data access
│   │   ├── database.go                # Database initialization
│   │   ├── message_repository.go      # Message data access
│   │   ├── post_repository.go         # Post data access
│   │   ├── repositories.go            # Repository factory
│   │   ├── session_repository.go      # Session data access
│   │   ├── user_repository.go         # User data access
│   │   └── user_repository_test.go    # User repository tests
│   │
│   ├── service/            # Business logic layer
│   │   ├── auth_service.go            # Authentication logic
│   │   ├── auth_service_test.go       # Authentication tests
│   │   ├── comment_service.go         # Comment logic
│   │   ├── message_service.go         # Messaging logic
│   │   ├── post_service.go            # Post logic
│   │   └── services.go                # Service factory
│   │
│   └── websocket/          # WebSocket management
│       ├── client.go                  # WebSocket client
│       └── hub.go                     # WebSocket hub
│
├── pkg/                    # Public packages
│   └── logger/            # Logging package
│       └── logger.go                  # Logger implementation
│
├── .env.example           # Environment variables template
├── Dockerfile             # Docker image definition
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
└── Makefile              # Build automation
```

## Frontend Structure

```
frontend/
├── app.js                 # Main application logic
│   ├── Configuration
│   ├── State management
│   ├── Initialization
│   ├── Authentication handlers
│   ├── Post management
│   ├── Comment handlers
│   ├── Message handlers
│   ├── WebSocket handlers
│   └── Utility functions
│
├── index.html            # HTML structure
│   ├── Authentication view
│   ├── Main forum view
│   ├── Post feed
│   ├── Message panel
│   └── Modal dialogs
│
├── styles.css            # Styling
│   ├── Global styles
│   ├── Layout
│   ├── Components
│   ├── Authentication
│   ├── Posts and comments
│   ├── Messaging
│   └── Responsive design
│
└── README.md            # Frontend documentation
```

## Key Architectural Patterns

### Backend Layers

```
┌─────────────────────────────────────┐
│         HTTP Handlers               │ ← Handle HTTP requests/responses
├─────────────────────────────────────┤
│         Service Layer                │ ← Business logic
├─────────────────────────────────────┤
│       Repository Layer               │ ← Data access
├─────────────────────────────────────┤
│          Database                    │ ← SQLite storage
└─────────────────────────────────────┘
```

### Request Flow

```
Client Request
    ↓
Middleware (Logging, CORS, Rate Limit, etc.)
    ↓
Router
    ↓
Handler (Parse request, validate)
    ↓
Service (Business logic)
    ↓
Repository (Database operations)
    ↓
Database
    ↓
Repository (Return data)
    ↓
Service (Transform data)
    ↓
Handler (Format response)
    ↓
Middleware (Security headers, etc.)
    ↓
Client Response
```

### WebSocket Flow

```
Client
    ↓
WebSocket Handler (Authenticate)
    ↓
Hub (Register client)
    ↓
Client (Bidirectional communication)
    ↓
Hub (Broadcast messages)
    ↓
All Connected Clients
```

## File Responsibilities

### Backend

#### `cmd/server/main.go`
- Application initialization
- Configuration loading
- Database setup
- Service initialization
- Router setup
- Graceful shutdown

#### `internal/api/handlers/`
- HTTP request handling
- Request validation
- Response formatting
- Error handling

#### `internal/service/`
- Business logic
- Data validation
- Data transformation
- Orchestration of repositories

#### `internal/repository/`
- Database queries
- Data persistence
- Data retrieval
- Transaction management

#### `internal/domain/`
- Data structures
- Request/Response DTOs
- Domain models

#### `internal/config/`
- Configuration loading
- Environment variable parsing
- Configuration validation

#### `pkg/logger/`
- Structured logging
- Log levels
- Log formatting

### Frontend

#### `app.js`
- State management
- API communication
- WebSocket handling
- DOM manipulation
- Event handling

#### `index.html`
- Page structure
- Semantic markup
- Accessibility features

#### `styles.css`
- Visual styling
- Layout and positioning
- Responsive design
- Animations

## Configuration Files

### `.env.example`
Template for environment variables with defaults and descriptions.

### `docker-compose.yml`
Docker Compose configuration for containerized deployment.

### `Dockerfile`
Multi-stage Docker build for optimized images.

### `Makefile`
Build automation and development tasks.

### `.github/workflows/ci.yml`
GitHub Actions CI/CD pipeline configuration.

## Data Flow

### Authentication Flow
```
1. User submits credentials
2. Handler receives request
3. Service validates credentials
4. Repository checks database
5. Service creates session
6. Handler returns token
7. Frontend stores token
8. Token used in subsequent requests
```

### Post Creation Flow
```
1. User creates post
2. Handler validates token
3. Service validates post data
4. Repository saves post
5. Handler returns created post
6. Frontend updates UI
```

### Real-time Message Flow
```
1. User sends message
2. Handler validates and saves
3. Service creates message
4. Repository saves to database
5. WebSocket hub broadcasts
6. Connected clients receive
7. Frontend updates UI
```

## Testing Structure

```
backend/
├── internal/
│   ├── repository/
│   │   └── *_test.go       # Repository tests
│   └── service/
│       └── *_test.go       # Service tests
```

### Test Patterns
- Table-driven tests
- In-memory database for isolation
- Test fixtures and helpers
- Mocking external dependencies

## Build Artifacts

### Development
```
backend/
└── data/
    └── forum.db           # Development database

frontend/
(No build artifacts, vanilla JS)
```

### Production
```
backend/
├── bin/
│   └── forum-backend      # Compiled binary
└── data/
    └── forum.db           # Production database

Docker:
- forum-backend:latest     # Docker image
- forum-data               # Docker volume
```

## Documentation Files

| File | Purpose |
|------|---------|
| `README.md` | Project overview and getting started |
| `API.md` | Complete API documentation |
| `CONTRIBUTING.md` | Contribution guidelines |
| `DEPLOYMENT.md` | Deployment instructions |
| `CHANGELOG.md` | Version history |
| `PROJECT_STRUCTURE.md` | This file |
| `LICENSE` | MIT License |
| `frontend/README.md` | Frontend-specific docs |

## Dependencies

### Backend (Go)
- `gorilla/mux` - HTTP router
- `gorilla/websocket` - WebSocket support
- `mattn/go-sqlite3` - SQLite driver
- `golang.org/x/crypto` - Cryptography (bcrypt)

### Frontend
- No external dependencies (vanilla JS)
- Modern browser APIs only

## Development Workflow

```
Development
    ├── Write code
    ├── Run tests (make test)
    ├── Format code (make fmt)
    ├── Lint code (make lint)
    └── Build (make build)

Deployment
    ├── Build Docker image
    ├── Run tests
    ├── Deploy to environment
    └── Health check
```

## Notes

### Code Organization Principles
1. **Separation of Concerns**: Each layer has a single responsibility
2. **Dependency Injection**: Dependencies passed explicitly
3. **Interface Segregation**: Small, focused interfaces
4. **Single Responsibility**: Functions do one thing well

### Naming Conventions
- **Handlers**: `*Handler` suffix
- **Services**: `*Service` suffix
- **Repositories**: `*Repository` suffix
- **Tests**: `*_test.go` suffix
- **Interfaces**: No special suffix

### Package Dependencies
```
cmd/server
    ↓
internal/api/router
    ↓
internal/api/handlers
    ↓
internal/service
    ↓
internal/repository
    ↓
internal/domain
```

No circular dependencies allowed.

