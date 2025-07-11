package dto

import (
	"sheep_farm_backend_go/internal/domain"
	"time"
)

type VaccinationDTO struct {
	Date        DateOnly `json:"date"`
	Vaccine     string   `json:"vaccine"`
	Vaccinator  string   `json:"vaccinator"`
	Description string   `json:"description,omitempty"`
}

type LambingDTO struct {
	Date    DateOnly `json:"date"`
	NumBorn int      `json:"numBorn"`
	Sexes   []string `json:"sexes"`
	NumDead int      `json:"numDead"`
}

type TreatmentDTO struct {
	Date               DateOnly `json:"date"`
	DiseaseDescription string   `json:"diseaseDescription"`
	TreatDescription   string   `json:"treatDescription"`
}

// CreateSheepRequest represents the data for creating a new sheep.
type CreateSheepRequest struct {
	EarNumber1       string           `json:"earNumber1"`
	EarNumber2       string           `json:"earNumber2,omitempty"`
	EarNumber3       string           `json:"earNumber3,omitempty"`
	NeckNumber       *string          `json:"neckNumber,omitempty"`
	FatherGen        string           `json:"fatherGen,omitempty"`
	BirthWeight      float64          `json:"birthWeight,omitempty"`
	Gender           string           `json:"gender"`
	DateOfBirth      DateOnly         `json:"dateOfBirth"`
	LastShearingDate *DateOnly        `json:"lastShearingDate,omitempty"`
	LastHoofTrimDate *DateOnly        `json:"lastHoofTrimDate,omitempty"`
	PhotoURL         string           `json:"photoUrl,omitempty"`
	Lambings         []LambingDTO     `json:"lambings"`
	Vaccinations     []VaccinationDTO `json:"vaccinations"`
	Treatments       []TreatmentDTO   `json:"treatments"`
}

// UpdateSheepRequest represents the data for updating an existing sheep.
type UpdateSheepRequest struct {
	EarNumber1       *string           `json:"earNumber1,omitempty"`
	EarNumber2       *string           `json:"earNumber2,omitempty"`
	EarNumber3       *string           `json:"earNumber3,omitempty"`
	NeckNumber       **string          `json:"neckNumber,omitempty"`
	FatherGen        *string           `json:"fatherGen,omitempty"`
	BirthWeight      *float64          `json:"birthWeight,omitempty"`
	Gender           *string           `json:"gender,omitempty"`
	DateOfBirth      *DateOnly         `json:"dateOfBirth,omitempty"`
	LastShearingDate **DateOnly        `json:"lastShearingDate,omitempty"`
	LastHoofTrimDate **DateOnly        `json:"lastHoofTrimDate,omitempty"`
	PhotoURL         *string           `json:"photoUrl,omitempty"`
	Lambings         *[]LambingDTO     `json:"lambings,omitempty"`
	Vaccinations     *[]VaccinationDTO `json:"vaccinations,omitempty"`
	Treatments       *[]TreatmentDTO   `json:"treatments,omitempty"`
}

// SheepResponse represents the sheep data returned in API responses.
type SheepResponse struct {
	ID               string           `json:"id"`
	EarNumber1       string           `json:"earNumber1"`
	EarNumber2       string           `json:"earNumber2,omitempty"`
	EarNumber3       string           `json:"earNumber3,omitempty"`
	NeckNumber       *string          `json:"neckNumber,omitempty"`
	FatherGen        string           `json:"fatherGen,omitempty"`
	BirthWeight      float64          `json:"birthWeight,omitempty"`
	Gender           string           `json:"gender"`
	DateOfBirth      DateOnly         `json:"dateOfBirth"`
	LastShearingDate *DateOnly        `json:"lastShearingDate,omitempty"`
	LastHoofTrimDate *DateOnly        `json:"lastHoofTrimDate,omitempty"`
	PhotoURL         string           `json:"photoUrl,omitempty"`
	Lambings         []LambingDTO     `json:"lambings"`
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
			Date:        time.Time(v.Date),
			Vaccine:     v.Vaccine,
			Vaccinator:  v.Vaccinator,
			Description: v.Description,
		}
	}

	domainTreatments := make([]domain.Treatment, len(req.Treatments))
	for i, t := range req.Treatments {
		domainTreatments[i] = domain.Treatment{
			Date:               time.Time(t.Date),
			DiseaseDescription: t.DiseaseDescription,
			TreatDescription:   t.TreatDescription,
		}
	}

	domainLambings := make([]domain.Lambing, len(req.Lambings))
	for i, l := range req.Lambings {
		domainLambings[i] = domain.Lambing{
			Date:    time.Time(l.Date),
			NumBorn: l.NumBorn,
			Sexes:   l.Sexes,
			NumDead: l.NumDead,
		}
	}

	var lastShearingDate *time.Time
	if req.LastShearingDate != nil {
		lastShearingDate = req.LastShearingDate.ToTimePtr()
	}

	var lastHoofTrimDate *time.Time
	if req.LastHoofTrimDate != nil {
		lastHoofTrimDate = req.LastHoofTrimDate.ToTimePtr()
	}

	return &domain.Sheep{
		EarNumber1:       req.EarNumber1,
		EarNumber2:       req.EarNumber2,
		EarNumber3:       req.EarNumber3,
		NeckNumber:       req.NeckNumber,
		FatherGen:        req.FatherGen,
		BirthWeight:      req.BirthWeight,
		Gender:           req.Gender,
		DateOfBirth:      time.Time(req.DateOfBirth),
		LastShearingDate: lastShearingDate,
		LastHoofTrimDate: lastHoofTrimDate,
		PhotoURL:         req.PhotoURL,
		Lambings:         domainLambings,
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
			Date:        DateOnly(v.Date),
			Vaccine:     v.Vaccine,
			Vaccinator:  v.Vaccinator,
			Description: v.Description,
		}
	}

	responseTreatments := make([]TreatmentDTO, len(s.Treatments))
	for i, t := range s.Treatments {
		responseTreatments[i] = TreatmentDTO{
			Date:               DateOnly(t.Date),
			DiseaseDescription: t.DiseaseDescription,
			TreatDescription:   t.TreatDescription,
		}
	}

	responseLambings := make([]LambingDTO, len(s.Lambings))
	for i, l := range s.Lambings {
		responseLambings[i] = LambingDTO{
			Date:    DateOnly(l.Date),
			NumBorn: l.NumBorn,
			Sexes:   l.Sexes,
			NumDead: l.NumDead,
		}
	}

	return &SheepResponse{
		ID:               s.ID,
		EarNumber1:       s.EarNumber1,
		EarNumber2:       s.EarNumber2,
		EarNumber3:       s.EarNumber3,
		NeckNumber:       s.NeckNumber,
		FatherGen:        s.FatherGen,
		BirthWeight:      s.BirthWeight,
		Gender:           s.Gender,
		DateOfBirth:      DateOnly(s.DateOfBirth),
		LastShearingDate: FromTimePtrPtr(s.LastShearingDate),
		LastHoofTrimDate: FromTimePtrPtr(s.LastHoofTrimDate),
		PhotoURL:         s.PhotoURL,
		Lambings:         responseLambings,
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
