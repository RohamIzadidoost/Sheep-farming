package services

import (
	"context"
	"time"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"

	"golang.org/x/crypto/bcrypt" // For password hashing
)

// UserService provides use cases for user management.
// This is an "application service" (use case).
type UserService struct {
	userRepo ports.UserRepository
}

// NewUserService creates a new UserService instance.
func NewUserService(userRepo ports.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser handles the creation of a new user.
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	// Check if user with this email already exists
	existingUser, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil && err != domain.ErrNotFound {
		return err // Other database error
	}
	if existingUser != nil {
		return domain.ErrEmailAlreadyExists
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return domain.ErrInternal // Failed to hash password
	}
	user.PasswordHash = string(hashedPassword) // Store hashed password

	// ID will be generated by the repository
	user.Role = domain.RoleUser // Default role for new users
	return s.userRepo.CreateUser(ctx, user)
}

// GetUserByID retrieves a user by their ID.
func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

// GetUserByEmail retrieves a user by their email.
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

// UpdateUser updates an existing user.
func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	// Add business rules specific to updating, e.g., only self or admin can update
	existingUser, err := s.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		return err
	}
	// Do not update password hash directly here unless new password is provided and hashed
	user.PasswordHash = existingUser.PasswordHash // Preserve existing hash if not changed in update request
	user.UpdatedAt = time.Now()
	return s.userRepo.UpdateUser(ctx, user)
}

// DeleteUser deletes a user by their ID.
func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	return s.userRepo.DeleteUser(ctx, userID)
}
