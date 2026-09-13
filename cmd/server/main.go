package main

import (
	"log"
	"net/http"
	"time"

	"go-ai-stream/internal/config"
	"go-ai-stream/internal/handlers"
	"go-ai-stream/internal/middleware"
	"go-ai-stream/internal/service"
	"go-ai-stream/internal/storage"
)

func main() {
	// 1. Tải cấu hình từ biến môi trường
	cfg := config.Load()

	// 2. Khởi tạo Session Store (Ưu tiên Redis nếu có cấu hình, ngược lại dùng RAM)
	var sessionStore storage.SessionStore
	if cfg.RedisURL != "" {
		rStore, err := storage.NewRedisStore(cfg.RedisURL)
		if err != nil {
			log.Printf("[WARN] Không thể kết nối Redis (%v), fallback về In-Memory RAM.", err)
			sessionStore = storage.NewMemoryStore()
		} else {
			log.Println("[INFO] Đã kết nối thành công tới Redis Cluster/Server.")
			sessionStore = rStore
		}
	} else {
		log.Println("[INFO] Sử dụng bộ nhớ In-Memory RAM cho Session.")
		sessionStore = storage.NewMemoryStore()
	}

	// 3. Khởi tạo Services & Handlers
	aiService := service.NewAIService(cfg)
	chatHandler := handlers.NewChatHandler(aiService, sessionStore)

	// 4. Khởi tạo Middleware chống Spam
	rateLimiter := middleware.NewRateLimiter(1000 * time.Millisecond)

	// 5. Định nghĩa Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HandleHome)
	mux.HandleFunc("/api/stream", rateLimiter.Limit(chatHandler.HandleAIStream))

	log.Printf("[INFO] Server Go AI Stream đang khởi động tại cổng %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("[FATAL] Server gặp sự cố: %v", err)
	}
}
