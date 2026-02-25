package repository

import (
	"testing"
)

func TestPostRepository_CreatePost(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID, err := userRepo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")

	postID, err := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General")
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	if postID == 0 {
		t.Error("Expected non-zero post ID")
	}
}

func TestPostRepository_GetPostByID(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID, _ := userRepo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")
	postID, _ := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General")

	post, err := postRepo.GetPostByID(int(postID))
	if err != nil {
		t.Fatalf("Failed to get post: %v", err)
	}

	if post.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got '%s'", post.Title)
	}
	if post.Content != "This is a test post content." {
		t.Errorf("Expected content 'This is a test post content.', got '%s'", post.Content)
	}
}

func TestPostRepository_ListPosts(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID, _ := userRepo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")
	for i := 1; i <= 5; i++ {
		postRepo.CreatePost(int(userID), "Test Post "+string(rune(i)), "This is test post content "+string(rune(i)), "General")
	}

	posts, err := postRepo.ListPosts("", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list posts: %v", err)
	}

	if len(posts) != 5 {
		t.Errorf("Expected 5 posts, got %d", len(posts))
	}
}

func TestPostRepository_GetPostsByUserID(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID1, _ := userRepo.CreateUser("user1", "user1@example.com", "hashedpass1", "User", "One", 30, "male")
	userID2, _ := userRepo.CreateUser("user2", "user2@example.com", "hashedpass2", "User", "Two", 35, "female")

	postRepo.CreatePost(int(userID1), "User 1 Post", "Content for user 1", "General")
	postRepo.CreatePost(int(userID2), "User 2 Post", "Content for user 2", "General")

	posts, err := postRepo.GetPostsByUserID(int(userID1), 10, 0)
	if err != nil {
		t.Fatalf("Failed to get posts by user ID: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("Expected 1 post for user ID 1, got %d", len(posts))
	}
	if posts[0].Title != "User 1 Post" {
		t.Errorf("Expected title 'User 1 Post', got '%s'", posts[0].Title)
	}
}

func TestPostRepository_Exists(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	postRepo := repos.Post

	userID, _ := userRepo.CreateUser("testuser", "test@example.com", "hashedpass", "John", "Doe", 25, "male")
	postID, _ := postRepo.CreatePost(int(userID), "Test Post", "This is a test post content.", "General")

	exists, err := postRepo.PostExists(int(postID))
	if err != nil {
		t.Fatalf("Failed to check if post exists: %v", err)
	}

	if !exists {
		t.Error("Expected post to exist")
	}
}
