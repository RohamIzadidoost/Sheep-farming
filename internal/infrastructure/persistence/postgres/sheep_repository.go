package postgres

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"sheep_farm_backend_go/internal/application/ports"
	"sheep_farm_backend_go/internal/domain"
)

var _ ports.SheepRepository = &SheepRepository{}
var _ ports.VaccinationRepository = &SheepRepository{}

type SheepRepository struct {
	db *gorm.DB
}

func NewSheepRepository(db *gorm.DB) *SheepRepository {
	return &SheepRepository{db: db}
}

func (r *SheepRepository) CreateSheep(ctx context.Context, s *domain.Sheep) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *SheepRepository) GetSheepByID(ctx context.Context, userID, id string) (*domain.Sheep, error) {
	var s domain.Sheep
	err := r.db.WithContext(ctx).Where("owner_user_id = ? AND id = ?", userID, id).
		Preload("Lambings").Preload("Vaccinations").Preload("Treatments").First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &s, err
}

func (r *SheepRepository) GetAllSheep(ctx context.Context, userID string) ([]domain.Sheep, error) {
	var list []domain.Sheep
	err := r.db.WithContext(ctx).Where("owner_user_id = ?", userID).Find(&list).Error
	return list, err
}

func (r *SheepRepository) FilterSheep(ctx context.Context, userID string, gender *string, minAgeDays, maxAgeDays *int) ([]domain.Sheep, error) {
	var list []domain.Sheep
	q := r.db.WithContext(ctx).Where("owner_user_id = ?", userID)
	if gender != nil {
		q = q.Where("gender = ?", *gender)
	}
	if minAgeDays != nil || maxAgeDays != nil {
		now := time.Now()
		if minAgeDays != nil {
			q = q.Where("date_of_birth <= ?", now.AddDate(0, 0, -*minAgeDays))
		}
		if maxAgeDays != nil {
			q = q.Where("date_of_birth >= ?", now.AddDate(0, 0, -*maxAgeDays))
		}
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *SheepRepository) UpdateSheep(ctx context.Context, s *domain.Sheep) error {
	s.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *SheepRepository) DeleteSheep(ctx context.Context, userID, id string) error {
	return r.db.WithContext(ctx).Where("owner_user_id = ? AND id = ?", userID, id).Delete(&domain.Sheep{}).Error
}

// VaccinationRepository
func (r *SheepRepository) CreateVaccination(ctx context.Context, userID, sheepID string, v domain.Vaccination) error {
	v.SheepID = 0
	return r.db.WithContext(ctx).Model(&domain.Sheep{}).
		Where("owner_user_id = ? AND id = ?", userID, sheepID).
		Association("Vaccinations").Append(&v)
}

func (r *SheepRepository) GetVaccinations(ctx context.Context, userID, sheepID string) ([]domain.Vaccination, error) {
	sheep, err := r.GetSheepByID(ctx, userID, sheepID)
	if err != nil {
		return nil, err
	}
	return sheep.Vaccinations, nil
}

func (r *SheepRepository) DeleteVaccination(ctx context.Context, userID, sheepID string, index int) error {
	sheep, err := r.GetSheepByID(ctx, userID, sheepID)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(sheep.Vaccinations) {
		return errors.New("index out of range")
	}
	return r.db.WithContext(ctx).Delete(&sheep.Vaccinations[index]).Error
}
