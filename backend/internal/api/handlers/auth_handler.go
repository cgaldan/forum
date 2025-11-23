package handlers

import (
	"encoding/json"
	"net/http"

	"forum-backend/internal/domain"
	"forum-backend/internal/service"
	"forum-backend/pkg/logger"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
	logger      *logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	user, token, err := h.authService.Register(req)
	if err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(domain.AuthResponse{
		Success: true,
		Message: "User registered successfully",
		User:    user,
		Token:   token,
	})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	user, token, err := h.authService.Login(req)
	if err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(domain.AuthResponse{
		Success: true,
		Message: "Login successful",
		User:    user,
		Token:   token,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(domain.AuthResponse{Success: false, Message: "No token provided"})
		return
	}

	if err := h.authService.Logout(token); err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(domain.AuthResponse{Success: true, Message: "Logout successful"})
}

// Me handles getting current user info
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(domain.AuthResponse{Success: false, Message: "No token provided"})
		return
	}

	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(domain.AuthResponse{Success: true, User: user})
}

