package service

import (
	"testing"
	"go-ai-stream/internal/config"
)

func TestNewAIService(t *testing.T) {
	cfg := &config.Config{
		AIAPIURL: "https://api.test.com",
		AIAPIKey: "test-key",
	}
	svc := NewAIService(cfg)
	if svc == nil {
		t.Fatal("Lỗi: Không thể khởi tạo AIService")
	}
}
