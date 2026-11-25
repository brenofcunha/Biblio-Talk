package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"lume_api/internal/models"
)

type DiaryRepository struct {
	DB *pgxpool.Pool
}

func NewDiaryRepository(db *pgxpool.Pool) *DiaryRepository {
	return &DiaryRepository{DB: db}
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanReadingSession(row rowScanner) (*models.ReadingSession, error) {
	session := &models.ReadingSession{}
	err := row.Scan(
		&session.ID, &session.UserID, &session.BookID, &session.Status,
		&session.CurrentPage, &session.TotalPages, &session.StartedAt,
		&session.FinishedAt, &session.IsReread, &session.CreatedAt, &session.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if session.TotalPages != nil && *session.TotalPages > 0 {
		progressPercent := float64(session.CurrentPage) * 100 / float64(*session.TotalPages)
		session.ProgressPercent = &progressPercent
	}
	return session, nil
}

const readingSessionColumns = `id,user_id,book_id,status,current_page,total_pages,started_at,finished_at,is_reread,created_at,updated_at`

func appendSessionHistory(ctx context.Context, tx pgx.Tx, userID, sessionID, event string, from, to *string, details map[string]interface{}) error {
	if details == nil {
		details = map[string]interface{}{}
	}
	raw, err := json.Marshal(details)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO reading_session_history (user_id,session_id,event_type,from_status,to_status,details)
		VALUES ($1,$2,$3,$4,$5,$6::jsonb)
	`, userID, sessionID, event, from, to, raw)
	return err
}

func upsertShelfTx(ctx context.Context, tx pgx.Tx, userID, bookID, status string) error {
	if status == "want_to_read" {
		status = "tbr"
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO shelves (id,user_id,book_id,status)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (user_id,book_id) DO UPDATE SET status=EXCLUDED.status,updated_at=now()
	`, uuid.NewString(), userID, bookID, status)
	return err
}

func (r *DiaryRepository) CreateReadingSession(ctx context.Context, userID string, session models.ReadingSession) (*models.ReadingSession, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	created, err := scanReadingSession(tx.QueryRow(ctx, `
		INSERT INTO reading_sessions (user_id,book_id,status,current_page,total_pages,started_at,finished_at,is_reread)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING `+readingSessionColumns,
		userID, session.BookID, session.Status, session.CurrentPage, session.TotalPages,
		session.StartedAt, session.FinishedAt, session.IsReread,
	))
	if err != nil {
		return nil, err
	}
	if err := upsertShelfTx(ctx, tx, userID, session.BookID, session.Status); err != nil {
		return nil, err
	}
	if err := appendSessionHistory(ctx, tx, userID, created.ID, "created", nil, &created.Status, map[string]interface{}{
		"currentPage": created.CurrentPage,
		"totalPages":  created.TotalPages,
	}); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return created, nil
}

func (r *DiaryRepository) GetReadingSession(ctx context.Context, userID, id string) (*models.ReadingSession, error) {
	return scanReadingSession(r.DB.QueryRow(ctx, `
		SELECT `+readingSessionColumns+`
		FROM reading_sessions WHERE id=$1 AND user_id=$2
	`, id, userID))
}

func (r *DiaryRepository) ListReadingSessions(ctx context.Context, userID string, limit int) ([]models.ReadingSession, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT `+readingSessionColumns+`
		FROM reading_sessions WHERE user_id=$1
		ORDER BY updated_at DESC,id DESC LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sessions := []models.ReadingSession{}
	for rows.Next() {
		session, err := scanReadingSession(rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, *session)
	}
	return sessions, rows.Err()
}

func (r *DiaryRepository) UpdateReadingSession(ctx context.Context, userID, id string, updated models.ReadingSession) (*models.ReadingSession, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	current, err := scanReadingSession(tx.QueryRow(ctx, `
		SELECT `+readingSessionColumns+` FROM reading_sessions WHERE id=$1 AND user_id=$2 FOR UPDATE
	`, id, userID))
	if err != nil {
		return nil, err
	}
	result, err := scanReadingSession(tx.QueryRow(ctx, `
		UPDATE reading_sessions
		SET status=$3,current_page=$4,total_pages=$5,started_at=$6,finished_at=$7,is_reread=$8,updated_at=now()
		WHERE id=$1 AND user_id=$2
		RETURNING `+readingSessionColumns,
		id, userID, updated.Status, updated.CurrentPage, updated.TotalPages,
		updated.StartedAt, updated.FinishedAt, updated.IsReread,
	))
	if err != nil {
		return nil, err
	}
	if err := upsertShelfTx(ctx, tx, userID, current.BookID, updated.Status); err != nil {
		return nil, err
	}
	var from *string
	if current.Status != updated.Status {
		from = &current.Status
	}
	if err := appendSessionHistory(ctx, tx, userID, id, "updated", from, &updated.Status, map[string]interface{}{
		"currentPage": updated.CurrentPage,
		"totalPages":  updated.TotalPages,
		"startedAt":   updated.StartedAt,
		"finishedAt":  updated.FinishedAt,
		"isReread":    updated.IsReread,
	}); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *DiaryRepository) UpdateReadingStatus(ctx context.Context, userID, id, status string) (*models.ReadingSession, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	current, err := scanReadingSession(tx.QueryRow(ctx, `
		SELECT `+readingSessionColumns+` FROM reading_sessions WHERE id=$1 AND user_id=$2 FOR UPDATE
	`, id, userID))
	if err != nil {
		return nil, err
	}
	from := current.Status
	if status == "reading" && current.StartedAt == nil {
		now := time.Now().UTC()
		current.StartedAt = &now
	}
	if status == "read" {
		now := time.Now().UTC()
		current.FinishedAt = &now
		if current.TotalPages != nil {
			current.CurrentPage = *current.TotalPages
		}
	}
	current.Status = status
	updated, err := scanReadingSession(tx.QueryRow(ctx, `
		UPDATE reading_sessions
		SET status=$3,current_page=$4,started_at=$5,finished_at=$6,updated_at=now()
		WHERE id=$1 AND user_id=$2
		RETURNING `+readingSessionColumns,
		id, userID, current.Status, current.CurrentPage, current.StartedAt, current.FinishedAt,
	))
	if err != nil {
		return nil, err
	}
	if err := upsertShelfTx(ctx, tx, userID, current.BookID, status); err != nil {
		return nil, err
	}
	if err := appendSessionHistory(ctx, tx, userID, id, "status_changed", &from, &status, nil); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *DiaryRepository) DeleteReadingSession(ctx context.Context, userID, id string) error {
	tag, err := r.DB.Exec(ctx, `DELETE FROM reading_sessions WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *DiaryRepository) AddReadingProgress(ctx context.Context, userID, sessionID string, progress models.ReadingProgress) (*models.ReadingProgress, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	session, err := scanReadingSession(tx.QueryRow(ctx, `
		SELECT `+readingSessionColumns+` FROM reading_sessions WHERE id=$1 AND user_id=$2 FOR UPDATE
	`, sessionID, userID))
	if err != nil {
		return nil, err
	}
	if session.TotalPages != nil && progress.CurrentPage > *session.TotalPages {
		return nil, errors.New("currentPage não pode exceder totalPages")
	}
	if session.TotalPages != nil && progress.PagesRead > *session.TotalPages {
		return nil, errors.New("pagesRead não pode exceder totalPages")
	}
	created := &models.ReadingProgress{}
	err = tx.QueryRow(ctx, `
		INSERT INTO reading_progress (session_id,user_id,progress_date,current_page,pages_read,note)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (session_id,progress_date) DO UPDATE
		SET current_page=EXCLUDED.current_page,pages_read=EXCLUDED.pages_read,note=EXCLUDED.note
		RETURNING id,session_id,user_id,progress_date,current_page,pages_read,note,created_at
	`, sessionID, userID, progress.ProgressDate, progress.CurrentPage, progress.PagesRead, progress.Note).Scan(
		&created.ID, &created.SessionID, &created.UserID, &created.ProgressDate,
		&created.CurrentPage, &created.PagesRead, &created.Note, &created.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE reading_sessions
		SET current_page=$3,status='reading',started_at=CASE WHEN $5 THEN $4 ELSE COALESCE(started_at,$4) END,finished_at=NULL,updated_at=now()
		WHERE id=$1 AND user_id=$2
	`, sessionID, userID, progress.CurrentPage, progress.ProgressDate, session.Status == "read"); err != nil {
		return nil, err
	}
	if err := appendSessionHistory(ctx, tx, userID, sessionID, "progress_recorded", nil, nil, map[string]interface{}{
		"date":        progress.ProgressDate.Format("2006-01-02"),
		"currentPage": progress.CurrentPage,
		"pagesRead":   progress.PagesRead,
	}); err != nil {
		return nil, err
	}
	if session.Status != "reading" {
		if err := upsertShelfTx(ctx, tx, userID, session.BookID, "reading"); err != nil {
			return nil, err
		}
		if err := appendSessionHistory(ctx, tx, userID, sessionID, "status_changed", &session.Status, strPtr("reading"), nil); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return created, nil
}

func strPtr(value string) *string { return &value }

func (r *DiaryRepository) ListReadingProgress(ctx context.Context, userID, sessionID string) ([]models.ReadingProgress, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id,session_id,user_id,progress_date,current_page,pages_read,note,created_at
		FROM reading_progress
		WHERE session_id=$1 AND user_id=$2
		ORDER BY progress_date DESC
	`, sessionID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	progress := []models.ReadingProgress{}
	for rows.Next() {
		var item models.ReadingProgress
		if err := rows.Scan(&item.ID, &item.SessionID, &item.UserID, &item.ProgressDate, &item.CurrentPage, &item.PagesRead, &item.Note, &item.CreatedAt); err != nil {
			return nil, err
		}
		progress = append(progress, item)
	}
	return progress, rows.Err()
}

func (r *DiaryRepository) ListReadingHistory(ctx context.Context, userID, sessionID string) ([]models.ReadingSessionHistory, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id,session_id,event_type,from_status,to_status,details,created_at
		FROM reading_session_history
		WHERE session_id=$1 AND user_id=$2
		ORDER BY created_at DESC
	`, sessionID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	history := []models.ReadingSessionHistory{}
	for rows.Next() {
		var item models.ReadingSessionHistory
		if err := rows.Scan(&item.ID, &item.SessionID, &item.EventType, &item.FromStatus, &item.ToStatus, &item.Details, &item.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, item)
	}
	return history, rows.Err()
}
