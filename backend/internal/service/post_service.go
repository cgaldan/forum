package service

import (
	"fmt"
	"real-time-forum/internal/domain"
	"real-time-forum/internal/repository"
	"real-time-forum/packages/logger"
)

type PostService struct {
	postRepo    repository.PostRepositoryInterface
	commentRepo repository.CommentRepositoryInterface
	logger      *logger.Logger
}

func NewPostService(postRepo repository.PostRepositoryInterface, commentRepo repository.CommentRepositoryInterface, logger *logger.Logger) *PostService {
	return &PostService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		logger:      logger,
	}
}

func (s *PostService) CreatePost(userID int, postData domain.CreatePostRequest) (*domain.Post, error) {
	if err := s.validatePost(postData); err != nil {
		return nil, err
	}

	postID, err := s.postRepo.CreatePost(userID, postData.Title, postData.Content, postData.Category)
	if err != nil {
		s.logger.Error("Failed to create post", "error", err)
		return nil, fmt.Errorf("failed to create post")
	}

	post, err := s.postRepo.GetPostByID(int(postID))
	if err != nil {
		s.logger.Error("Failed to retrieve created post", "error", err)
		return nil, fmt.Errorf("failed to retrieve created post")
	}

	s.logger.Info("Post created successfully", "postID", postID, "userID", userID)
	return post, nil
}

func (s *PostService) GetPostByID(postID int) (*domain.Post, []domain.Comment, error) {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		s.logger.Error("Failed to get post by ID", "error", err, "postID", postID)
		return nil, nil, fmt.Errorf("failed to get post")
	}

	comments, err := s.commentRepo.GetCommentsByPostID(postID)
	if err != nil {
		s.logger.Error("Failed to get comments for post", "error", err, "postID", postID)
		return post, nil, fmt.Errorf("failed to get comments for post")
	}

	return post, comments, nil
}

func (s *PostService) ListPosts(category string, limit, offset int) ([]domain.Post, error) {
	limit, offset = s.validateLimitAndOffset(limit, offset)

	posts, err := s.postRepo.ListPosts(category, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list posts", "error", err, "category", category, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("failed to list posts")
	}

	return posts, nil
}

func (s *PostService) GetPostsByUserID(userID int, limit, offset int) ([]domain.Post, error) {
	limit, offset = s.validateLimitAndOffset(limit, offset)

	posts, err := s.postRepo.GetPostsByUserID(userID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get posts by user ID", "error", err, "userID", userID, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("failed to get posts by user ID")
	}

	return posts, nil
}

func (s *PostService) validatePost(data domain.CreatePostRequest) error {
	if data.Title == "" || len(data.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters")
	}
	if data.Content == "" || len(data.Content) < 10 {
		return fmt.Errorf("content must be at least 10 characters")
	}
	if data.Category == "" {
		return fmt.Errorf("category is required")
	}
	return nil
}

func (s *PostService) validateLimitAndOffset(limit, offset int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
