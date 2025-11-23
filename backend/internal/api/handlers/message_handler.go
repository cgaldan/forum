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

// MessageHandler handles messaging endpoints
type MessageHandler struct {
	messageService *service.MessageService
	authService    *service.AuthService
	logger         *logger.Logger
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(messageService *service.MessageService, authService *service.AuthService, logger *logger.Logger) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		authService:    authService,
		logger:         logger,
	}
}

// SendMessage handles sending a message
func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get sender from session
	token := r.Header.Get("Authorization")
	sender, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{Success: false, Message: "Unauthorized"})
		return
	}

	// Get receiver ID from URL
	vars := mux.Vars(r)
	receiverID, err := strconv.Atoi(vars["id"])
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{Success: false, Message: "Invalid user ID"})
		return
	}

	var req domain.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	message, err := h.messageService.SendMessage(sender.ID, receiverID, req.Content)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(domain.MessageResponse{
		Success: true,
		Message: "Message sent",
		Msg:     message,
	})
}

// GetMessages handles getting message history
func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get current user from session
	token := r.Header.Get("Authorization")
	currentUser, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessagesResponse{Success: false, Message: "Unauthorized"})
		return
	}

	// Get other user ID from URL
	vars := mux.Vars(r)
	otherUserID, err := strconv.Atoi(vars["id"])
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessagesResponse{Success: false, Message: "Invalid user ID"})
		return
	}

	// Get pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
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

	messages, err := h.messageService.GetConversation(currentUser.ID, otherUserID, limit, offset)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessagesResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(domain.MessagesResponse{
		Success:  true,
		Messages: messages,
	})
}

// GetConversations handles getting all conversations
func (h *MessageHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get current user from session
	token := r.Header.Get("Authorization")
	currentUser, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.ConversationsResponse{Success: false, Message: "Unauthorized"})
		return
	}

	conversations, err := h.messageService.GetConversations(currentUser.ID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.ConversationsResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(domain.ConversationsResponse{
		Success:       true,
		Conversations: conversations,
	})
}

