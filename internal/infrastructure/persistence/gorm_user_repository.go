package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

var _ ports.UserRepository = &GormUserRepository{}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// CreateUser implements ports.UserRepository
func (r *GormUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	var existing domain.User
	err := r.db.WithContext(ctx).Where("email = ?", user.Email).First(&existing).Error
	if err == nil {
		return domain.ErrEmailAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check for existing user: %w", err)
	}

	return r.db.WithContext(ctx).Create(user).Error
}

// GetUserByID implements ports.UserRepository
func (r *GormUserRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &user, err
}

// GetUserByEmail implements ports.UserRepository
func (r *GormUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &user, err
}

// UpdateUser implements ports.UserRepository
func (r *GormUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(user).Error
}

// DeleteUser implements ports.UserRepository
func (r *GormUserRepository) DeleteUser(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Where("id = ?", userID).Delete(&domain.User{}).Error
}
