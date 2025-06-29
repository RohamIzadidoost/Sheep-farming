package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrAlreadyExists      = errors.New("already exists")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInternal           = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid email or password") // New: For login failures
	ErrEmailAlreadyExists = errors.New("email already exists")      // New: For registration
	ErrInvalidToken       = errors.New("invalid or expired token")  // New: For JWT validation
	ErrPermissionDenied   = errors.New("permission denied")         // New: For role-based access control
)
