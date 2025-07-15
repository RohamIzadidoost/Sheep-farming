package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/domain"
	"sheep_farm_backend_go/internal/infrastructure/http/dto"
	"sheep_farm_backend_go/internal/infrastructure/http/middleware"
)

// VaccineHandler handles HTTP requests related to vaccine definitions.
type VaccineHandler struct {
	vaccineService *services.VaccineService
}

// NewVaccineHandler creates a new VaccineHandler.
func NewVaccineHandler(vaccineService *services.VaccineService) *VaccineHandler {
	return &VaccineHandler{vaccineService: vaccineService}
}

// CreateVaccine handles POST /vaccines requests.
func (h *VaccineHandler) CreateVaccine(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateVaccineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vaccine := req.ToDomain(userID)
	if err := h.vaccineService.CreateVaccine(r.Context(), vaccine); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ToVaccineResponse(vaccine)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetVaccineByID handles GET /vaccines/{id} requests.
func (h *VaccineHandler) GetVaccineByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	vaccineID, _ := strconv.ParseUint(idStr, 10, 64)

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vaccine, err := h.vaccineService.GetVaccineByID(r.Context(), userID, uint(vaccineID))
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ToVaccineResponse(vaccine)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetAllVaccines handles GET /vaccines requests.
func (h *VaccineHandler) GetAllVaccines(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vaccineList, err := h.vaccineService.GetAllVaccines(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.VaccineResponse
	for _, vaccine := range vaccineList {
		responses = append(responses, *dto.ToVaccineResponse(&vaccine))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// UpdateVaccine handles PUT /vaccines/{id} requests.
func (h *VaccineHandler) UpdateVaccine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	vaccineID, _ := strconv.ParseUint(idStr, 10, 64)

	var req dto.UpdateVaccineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	existingVaccine, err := h.vaccineService.GetVaccineByID(r.Context(), userID, uint(vaccineID))
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Name != nil {
		existingVaccine.Name = *req.Name
	}
	if req.IntervalMonths != nil {
		existingVaccine.IntervalMonths = *req.IntervalMonths
	}

	if err := h.vaccineService.UpdateVaccine(r.Context(), existingVaccine); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ToVaccineResponse(existingVaccine)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteVaccine handles DELETE /vaccines/{id} requests.
func (h *VaccineHandler) DeleteVaccine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	vaccineID, _ := strconv.ParseUint(idStr, 10, 64)

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := h.vaccineService.DeleteVaccine(r.Context(), userID, uint(vaccineID)); err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
