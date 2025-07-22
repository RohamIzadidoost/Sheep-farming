package external

import (
	"context"
	"fmt"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

// Ensure ConsoleNotifier implements the ReminderNotifier interface
var _ ports.ReminderNotifier = &ConsoleNotifier{}

// ConsoleNotifier sends reminders to the console (for demonstration purposes).
type ConsoleNotifier struct{}

// NewConsoleNotifier creates a new ConsoleNotifier.
func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{}
}

// SendReminder implements ports.ReminderNotifier.
func (n *ConsoleNotifier) SendReminder(ctx context.Context, reminder domain.Reminder) error {
	fmt.Printf("--- REMINDER to User %d ---\n", reminder.OwnerUserID)
	fmt.Printf("Type: %s\n", reminder.Type)
	fmt.Printf("Sheep: %s\n", reminder.SheepName)
	if reminder.VaccineName != "" {
		fmt.Printf("Vaccine: %s\n", reminder.VaccineName)
	}
	fmt.Printf("Due Date: %s\n", reminder.DueDate.Format("2006-01-02"))
	fmt.Printf("Message: %s\n", reminder.Message)
	fmt.Println("---------------------------------")
	// In a real application, you would integrate with SMS/Email/Push Notification services here.
	return nil
}
