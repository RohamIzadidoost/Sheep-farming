package domain

import "time"

// Treatment represents a treatment record for a sheep.
type Treatment struct {
	Date               time.Time `json:"date" firestore:"date"`
	DiseaseDescription string    `json:"diseaseDescription" firestore:"diseaseDescription"`
	TreatDescription   string    `json:"treatDescription" firestore:"treatDescription"`
}
