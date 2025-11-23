package repository

import (
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *UserRepository {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return NewUserRepository(db)
}

func TestUserRepository_Create(t *testing.T) {
	repo := setupTestDB(t)

	userID, err := repo.Create("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if userID == 0 {
		t.Error("Expected non-zero user ID")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	repo := setupTestDB(t)

	// Create a user
	userID, _ := repo.Create("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")

	// Get the user
	user, err := repo.GetByID(int(userID))
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

func TestUserRepository_GetByIdentifier(t *testing.T) {
	repo := setupTestDB(t)

	// Create a user
	repo.Create("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")

	// Get by nickname
	user, hash, err := repo.GetByIdentifier("testuser")
	if err != nil {
		t.Fatalf("Failed to get user by nickname: %v", err)
	}
	if user.Nickname != "testuser" {
		t.Errorf("Expected nickname 'testuser', got '%s'", user.Nickname)
	}
	if hash != "hashedpass" {
		t.Errorf("Expected hash 'hashedpass', got '%s'", hash)
	}

	// Get by email
	user, hash, err = repo.GetByIdentifier("test@example.com")
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
}

func TestUserRepository_UpdateLastSeen(t *testing.T) {
	repo := setupTestDB(t)

	// Create a user
	userID, _ := repo.Create("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Update last seen
	err := repo.UpdateLastSeen(int(userID))
	if err != nil {
		t.Fatalf("Failed to update last seen: %v", err)
	}

	// Verify update
	user, _ := repo.GetByID(int(userID))
	if user.LastSeen.Before(time.Now().Add(-1 * time.Second)) {
		t.Error("LastSeen not updated")
	}
}

