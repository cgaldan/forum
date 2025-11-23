# Contributing to Real-Time Forum

Thank you for considering contributing to this project! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and considerate in all interactions. We aim to foster an open and welcoming environment.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- A clear and descriptive title
- Steps to reproduce the problem
- Expected behavior vs actual behavior
- Your environment (OS, Go version, browser, etc.)
- Any relevant logs or screenshots

### Suggesting Enhancements

Enhancement suggestions are welcome! Please create an issue with:
- A clear and descriptive title
- Detailed description of the proposed feature
- Rationale for why this enhancement would be useful
- Any implementation ideas you have

### Pull Requests

1. **Fork the repository** and create your branch from `main`

```bash
git checkout -b feature/amazing-feature
```

2. **Make your changes** following the coding standards

3. **Add tests** for your changes (if applicable)

4. **Ensure tests pass**

```bash
cd backend
make test
```

5. **Format your code**

```bash
make fmt
make lint
```

6. **Commit your changes** with a descriptive commit message

```bash
git commit -m "Add amazing feature"
```

7. **Push to your fork**

```bash
git push origin feature/amazing-feature
```

8. **Create a Pull Request** on GitHub

## Development Setup

### Backend Development

1. Install Go 1.21 or higher

2. Clone the repository

```bash
git clone <repository-url>
cd real-time-forum
```

3. Set up environment

```bash
cd backend
cp .env.example .env
```

4. Install dependencies

```bash
make deps
```

5. Run the application

```bash
make run
```

### Running Tests

```bash
make test
```

### Code Formatting

```bash
make fmt
```

### Linting

```bash
make lint
```

## Coding Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Write clear, self-documenting code
- Add comments for exported functions and types
- Keep functions small and focused
- Use meaningful variable names
- Handle errors properly

#### Example

```go
// UserService handles user-related business logic
type UserService struct {
    repo   *UserRepository
    logger *logger.Logger
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id int) (*User, error) {
    if id <= 0 {
        return nil, fmt.Errorf("invalid user ID: %d", id)
    }

    user, err := s.repo.GetByID(id)
    if err != nil {
        s.logger.Error("Failed to get user", "id", id, "error", err)
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    return user, nil
}
```

### JavaScript Code

- Use ES6+ features
- Write clear, self-documenting code
- Add comments for complex logic
- Use meaningful variable and function names
- Keep functions small and focused
- Handle errors gracefully

#### Example

```javascript
/**
 * Loads posts from the API
 * @param {string} category - Optional category filter
 * @returns {Promise<void>}
 */
async function loadPosts(category = '') {
    try {
        const url = category 
            ? `${API_URL}/posts?category=${category}`
            : `${API_URL}/posts`;
        
        const response = await fetch(url);
        const data = await response.json();
        
        if (data.success) {
            state.posts = data.posts || [];
            renderPosts();
        } else {
            throw new Error(data.message || 'Failed to load posts');
        }
    } catch (error) {
        console.error('Load posts error:', error);
        showToast('Failed to load posts', 'error');
    }
}
```

## Project Structure

### Backend

```
backend/
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── api/            # API layer (handlers, middleware, router)
│   ├── config/         # Configuration management
│   ├── domain/         # Domain models and DTOs
│   ├── repository/     # Data access layer
│   ├── service/        # Business logic layer
│   └── websocket/      # WebSocket hub and client
└── pkg/                # Public packages
    └── logger/         # Logging package
```

### Adding New Features

#### Adding a New API Endpoint

1. Define models in `internal/domain/`
2. Add repository methods in `internal/repository/`
3. Add service methods in `internal/service/`
4. Add handler in `internal/api/handlers/`
5. Register route in `internal/api/router/router.go`
6. Write tests for each layer
7. Update API documentation

#### Example: Adding a "Like Post" Feature

1. **Domain Model** (`internal/domain/models.go`)

```go
type PostLike struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
}
```

2. **Repository** (`internal/repository/like_repository.go`)

```go
func (r *LikeRepository) Create(postID, userID int) error {
    _, err := r.db.Exec(`
        INSERT INTO post_likes (post_id, user_id)
        VALUES (?, ?)`, postID, userID)
    return err
}
```

3. **Service** (`internal/service/post_service.go`)

```go
func (s *PostService) LikePost(postID, userID int) error {
    return s.likeRepo.Create(postID, userID)
}
```

4. **Handler** (`internal/api/handlers/post_handler.go`)

```go
func (h *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

5. **Router** (`internal/api/router/router.go`)

```go
api.HandleFunc("/posts/{id}/like", postHandler.LikePost).Methods("POST")
```

## Testing Guidelines

### Writing Tests

- Write tests for all new features
- Aim for high test coverage
- Test both success and error cases
- Use table-driven tests where appropriate
- Mock external dependencies

### Example Test

```go
func TestUserService_GetUserByID(t *testing.T) {
    tests := []struct {
        name    string
        userID  int
        want    *User
        wantErr bool
    }{
        {
            name:    "valid user",
            userID:  1,
            want:    &User{ID: 1, Nickname: "test"},
            wantErr: false,
        },
        {
            name:    "invalid user ID",
            userID:  -1,
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := setupTestService(t)
            got, err := service.GetUserByID(tt.userID)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetUserByID() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Documentation

- Update README.md for major changes
- Update API.md for API changes
- Add inline comments for complex code
- Write clear commit messages
- Update CHANGELOG.md

## Commit Message Format

Use clear, descriptive commit messages:

```
type: Brief description

Longer description if needed

Fixes #123
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Formatting, missing semicolons, etc.
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance tasks

**Examples:**
```
feat: Add post like functionality

Implements the ability for users to like posts.
Includes database schema, API endpoints, and frontend UI.

Closes #45
```

```
fix: Prevent duplicate likes on posts

Users were able to like the same post multiple times.
Added unique constraint on (post_id, user_id).

Fixes #67
```

## Review Process

All pull requests will be reviewed by maintainers. We look for:

1. **Code Quality**
   - Follows coding standards
   - Well-structured and maintainable
   - Properly tested

2. **Documentation**
   - Clear commit messages
   - Updated documentation
   - Code comments where needed

3. **Testing**
   - Tests pass
   - Good test coverage
   - Tests are meaningful

4. **Functionality**
   - Works as intended
   - No breaking changes (unless discussed)
   - Handles edge cases

## Questions?

If you have questions about contributing, feel free to:
- Open an issue with the `question` label
- Reach out to maintainers

Thank you for contributing! 🎉

