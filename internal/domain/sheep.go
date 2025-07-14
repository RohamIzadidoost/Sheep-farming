package domain

import (
	"time"

	"gorm.io/gorm" // Import GORM
)

// Possible reproduction states for a sheep
const (
	ReproductionNormal    = "normal"
	ReproductionPregnant  = "pregnant"
	ReproductionPostBirth = "post_birth"
)

// Possible health states for a sheep
const (
	HealthHealthy        = "healthy"
	HealthSick           = "sick"
	HealthUnderTreatment = "under_treatment"
)

// Sheep represents a sheep entity in the domain, adapted for GORM.
type Sheep struct {
	gorm.Model // Embed gorm.Model for ID (uint), CreatedAt, UpdatedAt, DeletedAt
	// Note: If you specifically need string IDs (UUIDs) for your primary key,
	// you would typically override `ID` like `ID string `gorm:"primaryKey;type:varchar(36)"`.
	// However, per your request, I'm not changing or adding fields, so gorm.Model's uint ID will be used.

	ID                string     `json:"id,omitempty" firestore:"id,omitempty"` // Your original ID field
	EarNumber1        string     `json:"earNumber1" firestore:"earNumber1"`
	EarNumber2        string     `json:"earNumber2,omitempty" firestore:"earNumber2,omitempty"`
	EarNumber3        string     `json:"earNumber3,omitempty" firestore:"earNumber3,omitempty"`
	NeckNumber        *string    `json:"neckNumber,omitempty" firestore:"neckNumber,omitempty"`
	FatherGen         string     `json:"fatherGen,omitempty" firestore:"fatherGen,omitempty"`
	BirthWeight       float64    `json:"birthWeight,omitempty" firestore:"birthWeight,omitempty"`
	Gender            string     `json:"gender" firestore:"gender"`
	ReproductionState string     `json:"reproductionState" firestore:"reproductionState"`
	HealthState       string     `json:"healthState" firestore:"healthState"`
	DateOfBirth       time.Time  `json:"dateOfBirth" firestore:"dateOfBirth"`
	LastShearingDate  *time.Time `json:"lastShearingDate,omitempty" firestore:"lastShearingDate,omitempty"`
	LastHoofTrimDate  *time.Time `json:"lastHoofTrimDate,omitempty" firestore:"lastHoofTrimDate,omitempty"`
	PhotoURL          string     `json:"photoUrl,omitempty" firestore:"photoUrl,omitempty"`

	// GORM relationships: You'll need to define the foreign key on the Lambing, Vaccination, and Treatment structs
	Lambings     []Lambing     `gorm:"foreignKey:SheepID"`     // GORM will expect a 'SheepID' in Lambing struct
	Vaccinations []Vaccination `gorm:"foreignKey:SheepID"` // GORM will expect a 'SheepID' in Vaccination struct
	Treatments   []Treatment   `gorm:"foreignKey:SheepID"`   // GORM will expect a 'SheepID' in Treatment struct

	OwnerUserID string    `gorm:"index" json:"ownerUserId" firestore:"ownerUserId"` // Added GORM index for OwnerUserID
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`                   // GORM will manage this, but keeping your original tag
	UpdatedAt   time.Time `json:"updatedAt" firestore:"updatedAt"`                   // GORM will manage this, but keeping your original tag
}