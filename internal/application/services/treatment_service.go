package services

import (
	"context"
	"time"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

// TreatmentService provides operations on treatment records
type TreatmentService struct {
	repo ports.TreatmentRepository
}

func NewTreatmentService(repo ports.TreatmentRepository) *TreatmentService {
	return &TreatmentService{repo: repo}
}

func (s *TreatmentService) Create(ctx context.Context, userID, sheepID string, t domain.Treatment) error {
	return s.repo.AddTreatment(ctx, userID, sheepID, t)
}

func (s *TreatmentService) List(ctx context.Context, userID string, from, to *time.Time) ([]domain.Treatment, error) {
	return s.repo.FilterTreatments(ctx, userID, from, to)
}

func (s *TreatmentService) Update(ctx context.Context, userID, sheepID string, index int, t domain.Treatment) error {
	return s.repo.UpdateTreatment(ctx, userID, sheepID, index, t)
}

func (s *TreatmentService) Delete(ctx context.Context, userID, sheepID string, index int) error {
	return s.repo.DeleteTreatment(ctx, userID, sheepID, index)
}
