package repository

import (
	"database/sql"
	"fmt"

	"forum-backend/internal/domain"
)

// PostRepository handles post data access
type PostRepository struct {
	db *sql.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create creates a new post
func (r *PostRepository) Create(userID int, title, content, category string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO posts (user_id, title, content, category)
		VALUES (?, ?, ?, ?)`, userID, title, content, category)
	
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %w", err)
	}

	return result.LastInsertId()
}

// GetByID gets a post by ID
func (r *PostRepository) GetByID(id int) (*domain.Post, error) {
	var post domain.Post
	err := r.db.QueryRow(`
		SELECT p.id, p.user_id, p.title, p.content, p.category, p.created_at, p.updated_at, u.nickname
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?`, id).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category,
		&post.CreatedAt, &post.UpdatedAt, &post.Author)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

// List returns a paginated list of posts
func (r *PostRepository) List(category string, limit, offset int) ([]domain.Post, error) {
	var rows *sql.Rows
	var err error

	if category != "" {
		rows, err = r.db.Query(`
			SELECT p.id, p.user_id, p.title, p.content, p.category, p.created_at, p.updated_at, u.nickname
			FROM posts p
			JOIN users u ON p.user_id = u.id
			WHERE p.category = ?
			ORDER BY p.created_at DESC LIMIT ? OFFSET ?`, category, limit, offset)
	} else {
		rows, err = r.db.Query(`
			SELECT p.id, p.user_id, p.title, p.content, p.category, p.created_at, p.updated_at, u.nickname
			FROM posts p
			JOIN users u ON p.user_id = u.id
			ORDER BY p.created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category,
			&post.CreatedAt, &post.UpdatedAt, &post.Author); err != nil {
			continue
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetByUserID returns posts by a specific user
func (r *PostRepository) GetByUserID(userID, limit, offset int) ([]domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.category, p.created_at, p.updated_at, u.nickname
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC LIMIT ? OFFSET ?`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category,
			&post.CreatedAt, &post.UpdatedAt, &post.Author); err != nil {
			continue
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Exists checks if a post exists
func (r *PostRepository) Exists(id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", id).Scan(&exists)
	return exists, err
}

