package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthOptionalAllowsAnonymousRequest(t *testing.T) {
	app := fiber.New()
	app.Get("/", AuthOptional("secret", func(context.Context, string) (bool, error) {
		t.Fatal("session validator called without credentials")
		return false, nil
	}), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	if err != nil {
		t.Fatalf("app.Test error = %v", err)
	}
	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusNoContent)
	}
}

func TestAuthOptionalRejectsInvalidProvidedToken(t *testing.T) {
	app := fiber.New()
	app.Get("/", AuthOptional("secret", func(context.Context, string) (bool, error) {
		t.Fatal("session validator called with invalid token")
		return false, nil
	}), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer invalid")
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("app.Test error = %v", err)
	}
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}
}
