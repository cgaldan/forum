package repository

import (
	"database/sql"
	"fmt"
	"real-time-forum/internal/domain"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateMessage(senderID, receiverID int, content string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO messages (sender_id, receiver_id, content)
		VALUES (?, ?, ?)`, senderID, receiverID, content)

	if err != nil {
		return 0, fmt.Errorf("failed to create message: %w", err)
	}

	return result.LastInsertId()
}

func (r *MessageRepository) GetConversation(userID1, userID2, limit, offset int) ([]domain.Message, error) {
	rows, err := r.db.Query(`
		SELECT m.id, m.sender_id, m.receiver_id, m.content, m.created_at, m.read_at, u.nickname
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?`, userID1, userID2, userID2, userID1, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var message domain.Message
		err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt, &message.ReadAt, &message.SenderName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *MessageRepository) MarkMessageAsRead(receiverID, senderID int) error {
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
		err := rows.Scan(&partnerID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan partner ID: %w", err)
		}
		partners = append(partners, partnerID)
	}

	return partners, nil
}

func (r *MessageRepository) GetUnreadMessagesCount(receiverID, senderID int) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM messages
		WHERE receiver_id = ? AND sender_id = ? AND read_at IS NULL`,
		receiverID, senderID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get unread messages count: %w", err)
	}

	return count, nil
}

func (r *MessageRepository) GetLastMessageBetweenUsers(user1ID, user2ID int) (*domain.Message, error) {
	var message domain.Message
	err := r.db.QueryRow(`
		SELECT m.id, m.sender_id, m.receiver_id, m.content, m.created_at, m.read_at, u.nickname
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?)
		   OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC LIMIT 1`,
		user1ID, user2ID, user2ID, user1ID).Scan(
		&message.ID, &message.SenderID, &message.ReceiverID, &message.Content,
		&message.CreatedAt, &message.ReadAt, &message.SenderName)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no messages found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get last message: %w", err)
	}

	return &message, nil
}
