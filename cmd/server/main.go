package main

import (
	"log"
	"net/http"

	"go-ai-stream/internal/client"
	"go-ai-stream/internal/middleware"
	"go-ai-stream/internal/storage"
)

func main() {
	// Khởi tạo storage và client
	store := storage.NewMemoryStore()
	serviceClient := client.NewServiceClient()

	mux := http.NewServeMux()

	// Khai báo các route API
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/api/ai", func(w http.ResponseWriter, r *http.Request) {
		// Sử dụng store để lưu hoặc lấy dữ liệu session
		_ = store.AddMessage(r.Context(), "default", storage.Message{Role: "user", Content: "ping"})

		if err := serviceClient.CallAIWorker(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("AI Worker connection successful"))
	})

	// Bọc mux bằng middleware RateLimiter và Recovery
	var handler http.Handler = mux
	handler = middleware.RateLimiter(handler)
	handler = middleware.Recovery(handler)

	log.Println("Server đang chạy tại cổng :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Lỗi khởi động server: %v", err)
	}
}
