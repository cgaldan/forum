package domain

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
	Token   string `json:"token,omitempty"`
}

type PostsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Posts   []Post `json:"posts,omitempty"`
}

type PostResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Post    *Post  `json:"post,omitempty"`
}

type PostDetailResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
	Post     *Post     `json:"post,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
}

type CommentResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Comment *Comment `json:"comment,omitempty"`
}

type MessagesResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type MessageResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Msg     *Message `json:"msg,omitempty"`
}

type ConversationsResponse struct {
	Success       bool           `json:"success"`
	Message       string         `json:"message,omitempty"`
	Conversations []Conversation `json:"conversations,omitempty"`
}

type OnlineUsersResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Users   []UserStatus `json:"users,omitempty"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}
