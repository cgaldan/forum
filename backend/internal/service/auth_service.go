package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"forum-backend/internal/domain"
	"forum-backend/internal/repository"
	"forum-backend/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	logger      *logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, logger *logger.Logger) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

// Register registers a new user
func (s *AuthService) Register(req domain.RegisterRequest) (*domain.User, string, error) {
	// Validate
	if err := s.validateRegistration(req); err != nil {
		return nil, "", err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return nil, "", fmt.Errorf("failed to process password")
	}

	// Create user
	userID, err := s.userRepo.Create(req.Nickname, req.Email, string(hashedPassword),
		req.FirstName, req.LastName, req.Age, req.Gender)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, "", fmt.Errorf("nickname or email already exists")
		}
		s.logger.Error("Failed to create user", "error", err)
		return nil, "", fmt.Errorf("failed to create user")
	}

	// Get created user
	user, err := s.userRepo.GetByID(int(userID))
	if err != nil {
		s.logger.Error("Failed to get created user", "error", err, "userID", userID)
		return nil, "", fmt.Errorf("user created but failed to retrieve")
	}

	// Create session
	sessionID, err := s.createSession(int(userID))
	if err != nil {
		s.logger.Error("Failed to create session", "error", err, "userID", userID)
		return nil, "", fmt.Errorf("user created but failed to create session")
	}

	s.logger.Info("User registered successfully", "userID", userID, "nickname", req.Nickname)
	return user, sessionID, nil
}

// Login authenticates a user
func (s *AuthService) Login(req domain.LoginRequest) (*domain.User, string, error) {
	if req.Identifier == "" || req.Password == "" {
		return nil, "", fmt.Errorf("identifier and password are required")
	}

	// Get user by identifier
	user, passwordHash, err := s.userRepo.GetByIdentifier(req.Identifier)
	if err != nil {
		s.logger.Debug("Login failed: user not found", "identifier", req.Identifier)
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		s.logger.Debug("Login failed: invalid password", "userID", user.ID)
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Update last seen
	s.userRepo.UpdateLastSeen(user.ID)

	// Create session
	sessionID, err := s.createSession(user.ID)
	if err != nil {
		s.logger.Error("Failed to create session", "error", err, "userID", user.ID)
		return nil, "", fmt.Errorf("failed to create session")
	}

	s.logger.Info("User logged in successfully", "userID", user.ID, "nickname", user.Nickname)
	return user, sessionID, nil
}

// Logout logs out a user
func (s *AuthService) Logout(sessionID string) error {
	if err := s.sessionRepo.Delete(sessionID); err != nil {
		s.logger.Error("Failed to delete session", "error", err, "sessionID", sessionID)
		return fmt.Errorf("failed to logout")
	}

	s.logger.Info("User logged out successfully", "sessionID", sessionID)
	return nil
}

// ValidateSession validates a session and returns the user
func (s *AuthService) ValidateSession(sessionID string) (*domain.User, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	// Get session
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		s.logger.Debug("Invalid session", "sessionID", sessionID, "error", err)
		return nil, fmt.Errorf("invalid or expired session")
	}

	// Get user
	user, err := s.userRepo.GetByID(session.UserID)
	if err != nil {
		s.logger.Error("Failed to get user for session", "error", err, "userID", session.UserID)
		return nil, fmt.Errorf("failed to get user")
	}

	return user, nil
}

// Helper functions

func (s *AuthService) validateRegistration(req domain.RegisterRequest) error {
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

func (s *AuthService) createSession(userID int) (string, error) {
	sessionID := generateSessionID()
	expiresAt := time.Now().Add(24 * time.Hour)

	if err := s.sessionRepo.Create(sessionID, userID, expiresAt); err != nil {
		return "", err
	}

	return sessionID, nil
}

func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

