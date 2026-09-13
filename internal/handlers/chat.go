package handlers

import (
	"encoding/json"
	"net/http"

	"go-ai-stream/internal/service"
	"go-ai-stream/internal/storage"
)

type ChatHandler struct {
	aiService    *service.AIService
	sessionStore storage.SessionStore
}

func NewChatHandler(aiService *service.AIService, sessionStore storage.SessionStore) *ChatHandler {
	return &ChatHandler{
		aiService:    aiService,
		sessionStore: sessionStore,
	}
}

type StreamRequest struct {
	Prompt string `json:"prompt"`
}

func (h *ChatHandler) HandleStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		sessionID = "default"
	}

	var req StreamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 1. Lưu tin nhắn của User dạng struct storage.Message
	userMsg := storage.Message{
		Role:    "user",
		Content: req.Prompt,
	}
	if err := h.sessionStore.AddMessage(ctx, sessionID, userMsg); err != nil {
		http.Error(w, "Failed to save user message", http.StatusInternalServerError)
		return
	}

	// 2. Lấy lịch sử bằng GetMessages
	_, err := h.sessionStore.GetMessages(ctx, sessionID)
	if err != nil {
		http.Error(w, "Failed to retrieve chat history", http.StatusInternalServerError)
		return
	}

	// Thiết lập Header cho Streaming Response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Gọi AI Worker và trả stream
	responseContent := "Phản hồi mẫu từ AI Service cho prompt: " + req.Prompt
	w.Write([]byte(responseContent))

	// 3. Lưu phản hồi của AI dạng struct storage.Message
	aiMsg := storage.Message{
		Role:    "assistant",
		Content: responseContent,
	}
	_ = h.sessionStore.AddMessage(ctx, sessionID, aiMsg)
}
