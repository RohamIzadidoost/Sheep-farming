package middleware

import (
	"context"
	"net/http"
	"strings"

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
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, domain.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
			return
		}

		tokenString := parts[1]
		user, err := m.authService.ValidateToken(r.Context(), tokenString)
		if err != nil {
			if err == domain.ErrInvalidToken {
				http.Error(w, domain.ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError)
			return
		}

		// Store user ID and Role in the request context
		ctx := context.WithValue(r.Context(), UserIDContextKey, user.ID)
		ctx = context.WithValue(ctx, UserRoleContextKey, user.Role) // Store user role
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext extracts the user ID from the request context.
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	if !ok || userID == "" {
		return "", domain.ErrUnauthorized // User ID not found in context
	}
	return userID, nil
}

// GetUserRoleFromContext extracts the user role from the request context.
func GetUserRoleFromContext(ctx context.Context) (domain.UserRole, error) {
	userRole, ok := ctx.Value(UserRoleContextKey).(domain.UserRole)
	if !ok || userRole == "" {
		return "", domain.ErrUnauthorized // User Role not found in context
	}
	return userRole, nil
}
