package domain

import "time"

// UserRole defines the roles a user can have.
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// User represents a user entity in the domain.
type User struct {
	ID           string    `json:"id,omitempty" firestore:"id,omitempty"`
	Email        string    `json:"email" firestore:"email"`    // User's email, used for login
	PasswordHash string    `json:"-" firestore:"passwordHash"` // Hashed password, never exposed via JSON
	Role         UserRole  `json:"role" firestore:"role"`      // User's role (admin, user)
	CreatedAt    time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" firestore:"updatedAt"`
}
