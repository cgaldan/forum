package repository

import (
	"database/sql"
	"fmt"
	"real-time-forum/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(nickname, email, passwordHash, firstName, lastName string, age int, gender string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO users (nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		nickname, email, passwordHash, firstName, lastName, age, gender)

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return result.LastInsertId()
}

func (r *UserRepository) GetUserByID(userID int) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(`
		SELECT id, nickname, email, first_name, last_name, age, gender, created_at, last_seen
		FROM users WHERE id = ?`, userID).Scan(
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

func (r *UserRepository) GetUserByIdentifier(identifier string) (*domain.User, string, error) {
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

func (r *UserRepository) UpdateLastSeen(userID int) error {
	_, err := r.db.Exec("UPDATE users SET last_seen = CURRENT_TIMESTAMP WHERE id = ?", userID)
	return err
}
