package service

import (
	"testing"

	"forum-backend/internal/domain"
	"forum-backend/internal/repository"
	"forum-backend/pkg/logger"
)

func setupTestAuthService(t *testing.T) *AuthService {
	db, err := repository.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := repository.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	log := logger.NewLogger(nil, logger.ErrorLevel)

	return NewAuthService(userRepo, sessionRepo, log)
}

func TestAuthService_Register(t *testing.T) {
	service := setupTestAuthService(t)

	req := domain.RegisterRequest{
		Nickname:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Age:       25,
		Gender:    "male",
	}

	user, token, err := service.Register(req)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	if user.Nickname != "testuser" {
		t.Errorf("Expected nickname 'testuser', got '%s'", user.Nickname)
	}
	if token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestAuthService_Register_ValidationError(t *testing.T) {
	service := setupTestAuthService(t)

	tests := []struct {
		name string
		req  domain.RegisterRequest
	}{
		{
			name: "Short nickname",
			req: domain.RegisterRequest{
				Nickname:  "ab",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Age:       25,
				Gender:    "male",
			},
		},
		{
			name: "Invalid email",
			req: domain.RegisterRequest{
				Nickname:  "testuser",
				Email:     "invalid",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Age:       25,
				Gender:    "male",
			},
		},
		{
			name: "Short password",
			req: domain.RegisterRequest{
				Nickname:  "testuser",
				Email:     "test@example.com",
				Password:  "pass",
				FirstName: "John",
				LastName:  "Doe",
				Age:       25,
				Gender:    "male",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.Register(tt.req)
			if err == nil {
				t.Error("Expected validation error, got nil")
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	service := setupTestAuthService(t)

	// Register a user first
	regReq := domain.RegisterRequest{
		Nickname:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Age:       25,
		Gender:    "male",
	}
	service.Register(regReq)

	// Login
	loginReq := domain.LoginRequest{
		Identifier: "testuser",
		Password:   "password123",
	}

	user, token, err := service.Login(loginReq)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if user.Nickname != "testuser" {
		t.Errorf("Expected nickname 'testuser', got '%s'", user.Nickname)
	}
	if token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	service := setupTestAuthService(t)

	// Register a user first
	regReq := domain.RegisterRequest{
		Nickname:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Age:       25,
		Gender:    "male",
	}
	service.Register(regReq)

	// Login with wrong password
	loginReq := domain.LoginRequest{
		Identifier: "testuser",
		Password:   "wrongpassword",
	}

	_, _, err := service.Login(loginReq)
	if err == nil {
		t.Error("Expected error for invalid credentials, got nil")
	}
}

