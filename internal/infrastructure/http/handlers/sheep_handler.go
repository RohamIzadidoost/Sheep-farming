package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
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

	var gender *string
	if g := r.URL.Query().Get("gender"); g != "" {
		gender = &g
	}

	var minAgeDaysPtr, maxAgeDaysPtr *int
	if v := r.URL.Query().Get("minAgeDays"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			minAgeDaysPtr = &n
		}
	}
	if v := r.URL.Query().Get("maxAgeDays"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			maxAgeDaysPtr = &n
		}
	}

	sheepList, err := h.sheepService.FilterSheep(r.Context(), userID, gender, minAgeDaysPtr, maxAgeDaysPtr)
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
	if req.EarNumber1 != nil {
		existingSheep.EarNumber1 = *req.EarNumber1
	}
	if req.EarNumber2 != nil {
		existingSheep.EarNumber2 = *req.EarNumber2
	}
	if req.EarNumber3 != nil {
		existingSheep.EarNumber3 = *req.EarNumber3
	}
	if req.NeckNumber != nil {
		existingSheep.NeckNumber = *req.NeckNumber
	}
	if req.FatherGen != nil {
		existingSheep.FatherGen = *req.FatherGen
	}
	if req.BirthWeight != nil {
		existingSheep.BirthWeight = *req.BirthWeight
	}
	if req.Gender != nil {
		existingSheep.Gender = *req.Gender
	}
	if req.ReproductionState != nil {
		existingSheep.ReproductionState = *req.ReproductionState
	}
	if req.HealthState != nil {
		existingSheep.HealthState = *req.HealthState
	}
	if req.DateOfBirth != nil {
		existingSheep.DateOfBirth = time.Time(*req.DateOfBirth)
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
				Date:        time.Time(v.Date),
				Vaccine:     v.Vaccine,
				Vaccinator:  v.Vaccinator,
				Description: v.Description,
			}
		}
		existingSheep.Vaccinations = domainVaccinations
	}
	if req.Treatments != nil {
		domainTreatments := make([]domain.Treatment, len(*req.Treatments))
		for i, t := range *req.Treatments {
			domainTreatments[i] = domain.Treatment{
				Date:               time.Time(t.Date),
				DiseaseDescription: t.DiseaseDescription,
				TreatDescription:   t.TreatDescription,
			}
		}
		existingSheep.Treatments = domainTreatments
	}
	if req.Lambings != nil {
		domainLambings := make([]domain.Lambing, len(*req.Lambings))
		for i, l := range *req.Lambings {
			domainLambings[i] = domain.Lambing{
				Date:    time.Time(l.Date),
				NumBorn: l.NumBorn,
				Sexes:   l.Sexes,
				NumDead: l.NumDead,
			}
		}
		existingSheep.Lambings = domainLambings
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

// AddVaccination handles POST /sheep/{id}/vaccinations requests.
func (h *SheepHandler) AddVaccination(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]

	var req dto.VaccinationDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = h.sheepService.AddVaccination(r.Context(), userID, sheepID, domain.Vaccination{
		Date:        time.Time(req.Date),
		Vaccine:     req.Vaccine,
		Vaccinator:  req.Vaccinator,
		Description: req.Description,
	})
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// AddTreatment handles POST /sheep/{id}/treatments requests.
func (h *SheepHandler) AddTreatment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]

	var req dto.TreatmentDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = h.sheepService.AddTreatment(r.Context(), userID, sheepID, domain.Treatment{
		Date:               time.Time(req.Date),
		DiseaseDescription: req.DiseaseDescription,
		TreatDescription:   req.TreatDescription,
	})
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// AddLambing handles POST /sheep/{id}/lambings requests.
func (h *SheepHandler) AddLambing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]

	var req dto.LambingDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = h.sheepService.AddLambing(r.Context(), userID, sheepID, domain.Lambing{
		Date:    time.Time(req.Date),
		NumBorn: req.NumBorn,
		Sexes:   req.Sexes,
		NumDead: req.NumDead,
	})
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateVaccination handles PUT /sheep/{id}/vaccinations/{idx} requests.
func (h *SheepHandler) UpdateVaccination(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]
	idxStr := vars["idx"]
	index, err := strconv.Atoi(idxStr)
	if err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	var req dto.VaccinationDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = h.sheepService.UpdateVaccination(r.Context(), userID, sheepID, index, domain.Vaccination{
		Date:        time.Time(req.Date),
		Vaccine:     req.Vaccine,
		Vaccinator:  req.Vaccinator,
		Description: req.Description,
	})
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteVaccination handles DELETE /sheep/{id}/vaccinations/{idx} requests.
func (h *SheepHandler) DeleteVaccination(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]
	idxStr := vars["idx"]
	index, err := strconv.Atoi(idxStr)
	if err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = h.sheepService.DeleteVaccination(r.Context(), userID, sheepID, index)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateTreatment handles PUT /sheep/{id}/treatments/{idx} requests.
func (h *SheepHandler) UpdateTreatment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]
	idxStr := vars["idx"]
	index, err := strconv.Atoi(idxStr)
	if err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	var req dto.TreatmentDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = h.sheepService.UpdateTreatment(r.Context(), userID, sheepID, index, domain.Treatment{
		Date:               time.Time(req.Date),
		DiseaseDescription: req.DiseaseDescription,
		TreatDescription:   req.TreatDescription,
	})
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteTreatment handles DELETE /sheep/{id}/treatments/{idx} requests.
func (h *SheepHandler) DeleteTreatment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]
	idxStr := vars["idx"]
	index, err := strconv.Atoi(idxStr)
	if err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = h.sheepService.DeleteTreatment(r.Context(), userID, sheepID, index)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateLambing handles PUT /sheep/{id}/lambings/{idx} requests.
func (h *SheepHandler) UpdateLambing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]
	idxStr := vars["idx"]
	index, err := strconv.Atoi(idxStr)
	if err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	var req dto.LambingDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = h.sheepService.UpdateLambing(r.Context(), userID, sheepID, index, domain.Lambing{
		Date:    time.Time(req.Date),
		NumBorn: req.NumBorn,
		Sexes:   req.Sexes,
		NumDead: req.NumDead,
	})
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteLambing handles DELETE /sheep/{id}/lambings/{idx} requests.
func (h *SheepHandler) DeleteLambing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sheepID := vars["id"]
	idxStr := vars["idx"]
	index, err := strconv.Atoi(idxStr)
	if err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = h.sheepService.DeleteLambing(r.Context(), userID, sheepID, index)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
