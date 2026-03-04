package handlers

import (
	"net/http"
	"real-time-forum/internal/service"
	"real-time-forum/internal/websocket"
	"real-time-forum/packages/logger"

	gorillaws "github.com/gorilla/websocket"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub         *websocket.Hub
	authService service.AuthServiceInterface
	logger      *logger.Logger
}

func NewWebSocketHandler(hub *websocket.Hub, authService service.AuthServiceInterface, logger *logger.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		authService: authService,
		logger:      logger,
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.ValidateSession(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := websocket.NewClient(h.hub, conn, user.ID, h.logger)

	h.hub.RegisterClientToHub(client)

	go client.ReadPump()
	go client.WritePump()
}
