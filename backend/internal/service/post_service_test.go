package service

import (
	"real-time-forum/internal/domain"
	"testing"
)

func TestPostService_CreatePost(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, "testuser", "test@example.com", "password123", "John", "Doe", 25, "male")

	tests := []struct {
		name        string
		postData    domain.CreatePostRequest
		expectError bool
	}{
		{
			name: "valid post",
			postData: domain.CreatePostRequest{
				Title:    "Test Post",
				Content:  "This is a test post content",
				Category: "general",
			},
			expectError: false,
		},
		{
			name: "empty title",
			postData: domain.CreatePostRequest{
				Title:    "",
				Content:  "This is a test post content",
				Category: "general",
			},
			expectError: true,
		},
		{
			name: "content too short",
			postData: domain.CreatePostRequest{
				Title:    "Test Post",
				Content:  "Short",
				Category: "general",
			},
			expectError: true,
		},
		{
			name: "title too short",
			postData: domain.CreatePostRequest{
				Title:    "Hi",
				Content:  "This is a valid content with enough characters",
				Category: "general",
			},
			expectError: true,
		},
		{
			name: "empty category",
			postData: domain.CreatePostRequest{
				Title:    "Test Post",
				Content:  "This is a valid content with enough characters",
				Category: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := services.Post.CreatePost(userID, tt.postData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if post == nil {
				t.Fatal("Expected post but got nil")
			}

			if post.Title != tt.postData.Title {
				t.Errorf("Expected title %s, got %s", tt.postData.Title, post.Title)
			}

			if post.Content != tt.postData.Content {
				t.Errorf("Expected content %s, got %s", tt.postData.Content, post.Content)
			}

			if post.Category != tt.postData.Category {
				t.Errorf("Expected category %s, got %s", tt.postData.Category, post.Category)
			}

			if post.UserID != userID {
				t.Errorf("Expected user ID %d, got %d", userID, post.UserID)
			}
		})
	}
}

func TestPostService_GetPostByID(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, "testuser", "test@example.com", "password123", "John", "Doe", 25, "male")

	post, err := services.Post.CreatePost(userID, domain.CreatePostRequest{
		Title:    "Test Post",
		Content:  "This is a test post content with enough characters",
		Category: "general",
	})
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	_, err = services.Comment.CreateComment(userID, post.ID, domain.CreateCommentRequest{
		Content: "This is a test comment",
	})
	if err != nil {
		t.Fatalf("Failed to create test comment: %v", err)
	}

	t.Run("get existing post", func(t *testing.T) {
		retrievedPost, comments, err := services.Post.GetPostByID(post.ID)
		if err != nil {
			t.Fatalf("Failed to get post: %v", err)
		}

		if retrievedPost == nil {
			t.Fatal("Expected post but got nil")
		}

		if retrievedPost.ID != post.ID {
			t.Errorf("Expected post ID %d, got %d", post.ID, retrievedPost.ID)
		}

		if retrievedPost.Title != post.Title {
			t.Errorf("Expected title %s, got %s", post.Title, retrievedPost.Title)
		}

		if len(comments) != 1 {
			t.Errorf("Expected 1 comment, got %d", len(comments))
		}

		if comments[0].Content != "This is a test comment" {
			t.Errorf("Expected comment content 'This is a test comment', got '%s'", comments[0].Content)
		}
	})

	t.Run("get non-existent post", func(t *testing.T) {
		_, _, err := services.Post.GetPostByID(99999)
		if err == nil {
			t.Error("Expected error for non-existent post")
		}
	})
}

func TestPostService_ListPosts(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, "testuser", "test@example.com", "password123", "John", "Doe", 25, "male")

	posts := []domain.CreatePostRequest{
		{Title: "Post 1", Content: "Content 1 with enough characters", Category: "general"},
		{Title: "Post 2", Content: "Content 2 with enough characters", Category: "tech"},
		{Title: "Post 3", Content: "Content 3 with enough characters", Category: "general"},
	}

	for _, postData := range posts {
		_, err := services.Post.CreatePost(userID, postData)
		if err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}
	}

	t.Run("list all posts", func(t *testing.T) {
		retrievedPosts, err := services.Post.ListPosts("", 10, 0)
		if err != nil {
			t.Fatalf("Failed to list posts: %v", err)
		}

		if len(retrievedPosts) != 3 {
			t.Errorf("Expected 3 posts, got %d", len(retrievedPosts))
		}
	})

	t.Run("list posts by category", func(t *testing.T) {
		retrievedPosts, err := services.Post.ListPosts("general", 10, 0)
		if err != nil {
			t.Fatalf("Failed to list posts by category: %v", err)
		}

		if len(retrievedPosts) != 2 {
			t.Errorf("Expected 2 general posts, got %d", len(retrievedPosts))
		}

		for _, post := range retrievedPosts {
			if post.Category != "general" {
				t.Errorf("Expected category 'general', got '%s'", post.Category)
			}
		}
	})

	t.Run("list posts with pagination", func(t *testing.T) {
		retrievedPosts, err := services.Post.ListPosts("", 2, 0)
		if err != nil {
			t.Fatalf("Failed to list posts with pagination: %v", err)
		}

		if len(retrievedPosts) != 2 {
			t.Errorf("Expected 2 posts with limit 2, got %d", len(retrievedPosts))
		}
	})
}

func TestPostService_GetPostsByUserID(t *testing.T) {
	services := SetupTestServices(t)

	user1ID := CreateTestUser(t, services, "user1", "user1@example.com", "password123", "John", "Doe", 25, "male")
	user2ID := CreateTestUser(t, services, "user2", "user2@example.com", "password123", "Jane", "Smith", 30, "female")

	_, err := services.Post.CreatePost(user1ID, domain.CreatePostRequest{
		Title: "User1 Post 1", Content: "Content 1 with enough characters for validation", Category: "general",
	})
	if err != nil {
		t.Fatalf("Failed to create post for user1: %v", err)
	}

	_, err = services.Post.CreatePost(user1ID, domain.CreatePostRequest{
		Title: "User1 Post 2", Content: "Content 2 with enough characters for validation", Category: "tech",
	})
	if err != nil {
		t.Fatalf("Failed to create post for user1: %v", err)
	}

	_, err = services.Post.CreatePost(user2ID, domain.CreatePostRequest{
		Title: "User2 Post", Content: "Content 3 with enough characters for validation", Category: "general",
	})
	if err != nil {
		t.Fatalf("Failed to create post for user2: %v", err)
	}

	t.Run("get posts by user ID", func(t *testing.T) {
		posts, err := services.Post.GetPostsByUserID(user1ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get posts by user ID: %v", err)
		}

		if len(posts) != 2 {
			t.Errorf("Expected 2 posts for user1, got %d", len(posts))
		}

		for _, post := range posts {
			if post.UserID != user1ID {
				t.Errorf("Expected user ID %d, got %d", user1ID, post.UserID)
			}
		}
	})
}
