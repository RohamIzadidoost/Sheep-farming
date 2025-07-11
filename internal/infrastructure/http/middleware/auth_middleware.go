package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

// CtxKeyUserID is the key for storing user ID in context.
type CtxKeyUserID string

const (
	UserIDContextKey   CtxKeyUserID = "userID"
	UserRoleContextKey CtxKeyUserID = "userRole" // New: for storing user role in context
)

// AuthMiddleware provides authentication middleware.
type AuthMiddleware struct {
	authService ports.AuthService
}

// NewAuthMiddleware creates a new AuthMiddleware instance.
func NewAuthMiddleware(authService ports.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// Authenticate is the middleware function to validate JWT tokens.
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized.Error()})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidInput.Error()})
			return
		}

		tokenString := parts[1]
		user, err := m.authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			if err == domain.ErrInvalidToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidToken.Error()})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternal.Error()})
			return
		}

		c.Set(string(UserIDContextKey), user.ID)
		c.Set(string(UserRoleContextKey), user.Role)
		c.Next()
	}
}

// GetUserIDFromContext extracts the user ID from the request context.
func GetUserIDFromContext(c *gin.Context) (string, error) {
	val, exists := c.Get(string(UserIDContextKey))
	if !exists {
		return "", domain.ErrUnauthorized
	}
	userID, ok := val.(string)
	if !ok || userID == "" {
		return "", domain.ErrUnauthorized
	}
	return userID, nil
}

// GetUserRoleFromContext extracts the user role from the request context.
func GetUserRoleFromContext(c *gin.Context) (domain.UserRole, error) {
	val, exists := c.Get(string(UserRoleContextKey))
	if !exists {
		return "", domain.ErrUnauthorized
	}
	role, ok := val.(domain.UserRole)
	if !ok || role == "" {
		return "", domain.ErrUnauthorized
	}
	return role, nil
}
