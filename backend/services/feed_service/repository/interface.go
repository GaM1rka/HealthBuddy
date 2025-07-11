package repository

import (
	"context"

	"feed_service/domain"
)

// FeedRepository defines methods for managing publications and comments
// along with health check capability.
type FeedRepository interface {
	CreatePublication(pub *domain.Publication) error
	GetPublication(postID string) (*domain.Publication, error)
	ListPublications() ([]domain.Publication, error)
	UpdatePublication(pub *domain.Publication) error
	DeletePublication(postID string) error

	CreateComment(cmt *domain.Comment) error
	ListComments(postID string) ([]domain.Comment, error)
	GetComment(commentID string) (*domain.Comment, error)
	UpdateComment(cmt *domain.Comment) error
	DeleteComment(commentID string) error
	Health(ctx context.Context) error
}
