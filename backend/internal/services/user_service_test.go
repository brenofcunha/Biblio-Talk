package services

import (
	"context"
	"errors"
	"testing"

	"lume_api/internal/models"
)

func TestUpdatePrivacyRejectsMissingOrUnsupportedVisibility(t *testing.T) {
	service := NewUserService(nil, nil)
	settings := models.PrivacySettings{
		Profile: "public", Diary: "private", Progress: "public", Shelf: "public",
		Notes: "private", Comments: "public", Lists: "public", Stats: "public", Favorites: "public",
	}
	settings.Comments = "followers"
	if err := service.UpdatePrivacy(context.Background(), "user-id", settings); !errors.Is(err, ErrInvalidPrivacy) {
		t.Fatalf("UpdatePrivacy error = %v, want ErrInvalidPrivacy", err)
	}
}

func TestUpdateFavoritesValidatesCountDuplicatesAndIDs(t *testing.T) {
	service := NewUserService(nil, nil)
	tests := []struct {
		name string
		ids  []string
	}{
		{name: "more than four", ids: []string{
			"00000000-0000-4000-8000-000000000001",
			"00000000-0000-4000-8000-000000000002",
			"00000000-0000-4000-8000-000000000003",
			"00000000-0000-4000-8000-000000000004",
			"00000000-0000-4000-8000-000000000005",
		}},
		{name: "duplicate", ids: []string{"00000000-0000-4000-8000-000000000001", "00000000-0000-4000-8000-000000000001"}},
		{name: "invalid id", ids: []string{"not-a-uuid"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.UpdateFavorites(context.Background(), "user-id", test.ids)
			if !errors.Is(err, ErrInvalidProfile) {
				t.Fatalf("UpdateFavorites error = %v, want ErrInvalidProfile", err)
			}
		})
	}
}

func TestRecordConsentRequiresKnownTypeAndVersion(t *testing.T) {
	service := NewUserService(nil, nil)
	_, err := service.RecordConsent(context.Background(), "user-id", models.Consent{ConsentType: "unknown", Version: "1", Granted: true})
	if !errors.Is(err, ErrInvalidConsent) {
		t.Fatalf("RecordConsent error = %v, want ErrInvalidConsent", err)
	}
	_, err = service.RecordConsent(context.Background(), "user-id", models.Consent{ConsentType: "privacy_policy", Version: " ", Granted: true})
	if !errors.Is(err, ErrInvalidConsent) {
		t.Fatalf("RecordConsent error = %v, want ErrInvalidConsent", err)
	}
}

func TestSaveAnnualChallengeValidatesBounds(t *testing.T) {
	service := NewUserService(nil, nil)
	if err := service.SaveAnnualChallenge(context.Background(), "user-id", models.AnnualChallenge{Year: 1999, TargetBooks: 10}); !errors.Is(err, ErrInvalidChallenge) {
		t.Fatalf("SaveAnnualChallenge error = %v, want ErrInvalidChallenge", err)
	}
	if err := service.SaveAnnualChallenge(context.Background(), "user-id", models.AnnualChallenge{Year: 2026, TargetBooks: 0}); !errors.Is(err, ErrInvalidChallenge) {
		t.Fatalf("SaveAnnualChallenge error = %v, want ErrInvalidChallenge", err)
	}
}
