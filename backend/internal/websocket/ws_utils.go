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
	clientsToSend := []*Client{}
	if client, ok := hub.clients[receiverID]; ok {
		clientsToSend = append(clientsToSend, client)
	}
	if client, ok := hub.clients[message.SenderID]; ok {
		clientsToSend = append(clientsToSend, client)
	}
	hub.mu.RUnlock()

	var toRemove []int
	for _, client := range clientsToSend {
		select {
		case client.Send <- data:
		default:
			toRemove = append(toRemove, client.UserID)
		}
	}

	if len(toRemove) > 0 {
		hub.mu.Lock()
		for _, id := range toRemove {
			if client, ok := hub.clients[id]; ok {
				close(client.Send)
				delete(hub.clients, id)
			}
		}
		hub.mu.Unlock()
	}
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

	hub.mu.RLock()
	clients := make([]*Client, 0, len(hub.clients))
	for _, client := range hub.clients {
		clients = append(clients, client)
	}
	hub.mu.RUnlock()

	var toRemove []int
	for _, client := range clients {
		select {
		case client.Send <- data:
		default:
			toRemove = append(toRemove, client.UserID)
		}
	}

	if len(toRemove) > 0 {
		hub.mu.Lock()
		for _, id := range toRemove {
			if client, ok := hub.clients[id]; ok {
				close(client.Send)
				delete(hub.clients, id)
			}
		}
		hub.mu.Unlock()
	}
}

func (hub *Hub) sendOnlineUsers(_ *Client) {
	hub.mu.RLock()
	userIDs := make([]int, 0, len(hub.clients))
	clientsSnapshot := make(map[int]*Client, len(hub.clients))
	for userID, c := range hub.clients {
		userIDs = append(userIDs, userID)
		clientsSnapshot[userID] = c
	}
	hub.mu.RUnlock()

	var users []domain.UserStatus
	for _, userID := range userIDs {
		user, err := hub.userRepo.GetUserByID(userID)
		if err != nil {
			continue
		}
		users = append(users, domain.UserStatus{
			UserID:   userID,
			Nickname: user.Nickname,
			Online:   true,
		})
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

	clientsToSend := make([]*Client, 0, len(clientsSnapshot))
	for _, c := range clientsSnapshot {
		clientsToSend = append(clientsToSend, c)
	}

	var toRemove []int
	for _, client := range clientsToSend {
		select {
		case client.Send <- data:
		default:
			toRemove = append(toRemove, client.UserID)
		}
	}

	if len(toRemove) > 0 {
		hub.mu.Lock()
		for _, id := range toRemove {
			if client, ok := hub.clients[id]; ok {
				close(client.Send)
				delete(hub.clients, id)
			}
		}
		hub.mu.Unlock()
	}
}
