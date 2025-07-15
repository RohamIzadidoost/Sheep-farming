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
	gorm.Model

	EarNumber1        string     `json:"earNumber1" firestore:"earNumber1" gorm:"column:ear_number1"`
	EarNumber2        string     `json:"earNumber2,omitempty" firestore:"earNumber2,omitempty" gorm:"column:ear_number2"`
	EarNumber3        string     `json:"earNumber3,omitempty" firestore:"earNumber3,omitempty" gorm:"column:ear_number3"`
	NeckNumber        *string    `json:"neckNumber,omitempty" firestore:"neckNumber,omitempty" gorm:"column:neck_number"`
	FatherGen         string     `json:"fatherGen,omitempty" firestore:"fatherGen,omitempty" gorm:"column:father_gen"`
	BirthWeight       float64    `json:"birthWeight,omitempty" firestore:"birthWeight,omitempty" gorm:"column:birth_weight"`
	Gender            string     `json:"gender" firestore:"gender" gorm:"column:gender"`
	ReproductionState string     `json:"reproductionState" firestore:"reproductionState" gorm:"column:reproduction_state"`
	HealthState       string     `json:"healthState" firestore:"healthState" gorm:"column:health_state"`
	DateOfBirth       time.Time  `json:"dateOfBirth" firestore:"dateOfBirth" gorm:"column:date_of_birth"`
	LastShearingDate  *time.Time `json:"lastShearingDate,omitempty" firestore:"lastShearingDate,omitempty" gorm:"column:last_shearing_date"`
	LastHoofTrimDate  *time.Time `json:"lastHoofTrimDate,omitempty" firestore:"lastHoofTrimDate,omitempty" gorm:"column:last_hoof_trim_date"`
	PhotoURL          string     `json:"photoUrl,omitempty" firestore:"photoUrl,omitempty" gorm:"column:photo_url"`

	// GORM relationships: You'll need to define the foreign key on the Lambing, Vaccination, and Treatment structs
	Lambings     []Lambing     `gorm:"foreignKey:SheepID"`
	Vaccinations []Vaccination `gorm:"foreignKey:SheepID"`
	Treatments   []Treatment   `gorm:"foreignKey:SheepID"`

	OwnerUserID uint `gorm:"index;column:owner_user_id" json:"ownerUserId" firestore:"ownerUserId"`
}
