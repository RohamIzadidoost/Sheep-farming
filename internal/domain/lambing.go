package domain

import (
	"time"

	"gorm.io/gorm" // Import GORM
)

// Lambing represents a lambing event for a sheep.
type Lambing struct {
	gorm.Model // Embed gorm.Model for ID (uint), CreatedAt, UpdatedAt, DeletedAt
	// No explicit ID field added here as gorm.Model's uint ID is likely sufficient for child records.

	Date    time.Time `json:"date" firestore:"date"`
	NumBorn int       `json:"numBorn" firestore:"numBorn"`
	Sexes   []string  `gorm:"type:text" json:"sexes" firestore:"sexes"` // GORM stores slices as TEXT (e.g., JSON or CSV)
	NumDead int       `json:"numDead" firestore:"numDead"`

	// Foreign key to Sheep (added to enable the relationship from Sheep struct)
	// You'll need to make sure this field exists in your actual Lambing struct
	// if it's in a separate file, or GORM won't be able to establish the link.
	SheepID uint `gorm:"index"` // Assuming Sheep's ID is uint from gorm.Model
}