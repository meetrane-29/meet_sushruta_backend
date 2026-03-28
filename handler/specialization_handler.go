package handler

import (
	"strconv"

	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SpecializationHandler struct {
	specializationService service.SpecializationService
}

func NewSpecializationHandler(specializationService service.SpecializationService) *SpecializationHandler {
	return &SpecializationHandler{
		specializationService: specializationService,
	}
}

type CreateSpecializationRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	Icon        string     `json:"icon"`
	ParentID    *uuid.UUID `json:"parent_id"`
}

type UpdateSpecializationRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	Icon        string     `json:"icon"`
	ParentID    *uuid.UUID `json:"parent_id"`
	IsActive    *bool      `json:"is_active"`
}

// CreateSpecialization creates a new specialization
// POST /api/v1/specializations
func (h *SpecializationHandler) CreateSpecialization(c *gin.Context) {
	var req CreateSpecializationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	specialization := &model.Specialization{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Icon:        req.Icon,
		ParentID:    req.ParentID,
		IsActive:    true,
	}

	err := h.specializationService.CreateSpecialization(specialization)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, specialization)
}

// GetAllSpecializations retrieves all specializations with pagination
// GET /api/v1/specializations
func (h *SpecializationHandler) GetAllSpecializations(c *gin.Context) {
	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	specializations, total, err := h.specializationService.ListSpecializations(page, limit)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"specializations": specializations,
		"total":           total,
		"page":            page,
		"limit":           limit,
	})
}

// GetMainSpecializations retrieves all main specializations (without parent)
// GET /api/v1/specializations/main
func (h *SpecializationHandler) GetMainSpecializations(c *gin.Context) {
	page := 1
	limit := 50

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	specializations, total, err := h.specializationService.ListMainSpecializations(page, limit)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"specializations": specializations,
		"total":           total,
		"page":            page,
		"limit":           limit,
	})
}

// GetSubSpecializations retrieves sub-specializations for a parent
// GET /api/v1/specializations/:id/sub
func (h *SpecializationHandler) GetSubSpecializations(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid specialization id")
		return
	}

	subSpecializations, err := h.specializationService.GetSubSpecializations(id)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"specializations": subSpecializations,
	})
}

// GetSpecialization retrieves a specialization by ID
// GET /api/v1/specializations/:id
func (h *SpecializationHandler) GetSpecialization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid specialization id")
		return
	}

	specialization, err := h.specializationService.GetSpecialization(id)
	if err != nil {
		utils.Fail(c, 404, "specialization not found")
		return
	}

	utils.OK(c, specialization)
}

// UpdateSpecialization updates a specialization
// PATCH /api/v1/specializations/:id
func (h *SpecializationHandler) UpdateSpecialization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid specialization id")
		return
	}

	var req UpdateSpecializationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	specialization, err := h.specializationService.GetSpecialization(id)
	if err != nil {
		utils.Fail(c, 404, "specialization not found")
		return
	}

	if req.Name != "" {
		specialization.Name = req.Name
	}
	if req.Description != "" {
		specialization.Description = req.Description
	}
	if req.Category != "" {
		specialization.Category = req.Category
	}
	if req.Icon != "" {
		specialization.Icon = req.Icon
	}
	if req.IsActive != nil {
		specialization.IsActive = *req.IsActive
	}

	err = h.specializationService.UpdateSpecialization(specialization)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, specialization)
}

// DeleteSpecialization deletes a specialization
// DELETE /api/v1/specializations/:id
func (h *SpecializationHandler) DeleteSpecialization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid specialization id")
		return
	}

	err = h.specializationService.DeleteSpecialization(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message": "specialization deleted successfully",
	})
}
