package service

import (
	"fmt"

	"forum-backend/internal/domain"
	"forum-backend/internal/repository"
	"forum-backend/pkg/logger"
)

// CommentService handles comment business logic
type CommentService struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
	logger      *logger.Logger
}

// NewCommentService creates a new comment service
func NewCommentService(commentRepo *repository.CommentRepository, postRepo *repository.PostRepository, logger *logger.Logger) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		logger:      logger,
	}
}

// CreateComment creates a new comment
func (s *CommentService) CreateComment(postID, userID int, req domain.CreateCommentRequest) (*domain.Comment, error) {
	// Validate
	if req.Content == "" || len(req.Content) < 1 {
		return nil, fmt.Errorf("comment cannot be empty")
	}

	// Check if post exists
	exists, err := s.postRepo.Exists(postID)
	if err != nil || !exists {
		return nil, fmt.Errorf("post not found")
	}

	// Create comment
	commentID, err := s.commentRepo.Create(postID, userID, req.Content)
	if err != nil {
		s.logger.Error("Failed to create comment", "error", err, "postID", postID, "userID", userID)
		return nil, fmt.Errorf("failed to create comment")
	}

	// Get created comment
	comment, err := s.commentRepo.GetByID(int(commentID))
	if err != nil {
		s.logger.Error("Failed to get created comment", "error", err, "commentID", commentID)
		return nil, fmt.Errorf("comment created but failed to retrieve")
	}

	s.logger.Info("Comment created successfully", "commentID", commentID, "postID", postID, "userID", userID)
	return comment, nil
}

// GetCommentsByPost gets all comments for a post
func (s *CommentService) GetCommentsByPost(postID int) ([]domain.Comment, error) {
	comments, err := s.commentRepo.GetByPostID(postID)
	if err != nil {
		s.logger.Error("Failed to get comments", "error", err, "postID", postID)
		return nil, fmt.Errorf("failed to fetch comments")
	}

	return comments, nil
}

// GetUserComments gets comments by a specific user
func (s *CommentService) GetUserComments(userID, limit, offset int) ([]domain.Comment, error) {
	comments, err := s.commentRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get user comments", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to fetch user comments")
	}

	return comments, nil
}

