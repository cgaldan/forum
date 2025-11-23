package domain

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Identifier string `json:"identifier"` // nickname or email
	Password   string `json:"password"`
}

// CreatePostRequest represents a request to create a post
type CreatePostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

// CreateCommentRequest represents a request to create a comment
type CreateCommentRequest struct {
	Content string `json:"content"`
}

// SendMessageRequest represents a request to send a message
type SendMessageRequest struct {
	Content string `json:"content"`
}

