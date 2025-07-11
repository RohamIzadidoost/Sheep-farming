package handlers

import (
	"encoding/json"
	"net/http"

	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/infrastructure/http/middleware"
)

// ReminderHandler provides reminder-related endpoints.
type ReminderHandler struct {
	reminderService *services.ReminderService
}

// NewReminderHandler creates a ReminderHandler.
func NewReminderHandler(reminderService *services.ReminderService) *ReminderHandler {
	return &ReminderHandler{reminderService: reminderService}
}

// GetReminders calculates reminders for the authenticated user.
func (h *ReminderHandler) GetReminders(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	reminders, err := h.reminderService.CalculateAndSendReminders(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}
