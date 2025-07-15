package services

import (
	"context"
	"time"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

// SheepService provides use cases for sheep management.
// This is an "application service" (use case).
type SheepService struct {
	repo      ports.SheepRepository
	treatRepo ports.TreatmentRepository
	lambRepo  ports.LambingRepository
}

// NewSheepService creates a new SheepService instance.
func NewSheepService(repo ports.SheepRepository, treatRepo ports.TreatmentRepository, lambRepo ports.LambingRepository) *SheepService {
	return &SheepService{repo: repo, treatRepo: treatRepo, lambRepo: lambRepo}
}

// CreateSheep handles the creation of a new sheep.
func (s *SheepService) CreateSheep(ctx context.Context, sheep *domain.Sheep) error {
	// Add business rules here before persisting
	return s.repo.CreateSheep(ctx, sheep)
}

// GetSheepByID retrieves a sheep by its ID.
func (s *SheepService) GetSheepByID(ctx context.Context, userID, sheepID uint) (*domain.Sheep, error) {
	return s.repo.GetSheepByID(ctx, userID, sheepID)
}

// GetAllSheep retrieves all sheep for a given user.
func (s *SheepService) GetAllSheep(ctx context.Context, userID uint) ([]domain.Sheep, error) {
	return s.repo.GetAllSheep(ctx, userID)
}

// FilterSheep retrieves sheep filtered by gender and age range (in days).
func (s *SheepService) FilterSheep(ctx context.Context, userID uint, gender *string, minAgeDays, maxAgeDays *int) ([]domain.Sheep, error) {
	return s.repo.FilterSheep(ctx, userID, gender, minAgeDays, maxAgeDays)
}

// UpdateSheep updates an existing sheep.
func (s *SheepService) UpdateSheep(ctx context.Context, sheep *domain.Sheep) error {
	// Add business rules specific to updating
	existingSheep, err := s.repo.GetSheepByID(ctx, sheep.OwnerUserID, sheep.ID)
	if err != nil {
		return err // Or return domain.ErrNotFound
	}
	if existingSheep.OwnerUserID != sheep.OwnerUserID {
		return domain.ErrUnauthorized // Ensure only owner can update
	}
	sheep.UpdatedAt = time.Now()
	return s.repo.UpdateSheep(ctx, sheep)
}

// DeleteSheep deletes a sheep by its ID.
func (s *SheepService) DeleteSheep(ctx context.Context, userID, sheepID uint) error {
	// Add business rules specific to deleting
	return s.repo.DeleteSheep(ctx, userID, sheepID)
}

// AddVaccination appends a vaccination record to the sheep.
func (s *SheepService) AddVaccination(ctx context.Context, userID, sheepID uint, v domain.Vaccination) error {
	sh, err := s.repo.GetSheepByID(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if sh.OwnerUserID != userID {
		return domain.ErrUnauthorized
	}
	sh.Vaccinations = append(sh.Vaccinations, v)
	return s.repo.UpdateSheep(ctx, sh)
}

// AddTreatment appends a treatment record to the sheep.
func (s *SheepService) AddTreatment(ctx context.Context, userID, sheepID uint, t domain.Treatment) error {
	sh, err := s.repo.GetSheepByID(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if sh.OwnerUserID != userID {
		return domain.ErrUnauthorized
	}
	return s.treatRepo.AddTreatment(ctx, userID, sheepID, t)
}

// AddLambing appends a lambing record to the sheep.
func (s *SheepService) AddLambing(ctx context.Context, userID, sheepID uint, l domain.Lambing) error {
	sh, err := s.repo.GetSheepByID(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if sh.OwnerUserID != userID {
		return domain.ErrUnauthorized
	}
	return s.lambRepo.AddLambing(ctx, userID, sheepID, l)
}

// UpdateVaccination updates a vaccination record by index.
func (s *SheepService) UpdateVaccination(ctx context.Context, userID, sheepID uint, index int, v domain.Vaccination) error {
	sh, err := s.repo.GetSheepByID(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if sh.OwnerUserID != userID {
		return domain.ErrUnauthorized
	}
	if index < 0 || index >= len(sh.Vaccinations) {
		return domain.ErrNotFound
	}
	sh.Vaccinations[index] = v
	return s.repo.UpdateSheep(ctx, sh)
}

// DeleteVaccination removes a vaccination record by index.
func (s *SheepService) DeleteVaccination(ctx context.Context, userID, sheepID uint, index int) error {
	sh, err := s.repo.GetSheepByID(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if sh.OwnerUserID != userID {
		return domain.ErrUnauthorized
	}
	if index < 0 || index >= len(sh.Vaccinations) {
		return domain.ErrNotFound
	}
	sh.Vaccinations = append(sh.Vaccinations[:index], sh.Vaccinations[index+1:]...)
	return s.repo.UpdateSheep(ctx, sh)
}

// UpdateTreatment updates a treatment record by index.
func (s *SheepService) UpdateTreatment(ctx context.Context, userID, sheepID uint, index int, t domain.Treatment) error {
	sh, err := s.repo.GetSheepByID(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if sh.OwnerUserID != userID {
		return domain.ErrUnauthorized
	}
	return s.treatRepo.UpdateTreatment(ctx, userID, sheepID, index, t)
}

// DeleteTreatment removes a treatment record by index.
func (s *SheepService) DeleteTreatment(ctx context.Context, userID, sheepID uint, index int) error {
	sh, err := s.repo.GetSheepByID(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if sh.OwnerUserID != userID {
		return domain.ErrUnauthorized
	}
	return s.treatRepo.DeleteTreatment(ctx, userID, sheepID, index)
}

// UpdateLambing updates a lambing record by index.
func (s *SheepService) UpdateLambing(ctx context.Context, userID, sheepID uint, index int, l domain.Lambing) error {
	sh, err := s.repo.GetSheepByID(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if sh.OwnerUserID != userID {
		return domain.ErrUnauthorized
	}
	return s.lambRepo.UpdateLambing(ctx, userID, sheepID, index, l)
}

// DeleteLambing removes a lambing record by index.
func (s *SheepService) DeleteLambing(ctx context.Context, userID, sheepID uint, index int) error {
	sh, err := s.repo.GetSheepByID(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if sh.OwnerUserID != userID {
		return domain.ErrUnauthorized
	}
	return s.lambRepo.DeleteLambing(ctx, userID, sheepID, index)
}
