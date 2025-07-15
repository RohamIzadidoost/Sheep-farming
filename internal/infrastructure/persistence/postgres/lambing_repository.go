package postgres

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

var _ ports.LambingRepository = &LambingRepository{}

type LambingRepository struct{ db *gorm.DB }

func NewLambingRepository(db *gorm.DB) *LambingRepository { return &LambingRepository{db: db} }

func (r *LambingRepository) AddLambing(ctx context.Context, userID, sheepID uint, l domain.Lambing) error {
	var sheep domain.Sheep
	if err := r.db.WithContext(ctx).Where("owner_user_id = ? AND id = ?", userID, sheepID).First(&sheep).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	l.SheepID = sheep.Model.ID
	return r.db.WithContext(ctx).Create(&l).Error
}

func (r *LambingRepository) GetLambings(ctx context.Context, userID, sheepID uint) ([]domain.Lambing, error) {
	var list []domain.Lambing
	err := r.db.WithContext(ctx).
		Joins("JOIN sheep ON sheep.id = lambings.sheep_id").
		Where("sheep.owner_user_id = ? AND sheep.id = ?", userID, sheepID).
		Find(&list).Error
	return list, err
}

func (r *LambingRepository) UpdateLambing(ctx context.Context, userID, sheepID uint, index int, l domain.Lambing) error {
	lambings, err := r.GetLambings(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(lambings) {
		return domain.ErrNotFound
	}
	l.ID = lambings[index].ID
	return r.db.WithContext(ctx).Model(&domain.Lambing{}).Where("id = ?", l.ID).Updates(l).Error
}

func (r *LambingRepository) DeleteLambing(ctx context.Context, userID, sheepID uint, index int) error {
	lambings, err := r.GetLambings(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(lambings) {
		return domain.ErrNotFound
	}
	return r.db.WithContext(ctx).Delete(&domain.Lambing{}, lambings[index].ID).Error
}

func (r *LambingRepository) FilterLambings(ctx context.Context, userID uint, from, to *time.Time) ([]domain.Lambing, error) {
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
