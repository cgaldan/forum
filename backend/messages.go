package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Message struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	ReadAt     *time.Time `json:"read_at,omitempty"`
	SenderName string    `json:"sender_name"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
}

type MessageResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
	Messages []Message `json:"messages,omitempty"`
	Msg      *Message  `json:"msg,omitempty"`
}

type Conversation struct {
	UserID       int       `json:"user_id"`
	Nickname     string    `json:"nickname"`
	LastMessage  string    `json:"last_message"`
	LastTime     time.Time `json:"last_time"`
	UnreadCount  int       `json:"unread_count"`
}

type ConversationsResponse struct {
	Success       bool           `json:"success"`
	Conversations []Conversation `json:"conversations"`
}

// Send a private message
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get sender from session
	token := r.Header.Get("Authorization")
	sender, err := getUserFromSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(MessageResponse{Success: false, Message: "Unauthorized"})
		return
	}

	// Get receiver ID from URL
	vars := mux.Vars(r)
	receiverID, err := strconv.Atoi(vars["id"])
	if err != nil {
		json.NewEncoder(w).Encode(MessageResponse{Success: false, Message: "Invalid user ID"})
		return
	}

	// Can't message yourself
	if receiverID == sender.ID {
		json.NewEncoder(w).Encode(MessageResponse{Success: false, Message: "Cannot message yourself"})
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(MessageResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	if req.Content == "" || len(req.Content) > 1000 {
		json.NewEncoder(w).Encode(MessageResponse{Success: false, Message: "Message must be 1-1000 characters"})
		return
	}

	// Check if receiver exists
	var receiverNickname string
	err = db.QueryRow("SELECT nickname FROM users WHERE id = ?", receiverID).Scan(&receiverNickname)
	if err != nil {
		json.NewEncoder(w).Encode(MessageResponse{Success: false, Message: "User not found"})
		return
	}

	// Insert message
	result, err := db.Exec(`
		INSERT INTO messages (sender_id, receiver_id, content)
		VALUES (?, ?, ?)`, sender.ID, receiverID, req.Content)

	if err != nil {
		json.NewEncoder(w).Encode(MessageResponse{Success: false, Message: "Failed to send message"})
		return
	}

	messageID, _ := result.LastInsertId()

	msg := &Message{
		ID:         int(messageID),
		SenderID:   sender.ID,
		ReceiverID: receiverID,
		Content:    req.Content,
		CreatedAt:  time.Now(),
		SenderName: sender.Nickname,
	}

	// Broadcast message via WebSocket
	broadcastMessage(msg, receiverID)

	json.NewEncoder(w).Encode(MessageResponse{
		Success: true,
		Message: "Message sent",
		Msg:     msg,
	})
}

// Get message history with a user
func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get current user from session
	token := r.Header.Get("Authorization")
	currentUser, err := getUserFromSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(MessageResponse{Success: false, Message: "Unauthorized"})
		return
	}

	// Get other user ID from URL
	vars := mux.Vars(r)
	otherUserID, err := strconv.Atoi(vars["id"])
	if err != nil {
		json.NewEncoder(w).Encode(MessageResponse{Success: false, Message: "Invalid user ID"})
		return
	}

	// Get pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // default
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

	// Get messages between the two users
	rows, err := db.Query(`
		SELECT m.id, m.sender_id, m.receiver_id, m.content, m.created_at, m.read_at, u.nickname
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?)
		   OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?`,
		currentUser.ID, otherUserID, otherUserID, currentUser.ID, limit, offset)

	if err != nil {
		json.NewEncoder(w).Encode(MessageResponse{Success: false, Message: "Failed to fetch messages"})
		return
	}
	defer rows.Close()

	messages := []Message{}
	for rows.Next() {
		var msg Message
		var readAt *time.Time
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content,
			&msg.CreatedAt, &readAt, &msg.SenderName)
		if err == nil {
			msg.ReadAt = readAt
			messages = append(messages, msg)
		}
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	// Mark messages as read
	_, err = db.Exec(`
		UPDATE messages
		SET read_at = CURRENT_TIMESTAMP
		WHERE receiver_id = ? AND sender_id = ? AND read_at IS NULL`,
		currentUser.ID, otherUserID)

	json.NewEncoder(w).Encode(MessageResponse{
		Success:  true,
		Messages: messages,
	})
}

// Get list of conversations
func getConversationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get current user from session
	token := r.Header.Get("Authorization")
	currentUser, err := getUserFromSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(ConversationsResponse{Success: false})
		return
	}

	// Get all users the current user has messaged with
	rows, err := db.Query(`
		SELECT DISTINCT
			CASE
				WHEN m.sender_id = ? THEN m.receiver_id
				ELSE m.sender_id
			END as other_user_id,
			u.nickname,
			(SELECT content FROM messages
			 WHERE (sender_id = ? AND receiver_id = other_user_id)
			    OR (sender_id = other_user_id AND receiver_id = ?)
			 ORDER BY created_at DESC LIMIT 1) as last_message,
			(SELECT created_at FROM messages
			 WHERE (sender_id = ? AND receiver_id = other_user_id)
			    OR (sender_id = other_user_id AND receiver_id = ?)
			 ORDER BY created_at DESC LIMIT 1) as last_time,
			(SELECT COUNT(*) FROM messages
			 WHERE receiver_id = ? AND sender_id = other_user_id AND read_at IS NULL) as unread_count
		FROM messages m
		JOIN users u ON (
			CASE
				WHEN m.sender_id = ? THEN m.receiver_id
				ELSE m.sender_id
			END = u.id
		)
		WHERE m.sender_id = ? OR m.receiver_id = ?
		ORDER BY last_time DESC`,
		currentUser.ID, currentUser.ID, currentUser.ID,
		currentUser.ID, currentUser.ID, currentUser.ID,
		currentUser.ID, currentUser.ID, currentUser.ID)

	if err != nil {
		json.NewEncoder(w).Encode(ConversationsResponse{Success: false})
		return
	}
	defer rows.Close()

	conversations := []Conversation{}
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(&conv.UserID, &conv.Nickname, &conv.LastMessage, &conv.LastTime, &conv.UnreadCount)
		if err == nil {
			conversations = append(conversations, conv)
		}
	}

	json.NewEncoder(w).Encode(ConversationsResponse{
		Success:       true,
		Conversations: conversations,
	})
}

// Broadcast message via WebSocket
func broadcastMessage(msg *Message, receiverID int) {
	if hub == nil {
		return
	}

	wsMsg := WSMessage{
		Type:    "new_message",
		Payload: msg,
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		return
	}

	// Send to receiver if online
	hub.mu.RLock()
	if client, ok := hub.clients[receiverID]; ok {
		select {
		case client.Send <- data:
		default:
		}
	}
	hub.mu.RUnlock()

	// Also send to sender (for multi-device sync)
	hub.mu.RLock()
	if client, ok := hub.clients[msg.SenderID]; ok {
		select {
		case client.Send <- data:
		default:
		}
	}
	hub.mu.RUnlock()
}

