package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"auth_service/cmd/config"
	"auth_service/domain"
	"auth_service/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type userRepo struct {
	db *gorm.DB
}

// NewPostgresUserRepo -> UserRepo
func NewUserRepo(db *gorm.DB) repository.UserRepo {
	return &userRepo{db: db}
}

// health
func (r *userRepo) Health(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Create user
func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepo) Delete(ctx context.Context, u *domain.User) error {
	result := r.db.WithContext(ctx).
		Unscoped().
		Where("email = ?", u.Email).
		Or("username = ?", u.Username).
		Delete(&domain.User{})
	return result.Error
}

// Search user by username
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&u).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// InitDB
func InitDB(dbConfig config.DBConfig) (*gorm.DB, error) {
	// DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name,
	)
	// logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}
