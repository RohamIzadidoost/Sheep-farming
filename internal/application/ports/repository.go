package ports

import (
	"context"
	"sheep_farm_backend_go/internal/domain"
)

// SheepRepository defines the interface for sheep data operations.
// This is a "driven port" (output port).
type SheepRepository interface {
	CreateSheep(ctx context.Context, sheep *domain.Sheep) error
	GetSheepByID(ctx context.Context, userID, sheepID string) (*domain.Sheep, error)
	GetAllSheep(ctx context.Context, userID string) ([]domain.Sheep, error)
	UpdateSheep(ctx context.Context, sheep *domain.Sheep) error
	DeleteSheep(ctx context.Context, userID, sheepID string) error
}

// VaccineRepository defines the interface for vaccine definition operations.
type VaccineRepository interface {
	CreateVaccine(ctx context.Context, vaccine *domain.Vaccine) error
	GetVaccineByID(ctx context.Context, userID, vaccineID string) (*domain.Vaccine, error)
	GetAllVaccines(ctx context.Context, userID string) ([]domain.Vaccine, error)
	UpdateVaccine(ctx context.Context, vaccine *domain.Vaccine) error
	DeleteVaccine(ctx context.Context, userID, vaccineID string) error
}

// VaccinationRepository defines operations for vaccination records.
type VaccinationRepository interface {
	CreateVaccination(ctx context.Context, userID, sheepID string, v domain.Vaccination) error
	GetVaccinations(ctx context.Context, userID, sheepID string) ([]domain.Vaccination, error)
	DeleteVaccination(ctx context.Context, userID, sheepID string, index int) error
}
