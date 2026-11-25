package dto

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailRequest struct {
	Email string `json:"email"`
}

type TokenRequest struct {
	Token string `json:"token"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type UpdateMeRequest struct {
	Bio         *string  `json:"bio"`
	AvatarURL   *string  `json:"avatarUrl"`
	Top4BookIDs []string `json:"top4BookIds"`
}

type UpdateProfileRequest struct {
	DisplayName string  `json:"displayName"`
	Bio         *string `json:"bio"`
	AvatarURL   *string `json:"avatarUrl"`
	BannerURL   *string `json:"bannerUrl"`
	Location    *string `json:"location"`
}

type UpdatePrivacyRequest struct {
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

type UpdateFavoritesRequest struct {
	BookIDs []string `json:"bookIds"`
}

type AnnualChallengeRequest struct {
	Year        int  `json:"year"`
	TargetBooks int  `json:"targetBooks"`
	IsPublic    bool `json:"isPublic"`
}

type ConsentRequest struct {
	ConsentType string `json:"consentType"`
	Version     string `json:"version"`
	Granted     *bool  `json:"granted"`
}

type ShelfUpsertRequest struct {
	BookID string `json:"bookId"`
	Status string `json:"status"`
}

type CreateLogRequest struct {
	BookID           string   `json:"bookId"`
	Rating           float64  `json:"rating"`
	ReviewText       *string  `json:"reviewText"`
	ReadDate         *string  `json:"readDate"`
	IsReread         bool     `json:"isReread"`
	ContainsSpoilers bool     `json:"containsSpoilers"`
	Tags             []string `json:"tags"`
	Liked            bool     `json:"liked"`
}

type UpdateLogRequest = CreateLogRequest

type ReadingSessionRequest struct {
	BookID      string  `json:"bookId"`
	Status      string  `json:"status"`
	CurrentPage int     `json:"currentPage"`
	TotalPages  *int    `json:"totalPages"`
	StartedAt   *string `json:"startedAt"`
	FinishedAt  *string `json:"finishedAt"`
	IsReread    bool    `json:"isReread"`
}

type ReadingSessionStatusRequest struct {
	Status string `json:"status"`
}

type ReadingProgressRequest struct {
	Date        string  `json:"date"`
	CurrentPage int     `json:"currentPage"`
	PagesRead   int     `json:"pagesRead"`
	Note        *string `json:"note"`
}

type PrivateNoteRequest struct {
	BookID           *string `json:"bookId"`
	ReadingSessionID *string `json:"readingSessionId"`
	Content          string  `json:"content"`
}

type UpdatePrivateNoteRequest struct {
	Content string `json:"content"`
}

type NoteTagRequest struct {
	Tag string `json:"tag"`
}

type NoteQuoteRequest struct {
	QuoteText string  `json:"quoteText"`
	Page      *int    `json:"page"`
	Chapter   *string `json:"chapter"`
}

type ReviewDraftRequest struct {
	BookID           string   `json:"bookId"`
	ReadingSessionID *string  `json:"readingSessionId"`
	Rating           *float64 `json:"rating"`
	ReviewText       string   `json:"reviewText"`
	ContainsSpoilers bool     `json:"containsSpoilers"`
}

type UpdateReviewRequest struct {
	Rating           *float64 `json:"rating"`
	ReviewText       string   `json:"reviewText"`
	ContainsSpoilers bool     `json:"containsSpoilers"`
}

type ReviewPrivacyRequest struct {
	AllowComments *bool `json:"allowComments"`
	AllowLikes    *bool `json:"allowLikes"`
}

type ReadingGoalRequest struct {
	Year        int    `json:"year"`
	Type        string `json:"type"`
	Month       *int   `json:"month"`
	TargetBooks *int   `json:"targetBooks"`
	TargetPages *int   `json:"targetPages"`
}

type ReadingHistoryImportRequest struct {
	Entries []ReadingSessionRequest `json:"entries"`
}

type CreateCommentRequest struct {
	CommentText string `json:"commentText"`
}

type CreateListRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	IsPrivate   bool    `json:"isPrivate"`
}

type UpdateListRequest = CreateListRequest

type UpsertListItemRequest struct {
	BookID   string  `json:"bookId"`
	Position int     `json:"position"`
	Notes    *string `json:"notes"`
}
