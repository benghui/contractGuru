package middleware

import (
	"net/http"
	"time"

	"github.com/contractGuru/pkg/logger"
)

// LoggingMiddleware logs request methods & path & duration
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		logger.Info.Printf("%s - %s %s %v", r.Proto, r.Method, r.URL.Path, time.Since(start))
	})
}
