package handlers

import (
	"encoding/json"
	"net/http"

	"forum-backend/internal/domain"
	"forum-backend/internal/service"
	"forum-backend/internal/websocket"
	"forum-backend/pkg/logger"

	gorillaws "github.com/gorilla/websocket"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (configure properly in production)
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub         *websocket.Hub
	authService *service.AuthService
	logger      *logger.Logger
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *websocket.Hub, authService *service.AuthService, logger *logger.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		authService: authService,
		logger:      logger,
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get token from query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verify token and get user
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade error", "error", err)
		return
	}

	client := &websocket.Client{
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: user.ID,
		Logger: h.logger,
	}

	h.hub.Register(client)

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

// GetOnlineUsers handles getting online users
func (h *WebSocketHandler) GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verify authentication
	token := r.Header.Get("Authorization")
	_, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.OnlineUsersResponse{Success: false, Message: "Unauthorized"})
		return
	}

	// userIDs := h.hub.GetOnlineUsers()
	var users []domain.UserStatus
	// TODO: Get user details for each userID

	json.NewEncoder(w).Encode(domain.OnlineUsersResponse{
		Success: true,
		Users:   users,
	})
}

