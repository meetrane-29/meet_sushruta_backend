package handler

import (
	"strconv"

	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HospitalHandler struct {
	hospitalService service.HospitalService
}

func NewHospitalHandler(hospitalService service.HospitalService) *HospitalHandler {
	return &HospitalHandler{
		hospitalService: hospitalService,
	}
}

type CreateHospitalRequest struct {
	Name         string  `json:"name" binding:"required"`
	Address      string  `json:"address" binding:"required"`
	City         string  `json:"city" binding:"required"`
	State        string  `json:"state" binding:"required"`
	PostalCode   string  `json:"postal_code"`
	Phone        string  `json:"phone" binding:"required"`
	Email        string  `json:"email" binding:"required"`
	Website      string  `json:"website"`
	Description  string  `json:"description"`
	TotalBeds    int     `json:"total_beds"`
	AvailableBeds int    `json:"available_beds"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

type UpdateHospitalRequest struct {
	Name         string  `json:"name"`
	Address      string  `json:"address"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	PostalCode   string  `json:"postal_code"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email"`
	Website      string  `json:"website"`
	Description  string  `json:"description"`
	TotalBeds    *int    `json:"total_beds"`
	AvailableBeds *int   `json:"available_beds"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	IsVerified   *bool   `json:"is_verified"`
	IsActive     *bool   `json:"is_active"`
}

// CreateHospital creates a new hospital
// POST /api/v1/hospitals
func (h *HospitalHandler) CreateHospital(c *gin.Context) {
	var req CreateHospitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	hospital := &model.Hospital{
		Name:          req.Name,
		Address:       req.Address,
		City:          req.City,
		State:         req.State,
		PostalCode:    req.PostalCode,
		Phone:         req.Phone,
		Email:         req.Email,
		Website:       req.Website,
		Description:   req.Description,
		TotalBeds:     req.TotalBeds,
		AvailableBeds: req.AvailableBeds,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		IsActive:      true,
	}

	err := h.hospitalService.CreateHospital(hospital)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, hospital)
}

// GetAllHospitals retrieves all hospitals with pagination
// GET /api/v1/hospitals
func (h *HospitalHandler) GetAllHospitals(c *gin.Context) {
	page := 1
	limit := 10

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

	hospitals, total, err := h.hospitalService.ListHospitals(page, limit)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"hospitals": hospitals,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

// GetHospitalsByCity retrieves hospitals in a specific city
// GET /api/v1/hospitals/city/:city
func (h *HospitalHandler) GetHospitalsByCity(c *gin.Context) {
	city := c.Param("city")
	page := 1
	limit := 10

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

	hospitals, total, err := h.hospitalService.ListHospitalsByCity(city, page, limit)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"hospitals": hospitals,
		"total":     total,
		"city":      city,
		"page":      page,
		"limit":     limit,
	})
}

// SearchHospitals searches hospitals by name, city, or address
// GET /api/v1/hospitals/search?q=query
func (h *HospitalHandler) SearchHospitals(c *gin.Context) {
	query := c.Query("q")
	page := 1
	limit := 10

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

	hospitals, total, err := h.hospitalService.SearchHospitals(query, page, limit)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"hospitals": hospitals,
		"total":     total,
		"query":     query,
		"page":      page,
		"limit":     limit,
	})
}

// GetHospital retrieves a hospital by ID
// GET /api/v1/hospitals/:id
func (h *HospitalHandler) GetHospital(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid hospital id")
		return
	}

	hospital, err := h.hospitalService.GetHospital(id)
	if err != nil {
		utils.Fail(c, 404, "hospital not found")
		return
	}

	utils.OK(c, hospital)
}

// UpdateHospital updates a hospital
// PATCH /api/v1/hospitals/:id
func (h *HospitalHandler) UpdateHospital(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid hospital id")
		return
	}

	var req UpdateHospitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	hospital, err := h.hospitalService.GetHospital(id)
	if err != nil {
		utils.Fail(c, 404, "hospital not found")
		return
	}

	// Update fields if provided
	if req.Name != "" {
		hospital.Name = req.Name
	}
	if req.Address != "" {
		hospital.Address = req.Address
	}
	if req.City != "" {
		hospital.City = req.City
	}
	if req.State != "" {
		hospital.State = req.State
	}
	if req.PostalCode != "" {
		hospital.PostalCode = req.PostalCode
	}
	if req.Phone != "" {
		hospital.Phone = req.Phone
	}
	if req.Email != "" {
		hospital.Email = req.Email
	}
	if req.Website != "" {
		hospital.Website = req.Website
	}
	if req.Description != "" {
		hospital.Description = req.Description
	}
	if req.TotalBeds != nil {
		hospital.TotalBeds = *req.TotalBeds
	}
	if req.AvailableBeds != nil {
		hospital.AvailableBeds = *req.AvailableBeds
	}
	if req.Latitude != nil {
		hospital.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		hospital.Longitude = *req.Longitude
	}
	if req.IsVerified != nil {
		hospital.IsVerified = *req.IsVerified
	}
	if req.IsActive != nil {
		hospital.IsActive = *req.IsActive
	}

	err = h.hospitalService.UpdateHospital(hospital)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, hospital)
}

// DeleteHospital deletes a hospital
// DELETE /api/v1/hospitals/:id
func (h *HospitalHandler) DeleteHospital(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid hospital id")
		return
	}

	err = h.hospitalService.DeleteHospital(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message": "hospital deleted successfully",
	})
}
