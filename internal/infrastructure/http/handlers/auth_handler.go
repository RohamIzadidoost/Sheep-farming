package handlers

import (
	"encoding/json"
	"net/http"

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
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == domain.ErrEmailAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Optionally log in the user immediately after registration
	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		// Log the error but don't fail registration just because auto-login failed
		// In production, might return 201 Created without token and require separate login
		http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.AuthResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Token:  token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Login handles POST /login requests.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			http.Error(w, err.Error(), http.StatusUnauthorized) // 401 Unauthorized
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user details to return in response
	user, err := h.userService.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError) // Should not happen if Login succeeded
		return
	}

	resp := dto.AuthResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Token:  token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
