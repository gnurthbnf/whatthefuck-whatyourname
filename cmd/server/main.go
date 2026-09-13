package main

import (
	"log"
	"net/http"

	"go-ai-stream/internal/middleware"
	"go-ai-stream/internal/storage"
)

func main() {
	// Khởi tạo bộ lưu trữ RAM nội bộ
	store := storage.NewMemoryStore()

	// Đăng ký các route xử lý HTTP
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "running", "message": "Go AI Stream Server is up!"}`))
	})

	// Áp dụng middleware chống spam (Rate limiter)
	handler := middleware.RateLimiter(mux)

	log.Println("Server đang khởi động tại cổng :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Lỗi khởi động server: %v", err)
	}
}
