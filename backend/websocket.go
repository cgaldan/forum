package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Client represents a WebSocket client connection
type Client struct {
	ID     int
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	UserID int
}

// Hub maintains active WebSocket connections
type Hub struct {
	clients    map[int]*Client // userID -> Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// WebSocket message types
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type UserStatusPayload struct {
	UserID   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	Online   bool   `json:"online"`
}

type OnlineUsersPayload struct {
	Users []UserStatusPayload `json:"users"`
}

var hub *Hub

func newHub() *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

			// Update user status to online
			updateUserStatus(client.UserID, true)

			// Broadcast user came online
			h.broadcastUserStatus(client.UserID, true)

			// Send current online users to the new client
			h.sendOnlineUsers(client)

			log.Printf("Client connected: UserID=%d, Total clients=%d", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()

			// Update user status to offline
			updateUserStatus(client.UserID, false)

			// Broadcast user went offline
			h.broadcastUserStatus(client.UserID, false)

			log.Printf("Client disconnected: UserID=%d, Total clients=%d", client.UserID, len(h.clients))

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

func (h *Hub) broadcastUserStatus(userID int, online bool) {
	// Get user nickname
	var nickname string
	err := db.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&nickname)
	if err != nil {
		return
	}

	payload := UserStatusPayload{
		UserID:   userID,
		Nickname: nickname,
		Online:   online,
	}

	msg := WSMessage{
		Type:    "user_status",
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.broadcast <- data
}

func (h *Hub) sendOnlineUsers(client *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := []UserStatusPayload{}
	for userID := range h.clients {
		var nickname string
		err := db.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&nickname)
		if err == nil {
			users = append(users, UserStatusPayload{
				UserID:   userID,
				Nickname: nickname,
				Online:   true,
			})
		}
	}

	msg := WSMessage{
		Type: "online_users",
		Payload: OnlineUsersPayload{
			Users: users,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	select {
	case client.Send <- data:
	default:
	}
}

func (h *Hub) getOnlineUsers() []int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userIDs := make([]int, 0, len(h.clients))
	for userID := range h.clients {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming messages
		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		// Handle different message types
		switch wsMsg.Type {
		case "ping":
			// Respond with pong
			pong := WSMessage{Type: "pong"}
			data, _ := json.Marshal(pong)
			c.Send <- data
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verify token and get user
	user, err := getUserFromSession(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    hub,
		UserID: user.ID,
	}

	client.Hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

func updateUserStatus(userID int, online bool) {
	_, err := db.Exec("UPDATE users SET last_seen = CURRENT_TIMESTAMP WHERE id = ?", userID)
	if err != nil {
		log.Printf("Error updating user status: %v", err)
	}
}

// Get online users API endpoint
func getOnlineUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verify authentication
	token := r.Header.Get("Authorization")
	_, err := getUserFromSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userIDs := hub.getOnlineUsers()
	users := []UserStatusPayload{}

	for _, userID := range userIDs {
		var nickname string
		err := db.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&nickname)
		if err == nil {
			users = append(users, UserStatusPayload{
				UserID:   userID,
				Nickname: nickname,
				Online:   true,
			})
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"users":   users,
	})
}

