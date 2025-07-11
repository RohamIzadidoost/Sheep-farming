package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
func (h *ReminderHandler) GetReminders(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	reminders, err := h.reminderService.CalculateAndSendReminders(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}
