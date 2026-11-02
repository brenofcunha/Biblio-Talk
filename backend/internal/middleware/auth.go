package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"lume_api/internal/utils"
)

func AuthRequired(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
			return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized", "Token JWT ausente", "envie Authorization: Bearer <token>")
		}
		token := strings.TrimSpace(header[7:])
		claims, err := utils.ParseToken(jwtSecret, token)
		if err != nil {
			return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized", "Token JWT inválido", err.Error())
		}
		c.Locals("userId", claims.UserID)
		c.Locals("username", claims.Username)
		return c.Next()
	}
}
