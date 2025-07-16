package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lib/pq"
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
	// validate incoming credentials
	if err := creds.Validate(); err != nil {
		return "", err
	}

	// hash password
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
	// attempt to create user; map unique‑constraint errors to ErrEmailTaken
	if err := s.repo.Create(ctx, user); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return "", usecases.ErrEmailTaken
		}
		return "", fmt.Errorf("repo create: %w", err)
	}

	userID := user.ID

	// build profile payload
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

	// call profile service
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/profile", s.profileSvcURL), bytes.NewReader(body))
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

	// generate JWT
	token, err := jwt.GenerateToken(userID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (s *authService) Login(ctx context.Context, creds domain.LoginCredentials) (string, error) {
	// find user by username
	user, err := s.repo.FindByUserName(ctx, creds.Username)
	if err != nil {
		// hide "not found" vs wrong password: always return invalid credentials
		return "", usecases.ErrInvalidCredentials
	}

	// compare passwords
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		return "", usecases.ErrInvalidCredentials
	}

	return jwt.GenerateToken(user.ID)
}

func (s *authService) Health(ctx context.Context) error {
	return s.repo.Health(ctx)
}

func (s *authService) DeleteUser(ctx context.Context, userID string) error {
	return s.repo.Delete(ctx, userID)
}

// FindByID
func (s *authService) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	return s.repo.FindByUserID(ctx, userID)
}
