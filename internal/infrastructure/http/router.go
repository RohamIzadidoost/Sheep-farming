package http

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/infrastructure/http/handlers"
	"sheep_farm_backend_go/internal/infrastructure/http/middleware"
)

// Server represents the HTTP server.
type Server struct {
	Engine          *gin.Engine
	SheepHandler    *handlers.SheepHandler
	VaccineHandler  *handlers.VaccineHandler
	ReminderHandler *handlers.ReminderHandler
	AuthHandler     *handlers.AuthHandler
	AuthMiddleware  *middleware.AuthMiddleware
}

// NewServer creates a new HTTP server instance.
func NewServer(
	sheepService *services.SheepService,
	vaccineService *services.VaccineService,
	authService ports.AuthService,
	userService *services.UserService,
	reminderService *services.ReminderService,
) *Server {
	authHandler := handlers.NewAuthHandler(authService, userService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	sheepHandler := handlers.NewSheepHandler(sheepService)
	vaccineHandler := handlers.NewVaccineHandler(vaccineService)
	reminderHandler := handlers.NewReminderHandler(reminderService)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	s := &Server{
		Engine:          engine,
		SheepHandler:    sheepHandler,
		VaccineHandler:  vaccineHandler,
		ReminderHandler: reminderHandler,
		AuthHandler:     authHandler,
		AuthMiddleware:  authMiddleware,
	}
	s.setupRoutes()
	return s
}

// setupRoutes defines all HTTP API endpoints.
func (s *Server) setupRoutes() {
	api := s.Engine.Group("/api/v1")

	// Public Auth Routes
	api.POST("/register", s.AuthHandler.Register)
	api.POST("/login", s.AuthHandler.Login)

	// Protected routes
	protected := api.Group("/")
	protected.Use(s.AuthMiddleware.Authenticate())

	protected.POST("/sheep", s.SheepHandler.CreateSheep)
	protected.GET("/sheep", s.SheepHandler.GetAllSheep)
	protected.GET("/sheep/:id", s.SheepHandler.GetSheepByID)
	protected.PUT("/sheep/:id", s.SheepHandler.UpdateSheep)
	protected.DELETE("/sheep/:id", s.SheepHandler.DeleteSheep)

	protected.GET("/reminders", s.ReminderHandler.GetReminders)

	protected.POST("/vaccines", s.VaccineHandler.CreateVaccine)
	protected.GET("/vaccines", s.VaccineHandler.GetAllVaccines)
	protected.GET("/vaccines/:id", s.VaccineHandler.GetVaccineByID)
	protected.PUT("/vaccines/:id", s.VaccineHandler.UpdateVaccine)
	protected.DELETE("/vaccines/:id", s.VaccineHandler.DeleteVaccine)

	s.Engine.GET("/", func(c *gin.Context) {
		c.String(200, "Sheep Farm API is running!")
	})
}

// Start runs the HTTP server.
func (s *Server) Start(addr string) {
	c := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://[::]:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	s.Engine.Use(c)

	log.Printf("Server starting on %s", addr)
	s.Engine.Run(addr)
}
