package handler

import (
	"strconv"

	"meet_sushruta/model"
	"meet_sushruta/repository"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdmissionHandler struct {
	admissionRepo repository.AdmissionRepository
}

func NewAdmissionHandler(admissionRepo repository.AdmissionRepository) *AdmissionHandler {
	return &AdmissionHandler{
		admissionRepo: admissionRepo,
	}
}

// AdmissionResponse represents an admission record in API responses
type AdmissionResponse struct {
	ID            string `json:"id"`
	PatientID     string `json:"patient_id"`
	PatientName   string `json:"patient_name"`
	DoctorID      string `json:"doctor_id"`
	DoctorName    string `json:"doctor_name"`
	BedID         string `json:"bed_id"`
	BedNumber     string `json:"bed_number"`
	Ward          string `json:"ward"`
	AdmissionDate string `json:"admission_date"`
	DischargeDate string `json:"discharge_date"`
	Reason        string `json:"reason"`
	Diagnosis     string `json:"diagnosis"`
	Status        string `json:"status"`
	RoomNumber    string `json:"room_number"`
	IsEmergency   bool   `json:"is_emergency"`
	Notes         string `json:"notes"`
}

// Helper function to convert admission model to response
func convertAdmissionToResponse(admission *model.AdmissionRecord) *AdmissionResponse {
	if admission == nil {
		return nil
	}

	patientName := "Unknown"
	if admission.Patient != nil && admission.Patient.User != nil {
		patientName = admission.Patient.User.FirstName + " " + admission.Patient.User.LastName
	}

	doctorName := "Unknown"
	if admission.Doctor != nil && admission.Doctor.User != nil {
		doctorName = admission.Doctor.User.FirstName + " " + admission.Doctor.User.LastName
	}

	bedNumber := ""
	if admission.Bed != nil {
		bedNumber = admission.Bed.BedNumber
	}

	bedID := ""
	if admission.BedID != uuid.Nil {
		bedID = admission.BedID.String()
	}

	dischargeDate := ""
	if admission.DischargeDate != nil {
		dischargeDate = *admission.DischargeDate
	}

	return &AdmissionResponse{
		ID:            admission.ID.String(),
		PatientID:     admission.PatientID.String(),
		PatientName:   patientName,
		DoctorID:      admission.DoctorID.String(),
		DoctorName:    doctorName,
		BedID:         bedID,
		BedNumber:     bedNumber,
		Ward:          admission.Ward,
		AdmissionDate: admission.AdmissionDate,
		DischargeDate: dischargeDate,
		Reason:        admission.Reason,
		Diagnosis:     admission.Diagnosis,
		Status:        string(admission.Status),
		RoomNumber:    admission.RoomNumber,
		IsEmergency:   admission.IsEmergency,
		Notes:         admission.Notes,
	}
}

// GetActiveAdmissions retrieves all active admission records
// GET /api/v1/admissions/active
// Query params:
//   - page: page number (default 1)
//   - limit: records per page (default 10, max 100)
func (h *AdmissionHandler) GetActiveAdmissions(c *gin.Context) {
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

	// Get all admissions and filter for active ones
	admissions, _, err := h.admissionRepo.GetAll(page, limit)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	// Filter to only active admissions
	activeAdmissions := make([]model.AdmissionRecord, 0)
	for _, admission := range admissions {
		if admission.Status == model.AdmissionActive {
			activeAdmissions = append(activeAdmissions, admission)
		}
	}

	responses := make([]AdmissionResponse, 0)
	for i := range activeAdmissions {
		responses = append(responses, *convertAdmissionToResponse(&activeAdmissions[i]))
	}

	utils.OK(c, gin.H{
		"admissions": responses,
		"total":      len(responses),
		"page":       page,
		"limit":      limit,
	})
}

// GetAllAdmissions retrieves all admission records (active and discharged)
// GET /api/v1/admissions
// Query params:
//   - page: page number (default 1)
//   - limit: records per page (default 10, max 100)
//   - status: filter by status (active, discharged, cancelled)
func (h *AdmissionHandler) GetAllAdmissions(c *gin.Context) {
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

	admissions, total, err := h.admissionRepo.GetAll(page, limit)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	responses := make([]AdmissionResponse, 0)
	for i := range admissions {
		responses = append(responses, *convertAdmissionToResponse(&admissions[i]))
	}

	utils.OK(c, gin.H{
		"admissions": responses,
		"total":      total,
		"page":       page,
		"limit":      limit,
	})
}

// GetAdmissionByID retrieves a specific admission record
// GET /api/v1/admissions/:id
func (h *AdmissionHandler) GetAdmissionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid admission id")
		return
	}

	admission, err := h.admissionRepo.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, "admission not found")
		return
	}

	utils.OK(c, convertAdmissionToResponse(admission))
}

// CreateAdmission creates a new admission record
// POST /api/v1/admissions
type CreateAdmissionRequest struct {
	PatientID     uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID      uuid.UUID `json:"doctor_id" binding:"required"`
	BedID         uuid.UUID `json:"bed_id"`
	AdmissionDate string    `json:"admission_date" binding:"required"`
	Reason        string    `json:"reason" binding:"required"`
	Diagnosis     string    `json:"diagnosis"`
	Ward          string    `json:"ward"`
	RoomNumber    string    `json:"room_number"`
	IsEmergency   bool      `json:"is_emergency"`
	Notes         string    `json:"notes"`
}

func (h *AdmissionHandler) CreateAdmission(c *gin.Context) {
	var req CreateAdmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	admission := &model.AdmissionRecord{
		PatientID:     req.PatientID,
		DoctorID:      req.DoctorID,
		BedID:         req.BedID,
		AdmissionDate: req.AdmissionDate,
		Reason:        req.Reason,
		Diagnosis:     req.Diagnosis,
		Status:        model.AdmissionActive,
		Ward:          req.Ward,
		RoomNumber:    req.RoomNumber,
		IsEmergency:   req.IsEmergency,
		Notes:         req.Notes,
	}

	if err := h.admissionRepo.Create(admission); err != nil {
		utils.Fail(c, 400, "failed to create admission: "+err.Error())
		return
	}

	utils.OK(c, convertAdmissionToResponse(admission))
}

// UpdateAdmission updates an admission record
// PATCH /api/v1/admissions/:id
type UpdateAdmissionRequest struct {
	Diagnosis     string  `json:"diagnosis"`
	Status        string  `json:"status"`
	DischargeDate *string `json:"discharge_date"`
	Notes         string  `json:"notes"`
}

func (h *AdmissionHandler) UpdateAdmission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid admission id")
		return
	}

	var req UpdateAdmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	admission, err := h.admissionRepo.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, "admission not found")
		return
	}

	if req.Diagnosis != "" {
		admission.Diagnosis = req.Diagnosis
	}
	if req.Status != "" {
		admission.Status = model.AdmissionStatus(req.Status)
	}
	if req.DischargeDate != nil {
		admission.DischargeDate = req.DischargeDate
	}
	if req.Notes != "" {
		admission.Notes = req.Notes
	}

	if err := h.admissionRepo.Update(admission); err != nil {
		utils.Fail(c, 400, "failed to update admission: "+err.Error())
		return
	}

	utils.OK(c, convertAdmissionToResponse(admission))
}

// DeleteAdmission deletes an admission record
// DELETE /api/v1/admissions/:id
func (h *AdmissionHandler) DeleteAdmission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid admission id")
		return
	}

	if err := h.admissionRepo.SoftDelete(id); err != nil {
		utils.Fail(c, 400, "failed to delete admission: "+err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message": "admission deleted successfully",
	})
}

// GetPatientAdmissions retrieves all admissions for a specific patient
// GET /api/v1/admissions/patient/:patient_id
func (h *AdmissionHandler) GetPatientAdmissions(c *gin.Context) {
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

	admissions, total, err := h.admissionRepo.GetByPatientID(patientID, page, limit)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch admissions: "+err.Error())
		return
	}

	responses := make([]AdmissionResponse, 0)
	for i := range admissions {
		responses = append(responses, *convertAdmissionToResponse(&admissions[i]))
	}

	utils.OK(c, gin.H{
		"admissions": responses,
		"total":      total,
		"page":       page,
		"limit":      limit,
	})
}
