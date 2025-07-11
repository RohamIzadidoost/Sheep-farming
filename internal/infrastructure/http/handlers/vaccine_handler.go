package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
func (h *VaccineHandler) CreateVaccine(c *gin.Context) {
	var req dto.CreateVaccineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidInput.Error()})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	vaccine := req.ToDomain(userID)
	if err := h.vaccineService.CreateVaccine(c.Request.Context(), vaccine); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ToVaccineResponse(vaccine)
	c.JSON(http.StatusCreated, resp)
}

// GetVaccineByID handles GET /vaccines/{id} requests.
func (h *VaccineHandler) GetVaccineByID(c *gin.Context) {
	vaccineID := c.Param("id")

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	vaccine, err := h.vaccineService.GetVaccineByID(c.Request.Context(), userID, vaccineID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ToVaccineResponse(vaccine)
	c.JSON(http.StatusOK, resp)
}

// GetAllVaccines handles GET /vaccines requests.
func (h *VaccineHandler) GetAllVaccines(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	vaccineList, err := h.vaccineService.GetAllVaccines(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.VaccineResponse
	for _, vaccine := range vaccineList {
		responses = append(responses, *dto.ToVaccineResponse(&vaccine))
	}

	c.JSON(http.StatusOK, responses)
}

// UpdateVaccine handles PUT /vaccines/{id} requests.
func (h *VaccineHandler) UpdateVaccine(c *gin.Context) {
	vaccineID := c.Param("id")

	var req dto.UpdateVaccineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidInput.Error()})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	existingVaccine, err := h.vaccineService.GetVaccineByID(c.Request.Context(), userID, vaccineID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.Name != nil {
		existingVaccine.Name = *req.Name
	}
	if req.IntervalMonths != nil {
		existingVaccine.IntervalMonths = *req.IntervalMonths
	}

	if err := h.vaccineService.UpdateVaccine(c.Request.Context(), existingVaccine); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ToVaccineResponse(existingVaccine)
	c.JSON(http.StatusOK, resp)
}

// DeleteVaccine handles DELETE /vaccines/{id} requests.
func (h *VaccineHandler) DeleteVaccine(c *gin.Context) {
	vaccineID := c.Param("id")

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.vaccineService.DeleteVaccine(c.Request.Context(), userID, vaccineID); err != nil {
		if err == domain.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
