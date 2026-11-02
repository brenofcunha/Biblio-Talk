package utils

import (
	"github.com/gofiber/fiber/v2"
	"lume_api/internal/dto"
)

func Success(c *fiber.Ctx, status int, data interface{}, message string) error {
	return c.Status(status).JSON(dto.APIResponse{Success: true, Message: message, Data: data})
}

func Fail(c *fiber.Ctx, status int, code, message, details string) error {
	return c.Status(status).JSON(dto.APIResponse{
		Success: false,
		Message: message,
		Error: dto.ErrorDetail{Code: code, Details: details},
	})
}
