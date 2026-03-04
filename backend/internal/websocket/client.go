package websocket

import (
	"encoding/json"
	"real-time-forum/packages/logger"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	Hub        *Hub
	Connection *websocket.Conn
	Send       chan []byte
	UserID     int
	Logger     *logger.Logger
}

type WsMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func NewClient(hub *Hub, connection *websocket.Conn, userID int, logger *logger.Logger) *Client {
	return &Client{
		Hub:        hub,
		Connection: connection,
		Send:       make(chan []byte, 256),
		UserID:     userID,
		Logger:     logger,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Connection.Close()
	}()
	c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	c.Connection.SetPongHandler(func(string) error {
		c.Connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Logger.Error("WebSocket read error", "error", err, "userID", c.UserID)
			}
			break
		}

		var wsMessage WsMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			c.Logger.Warn("Failed to unmarshal message", "error", err)
			continue
		}

		switch wsMessage.Type {
		case "ping":
			pong := WsMessage{
				Type: "pong",
			}
			data, _ := json.Marshal(pong)
			c.Send <- data
		default:
			c.Logger.Debug("Unknown message type", "type", wsMessage.Type)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			writer.Write(message)

			queueSize := len(c.Send)
			for i := 0; i < queueSize; i++ {
				writer.Write(<-c.Send)
			}

			if err := writer.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
