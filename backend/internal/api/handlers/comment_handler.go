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

// CommentHandler handles comment endpoints
type CommentHandler struct {
	commentService *service.CommentService
	authService    *service.AuthService
	logger         *logger.Logger
}

// NewCommentHandler creates a new comment handler
func NewCommentHandler(commentService *service.CommentService, authService *service.AuthService, logger *logger.Logger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		authService:    authService,
		logger:         logger,
	}
}

// CreateComment handles creating a new comment
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user from session
	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{Success: false, Message: "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{Success: false, Message: "Invalid post ID"})
		return
	}

	var req domain.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	comment, err := h.commentService.CreateComment(postID, user.ID, req)
	if err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(domain.CommentResponse{
		Success: true,
		Message: "Comment created successfully",
		Comment: comment,
	})
}

