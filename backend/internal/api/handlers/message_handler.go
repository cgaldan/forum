package handlers

import (
	"encoding/json"
	"net/http"
	"path"
	"real-time-forum/internal/domain"
	"real-time-forum/internal/service"
	"real-time-forum/packages/logger"
	"strconv"
)

type MessageHandler struct {
	messageService service.MessageServiceInterface
	authService    service.AuthServiceInterface
	logger         *logger.Logger
}

func NewMessageHandler(messageService service.MessageServiceInterface, authService service.AuthServiceInterface, logger *logger.Logger) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		authService:    authService,
		logger:         logger,
	}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	sender, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idStr := path.Base(r.URL.Path)
	receiverID, err := strconv.Atoi(idStr)

	// WITH GORILLA PKG IMPLEMENTATION
	// vars := mux.Vars(r)
	// receiverID, err := strconv.Atoi(vars["id"])

	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var req domain.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Invalid JSON",
		})
		return
	}

	message, err := h.messageService.SendMessage(sender.ID, receiverID, req.Content)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.MessageResponse{
		Success: true,
		Message: "Message sent",
		Msg:     message,
	})
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	currentUser, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idStr := path.Base(r.URL.Path)
	receiverID, err := strconv.Atoi(idStr)

	// WITH GORILLA PKG IMPLEMENTATION
	// vars := mux.Vars(r)
	// receiverID, err := strconv.Atoi(vars["id"])

	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
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

	messages, err := h.messageService.GetConversation(currentUser.ID, receiverID, limit, offset)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.MessagesResponse{
		Success:  true,
		Message:  "Messages retrieved",
		Messages: messages,
	})
}

func (h *MessageHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	currentUser, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.ConversationsResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	conversations, err := h.messageService.GetConversationsByUserID(currentUser.ID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.ConversationsResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.ConversationsResponse{
		Success:       true,
		Message:       "Conversations retrieved",
		Conversations: conversations,
	})
}
