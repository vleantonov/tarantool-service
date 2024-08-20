package domain

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrUsernamePasswordRequired = errors.New("username and password must be provided")
	ErrDataRequired             = errors.New(`field "data" must be specified and must be not empty`)
	ErrInvalidCredentials       = errors.New("invalid username or password")
	ErrInvalidToken             = errors.New("invalid token")
	ErrMultipleUser             = errors.New("multiple users found")
	ErrPasswordMismatch         = errors.New("user password mismatch")
)
