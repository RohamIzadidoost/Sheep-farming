package services

import (
	"context"
	"time"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

// LambingService provides operations on lambing records
type LambingService struct {
	repo ports.LambingRepository
}

func NewLambingService(repo ports.LambingRepository) *LambingService {
	return &LambingService{repo: repo}
}

func (s *LambingService) Create(ctx context.Context, userID, sheepID uint, l domain.Lambing) error {
	return s.repo.AddLambing(ctx, userID, sheepID, l)
}

func (s *LambingService) List(ctx context.Context, userID uint, from, to *time.Time) ([]domain.Lambing, error) {
	return s.repo.FilterLambings(ctx, userID, from, to)
}

func (s *LambingService) Update(ctx context.Context, userID, sheepID uint, index int, l domain.Lambing) error {
	return s.repo.UpdateLambing(ctx, userID, sheepID, index, l)
}

func (s *LambingService) Delete(ctx context.Context, userID, sheepID uint, index int) error {
	return s.repo.DeleteLambing(ctx, userID, sheepID, index)
}
