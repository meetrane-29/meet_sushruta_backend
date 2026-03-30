package handler

import (
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UHIDHandler struct {
	uhidService service.UHIDService
}

func NewUHIDHandler(uhidService service.UHIDService) *UHIDHandler {
	return &UHIDHandler{
		uhidService: uhidService,
	}
}

type GenerateUHIDRequest struct {
	PatientID uuid.UUID `json:"patient_id" binding:"required"`
}

type UHIDResponse struct {
	UHID string `json:"uhid"`
}

// GenerateUHID generates a new UHID for patient
// POST /api/v1/uhid/generate
func (h *UHIDHandler) GenerateUHID(c *gin.Context) {
	var req GenerateUHIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: patient_id is required")
		return
	}

	uhid, err := h.uhidService.GenerateUHID(req.PatientID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, UHIDResponse{UHID: uhid})
}

// GetUHIDByPatient retrieves UHID for a patient
// GET /api/v1/uhid/patient/:patient_id
func (h *UHIDHandler) GetUHIDByPatient(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient ID")
		return
	}

	uhid, err := h.uhidService.GetUHIDByPatientID(patientID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	if uhid == nil {
		utils.Fail(c, 404, "UHID not found for this patient")
		return
	}

	utils.OK(c, uhid)
}

// ValidateUHID validates if a UHID exists and is active
// GET /api/v1/uhid/validate/:uhid
func (h *UHIDHandler) ValidateUHID(c *gin.Context) {
	uhidStr := c.Param("uhid")
	if uhidStr == "" {
		utils.Fail(c, 400, "UHID parameter is required")
		return
	}

	uhid, err := h.uhidService.ValidateUHID(uhidStr)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, uhid)
}
