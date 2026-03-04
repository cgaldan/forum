package service

import (
	"real-time-forum/internal/domain"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/websocket"
	"real-time-forum/packages/logger"
)

type Services struct {
	Auth    AuthServiceInterface
	Post    PostServiceInterface
	Comment CommentServiceInterface
	Message MessageServiceInterface
}

func NewServices(repos *repository.Repositories, hub *websocket.Hub, logger *logger.Logger) *Services {
	return &Services{
		Auth:    NewAuthService(repos.User, repos.Session, logger),
		Post:    NewPostService(repos.Post, repos.Comment, logger),
		Comment: NewCommentService(repos.Comment, repos.Post, logger),
		Message: NewMessageService(repos.Message, repos.User, hub, logger),
	}
}

type AuthServiceInterface interface {
	Register(registrationData domain.RegisterRequest) (*domain.User, string, error)
	Login(loginData domain.LoginRequest) (*domain.User, string, error)
	Logout(sessionID string) error
	ValidateSession(sessionID string) (*domain.User, error)
}

type PostServiceInterface interface {
	CreatePost(userID int, postData domain.CreatePostRequest) (*domain.Post, error)
	GetPostByID(postID int) (*domain.Post, []domain.Comment, error)
	ListPosts(category string, limit, offset int) ([]domain.Post, error)
	GetPostsByUserID(userID int, limit, offset int) ([]domain.Post, error)
}

type CommentServiceInterface interface {
	CreateComment(userID int, postID int, commentData domain.CreateCommentRequest) (*domain.Comment, error)
	GetCommentsByPostID(postID int) ([]domain.Comment, error)
	GetCommentsByUserID(userID, limit, offset int) ([]domain.Comment, error)
}

type MessageServiceInterface interface {
	SendMessage(senderID, receiverID int, content string) (*domain.Message, error)
	GetConversation(userID1, userID2 int, limit, offset int) ([]domain.Message, error)
	GetConversationsByUserID(userID int) ([]domain.Conversation, error)
}
