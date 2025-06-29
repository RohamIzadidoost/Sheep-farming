package domain

import "time"

// Sheep represents a sheep entity in the domain.
type Sheep struct {
	ID               string        `json:"id,omitempty" firestore:"id,omitempty"` // omitempty: don't include if empty
	Name             string        `json:"name" firestore:"name"`
	Gender           string        `json:"gender" firestore:"gender"` // "male" or "female"
	DateOfBirth      time.Time     `json:"dateOfBirth" firestore:"dateOfBirth"`
	BreedingDate     *time.Time    `json:"breedingDate,omitempty" firestore:"breedingDate,omitempty"` // Pointer for nullable field
	LastShearingDate *time.Time    `json:"lastShearingDate,omitempty" firestore:"lastShearingDate,omitempty"`
	LastHoofTrimDate *time.Time    `json:"lastHoofTrimDate,omitempty" firestore:"lastHoofTrimDate,omitempty"`
	PhotoURL         string        `json:"photoUrl,omitempty" firestore:"photoUrl,omitempty"`
	Vaccinations     []Vaccination `json:"vaccinations" firestore:"vaccinations"`
	Treatments       []Treatment   `json:"treatments" firestore:"treatments"`
	OwnerUserID      string        `json:"ownerUserId" firestore:"ownerUserId"` // To link sheep to a user
	CreatedAt        time.Time     `json:"createdAt" firestore:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt" firestore:"updatedAt"`
}
