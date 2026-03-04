package websocket

import (
	"encoding/json"
	"real-time-forum/internal/domain"
)

func (hub *Hub) RegisterClientToHub(client *Client) {
	hub.register <- client
}

func (hub *Hub) UnregisterClientFromHub(client *Client) {
	hub.unregister <- client
}

func (hub *Hub) BroadcastMessage(message *domain.Message, receiverID int) {
	wsMessage := WsMessage{
		Type:    "new_message",
		Payload: message,
	}

	data, err := json.Marshal(wsMessage)
	if err != nil {
		hub.logger.Error("Failed to marshal message", "error", err)
		return
	}

	hub.mu.RLock()
	if client, ok := hub.clients[receiverID]; ok {
		client.Send <- data
	}

	if client, ok := hub.clients[message.SenderID]; ok {
		client.Send <- data
	}
	hub.mu.RUnlock()
}

func (hub *Hub) broadcastUserStatus(userID int, online bool) {
	user, err := hub.userRepo.GetUserByID(userID)
	if err != nil {
		hub.logger.Error("Failed to get user status", "userID", userID, "error", err)
		return
	}

	payload := domain.UserStatus{
		UserID:   userID,
		Nickname: user.Nickname,
		Online:   online,
	}

	message := WsMessage{
		Type:    "user_status",
		Payload: payload,
	}

	data, err := json.Marshal(message)
	if err != nil {
		hub.logger.Error("Failed to marshal message", "error", err)
		return
	}

	for id, client := range hub.clients {
		if id != userID {
			client.Send <- data
		}
	}
}

func (hub *Hub) sendOnlineUsers(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	var users []domain.UserStatus
	for userID := range hub.clients {
		user, err := hub.userRepo.GetUserByID(userID)
		if err == nil {
			users = append(users, domain.UserStatus{
				UserID:   userID,
				Nickname: user.Nickname,
				Online:   true,
			})
		}
	}

	message := WsMessage{
		Type:    "online_users",
		Payload: map[string]any{"users": users},
	}

	data, err := json.Marshal(message)
	if err != nil {
		hub.logger.Error("Failed to marshal message", "error", err)
		return
	}

	client.Send <- data
}
