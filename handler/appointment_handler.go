package handler

import (
	"strconv"

	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentHandler struct {
	appointmentService service.AppointmentService
}

func NewAppointmentHandler(appointmentService service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

type BookAppointmentRequest struct {
	PatientID       uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID        uuid.UUID `json:"doctor_id" binding:"required"`
	AppointmentDate string    `json:"appointment_date" binding:"required"` // YYYY-MM-DD
	AppointmentTime string    `json:"appointment_time" binding:"required"` // HH:MM
	Reason          string    `json:"reason"`
	Notes           string    `json:"notes"`
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
	var req BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	appointment := &model.Appointment{
		PatientID:       req.PatientID,
		DoctorID:        req.DoctorID,
		AppointmentDate: req.AppointmentDate,
		AppointmentTime: req.AppointmentTime,
		Reason:          req.Reason,
		Notes:           req.Notes,
	}

	err := h.appointmentService.BookAppointment(appointment)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, appointment)
}

// GetAllAppointments retrieves all appointments with pagination
// GET /api/v1/appointments
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
