package router

import (
	"net/http"
	"real-time-forum/internal/api/handlers"
	"real-time-forum/internal/api/middleware"
	"real-time-forum/internal/config"
	"real-time-forum/internal/service"
	"real-time-forum/internal/websocket"
	"real-time-forum/packages/logger"

	"github.com/gorilla/mux"
)

func NewRouter(services *service.Services, config *config.Config, hub *websocket.Hub, logger *logger.Logger) *mux.Router {
	r := mux.NewRouter()

	authHandler := handlers.NewAuthHandler(services.Auth, logger)
	postHandler := handlers.NewPostHandler(services.Post, services.Auth, logger)
	commentHandler := handlers.NewCommentHandler(services.Comment, services.Auth, logger)
	messageHandler := handlers.NewMessageHandler(services.Message, services.Auth, logger)
	websocketHandler := handlers.NewWebSocketHandler(hub, services.Auth, logger)
	healthHandler := handlers.NewHealthHandler("1.0.0")

	r.HandleFunc("/health", healthHandler.Health).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	api.HandleFunc("/auth/me", authHandler.GetCurrentUser).Methods("GET")

	// Post routes
	api.HandleFunc("/posts", postHandler.GetPosts).Methods("GET")
	api.HandleFunc("/posts", postHandler.CreatePost).Methods("POST")
	api.HandleFunc("/posts/{id}", postHandler.GetPostByID).Methods("GET")
	api.HandleFunc("/posts/{id}/comments", commentHandler.CreateComment).Methods("POST")

	// Message routes
	api.HandleFunc("/messages/conversations", messageHandler.GetConversations).Methods("GET")
	api.HandleFunc("/messages/{id}", messageHandler.GetMessages).Methods("GET")
	api.HandleFunc("/messages/{id}", messageHandler.SendMessage).Methods("POST")

	// Websocket routes
	r.HandleFunc("/ws", websocketHandler.HandleWebSocket)

	frontendPath := "../frontend"
	if config.Environment == "production" {
		frontendPath = config.Frontend.Path
	}
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(frontendPath))))

	r.Use(middleware.RecoveryMiddleware(logger))
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware(config))
	r.Use(middleware.RateLimiterMiddleware(config))

	return r
}
