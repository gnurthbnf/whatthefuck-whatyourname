package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-ai-stream/internal/config"
	"go-ai-stream/internal/storage"
)

type AIService struct {
	cfg *config.Config
}

func NewAIService(cfg *config.Config) *AIService {
	return &AIService{cfg: cfg}
}

type AIRequest struct {
	Model    string            `json:"model"`
	Messages []storage.Message `json:"messages"`
	Stream   bool              `json:"stream"`
}

func (s *AIService) CallAIWithFallback(ctx context.Context, messages []storage.Message) (*http.Response, string, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	aiReq := AIRequest{
		Model:    "llama-3.3-70b-versatile",
		Messages: messages,
		Stream:   true,
	}
	jsonData, err := json.Marshal(aiReq)
	if err != nil {
		return nil, "", err
	}

	// 1. Thử Provider chính (Groq)
	if s.cfg.AIAPIKey != "" {
		req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.AIAPIURL, bytes.NewBuffer(jsonData))
		if err == nil {
			req.Header.Set("Authorization", "Bearer "+s.cfg.AIAPIKey)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err == nil && resp.StatusCode == http.StatusOK {
				return resp, "Primary-Groq", nil
			}
			if resp != nil {
				resp.Body.Close()
			}
			log.Println("[WARN] AI chính lỗi, chuyển sang Fallback...")
		}
	}

	// 2. Thử Provider dự phòng (Fallback)
	if s.cfg.AIFallbackURL != "" && s.cfg.AIFallbackKey != "" {
		req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.AIFallbackURL, bytes.NewBuffer(jsonData))
		if err == nil {
			req.Header.Set("Authorization", "Bearer "+s.cfg.AIFallbackKey)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err == nil && resp.StatusCode == http.StatusOK {
				return resp, "Fallback-Provider", nil
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}

	return nil, "None", fmt.Errorf("tất cả các AI provider đều thất bại")
}
