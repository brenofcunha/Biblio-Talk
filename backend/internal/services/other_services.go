package services

import (
	"context"

	"lume_api/internal/models"
	"lume_api/internal/repositories"
)

type ShelfService struct{ repo *repositories.ShelfRepository }

func NewShelfService(repo *repositories.ShelfRepository) *ShelfService { return &ShelfService{repo: repo} }

func (s *ShelfService) Upsert(ctx context.Context, userID, bookID, status string) error {
	return s.repo.Upsert(ctx, userID, bookID, status)
}

func (s *ShelfService) ByStatus(ctx context.Context, username, status string) ([]models.Shelf, error) {
	return s.repo.ListByUsernameStatus(ctx, username, status)
}

type ListService struct{ repo *repositories.ListRepository }

func NewListService(repo *repositories.ListRepository) *ListService { return &ListService{repo: repo} }
func (s *ListService) Create(ctx context.Context, userID, title string, description *string, isPrivate bool) (*models.List, error) {
	return s.repo.Create(ctx, userID, title, description, isPrivate)
}
func (s *ListService) Mine(ctx context.Context, userID string) ([]models.List, error) { return s.repo.ListMine(ctx, userID) }
func (s *ListService) Detail(ctx context.Context, id string) (map[string]interface{}, error) { return s.repo.GetByID(ctx, id) }
func (s *ListService) Update(ctx context.Context, id, userID, title string, description *string, isPrivate bool) error {
	return s.repo.Update(ctx, id, userID, title, description, isPrivate)
}
func (s *ListService) Delete(ctx context.Context, id, userID string) error { return s.repo.Delete(ctx, id, userID) }
func (s *ListService) UpsertItem(ctx context.Context, listID, bookID string, position int, notes *string) error {
	return s.repo.UpsertItem(ctx, listID, bookID, position, notes)
}
func (s *ListService) DeleteItem(ctx context.Context, listID, bookID string) error { return s.repo.DeleteItem(ctx, listID, bookID) }
func (s *ListService) Trending(ctx context.Context, limit int) ([]map[string]interface{}, error) { return s.repo.Trending(ctx, limit) }

type FeedService struct{ repo *repositories.FeedRepository }

func NewFeedService(repo *repositories.FeedRepository) *FeedService { return &FeedService{repo: repo} }

func (s *FeedService) Following(ctx context.Context, userID string, limit int, cursor string) ([]map[string]interface{}, string, error) {
	return s.repo.ListFollowingFeed(ctx, userID, limit, cursor)
}
func (s *FeedService) Global(ctx context.Context, limit int) ([]map[string]interface{}, error) { return s.repo.ListGlobalFeed(ctx, limit) }
func (s *FeedService) Reviews(ctx context.Context, limit int) ([]map[string]interface{}, error) { return s.repo.ListReviewOnly(ctx, limit) }

type StatsService struct{ repo *repositories.StatsRepository }

func NewStatsService(repo *repositories.StatsRepository) *StatsService { return &StatsService{repo: repo} }

func (s *StatsService) Year(ctx context.Context, username string, year int) (models.YearStats, error) {
	return s.repo.GetYearStats(ctx, username, year)
}
