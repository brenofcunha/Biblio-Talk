package services

import (
	"context"
	"errors"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	for _, username := range []string{"abc", "reader_1", "123"} {
		if err := ValidateUsername(username); err != nil {
			t.Errorf("ValidateUsername(%q) returned an error: %v", username, err)
		}
	}
	for _, username := range []string{"ab", "with space", "reader.name", "a-very-long-username-that-is-too-long"} {
		if err := ValidateUsername(username); !errors.Is(err, ErrValidation) {
			t.Errorf("ValidateUsername(%q) error = %v, want ErrValidation", username, err)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	for _, email := range []string{"reader@example.com", "first.last+tag@example.org"} {
		if err := ValidateEmail(email); err != nil {
			t.Errorf("ValidateEmail(%q) returned an error: %v", email, err)
		}
	}
	for _, email := range []string{"", "not-an-email", "Reader <reader@example.com>"} {
		if err := ValidateEmail(email); !errors.Is(err, ErrValidation) {
			t.Errorf("ValidateEmail(%q) error = %v, want ErrValidation", email, err)
		}
	}
}

func TestValidateNewPassword(t *testing.T) {
	for _, password := range []string{"12345678", "senha-com-acentuação"} {
		if err := ValidateNewPassword(password); err != nil {
			t.Errorf("ValidateNewPassword(%q) returned an error: %v", password, err)
		}
	}
	for _, password := range []string{"short", string(make([]byte, 73))} {
		if err := ValidateNewPassword(password); !errors.Is(err, ErrValidation) {
			t.Errorf("ValidateNewPassword length %d error = %v, want ErrValidation", len(password), err)
		}
	}
}

func TestRevokeUserSessionRejectsInvalidID(t *testing.T) {
	service := &AuthService{}
	err := service.RevokeUserSession(context.Background(), "user-id", "not-a-uuid")
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("RevokeUserSession error = %v, want ErrValidation", err)
	}
}
