package repository

import (
	"database/sql"
	"fmt"

	"forum-backend/internal/domain"
)

// MessageRepository handles message data access
type MessageRepository struct {
	db *sql.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create creates a new message
func (r *MessageRepository) Create(senderID, receiverID int, content string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO messages (sender_id, receiver_id, content)
		VALUES (?, ?, ?)`, senderID, receiverID, content)
	
	if err != nil {
		return 0, fmt.Errorf("failed to create message: %w", err)
	}

	return result.LastInsertId()
}

// GetConversation gets messages between two users
func (r *MessageRepository) GetConversation(user1ID, user2ID, limit, offset int) ([]domain.Message, error) {
	rows, err := r.db.Query(`
		SELECT m.id, m.sender_id, m.receiver_id, m.content, m.created_at, m.read_at, u.nickname
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?)
		   OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?`,
		user1ID, user2ID, user2ID, user1ID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content,
			&msg.CreatedAt, &msg.ReadAt, &msg.SenderName); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// MarkAsRead marks messages as read
func (r *MessageRepository) MarkAsRead(receiverID, senderID int) error {
	_, err := r.db.Exec(`
		UPDATE messages
		SET read_at = CURRENT_TIMESTAMP
		WHERE receiver_id = ? AND sender_id = ? AND read_at IS NULL`,
		receiverID, senderID)
	
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	return nil
}

// GetConversationPartners gets list of users the user has conversations with
func (r *MessageRepository) GetConversationPartners(userID int) ([]int, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT
			CASE 
				WHEN sender_id = ? THEN receiver_id 
				ELSE sender_id 
			END as other_user_id
		FROM messages
		WHERE sender_id = ? OR receiver_id = ?`,
		userID, userID, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to get conversation partners: %w", err)
	}
	defer rows.Close()

	var partners []int
	for rows.Next() {
		var partnerID int
		if err := rows.Scan(&partnerID); err != nil {
			continue
		}
		partners = append(partners, partnerID)
	}

	return partners, nil
}

// GetUnreadCount gets the count of unread messages from a specific user
func (r *MessageRepository) GetUnreadCount(receiverID, senderID int) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM messages
		WHERE receiver_id = ? AND sender_id = ? AND read_at IS NULL`,
		receiverID, senderID).Scan(&count)
	
	return count, err
}

// GetLastMessage gets the last message between two users
func (r *MessageRepository) GetLastMessage(user1ID, user2ID int) (*domain.Message, error) {
	var msg domain.Message
	err := r.db.QueryRow(`
		SELECT m.id, m.sender_id, m.receiver_id, m.content, m.created_at, m.read_at, u.nickname
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?)
		   OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC LIMIT 1`,
		user1ID, user2ID, user2ID, user1ID).Scan(
		&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content,
		&msg.CreatedAt, &msg.ReadAt, &msg.SenderName)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no messages found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get last message: %w", err)
	}

	return &msg, nil
}

