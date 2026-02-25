package repository

import (
	"database/sql"
	"fmt"
	"real-time-forum/internal/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(userID int, title, content, category string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO posts (user_id, title, content, category)
		VALUES (?, ?, ?, ?)`, userID, title, content, category)

	if err != nil {
		return 0, fmt.Errorf("failed to create post: %w", err)
	}

	return result.LastInsertId()
}

func (r *PostRepository) GetPostByID(postID int) (*domain.Post, error) {
	var post domain.Post
	err := r.db.QueryRow(`
		SELECT p.id, p.user_id, p.title, p.content, p.category, p.created_at, p.updated_at, u.nickname
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?`, postID).Scan(
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

func (r *PostRepository) ListPosts(category string, limit, offset int) ([]domain.Post, error) {
	var rows *sql.Rows
	var err error

	if category != "" {
		rows, err = r.db.Query(`
			SELECT p.id, p.user_id, p.title, p.content, p.category, p.created_at, p.updated_at, u.nickname
			FROM posts p
			JOIN users u ON p.user_id = u.id
			WHERE p.category = ?
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?`, category, limit, offset)
	} else {
		rows, err = r.db.Query(`
			SELECT p.id, p.user_id, p.title, p.content, p.category, p.created_at, p.updated_at, u.nickname
			FROM posts p
			JOIN users u ON p.user_id = u.id
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?`, limit, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category,
			&post.CreatedAt, &post.UpdatedAt, &post.Author)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) GetPostsByUserID(userID int, limit, offset int) ([]domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.category, p.created_at, p.updated_at, u.nickname
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category,
			&post.CreatedAt, &post.UpdatedAt, &post.Author)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) PostExists(postID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)`, postID).Scan(&exists)
	return exists, err
}
