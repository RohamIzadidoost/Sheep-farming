package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors" // CORS middleware

	"sheep_farm_backend_go/internal/application/ports" // Import ports for auth service
	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/infrastructure/http/handlers"
	"sheep_farm_backend_go/internal/infrastructure/http/middleware" // New middleware
)

// Server represents the HTTP server.
type Server struct {
	Router          *mux.Router
	SheepHandler    *handlers.SheepHandler
	VaccineHandler  *handlers.VaccineHandler
	ReminderHandler *handlers.ReminderHandler
	AuthHandler     *handlers.AuthHandler      // New: Auth handler
	AuthMiddleware  *middleware.AuthMiddleware // New: Auth middleware
}

// NewServer creates a new HTTP server instance.
func NewServer(
	sheepService *services.SheepService,
	vaccineService *services.VaccineService,
	authService ports.AuthService, // New: Auth service
	userService *services.UserService, // New: User service for AuthHandler
	reminderService *services.ReminderService,
) *Server {

	authHandler := handlers.NewAuthHandler(authService, userService) // Pass user service
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// User ID will now be dynamically extracted from JWT in AuthMiddleware
	// The fixedUserID passed to Sheep/Vaccine handlers is for demonstration before auth is fully implemented for every route.
	// In production, SheepHandler and VaccineHandler would get user ID from context via middleware.
	sheepHandler := handlers.NewSheepHandler(sheepService)
	vaccineHandler := handlers.NewVaccineHandler(vaccineService)
	reminderHandler := handlers.NewReminderHandler(reminderService)

	router := mux.NewRouter()
	s := &Server{
		Router:          router,
		SheepHandler:    sheepHandler,
		VaccineHandler:  vaccineHandler,
		ReminderHandler: reminderHandler,
		AuthHandler:     authHandler,    // New
		AuthMiddleware:  authMiddleware, // New
	}
	s.setupRoutes()
	return s
}

// setupRoutes defines all HTTP API endpoints.
func (s *Server) setupRoutes() {
	apiRouter := s.Router.PathPrefix("/api/v1").Subrouter()

	// Public Auth Routes (no authentication required)
	apiRouter.HandleFunc("/register", s.AuthHandler.Register).Methods("POST")
	apiRouter.HandleFunc("/login", s.AuthHandler.Login).Methods("POST")

	// Protected Routes (authentication required)
	protectedRouter := apiRouter.PathPrefix("/").Subrouter()
	protectedRouter.Use(s.AuthMiddleware.Authenticate) // Apply authentication middleware to all routes below

	// Sheep Routes (now protected)
	protectedRouter.HandleFunc("/sheep", s.SheepHandler.CreateSheep).Methods("POST")
	protectedRouter.HandleFunc("/sheep", s.SheepHandler.GetAllSheep).Methods("GET")
	protectedRouter.HandleFunc("/sheep/{id}", s.SheepHandler.GetSheepByID).Methods("GET")
	protectedRouter.HandleFunc("/sheep/{id}", s.SheepHandler.UpdateSheep).Methods("PUT")
	protectedRouter.HandleFunc("/sheep/{id}", s.SheepHandler.DeleteSheep).Methods("DELETE")
	protectedRouter.HandleFunc("/sheep/{id}/vaccinations", s.SheepHandler.AddVaccination).Methods("POST")
	protectedRouter.HandleFunc("/sheep/{id}/treatments", s.SheepHandler.AddTreatment).Methods("POST")
	protectedRouter.HandleFunc("/sheep/{id}/lambings", s.SheepHandler.AddLambing).Methods("POST")

	// Reminder Route
	protectedRouter.HandleFunc("/reminders", s.ReminderHandler.GetReminders).Methods("GET")

	// Vaccine Routes (now protected)
	protectedRouter.HandleFunc("/vaccines", s.VaccineHandler.CreateVaccine).Methods("POST")
	protectedRouter.HandleFunc("/vaccines", s.VaccineHandler.GetAllVaccines).Methods("GET")
	protectedRouter.HandleFunc("/vaccines/{id}", s.VaccineHandler.GetVaccineByID).Methods("GET")
	protectedRouter.HandleFunc("/vaccines/{id}", s.VaccineHandler.UpdateVaccine).Methods("PUT")
	protectedRouter.HandleFunc("/vaccines/{id}", s.VaccineHandler.DeleteVaccine).Methods("DELETE")

	// Root endpoint for health check
	s.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Sheep Farm API is running!"))
	})
}

// Start runs the HTTP server.
func (s *Server) Start(addr string) {
	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://[::]:5500"}, // Allow your React app
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Authorization header is now important!
		AllowCredentials: true,
		Debug:            false, // Set to true for debugging CORS issues
	})

	handler := c.Handler(s.Router)

	log.Printf("Server starting on %s", addr)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
