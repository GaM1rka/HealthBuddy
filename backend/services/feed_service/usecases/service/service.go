package service

import (
	"context"
	"errors"
	"sort"

	"feed_service/domain"
	repository "feed_service/repository"
	"feed_service/usecases"
)

// feedService implements usecases.FeedService
// Handles publications and comments business logic
// along with health check capability.

type feedService struct {
	repository repository.FeedRepository
}

// NewFeedService creates a new FeedService
func NewFeedService(repo repository.FeedRepository) usecases.FeedService {
	return &feedService{repository: repo}
}

// Health performs a health check
func (s *feedService) Health(ctx context.Context) error {
	return s.repository.Health(ctx)
}

// CreatePublication creates a new post
func (s *feedService) CreatePublication(ctx context.Context, userID string, req domain.PublicationRequest) (domain.PublicationResponse, error) {
	pub := &domain.Publication{
		PostID:  domain.NewUUID(),
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}
	if err := s.repository.CreatePublication(pub); err != nil {
		return domain.PublicationResponse{}, err
	}
	return pub.ToResponse(), nil
}

// GetPublication returns a post by ID
func (s *feedService) GetPublication(ctx context.Context, postID string) (domain.PublicationResponse, error) {
	pub, err := s.repository.GetPublication(postID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.PublicationResponse{}, usecases.ErrNotFound
		}
		return domain.PublicationResponse{}, err
	}
	return pub.ToResponse(), nil
}

// ListPublications lists all posts sorted by newest first
func (s *feedService) ListPublications(ctx context.Context) ([]domain.PublicationResponse, error) {
	pubs, err := s.repository.ListPublications()
	if err != nil {
		return nil, err
	}
	// most fresh publications first
	sort.Slice(pubs, func(i, j int) bool {
		return pubs[i].CreatedAt.After(pubs[j].CreatedAt)
	})

	out := make([]domain.PublicationResponse, len(pubs))
	for i, p := range pubs {
		out[i] = p.ToResponse()
	}
	return out, nil
}

// UpdatePublication updates post fields
func (s *feedService) UpdatePublication(ctx context.Context, postID string, req domain.PublicationRequest) (domain.PublicationResponse, error) {
	pub, err := s.repository.GetPublication(postID)
	if err != nil {
		return domain.PublicationResponse{}, err
	}
	pub.Title = req.Title
	pub.Content = req.Content
	if err := s.repository.UpdatePublication(pub); err != nil {
		return domain.PublicationResponse{}, err
	}
	return pub.ToResponse(), nil
}

// DeletePublication permanently deletes a post
func (s *feedService) DeletePublication(ctx context.Context, postID string) error {
	return s.repository.DeletePublication(postID)
}

// CreateComment adds a new comment to a post
func (s *feedService) CreateComment(ctx context.Context, userID string, req domain.PostCommentRequest) (domain.CommentResponse, error) {
	comment := &domain.Comment{
		CommentID: domain.NewUUID(),
		PostID:    req.PostID,
		UserID:    userID,
		Content:   req.Content,
	}
	if err := s.repository.CreateComment(comment); err != nil {
		return domain.CommentResponse{}, err
	}
	return comment.ToResponse(), nil
}

// GetComment returns a comment by ID
func (s *feedService) GetComment(ctx context.Context, commentID string) (domain.CommentResponse, error) {
	c, err := s.repository.GetComment(commentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.CommentResponse{}, usecases.ErrNotFound
		}
		return domain.CommentResponse{}, err
	}
	return c.ToResponse(), nil
}

// ListComments lists comments for a post
func (s *feedService) ListComments(ctx context.Context, postID string) ([]domain.CommentResponse, error) {
	comments, err := s.repository.ListComments(postID)
	if err != nil {
		return nil, err
	}
	//most fresh comments
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.After(comments[j].CreatedAt)
	})
	out := make([]domain.CommentResponse, len(comments))
	for i, c := range comments {
		out[i] = c.ToResponse()
	}
	return out, nil
}

// UpdateComment updates comment content
func (s *feedService) UpdateComment(ctx context.Context, req domain.PutCommentRequest) (domain.CommentResponse, error) {
	c, err := s.repository.GetComment(req.CommentID)
	if err != nil {
		return domain.CommentResponse{}, err
	}
	c.Content = req.Content
	if err := s.repository.UpdateComment(c); err != nil {
		return domain.CommentResponse{}, err
	}
	return c.ToResponse(), nil
}

// DeleteComment deletes a comment
func (s *feedService) DeleteComment(ctx context.Context, commentID string) error {
	return s.repository.DeleteComment(commentID)
}

// ListPublicationsByUser returns all publications for a given user, newest first.
func (s *feedService) ListPublicationsByUser(ctx context.Context, userID string) ([]domain.PublicationResponse, error) {
	// Retrieve domain publications from the repository
	pubs, err := s.repository.ListPublicationsByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, usecases.ErrNotFound
		}
		return nil, err
	}

	// Map domain.Publication to PublicationResponse
	out := make([]domain.PublicationResponse, len(pubs))
	for i, p := range pubs {
		out[i] = p.ToResponse()
	}

	return out, nil
}
