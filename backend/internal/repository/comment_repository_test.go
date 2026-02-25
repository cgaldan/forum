package repository

import (
	"testing"
)

func TestCommentRepository_CreateComment(t *testing.T) {
	repos := SetupTestDB(t)
	commentRepo := repos.Comment
	postRepo := repos.Post
	userRepo := repos.User

	userID, _ := userRepo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")
	postID, _ := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General")

	commentID, err := commentRepo.CreateComment(int(postID), int(userID), "This is a test comment.")
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	if commentID == 0 {
		t.Error("Expected non-zero comment ID")
	}
}

func TestCommentRepository_GetCommentsByPostID(t *testing.T) {
	repos := SetupTestDB(t)
	commentRepo := repos.Comment
	postRepo := repos.Post
	userRepo := repos.User

	userID, _ := userRepo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")
	postID, _ := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General")

	commentRepo.CreateComment(int(postID), int(userID), "This is a test comment.")

	comments, err := commentRepo.GetCommentsByPostID(int(postID))
	if err != nil {
		t.Fatalf("Failed to get comments: %v", err)
	}

	if len(comments) == 0 {
		t.Error("Expected at least one comment")
	}
}

func TestCommentRepository_GetCommentByID(t *testing.T) {
	repos := SetupTestDB(t)
	commentRepo := repos.Comment
	postRepo := repos.Post
	userRepo := repos.User

	userID, _ := userRepo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")
	postID, _ := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General")

	commentID, err := commentRepo.CreateComment(int(postID), int(userID), "This is a test comment.")
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	comment, err := commentRepo.GetCommentByID(int(commentID))
	if err != nil {
		t.Fatalf("Failed to get comment by ID: %v", err)
	}

	if comment.ID != int(commentID) {
		t.Errorf("Expected comment ID %d, got %d", commentID, comment.ID)
	}
}

func TestCommentRepository_GetCommentsByUserID(t *testing.T) {
	repos := SetupTestDB(t)
	commentRepo := repos.Comment
	userRepo := repos.User
	postRepo := repos.Post

	userID, _ := userRepo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")
	postID, _ := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General")

	commentRepo.CreateComment(int(postID), int(userID), "This is a test comment.")

	comments, err := commentRepo.GetCommentsByUserID(int(userID), 10, 0)
	if err != nil {
		t.Fatalf("Failed to get comments by user ID: %v", err)
	}

	if len(comments) == 0 {
		t.Error("Expected at least one comment for user")
	}
}
