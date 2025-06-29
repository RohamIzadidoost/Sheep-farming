package ports

import (
	"context"
	"sheep_farm_backend_go/internal/domain"
)

// UserRepository defines the interface for user data operations.
// This is a "driven port" (output port).
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userID string) error
}
