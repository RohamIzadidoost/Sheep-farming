package dto

import (
	"fmt"
	"time"
)

// DateOnly is a custom type to handle date-only fields for JSON (YYYY-MM-DD).
// This is important because default time.Time marshals to full ISO 8601 with timezone.
type DateOnly time.Time

// MarshalJSON implements the json.Marshaler interface.
func (d DateOnly) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(d).Format("2006-01-02"))), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *DateOnly) UnmarshalJSON(data []byte) error {
	s := string(data)
	// Remove quotes from string
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}

	if s == "" { // Handle empty string for nullable dates
		*d = DateOnly(time.Time{})
		return nil
	}

	// Try parsing YYYY-MM-DD format
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		// As a fallback, try parsing RFC3339 (Firestore default for nullable dates)
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return fmt.Errorf("failed to parse date: %w", err)
		}
	}
	*d = DateOnly(t)
	return nil
}

// ToTime converts DateOnly to time.Time. Returns nil if zero value.
func (d DateOnly) ToTimePtr() *time.Time {
	t := time.Time(d)
	if t.IsZero() {
		return nil
	}
	return &t
}

// FromTimePtr converts *time.Time to DateOnly. Returns zero value if nil.
func FromTimePtr(t *time.Time) DateOnly {
	if t == nil {
		return DateOnly(time.Time{})
	}
	return DateOnly(*t)
}

// FromTime converts time.Time to DateOnly.
func FromTime(t time.Time) DateOnly {
	return DateOnly(t)
}
