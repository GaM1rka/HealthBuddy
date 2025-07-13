package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"profile_service/domain"
	"profile_service/repository"
	"profile_service/usecases"
)

// profileService
type profileService struct {
	repo       repository.ProfileRepository
	authUrl    string
	feedUrl    string
	httpClient *http.Client
}

// NewProfileService
func NewProfileService(repo repository.ProfileRepository, authUrl string, feedUrl string) usecases.ProfileService {
	return &profileService{
		repo:       repo,
		authUrl:    authUrl,
		feedUrl:    feedUrl,
		httpClient: http.DefaultClient,
	}
}

func (s *profileService) Health(ctx context.Context) error {
	return s.repo.Health(ctx)
}

func (s *profileService) CreateProfile(ctx context.Context, userID string, req domain.ProfileRequest) (domain.ProfileResponse, error) {
	name := ""
	if req.Name != nil {
		name = *req.Name
	}
	bio := ""
	if req.Bio != nil {
		bio = *req.Bio
	}
	avatar := ""
	if req.Avatar != nil {
		avatar = *req.Avatar
	}

	p := domain.Profile{
		UserID: userID,
		Name:   name,
		Bio:    bio,
		Avatar: avatar,
	}
	if err := s.repo.Create(&p); err != nil {
		return domain.ProfileResponse{}, err
	}
	return p.ToResponse(), nil
}

func (s *profileService) GetProfile(ctx context.Context, userID string) (domain.ProfileResponse, error) {
	p, err := s.repo.GetByUserID(userID)
	if err != nil {
		return domain.ProfileResponse{}, err
	}
	resp := p.ToResponse()
	feedURL := fmt.Sprintf("%s/feed/user/publications", s.feedUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		log.Printf("GetProfile: build feed request failed: %v", err)
		return resp, nil
	}
	req.Header.Set("X-User-ID", userID)
	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("GetProfile: feed request failed: %v", err)
		return resp, nil
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		var posts []domain.PublicationResponse
		if err := json.NewDecoder(res.Body).Decode(&posts); err != nil {
			log.Printf("GetProfile: decode feed response failed: %v", err)
		} else {
			resp.Posts = posts
		}
	} else {
		log.Printf("GetProfile: feed service returned %d", res.StatusCode)
	}

	return resp, nil
}

func (s *profileService) ListProfiles(ctx context.Context) ([]domain.ProfileResponse, error) {
	profiles, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	out := make([]domain.ProfileResponse, len(profiles))
	for i, p := range profiles {
		out[i] = p.ToResponse()
	}
	return out, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, userID string, req domain.ProfileRequest) (domain.ProfileResponse, error) {
	p, err := s.repo.GetByUserID(userID)
	if err != nil {
		return domain.ProfileResponse{}, err
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Bio != nil {
		p.Bio = *req.Bio
	}
	if req.Avatar != nil {
		p.Avatar = *req.Avatar
	}
	if err := s.repo.Update(p); err != nil {
		return domain.ProfileResponse{}, err
	}
	return p.ToResponse(), nil
}

func (s *profileService) DeleteProfile(ctx context.Context, userID string) error {
	if err := s.repo.Delete(userID); err != nil {
		return err
	}
	// DELETE к Auth-сервису
	// route DELETE /auth/users/:id
	authEndpoint := fmt.Sprintf("%s/auth/user/%s", s.authUrl, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, authEndpoint, nil)
	if err != nil {
		log.Printf("failed to build auth delete request: %v", err)
		return nil
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("failed to call auth service delete: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		log.Printf("auth service delete returned status %d for user %s", resp.StatusCode, userID)
	}
	return nil
}
