package service

import (
	"io"
	"real-time-forum/internal/database"
	"real-time-forum/internal/repository"
	"real-time-forum/packages/logger"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func SetupTestServices(t *testing.T) *Services {
	t.Helper()

	db, err := database.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	repos := repository.NewRepositories(db)

	testLogger := logger.NewLogger(io.Discard, logger.InfoLevel)

	services := NewServices(repos, testLogger)

	return services
}

func CreateTestUser(t *testing.T, services *Services, nickname, email, password, firstName, lastName string, age int, gender string) int {
	t.Helper()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	userID, err := services.Auth.(*AuthService).userRepo.CreateUser(
		nickname, email, string(hashedPassword), firstName, lastName, age, gender,
	)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return int(userID)
}

func CreateTestPost(t *testing.T, services *Services, userID int, title, content, category string) int {
	t.Helper()

	postID, err := services.Post.(*PostService).postRepo.CreatePost(userID, title, content, category)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	return int(postID)
}

func CreateTestMessage(t *testing.T, services *Services, senderID, receiverID int, content string) int {
	t.Helper()

	message, err := services.Message.SendMessage(senderID, receiverID, content)
	if err != nil {
		t.Fatalf("Failed to create test message: %v", err)
	}

	return message.ID
}
