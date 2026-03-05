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

	chain := ChainMiddleware(
		middleware.RecoveryMiddleware(logger),
		middleware.LoggingMiddleware(logger),
		middleware.SecurityHeadersMiddleware(),
		middleware.CORSMiddleware(config),
		middleware.RateLimiterMiddleware(config),
	)

	// Health check
	mux.Handle("GET /health", chain(healthHandler.Health))

	// Auth routes
	mux.Handle("POST /api/auth/register", chain(authHandler.Register))
	mux.Handle("POST /api/auth/login", chain(authHandler.Login))
	mux.Handle("POST /api/auth/logout", chain(authHandler.Logout))
	mux.Handle("GET /api/auth/me", chain(authHandler.GetCurrentUser))

	// Post routes
	mux.Handle("GET /api/posts", chain(postHandler.GetPosts))
	mux.Handle("POST /api/posts", chain(postHandler.CreatePost))
	mux.Handle("GET /api/posts/{id}", chain(postHandler.GetPostByID))
	mux.Handle("POST /api/posts/{id}/comments", chain(commentHandler.CreateComment))

	// Message routes
	mux.Handle("GET /api/messages/conversations", chain(messageHandler.GetConversations))
	mux.Handle("GET /api/messages/{id}", chain(messageHandler.GetMessages))
	mux.Handle("POST /api/messages/{id}", chain(messageHandler.SendMessage))

	// WebSocket routes
	mux.Handle("/ws", chain(websocketHandler.HandleWebSocket))

	frontendPath := "../frontend"
	if config.Environment == "production" {
		frontendPath = "./frontend"
	}
	mux.Handle("/", http.FileServer(http.Dir(frontendPath)))

	return mux
}

func ChainMiddleware(chain ...Middleware) func(http.HandlerFunc) http.Handler {
	return func(h http.HandlerFunc) http.Handler {
		var handler http.Handler = h
		for i := len(chain) - 1; i >= 0; i-- {
			handler = chain[i](handler)
		}
		return handler
	}
}

type Middleware func(http.Handler) http.Handler
