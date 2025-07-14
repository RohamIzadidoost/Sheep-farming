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
	// Note: If you specifically need string IDs (UUIDs) for your primary key,
	// you would typically override `ID` like `ID string `gorm:"primaryKey;type:varchar(36)"`.
	// However, per your request, I'm not changing or adding fields, so gorm.Model's uint ID will be used.

	ID                string     `json:"id,omitempty" firestore:"id,omitempty" gorm:"column:id"`
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

	OwnerUserID string    `gorm:"index;column:owner_user_id" json:"ownerUserId" firestore:"ownerUserId"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt" gorm:"column:created_at"` // GORM will manage this, but keeping your original tag
	UpdatedAt   time.Time `json:"updatedAt" firestore:"updatedAt" gorm:"column:updated_at"` // GORM will manage this, but keeping your original tag
}
