package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/domain"
	"sheep_farm_backend_go/internal/infrastructure/config"
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

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// --- 1. Initialize PostgreSQL ---
	db, err := postgres.New(cfg.DatabaseURL)
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
	userService := services.NewUserService(userRepo)                             // NEW: Initialize UserService first
	authService := services.NewAuthService(userRepo, userService, cfg.JWTSecret) // UPDATED: Pass userService to AuthService

	sheepService := services.NewSheepService(sheepRepo, treatmentRepo, lambingRepo)
	vaccineService := services.NewVaccineService(vaccineRepo)
	treatmentService := services.NewTreatmentService(treatmentRepo)
	lambingService := services.NewLambingService(lambingRepo)
	reminderService := services.NewReminderService(sheepRepo, vaccineRepo, reminderNotifier)

	// --- 5. Initialize Scheduler ---
	appScheduler := scheduler.NewScheduler(reminderService, cfg.SchedulerUserID)
	appScheduler.StartScheduler() // Start the scheduler in a goroutine

	// --- 6. Initialize and Start HTTP Server (Presentation Layer) ---
	// User ID for handlers will now come from context after authentication.
	// No need to pass fixedUserID to handlers directly anymore.
	server := http.NewServer(sheepService, vaccineService, lambingService, treatmentService, authService, userService, reminderService)
	serverAddr := fmt.Sprintf(":%s", cfg.APIPort)
	server.Start(serverAddr) // This call is blocking

	// The scheduler and HTTP server run concurrently.
	// Cleanup happens when main exits, due to defer firestoreClient.Close() and appScheduler.StopScheduler()
}
