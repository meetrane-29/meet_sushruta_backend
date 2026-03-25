package handler

import (
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PharmacyHandler struct {
	pharmacyService service.PharmacyService
}

func NewPharmacyHandler(pharmacyService service.PharmacyService) *PharmacyHandler {
	return &PharmacyHandler{
		pharmacyService: pharmacyService,
	}
}

type DispenseRequest struct {
	PrescriptionID uuid.UUID `json:"prescription_id" binding:"required"`
	StaffID        uuid.UUID `json:"staff_id" binding:"required"`
}

// Dispense dispenses medication from a prescription
// POST /api/v1/pharmacy/dispense
func (h *PharmacyHandler) Dispense(c *gin.Context) {
	var req DispenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	err := h.pharmacyService.Dispense(req.PrescriptionID, req.StaffID)
	if err != nil {
		// Check if it's an insufficient stock error
		if _, ok := err.(service.ErrInsufficientStock); ok {
			utils.Fail(c, 400, err.Error())
			return
		}
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message":         "medication dispensed successfully",
		"prescription_id": req.PrescriptionID,
	})
}

// GetDispenseHistory retrieves dispensing history for a prescription
// GET /api/v1/pharmacy/dispense/:prescription_id
func (h *PharmacyHandler) GetDispenseHistory(c *gin.Context) {
	prescriptionIDStr := c.Param("prescription_id")
	prescriptionID, err := uuid.Parse(prescriptionIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid prescription id")
		return
	}

	history, err := h.pharmacyService.GetDispenseHistory(prescriptionID)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, history)
}
