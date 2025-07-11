package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"sheep_farm_backend_go/internal/application/services"
	"sheep_farm_backend_go/internal/domain"
	"sheep_farm_backend_go/internal/infrastructure/http/dto"
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
func (h *SheepHandler) CreateSheep(c *gin.Context) {
	var req dto.CreateSheepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidInput.Error()})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	sheep := req.ToDomain(userID)
	if err := h.sheepService.CreateSheep(c.Request.Context(), sheep); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ToSheepResponse(sheep)
	c.JSON(http.StatusCreated, resp)
}

// GetSheepByID handles GET /sheep/{id} requests.
func (h *SheepHandler) GetSheepByID(c *gin.Context) {
	sheepID := c.Param("id")

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	sheep, err := h.sheepService.GetSheepByID(c.Request.Context(), userID, sheepID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ToSheepResponse(sheep)
	c.JSON(http.StatusOK, resp)
}

// GetAllSheep handles GET /sheep requests.
func (h *SheepHandler) GetAllSheep(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	sheepList, err := h.sheepService.GetAllSheep(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.SheepResponse
	for _, sheep := range sheepList {
		responses = append(responses, *dto.ToSheepResponse(&sheep))
	}

	c.JSON(http.StatusOK, responses)
}

// UpdateSheep handles PUT /sheep/{id} requests.
func (h *SheepHandler) UpdateSheep(c *gin.Context) {
	sheepID := c.Param("id")

	var req dto.UpdateSheepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidInput.Error()})
		return
	}

	// Get existing sheep to apply updates
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	existingSheep, err := h.sheepService.GetSheepByID(c.Request.Context(), userID, sheepID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	if err := h.sheepService.UpdateSheep(c.Request.Context(), existingSheep); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ToSheepResponse(existingSheep)
	c.JSON(http.StatusOK, resp)
}

// DeleteSheep handles DELETE /sheep/{id} requests.
func (h *SheepHandler) DeleteSheep(c *gin.Context) {
	sheepID := c.Param("id")

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.sheepService.DeleteSheep(c.Request.Context(), userID, sheepID); err != nil {
		if err == domain.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
