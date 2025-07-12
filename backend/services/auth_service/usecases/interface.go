package usecases

import (
	"auth_service/domain"
	"context"
)

// AuthService
type AuthService interface {
	Register(ctx context.Context, creds domain.Credentials) (string, error)
	Login(ctx context.Context, creds domain.Credentials) (string, error)
	Health(ctx context.Context) error
}
