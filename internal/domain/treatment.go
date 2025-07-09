package domain

import "time"

// Treatment represents a treatment record for a sheep.
type Treatment struct {
	Date        time.Time `json:"date" firestore:"date"`
	Description string    `json:"description" firestore:"description"` // Details of the treatment}
