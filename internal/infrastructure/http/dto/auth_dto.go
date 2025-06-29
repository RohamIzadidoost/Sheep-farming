package dto

import "sheep_farm_backend_go/internal/domain"

// RegisterRequest represents the request body for user registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the response body after successful registration or login.
type AuthResponse struct {
	UserID string          `json:"userId"`
	Email  string          `json:"email"`
	Role   domain.UserRole `json:"role"`
	Token  string          `json:"token,omitempty"` // JWT token
}
