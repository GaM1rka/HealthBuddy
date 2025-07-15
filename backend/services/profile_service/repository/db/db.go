package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"profile_service/cmd/config"
	"profile_service/domain"
	repository "profile_service/repository"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type pgProfileRepo struct {
	db *gorm.DB
}

func NewProfileRepo(db *gorm.DB) repository.ProfileRepository {
	return &pgProfileRepo{db: db}
}

func (r *pgProfileRepo) Health(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (r *pgProfileRepo) Create(p *domain.Profile) error {
	return r.db.Create(p).Error
}

func (r *pgProfileRepo) GetByUserID(userID string) (*domain.Profile, error) {
	var p domain.Profile
	if err := r.db.Where("user_id = ?", userID).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgProfileRepo) Update(p *domain.Profile) error {
	return r.db.Save(p).Error
}

func (r *pgProfileRepo) Delete(userID string) error {
	return r.db.
		Unscoped().
		Where("user_id = ?", userID).
		Delete(&domain.Profile{}).
		Error
}

func (r *pgProfileRepo) List() ([]domain.Profile, error) {
	var profiles []domain.Profile
	if err := r.db.Find(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
}

func InitDB(dbConfig config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository.ErrDBConnection, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository.ErrDBConnection, err)
	}

	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("%w: %v", repository.ErrDBConnection, err)
	}

	if err := db.AutoMigrate(&domain.Profile{}); err != nil {
		return nil, fmt.Errorf("%w: %v", repository.ErrDBMigration, err)
	}

	log.Println("Database connection established and migrations applied")
	return db, nil
}
