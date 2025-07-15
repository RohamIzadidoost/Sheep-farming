package domain

import "time"

// ReminderType defines the type of a reminder.
type ReminderType string

const (
	ReminderTypeLambing     ReminderType = "lambing"
	ReminderTypeShearing    ReminderType = "shearing"
	ReminderTypeHoofTrim    ReminderType = "hoof_trim"
	ReminderTypeVaccination ReminderType = "vaccination"
	ReminderTypeWeaning     ReminderType = "weaning"
)

// Reminder represents an upcoming event or task for a sheep.
type Reminder struct {
	ID          uint         `json:"id,omitempty"`
	Type        ReminderType `json:"type"`
	SheepID     uint         `json:"sheepId"`
	SheepName   string       `json:"sheepName"`
	VaccineName string       `json:"vaccineName,omitempty"` // Only for vaccination reminders
	DueDate     time.Time    `json:"dueDate"`
	Message     string       `json:"message"`
	OwnerUserID uint         `json:"ownerUserId"` // To link reminder to a user
}
