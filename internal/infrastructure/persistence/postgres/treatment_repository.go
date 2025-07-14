package postgres

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

var _ ports.TreatmentRepository = &TreatmentRepository{}

type TreatmentRepository struct{ db *gorm.DB }

func NewTreatmentRepository(db *gorm.DB) *TreatmentRepository { return &TreatmentRepository{db: db} }

func (r *TreatmentRepository) AddTreatment(ctx context.Context, userID, sheepID string, t domain.Treatment) error {
	var sheep domain.Sheep
	if err := r.db.WithContext(ctx).Where("owner_user_id = ? AND id = ?", userID, sheepID).First(&sheep).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	t.SheepID = sheep.Model.ID
	return r.db.WithContext(ctx).Create(&t).Error
}

func (r *TreatmentRepository) GetTreatments(ctx context.Context, userID, sheepID string) ([]domain.Treatment, error) {
	var list []domain.Treatment
	err := r.db.WithContext(ctx).
		Joins("JOIN sheep ON sheep.id = treatments.sheep_id").
		Where("sheep.owner_user_id = ? AND sheep.id = ?", userID, sheepID).
		Find(&list).Error
	return list, err
}

func (r *TreatmentRepository) UpdateTreatment(ctx context.Context, userID, sheepID string, index int, t domain.Treatment) error {
	treatments, err := r.GetTreatments(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(treatments) {
		return domain.ErrNotFound
	}
	t.ID = treatments[index].ID
	return r.db.WithContext(ctx).Model(&domain.Treatment{}).Where("id = ?", t.ID).Updates(t).Error
}

func (r *TreatmentRepository) DeleteTreatment(ctx context.Context, userID, sheepID string, index int) error {
	treatments, err := r.GetTreatments(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(treatments) {
		return domain.ErrNotFound
	}
	return r.db.WithContext(ctx).Delete(&domain.Treatment{}, treatments[index].ID).Error
}

func (r *TreatmentRepository) FilterTreatments(ctx context.Context, userID string, from, to *time.Time) ([]domain.Treatment, error) {
	var list []domain.Treatment
	q := r.db.WithContext(ctx).Joins("JOIN sheep ON sheep.id = treatments.sheep_id").Where("sheep.owner_user_id = ?", userID)
	if from != nil {
		q = q.Where("treatments.date >= ?", *from)
	}
	if to != nil {
		q = q.Where("treatments.date <= ?", *to)
	}
	err := q.Find(&list).Error
	return list, err
}
