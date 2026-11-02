package models

import "time"

type User struct {
	ID          string     `json:"id" db:"id"`
	Username    string     `json:"username" db:"username"`
	Email       string     `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Bio         *string    `json:"bio" db:"bio"`
	AvatarURL   *string    `json:"avatarUrl" db:"avatar_url"`
	Top4BookIDs []string   `json:"top4BookIds" db:"top_4_book_ids"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
}

type Book struct {
	ID           string    `json:"id" db:"id"`
	ExternalAPIID *string  `json:"externalApiId" db:"external_api_id"`
	Title        string    `json:"title" db:"title"`
	Authors      []string  `json:"authors" db:"authors"`
	ISBN         *string   `json:"isbn" db:"isbn"`
	Publisher    *string   `json:"publisher" db:"publisher"`
	PublishYear  *int      `json:"publishYear" db:"publish_year"`
	PageCount    *int      `json:"pageCount" db:"page_count"`
	CoverURL     *string   `json:"coverUrl" db:"cover_url"`
	Genres       []string  `json:"genres" db:"genres"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type Log struct {
	ID               string     `json:"id" db:"id"`
	UserID           string     `json:"userId" db:"user_id"`
	BookID           string     `json:"bookId" db:"book_id"`
	Rating           float64    `json:"rating" db:"rating"`
	ReviewText       *string    `json:"reviewText" db:"review_text"`
	ReadDate         *time.Time `json:"readDate" db:"read_date"`
	IsReread         bool       `json:"isReread" db:"is_reread"`
	ContainsSpoilers bool       `json:"containsSpoilers" db:"contains_spoilers"`
	LikesCount       int        `json:"likesCount" db:"likes_count"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
}

type Shelf struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	BookID    string    `json:"bookId" db:"book_id"`
	Status    string    `json:"status" db:"status"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type List struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"userId" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description" db:"description"`
	IsPrivate   bool      `json:"isPrivate" db:"is_private"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type ListItem struct {
	ListID   string  `json:"listId" db:"list_id"`
	BookID   string  `json:"bookId" db:"book_id"`
	Position int     `json:"position" db:"position"`
	Notes    *string `json:"notes" db:"notes"`
}

type Social struct {
	FollowerID  string    `json:"followerId" db:"follower_id"`
	FollowingID string    `json:"followingId" db:"following_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type Comment struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"userId" db:"user_id"`
	LogID       string    `json:"logId" db:"log_id"`
	CommentText string    `json:"commentText" db:"comment_text"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type BookRatingSummary struct {
	AverageRating float64         `json:"averageRating"`
	RatingsCount  int             `json:"ratingsCount"`
	Histogram     map[string]int  `json:"histogram"`
}

type YearStats struct {
	Year              int                `json:"year"`
	BooksRead         int                `json:"booksRead"`
	BooksReadDelta    int                `json:"booksReadDelta"`
	PagesRead         int                `json:"pagesRead"`
	PagesPerDay       float64            `json:"pagesPerDay"`
	CurrentStreakDays int                `json:"currentStreakDays"`
	BestStreakDays    int                `json:"bestStreakDays"`
	GenresPercent     map[string]float64 `json:"genresPercent"`
	BooksPerMonth     map[string]int     `json:"booksPerMonth"`
	MostReadAuthor    string             `json:"mostReadAuthor"`
	MostReadAuthorBooks int              `json:"mostReadAuthorBooks"`
	MostReadAuthorPages int              `json:"mostReadAuthorPages"`
	NetworkWeekBooks  int                `json:"networkWeekBooks"`
	NetworkWeekPeople int                `json:"networkWeekPeople"`
}
