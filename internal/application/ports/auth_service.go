package ports

import (
	"context"
	"sheep_farm_backend_go/internal/domain"
)

// AuthService defines the interface for authentication operations.
// This is a "driving port" (input port) from the perspective of the web layer,
// and also an "output port" if it interacts with external auth providers.
type AuthService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, error) // Returns JWT token
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
}
