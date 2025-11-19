package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

type RegisterRequest struct {
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"` // nickname or email
	Password   string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
	Token   string `json:"token,omitempty"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	// Validate required fields
	if err := validateRegistration(req); err != nil {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Failed to process password"})
		return
	}

	// Insert user
	result, err := db.Exec(`
		INSERT INTO users (nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		req.Nickname, req.Email, string(hashedPassword), req.FirstName, req.LastName, req.Age, req.Gender)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Nickname or email already exists"})
			return
		}
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Failed to create user"})
		return
	}

	userID, _ := result.LastInsertId()

	// Create session
	sessionID, err := createSession(int(userID))
	if err != nil {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "User created but failed to create session"})
		return
	}

	user := &User{
		ID:        int(userID),
		Nickname:  req.Nickname,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       req.Age,
		Gender:    req.Gender,
	}

	json.NewEncoder(w).Encode(AuthResponse{
		Success: true,
		Message: "User registered successfully",
		User:    user,
		Token:   sessionID,
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	if req.Identifier == "" || req.Password == "" {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Identifier and password are required"})
		return
	}

	// Find user by nickname or email
	var user User
	var passwordHash string
	err := db.QueryRow(`
		SELECT id, nickname, email, password_hash, first_name, last_name, age, gender
		FROM users WHERE nickname = ? OR email = ?`,
		req.Identifier, req.Identifier).Scan(
		&user.ID, &user.Nickname, &user.Email, &passwordHash,
		&user.FirstName, &user.LastName, &user.Age, &user.Gender)

	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Invalid credentials"})
		return
	} else if err != nil {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Database error"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Invalid credentials"})
		return
	}

	// Update last seen
	db.Exec("UPDATE users SET last_seen = CURRENT_TIMESTAMP WHERE id = ?", user.ID)

	// Create session
	sessionID, err := createSession(user.ID)
	if err != nil {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Failed to create session"})
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{
		Success: true,
		Message: "Login successful",
		User:    &user,
		Token:   sessionID,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "No token provided"})
		return
	}

	// Remove session
	_, err := db.Exec("DELETE FROM sessions WHERE id = ?", token)
	if err != nil {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Failed to logout"})
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{Success: true, Message: "Logout successful"})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "No token provided"})
		return
	}

	user, err := getUserFromSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "Invalid or expired session"})
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{Success: true, User: user})
}

func validateRegistration(req RegisterRequest) error {
	if req.Nickname == "" || len(req.Nickname) < 3 {
		return fmt.Errorf("nickname must be at least 3 characters")
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return fmt.Errorf("valid email is required")
	}
	if req.Password == "" || len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	if req.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if req.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if req.Age < 13 || req.Age > 120 {
		return fmt.Errorf("age must be between 13 and 120")
	}
	if req.Gender == "" {
		return fmt.Errorf("gender is required")
	}
	return nil
}

func createSession(userID int) (string, error) {
	sessionID := generateSessionID()
	expiresAt := time.Now().Add(24 * time.Hour) // 24 hour session

	_, err := db.Exec(`
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES (?, ?, ?)`, sessionID, userID, expiresAt)

	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func getUserFromSession(sessionID string) (*User, error) {
	var user User
	err := db.QueryRow(`
		SELECT u.id, u.nickname, u.email, u.first_name, u.last_name, u.age, u.gender
		FROM users u
		JOIN sessions s ON u.id = s.user_id
		WHERE s.id = ? AND s.expires_at > CURRENT_TIMESTAMP`,
		sessionID).Scan(&user.ID, &user.Nickname, &user.Email, &user.FirstName, &user.LastName, &user.Age, &user.Gender)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
