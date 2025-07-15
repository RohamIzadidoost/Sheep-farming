package dto

import "sheep_farm_backend_go/internal/domain"

// CreateVaccineRequest represents the data for creating a new vaccine definition.
type CreateVaccineRequest struct {
	Name           string `json:"name"`
	IntervalMonths int    `json:"intervalMonths"`
}

// UpdateVaccineRequest represents the data for updating an existing vaccine definition.
type UpdateVaccineRequest struct {
	Name           *string `json:"name,omitempty"`
	IntervalMonths *int    `json:"intervalMonths,omitempty"`
}

// VaccineResponse represents the vaccine definition data returned in API responses.
type VaccineResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	IntervalMonths int    `json:"intervalMonths"`
}

// ToDomain converts CreateVaccineRequest to domain.Vaccine
func (req *CreateVaccineRequest) ToDomain(ownerUserID uint) *domain.Vaccine {
	return &domain.Vaccine{
		Name:           req.Name,
		IntervalMonths: req.IntervalMonths,
		OwnerUserID:    ownerUserID,
	}
}

// ToResponse converts domain.Vaccine to VaccineResponse
func ToVaccineResponse(v *domain.Vaccine) *VaccineResponse {
	return &VaccineResponse{
		ID:             v.ID,
		Name:           v.Name,
		IntervalMonths: v.IntervalMonths,
	}
}
