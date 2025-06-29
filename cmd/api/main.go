package main

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go" // Firebase Admin SDK core
	"github.com/joho/godotenv"        // For loading .env file
	"google.golang.org/api/option"
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

	// --- 2. Define App ID and Fixed User ID ---
	// This appID should match the __app_id used in your React frontend.
	// For local development, you might set it as an environment variable or hardcode if always default.
	// This determines the root collection in Firestore for your application's data.
	appID := os.Getenv("APP_ID")
	if appID == "" {
		appID = "default-app-id" // Fallback for local testing if not set
		log.Printf("APP_ID not set in environment, using default: %s", appID)
	}

	// For simplicity in this guide, we use a fixed user ID.
	// In a real application, this user ID would come from your authentication system
	// (e.g., from a JWT token sent by the frontend after user login).
	fixedUserID := os.Getenv("FIXED_USER_ID")
	if fixedUserID == "" {
		fixedUserID = "test-user-id" // Default for local testing
		log.Printf("FIXED_USER_ID not set in environment, using default: %s", fixedUserID)
	}
	log.Printf("Using fixed user ID for operations: %s", fixedUserID)

	// --- 3. Initialize Infrastructure Layer (Repositories & Notifiers) ---
	// Firestore Repository implementations
	sheepRepo := persistence.NewFirestoreRepository(firestoreClient, appID)
	vaccineRepo := persistence.NewFirestoreRepository(firestoreClient, appID) // Same repository for vaccines

	// External Notifier (e.g., console output, could be Twilio/SendGrid)
	reminderNotifier := external.NewConsoleNotifier()

	// --- 4. Initialize Application Layer (Services/Use Cases) ---
	sheepService := services.NewSheepService(sheepRepo)
	vaccineService := services.NewVaccineService(vaccineRepo)
	reminderService := services.NewReminderService(sheepRepo, vaccineRepo, reminderNotifier)

	// --- 5. Initialize Scheduler ---
	// The scheduler will periodically call reminderService.CalculateAndSendReminders
	appScheduler := scheduler.NewScheduler(reminderService, fixedUserID)
	appScheduler.StartScheduler() // Start the scheduler in a goroutine

	// --- 6. Initialize and Start HTTP Server (Presentation Layer) ---
	server := http.NewServer(sheepService, vaccineService, fixedUserID)
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8080" // Default port for API
	}
	serverAddr := fmt.Sprintf(":%s", apiPort)
	server.Start(serverAddr) // This call is blocking

	// The scheduler and HTTP server run concurrently.
	// Cleanup happens when main exits, due to defer firestoreClient.Close() and appScheduler.StopScheduler()
}

// Ensure you have an .env file in the root of your Go backend project:
// .env example:
// FIREBASE_CREDENTIALS_PATH=/path/to/your/sheep-farm-app-firebase-adminsdk.json
// APP_ID=default-app-id
// FIXED_USER_ID=test-user-id-from-frontend-auth
// API_PORT=8080
