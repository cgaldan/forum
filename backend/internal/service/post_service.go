package service

import (
	"fmt"

	"forum-backend/internal/domain"
	"forum-backend/internal/repository"
	"forum-backend/pkg/logger"
)

// PostService handles post business logic
type PostService struct {
	postRepo    *repository.PostRepository
	commentRepo *repository.CommentRepository
	logger      *logger.Logger
}

// NewPostService creates a new post service
func NewPostService(postRepo *repository.PostRepository, commentRepo *repository.CommentRepository, logger *logger.Logger) *PostService {
	return &PostService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		logger:      logger,
	}
}

// CreatePost creates a new post
func (s *PostService) CreatePost(userID int, req domain.CreatePostRequest) (*domain.Post, error) {
	// Validate
	if err := s.validatePost(req); err != nil {
		return nil, err
	}

	// Create post
	postID, err := s.postRepo.Create(userID, req.Title, req.Content, req.Category)
	if err != nil {
		s.logger.Error("Failed to create post", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to create post")
	}

	// Get created post
	post, err := s.postRepo.GetByID(int(postID))
	if err != nil {
		s.logger.Error("Failed to get created post", "error", err, "postID", postID)
		return nil, fmt.Errorf("post created but failed to retrieve")
	}

	s.logger.Info("Post created successfully", "postID", postID, "userID", userID)
	return post, nil
}

// GetPost gets a post by ID with comments
func (s *PostService) GetPost(postID int) (*domain.Post, []domain.Comment, error) {
	// Get post
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, nil, fmt.Errorf("post not found")
	}

	// Get comments
	comments, err := s.commentRepo.GetByPostID(postID)
	if err != nil {
		s.logger.Warn("Failed to get comments for post", "error", err, "postID", postID)
		comments = []domain.Comment{} // Return empty comments on error
	}

	return post, comments, nil
}

// ListPosts lists posts with optional category filter
func (s *PostService) ListPosts(category string, limit, offset int) ([]domain.Post, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := s.postRepo.List(category, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list posts", "error", err)
		return nil, fmt.Errorf("failed to fetch posts")
	}

	return posts, nil
}

// GetUserPosts gets posts by a specific user
func (s *PostService) GetUserPosts(userID, limit, offset int) ([]domain.Post, error) {
	posts, err := s.postRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get user posts", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to fetch user posts")
	}

	return posts, nil
}

func (s *PostService) validatePost(req domain.CreatePostRequest) error {
	if req.Title == "" || len(req.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters")
	}
	if req.Content == "" || len(req.Content) < 10 {
		return fmt.Errorf("content must be at least 10 characters")
	}
	if req.Category == "" {
		return fmt.Errorf("category is required")
	}
	return nil
}

