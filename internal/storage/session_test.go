package storage_test

import (
	"context"
	"testing"

	"go-ai-stream/internal/storage"
)

func TestMemoryStore(t *testing.T) {
	store := storage.NewMemoryStore()
	ctx := context.Background()
	sessionID := "test-session-1"

	msg := storage.Message{
		Role:    "user",
		Content: "Hello AI Worker",
	}

	// Test thêm tin nhắn
	err := store.AddMessage(ctx, sessionID, msg)
	if err != nil {
		t.Fatalf("Thêm tin nhắn bị lỗi: %v", err)
	}

	// Test lấy danh sách tin nhắn
	messages, err := store.GetMessages(ctx, sessionID)
	if err != nil {
		t.Fatalf("Lấy tin nhắn bị lỗi: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("Mong đợi 1 tin nhắn, nhưng nhận được %d", len(messages))
	}

	if messages[0].Content != "Hello AI Worker" {
		t.Errorf("Nội dung tin nhắn không khớp. Mong đợi 'Hello AI Worker', nhận được '%s'", messages[0].Content)
	}
}
