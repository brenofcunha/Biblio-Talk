package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"lume_api/internal/models"
)

type AuthRepository struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) IsEmailVerified(ctx context.Context, userID string) (bool, error) {
	var verified bool
	err := r.DB.QueryRow(ctx, `SELECT email_verified_at IS NOT NULL FROM users WHERE id=$1`, userID).Scan(&verified)
	return verified, err
}

func (r *AuthRepository) MarkEmailVerified(ctx context.Context, tokenHash string) error {
	tag, err := r.DB.Exec(ctx, `
		WITH consumed AS (
			UPDATE auth_action_tokens
			SET consumed_at=now()
			WHERE token_hash=$1 AND purpose='verify_email'
			  AND consumed_at IS NULL AND expires_at > now()
			RETURNING user_id
		)
		UPDATE users u SET email_verified_at=COALESCE(u.email_verified_at, now())
		FROM consumed c WHERE u.id=c.user_id
	`, tokenHash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *AuthRepository) CreateSession(ctx context.Context, userID, refreshHash string, expiresAt time.Time) (string, error) {
	id := uuid.NewString()
	if _, err := r.DB.Exec(ctx, `
		DELETE FROM auth_sessions
		WHERE expires_at <= now() OR revoked_at < now() - interval '30 days'
	`); err != nil {
		return "", err
	}
	_, err := r.DB.Exec(ctx, `
		INSERT INTO auth_sessions (id, user_id, refresh_token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`, id, userID, refreshHash, expiresAt)
	return id, err
}

func (r *AuthRepository) RotateSession(ctx context.Context, oldHash, newHash string, expiresAt time.Time) (string, string, error) {
	var sessionID, userID string
	err := r.DB.QueryRow(ctx, `
		UPDATE auth_sessions
		SET refresh_token_hash=$2, expires_at=$3
		WHERE refresh_token_hash=$1 AND revoked_at IS NULL AND expires_at > now()
		RETURNING id, user_id
	`, oldHash, newHash, expiresAt).Scan(&sessionID, &userID)
	return sessionID, userID, err
}

func (r *AuthRepository) IsSessionActive(ctx context.Context, sessionID string) (bool, error) {
	var active bool
	err := r.DB.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM auth_sessions
			WHERE id=$1 AND revoked_at IS NULL AND expires_at > now()
		)
	`, sessionID).Scan(&active)
	return active, err
}

func (r *AuthRepository) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := r.DB.Exec(ctx, `UPDATE auth_sessions SET revoked_at=now() WHERE id=$1 AND revoked_at IS NULL`, sessionID)
	return err
}

func (r *AuthRepository) ListActiveSessions(ctx context.Context, userID string) ([]models.AuthSessionInfo, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, created_at, expires_at
		FROM auth_sessions
		WHERE user_id=$1 AND revoked_at IS NULL AND expires_at > now()
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]models.AuthSessionInfo, 0)
	for rows.Next() {
		var session models.AuthSessionInfo
		if err := rows.Scan(&session.ID, &session.CreatedAt, &session.ExpiresAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, rows.Err()
}

func (r *AuthRepository) RevokeUserSession(ctx context.Context, userID, sessionID string) error {
	tag, err := r.DB.Exec(ctx, `
		UPDATE auth_sessions SET revoked_at=now()
		WHERE id=$1 AND user_id=$2 AND revoked_at IS NULL
	`, sessionID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *AuthRepository) RevokeOtherSessions(ctx context.Context, userID, currentSessionID string) (int64, error) {
	tag, err := r.DB.Exec(ctx, `
		UPDATE auth_sessions SET revoked_at=now()
		WHERE user_id=$1 AND id<>$2 AND revoked_at IS NULL
	`, userID, currentSessionID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (r *AuthRepository) CreateActionToken(ctx context.Context, userID, purpose, tokenHash string, expiresAt time.Time) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		DELETE FROM auth_action_tokens
		WHERE expires_at <= now() OR consumed_at < now() - interval '30 days'
	`); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE auth_action_tokens SET consumed_at=now()
		WHERE user_id=$1 AND purpose=$2 AND consumed_at IS NULL
	`, userID, purpose); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO auth_action_tokens (token_hash, user_id, purpose, expires_at)
		VALUES ($1, $2, $3, $4)
	`, tokenHash, userID, purpose, expiresAt); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *AuthRepository) ResetPassword(ctx context.Context, tokenHash, passwordHash string) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var userID string
	err = tx.QueryRow(ctx, `
		UPDATE auth_action_tokens
		SET consumed_at=now()
		WHERE token_hash=$1 AND purpose='reset_password'
		  AND consumed_at IS NULL AND expires_at > now()
		RETURNING user_id
	`, tokenHash).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return pgx.ErrNoRows
	}
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE users SET password_hash=$2 WHERE id=$1`, userID, passwordHash); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE auth_sessions SET revoked_at=now() WHERE user_id=$1 AND revoked_at IS NULL`, userID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *AuthRepository) ChangePassword(ctx context.Context, userID, passwordHash, currentSessionID string) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `UPDATE users SET password_hash=$2 WHERE id=$1`, userID, passwordHash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	if _, err := tx.Exec(ctx, `
		UPDATE auth_sessions SET revoked_at=now()
		WHERE user_id=$1 AND id<>$2 AND revoked_at IS NULL
	`, userID, currentSessionID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
