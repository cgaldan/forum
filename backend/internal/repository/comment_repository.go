package repository

import (
	"database/sql"
	"fmt"
	"real-time-forum/internal/domain"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(postID, userID int, content string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO comments (post_id, user_id, content)
		VALUES (?, ?, ?)`, postID, userID, content)

	if err != nil {
		return 0, fmt.Errorf("failed to create comment: %w", err)
	}

	return result.LastInsertId()
}

func (r *CommentRepository) GetCommentsByPostID(postID int) ([]domain.Comment, error) {
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

func (r *CommentRepository) GetCommentByID(commentID int) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.QueryRow(`
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.nickname
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?`, commentID).Scan(
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

func (r *CommentRepository) GetCommentsByUserID(userID int, limit, offset int) ([]domain.Comment, error) {
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
