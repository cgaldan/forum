package handlers

import (
	"encoding/json"
	"net/http"
	"real-time-forum/internal/domain"
	"real-time-forum/internal/service"
	"real-time-forum/packages/logger"
	"strconv"
)

type CommentHandler struct {
	commentService service.CommentServiceInterface
	authService    service.AuthServiceInterface
	logger         *logger.Logger
}

func NewCommentHandler(commentService service.CommentServiceInterface, authService service.AuthServiceInterface, logger *logger.Logger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		authService:    authService,
		logger:         logger,
	}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	user, err := h.authService.ValidateSession(tokenString)
	if err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idStr := r.PathValue("id")
	postID, err := strconv.Atoi(idStr)

	// WITH GORILLA PKG IMPLEMANTATION
	// vars := mux.Vars(r)
	// postID, err := strconv.Atoi(vars["id"])

	if err != nil || postID <= 0 {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var commentRequest domain.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&commentRequest); err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	comment, err := h.commentService.CreateComment(user.ID, postID, commentRequest)
	if err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.CommentResponse{
		Success: true,
		Message: "Comment created successfully",
		Comment: comment,
	})
}
