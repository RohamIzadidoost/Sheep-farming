package domain

import (
	"gorm.io/gorm" // Import GORM
)

// UserRole defines the roles a user can have.
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// User represents a user entity in the domain, adapted for GORM.
type User struct {
	gorm.Model

	Email        string   `gorm:"unique;not null;column:email" json:"email" firestore:"email"`
	PasswordHash string   `gorm:"not null;column:password_hash" json:"-" firestore:"passwordHash"`
	Role         UserRole `gorm:"type:varchar(10);not null;column:role" json:"role" firestore:"role"`
}
