package services

import (
	"context"
	"time"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

// VaccineService provides use cases for vaccine definitions.
type VaccineService struct {
	repo ports.VaccineRepository
}

// NewVaccineService creates a new VaccineService instance.
func NewVaccineService(repo ports.VaccineRepository) *VaccineService {
	return &VaccineService{repo: repo}
}

// CreateVaccine handles the creation of a new vaccine definition.
func (s *VaccineService) CreateVaccine(ctx context.Context, vaccine *domain.Vaccine) error {
	return s.repo.CreateVaccine(ctx, vaccine)
}

// GetVaccineByID retrieves a vaccine definition by its ID.
func (s *VaccineService) GetVaccineByID(ctx context.Context, userID, vaccineID uint) (*domain.Vaccine, error) {
	return s.repo.GetVaccineByID(ctx, userID, vaccineID)
}

// GetAllVaccines retrieves all vaccine definitions for a given user.
func (s *VaccineService) GetAllVaccines(ctx context.Context, userID uint) ([]domain.Vaccine, error) {
	return s.repo.GetAllVaccines(ctx, userID)
}

// UpdateVaccine updates an existing vaccine definition.
func (s *VaccineService) UpdateVaccine(ctx context.Context, vaccine *domain.Vaccine) error {
	existingVaccine, err := s.repo.GetVaccineByID(ctx, vaccine.OwnerUserID, vaccine.ID)
	if err != nil {
		return err
	}
	if existingVaccine.OwnerUserID != vaccine.OwnerUserID {
		return domain.ErrUnauthorized
	}
	vaccine.UpdatedAt = time.Now()
	return s.repo.UpdateVaccine(ctx, vaccine)
}

// DeleteVaccine deletes a vaccine definition by its ID.
func (s *VaccineService) DeleteVaccine(ctx context.Context, userID, vaccineID uint) error {
	return s.repo.DeleteVaccine(ctx, userID, vaccineID)
}
