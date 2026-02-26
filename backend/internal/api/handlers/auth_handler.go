package handlers

import (
	"encoding/json"
	"net/http"
	"real-time-forum/internal/domain"
	"real-time-forum/internal/service"
	"real-time-forum/packages/logger"
)

type AuthHandler struct {
	authService service.AuthServiceInterface
	logger      *logger.Logger
}

func NewAuthHandler(authService service.AuthServiceInterface, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerRequest domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	user, token, err := h.authService.Register(registerRequest)
	if err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.AuthResponse{
		Success: true,
		Message: "Registration successful",
		User:    user,
		Token:   token,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginRequest domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	user, token, err := h.authService.Login(loginRequest)
	if err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.AuthResponse{
		Success: true,
		Message: "Login successful",
		User:    user,
		Token:   token,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(domain.AuthResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	if err := h.authService.Logout(token); err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.AuthResponse{
		Success: true,
		Message: "Logout successful",
	})
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(domain.AuthResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.AuthResponse{
		Success: true,
		Message: "User retrieved successfully",
		User:    user,
	})
}
