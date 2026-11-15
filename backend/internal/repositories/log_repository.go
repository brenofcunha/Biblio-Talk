package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"lume_api/internal/models"
)

type LogRepository struct {
	DB *pgxpool.Pool
}

func NewLogRepository(db *pgxpool.Pool) *LogRepository { return &LogRepository{DB: db} }

func (r *LogRepository) Create(ctx context.Context, tx pgx.Tx, log models.Log) (*models.Log, error) {
	if log.ID == "" {
		log.ID = uuid.NewString()
	}
	q := `INSERT INTO logs (id,user_id,book_id,rating,review_text,read_date,is_reread,contains_spoilers,likes_count)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,0)
	RETURNING id,user_id,book_id,rating,review_text,read_date,is_reread,contains_spoilers,likes_count,created_at`
	row := tx.QueryRow(ctx, q, log.ID, log.UserID, log.BookID, log.Rating, log.ReviewText, log.ReadDate, log.IsReread, log.ContainsSpoilers)
	out := &models.Log{}
	err := row.Scan(&out.ID, &out.UserID, &out.BookID, &out.Rating, &out.ReviewText, &out.ReadDate, &out.IsReread, &out.ContainsSpoilers, &out.LikesCount, &out.CreatedAt)
	return out, err
}

func (r *LogRepository) GetByID(ctx context.Context, id string) (*models.Log, error) {
	q := `SELECT id,user_id,book_id,rating,review_text,read_date,is_reread,contains_spoilers,likes_count,created_at FROM logs WHERE id=$1`
	log := &models.Log{}
	err := r.DB.QueryRow(ctx, q, id).Scan(&log.ID, &log.UserID, &log.BookID, &log.Rating, &log.ReviewText, &log.ReadDate, &log.IsReread, &log.ContainsSpoilers, &log.LikesCount, &log.CreatedAt)
	if err != nil {
		return nil, err
	}
	return log, nil
}

func (r *LogRepository) Update(ctx context.Context, id string, rating float64, reviewText *string, readDate *time.Time, isReread bool, containsSpoilers bool) error {
	_, err := r.DB.Exec(ctx, `
		UPDATE logs SET rating=$2,review_text=$3,read_date=$4,is_reread=$5,contains_spoilers=$6
		WHERE id=$1
	`, id, rating, reviewText, readDate, isReread, containsSpoilers)
	return err
}

func (r *LogRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM logs WHERE id=$1`, id)
	return err
}

func (r *LogRepository) Like(ctx context.Context, tx pgx.Tx, userID, logID string) error {
	if _, err := tx.Exec(ctx, `INSERT INTO log_likes (user_id, log_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, userID, logID); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, `UPDATE logs SET likes_count=(SELECT count(*) FROM log_likes WHERE log_id=$1) WHERE id=$1`, logID)
	return err
}

func (r *LogRepository) Unlike(ctx context.Context, tx pgx.Tx, userID, logID string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM log_likes WHERE user_id=$1 AND log_id=$2`, userID, logID); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, `UPDATE logs SET likes_count=(SELECT count(*) FROM log_likes WHERE log_id=$1) WHERE id=$1`, logID)
	return err
}

func (r *LogRepository) AddComment(ctx context.Context, userID, logID, text string) (*models.Comment, error) {
	id := uuid.NewString()
	q := `INSERT INTO comments (id,user_id,log_id,comment_text) VALUES ($1,$2,$3,$4)
	RETURNING id,user_id,log_id,comment_text,created_at`
	c := &models.Comment{}
	err := r.DB.QueryRow(ctx, q, id, userID, logID, text).Scan(&c.ID, &c.UserID, &c.LogID, &c.CommentText, &c.CreatedAt)
	return c, err
}

func (r *LogRepository) ListComments(ctx context.Context, logID string) ([]models.Comment, error) {
	rows, err := r.DB.Query(ctx, `SELECT id,user_id,log_id,comment_text,created_at FROM comments WHERE log_id=$1 ORDER BY created_at ASC`, logID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.LogID, &c.CommentText, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *LogRepository) ListUserLogsByUsername(ctx context.Context, username string, limit int) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT l.id,l.user_id,l.book_id,l.rating,l.review_text,l.read_date,l.is_reread,l.contains_spoilers,l.likes_count,l.created_at,
		b.title,b.authors,b.cover_url
		FROM logs l
		JOIN users u ON u.id=l.user_id
		JOIN books b ON b.id=l.book_id
		WHERE u.username=$1
		ORDER BY COALESCE(l.read_date, l.created_at) DESC
		LIMIT $2
	`, username, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]interface{}{}
	for rows.Next() {
		var l models.Log
		var title string
		var authors []string
		var cover *string
		if err := rows.Scan(&l.ID, &l.UserID, &l.BookID, &l.Rating, &l.ReviewText, &l.ReadDate, &l.IsReread, &l.ContainsSpoilers, &l.LikesCount, &l.CreatedAt, &title, &authors, &cover); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{"log": l, "book": map[string]interface{}{"title": title, "authors": authors, "coverUrl": cover}})
	}
	return out, rows.Err()
}
