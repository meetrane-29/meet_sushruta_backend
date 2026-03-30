package handler

import (
	"strconv"

	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PatientHandler struct {
	patientService     service.PatientService
	appointmentService service.AppointmentService
}

func NewPatientHandler(patientService service.PatientService, appointmentService service.AppointmentService) *PatientHandler {
	return &PatientHandler{
		patientService:     patientService,
		appointmentService: appointmentService,
	}
}

type CreatePatientRequest struct {
	UserID           uuid.UUID `json:"user_id" binding:"required"`
	DateOfBirth      string    `json:"date_of_birth" binding:"required"`
	Gender           string    `json:"gender" binding:"required"`
	BloodGroup       string    `json:"blood_group"`
	Address          string    `json:"address"`
	EmergencyContact string    `json:"emergency_contact"`
	MedicalHistory   string    `json:"medical_history"`
	Allergies        string    `json:"allergies"`
}

type UpdatePatientRequest struct {
	DateOfBirth      string `json:"date_of_birth"`
	Gender           string `json:"gender"`
	BloodGroup       string `json:"blood_group"`
	Address          string `json:"address"`
	EmergencyContact string `json:"emergency_contact"`
	MedicalHistory   string `json:"medical_history"`
	Allergies        string `json:"allergies"`
}

// CreatePatient creates a new patient
// POST /api/v1/patients
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	patient := &model.Patient{
		UserID:           req.UserID,
		DateOfBirth:      req.DateOfBirth,
		Gender:           req.Gender,
		BloodGroup:       req.BloodGroup,
		Address:          req.Address,
		EmergencyContact: req.EmergencyContact,
		MedicalHistory:   req.MedicalHistory,
		Allergies:        req.Allergies,
	}

	err := h.patientService.RegisterPatient(patient)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, patient)
}

// GetAllPatients retrieves all patients with pagination
// GET /api/v1/patients
func (h *PatientHandler) GetAllPatients(c *gin.Context) {
	page := 1
	limit := 10
	search := c.Query("search")

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

	patients, total, err := h.patientService.ListPatients(page, limit, search)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"patients": patients,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// GetPatient retrieves a patient by ID
// GET /api/v1/patients/:id
func (h *PatientHandler) GetPatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient id")
		return
	}

	// Try by patient ID first, fall back to user ID
	patient, err := h.patientService.GetPatient(id)
	if err != nil {
		patient, err = h.patientService.GetPatientByUserID(id)
		if err != nil {
			utils.Fail(c, 404, "patient not found")
			return
		}
	}

	utils.OK(c, patient)
}

// UpdatePatient updates a patient
// PATCH /api/v1/patients/:id
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient id")
		return
	}

	var req UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	// Get existing patient
	patient, err := h.patientService.GetPatient(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	// Update fields if provided
	if req.DateOfBirth != "" {
		patient.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != "" {
		patient.Gender = req.Gender
	}
	if req.BloodGroup != "" {
		patient.BloodGroup = req.BloodGroup
	}
	if req.Address != "" {
		patient.Address = req.Address
	}
	if req.EmergencyContact != "" {
		patient.EmergencyContact = req.EmergencyContact
	}
	if req.MedicalHistory != "" {
		patient.MedicalHistory = req.MedicalHistory
	}
	if req.Allergies != "" {
		patient.Allergies = req.Allergies
	}

	err = h.patientService.UpdatePatient(patient)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, patient)
}

// DeletePatient soft deletes a patient
// DELETE /api/v1/patients/:id
func (h *PatientHandler) DeletePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient id")
		return
	}

	err = h.patientService.SoftDeletePatient(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message": "patient deleted successfully",
	})
}

// GetPatientAppointments retrieves future appointments for a specific patient
// GET /api/v1/patients/:id/appointments
// The :id parameter can be either a PATIENT_ID or a USER_ID (for patient's own appointments)
func (h *PatientHandler) GetPatientAppointments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient id")
		return
	}

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

	// Try to get patient by USER_ID first (for case when frontend sends user ID)
	patientID := id
	patient, err := h.patientService.GetPatientByUserID(id)
	if err == nil && patient != nil {
		// Successfully found patient by USER_ID, use its ID
		patientID = patient.ID
	}
	// If not found by USER_ID, assume the provided ID is already a PATIENT_ID and proceed

	appointments, total, err := h.appointmentService.GetPatientAppointments(patientID, page, limit)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"appointments": appointments,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}
