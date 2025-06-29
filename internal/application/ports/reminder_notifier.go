package ports

import (
	"context"
	"sheep_farm_backend_go/internal/domain"
)

// ReminderNotifier defines the interface for sending reminders (e.g., via SMS, Email).
// This is another "driven port" (output port).
type ReminderNotifier interface {
	SendReminder(ctx context.Context, reminder domain.Reminder) error
}
