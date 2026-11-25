package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"lume_api/internal/services"
)

func TestAuthErrorHTTPStatuses(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "invalid credentials", err: services.ErrInvalidCredentials, want: http.StatusUnauthorized},
		{name: "invalid input", err: services.ErrValidation, want: http.StatusBadRequest},
		{name: "duplicate account", err: services.ErrConflict, want: http.StatusConflict},
		{name: "invalid action token", err: services.ErrInvalidToken, want: http.StatusBadRequest},
		{name: "invalid refresh session", err: services.ErrInvalidSession, want: http.StatusUnauthorized},
		{name: "mail unavailable", err: services.ErrEmailUnavailable, want: http.StatusServiceUnavailable},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/", func(c *fiber.Ctx) error {
				return authError(c, test.err, "authentication failed")
			})
			response, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != test.want {
				t.Fatalf("status = %d, want %d", response.StatusCode, test.want)
			}
		})
	}
}
