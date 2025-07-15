package ports

import (
	"context"
	"sheep_farm_backend_go/internal/domain"
	"time"
)

// SheepRepository defines the interface for sheep data operations.
// This is a "driven port" (output port).
type SheepRepository interface {
	CreateSheep(ctx context.Context, sheep *domain.Sheep) error
	GetSheepByID(ctx context.Context, userID, sheepID uint) (*domain.Sheep, error)
	GetAllSheep(ctx context.Context, userID uint) ([]domain.Sheep, error)
	// FilterSheep returns sheep filtered by gender and age in days
	FilterSheep(ctx context.Context, userID uint, gender *string, minAgeDays, maxAgeDays *int) ([]domain.Sheep, error)
	UpdateSheep(ctx context.Context, sheep *domain.Sheep) error
	DeleteSheep(ctx context.Context, userID, sheepID uint) error
}

// VaccineRepository defines the interface for vaccine definition operations.
type VaccineRepository interface {
	CreateVaccine(ctx context.Context, vaccine *domain.Vaccine) error
	GetVaccineByID(ctx context.Context, userID, vaccineID uint) (*domain.Vaccine, error)
	GetAllVaccines(ctx context.Context, userID uint) ([]domain.Vaccine, error)
	UpdateVaccine(ctx context.Context, vaccine *domain.Vaccine) error
	DeleteVaccine(ctx context.Context, userID, vaccineID uint) error
}

// VaccinationRepository defines operations for vaccination records.
type VaccinationRepository interface {
	CreateVaccination(ctx context.Context, userID, sheepID uint, v domain.Vaccination) error
	GetVaccinations(ctx context.Context, userID, sheepID uint) ([]domain.Vaccination, error)
	DeleteVaccination(ctx context.Context, userID, sheepID uint, index int) error
}

// LambingRepository defines operations for lambings
// Additional methods for filtering by date or related events can be implemented
// as needed.
type LambingRepository interface {
	AddLambing(ctx context.Context, userID, sheepID uint, l domain.Lambing) error
	GetLambings(ctx context.Context, userID, sheepID uint) ([]domain.Lambing, error)
	UpdateLambing(ctx context.Context, userID, sheepID uint, index int, l domain.Lambing) error
	DeleteLambing(ctx context.Context, userID, sheepID uint, index int) error
	// FilterLambings retrieves lambings across all sheep within a date range
	FilterLambings(ctx context.Context, userID uint, from, to *time.Time) ([]domain.Lambing, error)
}

// TreatmentRepository defines operations for treatments
// Basic retrieval for a sheep
type TreatmentRepository interface {
	AddTreatment(ctx context.Context, userID, sheepID uint, t domain.Treatment) error
	GetTreatments(ctx context.Context, userID, sheepID uint) ([]domain.Treatment, error)
	UpdateTreatment(ctx context.Context, userID, sheepID uint, index int, t domain.Treatment) error
	DeleteTreatment(ctx context.Context, userID, sheepID uint, index int) error
	// FilterTreatments retrieves treatments across all sheep within a date range
	FilterTreatments(ctx context.Context, userID uint, from, to *time.Time) ([]domain.Treatment, error)
}
