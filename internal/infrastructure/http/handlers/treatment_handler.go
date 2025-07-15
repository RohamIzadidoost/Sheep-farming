package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/domain"
	"sheep_farm_backend_go/internal/infrastructure/http/dto"
	"sheep_farm_backend_go/internal/infrastructure/http/middleware"
)

// TreatmentHandler manages treatment endpoints
type TreatmentHandler struct {
	service *services.TreatmentService
}

func NewTreatmentHandler(s *services.TreatmentService) *TreatmentHandler {
	return &TreatmentHandler{service: s}
}

func (h *TreatmentHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var fromPtr, toPtr *time.Time
	if v := r.URL.Query().Get("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			fromPtr = &t
		}
	}
	if v := r.URL.Query().Get("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			toPtr = &t
		}
	}
	list, err := h.service.List(r.Context(), userID, fromPtr, toPtr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resp []dto.TreatmentDTO
	for _, t := range list {
		resp = append(resp, dto.TreatmentDTO{
			Date:               dto.DateOnly(t.Date),
			DiseaseDescription: t.DiseaseDescription,
			TreatDescription:   t.TreatDescription,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *TreatmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var req struct {
		SheepID uint `json:"sheepId"`
		dto.TreatmentDTO
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.Create(r.Context(), userID, req.SheepID, domain.Treatment{
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

func (h *TreatmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var req struct {
		SheepID uint `json:"sheepId"`
		Index   int  `json:"index"`
		dto.TreatmentDTO
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.Update(r.Context(), userID, req.SheepID, req.Index, domain.Treatment{
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

func (h *TreatmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var req struct {
		SheepID uint `json:"sheepId"`
		Index   int  `json:"index"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.Delete(r.Context(), userID, req.SheepID, req.Index)
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
