package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"lume_api/internal/models"
)

type ListRepository struct{ DB *pgxpool.Pool }

func NewListRepository(db *pgxpool.Pool) *ListRepository { return &ListRepository{DB: db} }

func (r *ListRepository) Create(ctx context.Context, userID, title string, description *string, isPrivate bool) (*models.List, error) {
	id := uuid.NewString()
	l := &models.List{}
	err := r.DB.QueryRow(ctx, `INSERT INTO lists (id,user_id,title,description,is_private) VALUES ($1,$2,$3,$4,$5)
	RETURNING id,user_id,title,description,is_private,created_at`, id, userID, title, description, isPrivate).
		Scan(&l.ID, &l.UserID, &l.Title, &l.Description, &l.IsPrivate, &l.CreatedAt)
	return l, err
}

func (r *ListRepository) ListMine(ctx context.Context, userID string) ([]models.List, error) {
	rows, err := r.DB.Query(ctx, `SELECT id,user_id,title,description,is_private,created_at FROM lists WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	out := []models.List{}
	for rows.Next() {
		var l models.List
		if err := rows.Scan(&l.ID, &l.UserID, &l.Title, &l.Description, &l.IsPrivate, &l.CreatedAt); err != nil { return nil, err }
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *ListRepository) GetByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var l models.List
	err := r.DB.QueryRow(ctx, `SELECT id,user_id,title,description,is_private,created_at FROM lists WHERE id=$1`, id).
		Scan(&l.ID, &l.UserID, &l.Title, &l.Description, &l.IsPrivate, &l.CreatedAt)
	if err != nil { return nil, err }
	rows, err := r.DB.Query(ctx, `SELECT li.list_id,li.book_id,li.position,li.notes,b.title,b.cover_url FROM list_items li JOIN books b ON b.id=li.book_id WHERE li.list_id=$1 ORDER BY li.position ASC`, id)
	if err != nil { return nil, err }
	defer rows.Close()
	items := []map[string]interface{}{}
	for rows.Next() {
		var it models.ListItem
		var title string
		var cover *string
		if err := rows.Scan(&it.ListID, &it.BookID, &it.Position, &it.Notes, &title, &cover); err != nil { return nil, err }
		items = append(items, map[string]interface{}{"listItem": it, "book": map[string]interface{}{"title": title, "coverUrl": cover}})
	}
	return map[string]interface{}{"list": l, "items": items}, rows.Err()
}

func (r *ListRepository) Update(ctx context.Context, id, userID, title string, description *string, isPrivate bool) error {
	_, err := r.DB.Exec(ctx, `UPDATE lists SET title=$3,description=$4,is_private=$5 WHERE id=$1 AND user_id=$2`, id, userID, title, description, isPrivate)
	return err
}

func (r *ListRepository) Delete(ctx context.Context, id, userID string) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM lists WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}

func (r *ListRepository) UpsertItem(ctx context.Context, listID, bookID string, position int, notes *string) error {
	_, err := r.DB.Exec(ctx, `INSERT INTO list_items (list_id,book_id,position,notes) VALUES ($1,$2,$3,$4)
	ON CONFLICT (list_id,book_id) DO UPDATE SET position=EXCLUDED.position, notes=EXCLUDED.notes`, listID, bookID, position, notes)
	return err
}

func (r *ListRepository) DeleteItem(ctx context.Context, listID, bookID string) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM list_items WHERE list_id=$1 AND book_id=$2`, listID, bookID)
	return err
}

func (r *ListRepository) Trending(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT l.id,l.title,l.description,u.username,count(li.book_id) as books
		FROM lists l
		JOIN users u ON u.id=l.user_id
		LEFT JOIN list_items li ON li.list_id=l.id
		WHERE l.is_private=false
		GROUP BY l.id,u.username
		ORDER BY books DESC,l.created_at DESC
		LIMIT $1`, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	out := []map[string]interface{}{}
	for rows.Next() {
		var id, title, username string
		var desc *string
		var books int
		if err := rows.Scan(&id, &title, &desc, &username, &books); err != nil { return nil, err }
		out = append(out, map[string]interface{}{"id": id, "title": title, "description": desc, "ownerHandle": username, "booksCount": books})
	}
	return out, rows.Err()
}
