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

type UpdateMeRequest struct {
	Bio         *string  `json:"bio"`
	AvatarURL   *string  `json:"avatarUrl"`
	Top4BookIDs []string `json:"top4BookIds"`
}

type ShelfUpsertRequest struct {
	BookID string `json:"bookId"`
	Status string `json:"status"`
}

type CreateLogRequest struct {
	BookID            string   `json:"bookId"`
	Rating            float64  `json:"rating"`
	ReviewText        *string  `json:"reviewText"`
	ReadDate          *string  `json:"readDate"`
	IsReread          bool     `json:"isReread"`
	ContainsSpoilers  bool     `json:"containsSpoilers"`
	Tags              []string `json:"tags"`
	Liked             bool     `json:"liked"`
}

type UpdateLogRequest = CreateLogRequest

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
