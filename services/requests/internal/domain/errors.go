package domain

import "errors"

var (
	ErrInvalidStatus = errors.New("invalid Status")
	ErrUserNotFound  = errors.New("user not found")
)
