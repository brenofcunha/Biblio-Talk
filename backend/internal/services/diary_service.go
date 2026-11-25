package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"lume_api/internal/dto"
	"lume_api/internal/models"
	"lume_api/internal/repositories"
)

var (
	ErrInvalidReadingData = errors.New("dados de leitura inválidos")
	ErrInvalidDiaryData   = errors.New("dados do diário inválidos")
)

type DiaryService struct {
	repo *repositories.DiaryRepository
}

func NewDiaryService(repo *repositories.DiaryRepository) *DiaryService {
	return &DiaryService{repo: repo}
}

func validSessionStatus(status string) bool {
	switch status {
	case "want_to_read", "reading", "read", "abandoned":
		return true
	default:
		return false
	}
}

func parseDateField(value *string, field string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", *value)
	if err != nil {
		return nil, fmt.Errorf("%s deve estar no formato YYYY-MM-DD", field)
	}
	return &parsed, nil
}

func (s *DiaryService) sessionFromRequest(req dto.ReadingSessionRequest) (models.ReadingSession, error) {
	if strings.TrimSpace(req.BookID) == "" {
		return models.ReadingSession{}, fmt.Errorf("%w: bookId é obrigatório", ErrInvalidReadingData)
	}
	if req.Status == "" {
		req.Status = "reading"
	}
	if !validSessionStatus(req.Status) {
		return models.ReadingSession{}, fmt.Errorf("%w: status deve ser want_to_read, reading, read ou abandoned", ErrInvalidReadingData)
	}
	if req.CurrentPage < 0 {
		return models.ReadingSession{}, fmt.Errorf("%w: currentPage não pode ser negativo", ErrInvalidReadingData)
	}
	if req.TotalPages != nil && (*req.TotalPages <= 0 || req.CurrentPage > *req.TotalPages) {
		return models.ReadingSession{}, fmt.Errorf("%w: totalPages deve ser positivo e maior ou igual a currentPage", ErrInvalidReadingData)
	}
	startedAt, err := parseDateField(req.StartedAt, "startedAt")
	if err != nil {
		return models.ReadingSession{}, err
	}
	finishedAt, err := parseDateField(req.FinishedAt, "finishedAt")
	if err != nil {
		return models.ReadingSession{}, err
	}
	if startedAt != nil && finishedAt != nil && finishedAt.Before(*startedAt) {
		return models.ReadingSession{}, fmt.Errorf("%w: finishedAt não pode anteceder startedAt", ErrInvalidReadingData)
	}
	if req.Status == "reading" && startedAt == nil {
		today := time.Now().UTC().Truncate(24 * time.Hour)
		startedAt = &today
	}
	if req.Status == "read" && finishedAt == nil {
		today := time.Now().UTC().Truncate(24 * time.Hour)
		finishedAt = &today
	}
	return models.ReadingSession{
		BookID: req.BookID, Status: req.Status, CurrentPage: req.CurrentPage,
		TotalPages: req.TotalPages, StartedAt: startedAt, FinishedAt: finishedAt, IsReread: req.IsReread,
	}, nil
}

func (s *DiaryService) CreateReadingSession(ctx context.Context, userID string, req dto.ReadingSessionRequest) (*models.ReadingSession, error) {
	session, err := s.sessionFromRequest(req)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateReadingSession(ctx, userID, session)
}

func (s *DiaryService) GetReadingSession(ctx context.Context, userID, id string) (*models.ReadingSession, error) {
	return s.repo.GetReadingSession(ctx, userID, id)
}

func (s *DiaryService) ListReadingSessions(ctx context.Context, userID string, limit int) ([]models.ReadingSession, error) {
	return s.repo.ListReadingSessions(ctx, userID, limit)
}

func (s *DiaryService) UpdateReadingSession(ctx context.Context, userID, id string, req dto.ReadingSessionRequest) (*models.ReadingSession, error) {
	session, err := s.sessionFromRequest(req)
	if err != nil {
		return nil, err
	}
	current, err := s.repo.GetReadingSession(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if session.BookID != current.BookID {
		return nil, fmt.Errorf("%w: bookId não pode ser alterado na sessão", ErrInvalidReadingData)
	}
	return s.repo.UpdateReadingSession(ctx, userID, id, session)
}

func (s *DiaryService) DeleteReadingSession(ctx context.Context, userID, id string) error {
	return s.repo.DeleteReadingSession(ctx, userID, id)
}

func (s *DiaryService) SetReadingStatus(ctx context.Context, userID, id, status string) (*models.ReadingSession, error) {
	if !validSessionStatus(status) {
		return nil, fmt.Errorf("%w: status deve ser want_to_read, reading, read ou abandoned", ErrInvalidReadingData)
	}
	return s.repo.UpdateReadingStatus(ctx, userID, id, status)
}

func (s *DiaryService) RecordProgress(ctx context.Context, userID, sessionID string, req dto.ReadingProgressRequest) (*models.ReadingProgress, error) {
	if req.CurrentPage < 0 || req.PagesRead <= 0 {
		return nil, fmt.Errorf("%w: currentPage não pode ser negativo e pagesRead deve ser positivo", ErrInvalidReadingData)
	}
	progressDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("%w: date deve estar no formato YYYY-MM-DD", ErrInvalidReadingData)
	}
	if req.Note != nil && len(strings.TrimSpace(*req.Note)) > 500 {
		return nil, fmt.Errorf("%w: note excede 500 caracteres", ErrInvalidReadingData)
	}
	return s.repo.AddReadingProgress(ctx, userID, sessionID, models.ReadingProgress{
		ProgressDate: progressDate, CurrentPage: req.CurrentPage, PagesRead: req.PagesRead, Note: req.Note,
	})
}

func (s *DiaryService) ReadingProgress(ctx context.Context, userID, sessionID string) ([]models.ReadingProgress, error) {
	if _, err := s.repo.GetReadingSession(ctx, userID, sessionID); err != nil {
		return nil, err
	}
	return s.repo.ListReadingProgress(ctx, userID, sessionID)
}

func (s *DiaryService) ReadingHistory(ctx context.Context, userID, sessionID string) ([]models.ReadingSessionHistory, error) {
	if _, err := s.repo.GetReadingSession(ctx, userID, sessionID); err != nil {
		return nil, err
	}
	return s.repo.ListReadingHistory(ctx, userID, sessionID)
}

func (s *DiaryService) CreateNote(ctx context.Context, userID string, req dto.PrivateNoteRequest) (*models.PrivateNote, error) {
	if (req.BookID == nil) == (req.ReadingSessionID == nil) {
		return nil, fmt.Errorf("%w: informe bookId ou readingSessionId", ErrInvalidDiaryData)
	}
	if len(strings.TrimSpace(req.Content)) == 0 || len(req.Content) > 10000 {
		return nil, fmt.Errorf("%w: content é obrigatório e deve ter até 10000 caracteres", ErrInvalidDiaryData)
	}
	return s.repo.CreatePrivateNote(ctx, userID, models.PrivateNote{
		BookID: req.BookID, ReadingSessionID: req.ReadingSessionID, Content: req.Content,
	})
}

func (s *DiaryService) GetNote(ctx context.Context, userID, id string) (*models.PrivateNote, error) {
	return s.repo.GetPrivateNote(ctx, userID, id)
}

func (s *DiaryService) ListNotes(ctx context.Context, userID, query string, limit int) ([]models.PrivateNote, error) {
	return s.repo.ListPrivateNotes(ctx, userID, strings.TrimSpace(query), limit)
}

func (s *DiaryService) UpdateNote(ctx context.Context, userID, id, content string) (*models.PrivateNote, error) {
	if len(strings.TrimSpace(content)) == 0 || len(content) > 10000 {
		return nil, fmt.Errorf("%w: content é obrigatório e deve ter até 10000 caracteres", ErrInvalidDiaryData)
	}
	return s.repo.UpdatePrivateNote(ctx, userID, id, content)
}

func (s *DiaryService) DeleteNote(ctx context.Context, userID, id string) error {
	return s.repo.DeletePrivateNote(ctx, userID, id)
}

func (s *DiaryService) PublishNote(ctx context.Context, userID, id string) (*models.PrivateNote, error) {
	return s.repo.PublishPrivateNote(ctx, userID, id)
}

func (s *DiaryService) CreateNoteTag(ctx context.Context, userID, noteID, tag string) (*models.PrivateNoteTag, error) {
	tag = strings.TrimSpace(tag)
	if len(tag) == 0 || len(tag) > 50 {
		return nil, fmt.Errorf("%w: tag deve ter entre 1 e 50 caracteres", ErrInvalidDiaryData)
	}
	return s.repo.CreatePrivateNoteTag(ctx, userID, noteID, tag)
}

func (s *DiaryService) DeleteNoteTag(ctx context.Context, userID, noteID, tagID string) error {
	return s.repo.DeletePrivateNoteTag(ctx, userID, noteID, tagID)
}

func validateQuote(quote dto.NoteQuoteRequest) error {
	if len(strings.TrimSpace(quote.QuoteText)) == 0 || len(quote.QuoteText) > 5000 {
		return fmt.Errorf("%w: quoteText é obrigatório e deve ter até 5000 caracteres", ErrInvalidDiaryData)
	}
	if quote.Page != nil && *quote.Page <= 0 {
		return fmt.Errorf("%w: page deve ser positivo", ErrInvalidDiaryData)
	}
	if quote.Chapter != nil && len(*quote.Chapter) > 200 {
		return fmt.Errorf("%w: chapter excede 200 caracteres", ErrInvalidDiaryData)
	}
	return nil
}

func (s *DiaryService) CreateNoteQuote(ctx context.Context, userID, noteID string, req dto.NoteQuoteRequest) (*models.PrivateNoteQuote, error) {
	if err := validateQuote(req); err != nil {
		return nil, err
	}
	return s.repo.CreatePrivateNoteQuote(ctx, userID, noteID, models.PrivateNoteQuote{
		QuoteText: req.QuoteText, Page: req.Page, Chapter: req.Chapter,
	})
}

func (s *DiaryService) UpdateNoteQuote(ctx context.Context, userID, noteID, quoteID string, req dto.NoteQuoteRequest) (*models.PrivateNoteQuote, error) {
	if err := validateQuote(req); err != nil {
		return nil, err
	}
	return s.repo.UpdatePrivateNoteQuote(ctx, userID, noteID, quoteID, models.PrivateNoteQuote{
		QuoteText: req.QuoteText, Page: req.Page, Chapter: req.Chapter,
	})
}

func (s *DiaryService) DeleteNoteQuote(ctx context.Context, userID, noteID, quoteID string) error {
	return s.repo.DeletePrivateNoteQuote(ctx, userID, noteID, quoteID)
}

func (s *DiaryService) CreateReviewDraft(ctx context.Context, userID string, req dto.ReviewDraftRequest) (*models.Review, error) {
	if strings.TrimSpace(req.BookID) == "" || len(strings.TrimSpace(req.ReviewText)) > 10000 {
		return nil, fmt.Errorf("%w: bookId é obrigatório e reviewText deve ter até 10000 caracteres", ErrInvalidDiaryData)
	}
	if req.Rating != nil && !validateRating(*req.Rating) {
		return nil, fmt.Errorf("%w: rating deve ser entre 0.5 e 5.0 em passos de 0.5", ErrInvalidDiaryData)
	}
	return s.repo.CreateReviewDraft(ctx, userID, models.Review{
		BookID: req.BookID, ReadingSessionID: req.ReadingSessionID,
		Rating: req.Rating, ReviewText: req.ReviewText, ContainsSpoilers: req.ContainsSpoilers,
	})
}

func (s *DiaryService) GetReviewDraft(ctx context.Context, userID, id string) (*models.Review, error) {
	return s.repo.GetReviewDraft(ctx, userID, id)
}

func (s *DiaryService) GetPublishedReview(ctx context.Context, id string) (*models.Review, error) {
	return s.repo.GetPublishedReview(ctx, id)
}

func (s *DiaryService) UpdateReviewDraft(ctx context.Context, userID, id string, req dto.ReviewDraftRequest) (*models.Review, error) {
	if len(strings.TrimSpace(req.ReviewText)) > 10000 {
		return nil, fmt.Errorf("%w: reviewText deve ter até 10000 caracteres", ErrInvalidDiaryData)
	}
	if req.Rating != nil && !validateRating(*req.Rating) {
		return nil, fmt.Errorf("%w: rating deve ser entre 0.5 e 5.0 em passos de 0.5", ErrInvalidDiaryData)
	}
	current, err := s.repo.GetReviewDraft(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if req.BookID != "" && req.BookID != current.BookID {
		return nil, fmt.Errorf("%w: bookId não pode ser alterado no rascunho", ErrInvalidDiaryData)
	}
	return s.repo.UpdateReviewDraft(ctx, userID, id, models.Review{
		ReadingSessionID: req.ReadingSessionID, Rating: req.Rating,
		ReviewText: req.ReviewText, ContainsSpoilers: req.ContainsSpoilers,
	})
}

func (s *DiaryService) DeleteReview(ctx context.Context, userID, id string) error {
	return s.repo.DeleteReview(ctx, userID, id)
}

func (s *DiaryService) SetReviewPublished(ctx context.Context, userID, id string, published bool) (*models.Review, error) {
	return s.repo.SetReviewPublished(ctx, userID, id, published)
}

func (s *DiaryService) UpdateReview(ctx context.Context, userID, id string, req dto.UpdateReviewRequest) (*models.Review, error) {
	if len(strings.TrimSpace(req.ReviewText)) == 0 || len(req.ReviewText) > 10000 {
		return nil, fmt.Errorf("%w: reviewText é obrigatório e deve ter até 10000 caracteres", ErrInvalidDiaryData)
	}
	if req.Rating != nil && !validateRating(*req.Rating) {
		return nil, fmt.Errorf("%w: rating deve ser entre 0.5 e 5.0 em passos de 0.5", ErrInvalidDiaryData)
	}
	return s.repo.UpdateReview(ctx, userID, id, models.Review{
		Rating: req.Rating, ReviewText: req.ReviewText, ContainsSpoilers: req.ContainsSpoilers,
	})
}

func (s *DiaryService) UpdateReviewPrivacy(ctx context.Context, userID, id string, req dto.ReviewPrivacyRequest) (*models.Review, error) {
	if req.AllowComments == nil || req.AllowLikes == nil {
		return nil, fmt.Errorf("%w: allowComments e allowLikes são obrigatórios", ErrInvalidDiaryData)
	}
	return s.repo.UpdateReviewPrivacy(ctx, userID, id, *req.AllowComments, *req.AllowLikes)
}

func (s *DiaryService) ReviewHistory(ctx context.Context, userID, id string) ([]models.ReviewHistory, error) {
	return s.repo.ListReviewHistory(ctx, userID, id)
}

func (s *DiaryService) CreateReviewFromNote(ctx context.Context, userID, noteID string) (*models.Review, error) {
	return s.repo.CreateReviewFromNote(ctx, userID, noteID)
}

func (s *DiaryService) CreateReviewFromLog(ctx context.Context, userID, logID string) (*models.Review, error) {
	return s.repo.CreateReviewFromLog(ctx, userID, logID)
}

func (s *DiaryService) SaveReadingGoal(ctx context.Context, userID string, req dto.ReadingGoalRequest) (*models.ReadingGoal, error) {
	if req.Year < 2000 || req.Year > 2100 {
		return nil, fmt.Errorf("%w: year deve estar entre 2000 e 2100", ErrInvalidReadingData)
	}
	if (req.Type == "monthly" && (req.Month == nil || *req.Month < 1 || *req.Month > 12)) ||
		(req.Type == "annual" && req.Month != nil) || (req.Type != "monthly" && req.Type != "annual") {
		return nil, fmt.Errorf("%w: type deve ser monthly com month 1-12 ou annual sem month", ErrInvalidReadingData)
	}
	if req.TargetBooks == nil && req.TargetPages == nil {
		return nil, fmt.Errorf("%w: informe targetBooks e/ou targetPages", ErrInvalidReadingData)
	}
	if req.TargetBooks != nil && (*req.TargetBooks < 1 || *req.TargetBooks > 1000) {
		return nil, fmt.Errorf("%w: targetBooks deve estar entre 1 e 1000", ErrInvalidReadingData)
	}
	if req.TargetPages != nil && (*req.TargetPages < 1 || *req.TargetPages > 1000000) {
		return nil, fmt.Errorf("%w: targetPages deve estar entre 1 e 1000000", ErrInvalidReadingData)
	}
	return s.repo.UpsertReadingGoal(ctx, userID, models.ReadingGoal{
		Year: req.Year, Type: req.Type, Month: req.Month,
		TargetBooks: req.TargetBooks, TargetPages: req.TargetPages,
	})
}

func (s *DiaryService) ReadingGoals(ctx context.Context, userID string, year int) ([]models.ReadingGoal, error) {
	return s.repo.ListReadingGoals(ctx, userID, year)
}

func (s *DiaryService) ImportReadingHistory(ctx context.Context, userID string, req dto.ReadingHistoryImportRequest) (*models.ReadingHistoryImport, error) {
	if len(req.Entries) == 0 || len(req.Entries) > 500 {
		return nil, fmt.Errorf("%w: entries deve conter entre 1 e 500 leituras", ErrInvalidReadingData)
	}
	sessions := make([]models.ReadingSession, 0, len(req.Entries))
	for i, entry := range req.Entries {
		session, err := s.sessionFromRequest(entry)
		if err != nil {
			return nil, fmt.Errorf("entries[%d]: %w", i, err)
		}
		sessions = append(sessions, session)
	}
	return s.repo.CreateReadingHistoryImport(ctx, userID, sessions)
}

func (s *DiaryService) ReadingHistoryImport(ctx context.Context, userID, jobID string) (*models.ReadingHistoryImport, error) {
	return s.repo.GetReadingHistoryImport(ctx, userID, jobID)
}
