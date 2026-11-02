package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SlogRequest(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)
		rid, _ := c.Locals("requestId").(string)
		logger.Info("requisição", "requestId", rid, "method", c.Method(), "path", c.Path(), "status", c.Response().StatusCode(), "latencyMs", latency.Milliseconds())
		return err
	}
}
