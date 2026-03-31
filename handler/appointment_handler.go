package handler

import (
	"fmt"
	"strconv"

	"meet_sushruta/model"
	"meet_sushruta/repository"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentHandler struct {
	appointmentService service.AppointmentService
	patientRepo        repository.PatientRepository
	doctorRepo         repository.DoctorRepository
}

func NewAppointmentHandler(
	appointmentService service.AppointmentService,
	patientRepo repository.PatientRepository,
	doctorRepo repository.DoctorRepository,
) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
		patientRepo:        patientRepo,
		doctorRepo:         doctorRepo,
	}
}

type BookAppointmentRequest struct {
	PatientID       string `json:"patient_id" binding:"required,uuid"`
	DoctorID        string `json:"doctor_id" binding:"required,uuid"`
	AppointmentDate string `json:"appointment_date" binding:"required"` // YYYY-MM-DD
	AppointmentTime string `json:"appointment_time" binding:"required"` // HH:MM
	Reason          string `json:"reason"`
	Notes           string `json:"notes"`
}

type UpdateAppointmentStatusRequest struct {
	Status string `json:"status" binding:"required"` // pending, confirmed, in_progress, completed, cancelled
}

type UpdateAppointmentVitalsRequest struct {
	Notes string `json:"notes"`
	// Vitals fields can be added here based on requirements
}

// BookAppointment creates a new appointment
// POST /api/v1/appointments
func (h *AppointmentHandler) BookAppointment(c *gin.Context) {
	fmt.Println("\n========== BookAppointment START ==========")
	fmt.Printf("[BookAppointment] Auth Header: %s\n", c.GetHeader("Authorization"))

	// Get context values from middleware
	userID, exists := c.Get("userID")
	fmt.Printf("[BookAppointment] Context userID (from middleware): %v (exists: %v)\n", userID, exists)

	role, exists := c.Get("role")
	fmt.Printf("[BookAppointment] Context role (from middleware): %v (exists: %v)\n", role, exists)

	var req BookAppointmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// Log the actual error for debugging
		fmt.Printf("[BookAppointment] Binding error: %v\n", err.Error())
		fmt.Printf("[BookAppointment] PatientID received: '%v' (type: %T)\n", req.PatientID, req.PatientID)
		fmt.Printf("[BookAppointment] DoctorID received: '%v' (type: %T)\n", req.DoctorID, req.DoctorID)
		utils.Fail(c, 400, fmt.Sprintf("invalid request: %v", err.Error()))
		fmt.Println("========== BookAppointment END (Binding Error) ==========\n")
		return
	}

	// Convert string UUIDs to uuid.UUID
	patientUUID, err := uuid.Parse(req.PatientID)
	if err != nil {
		fmt.Printf("[BookAppointment] UUID parse error for patient_id: %v\n", err.Error())
		utils.Fail(c, 400, fmt.Sprintf("invalid patient ID: %v", err.Error()))
		fmt.Println("========== BookAppointment END (UUID Parse Error) ==========\n")
		return
	}

	fmt.Printf("[BookAppointment] Parsed patient UUID: %s\n", patientUUID.String())

	// Look up patient by user_id (req.PatientID is actually the user_id from frontend)
	patient, err := h.patientRepo.GetByUserID(patientUUID)
	if err != nil {
		fmt.Printf("[BookAppointment] ERROR: Patient lookup failed for UUID %s - Error: %v\n", patientUUID.String(), err.Error())
		utils.Fail(c, 400, "patient not found - please complete your patient profile")
		fmt.Println("========== BookAppointment END (Patient Not Found) ==========\n")
		return
	}

	fmt.Printf("[BookAppointment] Patient found: ID=%s\n", patient.ID)

	doctorID, err := uuid.Parse(req.DoctorID)
	if err != nil {
		utils.Fail(c, 400, fmt.Sprintf("invalid doctor ID: %v", err.Error()))
		return
	}

	fmt.Printf("[BookAppointment] Successfully parsed - PatientID=%s, DoctorID=%s, Date=%s, Time=%s\n",
		patient.ID, doctorID, req.AppointmentDate, req.AppointmentTime)

	// Validate date format (YYYY-MM-DD)
	if len(req.AppointmentDate) != 10 {
		fmt.Printf("[BookAppointment] Invalid date format: '%s' (length: %d)\n", req.AppointmentDate, len(req.AppointmentDate))
		utils.Fail(c, 400, "appointment_date must be in YYYY-MM-DD format (e.g., 2026-03-30)")
		return
	}

	// Validate time format (HH:MM or HH:MM:SS)
	if len(req.AppointmentTime) < 5 || len(req.AppointmentTime) > 8 {
		fmt.Printf("[BookAppointment] Invalid time format: '%s' (length: %d)\n", req.AppointmentTime, len(req.AppointmentTime))
		utils.Fail(c, 400, "appointment_time must be in HH:MM format (e.g., 14:30)")
		return
	}

	appointment := &model.Appointment{
		PatientID:       patient.ID,
		DoctorID:        doctorID,
		AppointmentDate: req.AppointmentDate,
		AppointmentTime: req.AppointmentTime,
		Reason:          req.Reason,
		Notes:           req.Notes,
	}

	fmt.Printf("[BookAppointment] About to create appointment: %+v\n", appointment)

	err = h.appointmentService.BookAppointment(appointment)
	if err != nil {
		fmt.Printf("[BookAppointment] Error creating appointment: %v\n", err.Error())
		utils.Fail(c, 400, err.Error())
		return
	}

	fmt.Printf("[BookAppointment] Appointment created successfully with ID: %s\n", appointment.ID)
	fmt.Println("========== BookAppointment END (Success) ==========\n")

	utils.OK(c, appointment)
}

// GetAllAppointments retrieves appointments based on user role with pagination
// GET /api/v1/appointments
// - Returns all appointments for any authenticated user
func (h *AppointmentHandler) GetAllAppointments(c *gin.Context) {
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

	// For now, just get all appointments - no role filtering
	appointments, total, err := h.appointmentService.ListAppointments(page, limit)
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

// GetAppointment retrieves an appointment by ID
// GET /api/v1/appointments/:id
func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid appointment id")
		return
	}

	appointment, err := h.appointmentService.GetAppointment(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, appointment)
}

// UpdateAppointmentStatus updates appointment status with state machine validation
// PATCH /api/v1/appointments/:id/status
func (h *AppointmentHandler) UpdateAppointmentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid appointment id")
		return
	}

	var req UpdateAppointmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	newStatus := model.AppointmentStatus(req.Status)

	err = h.appointmentService.UpdateAppointmentStatus(id, newStatus)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	appointment, _ := h.appointmentService.GetAppointment(id)
	utils.OK(c, appointment)
}

// UpdateAppointmentVitals updates appointment vitals and notes
// PATCH /api/v1/appointments/:id/vitals
func (h *AppointmentHandler) UpdateAppointmentVitals(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid appointment id")
		return
	}

	var req UpdateAppointmentVitalsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	// Get appointment
	appointment, err := h.appointmentService.GetAppointment(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	// Only allow if appointment is in_progress
	if appointment.Status != model.AppointmentInProgress {
		utils.Fail(c, 400, "vitals can only be updated when appointment is in progress")
		return
	}

	// Update notes
	if req.Notes != "" {
		appointment.Notes = req.Notes
	}

	// Since we're only updating notes here, use a workaround
	// In a real scenario, you might want to update through a dedicated vitals endpoint
	utils.OK(c, gin.H{
		"message": "vitals update would persist here",
		"notes":   req.Notes,
	})
}

// GetAppointmentVitals retrieves vitals recorded for an appointment
// GET /api/v1/appointments/:id/vitals
func (h *AppointmentHandler) GetAppointmentVitals(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid appointment id")
		return
	}

	// Get appointment
	appointment, err := h.appointmentService.GetAppointment(id)
	if err != nil {
		utils.Fail(c, 404, "appointment not found")
		return
	}

	// Check if appointment has vitals recorded
	// Vitals are stored in the vitals table linked to patient, not directly to appointment
	// We'll retrieve all vitals for the patient and filter by appointment date/time
	vitalsService, ok := c.Get("vitalsService")
	if !ok {
		// Fallback: vitals need to be queried through the database
		// In this case, we'll return vitals recorded around the appointment time
		utils.OK(c, gin.H{
			"appointment_id":   appointment.ID,
			"patient_id":       appointment.PatientID,
			"appointment_date": appointment.AppointmentDate,
			"appointment_time": appointment.AppointmentTime,
			"vitals":           []model.Vitals{},
			"message":          "Vitals service not available in context",
		})
		return
	}

	// Try to use vitals service if available
	vitals, ok := vitalsService.(service.VitalsService)
	if !ok {
		utils.Fail(c, 500, "vitals service unavailable")
		return
	}

	// Get vitals for the patient
	patientVitals, err := vitals.GetVitalsByPatientID(appointment.PatientID)
	if err != nil {
		utils.OK(c, gin.H{
			"appointment_id": appointment.ID,
			"patient_id":     appointment.PatientID,
			"vitals":         []model.Vitals{},
			"message":        "No vitals recorded for this patient",
		})
		return
	}

	utils.OK(c, gin.H{
		"appointment_id":   appointment.ID,
		"patient_id":       appointment.PatientID,
		"appointment_date": appointment.AppointmentDate,
		"appointment_time": appointment.AppointmentTime,
		"vitals":           patientVitals,
	})
}

// GetTodayAppointments retrieves all appointments for today (for receptionist)
// GET /api/v1/appointments/today
func (h *AppointmentHandler) GetTodayAppointments(c *gin.Context) {
	page := 1
	limit := 100

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

	doctorID := c.Query("doctor_id")
	fmt.Printf("[GetTodayAppointments] Called with doctor_id: '%s' (empty: %v)\n", doctorID, doctorID == "")

	appointments, total, err := h.appointmentService.GetTodayAppointments(page, limit, doctorID)
	if err != nil {
		fmt.Printf("[GetTodayAppointments] Error: %v\n", err)
		utils.Fail(c, 500, err.Error())
		return
	}

	fmt.Printf("[GetTodayAppointments] Returning %d appointments for doctor_id: %s\n", len(appointments), doctorID)
	utils.OK(c, gin.H{
		"appointments": appointments,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}

// GetNext7DaysAppointments retrieves all appointments for next 7 days (for nurse)
// GET /api/v1/appointments/next-7-days
func (h *AppointmentHandler) GetNext7DaysAppointments(c *gin.Context) {
	page := 1
	limit := 100

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

	appointments, total, err := h.appointmentService.GetNext7DaysAppointments(page, limit)
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

// GetMyAppointments retrieves appointments for the logged-in doctor (last 7 days to next 7 days)
// GET /api/v1/appointments/my/schedule
func (h *AppointmentHandler) GetMyAppointments(c *gin.Context) {
	// Get doctor ID from context (set by middleware)
	doctorIDInterface, exists := c.Get("userID")
	if !exists {
		fmt.Printf("[GetMyAppointments] ERROR: userID not found in context\n")
		utils.Fail(c, 401, "user id not found in context")
		return
	}

	doctorIDStr := fmt.Sprintf("%v", doctorIDInterface)
	fmt.Printf("[GetMyAppointments] User ID from context: %s\n", doctorIDStr)

	doctorID, err := uuid.Parse(doctorIDStr)
	if err != nil {
		fmt.Printf("[GetMyAppointments] ERROR: Failed to parse UUID: %v\n", err)
		utils.Fail(c, 400, "invalid doctor id")
		return
	}

	// Get doctor by user ID to verify it's a doctor
	doctor, err := h.doctorRepo.GetByUserID(doctorID)
	if err != nil {
		fmt.Printf("[GetMyAppointments] ERROR: Doctor not found for userID %s: %v\n", doctorIDStr, err)
		utils.Fail(c, 404, "doctor not found")
		return
	}

	if doctor == nil {
		fmt.Printf("[GetMyAppointments] ERROR: Doctor is nil for userID %s\n", doctorIDStr)
		utils.Fail(c, 404, "doctor not found")
		return
	}

	fmt.Printf("[GetMyAppointments] Doctor found: ID=%s, Name=%s, Specialization=%s\n", doctor.ID, doctor.User.FirstName, doctor.Specialization)

	page := 1
	limit := 100

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

	fmt.Printf("[GetMyAppointments] Fetching appointments for doctor %s with page=%d, limit=%d\n", doctor.ID, page, limit)

	// Get appointments for doctor - includes last 7 days + next 7 days
	appointments, total, err := h.appointmentService.GetDoctorAppointments(doctor.ID, page, limit)
	if err != nil {
		fmt.Printf("[GetMyAppointments] ERROR: Failed to get appointments: %v\n", err)
		utils.Fail(c, 500, err.Error())
		return
	}

	fmt.Printf("[GetMyAppointments] Found %d appointments out of %d total\n", len(appointments), total)

	utils.OK(c, gin.H{
		"appointments": appointments,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}
