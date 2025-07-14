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

// LambingHandler manages lambing endpoints
type LambingHandler struct {
	service *services.LambingService
}

func NewLambingHandler(s *services.LambingService) *LambingHandler {
	return &LambingHandler{service: s}
}

func (h *LambingHandler) List(w http.ResponseWriter, r *http.Request) {
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
	var resp []dto.LambingDTO
	for _, l := range list {
		resp = append(resp, dto.LambingDTO{
			Date:    dto.DateOnly(l.Date),
			NumBorn: l.NumBorn,
			Sexes:   l.Sexes,
			NumDead: l.NumDead,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *LambingHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var req struct {
		SheepID string `json:"sheepId"`
		dto.LambingDTO
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.Create(r.Context(), userID, req.SheepID, domain.Lambing{
		Date:    time.Time(req.LambingDTO.Date),
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

func (h *LambingHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var req struct {
		SheepID string `json:"sheepId"`
		Index   int    `json:"index"`
		dto.LambingDTO
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, domain.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.Update(r.Context(), userID, req.SheepID, req.Index, domain.Lambing{
		Date:    time.Time(req.LambingDTO.Date),
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

func (h *LambingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var req struct {
		SheepID string `json:"sheepId"`
		Index   int    `json:"index"`
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
