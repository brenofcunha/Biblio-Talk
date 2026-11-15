package services

import (
        "context"
        "errors"
        "strings"
        "time"

        "github.com/jackc/pgx/v5"
        "lume_api/internal/models"
        "lume_api/internal/repositories"
        "lume_api/internal/utils"
)

type AuthService struct {
        users       *repositories.UserRepository
        jwtSecret   string
        jwtDuration time.Duration
}

func NewAuthService(users *repositories.UserRepository, jwtSecret string, jwtExpDuration time.Duration) *AuthService {
        return &AuthService{users: users, jwtSecret: jwtSecret, jwtDuration: jwtExpDuration}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*models.User, string, error) {
        if strings.TrimSpace(username) == "" || strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
                return nil, "", errors.New("username, email e senha são obrigatórios")
        }
        hash, err := utils.HashPassword(password)
        if err != nil {
                return nil, "", err
        }
        u, err := s.users.Create(ctx, strings.ToLower(strings.TrimSpace(username)), strings.ToLower(strings.TrimSpace(email)), hash)
        if err != nil {
                return nil, "", err
        }
        token, err := utils.GenerateToken(s.jwtSecret, u.ID, u.Username, s.jwtDuration)
        if err != nil {
                return nil, "", err
        }
        return u, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
        u, err := s.users.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
        if err != nil {
                if err == pgx.ErrNoRows {
                        return nil, "", errors.New("credenciais inválidas")
                }
                return nil, "", err
        }
        if err := utils.CheckPassword(u.PasswordHash, password); err != nil {
                return nil, "", errors.New("credenciais inválidas")
        }
        token, err := utils.GenerateToken(s.jwtSecret, u.ID, u.Username, s.jwtDuration)
        if err != nil {
                return nil, "", err
        }
        return u, token, nil
}
