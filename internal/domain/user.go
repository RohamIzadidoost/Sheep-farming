package domain

// User represents a simplified user in the domain (for linking data ownership).
// In a real app, this would be managed by Firebase Authentication or a dedicated Auth service.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email,omitempty"`
}
