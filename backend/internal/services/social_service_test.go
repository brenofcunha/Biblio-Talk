package services

import (
	"context"
	"errors"
	"testing"
)

func TestFollowRejectsSelfBeforeRepositoryAccess(t *testing.T) {
	service := NewUserService(nil, nil)
	userID := "8b6a8a91-a36f-4eb6-9fac-96516893b31a"
	if _, err := service.Follow(context.Background(), userID, userID); !errors.Is(err, ErrSelfFollow) {
		t.Fatalf("Follow error = %v, want ErrSelfFollow", err)
	}
}

func TestFollowRejectsMalformedUserIDs(t *testing.T) {
	service := NewUserService(nil, nil)
	if _, err := service.Follow(context.Background(), "not-a-uuid", "another-id"); !errors.Is(err, ErrSocialUserNotFound) {
		t.Fatalf("Follow error = %v, want ErrSocialUserNotFound", err)
	}
}

func TestUnfollowRejectsSelfBeforeRepositoryAccess(t *testing.T) {
	service := NewUserService(nil, nil)
	userID := "8b6a8a91-a36f-4eb6-9fac-96516893b31a"
	if _, err := service.Unfollow(context.Background(), userID, userID); !errors.Is(err, ErrSelfFollow) {
		t.Fatalf("Unfollow error = %v, want ErrSelfFollow", err)
	}
}
