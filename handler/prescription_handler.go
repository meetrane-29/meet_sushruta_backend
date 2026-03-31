package handler

import (
	"strconv"

	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PrescriptionHandler struct {
	prescriptionService service.PrescriptionService
	pharmacyService     service.PharmacyService
	labService          service.LabService
}

func NewPrescriptionHandler(
	prescriptionService service.PrescriptionService,
	pharmacyService service.PharmacyService,
	labService service.LabService,
) *PrescriptionHandler {
	return &PrescriptionHandler{
		prescriptionService: prescriptionService,
		pharmacyService:     pharmacyService,
		labService:          labService,
	}
}

type CreatePrescriptionRequest struct {
	AppointmentID uuid.UUID           `json:"appointment_id" binding:"required"`
	DoctorID      uuid.UUID           `json:"doctor_id" binding:"required"`
	Items         []service.ItemInput `json:"items" binding:"required,min=1"`
	Notes         string              `json:"notes"`
}

type UpdatePrescriptionStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type SendToPharmacyRequest struct {
	PrescriptionID uuid.UUID `json:"prescription_id" binding:"required"`
	Priority       string    `json:"priority"`
	Notes          string    `json:"notes"`
}

type SendToLabRequest struct {
	PatientID     uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID      uuid.UUID `json:"doctor_id" binding:"required"`
	TestType      string    `json:"test_type" binding:"required"`
	TestName      string    `json:"test_name" binding:"required"`
	ScheduledDate *string   `json:"scheduled_date"`
	Priority      string    `json:"priority"`
	Notes         string    `json:"notes"`
}

// CreatePrescription creates a new prescription with items
// POST /api/v1/prescriptions
func (h *PrescriptionHandler) CreatePrescription(c *gin.Context) {
	var req CreatePrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	prescription, err := h.prescriptionService.CreatePrescription(
		req.AppointmentID,
		req.DoctorID,
		req.Items,
		req.Notes,
	)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, prescription)
}

// GetPrescription retrieves a prescription by ID
// GET /api/v1/prescriptions/:id
func (h *PrescriptionHandler) GetPrescription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid prescription id")
		return
	}

	prescription, err := h.prescriptionService.GetPrescription(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, prescription)
}

// ListPrescriptions retrieves all prescriptions with pagination
// GET /api/v1/prescriptions
func (h *PrescriptionHandler) ListPrescriptions(c *gin.Context) {
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

	prescriptions, total, err := h.prescriptionService.ListPrescriptions(page, limit)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"prescriptions": prescriptions,
		"total":         total,
		"page":          page,
		"limit":         limit,
	})
}

// UpdatePrescriptionStatus updates prescription status
// PATCH /api/v1/prescriptions/:id/status
func (h *PrescriptionHandler) UpdatePrescriptionStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid prescription id")
		return
	}

	var req UpdatePrescriptionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	err = h.prescriptionService.UpdatePrescriptionStatus(id, req.Status)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	prescription, _ := h.prescriptionService.GetPrescription(id)
	utils.OK(c, prescription)
}

// GetPatientPrescriptions retrieves prescriptions for a patient
// GET /api/v1/prescriptions/patient/:patient_id
func (h *PrescriptionHandler) GetPatientPrescriptions(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient id")
		return
	}

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

	prescriptions, total, err := h.prescriptionService.GetPatientPrescriptions(patientID, page, limit)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"prescriptions": prescriptions,
		"total":         total,
		"page":          page,
		"limit":         limit,
	})
}

// GetPrescriptionMedicines retrieves all medicines in a prescription
// GET /api/v1/prescriptions/:id/medicines
func (h *PrescriptionHandler) GetPrescriptionMedicines(c *gin.Context) {
	prescriptionIDStr := c.Param("id")
	prescriptionID, err := uuid.Parse(prescriptionIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid prescription id")
		return
	}

	medicines, err := h.prescriptionService.GetPrescriptionMedicines(prescriptionID)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"prescription_id": prescriptionID,
		"medicines":       medicines,
		"total":           len(medicines),
	})
}

// SendToPharma sends a prescription to pharmacy
// POST /api/v1/prescriptions/:id/send-to-pharmacy
func (h *PrescriptionHandler) SendToPharmacy(c *gin.Context) {
	prescriptionIDStr := c.Param("id")
	prescriptionID, err := uuid.Parse(prescriptionIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid prescription id")
		return
	}

	var req SendToPharmacyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	// Get prescription details
	prescription, err := h.prescriptionService.GetPrescription(prescriptionID)
	if err != nil {
		utils.Fail(c, 404, "prescription not found")
		return
	}

	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	// Create pharmacy request
	pharmacyReq := &service.CreatePharmacyRequestInput{
		PatientID:      prescription.PatientID,
		DoctorID:       prescription.DoctorID,
		PrescriptionID: prescriptionID,
		Priority:       priority,
		Notes:          req.Notes,
	}

	result, err := h.pharmacyService.SendPrescriptionToPharmacy(pharmacyReq)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, result)
}

// SendToLab sends a lab request to lab
// POST /api/v1/prescriptions/:id/send-to-lab
func (h *PrescriptionHandler) SendToLab(c *gin.Context) {
	var req SendToLabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	// Create lab order
	labOrderInput := &service.LabOrderInput{
		PatientID:     req.PatientID,
		DoctorID:      req.DoctorID,
		TestType:      req.TestType,
		TestName:      req.TestName,
		ScheduledDate: req.ScheduledDate,
		Priority:      priority,
	}

	labRequest, err := h.labService.CreateOrder(labOrderInput)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, labRequest)
}
