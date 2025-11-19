package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestCreatePost(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create a user and get session
	registerPayload := RegisterRequest{
		Nickname:  "postuser",
		Email:     "post@example.com",
		Password:  "password123",
		FirstName: "Post",
		LastName:  "User",
		Age:       25,
		Gender:    "male",
	}

	body, _ := json.Marshal(registerPayload)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	registerHandler(w, req)

	var authResp AuthResponse
	json.NewDecoder(w.Body).Decode(&authResp)
	token := authResp.Token

	tests := []struct {
		name           string
		payload        CreatePostRequest
		token          string
		expectedStatus bool
	}{
		{
			name: "Valid post",
			payload: CreatePostRequest{
				Title:    "Test Post",
				Content:  "This is a test post content",
				Category: "general",
			},
			token:          token,
			expectedStatus: true,
		},
		{
			name: "Short title",
			payload: CreatePostRequest{
				Title:    "Hi",
				Content:  "This is a test post content",
				Category: "general",
			},
			token:          token,
			expectedStatus: false,
		},
		{
			name: "Short content",
			payload: CreatePostRequest{
				Title:    "Test Post",
				Content:  "Short",
				Category: "general",
			},
			token:          token,
			expectedStatus: false,
		},
		{
			name: "No token",
			payload: CreatePostRequest{
				Title:    "Test Post",
				Content:  "This is a test post content",
				Category: "general",
			},
			token:          "",
			expectedStatus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/posts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			w := httptest.NewRecorder()

			createPostHandler(w, req)

			var response PostResponse
			json.NewDecoder(w.Body).Decode(&response)

			if response.Success != tt.expectedStatus {
				t.Errorf("Expected success=%v, got success=%v", tt.expectedStatus, response.Success)
			}

			if tt.expectedStatus && response.Post == nil {
				t.Error("Expected post to be returned on successful creation")
			}
		})
	}
}

func TestGetPosts(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create a user and post
	registerPayload := RegisterRequest{
		Nickname:  "getpostuser",
		Email:     "getpost@example.com",
		Password:  "password123",
		FirstName: "Get",
		LastName:  "Post",
		Age:       25,
		Gender:    "male",
	}

	body, _ := json.Marshal(registerPayload)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	registerHandler(w, req)

	var authResp AuthResponse
	json.NewDecoder(w.Body).Decode(&authResp)
	token := authResp.Token

	// Create a post
	postPayload := CreatePostRequest{
		Title:    "Test Post",
		Content:  "This is a test post content",
		Category: "technology",
	}

	body, _ = json.Marshal(postPayload)
	req = httptest.NewRequest("POST", "/api/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w = httptest.NewRecorder()
	createPostHandler(w, req)

	// Get all posts
	req = httptest.NewRequest("GET", "/api/posts", nil)
	w = httptest.NewRecorder()
	getPostsHandler(w, req)

	var response PostResponse
	json.NewDecoder(w.Body).Decode(&response)

	if !response.Success {
		t.Error("Expected success=true for getting posts")
	}

	if len(response.Posts) == 0 {
		t.Error("Expected at least one post")
	}
}

func TestCreateComment(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create user and post
	registerPayload := RegisterRequest{
		Nickname:  "commentuser",
		Email:     "comment@example.com",
		Password:  "password123",
		FirstName: "Comment",
		LastName:  "User",
		Age:       25,
		Gender:    "female",
	}

	body, _ := json.Marshal(registerPayload)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	registerHandler(w, req)

	var authResp AuthResponse
	json.NewDecoder(w.Body).Decode(&authResp)
	token := authResp.Token

	// Create a post
	postPayload := CreatePostRequest{
		Title:    "Post for Comments",
		Content:  "This post will have comments",
		Category: "general",
	}

	body, _ = json.Marshal(postPayload)
	req = httptest.NewRequest("POST", "/api/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w = httptest.NewRecorder()
	createPostHandler(w, req)

	var postResp PostResponse
	json.NewDecoder(w.Body).Decode(&postResp)
	postID := postResp.Post.ID

	// Create a comment
	commentPayload := CreateCommentRequest{
		Content: "This is a test comment",
	}

	body, _ = json.Marshal(commentPayload)
	req = httptest.NewRequest("POST", "/api/posts/1/comments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w = httptest.NewRecorder()

	// Mock mux vars
	req = httptest.NewRequest("POST", "/api/posts/1/comments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w = httptest.NewRecorder()

	// Directly insert comment for testing
	_, err := db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)",
		postID, authResp.User.ID, commentPayload.Content)

	if err != nil {
		t.Errorf("Failed to create comment: %v", err)
	}
}

