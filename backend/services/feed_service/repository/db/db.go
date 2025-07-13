package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"feed_service/cmd/config"
	"feed_service/domain"
	repository "feed_service/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// pgFeedRepo handles publications and comments
type pgFeedRepo struct {
	db *gorm.DB
}

// NewFeedRepo constructor
func NewFeedRepo(db *gorm.DB) repository.FeedRepository {
	return &pgFeedRepo{db: db}
}

func (r *pgFeedRepo) Health(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// CreatePublication inserts a new post
func (r *pgFeedRepo) CreatePublication(p *domain.Publication) error {
	return r.db.Create(p).Error
}

// GetPublication by post ID
func (r *pgFeedRepo) GetPublication(postID string) (*domain.Publication, error) {
	var p domain.Publication
	err := r.db.Where("post_id = ?", postID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	return &p, err
}

// ListPublications returns all posts
func (r *pgFeedRepo) ListPublications() ([]domain.Publication, error) {
	var pubs []domain.Publication
	err := r.db.Find(&pubs).Error
	return pubs, err
}

// UpdatePublication updates post content or name
func (r *pgFeedRepo) UpdatePublication(p *domain.Publication) error {
	return r.db.Save(p).Error
}

// DeletePublication soft-deletes a post by ID
func (r *pgFeedRepo) DeletePublication(postID string) error {
	return r.db.Unscoped().Where("post_id = ?", postID).Delete(&domain.Publication{}).Error
}

// CreateComment on a post
func (r *pgFeedRepo) CreateComment(c *domain.Comment) error {
	return r.db.Create(c).Error
}

// ListComments returns comments for a post
func (r *pgFeedRepo) ListComments(postID string) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := r.db.Where("post_id = ?", postID).Find(&comments).Error
	return comments, err
}

// GetComment by comment ID
func (r *pgFeedRepo) GetComment(commentID string) (*domain.Comment, error) {
	var c domain.Comment
	err := r.db.Where("comment_id = ?", commentID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	return &c, err
}

// UpdateComment updates comment content
func (r *pgFeedRepo) UpdateComment(c *domain.Comment) error {
	return r.db.Save(c).Error
}

// DeleteComment removes a comment by ID
func (r *pgFeedRepo) DeleteComment(commentID string) error {
	return r.db.Unscoped().Where("comment_id = ?", commentID).Delete(&domain.Comment{}).Error
}

func InitDB(dbConfig config.DBConfig) (*gorm.DB, error) {
	// Build DSN string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name,
	)

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// Open the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository.ErrDBConnection, err)
	}

	// Get raw DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository.ErrDBConnection, err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("%w: %v", repository.ErrDBConnection, err)
	}

	// Run migrations for all domain models
	if err := db.AutoMigrate(
		&domain.Publication{},
		&domain.Comment{},
	); err != nil {
		return nil, fmt.Errorf("%w: %v", repository.ErrDBMigration, err)
	}

	log.Println("Database connection established and migrations applied")
	return db, nil
}

func (r *pgFeedRepo) ListPublicationsByUser(ctx context.Context, userID string) ([]domain.Publication, error) {
	var pubs []domain.Publication
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&pubs).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return pubs, nil
}
