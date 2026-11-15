package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"lume_api/internal/models"
)

type ShelfRepository struct {
	DB *pgxpool.Pool
}

func NewShelfRepository(db *pgxpool.Pool) *ShelfRepository { return &ShelfRepository{DB: db} }

func (r *ShelfRepository) Upsert(ctx context.Context, userID, bookID, status string) error {
	q := `INSERT INTO shelves (id,user_id,book_id,status) VALUES ($1,$2,$3,$4)
	ON CONFLICT (user_id,book_id) DO UPDATE SET status=EXCLUDED.status, updated_at=now()`
	_, err := r.DB.Exec(ctx, q, uuid.NewString(), userID, bookID, status)
	return err
}

func (r *ShelfRepository) ListByUsernameStatus(ctx context.Context, username, status string) ([]models.Shelf, error) {
	q := `SELECT s.id,s.user_id,s.book_id,s.status,s.updated_at FROM shelves s JOIN users u ON u.id=s.user_id WHERE u.username=$1 AND s.status=$2 ORDER BY s.updated_at DESC`
	rows, err := r.DB.Query(ctx, q, username, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.Shelf{}
	for rows.Next() {
		var s models.Shelf
		if err := rows.Scan(&s.ID, &s.UserID, &s.BookID, &s.Status, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
