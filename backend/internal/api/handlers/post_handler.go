package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum-backend/internal/domain"
	"forum-backend/internal/service"
	"forum-backend/pkg/logger"

	"github.com/gorilla/mux"
)

// PostHandler handles post endpoints
type PostHandler struct {
	postService *service.PostService
	authService *service.AuthService
	logger      *logger.Logger
}

// NewPostHandler creates a new post handler
func NewPostHandler(postService *service.PostService, authService *service.AuthService, logger *logger.Logger) *PostHandler {
	return &PostHandler{
		postService: postService,
		authService: authService,
		logger:      logger,
	}
}

// GetPosts handles listing posts
func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters
	category := r.URL.Query().Get("category")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	posts, err := h.postService.ListPosts(category, limit, offset)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostsResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(domain.PostsResponse{Success: true, Posts: posts})
}

// GetPost handles getting a single post
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{Success: false, Message: "Invalid post ID"})
		return
	}

	post, comments, err := h.postService.GetPost(postID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(domain.PostDetailResponse{
		Success:  true,
		Post:     post,
		Comments: comments,
	})
}

// CreatePost handles creating a new post
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user from session
	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostResponse{Success: false, Message: "Unauthorized"})
		return
	}

	var req domain.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.PostResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	post, err := h.postService.CreatePost(user.ID, req)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(domain.PostResponse{
		Success: true,
		Message: "Post created successfully",
		Post:    post,
	})
}

