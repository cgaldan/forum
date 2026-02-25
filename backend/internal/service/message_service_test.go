package service

import (
	"testing"
)

func TestMessageService_SendMessage(t *testing.T) {
	services := SetupTestServices(t)

	senderID := CreateTestUser(t, services, "sender", "sender@example.com", "password123", "John", "Doe", 25, "male")
	receiverID := CreateTestUser(t, services, "receiver", "receiver@example.com", "password123", "Jane", "Smith", 30, "female")

	tests := []struct {
		name        string
		senderID    int
		receiverID  int
		content     string
		expectError bool
	}{
		{
			name:        "valid message",
			senderID:    senderID,
			receiverID:  receiverID,
			content:     "Hello, this is a test message!",
			expectError: false,
		},
		{
			name:        "empty content",
			senderID:    senderID,
			receiverID:  receiverID,
			content:     "",
			expectError: true,
		},
		{
			name:        "message too long",
			senderID:    senderID,
			receiverID:  receiverID,
			content:     string(make([]byte, 1001)),
			expectError: true,
		},
		{
			name:        "send to self",
			senderID:    senderID,
			receiverID:  senderID,
			content:     "Message to myself",
			expectError: true,
		},
		{
			name:        "non-existent receiver",
			senderID:    senderID,
			receiverID:  99999,
			content:     "Message to ghost",
			expectError: true,
		},
		{
			name:        "non-existent sender",
			senderID:    99999,
			receiverID:  receiverID,
			content:     "Message from ghost",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := services.Message.SendMessage(tt.senderID, tt.receiverID, tt.content)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if message == nil {
				t.Fatal("Expected message but got nil")
			}

			if message.SenderID != tt.senderID {
				t.Errorf("Expected sender ID %d, got %d", tt.senderID, message.SenderID)
			}

			if message.ReceiverID != tt.receiverID {
				t.Errorf("Expected receiver ID %d, got %d", tt.receiverID, message.ReceiverID)
			}

			if message.Content != tt.content {
				t.Errorf("Expected content %s, got %s", tt.content, message.Content)
			}

			if message.SenderName != "sender" {
				t.Errorf("Expected sender name 'sender', got '%s'", message.SenderName)
			}
		})
	}
}

func TestMessageService_GetConversation(t *testing.T) {
	services := SetupTestServices(t)

	user1ID := CreateTestUser(t, services, "user1", "user1@example.com", "password123", "John", "Doe", 25, "male")
	user2ID := CreateTestUser(t, services, "user2", "user2@example.com", "password123", "Jane", "Smith", 30, "female")

	CreateTestMessage(t, services, user1ID, user2ID, "Hello from user1")
	CreateTestMessage(t, services, user2ID, user1ID, "Hello from user2")
	CreateTestMessage(t, services, user1ID, user2ID, "How are you?")

	t.Run("get conversation from user1 perspective", func(t *testing.T) {
		messages, err := services.Message.GetConversation(user1ID, user2ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get conversation: %v", err)
		}

		if len(messages) != 3 {
			t.Errorf("Expected 3 messages, got %d", len(messages))
		}

		messageContents := make(map[string]bool)
		for _, msg := range messages {
			messageContents[msg.Content] = true
		}

		expectedContents := []string{"Hello from user1", "Hello from user2", "How are you?"}
		for _, content := range expectedContents {
			if !messageContents[content] {
				t.Errorf("Expected message with content '%s' not found", content)
			}
		}
	})

	t.Run("get conversation from user2 perspective", func(t *testing.T) {
		messages, err := services.Message.GetConversation(user2ID, user1ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get conversation: %v", err)
		}

		if len(messages) != 3 {
			t.Errorf("Expected 3 messages, got %d", len(messages))
		}
	})

	t.Run("get conversation with pagination", func(t *testing.T) {
		messages, err := services.Message.GetConversation(user1ID, user2ID, 2, 0)
		if err != nil {
			t.Fatalf("Failed to get conversation with pagination: %v", err)
		}

		if len(messages) != 2 {
			t.Errorf("Expected 2 messages with limit 2, got %d", len(messages))
		}
	})

	t.Run("get conversation between users with no messages", func(t *testing.T) {
		user3ID := CreateTestUser(t, services, "user3", "user3@example.com", "password123", "Bob", "Wilson", 35, "male")

		messages, err := services.Message.GetConversation(user1ID, user3ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to get empty conversation: %v", err)
		}

		if len(messages) != 0 {
			t.Errorf("Expected 0 messages for users with no conversation, got %d", len(messages))
		}
	})
}

func TestMessageService_GetConversationsByUserID(t *testing.T) {
	services := SetupTestServices(t)

	user1ID := CreateTestUser(t, services, "user1", "user1@example.com", "password123", "John", "Doe", 25, "male")
	user2ID := CreateTestUser(t, services, "user2", "user2@example.com", "password123", "Jane", "Smith", 30, "female")
	user3ID := CreateTestUser(t, services, "user3", "user3@example.com", "password123", "Bob", "Wilson", 35, "male")

	CreateTestMessage(t, services, user1ID, user2ID, "Hello user2")
	CreateTestMessage(t, services, user2ID, user1ID, "Hello user1")
	CreateTestMessage(t, services, user1ID, user3ID, "Hello user3")
	CreateTestMessage(t, services, user3ID, user1ID, "Hello user1 back")

	t.Run("get conversations for user1", func(t *testing.T) {
		conversations, err := services.Message.GetConversationsByUserID(user1ID)
		if err != nil {
			t.Fatalf("Failed to get conversations: %v", err)
		}

		if len(conversations) != 2 {
			t.Errorf("Expected 2 conversations for user1, got %d", len(conversations))
		}

		userIDs := make(map[int]bool)
		for _, conv := range conversations {
			userIDs[conv.UserID] = true
			if conv.UserID == user2ID && conv.Nickname != "user2" {
				t.Errorf("Expected nickname 'user2', got '%s'", conv.Nickname)
			}
			if conv.UserID == user3ID && conv.Nickname != "user3" {
				t.Errorf("Expected nickname 'user3', got '%s'", conv.Nickname)
			}
		}

		if !userIDs[user2ID] || !userIDs[user3ID] {
			t.Error("Expected conversations with both user2 and user3")
		}
	})

	t.Run("get conversations for user with no messages", func(t *testing.T) {
		user4ID := CreateTestUser(t, services, "user4", "user4@example.com", "password123", "Alice", "Brown", 28, "female")

		conversations, err := services.Message.GetConversationsByUserID(user4ID)
		if err != nil {
			t.Fatalf("Failed to get conversations for user with no messages: %v", err)
		}

		if len(conversations) != 0 {
			t.Errorf("Expected 0 conversations for user with no messages, got %d", len(conversations))
		}
	})
}
