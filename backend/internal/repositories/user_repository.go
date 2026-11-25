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

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository { return &UserRepository{DB: db} }

func (r *UserRepository) Create(ctx context.Context, username, email, passwordHash string) (*models.User, error) {
	id := uuid.NewString()
	q := `INSERT INTO users (id,username,email,password_hash) VALUES ($1,$2,$3,$4)
	RETURNING id,username,display_name,email,password_hash,bio,avatar_url,banner_url,location,
		ARRAY(SELECT unnest(top_4_book_ids)::text),created_at`
	u := &models.User{}
	err := r.DB.QueryRow(ctx, q, id, username, email, passwordHash).
		Scan(&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.BannerURL, &u.Location, &u.Top4BookIDs, &u.CreatedAt)
	return u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	q := `SELECT id,username,display_name,email,password_hash,bio,avatar_url,banner_url,location,
		ARRAY(SELECT unnest(top_4_book_ids)::text),created_at FROM users WHERE lower(email)=lower($1)`
	u := &models.User{}
	err := r.DB.QueryRow(ctx, q, email).
		Scan(&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.BannerURL, &u.Location, &u.Top4BookIDs, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	q := `SELECT id,username,display_name,email,password_hash,bio,avatar_url,banner_url,location,
		ARRAY(SELECT unnest(top_4_book_ids)::text),created_at FROM users WHERE id=$1`
	u := &models.User{}
	err := r.DB.QueryRow(ctx, q, id).
		Scan(&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.BannerURL, &u.Location, &u.Top4BookIDs, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	q := `SELECT id,username,display_name,email,password_hash,bio,avatar_url,banner_url,location,
		ARRAY(SELECT unnest(top_4_book_ids)::text),created_at FROM users WHERE lower(username)=lower($1)`
	u := &models.User{}
	err := r.DB.QueryRow(ctx, q, username).
		Scan(&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.BannerURL, &u.Location, &u.Top4BookIDs, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) UpdateMe(ctx context.Context, id string, bio, avatarURL *string, top4 []string) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var found int
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM books WHERE id::text=ANY($1::text[])`, top4).Scan(&found); err != nil {
		return err
	}
	if found != len(top4) {
		return pgx.ErrNoRows
	}
	tag, err := tx.Exec(ctx, `
		UPDATE users
		SET bio=$2, avatar_url=$3,
			top_4_book_ids=ARRAY(SELECT unnest($4::text[])::uuid)
		WHERE id=$1
	`, id, bio, avatarURL, top4)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return tx.Commit(ctx)
}

func (r *UserRepository) UpdateProfile(ctx context.Context, userID, displayName string, bio, avatarURL, bannerURL, location *string) error {
	tag, err := r.DB.Exec(ctx, `
		UPDATE users
		SET display_name=$2, bio=$3, avatar_url=$4, banner_url=$5, location=$6
		WHERE id=$1
	`, userID, displayName, bio, avatarURL, bannerURL, location)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *UserRepository) GetPrivacy(ctx context.Context, userID string) (models.PrivacySettings, error) {
	var raw []byte
	err := r.DB.QueryRow(ctx, `SELECT privacy_settings FROM users WHERE id=$1`, userID).Scan(&raw)
	if err != nil {
		return models.PrivacySettings{}, err
	}
	settings := defaultPrivacySettings()
	if err := json.Unmarshal(raw, &settings); err != nil {
		return models.PrivacySettings{}, err
	}
	return settings, nil
}

func (r *UserRepository) UpdatePrivacy(ctx context.Context, userID string, settings models.PrivacySettings) error {
	raw, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	tag, err := r.DB.Exec(ctx, `UPDATE users SET privacy_settings=$2 WHERE id=$1`, userID, raw)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *UserRepository) CanView(ctx context.Context, username, resource string) (bool, error) {
	var visibility string
	err := r.DB.QueryRow(ctx, `
		SELECT COALESCE(privacy_settings->>$2, 'private')
		FROM users WHERE lower(username)=lower($1)
	`, username, resource).Scan(&visibility)
	return visibility == "public", err
}

func (r *UserRepository) CanViewByID(ctx context.Context, userID, resource string) (bool, error) {
	var visibility string
	err := r.DB.QueryRow(ctx, `
		SELECT COALESCE(privacy_settings->>$2, 'private')
		FROM users WHERE id=$1
	`, userID, resource).Scan(&visibility)
	return visibility == "public", err
}

func (r *UserRepository) UpdateFavorites(ctx context.Context, userID string, bookIDs []string) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var found int
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM books WHERE id::text=ANY($1::text[])`, bookIDs).Scan(&found); err != nil {
		return err
	}
	if found != len(bookIDs) {
		return pgx.ErrNoRows
	}
	tag, err := tx.Exec(ctx, `
		UPDATE users
		SET top_4_book_ids=ARRAY(SELECT unnest($2::text[])::uuid)
		WHERE id=$1
	`, userID, bookIDs)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return tx.Commit(ctx)
}

func (r *UserRepository) ListPublicFavorites(ctx context.Context, username string) ([]models.Book, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT b.id,b.external_api_id,b.title,b.authors,b.description,b.isbn,b.publisher,
			b.publish_year,b.page_count,b.cover_url,b.genres,b.created_at
		FROM users u
		CROSS JOIN LATERAL unnest(u.top_4_book_ids) WITH ORDINALITY AS favorite(book_id, position)
		JOIN books b ON b.id=favorite.book_id
		WHERE lower(u.username)=lower($1)
		  AND COALESCE(u.privacy_settings->>'favorites', 'private')='public'
		ORDER BY favorite.position
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	books := make([]models.Book, 0, 4)
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.ExternalAPIID, &book.Title, &book.Authors, &book.Description, &book.ISBN, &book.Publisher, &book.PublishYear, &book.PageCount, &book.CoverURL, &book.Genres, &book.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

func (r *UserRepository) ListPublicLists(ctx context.Context, username string) ([]models.List, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT l.id,l.user_id,l.title,l.description,l.is_private,l.created_at
		FROM lists l JOIN users u ON u.id=l.user_id
		WHERE lower(u.username)=lower($1) AND l.is_private=false
		  AND COALESCE(u.privacy_settings->>'lists', 'private')='public'
		ORDER BY l.created_at DESC
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lists := make([]models.List, 0)
	for rows.Next() {
		var list models.List
		if err := rows.Scan(&list.ID, &list.UserID, &list.Title, &list.Description, &list.IsPrivate, &list.CreatedAt); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, rows.Err()
}

func (r *UserRepository) UpsertAnnualChallenge(ctx context.Context, userID string, challenge models.AnnualChallenge) error {
	_, err := r.DB.Exec(ctx, `
		INSERT INTO annual_challenges (user_id,year,target_books,is_public)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (user_id,year) DO UPDATE
		SET target_books=EXCLUDED.target_books,is_public=EXCLUDED.is_public,updated_at=now()
	`, userID, challenge.Year, challenge.TargetBooks, challenge.IsPublic)
	return err
}

func (r *UserRepository) GetPublicAnnualChallenge(ctx context.Context, username string, year int) (models.AnnualChallenge, error) {
	var challenge models.AnnualChallenge
	err := r.DB.QueryRow(ctx, `
		SELECT ac.year,ac.target_books,ac.is_public,count(l.id)::int
		FROM annual_challenges ac
		JOIN users u ON u.id=ac.user_id
		LEFT JOIN logs l ON l.user_id=u.id
			AND COALESCE(l.read_date::date,l.created_at::date) >= make_date(ac.year,1,1)
			AND COALESCE(l.read_date::date,l.created_at::date) < make_date(ac.year+1,1,1)
		WHERE lower(u.username)=lower($1) AND ac.year=$2 AND ac.is_public=true
			AND COALESCE(u.privacy_settings->>'stats','private')='public'
		GROUP BY ac.year,ac.target_books,ac.is_public
	`, username, year).Scan(&challenge.Year, &challenge.TargetBooks, &challenge.IsPublic, &challenge.BooksRead)
	return challenge, err
}

func (r *UserRepository) RecordConsent(ctx context.Context, userID string, consent models.Consent) (*models.Consent, error) {
	err := r.DB.QueryRow(ctx, `
		INSERT INTO user_consents (user_id,consent_type,version,granted)
		VALUES ($1,$2,$3,$4)
		RETURNING id,recorded_at
	`, userID, consent.ConsentType, consent.Version, consent.Granted).Scan(&consent.ID, &consent.RecordedAt)
	if err != nil {
		return nil, err
	}
	return &consent, nil
}

func (r *UserRepository) ExportUserData(ctx context.Context, userID string) (map[string]json.RawMessage, error) {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	queries := map[string]string{
		"user":             `SELECT to_jsonb(u) - 'password_hash' FROM users u WHERE id=$1`,
		"logs":             `SELECT COALESCE(jsonb_agg(to_jsonb(l) ORDER BY l.created_at), '[]'::jsonb) FROM logs l WHERE l.user_id=$1`,
		"shelves":          `SELECT COALESCE(jsonb_agg(to_jsonb(s) ORDER BY s.updated_at), '[]'::jsonb) FROM shelves s WHERE s.user_id=$1`,
		"lists":            `SELECT COALESCE(jsonb_agg(jsonb_build_object('list',to_jsonb(l),'items',COALESCE((SELECT jsonb_agg(to_jsonb(li) ORDER BY li.position) FROM list_items li WHERE li.list_id=l.id),'[]'::jsonb)) ORDER BY l.created_at), '[]'::jsonb) FROM lists l WHERE l.user_id=$1`,
		"following":        `SELECT COALESCE(jsonb_agg(to_jsonb(s) ORDER BY s.created_at), '[]'::jsonb) FROM social s WHERE s.follower_id=$1`,
		"followers":        `SELECT COALESCE(jsonb_agg(to_jsonb(s) ORDER BY s.created_at), '[]'::jsonb) FROM social s WHERE s.following_id=$1`,
		"likes":            `SELECT COALESCE(jsonb_agg(to_jsonb(ll) ORDER BY ll.created_at), '[]'::jsonb) FROM log_likes ll WHERE ll.user_id=$1`,
		"comments":         `SELECT COALESCE(jsonb_agg(to_jsonb(c) ORDER BY c.created_at), '[]'::jsonb) FROM comments c WHERE c.user_id=$1`,
		"yearStats":        `SELECT COALESCE(jsonb_agg(to_jsonb(s) ORDER BY s.year), '[]'::jsonb) FROM user_year_stats s WHERE s.user_id=$1`,
		"annualChallenges": `SELECT COALESCE(jsonb_agg(to_jsonb(ac) ORDER BY ac.year), '[]'::jsonb) FROM annual_challenges ac WHERE ac.user_id=$1`,
		"consents":         `SELECT COALESCE(jsonb_agg(to_jsonb(c) ORDER BY c.recorded_at), '[]'::jsonb) FROM user_consents c WHERE c.user_id=$1`,
		"readingSessions": `SELECT COALESCE(jsonb_agg(
			to_jsonb(rs) ||
			jsonb_build_object(
				'progress', COALESCE((SELECT jsonb_agg(to_jsonb(p) ORDER BY p.progress_date) FROM reading_progress p WHERE p.session_id=rs.id), '[]'::jsonb),
				'history', COALESCE((SELECT jsonb_agg(to_jsonb(h) ORDER BY h.created_at) FROM reading_session_history h WHERE h.session_id=rs.id), '[]'::jsonb)
			) ORDER BY rs.created_at
		), '[]'::jsonb) FROM reading_sessions rs WHERE rs.user_id=$1`,
		"notes": `SELECT COALESCE(jsonb_agg(
			to_jsonb(n) ||
			jsonb_build_object(
				'tags', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.tag) FROM private_note_tags t WHERE t.note_id=n.id), '[]'::jsonb),
				'quotes', COALESCE((SELECT jsonb_agg(to_jsonb(q) ORDER BY q.created_at) FROM private_note_quotes q WHERE q.note_id=n.id), '[]'::jsonb)
			) ORDER BY n.created_at
		), '[]'::jsonb) FROM private_notes n WHERE n.user_id=$1`,
		"reviews": `SELECT COALESCE(jsonb_agg(
			to_jsonb(r) ||
			jsonb_build_object(
				'history', COALESCE((SELECT jsonb_agg(to_jsonb(h) ORDER BY h.created_at) FROM review_history h WHERE h.review_id=r.id), '[]'::jsonb)
			) ORDER BY r.created_at
		), '[]'::jsonb) FROM reviews r WHERE r.user_id=$1`,
		"readingGoals":          `SELECT COALESCE(jsonb_agg(to_jsonb(g) ORDER BY g.goal_year,g.goal_type,g.goal_month), '[]'::jsonb) FROM reading_goals g WHERE g.user_id=$1`,
		"readingHistoryImports": `SELECT COALESCE(jsonb_agg(to_jsonb(i) ORDER BY i.created_at), '[]'::jsonb) FROM reading_history_imports i WHERE i.user_id=$1`,
	}
	export := make(map[string]json.RawMessage, len(queries))
	for key, query := range queries {
		var raw []byte
		if err := tx.QueryRow(ctx, query, userID).Scan(&raw); err != nil {
			return nil, err
		}
		if key == "user" && string(raw) == "null" {
			return nil, pgx.ErrNoRows
		}
		export[key] = json.RawMessage(raw)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return export, nil
}

func (r *UserRepository) DeleteAccount(ctx context.Context, userID string) error {
	tag, err := r.DB.Exec(ctx, `DELETE FROM users WHERE id=$1`, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func defaultPrivacySettings() models.PrivacySettings {
	return models.PrivacySettings{
		Profile: "public", Diary: "public", Progress: "public", Shelf: "public",
		Notes: "private", Comments: "public", Lists: "public", Stats: "public", Favorites: "public",
	}
}

func (r *UserRepository) Follow(ctx context.Context, followerID, followingID string) (bool, error) {
	var created bool
	err := r.DB.QueryRow(ctx, `
		INSERT INTO social (follower_id, following_id)
		VALUES ($1,$2)
		ON CONFLICT DO NOTHING
		RETURNING true
	`, followerID, followingID).Scan(&created)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return created, err
}

func (r *UserRepository) Unfollow(ctx context.Context, followerID, followingID string) (bool, error) {
	var deleted bool
	err := r.DB.QueryRow(ctx, `
		DELETE FROM social WHERE follower_id=$1 AND following_id=$2
		RETURNING true
	`, followerID, followingID).Scan(&deleted)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return deleted, err
}

func (r *UserRepository) ListFollowers(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	q := `SELECT u.id, u.username, s.created_at FROM social s JOIN users u ON u.id=s.follower_id WHERE s.following_id=$1 ORDER BY s.created_at DESC`
	rows, err := r.DB.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]interface{}{}
	for rows.Next() {
		var id, username string
		var followedAt time.Time
		if err := rows.Scan(&id, &username, &followedAt); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{"id": id, "username": username, "followedAt": followedAt})
	}
	return out, rows.Err()
}

func (r *UserRepository) ListFollowing(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	q := `SELECT u.id, u.username, s.created_at FROM social s JOIN users u ON u.id=s.following_id WHERE s.follower_id=$1 ORDER BY s.created_at DESC`
	rows, err := r.DB.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]interface{}{}
	for rows.Next() {
		var id, username string
		var followedAt time.Time
		if err := rows.Scan(&id, &username, &followedAt); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{"id": id, "username": username, "followedAt": followedAt})
	}
	return out, rows.Err()
}

func IsNotFound(err error) bool { return err == pgx.ErrNoRows }
