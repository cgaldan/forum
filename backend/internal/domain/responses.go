package domain

// APIResponse is a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
	Token   string `json:"token,omitempty"`
}

// PostsResponse represents a response containing posts
type PostsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Posts   []Post `json:"posts,omitempty"`
}

// PostResponse represents a response containing a single post
type PostResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Post    *Post  `json:"post,omitempty"`
}

// PostDetailResponse represents a response containing post with comments
type PostDetailResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
	Post     *Post     `json:"post,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
}

// CommentResponse represents a response containing a comment
type CommentResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Comment *Comment `json:"comment,omitempty"`
}

// MessagesResponse represents a response containing messages
type MessagesResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

// MessageResponse represents a response containing a single message
type MessageResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Msg     *Message `json:"msg,omitempty"`
}

// ConversationsResponse represents a response containing conversations
type ConversationsResponse struct {
	Success       bool           `json:"success"`
	Message       string         `json:"message,omitempty"`
	Conversations []Conversation `json:"conversations,omitempty"`
}

// OnlineUsersResponse represents a response containing online users
type OnlineUsersResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Users   []UserStatus `json:"users,omitempty"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

