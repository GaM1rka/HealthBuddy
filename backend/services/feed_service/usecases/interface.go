package usecases

import (
	"context"
	"errors"

	"feed_service/domain"
)

// FeedService defines business logic for publications and comments
// along with health check capability.
type FeedService interface {
	// Health performs a health check
	Health(ctx context.Context) error

	// Publication operations
	CreatePublication(ctx context.Context, userID string, req domain.PublicationRequest) (domain.PublicationResponse, error)
	GetPublication(ctx context.Context, postID string) (domain.PublicationResponse, error)
	ListPublications(ctx context.Context) ([]domain.PublicationResponse, error)
	UpdatePublication(ctx context.Context, postID string, req domain.PublicationRequest) (domain.PublicationResponse, error)
	DeletePublication(ctx context.Context, postID string) error

	// Comment operations
	CreateComment(ctx context.Context, userID string, req domain.PostCommentRequest) (domain.CommentResponse, error)
	GetComment(ctx context.Context, commentID string) (domain.CommentResponse, error)
	ListComments(ctx context.Context, postID string) ([]domain.CommentResponse, error)
	UpdateComment(ctx context.Context, req domain.PutCommentRequest) (domain.CommentResponse, error)
	DeleteComment(ctx context.Context, commentID string) error
}

// error
var (
	ErrNotFound = errors.New("recording not found")
)
