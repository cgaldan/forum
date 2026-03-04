package handlers

import (
	"encoding/json"
	"net/http"
	"real-time-forum/internal/domain"
	"real-time-forum/internal/service"
	"real-time-forum/packages/logger"
	"strconv"
)

type PostHandler struct {
	postService service.PostServiceInterface
	authService service.AuthServiceInterface
	logger      *logger.Logger
}

func NewPostHandler(postService service.PostServiceInterface, authService service.AuthServiceInterface, logger *logger.Logger) *PostHandler {
	return &PostHandler{
		postService: postService,
		authService: authService,
		logger:      logger,
	}
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	category := r.URL.Query().Get("category")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	if limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 && limitNum <= 100 {
			limit = limitNum
		}
	}

	offset := 0
	if offsetStr != "" {
		if offsetNum, err := strconv.Atoi(offsetStr); err == nil && offsetNum >= 0 {
			offset = offsetNum
		}
	}

	posts, err := h.postService.ListPosts(category, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list posts", "error", err)
		json.NewEncoder(w).Encode(domain.PostsResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.PostsResponse{
		Success: true,
		Message: "Posts retrieved successfully",
		Posts:   posts,
	})
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	postID, err := strconv.Atoi(idStr)

	// WITH GORILLA PKG IMPLEMENTATION
	// vars := mux.Vars(r)
	// postID, err := strconv.Atoi(vars["id"])

	if err != nil || postID <= 0 {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	post, comments, err := h.postService.GetPostByID(postID)
	if err != nil {
		h.logger.Error("Failed to get post by ID", "error", err, "postID", postID)
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.PostDetailResponse{
		Success:  true,
		Message:  "Post retrieved successfully",
		Post:     post,
		Comments: comments,
	})
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var postData domain.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	post, err := h.postService.CreatePost(user.ID, postData)
	if err != nil {
		h.logger.Error("Failed to create post", "error", err, "userID", user.ID)
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.PostDetailResponse{
		Success: true,
		Message: "Post created successfully",
		Post:    post,
	})
}
