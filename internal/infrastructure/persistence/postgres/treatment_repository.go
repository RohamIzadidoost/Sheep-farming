package postgres

import (
	"context"

	"gorm.io/gorm"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

var _ ports.TreatmentRepository = &TreatmentRepository{}

type TreatmentRepository struct{ db *gorm.DB }

func NewTreatmentRepository(db *gorm.DB) *TreatmentRepository { return &TreatmentRepository{db: db} }

func (r *TreatmentRepository) GetTreatments(ctx context.Context, userID, sheepID string) ([]domain.Treatment, error) {
	var list []domain.Treatment
	err := r.db.WithContext(ctx).
		Joins("JOIN sheep ON sheep.id = treatments.sheep_id").
		Where("sheep.owner_user_id = ? AND sheep.id = ?", userID, sheepID).
		Find(&list).Error
	return list, err
}
