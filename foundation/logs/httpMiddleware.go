package logs

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func HttpLoggerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			scheme := "http"
			if request.TLS != nil {
				scheme = "https"
			}
			logger.Debug(
				"HTTP request",
				zap.Time("timestamp", time.Now()),
				zap.String("http_method", request.Method),
				zap.String("http_uri", request.RequestURI),
				zap.String("http_addr", request.RemoteAddr),
				zap.String("http_scheme", scheme),
				zap.String("http_agent", request.UserAgent()),
			)
			next.ServeHTTP(writer, request)
		})
	}
}
