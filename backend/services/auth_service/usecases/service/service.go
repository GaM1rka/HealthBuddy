package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"auth_service/domain"
	"auth_service/jwt"
	"auth_service/repository"
	"auth_service/usecases"
)

type authService struct {
	repo          repository.UserRepo
	profileSvcURL string
}

func NewAuthService(r repository.UserRepo, profileSvcURL string) usecases.AuthService {
	return &authService{
		repo:          r,
		profileSvcURL: profileSvcURL,
	}
}

func (s *authService) Register(ctx context.Context, creds domain.RegistrCredentials) (string, error) {
	if err := creds.Validate(); err != nil {
		return "", err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		ID:       domain.NewUUID(),
		Username: creds.Username,
		Password: string(hashed),
		Email:    creds.Email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return "", usecases.ErrEmailTaken
		}
		return "", err
	}

	userID := user.ID
	type profilePayload struct {
		Name   string `json:"name"`
		Bio    string `json:"bio"`
		Avatar string `json:"avatar"`
	}
	payload := profilePayload{Name: creds.Username}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.profileSvcURL+"/profile", bytes.NewReader(body))
	if err != nil {
		_ = s.repo.Delete(ctx, userID)
		return "", fmt.Errorf("new profile request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		_ = s.repo.Delete(ctx, userID)
		return "", usecases.ErrProfileServiceDown
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		_ = s.repo.Delete(ctx, userID)
		return "", usecases.ErrProfileServiceDown
	}

	// 5) генерируем JWT
	token, err := jwt.GenerateToken(userID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}

func (s *authService) Login(ctx context.Context, creds domain.LoginCredentials) (string, error) {
	user, err := s.repo.FindByUserName(ctx, creds.Username)
	if err != nil {
		return "", usecases.ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		return "", usecases.ErrInvalidCredentials
	}
	return jwt.GenerateToken(user.ID)
}

func (s *authService) Health(ctx context.Context) error {
	return s.repo.Health(ctx)
}

func (s *authService) DeleteUser(ctx context.Context, userID string) error {
	// delete user; if not exist map to ErrUserNotFound
	if err := s.repo.Delete(ctx, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return usecases.ErrUserNotFound
		}
		return err
	}
	return nil
}

// FindByID
func (s *authService) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	// find user; if not exist map to ErrUserNotFound
	user, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, usecases.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
