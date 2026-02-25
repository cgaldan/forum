package repository

import (
	"database/sql"
	"real-time-forum/internal/domain"
	"time"
)

type Repositories struct {
	User    UserRepositoryInterface
	Post    PostRepositoryInterface
	Comment CommentRepositoryInterface
	Message MessageRepositoryInterface
	Session SessionRepositoryInterface
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:    NewUserRepository(db),
		Post:    NewPostRepository(db),
		Comment: NewCommentRepository(db),
		Message: NewMessageRepository(db),
		Session: NewSessionRepository(db),
	}
}

type UserRepositoryInterface interface {
	CreateUser(nickname, email, passwordHash, firstName, lastName string, age int, gender string) (int64, error)
	GetUserByID(userID int) (*domain.User, error)
	GetUserByIdentifier(identifier string) (*domain.User, string, error)
	UpdateLastSeen(userID int) error
}

type PostRepositoryInterface interface {
	CreatePost(userID int, title, content, category string) (int64, error)
	GetPostByID(postID int) (*domain.Post, error)
	ListPosts(category string, limit, offset int) ([]domain.Post, error)
	GetPostsByUserID(userID int, limit, offset int) ([]domain.Post, error)
	PostExists(postID int) (bool, error)
}

type CommentRepositoryInterface interface {
	CreateComment(postID, userID int, content string) (int64, error)
	GetCommentsByPostID(postID int) ([]domain.Comment, error)
	GetCommentByID(commentID int) (*domain.Comment, error)
	GetCommentsByUserID(userID int, limit, offset int) ([]domain.Comment, error)
}

type SessionRepositoryInterface interface {
	CreateSession(sessionID string, userID int, expiresAt time.Time) error
	GetSessionBySessionID(sessionID string) (*domain.Session, error)
	DeleteSession(sessionID string) error
}

type MessageRepositoryInterface interface {
	CreateMessage(senderID, receiverID int, content string) (int64, error)
	GetConversation(userID1, userID2, limit, offset int) ([]domain.Message, error)
	MarkMessageAsRead(receiverID, senderID int) error
	GetConversationPartners(userID int) ([]int, error)
	GetUnreadMessagesCount(receiverID, senderID int) (int, error)
	GetLastMessageBetweenUsers(user1ID, user2ID int) (*domain.Message, error)
}
