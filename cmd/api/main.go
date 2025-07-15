package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"sheep_farm_backend_go/internal/domain"

	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/infrastructure/external"
	"sheep_farm_backend_go/internal/infrastructure/http"
	postgres "sheep_farm_backend_go/internal/infrastructure/persistence/postgres"
	"sheep_farm_backend_go/internal/infrastructure/scheduler"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables directly.")
	}

	// --- 1. Initialize PostgreSQL ---

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}
	db, err := postgres.New(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto migrate tables
	if err := db.AutoMigrate(&domain.User{}, &domain.Sheep{}, &domain.Vaccine{}, &domain.Vaccination{}, &domain.Treatment{}, &domain.Lambing{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// --- 2. Initialize Infrastructure Layer (Repositories & Notifiers) ---
	userRepo := postgres.NewUserRepository(db)
	sheepRepo := postgres.NewSheepRepository(db)
	vaccineRepo := postgres.NewVaccineRepository(db)
	lambingRepo := postgres.NewLambingRepository(db)
	treatmentRepo := postgres.NewTreatmentRepository(db)

	_ = lambingRepo
	_ = treatmentRepo

	// External Notifier
	reminderNotifier := external.NewConsoleNotifier()

	// --- 4. Initialize Application Layer (Services/Use Cases) ---
	userService := services.NewUserService(userRepo)              // NEW: Initialize UserService first
	authService := services.NewAuthService(userRepo, userService) // UPDATED: Pass userService to AuthService

	sheepService := services.NewSheepService(sheepRepo, treatmentRepo, lambingRepo)
	vaccineService := services.NewVaccineService(vaccineRepo)
	treatmentService := services.NewTreatmentService(treatmentRepo)
	lambingService := services.NewLambingService(lambingRepo)
	reminderService := services.NewReminderService(sheepRepo, vaccineRepo, reminderNotifier)

	// --- 5. Initialize Scheduler ---
	// The scheduler needs the User ID to schedule reminders for a specific user.
	// In a full system, scheduler might get user IDs from database or a dedicated service.
	// You might fetch all user IDs and schedule reminders for each.
	fixedUserIDForSchedulerStr := os.Getenv("SCHEDULER_USER_ID")
	if fixedUserIDForSchedulerStr == "" {
		fixedUserIDForSchedulerStr = "1" // default ID for testing
		log.Printf("SCHEDULER_USER_ID not set, using default: %s", fixedUserIDForSchedulerStr)
	}
	fixedID, _ := strconv.ParseUint(fixedUserIDForSchedulerStr, 10, 64)

	appScheduler := scheduler.NewScheduler(reminderService, uint(fixedID))
	appScheduler.StartScheduler() // Start the scheduler in a goroutine

	// --- 6. Initialize and Start HTTP Server (Presentation Layer) ---
	// User ID for handlers will now come from context after authentication.
	// No need to pass fixedUserID to handlers directly anymore.
	server := http.NewServer(sheepService, vaccineService, lambingService, treatmentService, authService, userService, reminderService)
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8080" // Default port for API
	}
	serverAddr := fmt.Sprintf(":%s", apiPort)
	server.Start(serverAddr) // This call is blocking

	// The scheduler and HTTP server run concurrently.
	// Cleanup happens when main exits, due to defer firestoreClient.Close() and appScheduler.StopScheduler()
}
