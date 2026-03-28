package handler

import (
	"strconv"
	"time"

	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DoctorHandler struct {
	doctorService service.DoctorService
}

func NewDoctorHandler(doctorService service.DoctorService) *DoctorHandler {
	return &DoctorHandler{
		doctorService: doctorService,
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
