package services

import (
	"context"
	"fmt"
	"time"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

// ReminderService calculates and sends reminders.
type ReminderService struct {
	sheepRepo   ports.SheepRepository
	vaccineRepo ports.VaccineRepository
	notifier    ports.ReminderNotifier
}

// NewReminderService creates a new ReminderService instance.
func NewReminderService(sheepRepo ports.SheepRepository, vaccineRepo ports.VaccineRepository, notifier ports.ReminderNotifier) *ReminderService {
	return &ReminderService{
		sheepRepo:   sheepRepo,
		vaccineRepo: vaccineRepo,
		notifier:    notifier,
	}
}

// CalculateAndSendReminders fetches all sheep and vaccine data for a user,
// calculates upcoming reminders, and sends notifications.
// This function would typically be called by a scheduler.
func (s *ReminderService) CalculateAndSendReminders(ctx context.Context, userID string) ([]domain.Reminder, error) {
	sheepList, err := s.sheepRepo.GetAllSheep(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all sheep for reminders: %w", err)
	}

	vaccineDefs, err := s.vaccineRepo.GetAllVaccines(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all vaccine definitions for reminders: %w", err)
	}
	vaccineMap := make(map[string]domain.Vaccine)
	for _, v := range vaccineDefs {
		vaccineMap[v.ID] = v
	}

	var upcomingReminders []domain.Reminder
	now := time.Now()

	for _, sheep := range sheepList {

		// Shearing reminder (e.g., annually from last shearing)
		if sheep.LastShearingDate != nil {
			nextShearingDate := sheep.LastShearingDate.AddDate(1, 0, 0)
			if nextShearingDate.After(now) && nextShearingDate.Before(now.AddDate(0, 1, 0)) { // Remind within next month
				upcomingReminders = append(upcomingReminders, domain.Reminder{
					Type:        domain.ReminderTypeShearing,
					SheepID:     sheep.ID,
					SheepName:   sheep.EarNumber1,
					DueDate:     nextShearingDate,
					Message:     fmt.Sprintf("زمان پشم‌چینی گوسفند %s در تاریخ %s نزدیک است.", sheep.EarNumber1, toPersianDate(nextShearingDate)),
					OwnerUserID: userID,
				})
			}
		}

		// Hoof Trim reminder (e.g., every 6 months)
		if sheep.LastHoofTrimDate != nil {
			nextHoofTrimDate := sheep.LastHoofTrimDate.AddDate(0, 6, 0)
			if nextHoofTrimDate.After(now) && nextHoofTrimDate.Before(now.AddDate(0, 1, 0)) { // Remind within next month
				upcomingReminders = append(upcomingReminders, domain.Reminder{
					Type:        domain.ReminderTypeHoofTrim,
					SheepID:     sheep.ID,
					SheepName:   sheep.EarNumber1,
					DueDate:     nextHoofTrimDate,
					Message:     fmt.Sprintf("زمان سم‌چینی گوسفند %s در تاریخ %s نزدیک است.", sheep.EarNumber1, toPersianDate(nextHoofTrimDate)),
					OwnerUserID: userID,
				})
			}
		}

		// Vaccination reminders based on intervalMonths
		for _, vax := range sheep.Vaccinations {
			if vaxDef, ok := vaccineMap[vax.Vaccine]; ok {
				nextVaccinationDate := vax.Date.AddDate(0, vaxDef.IntervalMonths, 0)
				if nextVaccinationDate.After(now) && nextVaccinationDate.Before(now.AddDate(0, 1, 0)) { // Remind within next month
					upcomingReminders = append(upcomingReminders, domain.Reminder{
						Type:        domain.ReminderTypeVaccination,
						SheepID:     sheep.ID,
						SheepName:   sheep.EarNumber1,
						VaccineName: vaxDef.Name,
						DueDate:     nextVaccinationDate,
						Message:     fmt.Sprintf("زمان واکسن %s برای گوسفند %s در تاریخ %s نزدیک است.", vaxDef.Name, sheep.EarNumber1, toPersianDate(nextVaccinationDate)),
						OwnerUserID: userID,
					})
				}
			}
		}
	}

	// Send notifications for all collected reminders
	for _, reminder := range upcomingReminders {
		err := s.notifier.SendReminder(ctx, reminder)
		if err != nil {
			fmt.Printf("Error sending reminder for %s (sheep %s): %v\n", reminder.Type, reminder.SheepName, err)
			// Continue processing other reminders even if one fails
		}
	}

	return upcomingReminders, nil
}

// Converts a Gregorian Date object to a Persian date string (YYYY/MM/DD)
// This is a utility function, ideally placed in a shared utils package.
func toPersianDate(t time.Time) string {
	// This is a placeholder. For accurate Persian date conversion in Go,
	// you would typically use a dedicated library like 'github.com/arsham/persian'
	// or similar, as Go's standard library does not support it natively.
	// Example using pseudo-code or external library:
	// pc := persian.NewCalendar(t)
	// return fmt.Sprintf("%d/%02d/%02d", pc.Year(), pc.Month(), pc.Day())
	return t.Format("2006/01/02") // Fallback to Gregorian for now
}
