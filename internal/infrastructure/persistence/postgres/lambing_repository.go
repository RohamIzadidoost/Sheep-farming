package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

var _ ports.LambingRepository = &LambingRepository{}

type LambingRepository struct{ db *gorm.DB }

func NewLambingRepository(db *gorm.DB) *LambingRepository { return &LambingRepository{db: db} }

func (r *LambingRepository) GetLambings(ctx context.Context, userID string, from, to *time.Time) ([]domain.Lambing, error) {
	var list []domain.Lambing
	q := r.db.WithContext(ctx).Joins("JOIN sheep ON sheep.id = lambings.sheep_id").Where("sheep.owner_user_id = ?", userID)
	if from != nil {
		q = q.Where("lambings.date >= ?", *from)
	}
	if to != nil {
		q = q.Where("lambings.date <= ?", *to)
	}
	err := q.Find(&list).Error
	return list, err
}
