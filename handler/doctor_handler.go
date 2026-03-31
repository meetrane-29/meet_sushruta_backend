package handler

import (
	"strconv"
	"time"

	"meet_sushruta/config"
	"meet_sushruta/model"
	"meet_sushruta/repository"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DoctorHandler struct {
	doctorService              service.DoctorService
	admissionRepository        repository.AdmissionRepository
	progressNoteRepository     repository.ProgressNoteRepository
	nurseInstructionRepository repository.NurseInstructionRepository
	dischargeSummaryRepository repository.DischargeSummaryRepository
}

func NewDoctorHandler(
	doctorService service.DoctorService,
	admissionRepo repository.AdmissionRepository,
	progressNoteRepo repository.ProgressNoteRepository,
	nurseInstructionRepo repository.NurseInstructionRepository,
	dischargeSummaryRepo repository.DischargeSummaryRepository,
) *DoctorHandler {
	return &DoctorHandler{
		doctorService:              doctorService,
		admissionRepository:        admissionRepo,
		progressNoteRepository:     progressNoteRepo,
		nurseInstructionRepository: nurseInstructionRepo,
		dischargeSummaryRepository: dischargeSummaryRepo,
	}
}

type CreateDoctorRequest struct {
	UserID               uuid.UUID `json:"user_id" binding:"required"`
	Specialization       string    `json:"specialization" binding:"required"`
	LicenseNumber        string    `json:"license_number" binding:"required"`
	CertificationURL     string    `json:"certification_url"`
	Bio                  string    `json:"bio"`
	Department           string    `json:"department"`
	ConsultationFee      float64   `json:"consultation_fee"`
	JoiningDate          int64     `json:"joining_date"`
	Salary               float64   `json:"salary"`
	AttendancePercentage float64   `json:"attendance_percentage"`
	LeaveBalance         int       `json:"leave_balance"`
}

type UpdateDoctorRequest struct {
	Specialization       string  `json:"specialization"`
	Bio                  string  `json:"bio"`
	Department           string  `json:"department"`
	ConsultationFee      float64 `json:"consultation_fee"`
	JoiningDate          int64   `json:"joining_date"`
	Salary               float64 `json:"salary"`
	AttendancePercentage float64 `json:"attendance_percentage"`
	LeaveBalance         int     `json:"leave_balance"`
}

// CreateDoctor creates a new doctor
// POST /api/v1/doctors
func (h *DoctorHandler) CreateDoctor(c *gin.Context) {
	var req CreateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	doctor := &model.Doctor{
		UserID:               req.UserID,
		Specialization:       req.Specialization,
		LicenseNumber:        req.LicenseNumber,
		CertificationURL:     req.CertificationURL,
		Bio:                  req.Bio,
		Department:           req.Department,
		ConsultationFee:      req.ConsultationFee,
		JoiningDate:          req.JoiningDate,
		Salary:               req.Salary,
		AttendancePercentage: req.AttendancePercentage,
		LeaveBalance:         req.LeaveBalance,
	}

	err := h.doctorService.CreateDoctor(doctor)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, doctor)
}

// GetAllDoctors retrieves all doctors with pagination
// GET /api/v1/doctors
// Supports search by name, specialization, or slug
func (h *DoctorHandler) GetAllDoctors(c *gin.Context) {
	page := 1
	limit := 100
	search := c.Query("search")
	slug := c.Query("slug")

	// If slug is provided, use it as search term
	if slug != "" {
		search = slug
	}

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

	doctors, _, err := h.doctorService.ListDoctors(page, limit, search)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, doctors)
}

// GetDoctor retrieves a doctor by ID
// GET /api/v1/doctors/:id
func (h *DoctorHandler) GetDoctor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor id")
		return
	}

	doctor, err := h.doctorService.GetDoctor(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, doctor)
}

// UpdateDoctor updates a doctor
// PATCH /api/v1/doctors/:id
func (h *DoctorHandler) UpdateDoctor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor id")
		return
	}

	var req UpdateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	// Get existing doctor
	doctor, err := h.doctorService.GetDoctor(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	// Update fields if provided
	if req.Specialization != "" {
		doctor.Specialization = req.Specialization
	}
	if req.Bio != "" {
		doctor.Bio = req.Bio
	}
	if req.Department != "" {
		doctor.Department = req.Department
	}
	if req.ConsultationFee > 0 {
		doctor.ConsultationFee = req.ConsultationFee
	}
	if req.JoiningDate > 0 {
		doctor.JoiningDate = req.JoiningDate
	}
	if req.Salary > 0 {
		doctor.Salary = req.Salary
	}
	if req.AttendancePercentage >= 0 {
		doctor.AttendancePercentage = req.AttendancePercentage
	}
	if req.LeaveBalance >= 0 {
		doctor.LeaveBalance = req.LeaveBalance
	}

	err = h.doctorService.UpdateDoctor(doctor)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, doctor)
}

// DeleteDoctor deletes a doctor (soft delete)
// DELETE /api/v1/doctors/:id
func (h *DoctorHandler) DeleteDoctor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor id")
		return
	}

	err = h.doctorService.SoftDeleteDoctor(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message": "doctor deleted successfully",
	})
}

// GetAvailableSlots retrieves available appointment slots for a doctor
// GET /api/v1/doctors/:id/slots?date=2026-03-25
func (h *DoctorHandler) GetAvailableSlots(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor id")
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		utils.Fail(c, 400, "date parameter is required (format: YYYY-MM-DD)")
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.Fail(c, 400, "invalid date format (use YYYY-MM-DD)")
		return
	}

	slots, err := h.doctorService.GetAvailableSlots(id, date)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"doctor_id": id.String(),
		"date":      dateStr,
		"slots":     slots,
	})
}

// GetMe retrieves the current logged-in doctor
// GET /api/v1/doctors/me
func (h *DoctorHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Fail(c, 401, "unauthorized")
		return
	}

	doctor, err := h.doctorService.GetDoctorByUserID(userID.(uuid.UUID))
	if err != nil {
		utils.Fail(c, 404, "doctor not found")
		return
	}

	utils.OK(c, doctor)
}

// CreateIPDPatients - kept for compatibility
// Use admission handler for creating admissions instead
// Deprecated - use POST /api/v1/admissions instead

// GetIPDPatients retrieves all admitted patients for a doctor
// GET /api/v1/doctors/ipd-patients
func (h *DoctorHandler) GetIPDPatients(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, 401, "unauthorized")
		return
	}

	// Get doctor ID from user ID
	doctor, err := h.doctorService.GetDoctorByUserID(userID.(uuid.UUID))
	if err != nil {
		utils.Fail(c, 404, "doctor not found")
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

	admissions, total, err := h.admissionRepository.GetActiveByDoctorID(doctor.ID, page, limit)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch admissions: "+err.Error())
		return
	}

	utils.OK(c, gin.H{
		"data":  admissions,
		"total": total,
	})
}

// CreateProgressNote creates a progress note for admitted patient
// POST /api/v1/progress-notes
type CreateProgressNoteRequest struct {
	AdmissionID    *uuid.UUID `json:"admission_id"` // optional — nil for OPD consultations
	PatientID      uuid.UUID  `json:"patient_id" binding:"required"`
	UserID         uuid.UUID  `json:"user_id" binding:"required"` // doctor's user_id (will lookup actual doctor_id)
	RecordedDate   *string    `json:"recorded_date"`              // optional — defaults to today
	Subjective     string     `json:"subjective"`
	Objective      string     `json:"objective"`
	Assessment     string     `json:"assessment"`
	Plan           string     `json:"plan"`
	Vitals         string     `json:"vitals"`
	Medications    string     `json:"medications"`
	NextReviewDate *string    `json:"next_review_date"`
	Notes          string     `json:"notes"`
}

func (h *DoctorHandler) CreateProgressNote(c *gin.Context) {
	var req CreateProgressNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	// Look up doctor record using user_id to get the actual doctor_id
	doctor, err := h.doctorService.GetDoctorByUserID(req.UserID)
	if err != nil {
		utils.Fail(c, 400, "doctor not found: ensure you are logged in as a doctor")
		return
	}

	// Set default recorded_date to today if not provided
	var recordedDate *string
	if req.RecordedDate != nil {
		recordedDate = req.RecordedDate
	} else {
		today := time.Now().Format("2006-01-02")
		recordedDate = &today
	}

	note := &model.ProgressNote{
		AdmissionID:    req.AdmissionID,
		PatientID:      req.PatientID,
		DoctorID:       doctor.ID, // Use the actual doctor_id from doctors table
		RecordedDate:   recordedDate,
		Subjective:     req.Subjective,
		Objective:      req.Objective,
		Assessment:     req.Assessment,
		Plan:           req.Plan,
		Vitals:         req.Vitals,
		Medications:    req.Medications,
		NextReviewDate: req.NextReviewDate,
		Notes:          req.Notes,
	}

	if err := h.progressNoteRepository.Create(note); err != nil {
		utils.Fail(c, 400, "failed to create progress note: "+err.Error())
		return
	}

	utils.OK(c, note)
}

// CreateNurseInstruction creates nurse instructions for a patient
// POST /api/v1/nurse-instructions
type CreateNurseInstructionRequest struct {
	AdmissionID         uuid.UUID `json:"admission_id" binding:"required"`
	PatientID           uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID            uuid.UUID `json:"doctor_id" binding:"required"`
	IssuedDate          string    `json:"issued_date" binding:"required"` // YYYY-MM-DD
	VitalsFrequency     string    `json:"vitals_frequency"`
	VitalsParameters    string    `json:"vitals_parameters"`
	DietType            string    `json:"diet_type"`
	DietaryRestrictions string    `json:"dietary_restrictions"`
	FluidsRestriction   string    `json:"fluids_restriction"`
	MedicationNotes     string    `json:"medication_notes"`
	DrugAllergies       string    `json:"drug_allergies"`
	ActivityLevel       string    `json:"activity_level"`
	ActivityNotes       string    `json:"activity_notes"`
	SpecialMonitoring   string    `json:"special_monitoring"`
	CautionPoints       string    `json:"caution_points"`
	HygieneInstructions string    `json:"hygiene_instructions"`
	OtherInstructions   string    `json:"other_instructions"`
}

func (h *DoctorHandler) CreateNurseInstruction(c *gin.Context) {
	var req CreateNurseInstructionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	instruction := &model.NurseInstruction{
		AdmissionID:         req.AdmissionID,
		PatientID:           req.PatientID,
		DoctorID:            req.DoctorID,
		IssuedDate:          req.IssuedDate,
		VitalsFrequency:     req.VitalsFrequency,
		VitalsParameters:    req.VitalsParameters,
		DietType:            req.DietType,
		DietaryRestrictions: req.DietaryRestrictions,
		FluidsRestriction:   req.FluidsRestriction,
		MedicationNotes:     req.MedicationNotes,
		DrugAllergies:       req.DrugAllergies,
		ActivityLevel:       req.ActivityLevel,
		ActivityNotes:       req.ActivityNotes,
		SpecialMonitoring:   req.SpecialMonitoring,
		CautionPoints:       req.CautionPoints,
		HygieneInstructions: req.HygieneInstructions,
		OtherInstructions:   req.OtherInstructions,
		Status:              "active",
	}

	if err := h.nurseInstructionRepository.Create(instruction); err != nil {
		utils.Fail(c, 400, "failed to create nurse instruction: "+err.Error())
		return
	}

	utils.OK(c, instruction)
}

// CreateDischargeSummary creates a discharge summary
// POST /api/v1/discharge-summaries
type CreateDischargeSummaryRequest struct {
	AdmissionID            uuid.UUID `json:"admission_id" binding:"required"`
	PatientID              uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID               uuid.UUID `json:"doctor_id" binding:"required"`
	DischargeDate          string    `json:"discharge_date" binding:"required"` // YYYY-MM-DD HH:MM
	FinalDiagnosis         string    `json:"final_diagnosis" binding:"required"`
	ProceduresPerformed    string    `json:"procedures_performed"`
	ComplicationsIfAny     string    `json:"complications_if_any"`
	DischargeMedications   string    `json:"discharge_medications"`
	MedicationInstructions string    `json:"medication_instructions"`
	DietRecommendation     string    `json:"diet_recommendation"`
	ActivityRecommendation string    `json:"activity_recommendation"`
	FollowUpInstructions   string    `json:"follow_up_instructions"`
	FollowUpDate           *string   `json:"follow_up_date"`
	FollowUpDoctor         *string   `json:"follow_up_doctor"`
	FollowUpSpecialty      *string   `json:"follow_up_specialty"`
	WarningSymptoms        string    `json:"warning_symptoms"`
	WhenToReturnHospital   string    `json:"when_to_return_hospital"`
	Outcome                string    `json:"outcome" binding:"required"`
	PatientEducation       string    `json:"patient_education"`
	Notes                  string    `json:"notes"`
}

func (h *DoctorHandler) CreateDischargeSummary(c *gin.Context) {
	var req CreateDischargeSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	discharge := &model.DischargeSummary{
		AdmissionID:            req.AdmissionID,
		PatientID:              req.PatientID,
		DoctorID:               req.DoctorID,
		DischargeDate:          req.DischargeDate,
		FinalDiagnosis:         req.FinalDiagnosis,
		ProceduresPerformed:    req.ProceduresPerformed,
		ComplicationsIfAny:     req.ComplicationsIfAny,
		DischargeMedications:   req.DischargeMedications,
		MedicationInstructions: req.MedicationInstructions,
		DietRecommendation:     req.DietRecommendation,
		ActivityRecommendation: req.ActivityRecommendation,
		FollowUpInstructions:   req.FollowUpInstructions,
		FollowUpDate:           req.FollowUpDate,
		FollowUpDoctor:         req.FollowUpDoctor,
		FollowUpSpecialty:      req.FollowUpSpecialty,
		WarningSymptoms:        req.WarningSymptoms,
		WhenToReturnHospital:   req.WhenToReturnHospital,
		Outcome:                model.PatientOutcome(req.Outcome),
		PatientEducation:       req.PatientEducation,
		Notes:                  req.Notes,
	}

	if err := h.dischargeSummaryRepository.Create(discharge); err != nil {
		utils.Fail(c, 400, "failed to create discharge summary: "+err.Error())
		return
	}

	utils.OK(c, discharge)
}

// ==================== PROGRESS NOTES - CRUD OPERATIONS ====================

// GetProgressNotes retrieves all progress notes with optional filtering
// GET /api/v1/progress-notes
func (h *DoctorHandler) GetProgressNotes(c *gin.Context) {
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

	notes, total, err := h.progressNoteRepository.GetAll(page, limit)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch progress notes: "+err.Error())
		return
	}

	utils.OK(c, gin.H{
		"data":  notes,
		"total": total,
		"page":  page,
	})
}

// GetProgressNote retrieves a specific progress note
// GET /api/v1/progress-notes/:id
func (h *DoctorHandler) GetProgressNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid progress note id")
		return
	}

	note, err := h.progressNoteRepository.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, "progress note not found")
		return
	}

	utils.OK(c, note)
}

// UpdateProgressNote updates a progress note
// PATCH /api/v1/progress-notes/:id
type UpdateProgressNoteRequest struct {
	Subjective     string  `json:"subjective"`
	Objective      string  `json:"objective"`
	Assessment     string  `json:"assessment"`
	Plan           string  `json:"plan"`
	Vitals         string  `json:"vitals"`
	Medications    string  `json:"medications"`
	NextReviewDate *string `json:"next_review_date"`
	Notes          string  `json:"notes"`
}

func (h *DoctorHandler) UpdateProgressNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid progress note id")
		return
	}

	var req UpdateProgressNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	note, err := h.progressNoteRepository.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, "progress note not found")
		return
	}

	if req.Subjective != "" {
		note.Subjective = req.Subjective
	}
	if req.Objective != "" {
		note.Objective = req.Objective
	}
	if req.Assessment != "" {
		note.Assessment = req.Assessment
	}
	if req.Plan != "" {
		note.Plan = req.Plan
	}
	if req.Vitals != "" {
		note.Vitals = req.Vitals
	}
	if req.Medications != "" {
		note.Medications = req.Medications
	}
	if req.NextReviewDate != nil {
		note.NextReviewDate = req.NextReviewDate
	}
	if req.Notes != "" {
		note.Notes = req.Notes
	}

	if err := h.progressNoteRepository.Update(note); err != nil {
		utils.Fail(c, 400, "failed to update progress note: "+err.Error())
		return
	}

	utils.OK(c, note)
}

// DeleteProgressNote deletes a progress note
// DELETE /api/v1/progress-notes/:id
func (h *DoctorHandler) DeleteProgressNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid progress note id")
		return
	}

	if err := h.progressNoteRepository.SoftDelete(id); err != nil {
		utils.Fail(c, 400, "failed to delete progress note: "+err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message": "progress note deleted successfully",
	})
}

// GetAdmissionProgressNotes retrieves all progress notes for an admission
// GET /api/v1/progress-notes/admission/:admission_id
func (h *DoctorHandler) GetAdmissionProgressNotes(c *gin.Context) {
	admissionIDStr := c.Param("admission_id")
	admissionID, err := uuid.Parse(admissionIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid admission id")
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

	notes, total, err := h.progressNoteRepository.GetByAdmissionID(admissionID, page, limit)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch progress notes: "+err.Error())
		return
	}

	utils.OK(c, gin.H{
		"data":  notes,
		"total": total,
		"page":  page,
	})
}

// ==================== NURSE INSTRUCTIONS - CRUD OPERATIONS ====================

// GetNurseInstructions retrieves all nurse instructions
// GET /api/v1/nurse-instructions
func (h *DoctorHandler) GetNurseInstructions(c *gin.Context) {
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

	instructions, total, err := h.nurseInstructionRepository.GetAll(page, limit)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch nurse instructions: "+err.Error())
		return
	}

	utils.OK(c, gin.H{
		"data":  instructions,
		"total": total,
		"page":  page,
	})
}

// GetNurseInstruction retrieves a specific nurse instruction
// GET /api/v1/nurse-instructions/:id
func (h *DoctorHandler) GetNurseInstruction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid nurse instruction id")
		return
	}

	instruction, err := h.nurseInstructionRepository.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, "nurse instruction not found")
		return
	}

	utils.OK(c, instruction)
}

// UpdateNurseInstruction updates a nurse instruction
// PATCH /api/v1/nurse-instructions/:id
type UpdateNurseInstructionRequest struct {
	VitalsFrequency     string `json:"vitals_frequency"`
	VitalsParameters    string `json:"vitals_parameters"`
	DietType            string `json:"diet_type"`
	DietaryRestrictions string `json:"dietary_restrictions"`
	FluidsRestriction   string `json:"fluids_restriction"`
	MedicationNotes     string `json:"medication_notes"`
	DrugAllergies       string `json:"drug_allergies"`
	ActivityLevel       string `json:"activity_level"`
	ActivityNotes       string `json:"activity_notes"`
	SpecialMonitoring   string `json:"special_monitoring"`
	CautionPoints       string `json:"caution_points"`
	HygieneInstructions string `json:"hygiene_instructions"`
	OtherInstructions   string `json:"other_instructions"`
	Status              string `json:"status"`
}

func (h *DoctorHandler) UpdateNurseInstruction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid nurse instruction id")
		return
	}

	var req UpdateNurseInstructionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	instruction, err := h.nurseInstructionRepository.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, "nurse instruction not found")
		return
	}

	if req.VitalsFrequency != "" {
		instruction.VitalsFrequency = req.VitalsFrequency
	}
	if req.VitalsParameters != "" {
		instruction.VitalsParameters = req.VitalsParameters
	}
	if req.DietType != "" {
		instruction.DietType = req.DietType
	}
	if req.DietaryRestrictions != "" {
		instruction.DietaryRestrictions = req.DietaryRestrictions
	}
	if req.FluidsRestriction != "" {
		instruction.FluidsRestriction = req.FluidsRestriction
	}
	if req.MedicationNotes != "" {
		instruction.MedicationNotes = req.MedicationNotes
	}
	if req.DrugAllergies != "" {
		instruction.DrugAllergies = req.DrugAllergies
	}
	if req.ActivityLevel != "" {
		instruction.ActivityLevel = req.ActivityLevel
	}
	if req.ActivityNotes != "" {
		instruction.ActivityNotes = req.ActivityNotes
	}
	if req.SpecialMonitoring != "" {
		instruction.SpecialMonitoring = req.SpecialMonitoring
	}
	if req.CautionPoints != "" {
		instruction.CautionPoints = req.CautionPoints
	}
	if req.HygieneInstructions != "" {
		instruction.HygieneInstructions = req.HygieneInstructions
	}
	if req.OtherInstructions != "" {
		instruction.OtherInstructions = req.OtherInstructions
	}
	if req.Status != "" {
		instruction.Status = req.Status
	}

	if err := h.nurseInstructionRepository.Update(instruction); err != nil {
		utils.Fail(c, 400, "failed to update nurse instruction: "+err.Error())
		return
	}

	utils.OK(c, instruction)
}

// DeleteNurseInstruction deletes a nurse instruction
// DELETE /api/v1/nurse-instructions/:id
func (h *DoctorHandler) DeleteNurseInstruction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid nurse instruction id")
		return
	}

	if err := h.nurseInstructionRepository.SoftDelete(id); err != nil {
		utils.Fail(c, 400, "failed to delete nurse instruction: "+err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message": "nurse instruction deleted successfully",
	})
}

// GetAdmissionNurseInstructions retrieves all nurse instructions for an admission
// GET /api/v1/nurse-instructions/admission/:admission_id
func (h *DoctorHandler) GetAdmissionNurseInstructions(c *gin.Context) {
	admissionIDStr := c.Param("admission_id")
	admissionID, err := uuid.Parse(admissionIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid admission id")
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

	instructions, total, err := h.nurseInstructionRepository.GetByAdmissionID(admissionID, page, limit)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch nurse instructions: "+err.Error())
		return
	}

	utils.OK(c, gin.H{
		"data":  instructions,
		"total": total,
		"page":  page,
	})
}

// ==================== DISCHARGE SUMMARIES - CRUD OPERATIONS ====================

// GetDischargeSummaries retrieves all discharge summaries
// GET /api/v1/discharge-summaries
func (h *DoctorHandler) GetDischargeSummaries(c *gin.Context) {
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

	summaries, total, err := h.dischargeSummaryRepository.GetAll(page, limit)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch discharge summaries: "+err.Error())
		return
	}

	utils.OK(c, gin.H{
		"data":  summaries,
		"total": total,
		"page":  page,
	})
}

// GetDischargeSummary retrieves a specific discharge summary
// GET /api/v1/discharge-summaries/:id
func (h *DoctorHandler) GetDischargeSummary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid discharge summary id")
		return
	}

	summary, err := h.dischargeSummaryRepository.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, "discharge summary not found")
		return
	}

	utils.OK(c, summary)
}

// UpdateDischargeSummary updates a discharge summary
// PATCH /api/v1/discharge-summaries/:id
type UpdateDischargeSummaryRequest struct {
	FinalDiagnosis         string  `json:"final_diagnosis"`
	ProceduresPerformed    string  `json:"procedures_performed"`
	ComplicationsIfAny     string  `json:"complications_if_any"`
	DischargeMedications   string  `json:"discharge_medications"`
	MedicationInstructions string  `json:"medication_instructions"`
	DietRecommendation     string  `json:"diet_recommendation"`
	ActivityRecommendation string  `json:"activity_recommendation"`
	FollowUpInstructions   string  `json:"follow_up_instructions"`
	FollowUpDate           *string `json:"follow_up_date"`
	FollowUpDoctor         *string `json:"follow_up_doctor"`
	FollowUpSpecialty      *string `json:"follow_up_specialty"`
	WarningSymptoms        string  `json:"warning_symptoms"`
	WhenToReturnHospital   string  `json:"when_to_return_hospital"`
	Outcome                string  `json:"outcome"`
	PatientEducation       string  `json:"patient_education"`
	Notes                  string  `json:"notes"`
}

func (h *DoctorHandler) UpdateDischargeSummary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid discharge summary id")
		return
	}

	var req UpdateDischargeSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	summary, err := h.dischargeSummaryRepository.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, "discharge summary not found")
		return
	}

	if req.FinalDiagnosis != "" {
		summary.FinalDiagnosis = req.FinalDiagnosis
	}
	if req.ProceduresPerformed != "" {
		summary.ProceduresPerformed = req.ProceduresPerformed
	}
	if req.ComplicationsIfAny != "" {
		summary.ComplicationsIfAny = req.ComplicationsIfAny
	}
	if req.DischargeMedications != "" {
		summary.DischargeMedications = req.DischargeMedications
	}
	if req.MedicationInstructions != "" {
		summary.MedicationInstructions = req.MedicationInstructions
	}
	if req.DietRecommendation != "" {
		summary.DietRecommendation = req.DietRecommendation
	}
	if req.ActivityRecommendation != "" {
		summary.ActivityRecommendation = req.ActivityRecommendation
	}
	if req.FollowUpInstructions != "" {
		summary.FollowUpInstructions = req.FollowUpInstructions
	}
	if req.FollowUpDate != nil {
		summary.FollowUpDate = req.FollowUpDate
	}
	if req.FollowUpDoctor != nil {
		summary.FollowUpDoctor = req.FollowUpDoctor
	}
	if req.FollowUpSpecialty != nil {
		summary.FollowUpSpecialty = req.FollowUpSpecialty
	}
	if req.WarningSymptoms != "" {
		summary.WarningSymptoms = req.WarningSymptoms
	}
	if req.WhenToReturnHospital != "" {
		summary.WhenToReturnHospital = req.WhenToReturnHospital
	}
	if req.Outcome != "" {
		summary.Outcome = model.PatientOutcome(req.Outcome)
	}
	if req.PatientEducation != "" {
		summary.PatientEducation = req.PatientEducation
	}
	if req.Notes != "" {
		summary.Notes = req.Notes
	}

	if err := h.dischargeSummaryRepository.Update(summary); err != nil {
		utils.Fail(c, 400, "failed to update discharge summary: "+err.Error())
		return
	}

	utils.OK(c, summary)
}

// DeleteDischargeSummary deletes a discharge summary
// DELETE /api/v1/discharge-summaries/:id
func (h *DoctorHandler) DeleteDischargeSummary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid discharge summary id")
		return
	}

	if err := h.dischargeSummaryRepository.SoftDelete(id); err != nil {
		utils.Fail(c, 400, "failed to delete discharge summary: "+err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message": "discharge summary deleted successfully",
	})
}

// GetAdmissionDischargeSummary retrieves discharge summary for an admission
// GET /api/v1/discharge-summaries/admission/:admission_id
func (h *DoctorHandler) GetAdmissionDischargeSummary(c *gin.Context) {
	admissionIDStr := c.Param("admission_id")
	admissionID, err := uuid.Parse(admissionIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid admission id")
		return
	}

	summary, err := h.dischargeSummaryRepository.GetByAdmissionID(admissionID)
	if err != nil {
		utils.Fail(c, 404, "discharge summary not found for this admission")
		return
	}

	utils.OK(c, summary)
}

// ==================== ANALYTICS & REPORTING ====================

// GetDoctorAnalytics retrieves analytics dashboard data for doctors
// GET /api/v1/doctors/:id/analytics
func (h *DoctorHandler) GetDoctorAnalytics(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor id")
		return
	}

	// Get doctor info
	doctor, err := h.doctorService.GetDoctor(id)
	if err != nil {
		utils.Fail(c, 404, "doctor not found")
		return
	}

	// Get active admissions count
	admissions, _, _ := h.admissionRepository.GetActiveByDoctorID(id, 1, 1000)
	activeAdmissionsCount := len(admissions)

	// Get appointment count (would need appointmentRepo)
	// For now, just return basic stats
	analytics := gin.H{
		"doctor_id":               id.String(),
		"doctor_name":             doctor.User.FirstName + " " + doctor.User.LastName,
		"specialization":          doctor.Specialization,
		"consultation_fee":        doctor.ConsultationFee,
		"active_admissions":       activeAdmissionsCount,
		"total_appointments":      0, // Would query from appointment repository
		"total_prescriptions":     0, // Would query from prescription repository
		"total_lab_orders":        0, // Would query from lab repository
		"patient_satisfied_count": 0, // Would query from ratings
	}

	utils.OK(c, analytics)
}

// GetMyAnalytics retrieves analytics for the current logged-in doctor
// GET /api/v1/doctors/analytics/my
func (h *DoctorHandler) GetMyAnalytics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, 401, "unauthorized")
		return
	}

	doctor, err := h.doctorService.GetDoctorByUserID(userID.(uuid.UUID))
	if err != nil {
		utils.Fail(c, 404, "doctor not found")
		return
	}

	// Get active admissions
	admissions, _, _ := h.admissionRepository.GetActiveByDoctorID(doctor.ID, 1, 1000)
	activeAdmissionsCount := len(admissions)

	// Get progress notes for this doctor
	progressNotes, _, _ := h.progressNoteRepository.GetByDoctorID(doctor.ID, 1, 1000)
	progressNotesCount := len(progressNotes)

	// Get discharge summaries
	dischargeSummaries, _, _ := h.dischargeSummaryRepository.GetByDoctorID(doctor.ID, 1, 1000)
	dischargeSummariesCount := len(dischargeSummaries)

	analytics := gin.H{
		"doctor_id":             doctor.ID.String(),
		"doctor_name":           doctor.User.FirstName + " " + doctor.User.LastName,
		"specialization":        doctor.Specialization,
		"consultation_fee":      doctor.ConsultationFee,
		"active_admissions":     activeAdmissionsCount,
		"total_progress_notes":  progressNotesCount,
		"total_discharge_cases": dischargeSummariesCount,
		"attendance_percentage": doctor.AttendancePercentage,
		"leave_balance":         doctor.LeaveBalance,
	}

	utils.OK(c, analytics)
}

// GetAdmissionStats retrieves hospital admission statistics
// GET /api/v1/analytics/admissions
func (h *DoctorHandler) GetAdmissionStats(c *gin.Context) {
	// Get all admissions with pagination
	allAdmissions, total, err := h.admissionRepository.GetAll(1, 10000)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch admissions: "+err.Error())
		return
	}

	// Count by status
	activeCount := 0
	dischargedCount := 0
	cancelledCount := 0

	for _, admission := range allAdmissions {
		switch admission.Status {
		case model.AdmissionActive:
			activeCount++
		case model.AdmissionDischarge:
			dischargedCount++
		case model.AdmissionCancelled:
			cancelledCount++
		}
	}

	stats := gin.H{
		"total_admissions":    total,
		"active_admissions":   activeCount,
		"discharged_patients": dischargedCount,
		"cancelled_admission": cancelledCount,
		"occupancy_rate":      float64(activeCount) / float64(total) * 100,
	}

	utils.OK(c, stats)
}

// GetDoctorPerformanceMetrics retrieves performance metrics for a doctor
// GET /api/v1/doctors/:id/performance
func (h *DoctorHandler) GetDoctorPerformanceMetrics(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor id")
		return
	}

	// Get discharge summaries to check outcomes
	dischargeSummaries, total, _ := h.dischargeSummaryRepository.GetByDoctorID(id, 1, 10000)

	recoveredCount := 0
	improvedCount := 0
	stableCount := 0
	referredCount := 0

	for _, summary := range dischargeSummaries {
		switch summary.Outcome {
		case model.OutcomeRecovered:
			recoveredCount++
		case model.OutcomeImproved:
			improvedCount++
		case model.OutcomeStableCondition:
			stableCount++
		case model.OutcomeReferredToSpecialist:
			referredCount++
		}
	}

	recoveryRate := float64(0)
	improvedRate := float64(0)
	if total > 0 {
		recoveryRate = float64(recoveredCount) / float64(total) * 100
		improvedRate = float64(improvedCount) / float64(total) * 100
	}

	metrics := gin.H{
		"total_discharged_patients": total,
		"recovered_count":           recoveredCount,
		"improved_count":            improvedCount,
		"stable_count":              stableCount,
		"referred_count":            referredCount,
		"recovery_rate":             recoveryRate,
		"improved_rate":             improvedRate,
	}

	utils.OK(c, metrics)
}

// GetActiveDoctorsToday returns doctors who logged in today
// GET /api/v1/doctors/active-today
// Requires: admin or receptionist role
func (h *DoctorHandler) GetActiveDoctorsToday(c *gin.Context) {
	// Get start of today in milliseconds
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()

	type ActiveDoctorResponse struct {
		DoctorID       string `json:"doctor_id"`
		UserID         string `json:"user_id"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		Specialization string `json:"specialization"`
		Department     string `json:"department"`
		LastLogin      int64  `json:"last_login"`
	}

	var doctors []model.Doctor
	if err := config.DB.
		Joins("JOIN users ON users.id = doctors.user_id").
		Where("users.role = ? AND users.last_login >= ? AND users.deleted_at IS NULL", "doctor", startOfDay).
		Preload("User").
		Find(&doctors).Error; err != nil {
		utils.Fail(c, 500, "failed to fetch active doctors: "+err.Error())
		return
	}

	responses := make([]ActiveDoctorResponse, 0, len(doctors))
	for _, d := range doctors {
		var lastLogin int64
		if d.User != nil && d.User.LastLogin != nil {
			lastLogin = *d.User.LastLogin
		}
		resp := ActiveDoctorResponse{
			DoctorID:       d.ID.String(),
			Specialization: d.Specialization,
			Department:     d.Department,
			LastLogin:      lastLogin,
		}
		if d.User != nil {
			resp.UserID = d.User.ID.String()
			resp.FirstName = d.User.FirstName
			resp.LastName = d.User.LastName
		}
		responses = append(responses, resp)
	}

	utils.OK(c, gin.H{
		"doctors": responses,
		"total":   len(responses),
	})
}
