package domain

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}

// Post represents a forum post
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    string    `json:"author"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Author    string    `json:"author"`
}

// Message represents a private message
type Message struct {
	ID         int        `json:"id"`
	SenderID   int        `json:"sender_id"`
	ReceiverID int        `json:"receiver_id"`
	Content    string     `json:"content"`
	CreatedAt  time.Time  `json:"created_at"`
	ReadAt     *time.Time `json:"read_at,omitempty"`
	SenderName string     `json:"sender_name"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Conversation represents a messaging conversation
type Conversation struct {
	UserID      int       `json:"user_id"`
	Nickname    string    `json:"nickname"`
	LastMessage string    `json:"last_message"`
	LastTime    time.Time `json:"last_time"`
	UnreadCount int       `json:"unread_count"`
}

// PostDetail includes a post with its comments
type PostDetail struct {
	Post     *Post     `json:"post"`
	Comments []Comment `json:"comments"`
}

// UserStatus represents online/offline status
type UserStatus struct {
	UserID   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	Online   bool   `json:"online"`
}

