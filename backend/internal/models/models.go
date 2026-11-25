package models

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	DisplayName  string    `json:"displayName" db:"display_name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Bio          *string   `json:"bio" db:"bio"`
	AvatarURL    *string   `json:"avatarUrl" db:"avatar_url"`
	BannerURL    *string   `json:"bannerUrl" db:"banner_url"`
	Location     *string   `json:"location" db:"location"`
	Top4BookIDs  []string  `json:"top4BookIds" db:"top_4_book_ids"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type PrivacySettings struct {
	Profile   string `json:"profile"`
	Diary     string `json:"diary"`
	Progress  string `json:"progress"`
	Shelf     string `json:"shelf"`
	Notes     string `json:"notes"`
	Comments  string `json:"comments"`
	Lists     string `json:"lists"`
	Stats     string `json:"stats"`
	Favorites string `json:"favorites"`
}

type Consent struct {
	ID          string    `json:"id"`
	ConsentType string    `json:"consentType"`
	Version     string    `json:"version"`
	Granted     bool      `json:"granted"`
	RecordedAt  time.Time `json:"recordedAt"`
}

type AuthSessionInfo struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	IsCurrent bool      `json:"isCurrent"`
}

type AnnualChallenge struct {
	Year        int  `json:"year"`
	TargetBooks int  `json:"targetBooks"`
	BooksRead   int  `json:"booksRead"`
	IsPublic    bool `json:"isPublic"`
}

type Book struct {
	ID            string    `json:"id" db:"id"`
	ExternalAPIID *string   `json:"externalApiId" db:"external_api_id"`
	Title         string    `json:"title" db:"title"`
	Authors       []string  `json:"authors" db:"authors"`
	Description   *string   `json:"description" db:"description"`
	ISBN          *string   `json:"isbn" db:"isbn"`
	Publisher     *string   `json:"publisher" db:"publisher"`
	PublishYear   *int      `json:"publishYear" db:"publish_year"`
	PageCount     *int      `json:"pageCount" db:"page_count"`
	CoverURL      *string   `json:"coverUrl" db:"cover_url"`
	Genres        []string  `json:"genres" db:"genres"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
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

type ReadingSession struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"userId" db:"user_id"`
	BookID          string     `json:"bookId" db:"book_id"`
	Status          string     `json:"status" db:"status"`
	CurrentPage     int        `json:"currentPage" db:"current_page"`
	TotalPages      *int       `json:"totalPages" db:"total_pages"`
	ProgressPercent *float64   `json:"progressPercent,omitempty" db:"-"`
	StartedAt       *time.Time `json:"startedAt" db:"started_at"`
	FinishedAt      *time.Time `json:"finishedAt" db:"finished_at"`
	IsReread        bool       `json:"isReread" db:"is_reread"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`
}

type ReadingProgress struct {
	ID           string    `json:"id" db:"id"`
	SessionID    string    `json:"sessionId" db:"session_id"`
	UserID       string    `json:"userId" db:"user_id"`
	ProgressDate time.Time `json:"date" db:"progress_date"`
	CurrentPage  int       `json:"currentPage" db:"current_page"`
	PagesRead    int       `json:"pagesRead" db:"pages_read"`
	Note         *string   `json:"note" db:"note"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type ReadingSessionHistory struct {
	ID         string                 `json:"id" db:"id"`
	SessionID  string                 `json:"sessionId" db:"session_id"`
	EventType  string                 `json:"eventType" db:"event_type"`
	FromStatus *string                `json:"fromStatus" db:"from_status"`
	ToStatus   *string                `json:"toStatus" db:"to_status"`
	Details    map[string]interface{} `json:"details" db:"details"`
	CreatedAt  time.Time              `json:"createdAt" db:"created_at"`
}

type PrivateNote struct {
	ID               string             `json:"id" db:"id"`
	UserID           string             `json:"-" db:"user_id"`
	Username         string             `json:"username,omitempty" db:"-"`
	BookID           *string            `json:"bookId" db:"book_id"`
	ReadingSessionID *string            `json:"readingSessionId" db:"reading_session_id"`
	Content          string             `json:"content" db:"content"`
	IsPublished      bool               `json:"isPublished" db:"is_published"`
	PublishedAt      *time.Time         `json:"publishedAt" db:"published_at"`
	CreatedAt        time.Time          `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time          `json:"updatedAt" db:"updated_at"`
	Tags             []PrivateNoteTag   `json:"tags,omitempty"`
	Quotes           []PrivateNoteQuote `json:"quotes,omitempty"`
}

type PrivateNoteTag struct {
	ID        string    `json:"id" db:"id"`
	NoteID    string    `json:"noteId" db:"note_id"`
	Tag       string    `json:"tag" db:"tag"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type PrivateNoteQuote struct {
	ID        string    `json:"id" db:"id"`
	NoteID    string    `json:"noteId" db:"note_id"`
	QuoteText string    `json:"quoteText" db:"quote_text"`
	Page      *int      `json:"page" db:"page"`
	Chapter   *string   `json:"chapter" db:"chapter"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Review struct {
	ID               string     `json:"id" db:"id"`
	UserID           string     `json:"-" db:"user_id"`
	Username         string     `json:"username,omitempty" db:"-"`
	BookID           string     `json:"bookId" db:"book_id"`
	ReadingSessionID *string    `json:"readingSessionId" db:"reading_session_id"`
	Rating           *float64   `json:"rating" db:"rating"`
	ReviewText       string     `json:"reviewText" db:"review_text"`
	ContainsSpoilers bool       `json:"containsSpoilers" db:"contains_spoilers"`
	Status           string     `json:"status" db:"status"`
	AllowComments    bool       `json:"allowComments" db:"allow_comments"`
	AllowLikes       bool       `json:"allowLikes" db:"allow_likes"`
	PublishedAt      *time.Time `json:"publishedAt" db:"published_at"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time  `json:"updatedAt" db:"updated_at"`
}

type ReviewHistory struct {
	ID        string                 `json:"id" db:"id"`
	ReviewID  string                 `json:"reviewId" db:"review_id"`
	Action    string                 `json:"action" db:"action"`
	Snapshot  map[string]interface{} `json:"snapshot" db:"snapshot"`
	CreatedAt time.Time              `json:"createdAt" db:"created_at"`
}

type ReadingGoal struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"userId" db:"user_id"`
	Year        int       `json:"year" db:"goal_year"`
	Type        string    `json:"type" db:"goal_type"`
	Month       *int      `json:"month" db:"goal_month"`
	TargetBooks *int      `json:"targetBooks" db:"target_books"`
	TargetPages *int      `json:"targetPages" db:"target_pages"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type ReadingHistoryImport struct {
	ID            string                 `json:"id" db:"id"`
	Status        string                 `json:"status" db:"status"`
	TotalItems    int                    `json:"totalItems" db:"total_items"`
	ImportedItems int                    `json:"importedItems" db:"imported_items"`
	Result        map[string]interface{} `json:"result" db:"result"`
	CreatedAt     time.Time              `json:"createdAt" db:"created_at"`
	CompletedAt   *time.Time             `json:"completedAt" db:"completed_at"`
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
	AverageRating float64        `json:"averageRating"`
	RatingsCount  int            `json:"ratingsCount"`
	Histogram     map[string]int `json:"histogram"`
}

type YearStats struct {
	Year                int                `json:"year"`
	BooksRead           int                `json:"booksRead"`
	BooksReadDelta      int                `json:"booksReadDelta"`
	PagesRead           int                `json:"pagesRead"`
	PagesPerDay         float64            `json:"pagesPerDay"`
	CurrentStreakDays   int                `json:"currentStreakDays"`
	BestStreakDays      int                `json:"bestStreakDays"`
	GenresPercent       map[string]float64 `json:"genresPercent"`
	BooksPerMonth       map[string]int     `json:"booksPerMonth"`
	MostReadAuthor      string             `json:"mostReadAuthor"`
	MostReadAuthorBooks int                `json:"mostReadAuthorBooks"`
	MostReadAuthorPages int                `json:"mostReadAuthorPages"`
	NetworkWeekBooks    int                `json:"networkWeekBooks"`
	NetworkWeekPeople   int                `json:"networkWeekPeople"`
}
