package router

import (
	"net/http"

	"forum-backend/internal/api/handlers"
	"forum-backend/internal/api/middleware"
	"forum-backend/internal/config"
	"forum-backend/internal/service"
	"forum-backend/internal/websocket"
	"forum-backend/pkg/logger"

	"github.com/gorilla/mux"
)

// NewRouter creates and configures the application router
func NewRouter(services *service.Services, hub *websocket.Hub, cfg *config.Config, log *logger.Logger) *mux.Router {
	r := mux.NewRouter()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(services.Auth, log)
	postHandler := handlers.NewPostHandler(services.Post, services.Auth, log)
	commentHandler := handlers.NewCommentHandler(services.Comment, services.Auth, log)
	messageHandler := handlers.NewMessageHandler(services.Message, services.Auth, log)
	wsHandler := handlers.NewWebSocketHandler(hub, services.Auth, log)
	healthHandler := handlers.NewHealthHandler("1.0.0")

	// Health check endpoint (no middleware)
	r.HandleFunc("/health", healthHandler.Health).Methods("GET")

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	api.HandleFunc("/auth/me", authHandler.Me).Methods("GET")

	// Post routes
	api.HandleFunc("/posts", postHandler.GetPosts).Methods("GET")
	api.HandleFunc("/posts", postHandler.CreatePost).Methods("POST")
	api.HandleFunc("/posts/{id}", postHandler.GetPost).Methods("GET")
	api.HandleFunc("/posts/{id}/comments", commentHandler.CreateComment).Methods("POST")

	// Message routes
	api.HandleFunc("/messages/conversations", messageHandler.GetConversations).Methods("GET")
	api.HandleFunc("/messages/{id}", messageHandler.GetMessages).Methods("GET")
	api.HandleFunc("/messages/{id}", messageHandler.SendMessage).Methods("POST")

	// Online users route
	api.HandleFunc("/users/online", wsHandler.GetOnlineUsers).Methods("GET")

	// WebSocket route
	r.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Serve static files from frontend directory
	frontendPath := "../frontend"
	if cfg.Environment == "production" {
		frontendPath = "./frontend"
	}
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(frontendPath)))

	// Apply middleware (order matters!)
	r.Use(middleware.RecoveryMiddleware(log))
	r.Use(middleware.LoggingMiddleware(log))
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware(cfg))
	r.Use(middleware.RateLimitMiddleware(cfg))

	return r
}

