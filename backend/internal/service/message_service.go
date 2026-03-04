package service

import (
	"fmt"
	"real-time-forum/internal/domain"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/websocket"
	"real-time-forum/packages/logger"
	"time"
)

type MessageService struct {
	messageRepo repository.MessageRepositoryInterface
	userRepo    repository.UserRepositoryInterface
	hub         *websocket.Hub
	logger      *logger.Logger
}

func NewMessageService(messageRepo repository.MessageRepositoryInterface, userRepo repository.UserRepositoryInterface, hub *websocket.Hub, logger *logger.Logger) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		hub:         hub,
		logger:      logger,
	}
}

func (s *MessageService) SendMessage(senderID, receiverID int, content string) (*domain.Message, error) {
	if err := s.validateMessage(content); err != nil {
		return nil, err
	}

	if receiverID == senderID {
		return nil, fmt.Errorf("cannot send message to yourself")
	}

	_, err := s.userRepo.GetUserByID(receiverID)
	if err != nil {
		s.logger.Error("Failed to find receiver", "error", err, "receiverID", receiverID)
		return nil, fmt.Errorf("receiver not found")
	}

	sender, err := s.userRepo.GetUserByID(senderID)
	if err != nil {
		s.logger.Error("Failed to find sender", "error", err, "senderID", senderID)
		return nil, fmt.Errorf("sender not found")
	}

	messageID, err := s.messageRepo.CreateMessage(senderID, receiverID, content)
	if err != nil {
		s.logger.Error("Failed to create message", "error", err)
		return nil, fmt.Errorf("failed to send message")
	}

	message := &domain.Message{
		ID:         int(messageID),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  time.Now(),
		SenderName: sender.Nickname,
	}

	s.hub.BroadcastMessage(message, receiverID)

	s.logger.Info("Message sent successfully", "messageID", messageID, "senderID", senderID, "receiverID", receiverID)
	return message, nil
}

func (s *MessageService) GetConversation(userID1, userID2, limit, offset int) ([]domain.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	messages, err := s.messageRepo.GetConversation(userID1, userID2, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get conversation", "error", err, "userID1", userID1, "userID2", userID2)
		return nil, fmt.Errorf("failed to get conversation")
	}

	if err := s.messageRepo.MarkMessageAsRead(userID1, userID2); err != nil {
		s.logger.Error("Failed to mark messages as read", "error", err, "userID1", userID1, "userID2", userID2)
	}

	return messages, nil
}

func (s *MessageService) GetConversationsByUserID(userID int) ([]domain.Conversation, error) {
	partners, err := s.messageRepo.GetConversationPartners(userID)
	if err != nil {
		s.logger.Error("Failed to get conversation partners", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to get conversations")
	}

	var conversations []domain.Conversation
	for _, partnerID := range partners {
		user, err := s.userRepo.GetUserByID(partnerID)
		if err != nil {
			s.logger.Error("Failed to get conversation partner details", "error", err, "partnerID", partnerID)
			continue
		}

		lastMessage, err := s.messageRepo.GetLastMessageBetweenUsers(userID, partnerID)
		if err != nil {
			s.logger.Error("Failed to get last message for conversation", "error", err, "userID", userID, "partnerID", partnerID)
			continue
		}

		unreadCount, err := s.messageRepo.GetUnreadMessagesCount(userID, partnerID)
		if err != nil {
			s.logger.Error("Failed to get unread message count for conversation", "error", err, "userID", userID, "partnerID", partnerID)
			continue
		}

		conversations = append(conversations, domain.Conversation{
			UserID:      partnerID,
			Nickname:    user.Nickname,
			LastMessage: lastMessage.Content,
			LastTime:    lastMessage.CreatedAt,
			UnreadCount: unreadCount,
		})
	}

	return conversations, nil
}

func (s *MessageService) validateMessage(content string) error {
	if len(content) == 0 {
		return fmt.Errorf("message content cannot be empty")
	}
	if len(content) > 1000 {
		return fmt.Errorf("message content cannot exceed 1000 characters")
	}
	return nil
}
