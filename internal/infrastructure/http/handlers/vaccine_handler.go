package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/domain"
	"sheep_farm_backend_go/internal/infrastructure/http/dto"
)

// VaccineHandler handles HTTP requests related to vaccine definitions.
type VaccineHandler struct {
	vaccineService *services.VaccineService
	fixedUserID    string // For simplicity, fixed user ID
}

// NewVaccineHandler creates a new VaccineHandler.
func NewVaccineHandler(vaccineService *services.VaccineService, fixedUserID string) *VaccineHandler {
	return &VaccineHandler{vaccineService: vaccineService, fixedUserID: fixedUserID}
}

// CreateVaccine handles POST /vaccines requests.
func (h *VaccineHandler) CreateVaccine(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateVaccineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	vaccine := req.ToDomain(h.fixedUserID)
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
	vaccineID := vars["id"]

	vaccine, err := h.vaccineService.GetVaccineByID(r.Context(), h.fixedUserID, vaccineID)
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
	vaccineList, err := h.vaccineService.GetAllVaccines(r.Context(), h.fixedUserID)
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
	vaccineID := vars["id"]

	var req dto.UpdateVaccineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	existingVaccine, err := h.vaccineService.GetVaccineByID(r.Context(), h.fixedUserID, vaccineID)
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
	vaccineID := vars["id"]

	if err := h.vaccineService.DeleteVaccine(r.Context(), h.fixedUserID, vaccineID); err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
