package dto

import (
	"sheep_farm_backend_go/internal/domain"
	"time"
)

type VaccinationDTO struct {
	VaccineID   string   `json:"vaccineId"`
	Date        DateOnly `json:"date"`
	Description string   `json:"description,omitempty"`
}

type TreatmentDTO struct {
	Date        DateOnly `json:"date"`
	Description string   `json:"description"`
}

// CreateSheepRequest represents the data for creating a new sheep.
type CreateSheepRequest struct {
	Name             string           `json:"name"`
	Gender           string           `json:"gender"`
	DateOfBirth      DateOnly         `json:"dateOfBirth"`
	BreedingDate     *DateOnly        `json:"breedingDate,omitempty"`
	LastShearingDate *DateOnly        `json:"lastShearingDate,omitempty"`
	LastHoofTrimDate *DateOnly        `json:"lastHoofTrimDate,omitempty"`
	PhotoURL         string           `json:"photoUrl,omitempty"`
	Vaccinations     []VaccinationDTO `json:"vaccinations"`
	Treatments       []TreatmentDTO   `json:"treatments"`
}

// UpdateSheepRequest represents the data for updating an existing sheep.
type UpdateSheepRequest struct {
	Name             *string           `json:"name,omitempty"` // Pointer for optional updates
	Gender           *string           `json:"gender,omitempty"`
	DateOfBirth      *DateOnly         `json:"dateOfBirth,omitempty"`
	BreedingDate     **DateOnly        `json:"breedingDate,omitempty"` // Pointer to pointer to handle explicit null/empty date
	LastShearingDate **DateOnly        `json:"lastShearingDate,omitempty"`
	LastHoofTrimDate **DateOnly        `json:"lastHoofTrimDate,omitempty"`
	PhotoURL         *string           `json:"photoUrl,omitempty"`
	Vaccinations     *[]VaccinationDTO `json:"vaccinations,omitempty"`
	Treatments       *[]TreatmentDTO   `json:"treatments,omitempty"`
}

// SheepResponse represents the sheep data returned in API responses.
type SheepResponse struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Gender           string           `json:"gender"`
	DateOfBirth      DateOnly         `json:"dateOfBirth"`
	BreedingDate     *DateOnly        `json:"breedingDate,omitempty"`
	LastShearingDate *DateOnly        `json:"lastShearingDate,omitempty"`
	LastHoofTrimDate *DateOnly        `json:"lastHoofTrimDate,omitempty"`
	PhotoURL         string           `json:"photoUrl,omitempty"`
	Vaccinations     []VaccinationDTO `json:"vaccinations"`
	Treatments       []TreatmentDTO   `json:"treatments"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
}

// ToDomain converts CreateSheepRequest to domain.Sheep
func (req *CreateSheepRequest) ToDomain(ownerUserID string) *domain.Sheep {
	domainVaccinations := make([]domain.Vaccination, len(req.Vaccinations))
	for i, v := range req.Vaccinations {
		domainVaccinations[i] = domain.Vaccination{
			VaccineID:   v.VaccineID,
			Date:        time.Time(v.Date),
			Description: v.Description,
		}
	}

	domainTreatments := make([]domain.Treatment, len(req.Treatments))
	for i, t := range req.Treatments {
		domainTreatments[i] = domain.Treatment{
			Date:        time.Time(t.Date),
			Description: t.Description,
		}
	}

	return &domain.Sheep{
		Name:             req.Name,
		Gender:           req.Gender,
		DateOfBirth:      time.Time(req.DateOfBirth),
		BreedingDate:     req.BreedingDate.ToTimePtr(),
		LastShearingDate: req.LastShearingDate.ToTimePtr(),
		LastHoofTrimDate: req.LastHoofTrimDate.ToTimePtr(),
		PhotoURL:         req.PhotoURL,
		Vaccinations:     domainVaccinations,
		Treatments:       domainTreatments,
		OwnerUserID:      ownerUserID,
	}
}

// ToResponse converts domain.Sheep to SheepResponse
func ToSheepResponse(s *domain.Sheep) *SheepResponse {
	responseVaccinations := make([]VaccinationDTO, len(s.Vaccinations))
	for i, v := range s.Vaccinations {
		responseVaccinations[i] = VaccinationDTO{
			VaccineID:   v.VaccineID,
			Date:        DateOnly(v.Date),
			Description: v.Description,
		}
	}

	responseTreatments := make([]TreatmentDTO, len(s.Treatments))
	for i, t := range s.Treatments {
		responseTreatments[i] = TreatmentDTO{
			Date:        DateOnly(t.Date),
			Description: t.Description,
		}
	}

	return &SheepResponse{
		ID:               s.ID,
		Name:             s.Name,
		Gender:           s.Gender,
		DateOfBirth:      DateOnly(s.DateOfBirth),
		BreedingDate:     FromTimePtrPtr(s.BreedingDate),
		LastShearingDate: FromTimePtrPtr(s.LastShearingDate),
		LastHoofTrimDate: FromTimePtrPtr(s.LastHoofTrimDate),
		PhotoURL:         s.PhotoURL,
		Vaccinations:     responseVaccinations,
		Treatments:       responseTreatments,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}

// FromTimePtrPtr converts *time.Time to *DateOnly
func FromTimePtrPtr(t *time.Time) *DateOnly {
	if t == nil {
		return nil
	}
	d := DateOnly(*t)
	return &d
}
