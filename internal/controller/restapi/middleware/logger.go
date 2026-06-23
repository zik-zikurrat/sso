package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware(l *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		traceID, _ := c.Locals(TraceIDKey).(string)

		l.Info("http request",
			slog.String("trace_id", traceID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", time.Since(start)),
		)

		return err
	}
}
