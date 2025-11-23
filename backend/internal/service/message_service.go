package service

import (
	"fmt"
	"time"

	"forum-backend/internal/domain"
	"forum-backend/internal/repository"
	"forum-backend/internal/websocket"
	"forum-backend/pkg/logger"
)

// MessageService handles messaging business logic
type MessageService struct {
	messageRepo *repository.MessageRepository
	userRepo    *repository.UserRepository
	hub         *websocket.Hub
	logger      *logger.Logger
}

// NewMessageService creates a new message service
func NewMessageService(messageRepo *repository.MessageRepository, userRepo *repository.UserRepository, hub *websocket.Hub, logger *logger.Logger) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		hub:         hub,
		logger:      logger,
	}
}

// SendMessage sends a message to a user
func (s *MessageService) SendMessage(senderID, receiverID int, content string) (*domain.Message, error) {
	// Validate
	if content == "" || len(content) > 1000 {
		return nil, fmt.Errorf("message must be 1-1000 characters")
	}

	if receiverID == senderID {
		return nil, fmt.Errorf("cannot message yourself")
	}

	// // Check if receiver exists
	// receiver, err := s.userRepo.GetByID(receiverID)
	// if err != nil {
	// 	return nil, fmt.Errorf("user not found")
	// }

	// Get sender info
	sender, err := s.userRepo.GetByID(senderID)
	if err != nil {
		return nil, fmt.Errorf("sender not found")
	}

	// Create message
	messageID, err := s.messageRepo.Create(senderID, receiverID, content)
	if err != nil {
		s.logger.Error("Failed to create message", "error", err, "senderID", senderID, "receiverID", receiverID)
		return nil, fmt.Errorf("failed to send message")
	}

	// Build message object
	message := &domain.Message{
		ID:         int(messageID),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  time.Now(),
		SenderName: sender.Nickname,
	}

	// Broadcast via WebSocket
	s.hub.BroadcastMessage(message, receiverID)

	s.logger.Info("Message sent successfully", "messageID", messageID, "from", senderID, "to", receiverID)
	return message, nil
}

// GetConversation gets messages between two users
func (s *MessageService) GetConversation(user1ID, user2ID, limit, offset int) ([]domain.Message, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	messages, err := s.messageRepo.GetConversation(user1ID, user2ID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get conversation", "error", err)
		return nil, fmt.Errorf("failed to fetch messages")
	}

	// Mark messages as read
	if err := s.messageRepo.MarkAsRead(user1ID, user2ID); err != nil {
		s.logger.Warn("Failed to mark messages as read", "error", err)
	}

	return messages, nil
}

// GetConversations gets all conversations for a user
func (s *MessageService) GetConversations(userID int) ([]domain.Conversation, error) {
	// Get conversation partners
	partners, err := s.messageRepo.GetConversationPartners(userID)
	if err != nil {
		s.logger.Error("Failed to get conversation partners", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to fetch conversations")
	}

	var conversations []domain.Conversation
	for _, partnerID := range partners {
		// Get user details
		user, err := s.userRepo.GetByID(partnerID)
		if err != nil {
			continue
		}

		// Get last message
		lastMsg, err := s.messageRepo.GetLastMessage(userID, partnerID)
		if err != nil {
			continue
		}

		// Get unread count
		unreadCount, _ := s.messageRepo.GetUnreadCount(userID, partnerID)

		conversations = append(conversations, domain.Conversation{
			UserID:      partnerID,
			Nickname:    user.Nickname,
			LastMessage: lastMsg.Content,
			LastTime:    lastMsg.CreatedAt,
			UnreadCount: unreadCount,
		})
	}

	return conversations, nil
}

