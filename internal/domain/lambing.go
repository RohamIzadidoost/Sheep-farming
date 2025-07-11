package domain

import "time"

// Lambing represents a lambing event for a sheep.
type Lambing struct {
	Date    time.Time `json:"date" firestore:"date"`
	NumBorn int       `json:"numBorn" firestore:"numBorn"`
	Sexes   []string  `json:"sexes" firestore:"sexes"`
	NumDead int       `json:"numDead" firestore:"numDead"`
}
