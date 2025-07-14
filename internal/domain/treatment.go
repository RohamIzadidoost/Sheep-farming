package domain

import (
	"time"

	"gorm.io/gorm" // Import GORM
)

// Treatment represents a treatment record for a sheep.
type Treatment struct {
	gorm.Model
	// No explicit ID field added here as gorm.Model's uint ID is likely sufficient for child records.

	Date               time.Time `json:"date" firestore:"date" gorm:"column:date"`
	DiseaseDescription string    `json:"diseaseDescription" firestore:"diseaseDescription" gorm:"column:disease_description"`
	TreatDescription   string    `json:"treatDescription" firestore:"treatDescription" gorm:"column:treat_description"`

	// Foreign key to Sheep (added to enable the relationship from Sheep struct)
	// You'll need to make sure this field exists in your actual Treatment struct
	// if it's in a separate file, or GORM won't be able to establish the link.
	SheepID uint `gorm:"index;column:sheep_id"`
}
