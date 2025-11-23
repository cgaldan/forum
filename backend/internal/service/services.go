package service

import (
	"forum-backend/internal/repository"
	"forum-backend/internal/websocket"
	"forum-backend/pkg/logger"
)

// Services holds all service instances
type Services struct {
	Auth    *AuthService
	Post    *PostService
	Comment *CommentService
	Message *MessageService
}

// NewServices creates a new Services instance
func NewServices(repos *repository.Repositories, hub *websocket.Hub, logger *logger.Logger) *Services {
	return &Services{
		Auth:    NewAuthService(repos.User, repos.Session, logger),
		Post:    NewPostService(repos.Post, repos.Comment, logger),
		Comment: NewCommentService(repos.Comment, repos.Post, logger),
		Message: NewMessageService(repos.Message, repos.User, hub, logger),
	}
}

