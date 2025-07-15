package domain

import (
	"time"

	"gorm.io/gorm" // Import GORM
)

// Vaccination represents a vaccine record for a sheep.
type Vaccination struct {
	gorm.Model
	// No explicit ID field added here as gorm.Model's uint ID is likely sufficient for child records.

	Date        time.Time `json:"date" firestore:"date" gorm:"column:date"`
	Vaccine     string    `json:"vaccine" firestore:"vaccine" gorm:"column:vaccine"`
	Vaccinator  string    `json:"vaccinator" firestore:"vaccinator" gorm:"column:vaccinator"`
	Description string    `json:"description,omitempty" firestore:"description,omitempty" gorm:"column:description"`

	// Foreign key to Sheep (added to enable the relationship from Sheep struct)
	// You'll need to make sure this field exists in your actual Vaccination struct
	// if it's in a separate file, or GORM won't be able to establish the link.
	SheepID uint `gorm:"index;column:sheep_id"`
}

// Vaccine represents a type of vaccine defined by the user.
type Vaccine struct {
	gorm.Model

	Name           string `json:"name" firestore:"name" gorm:"column:name"`
	IntervalMonths int    `json:"intervalMonths" firestore:"intervalMonths" gorm:"column:interval_months"`
	OwnerUserID    uint   `gorm:"index;column:owner_user_id" json:"ownerUserId" firestore:"ownerUserId"`
}
