package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors" // CORS middleware

	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/infrastructure/http/handlers"
)

// Server represents the HTTP server.
type Server struct {
	Router         *mux.Router
	SheepHandler   *handlers.SheepHandler
	VaccineHandler *handlers.VaccineHandler
}

// NewServer creates a new HTTP server instance.
func NewServer(sheepService *services.SheepService, vaccineService *services.VaccineService, fixedUserID string) *Server {
	sheepHandler := handlers.NewSheepHandler(sheepService, fixedUserID)
	vaccineHandler := handlers.NewVaccineHandler(vaccineService, fixedUserID)

	router := mux.NewRouter()
	s := &Server{
		Router:         router,
		SheepHandler:   sheepHandler,
		VaccineHandler: vaccineHandler,
	}
	s.setupRoutes()
	return s
}

// setupRoutes defines all HTTP API endpoints.
func (s *Server) setupRoutes() {
	apiRouter := s.Router.PathPrefix("/api/v1").Subrouter()

	// Sheep Routes
	apiRouter.HandleFunc("/sheep", s.SheepHandler.CreateSheep).Methods("POST")
	apiRouter.HandleFunc("/sheep", s.SheepHandler.GetAllSheep).Methods("GET")
	apiRouter.HandleFunc("/sheep/{id}", s.SheepHandler.GetSheepByID).Methods("GET")
	apiRouter.HandleFunc("/sheep/{id}", s.SheepHandler.UpdateSheep).Methods("PUT")
	apiRouter.HandleFunc("/sheep/{id}", s.SheepHandler.DeleteSheep).Methods("DELETE")

	// Vaccine Routes
	apiRouter.HandleFunc("/vaccines", s.VaccineHandler.CreateVaccine).Methods("POST")
	apiRouter.HandleFunc("/vaccines", s.VaccineHandler.GetAllVaccines).Methods("GET")
	apiRouter.HandleFunc("/vaccines/{id}", s.VaccineHandler.GetVaccineByID).Methods("GET")
	apiRouter.HandleFunc("/vaccines/{id}", s.VaccineHandler.UpdateVaccine).Methods("PUT")
	apiRouter.HandleFunc("/vaccines/{id}", s.VaccineHandler.DeleteVaccine).Methods("DELETE")

	// Root endpoint for health check
	s.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Sheep Farm API is running!"))
	})
}

// Start runs the HTTP server.
func (s *Server) Start(addr string) {
	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your React app
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
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
