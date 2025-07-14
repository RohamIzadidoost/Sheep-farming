package ports

import (
	"context"
	"sheep_farm_backend_go/internal/domain"
)

// ReminderService defines reminder related operations.
type ReminderService interface {
	CalculateAndSendReminders(ctx context.Context, userID string) ([]domain.Reminder, error)
}
