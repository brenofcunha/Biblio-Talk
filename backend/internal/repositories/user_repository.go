package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"lume_api/internal/models"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository { return &UserRepository{DB: db} }

func (r *UserRepository) Create(ctx context.Context, username, email, passwordHash string) (*models.User, error) {
	id := uuid.NewString()
	q := `INSERT INTO users (id,username,email,password_hash) VALUES ($1,$2,$3,$4)
	RETURNING id,username,email,password_hash,bio,avatar_url,top_4_book_ids,created_at`
	u := &models.User{}
	err := r.DB.QueryRow(ctx, q, id, username, email, passwordHash).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.Top4BookIDs, &u.CreatedAt)
	return u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	q := `SELECT id,username,email,password_hash,bio,avatar_url,top_4_book_ids,created_at FROM users WHERE email=$1`
	u := &models.User{}
	err := r.DB.QueryRow(ctx, q, email).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.Top4BookIDs, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	q := `SELECT id,username,email,password_hash,bio,avatar_url,top_4_book_ids,created_at FROM users WHERE id=$1`
	u := &models.User{}
	err := r.DB.QueryRow(ctx, q, id).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.Top4BookIDs, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	q := `SELECT id,username,email,password_hash,bio,avatar_url,top_4_book_ids,created_at FROM users WHERE username=$1`
	u := &models.User{}
	err := r.DB.QueryRow(ctx, q, username).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.Top4BookIDs, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) UpdateMe(ctx context.Context, id string, bio, avatarURL *string, top4 []string) error {
	_, err := r.DB.Exec(ctx, `
		UPDATE users
		SET bio=$2, avatar_url=$3, top_4_book_ids=$4::text[]
		WHERE id=$1
	`, id, bio, avatarURL, top4)
	return err
}

func (r *UserRepository) Follow(ctx context.Context, followerID, followingID string) error {
	_, err := r.DB.Exec(ctx, `INSERT INTO social (follower_id, following_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, followerID, followingID)
	return err
}

func (r *UserRepository) Unfollow(ctx context.Context, followerID, followingID string) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM social WHERE follower_id=$1 AND following_id=$2`, followerID, followingID)
	return err
}

func (r *UserRepository) ListFollowers(ctx context.Context, userID string) ([]map[string]string, error) {
	q := `SELECT u.id, u.username FROM social s JOIN users u ON u.id=s.follower_id WHERE s.following_id=$1 ORDER BY s.created_at DESC`
	rows, err := r.DB.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]string{}
	for rows.Next() {
		var id, username string
		if err := rows.Scan(&id, &username); err != nil {
			return nil, err
		}
		out = append(out, map[string]string{"id": id, "username": username})
	}
	return out, rows.Err()
}

func (r *UserRepository) ListFollowing(ctx context.Context, userID string) ([]map[string]string, error) {
	q := `SELECT u.id, u.username FROM social s JOIN users u ON u.id=s.following_id WHERE s.follower_id=$1 ORDER BY s.created_at DESC`
	rows, err := r.DB.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]string{}
	for rows.Next() {
		var id, username string
		if err := rows.Scan(&id, &username); err != nil {
			return nil, err
		}
		out = append(out, map[string]string{"id": id, "username": username})
	}
	return out, rows.Err()
}

func IsNotFound(err error) bool { return err == pgx.ErrNoRows }
