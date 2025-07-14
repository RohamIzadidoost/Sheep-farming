package domain

import (
	"time"

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

	// Your original string ID field.
	// Since gorm.Model also provides an 'ID' field (of type uint),
	// this 'ID' field might be treated as a regular column if not explicitly made the primary key.
	// To make this 'ID' field the primary key and override gorm.Model's ID, you'd add `gorm:"primaryKey;type:varchar(36)"`.
	// However, per your explicit instruction not to change fields, I'm just adding it here.
	ID string `json:"id,omitempty" firestore:"id,omitempty" gorm:"column:id"`

	Email        string   `gorm:"unique;not null;column:email" json:"email" firestore:"email"`
	PasswordHash string   `gorm:"not null;column:password_hash" json:"-" firestore:"passwordHash"`
	Role         UserRole `gorm:"type:varchar(10);not null;column:role" json:"role" firestore:"role"`

	// Your original CreatedAt and UpdatedAt fields.
	// GORM's gorm.Model also provides these, so you will effectively have two sets of these fields
	// in your Go struct, though GORM will primarily manage the ones from gorm.Model for database operations.
	// The JSON/Firestore tags on these fields will still work for serialization.
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt" gorm:"column:updated_at"`
}
