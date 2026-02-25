package service

import (
	"real-time-forum/internal/domain"
	"testing"
)

func TestAuthService_Register(t *testing.T) {
	services := SetupTestServices(t)

	tests := []struct {
		name        string
		userData    domain.RegisterRequest
		expectError bool
	}{
		{
			name: "valid registration",
			userData: domain.RegisterRequest{
				Nickname:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Age:       25,
				Gender:    "male",
			},
			expectError: false,
		},
		{
			name: "duplicate nickname",
			userData: domain.RegisterRequest{
				Nickname:  "testuser",
				Email:     "different@example.com",
				Password:  "password123",
				FirstName: "Jane",
				LastName:  "Smith",
				Age:       30,
				Gender:    "female",
			},
			expectError: true,
		},
		{
			name: "invalid email",
			userData: domain.RegisterRequest{
				Nickname:  "user2",
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Age:       25,
				Gender:    "male",
			},
			expectError: true,
		},
		{
			name: "weak password",
			userData: domain.RegisterRequest{
				Nickname:  "user3",
				Email:     "user3@example.com",
				Password:  "123",
				FirstName: "John",
				LastName:  "Doe",
				Age:       25,
				Gender:    "male",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, sessionID, err := services.Auth.Register(tt.userData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if user == nil {
				t.Fatal("Expected user but got nil")
			}

			if user.Nickname != tt.userData.Nickname {
				t.Errorf("Expected nickname %s, got %s", tt.userData.Nickname, user.Nickname)
			}

			if user.Email != tt.userData.Email {
				t.Errorf("Expected email %s, got %s", tt.userData.Email, user.Email)
			}

			if sessionID == "" {
				t.Error("Expected non-empty session ID")
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, "testuser", "test@example.com", "password123", "John", "Doe", 25, "male")

	tests := []struct {
		name        string
		loginData   domain.LoginRequest
		expectError bool
	}{
		{
			name: "valid login with nickname",
			loginData: domain.LoginRequest{
				Identifier: "testuser",
				Password:   "password123",
			},
			expectError: false,
		},
		{
			name: "valid login with email",
			loginData: domain.LoginRequest{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			expectError: false,
		},
		{
			name: "invalid password",
			loginData: domain.LoginRequest{
				Identifier: "testuser",
				Password:   "wrongpassword",
			},
			expectError: true,
		},
		{
			name: "non-existent user",
			loginData: domain.LoginRequest{
				Identifier: "nonexistent",
				Password:   "password123",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, sessionID, err := services.Auth.Login(tt.loginData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if user == nil {
				t.Fatal("Expected user but got nil")
			}

			if user.ID != userID {
				t.Errorf("Expected user ID %d, got %d", userID, user.ID)
			}

			if sessionID == "" {
				t.Error("Expected non-empty session ID")
			}
		})
	}
}

func TestAuthService_ValidateSession(t *testing.T) {
	services := SetupTestServices(t)

	CreateTestUser(t, services, "testuser", "test@example.com", "password123", "John", "Doe", 25, "male")

	_, sessionID, err := services.Auth.Login(domain.LoginRequest{
		Identifier: "testuser",
		Password:   "password123",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	t.Run("valid session", func(t *testing.T) {
		user, err := services.Auth.ValidateSession(sessionID)
		if err != nil {
			t.Fatalf("Failed to validate session: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user but got nil")
		}

		if user.Nickname != "testuser" {
			t.Errorf("Expected nickname 'testuser', got '%s'", user.Nickname)
		}
	})

	t.Run("invalid session", func(t *testing.T) {
		_, err := services.Auth.ValidateSession("invalid-session-id")
		if err == nil {
			t.Error("Expected error for invalid session")
		}
	})
}

func TestAuthService_Logout(t *testing.T) {
	services := SetupTestServices(t)

	CreateTestUser(t, services, "testuser", "test@example.com", "password123", "John", "Doe", 25, "male")

	_, sessionID, err := services.Auth.Login(domain.LoginRequest{
		Identifier: "testuser",
		Password:   "password123",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	t.Run("successful logout", func(t *testing.T) {
		err := services.Auth.Logout(sessionID)
		if err != nil {
			t.Fatalf("Failed to logout: %v", err)
		}

		_, err = services.Auth.ValidateSession(sessionID)
		if err == nil {
			t.Error("Expected error after logout")
		}
	})

	t.Run("logout with invalid session", func(t *testing.T) {
		err := services.Auth.Logout("invalid-session-id")
		if err != nil {
			t.Errorf("Logout with invalid session should not error, got: %v", err)
		}
	})
}
