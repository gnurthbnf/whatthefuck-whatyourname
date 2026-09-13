package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-ai-stream/internal/middleware"
)

func TestRateLimiter(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handlerToTest := middleware.RateLimiter(nextHandler)

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	rec := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Mong đợi status code %d, nhưng nhận được %d", http.StatusOK, rec.Code)
	}
}

func TestRecovery(t *testing.T) {
	// Giả lập một handler bị panic
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Lỗi giả lập!")
	})

	handlerToTest := middleware.Recovery(panicHandler)

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	rec := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rec, req)

	// Middleware Recovery phải bắt được panic và trả về 500 Internal Server Error
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Mong đợi status code %d khi xảy ra panic, nhưng nhận được %d", http.StatusInternalServerError, rec.Code)
	}
}
