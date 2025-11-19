package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    string    `json:"author"` // nickname of the author
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Author    string    `json:"author"` // nickname of the author
}

type CreatePostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

type PostResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Post    *Post    `json:"post,omitempty"`
	Posts   []Post   `json:"posts,omitempty"`
	Comment *Comment `json:"comment,omitempty"`
}

type PostDetailResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
	Post     *Post     `json:"post,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
}

// Get all posts with optional category filter and pagination
func getPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters
	category := r.URL.Query().Get("category")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Build query
	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.category, p.created_at, p.updated_at, u.nickname
		FROM posts p
		JOIN users u ON p.user_id = u.id`

	args := []interface{}{}
	if category != "" {
		query += " WHERE p.category = ?"
		args = append(args, category)
	}

	query += " ORDER BY p.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Failed to fetch posts"})
		return
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category,
			&post.CreatedAt, &post.UpdatedAt, &post.Author)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(PostResponse{Success: true, Posts: posts})
}

// Get a single post with its comments
func getPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	postID := vars["id"]

	// Get post
	var post Post
	err := db.QueryRow(`
		SELECT p.id, p.user_id, p.title, p.content, p.category, p.created_at, p.updated_at, u.nickname
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?`, postID).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category,
		&post.CreatedAt, &post.UpdatedAt, &post.Author)

	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(PostDetailResponse{Success: false, Message: "Post not found"})
		return
	} else if err != nil {
		json.NewEncoder(w).Encode(PostDetailResponse{Success: false, Message: "Database error"})
		return
	}

	// Get comments for this post
	rows, err := db.Query(`
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.nickname
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC`, postID)

	comments := []Comment{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var comment Comment
			err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content,
				&comment.CreatedAt, &comment.Author)
			if err == nil {
				comments = append(comments, comment)
			}
		}
	}

	json.NewEncoder(w).Encode(PostDetailResponse{
		Success:  true,
		Post:     &post,
		Comments: comments,
	})
}

// Create a new post
func createPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user from session
	token := r.Header.Get("Authorization")
	user, err := getUserFromSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Unauthorized"})
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	// Validate
	if req.Title == "" || len(req.Title) < 3 {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Title must be at least 3 characters"})
		return
	}
	if req.Content == "" || len(req.Content) < 10 {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Content must be at least 10 characters"})
		return
	}
	if req.Category == "" {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Category is required"})
		return
	}

	// Insert post
	result, err := db.Exec(`
		INSERT INTO posts (user_id, title, content, category)
		VALUES (?, ?, ?, ?)`, user.ID, req.Title, req.Content, req.Category)

	if err != nil {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Failed to create post"})
		return
	}

	postID, _ := result.LastInsertId()

	post := &Post{
		ID:       int(postID),
		UserID:   user.ID,
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Author:   user.Nickname,
	}

	json.NewEncoder(w).Encode(PostResponse{
		Success: true,
		Message: "Post created successfully",
		Post:    post,
	})
}

// Create a comment on a post
func createCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user from session
	token := r.Header.Get("Authorization")
	user, err := getUserFromSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	postID := vars["id"]

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Invalid JSON"})
		return
	}

	// Validate
	if req.Content == "" || len(req.Content) < 1 {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Comment cannot be empty"})
		return
	}

	// Check if post exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
	if err != nil || !exists {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Post not found"})
		return
	}

	// Insert comment
	result, err := db.Exec(`
		INSERT INTO comments (post_id, user_id, content)
		VALUES (?, ?, ?)`, postID, user.ID, req.Content)

	if err != nil {
		json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Failed to create comment"})
		return
	}

	commentID, _ := result.LastInsertId()
	postIDInt, _ := strconv.Atoi(postID)

	comment := &Comment{
		ID:      int(commentID),
		PostID:  postIDInt,
		UserID:  user.ID,
		Content: req.Content,
		Author:  user.Nickname,
	}

	json.NewEncoder(w).Encode(PostResponse{
		Success: true,
		Message: "Comment created successfully",
		Comment: comment,
	})
}

