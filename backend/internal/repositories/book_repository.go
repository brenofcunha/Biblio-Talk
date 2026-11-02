package repositories

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"lume_api/internal/models"
)

type BookRepository struct {
	DB *pgxpool.Pool
}

func NewBookRepository(db *pgxpool.Pool) *BookRepository { return &BookRepository{DB: db} }

func (r *BookRepository) SearchLocal(ctx context.Context, q string, limit int) ([]models.Book, error) {
	query := `SELECT id,external_api_id,title,authors,isbn,publisher,publish_year,page_count,cover_url,genres,created_at
	FROM books WHERE title ILIKE $1 OR EXISTS (SELECT 1 FROM unnest(authors) a WHERE a ILIKE $1)
	ORDER BY created_at DESC LIMIT $2`
	rows, err := r.DB.Query(ctx, query, "%"+q+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.ExternalAPIID, &b.Title, &b.Authors, &b.ISBN, &b.Publisher, &b.PublishYear, &b.PageCount, &b.CoverURL, &b.Genres, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *BookRepository) UpsertByExternalID(ctx context.Context, b models.Book) (*models.Book, error) {
	if b.ID == "" {
		b.ID = uuid.NewString()
	}
	if b.Authors == nil {
		b.Authors = []string{}
	}
	if b.Genres == nil {
		b.Genres = []string{}
	}
	q := `INSERT INTO books (id,external_api_id,title,authors,isbn,publisher,publish_year,page_count,cover_url,genres)
	VALUES ($1,$2,$3,$4::text[],$5,$6,$7,$8,$9,$10::text[])
	ON CONFLICT (external_api_id) DO UPDATE
	SET title=EXCLUDED.title, authors=EXCLUDED.authors, isbn=EXCLUDED.isbn, publisher=EXCLUDED.publisher,
	publish_year=EXCLUDED.publish_year, page_count=EXCLUDED.page_count, cover_url=EXCLUDED.cover_url, genres=EXCLUDED.genres
	RETURNING id,external_api_id,title,authors,isbn,publisher,publish_year,page_count,cover_url,genres,created_at`
	book := &models.Book{}
	err := r.DB.QueryRow(ctx, q, b.ID, b.ExternalAPIID, strings.TrimSpace(b.Title), b.Authors, b.ISBN, b.Publisher, b.PublishYear, b.PageCount, b.CoverURL, b.Genres).
		Scan(&book.ID, &book.ExternalAPIID, &book.Title, &book.Authors, &book.ISBN, &book.Publisher, &book.PublishYear, &book.PageCount, &book.CoverURL, &book.Genres, &book.CreatedAt)
	return book, err
}

func (r *BookRepository) GetByID(ctx context.Context, id string) (*models.Book, error) {
	b := &models.Book{}
	err := r.DB.QueryRow(ctx, `SELECT id,external_api_id,title,authors,isbn,publisher,publish_year,page_count,cover_url,genres,created_at FROM books WHERE id=$1`, id).
		Scan(&b.ID, &b.ExternalAPIID, &b.Title, &b.Authors, &b.ISBN, &b.Publisher, &b.PublishYear, &b.PageCount, &b.CoverURL, &b.Genres, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *BookRepository) GetBookDetail(ctx context.Context, id string) (map[string]interface{}, error) {
	book, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	q := `SELECT COALESCE(avg(rating),0), count(*) FROM logs WHERE book_id=$1`
	var avg float64
	var count int
	if err := r.DB.QueryRow(ctx, q, id).Scan(&avg, &count); err != nil {
		return nil, err
	}
	hist := map[string]int{"1": 0, "2": 0, "3": 0, "4": 0, "5": 0}
	rows, err := r.DB.Query(ctx, `SELECT floor(rating)::int as star, count(*) FROM logs WHERE book_id=$1 GROUP BY floor(rating)`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var star, c int
		if err := rows.Scan(&star, &c); err != nil {
			return nil, err
		}
		hist["12345"[star-1:star]] = c
	}

	reviewRows, err := r.DB.Query(ctx, `
		SELECT l.id,l.user_id,l.book_id,l.rating,l.review_text,l.read_date,l.is_reread,l.contains_spoilers,l.likes_count,l.created_at,u.username
		FROM logs l JOIN users u ON u.id=l.user_id
		WHERE l.book_id=$1
		ORDER BY l.created_at DESC LIMIT 10
	`, id)
	if err != nil {
		return nil, err
	}
	defer reviewRows.Close()
	reviews := []map[string]interface{}{}
	for reviewRows.Next() {
		var l models.Log
		var username string
		if err := reviewRows.Scan(&l.ID, &l.UserID, &l.BookID, &l.Rating, &l.ReviewText, &l.ReadDate, &l.IsReread, &l.ContainsSpoilers, &l.LikesCount, &l.CreatedAt, &username); err != nil {
			return nil, err
		}
		reviews = append(reviews, map[string]interface{}{"log": l, "username": username})
	}

	return map[string]interface{}{
		"book": book,
		"summary": map[string]interface{}{"averageRating": avg, "ratingsCount": count, "histogram": hist},
		"reviews": reviews,
		"fetchedAt": time.Now(),
	}, nil
}

func (r *BookRepository) PopularWeek(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT b.id,b.title,b.authors,b.cover_url,COALESCE(avg(l.rating),0) as avg_rating,count(l.id) as reads
		FROM books b
		JOIN logs l ON l.book_id=b.id
		WHERE l.created_at >= now() - interval '7 days'
		GROUP BY b.id,b.title,b.authors,b.cover_url
		ORDER BY reads DESC, avg_rating DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]interface{}{}
	for rows.Next() {
		var id, title string
		var authors []string
		var cover *string
		var avg float64
		var reads int
		if err := rows.Scan(&id, &title, &authors, &cover, &avg, &reads); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{"id": id, "title": title, "authors": authors, "coverUrl": cover, "averageRating": avg, "reads": reads})
	}
	return out, rows.Err()
}
