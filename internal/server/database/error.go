package database

import "errors"

// Describes possible database error.
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrDatabase          = errors.New("database error")
)
