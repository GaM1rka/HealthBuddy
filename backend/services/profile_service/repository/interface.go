package repository

import (
	"context"
	"profile_service/domain"
)

type ProfileRepository interface {
	Create(profile *domain.Profile) error
	GetByUserID(userID string) (*domain.Profile, error)
	Update(profile *domain.Profile) error
	Delete(userID string) error
	List() ([]domain.Profile, error)
	Health(ctx context.Context) error
}
