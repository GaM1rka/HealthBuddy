package repository

import (
	"context"

	"auth_service/domain"
)

// UserRepo
type UserRepo interface {
	Health(ctx context.Context) error
	Create(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, UserID string) error
	FindByUserName(ctx context.Context, username string) (*domain.User, error)
	FindByUserID(ctx context.Context, userID string) (*domain.User, error)
}
