package repository

import "database/sql"

// Repositories holds all repository instances
type Repositories struct {
	User    *UserRepository
	Post    *PostRepository
	Comment *CommentRepository
	Message *MessageRepository
	Session *SessionRepository
}

// NewRepositories creates a new Repositories instance
func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:    NewUserRepository(db),
		Post:    NewPostRepository(db),
		Comment: NewCommentRepository(db),
		Message: NewMessageRepository(db),
		Session: NewSessionRepository(db),
	}
}

