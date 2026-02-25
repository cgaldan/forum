package service

import (
	"real-time-forum/internal/domain"
	"testing"
)

func TestCommentService_CreateComment(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, "testuser", "test@example.com", "password123", "John", "Doe", 25, "male")

	post, err := services.Post.CreatePost(userID, domain.CreatePostRequest{
		Title:    "Test Post",
		Content:  "This is a test post content",
		Category: "general",
	})
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name        string
		commentData domain.CreateCommentRequest
		expectError bool
	}{
		{
			name: "valid comment",
			commentData: domain.CreateCommentRequest{
				Content: "This is a valid comment",
			},
			expectError: false,
		},
		{
			name: "empty content",
			commentData: domain.CreateCommentRequest{
				Content: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := services.Comment.CreateComment(userID, post.ID, tt.commentData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if comment == nil {
				t.Fatal("Expected comment but got nil")
			}

			if comment.Content != tt.commentData.Content {
				t.Errorf("Expected content %s, got %s", tt.commentData.Content, comment.Content)
			}

			if comment.UserID != userID {
				t.Errorf("Expected user ID %d, got %d", userID, comment.UserID)
			}

			if comment.PostID != post.ID {
				t.Errorf("Expected post ID %d, got %d", post.ID, comment.PostID)
			}
		})
	}
}

func TestCommentService_GetCommentsByPostID(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, "testuser", "test@example.com", "password123", "John", "Doe", 25, "male")

	post, err := services.Post.CreatePost(userID, domain.CreatePostRequest{
		Title:    "Test Post",
		Content:  "This is a test post content",
		Category: "general",
	})
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	comments := []string{
		"First comment",
		"Second comment",
		"Third comment",
	}

	for _, content := range comments {
		_, err := services.Comment.CreateComment(userID, post.ID, domain.CreateCommentRequest{
			Content: content,
		})
		if err != nil {
			t.Fatalf("Failed to create test comment: %v", err)
		}
	}

	t.Run("get comments by post ID", func(t *testing.T) {
		retrievedComments, err := services.Comment.GetCommentsByPostID(post.ID)
		if err != nil {
			t.Fatalf("Failed to get comments by post ID: %v", err)
		}

		if len(retrievedComments) != 3 {
			t.Errorf("Expected 3 comments, got %d", len(retrievedComments))
		}

		expectedContents := []string{"First comment", "Second comment", "Third comment"}
		for i, comment := range retrievedComments {
			if comment.Content != expectedContents[i] {
				t.Errorf("Expected comment content '%s', got '%s'", expectedContents[i], comment.Content)
			}
			if comment.PostID != post.ID {
				t.Errorf("Expected post ID %d, got %d", post.ID, comment.PostID)
			}
		}
	})

	t.Run("get comments for non-existent post", func(t *testing.T) {
		comments, err := services.Comment.GetCommentsByPostID(99999)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(comments) != 0 {
			t.Errorf("Expected 0 comments for non-existent post, got %d", len(comments))
		}
	})
}

func TestCommentService_GetCommentsByUserID(t *testing.T) {
	services := SetupTestServices(t)

	user1ID := CreateTestUser(t, services, "user1", "user1@example.com", "password123", "John", "Doe", 25, "male")
	user2ID := CreateTestUser(t, services, "user2", "user2@example.com", "password123", "Jane", "Smith", 30, "female")

	post1, err := services.Post.CreatePost(user1ID, domain.CreatePostRequest{
		Title: "Post 1", Content: "Content 1 with enough characters", Category: "general",
	})
	if err != nil {
		t.Fatalf("Failed to create post1: %v", err)
	}

	post2, err := services.Post.CreatePost(user2ID, domain.CreatePostRequest{
		Title: "Post 2", Content: "Content 2 with enough characters", Category: "general",
	})
	if err != nil {
		t.Fatalf("Failed to create post2: %v", err)
	}

	_, err = services.Comment.CreateComment(user1ID, post1.ID, domain.CreateCommentRequest{
		Content: "User1 comment on post1",
	})
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	_, err = services.Comment.CreateComment(user1ID, post2.ID, domain.CreateCommentRequest{
		Content: "User1 comment on post2",
	})
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	_, err = services.Comment.CreateComment(user2ID, post1.ID, domain.CreateCommentRequest{
		Content: "User2 comment on post1",
	})
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	t.Run("get comments by user ID", func(t *testing.T) {
		comments, err := services.Comment.GetCommentsByUserID(user1ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get comments by user ID: %v", err)
		}

		if len(comments) != 2 {
			t.Errorf("Expected 2 comments for user1, got %d", len(comments))
		}

		for _, comment := range comments {
			if comment.UserID != user1ID {
				t.Errorf("Expected user ID %d, got %d", user1ID, comment.UserID)
			}
		}
	})

	t.Run("get comments by user ID with pagination", func(t *testing.T) {
		comments, err := services.Comment.GetCommentsByUserID(user1ID, 1, 0)
		if err != nil {
			t.Fatalf("Failed to get comments by user ID with pagination: %v", err)
		}

		if len(comments) != 1 {
			t.Errorf("Expected 1 comment with limit 1, got %d", len(comments))
		}
	})
}
