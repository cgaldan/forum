package repository

import (
	"database/sql"
	"fmt"
	"time"

	"forum-backend/internal/domain"
)

// UserRepository handles user data access
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(nickname, email, passwordHash, firstName, lastName string, age int, gender string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO users (nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		nickname, email, passwordHash, firstName, lastName, age, gender)
	
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return result.LastInsertId()
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(id int) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(`
		SELECT id, nickname, email, first_name, last_name, age, gender, created_at, last_seen
		FROM users WHERE id = ?`, id).Scan(
		&user.ID, &user.Nickname, &user.Email, &user.FirstName, &user.LastName,
		&user.Age, &user.Gender, &user.CreatedAt, &user.LastSeen)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByIdentifier gets a user by nickname or email
func (r *UserRepository) GetByIdentifier(identifier string) (*domain.User, string, error) {
	var user domain.User
	var passwordHash string
	
	err := r.db.QueryRow(`
		SELECT id, nickname, email, password_hash, first_name, last_name, age, gender, created_at, last_seen
		FROM users WHERE nickname = ? OR email = ?`, identifier, identifier).Scan(
		&user.ID, &user.Nickname, &user.Email, &passwordHash, &user.FirstName, &user.LastName,
		&user.Age, &user.Gender, &user.CreatedAt, &user.LastSeen)

	if err == sql.ErrNoRows {
		return nil, "", fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	return &user, passwordHash, nil
}

// UpdateLastSeen updates user's last seen timestamp
func (r *UserRepository) UpdateLastSeen(userID int) error {
	_, err := r.db.Exec("UPDATE users SET last_seen = CURRENT_TIMESTAMP WHERE id = ?", userID)
	return err
}

// GetByNickname gets a user by nickname
func (r *UserRepository) GetByNickname(nickname string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(`
		SELECT id, nickname, email, first_name, last_name, age, gender, created_at, last_seen
		FROM users WHERE nickname = ?`, nickname).Scan(
		&user.ID, &user.Nickname, &user.Email, &user.FirstName, &user.LastName,
		&user.Age, &user.Gender, &user.CreatedAt, &user.LastSeen)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// List returns a paginated list of users
func (r *UserRepository) List(limit, offset int) ([]domain.User, error) {
	rows, err := r.db.Query(`
		SELECT id, nickname, email, first_name, last_name, age, gender, created_at, last_seen
		FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Nickname, &user.Email, &user.FirstName, &user.LastName,
			&user.Age, &user.Gender, &user.CreatedAt, &user.LastSeen); err != nil {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

// GetOnlineUsers returns users who are currently online
func (r *UserRepository) GetOnlineUsers(since time.Duration) ([]domain.User, error) {
	threshold := time.Now().Add(-since)
	
	rows, err := r.db.Query(`
		SELECT id, nickname, email, first_name, last_name, age, gender, created_at, last_seen
		FROM users WHERE last_seen > ? ORDER BY nickname`, threshold)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get online users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Nickname, &user.Email, &user.FirstName, &user.LastName,
			&user.Age, &user.Gender, &user.CreatedAt, &user.LastSeen); err != nil {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

