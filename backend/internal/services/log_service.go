package services

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"lume_api/internal/async"
	"lume_api/internal/models"
	"lume_api/internal/repositories"
)

type LogService struct {
	db      *pgxpool.Pool
	logs    *repositories.LogRepository
	shelves *repositories.ShelfRepository
	users   *repositories.UserRepository
	worker  *async.StatsWorkerPool
}

func NewLogService(db *pgxpool.Pool, logs *repositories.LogRepository, shelves *repositories.ShelfRepository, users *repositories.UserRepository, worker *async.StatsWorkerPool) *LogService {
	return &LogService{db: db, logs: logs, shelves: shelves, users: users, worker: worker}
}

func validateRating(r float64) bool {
	if r < 0.5 || r > 5.0 {
		return false
	}
	v := r * 10
	return int(v)%5 == 0
}

func parseReadDate(readDate *string) (*time.Time, error) {
	if readDate == nil || *readDate == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *readDate)
	if err != nil {
		return nil, errors.New("readDate deve estar no formato YYYY-MM-DD")
	}
	return &t, nil
}

func (s *LogService) Create(ctx context.Context, userID, username string, req models.Log, readDate *string) (*models.Log, error) {
	if !validateRating(req.Rating) {
		return nil, errors.New("rating deve ser entre 0.5 e 5.0 com passos de 0.5")
	}
	rd, err := parseReadDate(readDate)
	if err != nil {
		return nil, err
	}
	req.UserID = userID
	req.ReadDate = rd
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	created, err := s.logs.Create(ctx, tx, req)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO shelves (id,user_id,book_id,status) VALUES (gen_random_uuid(),$1,$2,'read') ON CONFLICT (user_id,book_id) DO UPDATE SET status='read',updated_at=now()`, userID, req.BookID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	s.worker.Enqueue(async.StatsJob{UserID: userID, Username: username, Year: time.Now().Year()})
	return created, nil
}

func (s *LogService) Get(ctx context.Context, id string) (*models.Log, error) { return s.logs.GetByID(ctx, id) }

func (s *LogService) Update(ctx context.Context, userID, username, logID string, req models.Log, readDate *string) error {
	current, err := s.logs.GetByID(ctx, logID)
	if err != nil {
		return err
	}
	if current.UserID != userID {
		return errors.New("você não pode editar este log")
	}
	if !validateRating(req.Rating) {
		return errors.New("rating deve ser entre 0.5 e 5.0 com passos de 0.5")
	}
	rd, err := parseReadDate(readDate)
	if err != nil {
		return err
	}
	if req.ReviewText != nil && len(*req.ReviewText) > 4000 {
		return errors.New("reviewText excede o limite de 4000 caracteres")
	}
	if err := s.logs.Update(ctx, logID, req.Rating, req.ReviewText, rd, req.IsReread, req.ContainsSpoilers); err != nil {
		return err
	}
	s.worker.Enqueue(async.StatsJob{UserID: userID, Username: username, Year: time.Now().Year()})
	return nil
}

func (s *LogService) Delete(ctx context.Context, userID, username, logID string) error {
	current, err := s.logs.GetByID(ctx, logID)
	if err != nil {
		return err
	}
	if current.UserID != userID {
		return errors.New("você não pode excluir este log")
	}
	if err := s.logs.Delete(ctx, logID); err != nil {
		return err
	}
	s.worker.Enqueue(async.StatsJob{UserID: userID, Username: username, Year: time.Now().Year()})
	return nil
}

func (s *LogService) Like(ctx context.Context, userID, logID string) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := s.logs.Like(ctx, tx, userID, logID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *LogService) Unlike(ctx context.Context, userID, logID string) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := s.logs.Unlike(ctx, tx, userID, logID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *LogService) AddComment(ctx context.Context, userID, logID, text string) (*models.Comment, error) {
	if len(text) == 0 {
		return nil, errors.New("commentText é obrigatório")
	}
	return s.logs.AddComment(ctx, userID, logID, text)
}

func (s *LogService) Comments(ctx context.Context, logID string) ([]models.Comment, error) {
	return s.logs.ListComments(ctx, logID)
}

func (s *LogService) UserLogs(ctx context.Context, username string, limit int) ([]map[string]interface{}, error) {
	return s.logs.ListUserLogsByUsername(ctx, username, limit)
}
