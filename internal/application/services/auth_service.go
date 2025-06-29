package services

import (
	"context"
	"errors" // Make sure errors is imported here
	"fmt"
	"os" // For JWT_SECRET_KEY
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt" // For password comparison

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

// claims represents the JWT claims.
type Claims struct {
	UserID string          `json:"userId"`
	Email  string          `json:"email"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// AuthService implements ports.AuthService.
type AuthService struct {
	userRepo     ports.UserRepository // Still need userRepo for Login/ValidateToken operations
	userService  *UserService         // NEW: Dependency on UserService for user creation
	jwtSecretKey []byte
}

// NewAuthService creates a new AuthService instance.
// UPDATED: Now accepts UserService as a dependency.
func NewAuthService(userRepo ports.UserRepository, userService *UserService) *AuthService {
	// Get JWT secret key from environment variable
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		// In production, this should be a strong, random key.
		// For development, it's okay to have a default, but warn the user.
		jwtSecret = "supersecretjwtkeythatshouldbechangedinproduction"
		fmt.Println("WARNING: JWT_SECRET_KEY not set, using default. Change in production!")
	}
	return &AuthService{
		userRepo:     userRepo,
		userService:  userService, // Assign the passed UserService
		jwtSecretKey: []byte(jwtSecret),
	}
}

// Register registers a new user.
// UPDATED: Calls userService.CreateUser to hash password and persist.
func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// Basic validation for email and password length
	if email == "" || password == "" {
		return nil, domain.ErrInvalidInput
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: password,        // PasswordHash will contain the plaintext password here, UserService will hash it
		Role:         domain.RoleUser, // Default role
	}

	// CORRECTED: Call UserService.CreateUser to handle user creation and password hashing
	err := s.userService.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", domain.ErrInvalidInput
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email) // Still uses userRepo for getting user by email
	if err != nil {
		if err == domain.ErrNotFound {
			return "", domain.ErrInvalidCredentials // User not found
		}
		return "", err
	}

	// Compare provided password with hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidCredentials // Password does not match
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecretKey)
	if err != nil {
		return "", domain.ErrInternal // Failed to sign token
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user if valid.
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, domain.ErrInvalidToken // Signature is invalid
		}
		return nil, domain.ErrInvalidToken // Other token parsing errors (e.g., malformed, expired)
	}

	if !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	// Optionally, fetch user from repo to ensure they still exist and token is not revoked
	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrInvalidToken // User associated with token not found
		}
		return nil, err
	}

	return user, nil
}
