package logs

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// FiberLogger is a gofiber logger
type FiberLogger struct {
	logger *zap.Logger
}

// NewFiberLogger creates a new FiberLogger with zap.Logger
func NewFiberLogger(logger *zap.Logger) *FiberLogger {
	return &FiberLogger{logger: logger}
}

// Middleware is an fiber router middleware
func (f *FiberLogger) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		f.logger.Debug(
			"HTTP request",
			zap.Time("timestamp", time.Now()),
			zap.ByteString("http_method", c.Request().Header.Method()),
			zap.ByteString("http_uri", c.Request().Header.RequestURI()),
			zap.ByteString("http_addr", c.Request().Header.RequestURI()),
			zap.ByteString("http_scheme", c.Request().Header.Protocol()),
			zap.ByteString("http_agent", c.Request().Header.UserAgent()),
		)
		return c.Next()
	}
}
