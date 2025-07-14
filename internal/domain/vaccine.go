package domain

import (
	"time"

	"gorm.io/gorm" // Import GORM
)

// Vaccination represents a vaccine record for a sheep.
type Vaccination struct {
	gorm.Model // Embed gorm.Model for ID (uint), CreatedAt, UpdatedAt, DeletedAt
	// No explicit ID field added here as gorm.Model's uint ID is likely sufficient for child records.

	Date        time.Time `json:"date" firestore:"date"`
	Vaccine     string    `json:"vaccine" firestore:"vaccine"`
	Vaccinator  string    `json:"vaccinator" firestore:"vaccinator"`
	Description string    `json:"description,omitempty" firestore:"description,omitempty"`

	// Foreign key to Sheep (added to enable the relationship from Sheep struct)
	// You'll need to make sure this field exists in your actual Vaccination struct
	// if it's in a separate file, or GORM won't be able to establish the link.
	SheepID uint `gorm:"index"` // Assuming Sheep's ID is uint from gorm.Model
}
// Vaccine represents a type of vaccine defined by the user.
type Vaccine struct {
	gorm.Model // Embed gorm.Model for ID (uint), CreatedAt, UpdatedAt, DeletedAt
	// Your existing ID field is a string, which means gorm.Model's ID (uint) will also exist.
	// This might lead to two ID columns (one uint, one string).
	// To use your string ID as primary, you'd typically remove gorm.Model and define `ID string `gorm:"primaryKey"`.
	// Or, if gorm.Model's ID is enough, remove your string ID field.
	// Per your request, no fields are changed or removed.

	ID           string    `json:"id,omitempty" firestore:"id,omitempty"` // Your original ID field
	Name         string    `json:"name" firestore:"name"`
	IntervalMonths int       `json:"intervalMonths" firestore:"intervalMonths"` // How often it should be administered (in months)
	OwnerUserID  string    `gorm:"index" json:"ownerUserId" firestore:"ownerUserId"` // Added GORM index for OwnerUserID
	CreatedAt    time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" firestore:"updatedAt"`
}
