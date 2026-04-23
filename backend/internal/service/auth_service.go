package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"real-time-forum/internal/domain"
	"real-time-forum/internal/repository"
	"real-time-forum/packages/logger"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    repository.UserRepositoryInterface
	sessionRepo repository.SessionRepositoryInterface
	logger      *logger.Logger
}

func NewAuthService(userRepo repository.UserRepositoryInterface, sessionRepo repository.SessionRepositoryInterface, logger *logger.Logger) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

func (s *AuthService) Register(registrationData domain.RegisterRequest) (*domain.User, string, error) {
	if err := s.validateRegistrationData(registrationData); err != nil {
		return nil, "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationData.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return nil, "", fmt.Errorf("failed to process password")
	}

	userID, err := s.userRepo.CreateUser(
		registrationData.Nickname,
		registrationData.Email,
		string(hashedPassword),
		registrationData.FirstName,
		registrationData.LastName,
		registrationData.Age,
		registrationData.Gender,
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, "", fmt.Errorf("nickname or email already in use")
		}
		s.logger.Error("Failed to create user", "error", err)
		return nil, "", fmt.Errorf("failed to create user")
	}

	user, err := s.userRepo.GetUserByID(int(userID))
	if err != nil {
		s.logger.Error("Failed to retrieve created user", "error", err)
		return nil, "", fmt.Errorf("failed to retrieve user after creation")
	}

	sessionID, err := s.createSession(int(userID))
	if err != nil {
		s.logger.Error("Failed to create session for new user", "error", err)
		return nil, "", fmt.Errorf("failed to create session")
	}

	s.logger.Info("User registered successfully", "userID", userID, "nickname", registrationData.Nickname)
	return user, sessionID, nil
}

func (s *AuthService) Login(loginData domain.LoginRequest) (*domain.User, string, error) {
	if loginData.Identifier == "" || loginData.Password == "" {
		return nil, "", fmt.Errorf("identifier and password are required")
	}

	user, passwordHash, err := s.userRepo.GetUserByIdentifier(loginData.Identifier)
	if err != nil {
		s.logger.Debug("Login failed - user not found", "identifier", loginData.Identifier)
		return nil, "", fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginData.Password)); err != nil {
		s.logger.Debug("Login failed - password mismatch", "identifier", loginData.Identifier)
		return nil, "", fmt.Errorf("invalid identifier or password")
	}

	if err := s.userRepo.UpdateLastSeen(user.ID); err != nil {
		s.logger.Error("Failed to update last seen for user", "userID", user.ID, "error", err)
	}

	sessionID, err := s.createSession(user.ID)
	if err != nil {
		s.logger.Error("Failed to create session for logged in user", "error", err)
		return nil, "", fmt.Errorf("failed to create session")
	}

	s.logger.Info("User logged in successfully", "userID", user.ID, "nickname", user.Nickname)
	return user, sessionID, nil
}

func (s *AuthService) Logout(sessionID string) error {
	if err := s.sessionRepo.DeleteSession(sessionID); err != nil {
		s.logger.Error("Failed to delete session during logout", "sessionID", sessionID, "error", err)
		return fmt.Errorf("failed to logout")
	}

	s.logger.Info("User logged out successfully", "sessionID", sessionID)
	return nil
}

func (s *AuthService) ValidateSession(sessionID string) (*domain.User, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	session, err := s.sessionRepo.GetSessionBySessionID(sessionID)
	if err != nil {
		s.logger.Debug("Session validation failed - session not found", "sessionID", sessionID)
		return nil, fmt.Errorf("invalid session")
	}

	user, err := s.userRepo.GetUserByID(session.UserID)
	if err != nil {
		s.logger.Error("Failed to retrieve user for session validation", "userID", session.UserID, "error", err)
		return nil, fmt.Errorf("invalid session")
	}

	return user, nil
}

// Helper functions

func (s *AuthService) validateRegistrationData(data domain.RegisterRequest) error {
	if data.Nickname == "" || len(data.Nickname) < 3 {
		return fmt.Errorf("nickname must be at least 3 characters")
	}
	if data.Email == "" || !strings.Contains(data.Email, "@") {
		return fmt.Errorf("valid email is required")
	}
	if data.Password == "" || len(data.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	if data.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if data.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if data.Age < 13 || data.Age > 120 {
		return fmt.Errorf("age must be between 13 and 120")
	}
	if data.Gender == "" {
		return fmt.Errorf("gender is required")
	}
	return nil
}

func (s *AuthService) createSession(userID int) (string, error) {
	sessionID := generateSessionID()
	expiresAt := time.Now().Add(24 * time.Hour)

	if err := s.sessionRepo.CreateSession(sessionID, userID, expiresAt); err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	return sessionID, nil
}

func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
