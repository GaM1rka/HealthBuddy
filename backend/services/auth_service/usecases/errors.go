package usecases

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrEmailTaken         = errors.New("email/username already in use")
	ErrProfileServiceDown = errors.New("cannot create profile")
)
