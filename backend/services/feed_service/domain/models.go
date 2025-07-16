package domain

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	validate      *validator.Validate
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
)

func init() {
	validate = validator.New()
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return usernameRegex.MatchString(fl.Field().String())
	})
}

// publication structure
type Publication struct {
	gorm.Model
	PostID  string `gorm:"type:char(36);uniqueIndex" json:"post_id" validate:"required,uuid4"`
	UserID  string `gorm:"type:char(36)" json:"user_id" validate:"required,uuid4"`
	Name    string `gorm:"size:30" json:"name"`
	Title   string `gorm:"size:100" json:"title"`
	Content string `gorm:"size:10000" json:"content"`
}

// comment structure
type Comment struct {
	gorm.Model
	CommentID string `gorm:"type:char(36);uniqueIndex" json:"comment_id" validate:"required,uuid4"`
	PostID    string `gorm:"type:char(36)" json:"post_id" validate:"required,uuid4"`
	UserID    string `gorm:"type:char(36)" json:"user_id" validate:"required,uuid4"`
	Name      string `gorm:"size:30" json:"name"`
	Content   string `gorm:"size:10000" json:"content" validate:"required,max=10000"`
}

// post/put publication requests
type PublicationRequest struct {
	Title   string `json:"title" validate:"required,max=300"`
	Content string `json:"content" validate:"required,max=10000"`
}

func (r *PublicationRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return fmt.Errorf("field %q failed on the %q tag", e.Field(), e.Tag())
		}
	}
	return nil
}

// get publication responce
type PublicationResponse struct {
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// post comment request
type PostCommentRequest struct {
	PostID  string `json:"post_id" validate:"required,uuid4"`
	Content string `json:"content" validate:"required,max=10000"`
}

func (r *PostCommentRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return fmt.Errorf("field %q failed on the %q tag", e.Field(), e.Tag())
		}
	}
	return nil
}

// put comment request
type PutCommentRequest struct {
	CommentID string `json:"comment_id" validate:"required,uuid4"`
	Content   string `json:"content" validate:"required,max=10000"`
}

func (r *PutCommentRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return fmt.Errorf("field %q failed on the %q tag", e.Field(), e.Tag())
		}
	}
	return nil
}

// comment response
type CommentResponse struct {
	CommentID string    `json:"comment_id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Converter
func (p *Publication) ToResponse() PublicationResponse {
	return PublicationResponse{
		PostID:    p.PostID,
		UserID:    p.UserID,
		Name:      p.Name,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
	}
}

func (p *Comment) ToResponse() CommentResponse {
	return CommentResponse{
		CommentID: p.CommentID,
		UserID:    p.UserID,
		Name:      p.Name,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
	}
}

func NewUUID() string {
	return uuid.New().String()
}
