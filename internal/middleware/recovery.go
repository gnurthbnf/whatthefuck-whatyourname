package middleware

import (
	"log"
	"net/http"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("CRITICAL PANIC RECOVERED: %v", err)
				http.Error(w, "Lỗi hệ thống nội bộ, vui lòng thử lại sau!", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
