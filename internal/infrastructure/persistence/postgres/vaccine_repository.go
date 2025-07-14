package postgres

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

var _ ports.VaccineRepository = &VaccineRepository{}

type VaccineRepository struct {
	db *gorm.DB
}

func NewVaccineRepository(db *gorm.DB) *VaccineRepository {
	return &VaccineRepository{db: db}
}

func (r *VaccineRepository) CreateVaccine(ctx context.Context, v *domain.Vaccine) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *VaccineRepository) GetVaccineByID(ctx context.Context, userID, id string) (*domain.Vaccine, error) {
	var v domain.Vaccine
	err := r.db.WithContext(ctx).Where("owner_user_id = ? AND id = ?", userID, id).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &v, err
}

func (r *VaccineRepository) GetAllVaccines(ctx context.Context, userID string) ([]domain.Vaccine, error) {
	var list []domain.Vaccine
	err := r.db.WithContext(ctx).Where("owner_user_id = ?", userID).Find(&list).Error
	return list, err
}

func (r *VaccineRepository) UpdateVaccine(ctx context.Context, v *domain.Vaccine) error {
	v.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(v).Error
}

func (r *VaccineRepository) DeleteVaccine(ctx context.Context, userID, id string) error {
	return r.db.WithContext(ctx).Where("owner_user_id = ? AND id = ?", userID, id).Delete(&domain.Vaccine{}).Error
}
