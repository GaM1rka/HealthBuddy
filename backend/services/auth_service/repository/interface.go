package repository

import (
	"context"

	"auth_service/domain"
)

// UserRepo
type UserRepo interface {
	Health(ctx context.Context) error
	Create(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, u *domain.User) error
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}
