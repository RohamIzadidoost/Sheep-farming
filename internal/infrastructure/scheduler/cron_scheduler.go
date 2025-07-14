package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3" // Popular Go scheduler library

	"sheep_farm_backend_go/internal/application/ports" // ReminderService port
)

// Scheduler manages periodic tasks like sending reminders.
type Scheduler struct {
	cron            *cron.Cron
	reminderService ports.ReminderService
	// In a real app, you might have an authentication service to get user IDs
	// Or iterate over all users in your database.
	fixedUserID string // For simplicity, we'll use a fixed user ID for scheduling reminders
}

// NewScheduler creates a new Scheduler instance.
func NewScheduler(reminderService ports.ReminderService, fixedUserID string) *Scheduler {
	return &Scheduler{
		cron:            cron.New(), // cron.New(cron.WithChain(cron.Recover(log.New(os.Stdout, "", log.LstdFlags)))), for robust error handling
		reminderService: reminderService,
		fixedUserID:     fixedUserID,
	}
}

// StartScheduler initializes and starts the cron scheduler.
func (s *Scheduler) StartScheduler() {
	// Schedule the reminder calculation task to run daily at a specific time (e.g., 08:00 AM)
	// Cron string format: "minute hour day_of_month month day_of_week"
	// Example: "0 8 * * *" means "at 08:00 AM every day"
	// For testing, you might use a more frequent interval like "*/5 * * * * *" (every 5 seconds)
	_, err := s.cron.AddFunc("0 8 * * *", func() {
		log.Println("Running daily reminder check...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Give it 30 seconds
		defer cancel()

		// Call the application layer service to calculate and send reminders
		reminders, err := s.reminderService.CalculateAndSendReminders(ctx, s.fixedUserID)
		if err != nil {
			log.Printf("Error during reminder calculation: %v", err)
		} else {
			log.Printf("Finished reminder check. Found %d upcoming reminders.", len(reminders))
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule reminder job: %v", err)
	}

	s.cron.Start()
	log.Println("Scheduler started successfully. Daily reminder check scheduled for 08:00 AM.")
}

// StopScheduler stops the cron scheduler.
func (s *Scheduler) StopScheduler() {
	s.cron.Stop()
	log.Println("Scheduler stopped.")
}
