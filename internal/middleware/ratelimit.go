package middleware

import "net/http"

// RateLimiter bọc http.Handler để kiểm soát tần suất request
func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logic rate limit xử lý tại đây
		next.ServeHTTP(w, r)
	})
}
