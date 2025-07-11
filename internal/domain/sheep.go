package domain

import "time"

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

// Sheep represents a sheep entity in the domain.
type Sheep struct {
	ID                string        `json:"id,omitempty" firestore:"id,omitempty"`
	EarNumber1        string        `json:"earNumber1" firestore:"earNumber1"`
	EarNumber2        string        `json:"earNumber2,omitempty" firestore:"earNumber2,omitempty"`
	EarNumber3        string        `json:"earNumber3,omitempty" firestore:"earNumber3,omitempty"`
	NeckNumber        *string       `json:"neckNumber,omitempty" firestore:"neckNumber,omitempty"`
	FatherGen         string        `json:"fatherGen,omitempty" firestore:"fatherGen,omitempty"`
	BirthWeight       float64       `json:"birthWeight,omitempty" firestore:"birthWeight,omitempty"`
	Gender            string        `json:"gender" firestore:"gender"`
	ReproductionState string        `json:"reproductionState" firestore:"reproductionState"`
	HealthState       string        `json:"healthState" firestore:"healthState"`
	DateOfBirth       time.Time     `json:"dateOfBirth" firestore:"dateOfBirth"`
	LastShearingDate  *time.Time    `json:"lastShearingDate,omitempty" firestore:"lastShearingDate,omitempty"`
	LastHoofTrimDate  *time.Time    `json:"lastHoofTrimDate,omitempty" firestore:"lastHoofTrimDate,omitempty"`
	PhotoURL          string        `json:"photoUrl,omitempty" firestore:"photoUrl,omitempty"`
	Lambings          []Lambing     `json:"lambings" firestore:"lambings"`
	Vaccinations      []Vaccination `json:"vaccinations" firestore:"vaccinations"`
	Treatments        []Treatment   `json:"treatments" firestore:"treatments"`
	OwnerUserID       string        `json:"ownerUserId" firestore:"ownerUserId"`
	CreatedAt         time.Time     `json:"createdAt" firestore:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt" firestore:"updatedAt"`
}
