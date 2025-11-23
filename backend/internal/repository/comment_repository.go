package repository

import (
	"database/sql"
	"fmt"

	"forum-backend/internal/domain"
)

// CommentRepository handles comment data access
type CommentRepository struct {
	db *sql.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create creates a new comment
func (r *CommentRepository) Create(postID, userID int, content string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO comments (post_id, user_id, content)
		VALUES (?, ?, ?)`, postID, userID, content)
	
	if err != nil {
		return 0, fmt.Errorf("failed to create comment: %w", err)
	}

	return result.LastInsertId()
}

// GetByPostID gets all comments for a post
func (r *CommentRepository) GetByPostID(postID int) ([]domain.Comment, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.nickname
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC`, postID)

	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content,
			&comment.CreatedAt, &comment.Author); err != nil {
			continue
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// GetByID gets a comment by ID
func (r *CommentRepository) GetByID(id int) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.QueryRow(`
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.nickname
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?`, id).Scan(
		&comment.ID, &comment.PostID, &comment.UserID, &comment.Content,
		&comment.CreatedAt, &comment.Author)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("comment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return &comment, nil
}

// GetByUserID gets comments by a specific user
func (r *CommentRepository) GetByUserID(userID, limit, offset int) ([]domain.Comment, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.nickname
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.user_id = ?
		ORDER BY c.created_at DESC LIMIT ? OFFSET ?`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get user comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content,
			&comment.CreatedAt, &comment.Author); err != nil {
			continue
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

