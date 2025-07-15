package firebase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

// Ensure FirestoreUserRepository implements the ports.UserRepository interface
var _ ports.UserRepository = &FirestoreUserRepository{}

// FirestoreUserRepository implements the UserRepository interface using Firestore.
type FirestoreUserRepository struct {
	client *firestore.Client
	appID  string // The __app_id from the frontend context
}

// NewFirestoreUserRepository creates a new FirestoreUserRepository.
func NewFirestoreUserRepository(client *firestore.Client, appID string) *FirestoreUserRepository {
	return &FirestoreUserRepository{client: client, appID: appID}
}

// Helper to get the correct collection path for users
func (r *FirestoreUserRepository) getUserCollection() *firestore.CollectionRef {
	// Users collection is usually directly under the appID in Firestore
	// You might want a different structure based on your security rules
	// For example: /artifacts/{appId}/users (contains user profiles)
	return r.client.Collection(fmt.Sprintf("artifacts/%s/users", r.appID))
}

// CreateUser implements ports.UserRepository
func (r *FirestoreUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	// Check if user with this email already exists directly in Firestore (for race conditions)
	iter := r.getUserCollection().Where("email", "==", user.Email).Limit(1).Documents(ctx)
	_, err := iter.Next()
	if err == nil {
		return domain.ErrEmailAlreadyExists // User with this email already exists
	}
	if err != iterator.Done {
		return fmt.Errorf("failed to check for existing user by email: %w", err)
	}

	// Firestore generates ID automatically if not provided
	docRef, _, err := r.getUserCollection().Add(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user in Firestore: %w", err)
	}
	id, _ := strconv.ParseUint(docRef.ID, 10, 64)
	user.ID = uint(id)
	return nil
}

// GetUserByID implements ports.UserRepository
func (r *FirestoreUserRepository) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	docID := strconv.FormatUint(uint64(userID), 10)
	docSnap, err := r.getUserCollection().Doc(docID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user from Firestore: %w", err)
	}

	var user domain.User
	if err := docSnap.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to convert Firestore data to user: %w", err)
	}
	id, _ := strconv.ParseUint(docSnap.Ref.ID, 10, 64)
	user.ID = uint(id)
	return &user, nil
}

// GetUserByEmail implements ports.UserRepository
func (r *FirestoreUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	iter := r.getUserCollection().Where("email", "==", email).Limit(1).Documents(ctx)
	docSnap, err := iter.Next()
	if err == iterator.Done {
		return nil, domain.ErrNotFound // No user with this email
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email from Firestore: %w", err)
	}

	var user domain.User
	if err := docSnap.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to convert Firestore data to user: %w", err)
	}
	uid, _ := strconv.ParseUint(docSnap.Ref.ID, 10, 64)
	user.ID = uint(uid)
	return &user, nil
}

// UpdateUser implements ports.UserRepository
func (r *FirestoreUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	updateMap := map[string]interface{}{
		"email":        user.Email,
		"passwordHash": user.PasswordHash, // This should be the already hashed password
		"role":         user.Role,
		"updatedAt":    time.Now(),
	}
	docID := strconv.FormatUint(uint64(user.ID), 10)
	_, err := r.getUserCollection().Doc(docID).Set(ctx, updateMap, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to update user in Firestore: %w", err)
	}
	return nil
}

// DeleteUser implements ports.UserRepository
func (r *FirestoreUserRepository) DeleteUser(ctx context.Context, userID uint) error {
	docID := strconv.FormatUint(uint64(userID), 10)
	_, err := r.getUserCollection().Doc(docID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user from Firestore: %w", err)
	}
	return nil
}
