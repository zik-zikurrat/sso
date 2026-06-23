package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ctxKey string

const TraceIDKey ctxKey = "trace_id"

func TracingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.NewString()
		}
		c.Locals(TraceIDKey, traceID)
		c.Response().Header.Set("X-Trace-ID", traceID)
		ctx := context.WithValue(c.UserContext(), TraceIDKey, traceID)
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func GetTraceID(c *fiber.Ctx) string {
	v := c.Locals(TraceIDKey)
	if v == nil {
		return ""
	}
	return v.(string)
}
