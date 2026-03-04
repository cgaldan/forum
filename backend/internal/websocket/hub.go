package websocket

import (
	"real-time-forum/internal/repository"
	"real-time-forum/packages/logger"
	"sync"
)

type Hub struct {
	clients    map[int]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	logger     *logger.Logger
	userRepo   repository.UserRepositoryInterface
}

func NewHub(logger *logger.Logger, userRepo repository.UserRepositoryInterface) *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mu:         sync.RWMutex{},
		logger:     logger,
		userRepo:   userRepo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClientCase(client)
		case client := <-h.unregister:
			h.unregisterClientCase(client)
		case message := <-h.broadcast:
			h.broadcastMessageCase(message)
		}
	}
}

func (h *Hub) registerClientCase(client *Client) {
	h.mu.Lock()
	h.clients[client.UserID] = client
	h.mu.Unlock()

	h.userRepo.UpdateLastSeen(client.UserID)

	h.broadcastUserStatus(client.UserID, true)

	h.sendOnlineUsers(client)

	h.logger.Info("Websocket client connected", "userID", client.UserID, "totalClients", len(h.clients))
}

func (h *Hub) unregisterClientCase(client *Client) {
	h.mu.Lock()
	if _, ok := h.clients[client.UserID]; ok {
		delete(h.clients, client.UserID)
		close(client.Send)
	}
	h.mu.Unlock()

	h.userRepo.UpdateLastSeen(client.UserID)

	h.broadcastUserStatus(client.UserID, false)

	h.logger.Info("Websocket client disconnected", "userID", client.UserID, "totalClients", len(h.clients))
}

func (h *Hub) broadcastMessageCase(message []byte) {
	h.mu.RLock()
	for _, client := range h.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			h.unregister <- client
		}
	}
	h.mu.RUnlock()
}
