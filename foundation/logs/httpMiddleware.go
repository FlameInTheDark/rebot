package logs

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LoggerKey is an key for the context.Context
const LoggerKey = "logger"

// HTTPLoggerMiddleware create logger middleware to http requests
func HTTPLoggerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			logger.Debug(
				"HTTP request",
				zap.Time("timestamp", time.Now()),
				zap.String("http_method", r.Method),
				zap.String("http_uri", r.RequestURI),
				zap.String("http_addr", r.RemoteAddr),
				zap.String("http_scheme", scheme),
				zap.String("http_agent", r.UserAgent()),
			)
			ctx := context.WithValue(r.Context(), LoggerKey, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
