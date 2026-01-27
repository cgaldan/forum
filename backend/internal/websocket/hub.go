package websocket

import (
	"encoding/json"
	"sync"

	"forum-backend/internal/domain"
	"forum-backend/internal/repository"
	"forum-backend/pkg/logger"
)

type Hub struct {
	clients    map[int]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	logger     *logger.Logger
	userRepo   *repository.UserRepository
}

func NewHub(logger *logger.Logger, userRepo *repository.UserRepository) *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
		userRepo:   userRepo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

			h.userRepo.UpdateLastSeen(client.UserID)

			h.broadcastUserStatus(client.UserID, true)

			h.sendOnlineUsers(client)

			h.logger.Info("WebSocket client connected", "userID", client.UserID, "totalClients", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()

			h.userRepo.UpdateLastSeen(client.UserID)

			h.broadcastUserStatus(client.UserID, false)

			h.logger.Info("WebSocket client disconnected", "userID", client.UserID, "totalClients", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.UserID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register registers a client
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister unregisters a client
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// BroadcastMessage broadcasts a message to specific users
func (h *Hub) BroadcastMessage(msg *domain.Message, receiverID int) {
	wsMsg := WSMessage{
		Type:    "new_message",
		Payload: msg,
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		h.logger.Error("Failed to marshal message", "error", err)
		return
	}

	// Send to receiver if online
	h.mu.RLock()
	if client, ok := h.clients[receiverID]; ok {
		select {
		case client.Send <- data:
		default:
		}
	}

	// Also send to sender (for multi-device sync)
	if client, ok := h.clients[msg.SenderID]; ok {
		select {
		case client.Send <- data:
		default:
		}
	}
	h.mu.RUnlock()
}

func (h *Hub) broadcastUserStatus(userID int, online bool) {
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		h.logger.Error("Failed to get user for status broadcast", "error", err, "userID", userID)
		return
	}

	payload := domain.UserStatus{
		UserID:   userID,
		Nickname: user.Nickname,
		Online:   online,
	}

	msg := WSMessage{
		Type:    "user_status",
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("Failed to marshal user status", "error", err)
		return
	}

	for id, client := range h.clients {
		if id != userID {
			client.Send <- data
		}
	}
}

func (h *Hub) sendOnlineUsers(client *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var users []domain.UserStatus
	for userID := range h.clients {
		user, err := h.userRepo.GetByID(userID)
		if err == nil {
			users = append(users, domain.UserStatus{
				UserID:   userID,
				Nickname: user.Nickname,
				Online:   true,
			})
		}
	}

	msg := WSMessage{
		Type: "online_users",
		Payload: map[string]interface{}{
			"users": users,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("Failed to marshal online users", "error", err)
		return
	}

	select {
	case client.Send <- data:
	default:
	}
}

// GetOnlineUsers returns list of online user IDs
func (h *Hub) GetOnlineUsers() []int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userIDs := make([]int, 0, len(h.clients))
	for userID := range h.clients {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}
