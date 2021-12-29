package logs

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type FiberLogger struct {
	logger *zap.Logger
}

func NewFiberLogger(logger *zap.Logger) *FiberLogger {
	return &FiberLogger{logger: logger}
}

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
