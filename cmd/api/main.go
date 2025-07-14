package main

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go" // Firebase Admin SDK core
	"github.com/joho/godotenv"        // For loading .env file
	"google.golang.org/api/option"    // For credentials option

	// Firestore client
	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/infrastructure/external"
	"sheep_farm_backend_go/internal/infrastructure/http"
	"sheep_farm_backend_go/internal/infrastructure/persistence"
	"sheep_farm_backend_go/internal/infrastructure/scheduler"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables directly.")
	}

	// --- 1. Initialize Firebase Admin SDK ---
	ctx := context.Background()

	// Get Firebase credentials path from environment variable
	// In production, consider alternative secure ways to provide credentials (e.g., Google Cloud service account environment variable)
	firebaseCredentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if firebaseCredentialsPath == "" {
		log.Fatal("FIREBASE_CREDENTIALS_PATH environment variable not set. Please provide path to your Firebase service account JSON file.")
	}

	sa := option.WithCredentialsFile(firebaseCredentialsPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v\n", err)
	}
	defer firestoreClient.Close() // Close the Firestore client when main exits

	log.Println("Firebase Firestore client initialized.")

	// --- 2. Define App ID ---
	// This appID should match the __app_id used in your React frontend.
	// For local development, you might set it as an environment variable or hardcode if always default.
	// This determines the root collection in Firestore for your application's data.
	appID := os.Getenv("APP_ID")
	if appID == "" {
		appID = "default-app-id" // Fallback for local testing if not set
		log.Printf("APP_ID not set in environment, using default: %s", appID)
	}

	// --- 3. Initialize Infrastructure Layer (Repositories & Notifiers) ---
	// Persistence Repositories
	userRepo := persistence.NewFirestoreUserRepository(firestoreClient, appID) // New: User repository
	sheepRepo := persistence.NewFirestoreRepository(firestoreClient, appID)
	vaccineRepo := persistence.NewFirestoreRepository(firestoreClient, appID) // Same repository for vaccines

	// External Notifier
	reminderNotifier := external.NewConsoleNotifier()

	// --- 4. Initialize Application Layer (Services/Use Cases) ---
	userService := services.NewUserService(userRepo)              // NEW: Initialize UserService first
	authService := services.NewAuthService(userRepo, userService) // UPDATED: Pass userService to AuthService

	sheepService := services.NewSheepService(sheepRepo)
	vaccineService := services.NewVaccineService(vaccineRepo)
	reminderService := services.NewReminderService(sheepRepo, vaccineRepo, reminderNotifier)

	// --- 5. Initialize Scheduler ---
	// The scheduler needs the User ID to schedule reminders for a specific user.
	// In a full system, scheduler might get user IDs from database or a dedicated service.
	// You might fetch all user IDs and schedule reminders for each.
	fixedUserIDForScheduler := os.Getenv("SCHEDULER_USER_ID")
	if fixedUserIDForScheduler == "" {
		fixedUserIDForScheduler = "test-user-id-for-scheduler" // Default for local testing
		log.Printf("SCHEDULER_USER_ID not set, using default: %s", fixedUserIDForScheduler)
	}

	appScheduler := scheduler.NewScheduler(reminderService, fixedUserIDForScheduler)

	runMode := os.Getenv("RUN_MODE")
	if runMode == "server" {
		appScheduler.StartScheduler() // Start scheduler only on server
	}

	// --- 6. Initialize and Start HTTP Server (Presentation Layer) ---
	// User ID for handlers will now come from context after authentication.
	// No need to pass fixedUserID to handlers directly anymore.
	server := http.NewServer(sheepService, vaccineService, authService, userService, reminderService) // Pass reminder service too
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8080" // Default port for API
	}
	serverAddr := fmt.Sprintf(":%s", apiPort)
	server.Start(serverAddr) // This call is blocking

	// The scheduler and HTTP server run concurrently.
	// Cleanup happens when main exits, due to defer firestoreClient.Close() and appScheduler.StopScheduler()
}
