package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	// Initialize database
	var err error
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "../db/forum.db"
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Test database connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize database schema
	if err = initDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize WebSocket hub
	hub = newHub()
	go hub.run()

	// Initialize rate limiter (100 requests per minute)
	rateLimiter := newRateLimiter(100, time.Minute)

	// Initialize router
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/api/auth/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/auth/logout", logoutHandler).Methods("POST")
	r.HandleFunc("/api/auth/me", meHandler).Methods("GET")

	// Post routes
	r.HandleFunc("/api/posts", getPostsHandler).Methods("GET")
	r.HandleFunc("/api/posts", createPostHandler).Methods("POST")
	r.HandleFunc("/api/posts/{id}", getPostHandler).Methods("GET")
	r.HandleFunc("/api/posts/{id}/comments", createCommentHandler).Methods("POST")

	// WebSocket route
	r.HandleFunc("/ws", wsHandler)

	// Online users route
	r.HandleFunc("/api/users/online", getOnlineUsersHandler).Methods("GET")

	// Message routes
	r.HandleFunc("/api/messages/conversations", getConversationsHandler).Methods("GET")
	r.HandleFunc("/api/messages/{id}", getMessagesHandler).Methods("GET")
	r.HandleFunc("/api/messages/{id}", sendMessageHandler).Methods("POST")

	// Serve static files from frontend directory
	frontendPath := "../frontend"
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(frontendPath)))

	// Apply middleware
	r.Use(loggingMiddleware)
	r.Use(securityHeadersMiddleware)
	r.Use(corsMiddleware)
	r.Use(rateLimitMiddleware(rateLimiter))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Database: %s\n", dbPath)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
