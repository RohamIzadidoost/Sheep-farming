package domain

import "time"

// Vaccination represents a vaccine record for a sheep.
type Vaccination struct {
	Date        time.Time `json:"date" firestore:"date"`
	Vaccine     string    `json:"vaccine" firestore:"vaccine"`
	Vaccinator  string    `json:"vaccinator" firestore:"vaccinator"`
	Description string    `json:"description,omitempty" firestore:"description,omitempty"`
}

// Vaccine represents a type of vaccine defined by the user.
type Vaccine struct {
	ID             string    `json:"id,omitempty" firestore:"id,omitempty"`
	Name           string    `json:"name" firestore:"name"`
	IntervalMonths int       `json:"intervalMonths" firestore:"intervalMonths"` // How often it should be administered (in months)
	OwnerUserID    string    `json:"ownerUserId" firestore:"ownerUserId"`
	CreatedAt      time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt" firestore:"updatedAt"`
}
