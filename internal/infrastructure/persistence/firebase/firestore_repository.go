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
func (r *FirestoreRepository) getUserCollection(userID uint, collectionName string) *firestore.CollectionRef {
	// This path must match your Firestore security rules: /artifacts/{appId}/users/{userId}/{collectionName}
	uid := strconv.FormatUint(uint64(userID), 10)
	return r.client.Collection(fmt.Sprintf("artifacts/%s/users/%s/%s", r.appID, uid, collectionName))
}

// CreateSheep implements ports.SheepRepository
func (r *FirestoreRepository) CreateSheep(ctx context.Context, sheep *domain.Sheep) error {
	if sheep.OwnerUserID == 0 {
		return domain.ErrUnauthorized // OwnerUserID is required
	}
	docRef, _, err := r.getUserCollection(sheep.OwnerUserID, "sheep").Add(ctx, sheep)
	if err != nil {
		return fmt.Errorf("failed to create sheep in Firestore: %w", err)
	}
	id, _ := strconv.ParseUint(docRef.ID, 10, 64)
	sheep.ID = uint(id)
	return nil
}

// GetSheepByID implements ports.SheepRepository
func (r *FirestoreRepository) GetSheepByID(ctx context.Context, userID, sheepID uint) (*domain.Sheep, error) {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	docSnap, err := r.getUserCollection(userID, "sheep").Doc(docID).Get(ctx)
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
	sid, _ := strconv.ParseUint(docSnap.Ref.ID, 10, 64)
	sheep.ID = uint(sid)

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
func (r *FirestoreRepository) GetAllSheep(ctx context.Context, userID uint) ([]domain.Sheep, error) {
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
		sid, _ := strconv.ParseUint(docSnap.Ref.ID, 10, 64)
		sheep.ID = uint(sid)

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

	docID := strconv.FormatUint(uint64(sheep.ID), 10)
	_, err := r.getUserCollection(sheep.OwnerUserID, "sheep").Doc(docID).Set(ctx, updateMap, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to update sheep in Firestore: %w", err)
	}
	return nil
}

// DeleteSheep implements ports.SheepRepository
func (r *FirestoreRepository) DeleteSheep(ctx context.Context, userID, sheepID uint) error {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	_, err := r.getUserCollection(userID, "sheep").Doc(docID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete sheep from Firestore: %w", err)
	}
	return nil
}

// CreateVaccine implements ports.VaccineRepository
func (r *FirestoreRepository) CreateVaccine(ctx context.Context, vaccine *domain.Vaccine) error {
	if vaccine.OwnerUserID == 0 {
		return domain.ErrUnauthorized // OwnerUserID is required
	}
	docRef, _, err := r.getUserCollection(vaccine.OwnerUserID, "vaccines").Add(ctx, vaccine)
	if err != nil {
		return fmt.Errorf("failed to create vaccine in Firestore: %w", err)
	}
	vid, _ := strconv.ParseUint(docRef.ID, 10, 64)
	vaccine.ID = uint(vid)
	return nil
}

// GetVaccineByID implements ports.VaccineRepository
func (r *FirestoreRepository) GetVaccineByID(ctx context.Context, userID, vaccineID uint) (*domain.Vaccine, error) {
	docID := strconv.FormatUint(uint64(vaccineID), 10)
	docSnap, err := r.getUserCollection(userID, "vaccines").Doc(docID).Get(ctx)
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
	vid, _ := strconv.ParseUint(docSnap.Ref.ID, 10, 64)
	vaccine.ID = uint(vid)
	return &vaccine, nil
}

// GetAllVaccines implements ports.VaccineRepository
func (r *FirestoreRepository) GetAllVaccines(ctx context.Context, userID uint) ([]domain.Vaccine, error) {
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
		vid, _ := strconv.ParseUint(docSnap.Ref.ID, 10, 64)
		vaccine.ID = uint(vid)
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
	docID := strconv.FormatUint(uint64(vaccine.ID), 10)
	_, err := r.getUserCollection(vaccine.OwnerUserID, "vaccines").Doc(docID).Set(ctx, updateMap, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to update vaccine in Firestore: %w", err)
	}
	return nil
}

// DeleteVaccine implements ports.VaccineRepository
func (r *FirestoreRepository) DeleteVaccine(ctx context.Context, userID, vaccineID uint) error {
	docID := strconv.FormatUint(uint64(vaccineID), 10)
	_, err := r.getUserCollection(userID, "vaccines").Doc(docID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete vaccine from Firestore: %w", err)
	}
	return nil
}

// CreateVaccination implements ports.VaccinationRepository
func (r *FirestoreRepository) CreateVaccination(ctx context.Context, userID, sheepID uint, v domain.Vaccination) error {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	_, _, err := r.getUserCollection(userID, "sheep").Doc(docID).Collection("vaccinations").Add(ctx, v)
	if err != nil {
		return fmt.Errorf("failed to create vaccination in Firestore: %w", err)
	}
	return nil
}

// GetVaccinations implements ports.VaccinationRepository
func (r *FirestoreRepository) GetVaccinations(ctx context.Context, userID, sheepID uint) ([]domain.Vaccination, error) {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	iter := r.getUserCollection(userID, "sheep").Doc(docID).Collection("vaccinations").Documents(ctx)
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
func (r *FirestoreRepository) DeleteVaccination(ctx context.Context, userID, sheepID uint, index int) error {
	// Retrieve documents to find the one at given index
	docID := strconv.FormatUint(uint64(sheepID), 10)
	iter := r.getUserCollection(userID, "sheep").Doc(docID).Collection("vaccinations").Documents(ctx)
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
func (r *FirestoreRepository) AddTreatment(ctx context.Context, userID, sheepID uint, t domain.Treatment) error {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	doc := r.getUserCollection(userID, "sheep").Doc(docID)
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
func (r *FirestoreRepository) GetTreatments(ctx context.Context, userID, sheepID uint) ([]domain.Treatment, error) {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	coll := r.getUserCollection(userID, "sheep").Doc(docID).Collection("treatments")
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
func (r *FirestoreRepository) UpdateTreatment(ctx context.Context, userID, sheepID uint, index int, t domain.Treatment) error {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	coll := r.getUserCollection(userID, "sheep").Doc(docID).Collection("treatments")
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
func (r *FirestoreRepository) DeleteTreatment(ctx context.Context, userID, sheepID uint, index int) error {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	coll := r.getUserCollection(userID, "sheep").Doc(docID).Collection("treatments")
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
func (r *FirestoreRepository) AddLambing(ctx context.Context, userID, sheepID uint, l domain.Lambing) error {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	doc := r.getUserCollection(userID, "sheep").Doc(docID)
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
func (r *FirestoreRepository) GetLambings(ctx context.Context, userID, sheepID uint) ([]domain.Lambing, error) {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	coll := r.getUserCollection(userID, "sheep").Doc(docID).Collection("lambings")
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
func (r *FirestoreRepository) UpdateLambing(ctx context.Context, userID, sheepID uint, index int, l domain.Lambing) error {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	coll := r.getUserCollection(userID, "sheep").Doc(docID).Collection("lambings")
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
func (r *FirestoreRepository) DeleteLambing(ctx context.Context, userID, sheepID uint, index int) error {
	docID := strconv.FormatUint(uint64(sheepID), 10)
	coll := r.getUserCollection(userID, "sheep").Doc(docID).Collection("lambings")
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

// FilterTreatments implements ports.TreatmentRepository
func (r *FirestoreRepository) FilterTreatments(ctx context.Context, userID uint, from, to *time.Time) ([]domain.Treatment, error) {
	sheep, err := r.GetAllSheep(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []domain.Treatment
	for _, s := range sheep {
		list, err := r.GetTreatments(ctx, userID, s.ID)
		if err != nil {
			return nil, err
		}
		for _, t := range list {
			if from != nil && t.Date.Before(*from) {
				continue
			}
			if to != nil && t.Date.After(*to) {
				continue
			}
			result = append(result, t)
		}
	}
	return result, nil
}

// FilterLambings implements ports.LambingRepository
func (r *FirestoreRepository) FilterLambings(ctx context.Context, userID uint, from, to *time.Time) ([]domain.Lambing, error) {
	sheep, err := r.GetAllSheep(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []domain.Lambing
	for _, s := range sheep {
		list, err := r.GetLambings(ctx, userID, s.ID)
		if err != nil {
			return nil, err
		}
		for _, l := range list {
			if from != nil && l.Date.Before(*from) {
				continue
			}
			if to != nil && l.Date.After(*to) {
				continue
			}
			result = append(result, l)
		}
	}
	return result, nil
}

// FilterSheep implements ports.SheepRepository
func (r *FirestoreRepository) FilterSheep(ctx context.Context, userID uint, gender *string, minAgeDays, maxAgeDays *int) ([]domain.Sheep, error) {
	sheepList, err := r.GetAllSheep(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []domain.Sheep
	now := time.Now()
	for _, sh := range sheepList {
		if gender != nil && sh.Gender != *gender {
			continue
		}
		ageDays := int(now.Sub(sh.DateOfBirth).Hours() / 24)
		if minAgeDays != nil && ageDays < *minAgeDays {
			continue
		}
		if maxAgeDays != nil && ageDays > *maxAgeDays {
			continue
		}
		result = append(result, sh)
	}
	return result, nil
}
