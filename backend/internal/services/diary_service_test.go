package services

import (
	"errors"
	"testing"

	"lume_api/internal/dto"
)

func TestReadingSessionRequestValidatesStatusAndPageBounds(t *testing.T) {
	service := NewDiaryService(nil)

	_, err := service.sessionFromRequest(dto.ReadingSessionRequest{
		BookID: "book-id", Status: "finished", TotalPages: intPointer(12),
	})
	if !errors.Is(err, ErrInvalidReadingData) {
		t.Fatalf("invalid status error = %v, want ErrInvalidReadingData", err)
	}

	_, err = service.sessionFromRequest(dto.ReadingSessionRequest{
		BookID: "book-id", Status: "reading", CurrentPage: 13, TotalPages: intPointer(12),
	})
	if !errors.Is(err, ErrInvalidReadingData) {
		t.Fatalf("invalid page bounds error = %v, want ErrInvalidReadingData", err)
	}
}

func TestReadingSessionRequestDefaultsStatusAndStartDate(t *testing.T) {
	service := NewDiaryService(nil)

	session, err := service.sessionFromRequest(dto.ReadingSessionRequest{BookID: "book-id"})
	if err != nil {
		t.Fatalf("sessionFromRequest error = %v", err)
	}
	if session.Status != "reading" {
		t.Fatalf("status = %q, want reading", session.Status)
	}
	if session.StartedAt == nil {
		t.Fatal("startedAt is nil for an active session")
	}
}

func TestCreatePrivateNoteRequiresExactlyOneParent(t *testing.T) {
	service := NewDiaryService(nil)

	_, err := service.CreateNote(t.Context(), "user-id", dto.PrivateNoteRequest{Content: "note"})
	if !errors.Is(err, ErrInvalidDiaryData) {
		t.Fatalf("CreateNote error = %v, want ErrInvalidDiaryData", err)
	}
}

func TestSaveReadingGoalRejectsInvalidPeriodBeforePersistence(t *testing.T) {
	service := NewDiaryService(nil)

	_, err := service.SaveReadingGoal(t.Context(), "user-id", dto.ReadingGoalRequest{
		Year: 2026, Type: "monthly", TargetBooks: intPointer(12),
	})
	if !errors.Is(err, ErrInvalidReadingData) {
		t.Fatalf("SaveReadingGoal error = %v, want ErrInvalidReadingData", err)
	}
}

func intPointer(value int) *int { return &value }
