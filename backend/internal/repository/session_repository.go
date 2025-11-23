package repository

import (
	"database/sql"
	"fmt"
	"time"

	"forum-backend/internal/domain"
)

// SessionRepository handles session data access
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session
func (r *SessionRepository) Create(sessionID string, userID int, expiresAt time.Time) error {
	_, err := r.db.Exec(`
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES (?, ?, ?)`, sessionID, userID, expiresAt)
	
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetByID gets a session by ID
func (r *SessionRepository) GetByID(sessionID string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.QueryRow(`
		SELECT id, user_id, created_at, expires_at
		FROM sessions WHERE id = ? AND expires_at > CURRENT_TIMESTAMP`, sessionID).Scan(
		&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

// Delete deletes a session
func (r *SessionRepository) Delete(sessionID string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// DeleteExpired deletes all expired sessions
func (r *SessionRepository) DeleteExpired() error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP")
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}
	return nil
}

// DeleteByUserID deletes all sessions for a user
func (r *SessionRepository) DeleteByUserID(userID int) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}
	return nil
}

