package domain

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrUsernamePasswordRequired = errors.New("username and password must be provided")
	ErrInvalidCredentials       = errors.New("invalid username or password")
	ErrInvalidToken             = errors.New("invalid token")
	ErrMultipleUser             = errors.New("multiple users found")
	ErrPasswordMismatch         = errors.New("user password mismatch")
)
