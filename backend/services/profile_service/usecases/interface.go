package usecases

import (
	"context"

	"profile_service/domain"
	"profile_service/repository"
)

var ErrNotFound = repository.ErrNotFound

type ProfileService interface {
	Health(ctx context.Context) error
	CreateProfile(ctx context.Context, userID string, req domain.ProfileRequest) (domain.ProfileResponse, error)
	GetProfile(ctx context.Context, userID string) (domain.ProfileResponse, error)
	ListProfiles(ctx context.Context) ([]domain.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req domain.ProfileRequest) (domain.ProfileResponse, error)
	DeleteProfile(ctx context.Context, userID string) error
}
