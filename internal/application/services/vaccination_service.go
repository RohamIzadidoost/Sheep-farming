package services

import (
	"context"
	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

// VaccinationService provides operations on vaccination records.
type VaccinationService struct {
	repo ports.VaccinationRepository
}

func NewVaccinationService(repo ports.VaccinationRepository) *VaccinationService {
	return &VaccinationService{repo: repo}
}

func (s *VaccinationService) AddVaccination(ctx context.Context, userID, sheepID uint, v domain.Vaccination) error {
	return s.repo.CreateVaccination(ctx, userID, sheepID, v)
}

func (s *VaccinationService) ListVaccinations(ctx context.Context, userID, sheepID uint) ([]domain.Vaccination, error) {
	return s.repo.GetVaccinations(ctx, userID, sheepID)
}

func (s *VaccinationService) DeleteVaccination(ctx context.Context, userID, sheepID uint, index int) error {
	return s.repo.DeleteVaccination(ctx, userID, sheepID, index)
}
