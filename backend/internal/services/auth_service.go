package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"lume_api/internal/models"
	"lume_api/internal/repositories"
	"lume_api/internal/utils"
)

var (
	ErrValidation         = errors.New("dados de autenticação inválidos")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrConflict           = errors.New("conta já existente")
	ErrEmailNotVerified   = errors.New("e-mail ainda não confirmado")
	ErrInvalidToken       = errors.New("token inválido ou expirado")
	ErrInvalidSession     = errors.New("sessão inválida ou expirada")
	ErrEmailUnavailable   = errors.New("serviço de e-mail indisponível")
)

const (
	emailVerificationPurpose = "verify_email"
	passwordResetPurpose     = "reset_password"
	emailTokenLifetime       = time.Hour
)

type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
	Ready() bool
}

type AuthSession struct {
	User         *models.User
	AccessToken  string
	RefreshToken string
	SessionID    string
}

type AuthService struct {
	users           *repositories.UserRepository
	authRepo        *repositories.AuthRepository
	mailer          EmailSender
	jwtSecret       string
	jwtDuration     time.Duration
	refreshDuration time.Duration
	frontendURL     string
}

func NewAuthService(
	users *repositories.UserRepository,
	authRepo *repositories.AuthRepository,
	mailer EmailSender,
	jwtSecret string,
	jwtDuration, refreshDuration time.Duration,
	frontendURL string,
) *AuthService {
	return &AuthService{
		users:           users,
		authRepo:        authRepo,
		mailer:          mailer,
		jwtSecret:       jwtSecret,
		jwtDuration:     jwtDuration,
		refreshDuration: refreshDuration,
		frontendURL:     strings.TrimRight(frontendURL, "/"),
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	email = normalizeEmail(email)
	if err := ValidateUsername(username); err != nil {
		return nil, err
	}
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := ValidateNewPassword(password); err != nil {
		return nil, err
	}
	if s.mailer == nil || !s.mailer.Ready() {
		return nil, ErrEmailUnavailable
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user, err := s.users.Create(ctx, username, email, hash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("%w: nome de usuário ou e-mail já cadastrado", ErrConflict)
		}
		return nil, err
	}
	if err := s.sendActionToken(ctx, user, emailVerificationPurpose, "verifyToken", "Confirme seu e-mail", "Confirme seu endereço de e-mail neste link:"); err != nil {
		return user, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthSession, error) {
	email = normalizeEmail(email)
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	if strings.TrimSpace(password) == "" || len(password) > 72 {
		return nil, fmt.Errorf("%w: senha inválida", ErrValidation)
	}
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := utils.CheckPassword(user.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
	}
	verified, err := s.authRepo.IsEmailVerified(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if !verified {
		return nil, ErrEmailNotVerified
	}
	return s.createSession(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*AuthSession, error) {
	if refreshToken == "" {
		return nil, ErrInvalidSession
	}
	newRefreshToken, err := randomToken()
	if err != nil {
		return nil, err
	}
	sessionID, userID, err := s.authRepo.RotateSession(ctx, hashToken(refreshToken), hashToken(newRefreshToken), time.Now().Add(s.refreshDuration))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidSession
		}
		return nil, err
	}
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	accessToken, err := utils.GenerateToken(s.jwtSecret, user.ID, user.Username, sessionID, s.jwtDuration)
	if err != nil {
		return nil, err
	}
	return &AuthSession{User: user, AccessToken: accessToken, RefreshToken: newRefreshToken, SessionID: sessionID}, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, sessionID string) error {
	return s.authRepo.RevokeSession(ctx, sessionID)
}

func (s *AuthService) ListActiveSessions(ctx context.Context, userID, currentSessionID string) ([]models.AuthSessionInfo, error) {
	sessions, err := s.authRepo.ListActiveSessions(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i := range sessions {
		sessions[i].IsCurrent = sessions[i].ID == currentSessionID
	}
	return sessions, nil
}

func (s *AuthService) RevokeUserSession(ctx context.Context, userID, sessionID string) error {
	if _, err := uuid.Parse(sessionID); err != nil {
		return fmt.Errorf("%w: identificador de sessão inválido", ErrValidation)
	}
	return s.authRepo.RevokeUserSession(ctx, userID, sessionID)
}

func (s *AuthService) RevokeOtherSessions(ctx context.Context, userID, currentSessionID string) (int64, error) {
	return s.authRepo.RevokeOtherSessions(ctx, userID, currentSessionID)
}

func (s *AuthService) ValidateSession(ctx context.Context, sessionID string) (bool, error) {
	if sessionID == "" {
		return false, nil
	}
	return s.authRepo.IsSessionActive(ctx, sessionID)
}

func (s *AuthService) RequestEmailVerification(ctx context.Context, email string) error {
	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	verified, err := s.authRepo.IsEmailVerified(ctx, user.ID)
	if err != nil || verified {
		return err
	}
	return s.sendActionToken(ctx, user, emailVerificationPurpose, "verifyToken", "Confirme seu e-mail", "Confirme seu endereço de e-mail neste link:")
}

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidToken
	}
	err := s.authRepo.MarkEmailVerified(ctx, hashToken(token))
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidToken
	}
	return err
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	return s.sendActionToken(ctx, user, passwordResetPurpose, "resetToken", "Redefina sua senha", "Redefina sua senha neste link:")
}

func (s *AuthService) ResetPassword(ctx context.Context, token, password string) error {
	if token == "" {
		return ErrInvalidToken
	}
	if err := ValidateNewPassword(password); err != nil {
		return err
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	err = s.authRepo.ResetPassword(ctx, hashToken(token), hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidToken
	}
	return err
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, currentSessionID, currentPassword, newPassword string) error {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := utils.CheckPassword(user.PasswordHash, currentPassword); err != nil {
		return ErrInvalidCredentials
	}
	if err := ValidateNewPassword(newPassword); err != nil {
		return err
	}
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.authRepo.ChangePassword(ctx, userID, hash, currentSessionID)
}

func (s *AuthService) createSession(ctx context.Context, user *models.User) (*AuthSession, error) {
	refreshToken, err := randomToken()
	if err != nil {
		return nil, err
	}
	sessionID, err := s.authRepo.CreateSession(ctx, user.ID, hashToken(refreshToken), time.Now().Add(s.refreshDuration))
	if err != nil {
		return nil, err
	}
	accessToken, err := utils.GenerateToken(s.jwtSecret, user.ID, user.Username, sessionID, s.jwtDuration)
	if err != nil {
		return nil, err
	}
	return &AuthSession{User: user, AccessToken: accessToken, RefreshToken: refreshToken, SessionID: sessionID}, nil
}

func (s *AuthService) sendActionToken(ctx context.Context, user *models.User, purpose, queryKey, subject, intro string) error {
	if s.mailer == nil || !s.mailer.Ready() {
		return ErrEmailUnavailable
	}
	token, err := randomToken()
	if err != nil {
		return err
	}
	if err := s.authRepo.CreateActionToken(ctx, user.ID, purpose, hashToken(token), time.Now().Add(emailTokenLifetime)); err != nil {
		return err
	}
	link := s.frontendURL + "/?" + queryKey + "=" + token
	body := intro + "\n\n" + link + "\n\nEste link expira em uma hora. Se você não solicitou esta ação, ignore esta mensagem."
	if err := s.mailer.Send(ctx, user.Email, subject, body); err != nil {
		return fmt.Errorf("%w: %v", ErrEmailUnavailable, err)
	}
	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func randomToken() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
