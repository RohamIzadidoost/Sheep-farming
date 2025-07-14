package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

var _ ports.SheepRepository = &GormRepository{}
var _ ports.VaccineRepository = &GormRepository{}
var _ ports.VaccinationRepository = &GormRepository{}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// --- SheepRepository Implementation ---

func (r *GormRepository) CreateSheep(ctx context.Context, sheep *domain.Sheep) error {
	return r.db.WithContext(ctx).Create(sheep).Error
}

func (r *GormRepository) GetSheepByID(ctx context.Context, userID, sheepID string) (*domain.Sheep, error) {
	var sheep domain.Sheep
	err := r.db.WithContext(ctx).
		Where("owner_user_id = ? AND id = ?", userID, sheepID).
		Preload("Lambings").
		Preload("Vaccinations").
		Preload("Treatments").
		First(&sheep).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &sheep, err
}

func (r *GormRepository) GetAllSheep(ctx context.Context, userID string) ([]domain.Sheep, error) {
	var sheep []domain.Sheep
	err := r.db.WithContext(ctx).Where("owner_user_id = ?", userID).Find(&sheep).Error
	return sheep, err
}

func (r *GormRepository) UpdateSheep(ctx context.Context, sheep *domain.Sheep) error {
	sheep.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(sheep).Error
}

func (r *GormRepository) DeleteSheep(ctx context.Context, userID, sheepID string) error {
	return r.db.WithContext(ctx).
		Where("owner_user_id = ? AND id = ?", userID, sheepID).
		Delete(&domain.Sheep{}).Error
}

// --- VaccineRepository Implementation ---

func (r *GormRepository) CreateVaccine(ctx context.Context, vaccine *domain.Vaccine) error {
	return r.db.WithContext(ctx).Create(vaccine).Error
}

func (r *GormRepository) GetVaccineByID(ctx context.Context, userID, vaccineID string) (*domain.Vaccine, error) {
	var vaccine domain.Vaccine
	err := r.db.WithContext(ctx).
		Where("owner_user_id = ? AND id = ?", userID, vaccineID).
		First(&vaccine).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &vaccine, err
}

func (r *GormRepository) GetAllVaccines(ctx context.Context, userID string) ([]domain.Vaccine, error) {
	var vaccines []domain.Vaccine
	err := r.db.WithContext(ctx).Where("owner_user_id = ?", userID).Find(&vaccines).Error
	return vaccines, err
}

func (r *GormRepository) UpdateVaccine(ctx context.Context, vaccine *domain.Vaccine) error {
	vaccine.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(vaccine).Error
}

func (r *GormRepository) DeleteVaccine(ctx context.Context, userID, vaccineID string) error {
	return r.db.WithContext(ctx).
		Where("owner_user_id = ? AND id = ?", userID, vaccineID).
		Delete(&domain.Vaccine{}).Error
}

// --- VaccinationRepository Implementation ---

func (r *GormRepository) CreateVaccination(ctx context.Context, userID, sheepID string, v domain.Vaccination) error {
	// You may validate sheep ownership here if necessary
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(&v).Error
}

func (r *GormRepository) GetVaccinations(ctx context.Context, userID, sheepID string) ([]domain.Vaccination, error) {
	var vaccinations []domain.Vaccination
	err := r.db.WithContext(ctx).
		Joins("JOIN sheep ON sheep.id = vaccinations.sheep_id").
		Where("sheep.owner_user_id = ? AND sheep.id = ?", userID, sheepID).
		Find(&vaccinations).Error
	return vaccinations, err
}

func (r *GormRepository) DeleteVaccination(ctx context.Context, userID, sheepID string, index int) error {
	vaccinations, err := r.GetVaccinations(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(vaccinations) {
		return fmt.Errorf("vaccination index out of bounds")
	}
	return r.db.WithContext(ctx).Delete(&vaccinations[index]).Error
}