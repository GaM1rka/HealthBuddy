package service

import (
	"context"

	"profile_service/domain"
	"profile_service/repository"
	"profile_service/usecases"
)

// profileService
type profileService struct {
	repo repository.ProfileRepository
}

// NewProfileService
func NewProfileService(repo repository.ProfileRepository) usecases.ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) Health(ctx context.Context) error {
	return s.repo.Health(ctx)
}

func (s *profileService) CreateProfile(ctx context.Context, userID string, req domain.ProfileRequest) (domain.ProfileResponse, error) {
	p := domain.Profile{
		UserID:   userID,
		Username: req.Name,
		Bio:      req.Bio,
		Avatar:   req.Avatar,
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
	return p.ToResponse(), nil
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
	p.Username = req.Name
	p.Bio = req.Bio
	p.Avatar = req.Avatar
	if err := s.repo.Update(p); err != nil {
		return domain.ProfileResponse{}, err
	}
	return p.ToResponse(), nil
}

func (s *profileService) DeleteProfile(ctx context.Context, userID string) error {
	return s.repo.Delete(userID)
}
