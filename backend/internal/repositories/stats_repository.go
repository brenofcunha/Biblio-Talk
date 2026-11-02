package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"lume_api/internal/models"
)

type StatsRepository struct { DB *pgxpool.Pool }

func NewStatsRepository(db *pgxpool.Pool) *StatsRepository { return &StatsRepository{DB: db} }

func (r *StatsRepository) GetYearStats(ctx context.Context, username string, year int) (models.YearStats, error) {
	stats := models.YearStats{Year: year, GenresPercent: map[string]float64{}, BooksPerMonth: map[string]int{}}
	var userID string
	if err := r.DB.QueryRow(ctx, `SELECT id FROM users WHERE username=$1`, username).Scan(&userID); err != nil {
		return stats, err
	}

	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)
	prevStart := start.AddDate(-1, 0, 0)
	prevEnd := start

	if err := r.DB.QueryRow(ctx, `SELECT count(*) FROM logs WHERE user_id=$1 AND COALESCE(read_date::date, created_at::date) >= $2 AND COALESCE(read_date::date, created_at::date) < $3`, userID, start, end).Scan(&stats.BooksRead); err != nil {
		return stats, err
	}
	var prevBooks int
	if err := r.DB.QueryRow(ctx, `SELECT count(*) FROM logs WHERE user_id=$1 AND COALESCE(read_date::date, created_at::date) >= $2 AND COALESCE(read_date::date, created_at::date) < $3`, userID, prevStart, prevEnd).Scan(&prevBooks); err != nil {
		return stats, err
	}
	stats.BooksReadDelta = stats.BooksRead - prevBooks

	if err := r.DB.QueryRow(ctx, `
		SELECT COALESCE(sum(COALESCE(b.page_count,0)),0)
		FROM logs l JOIN books b ON b.id=l.book_id
		WHERE l.user_id=$1 AND COALESCE(l.read_date::date, l.created_at::date) >= $2 AND COALESCE(l.read_date::date, l.created_at::date) < $3`,
		userID, start, end).Scan(&stats.PagesRead); err != nil {
		return stats, err
	}
	stats.PagesPerDay = float64(stats.PagesRead) / 365.0

	// Streak
	rows, err := r.DB.Query(ctx, `SELECT DISTINCT COALESCE(read_date::date, created_at::date) as d FROM logs WHERE user_id=$1 ORDER BY d DESC`, userID)
	if err != nil { return stats, err }
	defer rows.Close()
	days := []time.Time{}
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err != nil { return stats, err }
		days = append(days, d)
	}
	stats.CurrentStreakDays, stats.BestStreakDays = calculateStreak(days)

	genreRows, err := r.DB.Query(ctx, `
		SELECT g, count(*)
		FROM logs l JOIN books b ON b.id=l.book_id, unnest(b.genres) g
		WHERE l.user_id=$1 AND COALESCE(l.read_date::date, l.created_at::date) >= $2 AND COALESCE(l.read_date::date, l.created_at::date) < $3
		GROUP BY g`, userID, start, end)
	if err != nil { return stats, err }
	defer genreRows.Close()
	totalGenre := 0
	genreCounts := map[string]int{}
	for genreRows.Next() {
		var g string
		var c int
		if err := genreRows.Scan(&g, &c); err != nil { return stats, err }
		genreCounts[g] = c
		totalGenre += c
	}
	for g, c := range genreCounts {
		if totalGenre > 0 {
			stats.GenresPercent[g] = (float64(c) / float64(totalGenre)) * 100
		}
	}

	monthRows, err := r.DB.Query(ctx, `
		SELECT to_char(COALESCE(read_date::date, created_at::date), 'Mon'), count(*)
		FROM logs WHERE user_id=$1 AND COALESCE(read_date::date, created_at::date) >= $2 AND COALESCE(read_date::date, created_at::date) < $3
		GROUP BY 1`, userID, start, end)
	if err != nil { return stats, err }
	defer monthRows.Close()
	for monthRows.Next() {
		var m string
		var c int
		if err := monthRows.Scan(&m, &c); err != nil { return stats, err }
		stats.BooksPerMonth[m] = c
	}

	if err := r.DB.QueryRow(ctx, `
		SELECT a.author, count(*) as books, COALESCE(sum(COALESCE(b.page_count,0)),0) pages
		FROM logs l
		JOIN books b ON b.id=l.book_id
		LEFT JOIN LATERAL (SELECT unnest(b.authors) as author LIMIT 1) a ON true
		WHERE l.user_id=$1 AND COALESCE(l.read_date::date, l.created_at::date) >= $2 AND COALESCE(l.read_date::date, l.created_at::date) < $3
		GROUP BY a.author
		ORDER BY books DESC,pages DESC
		LIMIT 1`, userID, start, end).Scan(&stats.MostReadAuthor, &stats.MostReadAuthorBooks, &stats.MostReadAuthorPages); err != nil {
		stats.MostReadAuthor = ""
	}

	if err := r.DB.QueryRow(ctx, `
		SELECT count(*)
		FROM logs l
		JOIN social s ON s.following_id=l.user_id
		WHERE s.follower_id=$1 AND l.created_at >= now() - interval '7 days'`, userID).Scan(&stats.NetworkWeekBooks); err != nil {
		return stats, err
	}
	if err := r.DB.QueryRow(ctx, `
		SELECT count(distinct l.user_id)
		FROM logs l
		JOIN social s ON s.following_id=l.user_id
		WHERE s.follower_id=$1 AND l.created_at >= now() - interval '7 days'`, userID).Scan(&stats.NetworkWeekPeople); err != nil {
		return stats, err
	}
	return stats, nil
}

func calculateStreak(days []time.Time) (current int, best int) {
	if len(days) == 0 {
		return 0, 0
	}
	normalized := map[string]bool{}
	for _, d := range days {
		normalized[d.UTC().Format("2006-01-02")] = true
	}
	// current
	today := time.Now().UTC()
	if !normalized[today.Format("2006-01-02")] {
		today = today.AddDate(0, 0, -1)
	}
	for d := today; normalized[d.Format("2006-01-02")]; d = d.AddDate(0, 0, -1) {
		current++
	}
	// best
	list := make([]time.Time, 0, len(normalized))
	for key := range normalized {
		t, _ := time.Parse("2006-01-02", key)
		list = append(list, t)
	}
	for _, t := range list {
		streak := 1
		for d := t.AddDate(0, 0, 1); normalized[d.Format("2006-01-02")]; d = d.AddDate(0, 0, 1) {
			streak++
		}
		if streak > best {
			best = streak
		}
	}
	return
}

func (r *StatsRepository) UpsertCachedStats(ctx context.Context, userID string, year int, payload []byte) error {
	_, err := r.DB.Exec(ctx, `
		INSERT INTO user_year_stats (id,user_id,year,payload,updated_at)
		VALUES ($1,$2,$3,$4,now())
		ON CONFLICT (user_id,year) DO UPDATE SET payload=EXCLUDED.payload, updated_at=now()`, uuid.NewString(), userID, year, payload)
	return err
}
