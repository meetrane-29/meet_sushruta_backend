package handler

import (
	"strconv"
	"time"

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
	limit := 10
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

	utils.OK(c, doctor)
}

// CreateAdmission creates a new admission record
// POST /api/v1/admissions
type CreateAdmissionRequest struct {
	PatientID     uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID      uuid.UUID `json:"doctor_id" binding:"required"`
	BedID         uuid.UUID `json:"bed_id"`
	AdmissionDate string    `json:"admission_date" binding:"required"` // YYYY-MM-DD HH:MM
	Reason        string    `json:"reason" binding:"required"`
	Diagnosis     string    `json:"diagnosis"`
	Ward          string    `json:"ward"`
	RoomNumber    string    `json:"room_number"`
	IsEmergency   bool      `json:"is_emergency"`
	Notes         string    `json:"notes"`
}

func (h *DoctorHandler) CreateAdmission(c *gin.Context) {
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

	if err := h.admissionRepository.Create(admission); err != nil {
		utils.Fail(c, 400, "failed to create admission: "+err.Error())
		return
	}

	utils.OK(c, admission)
}

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
	AdmissionID    uuid.UUID `json:"admission_id" binding:"required"`
	PatientID      uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID       uuid.UUID `json:"doctor_id" binding:"required"`
	RecordedDate   string    `json:"recorded_date" binding:"required"` // YYYY-MM-DD
	Subjective     string    `json:"subjective"`
	Objective      string    `json:"objective"`
	Assessment     string    `json:"assessment"`
	Plan           string    `json:"plan"`
	Vitals         string    `json:"vitals"`
	Medications    string    `json:"medications"`
	NextReviewDate *string   `json:"next_review_date"`
	Notes          string    `json:"notes"`
}

func (h *DoctorHandler) CreateProgressNote(c *gin.Context) {
	var req CreateProgressNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	note := &model.ProgressNote{
		AdmissionID:    req.AdmissionID,
		PatientID:      req.PatientID,
		DoctorID:       req.DoctorID,
		RecordedDate:   req.RecordedDate,
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
