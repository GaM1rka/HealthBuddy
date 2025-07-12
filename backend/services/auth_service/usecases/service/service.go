package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

func (s *authService) Register(ctx context.Context, creds domain.Credentials) (string, error) {
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
		return "", err
	}

	userID := user.ID

	type profilePayload struct {
		Name   string `json:"name"`
		Bio    string `json:"bio"`
		Avatar string `json:"avatar"`
	}
	payload := profilePayload{
		Name:   creds.Username,
		Bio:    "",
		Avatar: "",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/profile", s.profileSvcURL), bytes.NewReader(body))
	if err != nil {
		_ = s.repo.Delete(ctx, user)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		_ = s.repo.Delete(ctx, user)
		return "", fmt.Errorf("failed to create profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		_ = s.repo.Delete(ctx, user)
		return "", fmt.Errorf("profile service returned status %d", resp.StatusCode)
	}

	token, err := jwt.GenerateToken(userID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (s *authService) Login(ctx context.Context, creds domain.Credentials) (string, error) {
	user, err := s.repo.FindByUsername(ctx, creds.Username)
	if err != nil {
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		return "", errors.New("invalid credentials")
	}

	return jwt.GenerateToken(user.ID)
}

func (s *authService) Health(ctx context.Context) error {
	return s.repo.Health(ctx)
}
