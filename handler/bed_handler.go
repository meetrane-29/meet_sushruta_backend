package handler

import (
	"strconv"
	"time"

	"meet_sushruta/config"
	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BedHandler handles bed-related API requests
type BedHandler struct {
	bedService service.BedService
}

// NewBedHandler creates a new bed handler
func NewBedHandler(bedService service.BedService) *BedHandler {
	return &BedHandler{
		bedService: bedService,
	}
}

// AdmitPatientRequest represents the admission request payload
type AdmitPatientRequest struct {
	BedID     string `json:"bed_id" binding:"required"`
	PatientID string `json:"patient_id" binding:"required"`
	Notes     string `json:"notes"`
}

// DischargePatientRequest represents the discharge request payload
type DischargePatientRequest struct {
	BedID          string `json:"bed_id" binding:"required"`
	DischargeNotes string `json:"discharge_notes"`
}

// BedResponse represents a bed in API responses
type BedResponse struct {
	ID                    string  `json:"id"`
	BedNumber             string  `json:"bed_number"`
	Ward                  string  `json:"ward"`
	Room                  string  `json:"room"`
	Floor                 int     `json:"floor"`
	BedType               string  `json:"bed_type"`
	Status                string  `json:"status"`
	PatientID             *string `json:"patient_id"`
	AdmittedAt            *string `json:"admitted_at"`
	DischargedAt          *string `json:"discharged_at"`
	Features              string  `json:"features"`
	DailyRate             float64 `json:"daily_rate"`
	IsMaintenanceRequired bool    `json:"is_maintenance_required"`
	MaintenanceNotes      string  `json:"maintenance_notes"`
	CreatedAt             int64   `json:"created_at"`
	UpdatedAt             int64   `json:"updated_at"`
}

// AdmitPatient admits a patient to a bed
// POST /api/v1/beds/admit
func (h *BedHandler) AdmitPatient(c *gin.Context) {
	var req AdmitPatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	// Parse UUIDs
	bedID, err := uuid.Parse(req.BedID)
	if err != nil {
		utils.Fail(c, 400, "invalid bed id")
		return
	}

	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		utils.Fail(c, 400, "invalid patient id")
		return
	}

	// Admit patient to bed
	if err := h.bedService.AdmitPatient(bedID, patientID, req.Notes); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message":     "patient admitted successfully",
		"bed_id":      bedID.String(),
		"patient_id":  patientID.String(),
		"admitted_at": time.Now().Format("2006-01-02 15:04"),
	})
}

// DischargePatient discharges a patient from a bed
// POST /api/v1/beds/discharge
func (h *BedHandler) DischargePatient(c *gin.Context) {
	var req DischargePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	// Parse bed ID
	bedID, err := uuid.Parse(req.BedID)
	if err != nil {
		utils.Fail(c, 400, "invalid bed id")
		return
	}

	// Discharge patient
	if err := h.bedService.DischargePatient(bedID, req.DischargeNotes); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message":       "patient discharged successfully",
		"bed_id":        bedID.String(),
		"discharged_at": time.Now().Format("2006-01-02 15:04"),
	})
}

// GetPatientBed retrieves the current bed for a patient
// GET /api/v1/patients/:patient_id/bed
func (h *BedHandler) GetPatientBed(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient id")
		return
	}

	bed, err := h.bedService.GetPatientBed(patientID)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	if bed == nil {
		utils.OK(c, gin.H{
			"message": "patient has no active bed",
			"bed":     nil,
		})
		return
	}

	response := convertBedToResponse(bed)
	utils.OK(c, response)
}

// GetAvailableBeds retrieves all available beds
// GET /api/v1/beds/available
// Query params:
//   - bed_type: optional filter by bed type
func (h *BedHandler) GetAvailableBeds(c *gin.Context) {
	bedType := c.Query("bed_type")

	beds, err := h.bedService.GetAvailableBeds(bedType)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	if beds == nil {
		beds = []model.Bed{}
	}

	responses := make([]BedResponse, 0)
	for i := range beds {
		responses = append(responses, *convertBedToResponse(&beds[i]))
	}

	utils.OK(c, gin.H{
		"total": len(responses),
		"beds":  responses,
	})
}

// GetBedsByWard retrieves all beds in a specific ward
// GET /api/v1/beds/ward/:ward
func (h *BedHandler) GetBedsByWard(c *gin.Context) {
	ward := c.Param("ward")

	beds, err := h.bedService.GetBedsByWard(ward)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	responses := make([]BedResponse, 0)
	for i := range beds {
		responses = append(responses, *convertBedToResponse(&beds[i]))
	}

	utils.OK(c, gin.H{
		"ward":  ward,
		"total": len(responses),
		"beds":  responses,
	})
}

// GetBedsByStatus retrieves all beds with a specific status
// GET /api/v1/beds/status/:status
func (h *BedHandler) GetBedsByStatus(c *gin.Context) {
	status := c.Param("status")

	beds, err := h.bedService.GetBedsByStatus(status)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	responses := make([]BedResponse, 0)
	for i := range beds {
		responses = append(responses, *convertBedToResponse(&beds[i]))
	}

	utils.OK(c, gin.H{
		"status": status,
		"total":  len(responses),
		"beds":   responses,
	})
}

// UpdateBedStatus updates the status of a bed
// PATCH /api/v1/beds/:id/status
type UpdateBedStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *BedHandler) UpdateBedStatus(c *gin.Context) {
	bedIDStr := c.Param("id")
	bedID, err := uuid.Parse(bedIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid bed id")
		return
	}

	var req UpdateBedStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	if err := h.bedService.UpdateBedStatus(bedID, req.Status); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message": "bed status updated successfully",
		"bed_id":  bedID.String(),
		"status":  req.Status,
	})
}

// GetWardStats retrieves statistics for a specific ward
// GET /api/v1/beds/ward/:ward/stats
func (h *BedHandler) GetWardStats(c *gin.Context) {
	ward := c.Param("ward")

	stats, err := h.bedService.GetWardStats(ward)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"ward":  ward,
		"stats": stats,
	})
}

// GetBedStats retrieves overall bed statistics
// GET /api/v1/beds/stats
func (h *BedHandler) GetBedStats(c *gin.Context) {
	stats, err := h.bedService.GetBedStats()
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"stats": stats,
	})
}

// GetAllBedsWithOccupancy retrieves all beds with occupancy information
// GET /api/v1/beds/all
// Query params:
//   - page: page number (default 1)
//   - limit: records per page (default 20)
func (h *BedHandler) GetAllBedsWithOccupancy(c *gin.Context) {
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

	offset := (page - 1) * limit

	// Get all beds from database
	var beds []model.Bed
	var total int64

	query := config.DB.Model(&model.Bed{})
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&beds).Error; err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	responses := make([]BedResponse, 0)
	for i := range beds {
		responses = append(responses, *convertBedToResponse(&beds[i]))
	}

	utils.OK(c, gin.H{
		"beds":  responses,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// Helper function to convert bed model to response
func convertBedToResponse(bed *model.Bed) *BedResponse {
	if bed == nil {
		return nil
	}

	var patientID *string
	if bed.PatientID != nil {
		patientIDStr := bed.PatientID.String()
		patientID = &patientIDStr
	}

	return &BedResponse{
		ID:                    bed.ID.String(),
		BedNumber:             bed.BedNumber,
		Ward:                  bed.Ward,
		Room:                  bed.Room,
		Floor:                 bed.Floor,
		BedType:               bed.BedType,
		Status:                bed.Status,
		PatientID:             patientID,
		AdmittedAt:            bed.AdmittedAt,
		DischargedAt:          bed.DischargedAt,
		Features:              bed.Features,
		DailyRate:             bed.DailyRate,
		IsMaintenanceRequired: bed.IsMaintenanceRequired,
		MaintenanceNotes:      bed.MaintenanceNotes,
		CreatedAt:             bed.CreatedAt,
		UpdatedAt:             bed.UpdatedAt,
	}
}
