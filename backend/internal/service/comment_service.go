package service

import (
	"fmt"
	"real-time-forum/internal/domain"
	"real-time-forum/internal/repository"
	"real-time-forum/packages/logger"
)

type CommentService struct {
	commentRepo repository.CommentRepositoryInterface
	postRepo    repository.PostRepositoryInterface
	logger      *logger.Logger
}

func NewCommentService(commentRepo repository.CommentRepositoryInterface, postRepo repository.PostRepositoryInterface, logger *logger.Logger) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		logger:      logger,
	}
}

func (s *CommentService) CreateComment(userID int, postID int, commentData domain.CreateCommentRequest) (*domain.Comment, error) {
	if err := s.validateComment(commentData); err != nil {
		return nil, err
	}

	exists, err := s.postRepo.PostExists(postID)
	if err != nil || !exists {
		s.logger.Error("Failed to check if post exists", "error", err, "postID", postID)
		return nil, fmt.Errorf("failed to create comment")
	}

	commentID, err := s.commentRepo.CreateComment(userID, postID, commentData.Content)
	if err != nil {
		s.logger.Error("Failed to create comment", "error", err)
		return nil, fmt.Errorf("failed to create comment")
	}

	comment, err := s.commentRepo.GetCommentByID(int(commentID))
	if err != nil {
		s.logger.Error("Failed to retrieve created comment", "error", err, "commentID", commentID)
		return nil, fmt.Errorf("failed to retrieve created comment")
	}

	s.logger.Info("Comment created successfully", "commentID", commentID, "userID", userID, "postID", postID)
	return comment, nil
}

func (s *CommentService) GetCommentsByPostID(postID int) ([]domain.Comment, error) {
	comments, err := s.commentRepo.GetCommentsByPostID(postID)
	if err != nil {
		s.logger.Error("Failed to get comments by post ID", "error", err, "postID", postID)
		return nil, fmt.Errorf("failed to get comments for post")
	}

	return comments, nil
}

func (s *CommentService) GetCommentsByUserID(userID, limit, offset int) ([]domain.Comment, error) {
	comments, err := s.commentRepo.GetCommentsByUserID(userID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get comments by user ID", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to get comments for user")
	}

	return comments, nil
}

func (s *CommentService) validateComment(commentData domain.CreateCommentRequest) error {
	if commentData.Content == "" || len(commentData.Content) < 1 {
		return fmt.Errorf("comment cannot be empty")
	}
	return nil
}
