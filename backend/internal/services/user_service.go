package services

import (
	"context"

	"lume_api/internal/repositories"
)

type UserService struct {
	users *repositories.UserRepository
	stats *repositories.StatsRepository
}

func NewUserService(users *repositories.UserRepository, stats *repositories.StatsRepository) *UserService {
	return &UserService{users: users, stats: stats}
}

func (s *UserService) GetPublicProfile(ctx context.Context, username string, year int) (map[string]interface{}, error) {
	u, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	st, err := s.stats.GetYearStats(ctx, username, year)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": u.ID,
		"username": u.Username,
		"email": u.Email,
		"bio": u.Bio,
		"avatarUrl": u.AvatarURL,
		"top4BookIds": u.Top4BookIDs,
		"createdAt": u.CreatedAt,
		"stats": st,
	}, nil
}

func (s *UserService) UpdateMe(ctx context.Context, userID string, bio, avatarURL *string, top4 []string) error {
	return s.users.UpdateMe(ctx, userID, bio, avatarURL, top4)
}
