package repository

import (
	"testing"
	"time"
)

func TestMessageRepository_CreateMessage(t *testing.T) {
	repos := SetupTestDB(t)
	msgRepo := repos.Message
	userRepo := repos.User

	userID1, _ := userRepo.CreateUser("sender", "sender@example.com", "hashedpass1", "User", "One", 20, "male")
	userID2, _ := userRepo.CreateUser("receiver", "receiver@example.com", "hashedpass2", "User", "Two", 21, "female")

	id, err := msgRepo.CreateMessage(int(userID1), int(userID2), "hello world")
	if err != nil {
		t.Fatalf("unexpected error creating message: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero message ID")
	}
}

func TestMessageRepository_GetConversation(t *testing.T) {
	repos := SetupTestDB(t)
	msgRepo := repos.Message
	userRepo := repos.User

	userID1, _ := userRepo.CreateUser("sender", "sender@example.com", "hashedpass1", "User", "One", 20, "male")
	userID2, _ := userRepo.CreateUser("receiver", "receiver@example.com", "hashedpass2", "User", "Two", 21, "female")

	conversation, err := msgRepo.GetConversation(int(userID1), int(userID2), 10, 0)
	if err != nil {
		t.Fatalf("unexpected error getting empty conversation: %v", err)
	}
	if len(conversation) != 0 {
		t.Errorf("expected 0 messages, got %d", len(conversation))
	}

	msgRepo.CreateMessage(int(userID1), int(userID2), "first")
	time.Sleep(1 * time.Millisecond)
	msgRepo.CreateMessage(int(userID2), int(userID1), "second")

	conversation, err = msgRepo.GetConversation(int(userID1), int(userID2), 10, 0)
	if err != nil {
		t.Fatalf("failed to get conversation: %v", err)
	}
	if len(conversation) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(conversation))
	}
	if conversation[0].Content != "first" || conversation[1].Content != "second" {
		t.Errorf("expected chronological order, got %v", conversation)
	}

	conversation, _ = msgRepo.GetConversation(int(userID1), int(userID2), 1, 0)
	if len(conversation) != 1 {
		t.Errorf("limit did not apply: expected 1, got %d", len(conversation))
	}
}

func TestMessageRepository_MarkAndCountUnread(t *testing.T) {
	repos := SetupTestDB(t)
	msgRepo := repos.Message
	userRepo := repos.User

	userID1, _ := userRepo.CreateUser("carl", "c@example.com", "pass", "Carl", "C", 22, "male")
	userID2, _ := userRepo.CreateUser("dana", "d@example.com", "pass", "Dana", "D", 23, "female")

	msgRepo.CreateMessage(int(userID2), int(userID1), "msg1")
	msgRepo.CreateMessage(int(userID2), int(userID1), "msg2")

	count, err := msgRepo.GetUnreadMessagesCount(int(userID1), int(userID2))
	if err != nil {
		t.Fatalf("error counting unread: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 unread, got %d", count)
	}

	if err := msgRepo.MarkMessageAsRead(int(userID1), int(userID2)); err != nil {
		t.Fatalf("failed to mark as read: %v", err)
	}
	count, _ = msgRepo.GetUnreadMessagesCount(int(userID1), int(userID2))
	if count != 0 {
		t.Errorf("expected 0 unread after mark, got %d", count)
	}
}

func TestMessageRepository_PartnersAndLastMessage(t *testing.T) {
	repos := SetupTestDB(t)
	msgRepo := repos.Message
	userRepo := repos.User

	userID1, _ := userRepo.CreateUser("e1", "e1@example.com", "pass", "E", "One", 40, "male")
	userID2, _ := userRepo.CreateUser("e2", "e2@example.com", "pass", "E", "Two", 41, "female")
	userID3, _ := userRepo.CreateUser("e3", "e3@example.com", "pass", "E", "Three", 42, "other")

	msgRepo.CreateMessage(int(userID1), int(userID2), "hey")
	msgRepo.CreateMessage(int(userID3), int(userID1), "yo")

	partners, err := msgRepo.GetConversationPartners(int(userID1))
	if err != nil {
		t.Fatalf("error getting partners: %v", err)
	}
	if len(partners) != 2 {
		t.Errorf("expected 2 partners, got %d", len(partners))
	}

	last, err := msgRepo.GetLastMessageBetweenUsers(int(userID1), int(userID2))
	if err != nil {
		t.Fatalf("error getting last message: %v", err)
	}
	if last == nil || last.Content != "hey" {
		t.Errorf("unexpected last message: %+v", last)
	}

	_, err = msgRepo.GetLastMessageBetweenUsers(int(userID2), int(userID3))
	if err == nil {
		t.Error("expected error when no messages exist")
	}
}
