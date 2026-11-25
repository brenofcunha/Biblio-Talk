package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"lume_api/internal/utils"
)

func AuthRequired(jwtSecret string, validateSession func(ctx context.Context, sessionID string) (bool, error)) fiber.Handler {
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
		if claims.SessionID == "" {
			return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized", "Sessão inválida", "faça login novamente")
		}
		active, err := validateSession(c.UserContext(), claims.SessionID)
		if err != nil {
			return utils.Fail(c, fiber.StatusInternalServerError, "session_check_failed", "Não foi possível validar a sessão", "")
		}
		if !active {
			return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized", "Sessão revogada ou expirada", "faça login novamente")
		}
		c.Locals("userId", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("sessionId", claims.SessionID)
		return c.Next()
	}
}

func AuthOptional(jwtSecret string, validateSession func(ctx context.Context, sessionID string) (bool, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Next()
		}
		if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
			return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized", "Token JWT inválido", "envie Authorization: Bearer <token>")
		}
		claims, err := utils.ParseToken(jwtSecret, strings.TrimSpace(header[7:]))
		if err != nil {
			return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized", "Token JWT inválido", err.Error())
		}
		if claims.SessionID == "" {
			return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized", "Sessão inválida", "faça login novamente")
		}
		active, err := validateSession(c.UserContext(), claims.SessionID)
		if err != nil {
			return utils.Fail(c, fiber.StatusInternalServerError, "session_check_failed", "Não foi possível validar a sessão", "")
		}
		if !active {
			return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized", "Sessão revogada ou expirada", "faça login novamente")
		}
		c.Locals("userId", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("sessionId", claims.SessionID)
		return c.Next()
	}
}
