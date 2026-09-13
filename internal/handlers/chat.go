package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go-ai-stream/internal/service"
	"go-ai-stream/internal/storage"
)

type AIStreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

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

func (h *ChatHandler) HandleAIStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ phương thức POST", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		sessionID = "default_session"
	}

	var reqBody struct {
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Dữ liệu JSON không hợp lệ", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Lưu tin nhắn của user vào session store (RAM hoặc Redis)
	if err := h.sessionStore.AddMessage(ctx, sessionID, "user", reqBody.Prompt); err != nil {
		log.Printf("[ERROR] Không thể lưu tin nhắn user: %v", err)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming không được hỗ trợ", http.StatusInternalServerError)
		return
	}

	// Lấy lịch sử hội thoại
	history, err := h.sessionStore.GetHistory(ctx, sessionID)
	if err != nil {
		history = []storage.Message{{Role: "user", Content: reqBody.Prompt}}
	}

	resp, providerName, err := h.aiService.CallAIWithFallback(ctx, history)
	if err != nil {
		log.Printf("[ERROR] Tất cả AI providers thất bại: %v", err)
		http.Error(w, "Hệ thống AI bận, vui lòng thử lại sau.", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	log.Printf("[INFO] Đang stream cho session [%s] dùng provider: %s", sessionID, providerName)

	reader := bufio.NewReader(resp.Body)
	var fullContent string

	for {
		select {
		case <-ctx.Done():
			log.Printf("[INFO] Client [%s] ngắt kết nối", sessionID)
			return
		default:
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		lineStr := string(line)
		if len(lineStr) > 5 && lineStr[:5] == "data:" {
			payload := lineStr[5:]
			if bytes.Equal(bytes.TrimSpace([]byte(payload)), []byte("[DONE]")) {
				break
			}

			var streamResp AIStreamResponse
			if err := json.Unmarshal([]byte(payload), &streamResp); err == nil {
				if len(streamResp.Choices) > 0 {
					content := streamResp.Choices[0].Delta.Content
					if content != "" {
						w.Write([]byte(content))
						flusher.Flush()
						fullContent += content
					}
				}
			}
		}
	}

	// Lưu phản hồi của AI vào session store
	if fullContent != "" {
		_ = h.sessionStore.AddMessage(ctx, sessionID, "assistant", fullContent)
	}
}
