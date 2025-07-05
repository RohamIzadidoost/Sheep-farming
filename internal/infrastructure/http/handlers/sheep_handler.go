package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/domain"
	"sheep_farm_backend_go/internal/infrastructure/http/dto" // Import DTOs
	"sheep_farm_backend_go/internal/infrastructure/http/middleware"
)

// SheepHandler handles HTTP requests related to sheep.
type SheepHandler struct {
	sheepService *services.SheepService
}

// NewSheepHandler creates a new SheepHandler.
func NewSheepHandler(sheepService *services.SheepService) *SheepHandler {
	return &SheepHandler{sheepService: sheepService}
}

// CreateSheep handles POST /sheep requests.
func (h *SheepHandler) CreateSheep(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSheepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sheep := req.ToDomain(userID) // Convert DTO to domain entity
	if err := h.sheepService.CreateSheep(r.Context(), sheep); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ToSheepResponse(sheep) // Convert domain entity to DTO for response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetSheepByID handles GET /sheep/{id} requests.
func (h *SheepHandler) GetSheepByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sheep, err := h.sheepService.GetSheepByID(r.Context(), userID, sheepID)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ToSheepResponse(sheep)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetAllSheep handles GET /sheep requests.
func (h *SheepHandler) GetAllSheep(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sheepList, err := h.sheepService.GetAllSheep(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert domain entities to DTOs for response
	var responses []dto.SheepResponse
	for _, sheep := range sheepList {
		responses = append(responses, *dto.ToSheepResponse(&sheep))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// UpdateSheep handles PUT /sheep/{id} requests.
func (h *SheepHandler) UpdateSheep(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]

	var req dto.UpdateSheepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	// Get existing sheep to apply updates
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	existingSheep, err := h.sheepService.GetSheepByID(r.Context(), userID, sheepID)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Apply updates from DTO to domain entity
	if req.Name != nil {
		existingSheep.Name = *req.Name
	}
	if req.Gender != nil {
		existingSheep.Gender = *req.Gender
	}
	if req.DateOfBirth != nil {
		existingSheep.DateOfBirth = time.Time(*req.DateOfBirth)
	}
	// Handle nullable pointer to pointer for dates. If `nil` means no change, `&DateOnly(time.Time{})` means set to null.
	if req.BreedingDate != nil {
		if *req.BreedingDate == nil { // Explicitly set to null
			existingSheep.BreedingDate = nil
		} else {
			t := time.Time(**req.BreedingDate)
			existingSheep.BreedingDate = &t
		}
	}
	if req.LastShearingDate != nil {
		if *req.LastShearingDate == nil {
			existingSheep.LastShearingDate = nil
		} else {
			t := time.Time(**req.LastShearingDate)
			existingSheep.LastShearingDate = &t
		}
	}
	if req.LastHoofTrimDate != nil {
		if *req.LastHoofTrimDate == nil {
			existingSheep.LastHoofTrimDate = nil
		} else {
			t := time.Time(**req.LastHoofTrimDate)
			existingSheep.LastHoofTrimDate = &t
		}
	}
	if req.PhotoURL != nil {
		existingSheep.PhotoURL = *req.PhotoURL
	}
	if req.Vaccinations != nil {
		domainVaccinations := make([]domain.Vaccination, len(*req.Vaccinations))
		for i, v := range *req.Vaccinations {
			domainVaccinations[i] = domain.Vaccination{
				VaccineID:   v.VaccineID,
				Date:        time.Time(v.Date),
				Description: v.Description,
			}
		}
		existingSheep.Vaccinations = domainVaccinations
	}
	if req.Treatments != nil {
		domainTreatments := make([]domain.Treatment, len(*req.Treatments))
		for i, t := range *req.Treatments {
			domainTreatments[i] = domain.Treatment{
				Date:        time.Time(t.Date),
				Description: t.Description,
			}
		}
		existingSheep.Treatments = domainTreatments
	}

	if err := h.sheepService.UpdateSheep(r.Context(), existingSheep); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ToSheepResponse(existingSheep)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteSheep handles DELETE /sheep/{id} requests.
func (h *SheepHandler) DeleteSheep(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := h.sheepService.DeleteSheep(r.Context(), userID, sheepID); err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
