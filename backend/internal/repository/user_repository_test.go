package repository

import (
	"testing"
	"time"
)

func TestUserRepository_CreateUser(t *testing.T) {
	repos := SetupTestDB(t)
	repo := repos.User

	userID, err := repo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if userID == 0 {
		t.Error("Expected non-zero user ID")
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	repos := SetupTestDB(t)
	repo := repos.User

	userID, _ := repo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")

	user, err := repo.GetUserByID(int(userID))
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if user.Nickname != "testuser" {
		t.Errorf("Expected nickname 'testuser', got '%s'", user.Nickname)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
}

func TestUserRepository_GetUserByIdentifier(t *testing.T) {
	repos := SetupTestDB(t)
	repo := repos.User

	userID, _ := repo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")

	user, pass, err := repo.GetUserByIdentifier("testuser")
	if err != nil {
		t.Fatalf("Failed to get user by identifier: %v", err)
	}

	if user.ID != int(userID) {
		t.Errorf("Expected user ID %d, got %d", userID, user.ID)
	}
	if pass != "hashedpass" {
		t.Errorf("Expected password hash 'hashedpass', got '%s'", pass)
	}
}

func TestUserRepository_UpdateLastSeen(t *testing.T) {
	repos := SetupTestDB(t)
	repo := repos.User

	userID, _ := repo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")

	time.Sleep(100 * time.Millisecond)

	err := repo.UpdateLastSeen(int(userID))
	if err != nil {
		t.Fatalf("Failed to update last seen: %v", err)
	}

	user, _ := repo.GetUserByID(int(userID))
	if user.LastSeen.Before(time.Now().Add(-1 * time.Minute)) {
		t.Error("Expected last seen to be updated to recent time")
	}
}
