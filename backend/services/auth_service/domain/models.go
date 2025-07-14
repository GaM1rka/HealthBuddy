package domain

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// Credentials holds login or registration data
type RegistrCredentials struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// User represents a registered user
type User struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Username  string    `gorm:"size:30;uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TokenResponse is returned after successful authentication
type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// Validate Credentials
func (c *RegistrCredentials) Validate() error {
	if len(c.Username) < 3 || len(c.Username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}
	if len(c.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if !emailRegex.MatchString(c.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

func (c *LoginCredentials) Validate() error {
	if len(c.Username) < 3 || len(c.Username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}
	if len(c.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

// Validate User
func (u *User) Validate() error {
	if u.ID == "" {
		return errors.New("user ID is required")
	}
	if len(u.Username) < 3 || len(u.Username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}
	if u.Password == "" {
		return errors.New("password hash is required")
	}
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NewUUID() string {
	return uuid.New().String()
}
