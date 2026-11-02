package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"lume_api/internal/utils"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var e *fiber.Error
	if errors.As(err, &e) {
		return utils.Fail(c, e.Code, "fiber_error", e.Message, "")
	}
	return utils.Fail(c, fiber.StatusInternalServerError, "internal_error", "Erro interno do servidor", err.Error())
}
