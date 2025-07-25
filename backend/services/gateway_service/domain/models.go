package domain

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// main structure
type Profile struct {
	gorm.Model
	UserID   string `gorm:"type:char(36);uniqueIndex" json:"user_id" validate:"required,uuid4"`
	Username string `gorm:"size:30" json:"username" validate:"username"`
	Bio      string `gorm:"size:500" json:"bio,omitempty" validate:"max=500"`
	Avatar   string `json:"avatar,omitempty" validate:"avatar_url"`
}

// post/put requests
type ProfileRequest struct {
	Name   string `json:"name" validate:"required,username"`
	Bio    string `json:"bio" validate:"max=500"`
	Avatar string `json:"avatar" validate:"avatar_url"`
}

// get request
type ProfileResponse struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Bio       string    `json:"bio,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Converter
func (p *Profile) ToResponse() ProfileResponse {
	return ProfileResponse{
		UserID:    p.UserID,
		Username:  p.Username,
		Bio:       p.Bio,
		Avatar:    p.Avatar,
		CreatedAt: p.CreatedAt,
	}
}

var (
	validate      *validator.Validate
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
)

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("username", validateUsername)
	_ = validate.RegisterValidation("avatar_url", validateAvatarURL)
}

// validation username
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return len(username) >= 3 && len(username) <= 30 &&
		usernameRegex.MatchString(username)
}

// validation URL ava
func validateAvatarURL(fl validator.FieldLevel) bool {
	avatarURL := fl.Field().String()
	if avatarURL == "" {
		return true // omitempty
	}

	u, err := url.Parse(avatarURL)
	return err == nil &&
		(u.Scheme == "http" || u.Scheme == "https") &&
		strings.HasSuffix(u.Path, ".jpg") ||
		strings.HasSuffix(u.Path, ".png") ||
		strings.HasSuffix(u.Path, ".jpeg")
}

// Validate main structure
func Validate(s interface{}) error {
	return validate.Struct(s)
}
