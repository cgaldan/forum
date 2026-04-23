# Project Structure

This document provides a detailed overview of the project's file and directory structure.

## Root Directory

```
real-time-forum/
├── .github/                  # GitHub-specific files
│   └── workflows/            # CI/CD workflows
│       └── ci.yml            # Continuous Integration pipeline
├── backend/                  # Backend Go application
├── frontend/                 # Frontend application
├── .gitignore                # Git ignore rules
├── API.md                    # API documentation
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
├── cmd/                     # Application entry points
│   └── server/              # Main server application
│       └── main.go          # Entry point with initialization
├── data/
│   └── database/             # Database entry point
│       └── forum.db          # Database file after project initialization
│
├── internal/               # Private application code
│   ├── api/                # API layer
│   │   ├── handlers/       # HTTP request handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── comment_handler.go
│   │   │   ├── health_handler.go
│   │   │   ├── message_handler.go
│   │   │   ├── post_handler.go
│   │   │   └── websocket_handler.go
│   │   │
│   │   ├── middleware/     # HTTP middlewares
│   │   │   ├── CORS.go
│   │   │   ├── logging.go
│   │   │   ├── middleware_helpers.go
│   │   │   ├── rate_limiter.go
│   │   │   ├── recovery.go
│   │   │   └── security_headers.go
│   │   │
│   │   └── router/         # Routing configuration
│   │       ├── router_without_gorilla.go       # Route definitions without using gorilla mux package
│   │       ├── router.go                       # Route definitions

│   │
│   ├── config/             # Configuration management
│   │   ├── config_models.go           # All configuration models
│   │   ├── config_utils.go            # Configuration helper functions
│   │   └── config.go                  # Configuration loading and validation
│   │
│   ├── database/             # Domain models and DTOs
│   │   └── database.go                # Database initialization functions and migrations
│   │
│   ├── domain/             # Domain models and DTOs
│   │   ├── models.go                  # Core domain models
│   │   ├── requests.go                # Request DTOs
│   │   └── responses.go               # Response DTOs
│   │
│   ├── repository/         # Data access layer
│   │   ├── comment_repository.go      # Comment data access
│   │   ├── comment_repository_test.go # Comment repository tests
│   │   ├── message_repository_test.go # Message repository tests
│   │   ├── message_repository.go      # Message data access
│   │   ├── post_repository_test.go    # Post repository tests
│   │   ├── post_repository.go         # Post data access
│   │   ├── repositories.go            # Repository factory
│   │   ├── session_repository_test.go # Session repository tests
│   │   ├── session_repository.go      # Session data access
│   │   ├── test_utils.go              # Tests helper functions
│   │   ├── user_repository_test.go    # User repository tests
│   │   └── user_repository.go         # User data access
│   │
│   ├── service/            # Business logic layer
│   │   ├── auth_service_test.go       # Authentication tests
│   │   ├── auth_service.go            # Authentication logic
│   │   ├── comment_service_test.go    # Comment tests
│   │   ├── comment_service.go         # Comment logic
│   │   ├── message_service_test.go    # Messaging tests
│   │   ├── message_service.go         # Messaging logic
│   │   ├── post_service_test.go       # Post tests
│   │   ├── post_service.go            # Post logic
│   │   ├── service_test_helpers.go    # Service test helper functions
│   │   └── services.go                # Service factory
│   │
│   └── websocket/          # WebSocket management
│       ├── client.go                  # WebSocket client
│       ├── hub.go                     # WebSocket hub
│       └── ws_utils.go                     # WebSocket helper functions
│
├── packages/                    # Public packages
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
├── js/                   # JavaScript modules
│   ├── config.js         # Configuration settings
│   ├── main.js           # Main application entry point
│   │   ├── Configuration
│   │   ├── State management
│   │   ├── Initialization
│   │   ├── Authentication handlers
│   │   ├── Post management
│   │   ├── Comment handlers
│   │   ├── Message handlers
│   │   ├── WebSocket handlers
│   │   └── Utility functions
│   │
│   ├── api/
│   │   └── client.js     # API client for HTTP requests
│   │
│   ├── modules/          # Feature-specific modules
│   │   ├── auth.js       # Authentication module
│   │   ├── messages.js   # Messaging module
│   │   ├── posts.js      # Posts module
│   │   └── websocket.js  # WebSocket module
│   │
│   ├── state/
│   │   └── store.js      # State management store
│   │
│   ├── ui/               # UI-related modules
│   │   ├── events.js     # Event handlers
│   │   └── ui.js         # UI manipulation functions
│   │
│   └── utils/
│       └── helpers.js    # Utility helper functions
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

## Testing Structure

```
backend/
├── internal/
│   ├── repository/
│   │   └── *_test.go       # Repository tests
│   └── service/
│       └── *_test.go       # Service tests
```

## Documentation Files

| File | Purpose |
|------|---------|
| [`README.md`](/README.md) | Project overview and getting started |
| [`DEPLOYMENT.md`](/DEPLOYMENT.md) | Deployment instructions |
| [`PROJECT_STRUCTURE.md`](/PROJECT_STRUCTURE.md) | This file |
| [`LICENSE`](/LICENSE) | MIT License |
