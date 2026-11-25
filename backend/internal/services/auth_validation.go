package services

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

var usernamePattern = regexp.MustCompile(`^[a-z0-9_]{3,30}$`)

func ValidateUsername(username string) error {
	if !usernamePattern.MatchString(username) {
		return fmt.Errorf("%w: o nome de usuário deve ter de 3 a 30 caracteres (letras, números e _)", ErrValidation)
	}
	return nil
}

func ValidateEmail(email string) error {
	if len(email) > 254 {
		return fmt.Errorf("%w: e-mail inválido", ErrValidation)
	}
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email || strings.TrimSpace(email) == "" {
		return fmt.Errorf("%w: e-mail inválido", ErrValidation)
	}
	return nil
}

func ValidateNewPassword(password string) error {
	if !utf8.ValidString(password) || utf8.RuneCountInString(password) < 8 || len(password) > 72 {
		return fmt.Errorf("%w: a senha deve ter entre 8 e 72 bytes", ErrValidation)
	}
	return nil
}
