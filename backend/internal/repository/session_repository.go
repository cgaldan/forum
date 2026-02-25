package repository

import (
	"database/sql"
	"fmt"
	"real-time-forum/internal/domain"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(sessionID string, userID int, expiresAt time.Time) error {
	_, err := r.db.Exec(`
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES (?, ?, ?)`, sessionID, userID, expiresAt)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *SessionRepository) GetSessionBySessionID(sessionID string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.QueryRow(`
		SELECT id, user_id, created_at, expires_at
		FROM sessions WHERE id = ? AND expires_at > CURRENT_TIMESTAMP`, sessionID).Scan(
		&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session by session ID: %w", err)
	}

	return &session, nil
}

func (r *SessionRepository) DeleteSession(sessionID string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
