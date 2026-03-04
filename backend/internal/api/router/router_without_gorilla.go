package router

import (
	"net/http"
	"real-time-forum/internal/api/handlers"
	"real-time-forum/internal/api/middleware"
	"real-time-forum/internal/config"
	"real-time-forum/internal/service"
	"real-time-forum/internal/websocket"
	"real-time-forum/packages/logger"
)

func NewRouterNoGorilla(services *service.Services, config *config.Config, hub *websocket.Hub, logger *logger.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	authHandler := handlers.NewAuthHandler(services.Auth, logger)
	postHandler := handlers.NewPostHandler(services.Post, services.Auth, logger)
	commentHandler := handlers.NewCommentHandler(services.Comment, services.Auth, logger)
	messageHandler := handlers.NewMessageHandler(services.Message, services.Auth, logger)
	websocketHandler := handlers.NewWebSocketHandler(hub, services.Auth, logger)
	healthHandler := handlers.NewHealthHandler("1.0.0")

	middleware := ChainMiddleware(
		middleware.RecoveryMiddleware(logger),
		middleware.LoggingMiddleware(logger),
		middleware.SecurityHeadersMiddleware(),
		middleware.CORSMiddleware(config),
		middleware.RateLimiterMiddleware(config),
	)

	// Health check
	mux.Handle("GET /health", middleware(healthHandler.Health))

	// Auth routes
	mux.Handle("POST /api/auth/register", middleware(authHandler.Register))
	mux.Handle("POST /api/auth/login", middleware(authHandler.Login))
	mux.Handle("POST /api/auth/logout", middleware(authHandler.Logout))
	mux.Handle("GET /api/auth/me", middleware(authHandler.GetCurrentUser))

	// Post routes
	mux.Handle("GET /api/posts", middleware(postHandler.GetPosts))
	mux.Handle("POST /api/posts", middleware(postHandler.CreatePost))
	mux.Handle("GET /api/posts/{id}", middleware(postHandler.GetPostByID))
	mux.Handle("POST /api/posts/{id}/comments", middleware(commentHandler.CreateComment))

	// Message routes
	mux.Handle("GET /api/messages/conversations", middleware(messageHandler.GetConversations))
	mux.Handle("GET /api/messages/{id}", middleware(messageHandler.GetMessages))
	mux.Handle("POST /api/messages/{id}", middleware(messageHandler.SendMessage))

	// WebSocket routes
	mux.Handle("/api/ws", middleware(websocketHandler.HandleWebSocket))

	return mux
}

func ChainMiddleware(middleware ...Middleware) func(http.HandlerFunc) http.Handler {
	return func(h http.HandlerFunc) http.Handler {
		var handler http.Handler = h
		for i := len(middleware) - 1; i >= 0; i-- {
			handler = middleware[i](handler)
		}
		return handler
	}
}

type Middleware func(http.Handler) http.Handler
