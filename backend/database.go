package main

import (
	"fmt"
)

func initDatabase() error {
	// Create users table
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		age INTEGER NOT NULL,
		gender TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(usersTable); err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	// Create sessions table for session management
	sessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(sessionsTable); err != nil {
		return fmt.Errorf("failed to create sessions table: %v", err)
	}

	// Create posts table
	postsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		category TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(postsTable); err != nil {
		return fmt.Errorf("failed to create posts table: %v", err)
	}

	// Create comments table
	commentsTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(commentsTable); err != nil {
		return fmt.Errorf("failed to create comments table: %v", err)
	}

	// Create messages table for private messaging
	messagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		read_at DATETIME,
		FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(messagesTable); err != nil {
		return fmt.Errorf("failed to create messages table: %v", err)
	}

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_nickname ON users(nickname);",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);",
		"CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);",
		"CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_posts_category ON posts(category);",
		"CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);",
		"CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);",
		"CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);",
		"CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages(receiver_id);",
		"CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);",
	}

	for _, idx := range indexes {
		if _, err := db.Exec(idx); err != nil {
			return fmt.Errorf("failed to create index: %v", err)
		}
	}

	return nil
}
