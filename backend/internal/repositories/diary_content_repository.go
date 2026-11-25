package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"lume_api/internal/models"
)

func scanPrivateNote(row rowScanner) (*models.PrivateNote, error) {
	note := &models.PrivateNote{}
	err := row.Scan(
		&note.ID, &note.UserID, &note.BookID, &note.ReadingSessionID, &note.Content,
		&note.IsPublished, &note.PublishedAt, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return note, nil
}

const privateNoteColumns = `id,user_id,book_id,reading_session_id,content,is_published,published_at,created_at,updated_at`

func (r *DiaryRepository) CreatePrivateNote(ctx context.Context, userID string, note models.PrivateNote) (*models.PrivateNote, error) {
	if note.ReadingSessionID != nil {
		var ownsSession bool
		if err := r.DB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM reading_sessions WHERE id=$1 AND user_id=$2)`, *note.ReadingSessionID, userID).Scan(&ownsSession); err != nil {
			return nil, err
		}
		if !ownsSession {
			return nil, pgx.ErrNoRows
		}
	}
	return scanPrivateNote(r.DB.QueryRow(ctx, `
		INSERT INTO private_notes (user_id,book_id,reading_session_id,content)
		VALUES ($1,$2,$3,$4)
		RETURNING `+privateNoteColumns,
		userID, note.BookID, note.ReadingSessionID, strings.TrimSpace(note.Content),
	))
}

func (r *DiaryRepository) GetPrivateNote(ctx context.Context, userID, id string) (*models.PrivateNote, error) {
	note, err := scanPrivateNote(r.DB.QueryRow(ctx, `
		SELECT `+privateNoteColumns+`
		FROM private_notes
		WHERE id=$1 AND (user_id=$2 OR is_published=true)
	`, id, userID))
	if err != nil {
		return nil, err
	}
	if err := r.DB.QueryRow(ctx, `SELECT username FROM users WHERE id=$1`, note.UserID).Scan(&note.Username); err != nil {
		return nil, err
	}
	if note.UserID == userID {
		note.Tags, err = r.listPrivateNoteTags(ctx, id)
		if err != nil {
			return nil, err
		}
		note.Quotes, err = r.listPrivateNoteQuotes(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return note, nil
}

func (r *DiaryRepository) listPrivateNoteTags(ctx context.Context, noteID string) ([]models.PrivateNoteTag, error) {
	rows, err := r.DB.Query(ctx, `SELECT id,note_id,tag,created_at FROM private_note_tags WHERE note_id=$1 ORDER BY tag`, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := []models.PrivateNoteTag{}
	for rows.Next() {
		var tag models.PrivateNoteTag
		if err := rows.Scan(&tag.ID, &tag.NoteID, &tag.Tag, &tag.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func (r *DiaryRepository) listPrivateNoteQuotes(ctx context.Context, noteID string) ([]models.PrivateNoteQuote, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id,note_id,quote_text,page,chapter,created_at,updated_at
		FROM private_note_quotes WHERE note_id=$1 ORDER BY created_at
	`, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	quotes := []models.PrivateNoteQuote{}
	for rows.Next() {
		var quote models.PrivateNoteQuote
		if err := rows.Scan(&quote.ID, &quote.NoteID, &quote.QuoteText, &quote.Page, &quote.Chapter, &quote.CreatedAt, &quote.UpdatedAt); err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}
	return quotes, rows.Err()
}

func (r *DiaryRepository) ListPrivateNotes(ctx context.Context, userID, query string, limit int) ([]models.PrivateNote, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT `+privateNoteColumns+`
		FROM private_notes
		WHERE user_id=$1 AND ($2='' OR content ILIKE '%' || $2 || '%')
		ORDER BY updated_at DESC,id DESC LIMIT $3
	`, userID, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	notes := []models.PrivateNote{}
	for rows.Next() {
		note, err := scanPrivateNote(rows)
		if err != nil {
			return nil, err
		}
		notes = append(notes, *note)
	}
	return notes, rows.Err()
}

func (r *DiaryRepository) UpdatePrivateNote(ctx context.Context, userID, id, content string) (*models.PrivateNote, error) {
	return scanPrivateNote(r.DB.QueryRow(ctx, `
		UPDATE private_notes SET content=$3,updated_at=now()
		WHERE id=$1 AND user_id=$2
		RETURNING `+privateNoteColumns,
		id, userID, strings.TrimSpace(content),
	))
}

func (r *DiaryRepository) DeletePrivateNote(ctx context.Context, userID, id string) error {
	tag, err := r.DB.Exec(ctx, `DELETE FROM private_notes WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *DiaryRepository) PublishPrivateNote(ctx context.Context, userID, id string) (*models.PrivateNote, error) {
	return scanPrivateNote(r.DB.QueryRow(ctx, `
		UPDATE private_notes
		SET is_published=true,published_at=now(),updated_at=now()
		WHERE id=$1 AND user_id=$2
		RETURNING `+privateNoteColumns,
		id, userID,
	))
}

func (r *DiaryRepository) CreatePrivateNoteTag(ctx context.Context, userID, noteID, tag string) (*models.PrivateNoteTag, error) {
	var ownsNote bool
	if err := r.DB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM private_notes WHERE id=$1 AND user_id=$2)`, noteID, userID).Scan(&ownsNote); err != nil {
		return nil, err
	}
	if !ownsNote {
		return nil, pgx.ErrNoRows
	}
	item := &models.PrivateNoteTag{}
	err := r.DB.QueryRow(ctx, `
		INSERT INTO private_note_tags (note_id,tag) VALUES ($1,$2)
		RETURNING id,note_id,tag,created_at
	`, noteID, strings.TrimSpace(tag)).Scan(&item.ID, &item.NoteID, &item.Tag, &item.CreatedAt)
	return item, err
}

func (r *DiaryRepository) DeletePrivateNoteTag(ctx context.Context, userID, noteID, tagID string) error {
	tag, err := r.DB.Exec(ctx, `
		DELETE FROM private_note_tags t
		WHERE t.id=$1 AND t.note_id=$2
		  AND EXISTS (SELECT 1 FROM private_notes n WHERE n.id=t.note_id AND n.user_id=$3)
	`, tagID, noteID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *DiaryRepository) CreatePrivateNoteQuote(ctx context.Context, userID, noteID string, quote models.PrivateNoteQuote) (*models.PrivateNoteQuote, error) {
	var ownsNote bool
	if err := r.DB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM private_notes WHERE id=$1 AND user_id=$2)`, noteID, userID).Scan(&ownsNote); err != nil {
		return nil, err
	}
	if !ownsNote {
		return nil, pgx.ErrNoRows
	}
	created := &models.PrivateNoteQuote{}
	err := r.DB.QueryRow(ctx, `
		INSERT INTO private_note_quotes (note_id,quote_text,page,chapter)
		VALUES ($1,$2,$3,$4)
		RETURNING id,note_id,quote_text,page,chapter,created_at,updated_at
	`, noteID, strings.TrimSpace(quote.QuoteText), quote.Page, quote.Chapter).Scan(
		&created.ID, &created.NoteID, &created.QuoteText, &created.Page, &created.Chapter, &created.CreatedAt, &created.UpdatedAt,
	)
	return created, err
}

func (r *DiaryRepository) UpdatePrivateNoteQuote(ctx context.Context, userID, noteID, quoteID string, quote models.PrivateNoteQuote) (*models.PrivateNoteQuote, error) {
	updated := &models.PrivateNoteQuote{}
	err := r.DB.QueryRow(ctx, `
		UPDATE private_note_quotes q
		SET quote_text=$4,page=$5,chapter=$6,updated_at=now()
		WHERE q.id=$1 AND q.note_id=$2
		  AND EXISTS (SELECT 1 FROM private_notes n WHERE n.id=q.note_id AND n.user_id=$3)
		RETURNING q.id,q.note_id,q.quote_text,q.page,q.chapter,q.created_at,q.updated_at
	`, quoteID, noteID, userID, strings.TrimSpace(quote.QuoteText), quote.Page, quote.Chapter).Scan(
		&updated.ID, &updated.NoteID, &updated.QuoteText, &updated.Page, &updated.Chapter, &updated.CreatedAt, &updated.UpdatedAt,
	)
	return updated, err
}

func (r *DiaryRepository) DeletePrivateNoteQuote(ctx context.Context, userID, noteID, quoteID string) error {
	tag, err := r.DB.Exec(ctx, `
		DELETE FROM private_note_quotes q
		WHERE q.id=$1 AND q.note_id=$2
		  AND EXISTS (SELECT 1 FROM private_notes n WHERE n.id=q.note_id AND n.user_id=$3)
	`, quoteID, noteID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func scanReview(row rowScanner) (*models.Review, error) {
	review := &models.Review{}
	err := row.Scan(
		&review.ID, &review.UserID, &review.BookID, &review.ReadingSessionID, &review.Rating,
		&review.ReviewText, &review.ContainsSpoilers, &review.Status, &review.AllowComments,
		&review.AllowLikes, &review.PublishedAt, &review.CreatedAt, &review.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return review, nil
}

const reviewColumns = `id,user_id,book_id,reading_session_id,rating,review_text,contains_spoilers,status,allow_comments,allow_likes,published_at,created_at,updated_at`

func appendReviewHistory(ctx context.Context, tx pgx.Tx, userID, reviewID, action string, review *models.Review) error {
	snapshot, err := json.Marshal(review)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO review_history (review_id,user_id,action,snapshot)
		VALUES ($1,$2,$3,$4::jsonb)
	`, reviewID, userID, action, snapshot)
	return err
}

func (r *DiaryRepository) createReviewTx(ctx context.Context, tx pgx.Tx, userID string, review models.Review, action string) (*models.Review, error) {
	if review.ReadingSessionID != nil {
		var bookID string
		if err := tx.QueryRow(ctx, `
			SELECT book_id FROM reading_sessions WHERE id=$1 AND user_id=$2
		`, *review.ReadingSessionID, userID).Scan(&bookID); err != nil {
			return nil, err
		}
		if bookID != review.BookID {
			return nil, errors.New("readingSessionId não corresponde ao bookId")
		}
	}
	created, err := scanReview(tx.QueryRow(ctx, `
		INSERT INTO reviews (user_id,book_id,reading_session_id,rating,review_text,contains_spoilers,status)
		VALUES ($1,$2,$3,$4,$5,$6,'draft')
		RETURNING `+reviewColumns,
		userID, review.BookID, review.ReadingSessionID, review.Rating, strings.TrimSpace(review.ReviewText), review.ContainsSpoilers,
	))
	if err != nil {
		return nil, err
	}
	if err := appendReviewHistory(ctx, tx, userID, created.ID, action, created); err != nil {
		return nil, err
	}
	return created, nil
}

func (r *DiaryRepository) CreateReviewDraft(ctx context.Context, userID string, review models.Review) (*models.Review, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	created, err := r.createReviewTx(ctx, tx, userID, review, "draft_created")
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return created, nil
}

func (r *DiaryRepository) CreateReviewFromNote(ctx context.Context, userID, noteID string) (*models.Review, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var bookID string
	var sessionID *string
	var content string
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(n.book_id, s.book_id), n.reading_session_id, n.content
		FROM private_notes n
		LEFT JOIN reading_sessions s ON s.id=n.reading_session_id
		WHERE n.id=$1 AND n.user_id=$2
	`, noteID, userID).Scan(&bookID, &sessionID, &content)
	if err != nil {
		return nil, err
	}
	review, err := r.createReviewTx(ctx, tx, userID, models.Review{
		BookID: bookID, ReadingSessionID: sessionID, ReviewText: content,
	}, "created_from_note")
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return review, nil
}

func (r *DiaryRepository) CreateReviewFromLog(ctx context.Context, userID, logID string) (*models.Review, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var review models.Review
	err = tx.QueryRow(ctx, `
		SELECT book_id,rating,COALESCE(review_text,''),contains_spoilers
		FROM logs WHERE id=$1 AND user_id=$2
	`, logID, userID).Scan(&review.BookID, &review.Rating, &review.ReviewText, &review.ContainsSpoilers)
	if err != nil {
		return nil, err
	}
	created, err := r.createReviewTx(ctx, tx, userID, review, "created_from_log")
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return created, nil
}

func (r *DiaryRepository) GetReviewDraft(ctx context.Context, userID, id string) (*models.Review, error) {
	return scanReview(r.DB.QueryRow(ctx, `
		SELECT `+reviewColumns+` FROM reviews
		WHERE id=$1 AND user_id=$2 AND status='draft'
	`, id, userID))
}

func (r *DiaryRepository) GetPublishedReview(ctx context.Context, id string) (*models.Review, error) {
	review, err := scanReview(r.DB.QueryRow(ctx, `
		SELECT `+reviewColumns+` FROM reviews
		WHERE id=$1 AND status='published'
	`, id))
	if err != nil {
		return nil, err
	}
	if err := r.DB.QueryRow(ctx, `SELECT username FROM users WHERE id=$1`, review.UserID).Scan(&review.Username); err != nil {
		return nil, err
	}
	return review, nil
}

func (r *DiaryRepository) UpdateReviewDraft(ctx context.Context, userID, id string, updated models.Review) (*models.Review, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	current, err := scanReview(tx.QueryRow(ctx, `
		SELECT `+reviewColumns+` FROM reviews WHERE id=$1 AND user_id=$2 AND status='draft' FOR UPDATE
	`, id, userID))
	if err != nil {
		return nil, err
	}
	if updated.ReadingSessionID != nil {
		var bookID string
		if err := tx.QueryRow(ctx, `SELECT book_id FROM reading_sessions WHERE id=$1 AND user_id=$2`, *updated.ReadingSessionID, userID).Scan(&bookID); err != nil {
			return nil, err
		}
		if bookID != current.BookID {
			return nil, errors.New("readingSessionId não corresponde ao bookId")
		}
	}
	result, err := scanReview(tx.QueryRow(ctx, `
		UPDATE reviews SET reading_session_id=$3,rating=$4,review_text=$5,contains_spoilers=$6,updated_at=now()
		WHERE id=$1 AND user_id=$2 AND status='draft'
		RETURNING `+reviewColumns,
		id, userID, updated.ReadingSessionID, updated.Rating, strings.TrimSpace(updated.ReviewText), updated.ContainsSpoilers,
	))
	if err != nil {
		return nil, err
	}
	if err := appendReviewHistory(ctx, tx, userID, id, "draft_updated", result); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *DiaryRepository) DeleteReview(ctx context.Context, userID, id string) error {
	tag, err := r.DB.Exec(ctx, `DELETE FROM reviews WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *DiaryRepository) SetReviewPublished(ctx context.Context, userID, id string, published bool) (*models.Review, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	current, err := scanReview(tx.QueryRow(ctx, `
		SELECT `+reviewColumns+` FROM reviews WHERE id=$1 AND user_id=$2 FOR UPDATE
	`, id, userID))
	if err != nil {
		return nil, err
	}
	status, action := "draft", "unpublished"
	if published {
		status, action = "published", "published"
		if strings.TrimSpace(current.ReviewText) == "" {
			return nil, errors.New("resenha vazia não pode ser publicada")
		}
	}
	result, err := scanReview(tx.QueryRow(ctx, `
		UPDATE reviews
		SET status=$3,published_at=CASE WHEN $3='published' THEN now() ELSE NULL END,updated_at=now()
		WHERE id=$1 AND user_id=$2
		RETURNING `+reviewColumns,
		id, userID, status,
	))
	if err != nil {
		return nil, err
	}
	if err := appendReviewHistory(ctx, tx, userID, id, action, result); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *DiaryRepository) UpdateReview(ctx context.Context, userID, id string, updated models.Review) (*models.Review, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	current, err := scanReview(tx.QueryRow(ctx, `
		SELECT `+reviewColumns+` FROM reviews WHERE id=$1 AND user_id=$2 FOR UPDATE
	`, id, userID))
	if err != nil {
		return nil, err
	}
	result, err := scanReview(tx.QueryRow(ctx, `
		UPDATE reviews SET rating=$3,review_text=$4,contains_spoilers=$5,updated_at=now()
		WHERE id=$1 AND user_id=$2
		RETURNING `+reviewColumns,
		id, userID, updated.Rating, strings.TrimSpace(updated.ReviewText), updated.ContainsSpoilers,
	))
	if err != nil {
		return nil, err
	}
	action := "review_updated"
	if current.Status == "published" {
		action = "published_review_edited"
	}
	if err := appendReviewHistory(ctx, tx, userID, id, action, result); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *DiaryRepository) UpdateReviewPrivacy(ctx context.Context, userID, id string, allowComments, allowLikes bool) (*models.Review, error) {
	return scanReview(r.DB.QueryRow(ctx, `
		UPDATE reviews SET allow_comments=$3,allow_likes=$4,updated_at=now()
		WHERE id=$1 AND user_id=$2
		RETURNING `+reviewColumns,
		id, userID, allowComments, allowLikes,
	))
}

func (r *DiaryRepository) ListReviewHistory(ctx context.Context, userID, reviewID string) ([]models.ReviewHistory, error) {
	var ownsReview bool
	if err := r.DB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM reviews WHERE id=$1 AND user_id=$2)`, reviewID, userID).Scan(&ownsReview); err != nil {
		return nil, err
	}
	if !ownsReview {
		return nil, pgx.ErrNoRows
	}
	rows, err := r.DB.Query(ctx, `
		SELECT id,review_id,action,snapshot,created_at
		FROM review_history WHERE review_id=$1
		ORDER BY created_at DESC
	`, reviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	history := []models.ReviewHistory{}
	for rows.Next() {
		var item models.ReviewHistory
		if err := rows.Scan(&item.ID, &item.ReviewID, &item.Action, &item.Snapshot, &item.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, item)
	}
	return history, rows.Err()
}

func (r *DiaryRepository) UpsertReadingGoal(ctx context.Context, userID string, goal models.ReadingGoal) (*models.ReadingGoal, error) {
	result := &models.ReadingGoal{}
	err := r.DB.QueryRow(ctx, `
		INSERT INTO reading_goals (user_id,goal_year,goal_type,goal_month,target_books,target_pages)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (user_id,goal_year,goal_type,(COALESCE(goal_month,0)))
		DO UPDATE SET target_books=EXCLUDED.target_books,target_pages=EXCLUDED.target_pages,updated_at=now()
		RETURNING id,user_id,goal_year,goal_type,goal_month,target_books,target_pages,created_at,updated_at
	`, userID, goal.Year, goal.Type, goal.Month, goal.TargetBooks, goal.TargetPages).Scan(
		&result.ID, &result.UserID, &result.Year, &result.Type, &result.Month,
		&result.TargetBooks, &result.TargetPages, &result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *DiaryRepository) ListReadingGoals(ctx context.Context, userID string, year int) ([]models.ReadingGoal, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id,user_id,goal_year,goal_type,goal_month,target_books,target_pages,created_at,updated_at
		FROM reading_goals WHERE user_id=$1 AND goal_year=$2
		ORDER BY goal_type,goal_month NULLS FIRST
	`, userID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	goals := []models.ReadingGoal{}
	for rows.Next() {
		var goal models.ReadingGoal
		if err := rows.Scan(&goal.ID, &goal.UserID, &goal.Year, &goal.Type, &goal.Month, &goal.TargetBooks, &goal.TargetPages, &goal.CreatedAt, &goal.UpdatedAt); err != nil {
			return nil, err
		}
		goals = append(goals, goal)
	}
	return goals, rows.Err()
}

func (r *DiaryRepository) CreateReadingHistoryImport(ctx context.Context, userID string, sessions []models.ReadingSession) (*models.ReadingHistoryImport, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	sessionIDs := make([]string, 0, len(sessions))
	type importEntry struct {
		ID          string  `json:"id"`
		BookID      string  `json:"book_id"`
		Status      string  `json:"status"`
		CurrentPage int     `json:"current_page"`
		TotalPages  *int    `json:"total_pages"`
		StartedAt   *string `json:"started_at"`
		FinishedAt  *string `json:"finished_at"`
		IsReread    bool    `json:"is_reread"`
	}
	entries := make([]importEntry, 0, len(sessions))
	for _, session := range sessions {
		id := uuid.NewString()
		sessionIDs = append(sessionIDs, id)
		var startedAt, finishedAt *string
		if session.StartedAt != nil {
			value := session.StartedAt.Format("2006-01-02")
			startedAt = &value
		}
		if session.FinishedAt != nil {
			value := session.FinishedAt.Format("2006-01-02")
			finishedAt = &value
		}
		entries = append(entries, importEntry{
			ID: id, BookID: session.BookID, Status: session.Status,
			CurrentPage: session.CurrentPage, TotalPages: session.TotalPages,
			StartedAt: startedAt, FinishedAt: finishedAt, IsReread: session.IsReread,
		})
	}
	entryJSON, err := json.Marshal(entries)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO reading_sessions (id,user_id,book_id,status,current_page,total_pages,started_at,finished_at,is_reread)
		SELECT x.id,$1,x.book_id,x.status,x.current_page,x.total_pages,x.started_at,x.finished_at,x.is_reread
		FROM jsonb_to_recordset($2::jsonb) AS x(
			id uuid,book_id uuid,status text,current_page int,total_pages int,
			started_at date,finished_at date,is_reread boolean
		)
	`, userID, entryJSON); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO shelves (id,user_id,book_id,status)
		SELECT gen_random_uuid(),$1,book_id,
			CASE WHEN status='want_to_read' THEN 'tbr' ELSE status END
		FROM (
			SELECT DISTINCT ON ((entry.value->>'book_id')::uuid)
				(entry.value->>'book_id')::uuid AS book_id,entry.value->>'status' AS status
			FROM jsonb_array_elements($2::jsonb) WITH ORDINALITY AS entry(value,ordinality)
			ORDER BY (entry.value->>'book_id')::uuid,entry.ordinality DESC
		) latest_per_book
		ON CONFLICT (user_id,book_id) DO UPDATE
		SET status=EXCLUDED.status,updated_at=now()
	`, userID, entryJSON); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO reading_session_history (user_id,session_id,event_type,from_status,to_status,details)
		SELECT $1,x.id,'imported',NULL,x.status,'{}'::jsonb
		FROM jsonb_to_recordset($2::jsonb) AS x(id uuid,status text)
	`, userID, entryJSON); err != nil {
		return nil, err
	}
	result, err := json.Marshal(map[string]interface{}{"sessionIds": sessionIDs})
	if err != nil {
		return nil, err
	}
	imported := &models.ReadingHistoryImport{}
	err = tx.QueryRow(ctx, `
		INSERT INTO reading_history_imports (user_id,status,total_items,imported_items,result,completed_at)
		VALUES ($1,'completed',$2,$2,$3::jsonb,now())
		RETURNING id,status,total_items,imported_items,result,created_at,completed_at
	`, userID, len(sessions), result).Scan(
		&imported.ID, &imported.Status, &imported.TotalItems, &imported.ImportedItems,
		&imported.Result, &imported.CreatedAt, &imported.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return imported, nil
}

func (r *DiaryRepository) GetReadingHistoryImport(ctx context.Context, userID, jobID string) (*models.ReadingHistoryImport, error) {
	job := &models.ReadingHistoryImport{}
	err := r.DB.QueryRow(ctx, `
		SELECT id,status,total_items,imported_items,result,created_at,completed_at
		FROM reading_history_imports WHERE id=$1 AND user_id=$2
	`, jobID, userID).Scan(
		&job.ID, &job.Status, &job.TotalItems, &job.ImportedItems,
		&job.Result, &job.CreatedAt, &job.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return job, nil
}
