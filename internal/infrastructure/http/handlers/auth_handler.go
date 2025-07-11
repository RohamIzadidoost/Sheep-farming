package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/domain"
	"sheep_farm_backend_go/internal/infrastructure/http/dto"
)

// AuthHandler handles HTTP requests related to authentication.
type AuthHandler struct {
	authService ports.AuthService
	userService *services.UserService // To get user details if needed after token validation
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService ports.AuthService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{authService: authService, userService: userService}
}

// Register handles POST /register requests.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidInput.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == domain.ErrEmailAlreadyExists {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Optionally log in the user immediately after registration
	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// Log the error but don't fail registration just because auto-login failed
		// In production, might return 201 Created without token and require separate login
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternal.Error()})
		return
	}

	resp := dto.AuthResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Token:  token,
	}
	c.JSON(http.StatusCreated, resp)
}

// Login handles POST /login requests.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidInput.Error()})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get user details to return in response
	user, err := h.userService.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternal.Error()})
		return
	}

	resp := dto.AuthResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Token:  token,
	}
	c.JSON(http.StatusOK, resp)
}
