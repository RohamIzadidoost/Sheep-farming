package domain

import (
	"time"

	"gorm.io/gorm" // Import GORM
)

// Lambing represents a lambing event for a sheep.
type Lambing struct {
	gorm.Model
	// No explicit ID field added here as gorm.Model's uint ID is likely sufficient for child records.

	Date    time.Time `json:"date" firestore:"date" gorm:"column:date"`
	NumBorn int       `json:"numBorn" firestore:"numBorn" gorm:"column:num_born"`
	Sexes   []string  `json:"sexes" firestore:"sexes" gorm:"type:text;column:sexes"`
	NumDead int       `json:"numDead" firestore:"numDead" gorm:"column:num_dead"`

	// Foreign key to Sheep (added to enable the relationship from Sheep struct)
	// You'll need to make sure this field exists in your actual Lambing struct
	// if it's in a separate file, or GORM won't be able to establish the link.
	SheepID uint `gorm:"index;column:sheep_id"`
}
