package handler

import (
	"strconv"

	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VitalsHandler struct {
	vitalsService service.VitalsService
}

func NewVitalsHandler(vitalsService service.VitalsService) *VitalsHandler {
	return &VitalsHandler{
		vitalsService: vitalsService,
	}
}

type CreateVitalsRequest struct {
	PatientID       uuid.UUID `json:"patient_id" binding:"required"`
	Temperature     float64   `json:"temperature" binding:"required"`
	BloodPressure   string    `json:"blood_pressure" binding:"required"`
	HeartRate       int       `json:"heart_rate" binding:"required"`
	RespiratoryRate int       `json:"respiratory_rate" binding:"required"`
	Weight          float64   `json:"weight" binding:"required"`
	Height          float64   `json:"height"`
	BloodSugar      float64   `json:"blood_sugar"`
	Oxygen          int       `json:"oxygen"` // SpO2
	RecordedAt      string    `json:"recorded_at" binding:"required"`
	RecordedBy      uuid.UUID `json:"recorded_by" binding:"required"`
}

// CreateVitals creates a new vitals record
// POST /api/v1/vitals
func (h *VitalsHandler) CreateVitals(c *gin.Context) {
	var req CreateVitalsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	vitals := &model.Vitals{
		PatientID:       req.PatientID,
		Temperature:     req.Temperature,
		BloodPressure:   req.BloodPressure,
		HeartRate:       req.HeartRate,
		RespiratoryRate: req.RespiratoryRate,
		Weight:          req.Weight,
		Height:          req.Height,
		BloodSugar:      req.BloodSugar,
		Oxygen:          req.Oxygen,
		RecordedAt:      req.RecordedAt,
		RecordedBy:      req.RecordedBy,
	}

	err := h.vitalsService.CreateVitals(vitals)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, vitals)
}

// GetVitalsByID retrieves vitals by ID
// GET /api/v1/vitals/:id
func (h *VitalsHandler) GetVitalsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid vitals id")
		return
	}

	vitals, err := h.vitalsService.GetVitalsByID(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, vitals)
}

// GetPatientVitals retrieves vitals for a patient
// GET /api/v1/patients/:patient_id/vitals
// Query params:
//   - limit: number of records to return (default 10, max 100)
//   - offset: number of records to skip (default 0)
func (h *VitalsHandler) GetPatientVitals(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient id")
		return
	}

	// Check if requesting latest only
	if c.Query("latest") == "true" {
		vitals, err := h.vitalsService.GetLatestVitalsByPatientID(patientID)
		if err != nil {
			utils.Fail(c, 500, err.Error())
			return
		}
		utils.OK(c, vitals)
		return
	}

	// Parse pagination params
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	vitals, err := h.vitalsService.GetVitalsHistoryByPatientID(patientID, limit, offset)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	if vitals == nil {
		vitals = []model.Vitals{}
	}

	utils.OK(c, vitals)
}

// UpdateVitals updates an existing vitals record
// PATCH /api/v1/vitals/:id
func (h *VitalsHandler) UpdateVitals(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid vitals id")
		return
	}

	var req CreateVitalsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	vitals := &model.Vitals{
		ID:              id,
		Temperature:     req.Temperature,
		BloodPressure:   req.BloodPressure,
		HeartRate:       req.HeartRate,
		RespiratoryRate: req.RespiratoryRate,
		Weight:          req.Weight,
		Height:          req.Height,
		BloodSugar:      req.BloodSugar,
		Oxygen:          req.Oxygen,
		RecordedAt:      req.RecordedAt,
	}

	err = h.vitalsService.UpdateVitals(vitals)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, vitals)
}

// DeleteVitals deletes a vitals record
// DELETE /api/v1/vitals/:id
func (h *VitalsHandler) DeleteVitals(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid vitals id")
		return
	}

	err = h.vitalsService.DeleteVitals(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, gin.H{"message": "vitals record deleted successfully"})
}
