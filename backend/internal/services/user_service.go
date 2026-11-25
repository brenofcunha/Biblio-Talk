package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"lume_api/internal/models"
	"lume_api/internal/repositories"
)

var (
	ErrInvalidProfile     = errors.New("dados de perfil inválidos")
	ErrInvalidPrivacy     = errors.New("configuração de privacidade inválida")
	ErrInvalidConsent     = errors.New("consentimento inválido")
	ErrInvalidChallenge   = errors.New("desafio anual inválido")
	ErrSelfFollow         = errors.New("não é possível seguir a própria conta")
	ErrSocialUserNotFound = errors.New("usuário não encontrado")
)

type UserService struct {
	users *repositories.UserRepository
	stats *repositories.StatsRepository
}

func NewUserService(users *repositories.UserRepository, stats *repositories.StatsRepository) *UserService {
	return &UserService{users: users, stats: stats}
}

func (s *UserService) GetPublicProfile(ctx context.Context, username string, year int) (map[string]interface{}, error) {
	canView, err := s.users.CanView(ctx, username, "profile")
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, pgx.ErrNoRows
	}
	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	privacy, err := s.users.GetPrivacy(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	profile := userProfileData(user, false)
	if privacy.Favorites != "public" {
		delete(profile, "top4BookIds")
	}
	if privacy.Stats == "public" {
		stats, err := s.stats.GetYearStats(ctx, user.Username, year)
		if err != nil {
			return nil, err
		}
		profile["stats"] = stats
	}
	return profile, nil
}

func (s *UserService) GetMyProfile(ctx context.Context, userID string, year int) (map[string]interface{}, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats, err := s.stats.GetYearStats(ctx, user.Username, year)
	if err != nil {
		return nil, err
	}
	profile := userProfileData(user, true)
	profile["stats"] = stats
	privacy, err := s.users.GetPrivacy(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	profile["privacy"] = privacy
	return profile, nil
}

func (s *UserService) UpdateMe(ctx context.Context, userID string, bio, avatarURL *string, top4 []string) error {
	if err := validateFavorites(top4); err != nil {
		return err
	}
	return s.users.UpdateMe(ctx, userID, bio, avatarURL, top4)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, profile models.User) error {
	if len(profile.DisplayName) > 80 || stringLength(profile.Bio) > 500 ||
		stringLength(profile.AvatarURL) > 2048 || stringLength(profile.BannerURL) > 2048 ||
		stringLength(profile.Location) > 120 {
		return fmt.Errorf("%w: displayName até 80, bio até 500, URL até 2048 e location até 120 caracteres", ErrInvalidProfile)
	}
	if profile.Bio != nil {
		value := strings.TrimSpace(*profile.Bio)
		profile.Bio = &value
	}
	return s.users.UpdateProfile(ctx, userID, strings.TrimSpace(profile.DisplayName), profile.Bio, profile.AvatarURL, profile.BannerURL, profile.Location)
}

func (s *UserService) UpdatePrivacy(ctx context.Context, userID string, settings models.PrivacySettings) error {
	for _, value := range []string{
		settings.Profile, settings.Diary, settings.Progress, settings.Shelf,
		settings.Notes, settings.Comments, settings.Lists, settings.Stats, settings.Favorites,
	} {
		if value != "public" && value != "private" {
			return ErrInvalidPrivacy
		}
	}
	return s.users.UpdatePrivacy(ctx, userID, settings)
}

func (s *UserService) CanView(ctx context.Context, username, resource string) (bool, error) {
	profileVisible, err := s.users.CanView(ctx, username, "profile")
	if err != nil || !profileVisible {
		return profileVisible, err
	}
	return s.users.CanView(ctx, username, resource)
}

func (s *UserService) CanViewByID(ctx context.Context, userID, resource string) (bool, error) {
	profileVisible, err := s.users.CanViewByID(ctx, userID, "profile")
	if err != nil || !profileVisible {
		return profileVisible, err
	}
	return s.users.CanViewByID(ctx, userID, resource)
}

func (s *UserService) Follow(ctx context.Context, followerID, followingID string) (bool, error) {
	follower, err := uuid.Parse(followerID)
	if err != nil {
		return false, ErrSocialUserNotFound
	}
	following, err := uuid.Parse(followingID)
	if err != nil {
		return false, ErrSocialUserNotFound
	}
	if follower == following {
		return false, ErrSelfFollow
	}
	user, err := s.users.FindByID(ctx, following.String())
	if errors.Is(err, pgx.ErrNoRows) {
		return false, ErrSocialUserNotFound
	}
	if err != nil {
		return false, err
	}
	canView, err := s.users.CanViewByID(ctx, user.ID, "profile")
	if err != nil {
		return false, err
	}
	if !canView {
		return false, ErrSocialUserNotFound
	}
	return s.users.Follow(ctx, follower.String(), following.String())
}

func (s *UserService) Unfollow(ctx context.Context, followerID, followingID string) (bool, error) {
	follower, err := uuid.Parse(followerID)
	if err != nil {
		return false, ErrSocialUserNotFound
	}
	following, err := uuid.Parse(followingID)
	if err != nil {
		return false, ErrSocialUserNotFound
	}
	if follower == following {
		return false, ErrSelfFollow
	}
	if _, err := s.users.FindByID(ctx, following.String()); errors.Is(err, pgx.ErrNoRows) {
		return false, ErrSocialUserNotFound
	} else if err != nil {
		return false, err
	}
	return s.users.Unfollow(ctx, follower.String(), following.String())
}

func (s *UserService) UpdateFavorites(ctx context.Context, userID string, bookIDs []string) error {
	if err := validateFavorites(bookIDs); err != nil {
		return err
	}
	return s.users.UpdateFavorites(ctx, userID, bookIDs)
}

func (s *UserService) PublicFavorites(ctx context.Context, username string) ([]models.Book, error) {
	canView, err := s.users.CanView(ctx, username, "profile")
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, pgx.ErrNoRows
	}
	canView, err = s.users.CanView(ctx, username, "favorites")
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, pgx.ErrNoRows
	}
	return s.users.ListPublicFavorites(ctx, username)
}

func (s *UserService) PublicLists(ctx context.Context, username string) ([]models.List, error) {
	canView, err := s.users.CanView(ctx, username, "profile")
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, pgx.ErrNoRows
	}
	canView, err = s.users.CanView(ctx, username, "lists")
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, pgx.ErrNoRows
	}
	return s.users.ListPublicLists(ctx, username)
}

func (s *UserService) SaveAnnualChallenge(ctx context.Context, userID string, challenge models.AnnualChallenge) error {
	currentYear := time.Now().Year()
	if challenge.Year < 2000 || challenge.Year > currentYear+1 || challenge.TargetBooks < 1 || challenge.TargetBooks > 1000 {
		return ErrInvalidChallenge
	}
	return s.users.UpsertAnnualChallenge(ctx, userID, challenge)
}

func (s *UserService) PublicAnnualChallenge(ctx context.Context, username string, year int) (models.AnnualChallenge, error) {
	if year < 2000 || year > 2100 {
		return models.AnnualChallenge{}, ErrInvalidChallenge
	}
	canView, err := s.users.CanView(ctx, username, "profile")
	if err != nil {
		return models.AnnualChallenge{}, err
	}
	if !canView {
		return models.AnnualChallenge{}, pgx.ErrNoRows
	}
	return s.users.GetPublicAnnualChallenge(ctx, username, year)
}

func (s *UserService) RecordConsent(ctx context.Context, userID string, consent models.Consent) (*models.Consent, error) {
	if consent.ConsentType != "terms_of_service" && consent.ConsentType != "privacy_policy" && consent.ConsentType != "marketing" {
		return nil, ErrInvalidConsent
	}
	if strings.TrimSpace(consent.Version) == "" || len(consent.Version) > 64 {
		return nil, ErrInvalidConsent
	}
	return s.users.RecordConsent(ctx, userID, consent)
}

func (s *UserService) Export(ctx context.Context, userID string) (map[string]interface{}, error) {
	data, err := s.users.ExportUserData(ctx, userID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"exportedAt": time.Now().UTC(),
		"data":       data,
	}, nil
}

func (s *UserService) DeleteAccount(ctx context.Context, userID string) error {
	return s.users.DeleteAccount(ctx, userID)
}

func userProfileData(user *models.User, includeEmail bool) map[string]interface{} {
	data := map[string]interface{}{
		"id":          user.ID,
		"username":    user.Username,
		"displayName": user.DisplayName,
		"bio":         user.Bio,
		"avatarUrl":   user.AvatarURL,
		"bannerUrl":   user.BannerURL,
		"location":    user.Location,
		"top4BookIds": user.Top4BookIDs,
		"createdAt":   user.CreatedAt,
	}
	if includeEmail {
		data["email"] = user.Email
	}
	return data
}

func stringLength(value *string) int {
	if value == nil {
		return 0
	}
	return len([]rune(*value))
}

func validateFavorites(bookIDs []string) error {
	if len(bookIDs) > 4 {
		return fmt.Errorf("%w: selecione no máximo quatro livros", ErrInvalidProfile)
	}
	seen := make(map[string]struct{}, len(bookIDs))
	for _, bookID := range bookIDs {
		id, err := uuid.Parse(bookID)
		if err != nil {
			return fmt.Errorf("%w: identificador de livro inválido", ErrInvalidProfile)
		}
		normalized := id.String()
		if _, exists := seen[normalized]; exists {
			return fmt.Errorf("%w: não repita livros favoritos", ErrInvalidProfile)
		}
		seen[normalized] = struct{}{}
	}
	return nil
}
