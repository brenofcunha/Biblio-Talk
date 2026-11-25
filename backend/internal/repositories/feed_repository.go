package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"lume_api/internal/utils"
)

type FeedRepository struct{ DB *pgxpool.Pool }

func NewFeedRepository(db *pgxpool.Pool) *FeedRepository { return &FeedRepository{DB: db} }

func (r *FeedRepository) ListFollowingFeed(ctx context.Context, userID string, limit int, cursor string) ([]map[string]interface{}, string, error) {
	base := `
	SELECT l.id,l.created_at,u.username,b.id,b.title,b.authors,b.cover_url,l.rating,l.review_text,l.contains_spoilers,l.likes_count,
	(SELECT count(*) FROM comments c WHERE c.log_id=l.id) AS comments_count
	FROM logs l
	JOIN users u ON u.id=l.user_id
	JOIN books b ON b.id=l.book_id
	JOIN social s ON s.following_id=l.user_id AND s.follower_id=$1
	WHERE COALESCE(u.privacy_settings->>'diary','private')='public'`
	args := []interface{}{userID}
	if cursor != "" {
		cur, err := utils.DecodeFeedCursor(cursor)
		if err != nil {
			return nil, "", err
		}
		base += ` AND (l.created_at,l.id) < ($2,$3)`
		args = append(args, cur.CreatedAt, cur.ID)
	}
	base += fmt.Sprintf(" ORDER BY l.created_at DESC,l.id DESC LIMIT %d", limit+1)
	rows, err := r.DB.Query(ctx, base, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	items := []map[string]interface{}{}
	var next string
	for rows.Next() {
		var id, username, bookID, title string
		var createdAt time.Time
		var authors []string
		var cover, review *string
		var rating float64
		var spoilers bool
		var likes, comments int
		if err := rows.Scan(&id, &createdAt, &username, &bookID, &title, &authors, &cover, &rating, &review, &spoilers, &likes, &comments); err != nil {
			return nil, "", err
		}
		items = append(items, map[string]interface{}{"id": id, "bookId": bookID, "createdAt": createdAt, "user": map[string]interface{}{"username": username}, "book": map[string]interface{}{"title": title, "authors": authors, "coverUrl": cover}, "rating": rating, "reviewText": review, "containsSpoilers": spoilers, "likesCount": likes, "commentsCount": comments})
	}
	if len(items) > limit {
		last := items[limit-1]
		next = utils.EncodeFeedCursor(last["createdAt"].(time.Time), last["id"].(string))
		items = items[:limit]
	}
	return items, next, rows.Err()
}

func (r *FeedRepository) ListGlobalFeed(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(ctx, `SELECT l.id,l.created_at,u.username,b.id,b.title,b.authors,b.cover_url,l.rating,l.review_text,l.contains_spoilers,l.likes_count FROM logs l JOIN users u ON u.id=l.user_id JOIN books b ON b.id=l.book_id WHERE COALESCE(u.privacy_settings->>'diary','private')='public' ORDER BY l.created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []map[string]interface{}{}
	for rows.Next() {
		var id, username, bookID, title string
		var createdAt time.Time
		var authors []string
		var cover, review *string
		var rating float64
		var spoilers bool
		var likes int
		if err := rows.Scan(&id, &createdAt, &username, &bookID, &title, &authors, &cover, &rating, &review, &spoilers, &likes); err != nil {
			return nil, err
		}
		items = append(items, map[string]interface{}{"id": id, "bookId": bookID, "createdAt": createdAt, "user": map[string]interface{}{"username": username}, "book": map[string]interface{}{"title": title, "authors": authors, "coverUrl": cover}, "rating": rating, "reviewText": review, "containsSpoilers": spoilers, "likesCount": likes})
	}
	return items, rows.Err()
}

func (r *FeedRepository) ListReviewOnly(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id,book_id,created_at,username,title,authors,cover_url,rating,review_text,contains_spoilers,likes_count,source
		FROM (
			SELECT l.id,b.id AS book_id,l.created_at,u.username,b.title,b.authors,b.cover_url,l.rating,l.review_text,l.contains_spoilers,l.likes_count,'log'::text AS source
			FROM logs l
			JOIN users u ON u.id=l.user_id
			JOIN books b ON b.id=l.book_id
			WHERE COALESCE(u.privacy_settings->>'diary','private')='public'
			  AND l.review_text IS NOT NULL AND length(trim(l.review_text)) > 0
			UNION ALL
			SELECT r.id,b.id AS book_id,r.published_at,u.username,b.title,b.authors,b.cover_url,r.rating,r.review_text,r.contains_spoilers,0,'review'::text AS source
			FROM reviews r
			JOIN users u ON u.id=r.user_id
			JOIN books b ON b.id=r.book_id
			WHERE r.status='published'
		) published_reviews
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []map[string]interface{}{}
	for rows.Next() {
		var id, bookID, username, title string
		var createdAt time.Time
		var authors []string
		var cover, review *string
		var rating *float64
		var spoilers bool
		var likes int
		var source string
		if err := rows.Scan(&id, &bookID, &createdAt, &username, &title, &authors, &cover, &rating, &review, &spoilers, &likes, &source); err != nil {
			return nil, err
		}
		items = append(items, map[string]interface{}{"id": id, "bookId": bookID, "createdAt": createdAt, "user": map[string]interface{}{"username": username}, "book": map[string]interface{}{"title": title, "authors": authors, "coverUrl": cover}, "rating": rating, "reviewText": review, "containsSpoilers": spoilers, "likesCount": likes, "source": source})
	}
	return items, rows.Err()
}
