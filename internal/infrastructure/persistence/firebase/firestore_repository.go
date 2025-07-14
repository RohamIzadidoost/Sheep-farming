package firebase

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sheep_farm_backend_go/internal/application/ports" // Implement ports
	"sheep_farm_backend_go/internal/domain"            // Import domain entities
)

// Ensure FirestoreRepository implements the ports interfaces
var _ ports.SheepRepository = &FirestoreRepository{}
var _ ports.VaccineRepository = &FirestoreRepository{}
var _ ports.VaccinationRepository = &FirestoreRepository{}
var _ ports.TreatmentRepository = &FirestoreRepository{}
var _ ports.LambingRepository = &FirestoreRepository{}

// FirestoreRepository implements the SheepRepository and VaccineRepository interfaces using Firestore.
type FirestoreRepository struct {
	client *firestore.Client
	appID  string // The __app_id from the frontend context
}

// NewFirestoreRepository creates a new FirestoreRepository.
func NewFirestoreRepository(client *firestore.Client, appID string) *FirestoreRepository {
	return &FirestoreRepository{client: client, appID: appID}
}

// Helper to get the correct collection path for a user
func (r *FirestoreRepository) getUserCollection(userID, collectionName string) *firestore.CollectionRef {
	// This path must match your Firestore security rules: /artifacts/{appId}/users/{userId}/{collectionName}
	return r.client.Collection(fmt.Sprintf("artifacts/%s/users/%s/%s", r.appID, userID, collectionName))
}

// CreateSheep implements ports.SheepRepository
func (r *FirestoreRepository) CreateSheep(ctx context.Context, sheep *domain.Sheep) error {
	if sheep.OwnerUserID == "" {
		return domain.ErrUnauthorized // OwnerUserID is required
	}
	// Firestore generates ID automatically if not provided
	docRef, _, err := r.getUserCollection(sheep.OwnerUserID, "sheep").Add(ctx, sheep)
	if err != nil {
		return fmt.Errorf("failed to create sheep in Firestore: %w", err)
	}
	sheep.ID = docRef.ID // Update sheep object with generated ID
	return nil
}

// GetSheepByID implements ports.SheepRepository
func (r *FirestoreRepository) GetSheepByID(ctx context.Context, userID, sheepID string) (*domain.Sheep, error) {
	docSnap, err := r.getUserCollection(userID, "sheep").Doc(sheepID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get sheep from Firestore: %w", err)
	}

	var sheep domain.Sheep
	if err := docSnap.DataTo(&sheep); err != nil {
		return nil, fmt.Errorf("failed to convert Firestore data to sheep: %w", err)
	}
	sheep.ID = docSnap.Ref.ID // Ensure ID is populated

	lambings, err := r.GetLambings(ctx, userID, sheepID)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	treatments, err := r.GetTreatments(ctx, userID, sheepID)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	vaccinations, err := r.GetVaccinations(ctx, userID, sheepID)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	sheep.Lambings = lambings
	sheep.Treatments = treatments
	sheep.Vaccinations = vaccinations
	return &sheep, nil
}

// GetAllSheep implements ports.SheepRepository
func (r *FirestoreRepository) GetAllSheep(ctx context.Context, userID string) ([]domain.Sheep, error) {
	var sheepList []domain.Sheep
	iter := r.getUserCollection(userID, "sheep").Documents(ctx)
	for {
		docSnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate sheep documents: %w", err)
		}

		var sheep domain.Sheep
		if err := docSnap.DataTo(&sheep); err != nil {
			return nil, fmt.Errorf("failed to convert Firestore data to sheep: %w", err)
		}
		sheep.ID = docSnap.Ref.ID

		lambings, err := r.GetLambings(ctx, userID, sheep.ID)
		if err != nil && err != domain.ErrNotFound {
			return nil, err
		}
		treatments, err := r.GetTreatments(ctx, userID, sheep.ID)
		if err != nil && err != domain.ErrNotFound {
			return nil, err
		}
		vaccinations, err := r.GetVaccinations(ctx, userID, sheep.ID)
		if err != nil && err != domain.ErrNotFound {
			return nil, err
		}
		sheep.Lambings = lambings
		sheep.Treatments = treatments
		sheep.Vaccinations = vaccinations

		sheepList = append(sheepList, sheep)
	}
	return sheepList, nil
}

// UpdateSheep implements ports.SheepRepository
func (r *FirestoreRepository) UpdateSheep(ctx context.Context, sheep *domain.Sheep) error {
	// Use Map to allow partial updates with firestore.Set(ctx, map, firestore.MergeAll)
	// Or define custom struct for updates. For simplicity, we'll update all fields from the sheep struct.
	updateMap := map[string]interface{}{
		"earNumber1":        sheep.EarNumber1,
		"earNumber2":        sheep.EarNumber2,
		"earNumber3":        sheep.EarNumber3,
		"neckNumber":        sheep.NeckNumber,
		"fatherGen":         sheep.FatherGen,
		"birthWeight":       sheep.BirthWeight,
		"gender":            sheep.Gender,
		"reproductionState": sheep.ReproductionState,
		"healthState":       sheep.HealthState,
		"dateOfBirth":       sheep.DateOfBirth,
		"lastShearingDate":  sheep.LastShearingDate,
		"lastHoofTrimDate":  sheep.LastHoofTrimDate,
		"photoUrl":          sheep.PhotoURL,
		"updatedAt":         time.Now(),
	}

	_, err := r.getUserCollection(sheep.OwnerUserID, "sheep").Doc(sheep.ID).Set(ctx, updateMap, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to update sheep in Firestore: %w", err)
	}
	return nil
}

// DeleteSheep implements ports.SheepRepository
func (r *FirestoreRepository) DeleteSheep(ctx context.Context, userID, sheepID string) error {
	_, err := r.getUserCollection(userID, "sheep").Doc(sheepID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete sheep from Firestore: %w", err)
	}
	return nil
}

// CreateVaccine implements ports.VaccineRepository
func (r *FirestoreRepository) CreateVaccine(ctx context.Context, vaccine *domain.Vaccine) error {
	if vaccine.OwnerUserID == "" {
		return domain.ErrUnauthorized // OwnerUserID is required
	}
	docRef, _, err := r.getUserCollection(vaccine.OwnerUserID, "vaccines").Add(ctx, vaccine)
	if err != nil {
		return fmt.Errorf("failed to create vaccine in Firestore: %w", err)
	}
	vaccine.ID = docRef.ID
	return nil
}

// GetVaccineByID implements ports.VaccineRepository
func (r *FirestoreRepository) GetVaccineByID(ctx context.Context, userID, vaccineID string) (*domain.Vaccine, error) {
	docSnap, err := r.getUserCollection(userID, "vaccines").Doc(vaccineID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get vaccine from Firestore: %w", err)
	}

	var vaccine domain.Vaccine
	if err := docSnap.DataTo(&vaccine); err != nil {
		return nil, fmt.Errorf("failed to convert Firestore data to vaccine: %w", err)
	}
	vaccine.ID = docSnap.Ref.ID
	return &vaccine, nil
}

// GetAllVaccines implements ports.VaccineRepository
func (r *FirestoreRepository) GetAllVaccines(ctx context.Context, userID string) ([]domain.Vaccine, error) {
	var vaccinesList []domain.Vaccine
	iter := r.getUserCollection(userID, "vaccines").Documents(ctx)
	for {
		docSnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate vaccine documents: %w", err)
		}

		var vaccine domain.Vaccine
		if err := docSnap.DataTo(&vaccine); err != nil {
			return nil, fmt.Errorf("failed to convert Firestore data to vaccine: %w", err)
		}
		vaccine.ID = docSnap.Ref.ID
		vaccinesList = append(vaccinesList, vaccine)
	}
	return vaccinesList, nil
}

// UpdateVaccine implements ports.VaccineRepository
func (r *FirestoreRepository) UpdateVaccine(ctx context.Context, vaccine *domain.Vaccine) error {
	updateMap := map[string]interface{}{
		"name":           vaccine.Name,
		"intervalMonths": vaccine.IntervalMonths,
		"updatedAt":      time.Now(),
	}
	_, err := r.getUserCollection(vaccine.OwnerUserID, "vaccines").Doc(vaccine.ID).Set(ctx, updateMap, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to update vaccine in Firestore: %w", err)
	}
	return nil
}

// DeleteVaccine implements ports.VaccineRepository
func (r *FirestoreRepository) DeleteVaccine(ctx context.Context, userID, vaccineID string) error {
	_, err := r.getUserCollection(userID, "vaccines").Doc(vaccineID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete vaccine from Firestore: %w", err)
	}
	return nil
}

// CreateVaccination implements ports.VaccinationRepository
func (r *FirestoreRepository) CreateVaccination(ctx context.Context, userID, sheepID string, v domain.Vaccination) error {
	_, _, err := r.getUserCollection(userID, "sheep").Doc(sheepID).Collection("vaccinations").Add(ctx, v)
	if err != nil {
		return fmt.Errorf("failed to create vaccination in Firestore: %w", err)
	}
	return nil
}

// GetVaccinations implements ports.VaccinationRepository
func (r *FirestoreRepository) GetVaccinations(ctx context.Context, userID, sheepID string) ([]domain.Vaccination, error) {
	iter := r.getUserCollection(userID, "sheep").Doc(sheepID).Collection("vaccinations").Documents(ctx)
	var list []domain.Vaccination
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate vaccinations: %w", err)
		}
		var v domain.Vaccination
		if err := doc.DataTo(&v); err != nil {
			return nil, fmt.Errorf("failed to decode vaccination: %w", err)
		}
		list = append(list, v)
	}
	return list, nil
}

// DeleteVaccination implements ports.VaccinationRepository. The index corresponds to order returned by GetVaccinations.
func (r *FirestoreRepository) DeleteVaccination(ctx context.Context, userID, sheepID string, index int) error {
	// Retrieve documents to find the one at given index
	iter := r.getUserCollection(userID, "sheep").Doc(sheepID).Collection("vaccinations").Documents(ctx)
	i := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate vaccinations: %w", err)
		}
		if i == index {
			_, err := doc.Ref.Delete(ctx)
			return err
		}
		i++
	}
	return fmt.Errorf("vaccination index not found")
}

// AddTreatment implements ports.TreatmentRepository
func (r *FirestoreRepository) AddTreatment(ctx context.Context, userID, sheepID string, t domain.Treatment) error {
	doc := r.getUserCollection(userID, "sheep").Doc(sheepID)
	if _, err := doc.Get(ctx); err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to get sheep for treatment: %w", err)
	}
	_, _, err := doc.Collection("treatments").Add(ctx, t)
	if err != nil {
		return fmt.Errorf("failed to add treatment: %w", err)
	}
	return nil
}

// GetTreatments implements ports.TreatmentRepository
func (r *FirestoreRepository) GetTreatments(ctx context.Context, userID, sheepID string) ([]domain.Treatment, error) {
	coll := r.getUserCollection(userID, "sheep").Doc(sheepID).Collection("treatments")
	iter := coll.Documents(ctx)
	var list []domain.Treatment
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return nil, domain.ErrNotFound
			}
			return nil, fmt.Errorf("failed to iterate treatments: %w", err)
		}
		var t domain.Treatment
		if err := doc.DataTo(&t); err != nil {
			return nil, fmt.Errorf("failed to decode treatment: %w", err)
		}
		list = append(list, t)
	}
	return list, nil
}

// UpdateTreatment implements ports.TreatmentRepository
func (r *FirestoreRepository) UpdateTreatment(ctx context.Context, userID, sheepID string, index int, t domain.Treatment) error {
	coll := r.getUserCollection(userID, "sheep").Doc(sheepID).Collection("treatments")
	iter := coll.Documents(ctx)
	i := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate treatments: %w", err)
		}
		if i == index {
			_, err := doc.Ref.Set(ctx, t)
			return err
		}
		i++
	}
	return domain.ErrNotFound
}

// DeleteTreatment implements ports.TreatmentRepository
func (r *FirestoreRepository) DeleteTreatment(ctx context.Context, userID, sheepID string, index int) error {
	coll := r.getUserCollection(userID, "sheep").Doc(sheepID).Collection("treatments")
	iter := coll.Documents(ctx)
	i := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate treatments: %w", err)
		}
		if i == index {
			_, err := doc.Ref.Delete(ctx)
			return err
		}
		i++
	}
	return domain.ErrNotFound
}

// AddLambing implements ports.LambingRepository
func (r *FirestoreRepository) AddLambing(ctx context.Context, userID, sheepID string, l domain.Lambing) error {
	doc := r.getUserCollection(userID, "sheep").Doc(sheepID)
	if _, err := doc.Get(ctx); err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to get sheep for lambing: %w", err)
	}
	_, _, err := doc.Collection("lambings").Add(ctx, l)
	if err != nil {
		return fmt.Errorf("failed to add lambing: %w", err)
	}
	return nil
}

// GetLambings implements ports.LambingRepository
func (r *FirestoreRepository) GetLambings(ctx context.Context, userID, sheepID string) ([]domain.Lambing, error) {
	coll := r.getUserCollection(userID, "sheep").Doc(sheepID).Collection("lambings")
	iter := coll.Documents(ctx)
	var list []domain.Lambing
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return nil, domain.ErrNotFound
			}
			return nil, fmt.Errorf("failed to iterate lambings: %w", err)
		}
		var l domain.Lambing
		if err := doc.DataTo(&l); err != nil {
			return nil, fmt.Errorf("failed to decode lambing: %w", err)
		}
		list = append(list, l)
	}
	return list, nil
}

// UpdateLambing implements ports.LambingRepository
func (r *FirestoreRepository) UpdateLambing(ctx context.Context, userID, sheepID string, index int, l domain.Lambing) error {
	coll := r.getUserCollection(userID, "sheep").Doc(sheepID).Collection("lambings")
	iter := coll.Documents(ctx)
	i := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate lambings: %w", err)
		}
		if i == index {
			_, err := doc.Ref.Set(ctx, l)
			return err
		}
		i++
	}
	return domain.ErrNotFound
}

// DeleteLambing implements ports.LambingRepository
func (r *FirestoreRepository) DeleteLambing(ctx context.Context, userID, sheepID string, index int) error {
	coll := r.getUserCollection(userID, "sheep").Doc(sheepID).Collection("lambings")
	iter := coll.Documents(ctx)
	i := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate lambings: %w", err)
		}
		if i == index {
			_, err := doc.Ref.Delete(ctx)
			return err
		}
		i++
	}
	return domain.ErrNotFound
}
