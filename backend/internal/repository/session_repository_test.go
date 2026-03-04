package repository

import (
	"testing"
	"time"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	repos := SetupTestDB(t)
	sessionRepo := repos.Session
	userRepo := repos.User

	userID, _ := userRepo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")

	err := sessionRepo.CreateSession("session123", int(userID), time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
}

func TestSessionRepository_GetSessionBySessionID(t *testing.T) {
	repos := SetupTestDB(t)
	sessionRepo := repos.Session
	userRepo := repos.User
	sessionID := "session123"

	userID, _ := userRepo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")

	err := sessionRepo.CreateSession(sessionID, int(userID), time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	session, err := sessionRepo.GetSessionBySessionID(sessionID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if session.ID != sessionID {
		t.Errorf("Expected session ID '%s', got '%s'", sessionID, session.ID)
	}
	if session.UserID != int(userID) {
		t.Errorf("Expected user ID %d, got %d", userID, session.UserID)
	}
}

func TestSessionRepository_DeleteSession(t *testing.T) {
	repos := SetupTestDB(t)
	sessionRepo := repos.Session
	userRepo := repos.User
	sessionID := "session123"

	userID, _ := userRepo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")

	err := sessionRepo.CreateSession(sessionID, int(userID), time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	err = sessionRepo.DeleteSession(sessionID)
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	session, err := sessionRepo.GetSessionBySessionID(sessionID)
	if err == nil {
		t.Fatalf("Expected error when getting deleted session, got nil")
	}

	if session != nil {
		t.Error("Expected nil session after deletion")
	}
}
