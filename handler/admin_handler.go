package handler

import (
	"fmt"
	"meet_sushruta/config"
	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// GetDashboardStats returns admin dashboard statistics
// GET /api/v1/admin/dashboard
// Requires: JWT token + admin role
func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.adminService.GetDashboardStats()
	if err != nil {
		utils.Fail(c, 500, "Failed to fetch dashboard stats: "+err.Error())
		return
	}

	utils.OK(c, stats)
}

// RegisterDoctor creates a new doctor user
// POST /api/v1/admin/doctors/register
// Requires: JWT token + admin role
type RegisterDoctorRequest struct {
	FirstName            string  `json:"first_name" binding:"required,min=2,max=100"`
	LastName             string  `json:"last_name" binding:"required,min=2,max=100"`
	Email                string  `json:"email" binding:"required,email"`
	Password             string  `json:"password" binding:"required,min=8"`
	Phone                string  `json:"phone" binding:"required"`
	Specialization       string  `json:"specialization" binding:"required"`
	LicenseNumber        string  `json:"license_number" binding:"required"`
	CertificationURL     string  `json:"certification_url"`
	Bio                  string  `json:"bio"`
	Department           string  `json:"department"`
	Fees                 float64 `json:"fees"`
	JoiningDate          int64   `json:"joining_date"`
	Salary               float64 `json:"salary"`
	AttendancePercentage float64 `json:"attendance_percentage"`
	LeaveBalance         int     `json:"leave_balance"`
}

func (h *AdminHandler) RegisterDoctor(c *gin.Context) {
	var req RegisterDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" {
		utils.Fail(c, 400, "first_name, last_name, and email are required")
		return
	}

	if req.LicenseNumber == "" {
		utils.Fail(c, 400, "license_number is required")
		return
	}

	// Convert to service request type
	serviceReq := service.RegisterDoctorRequest{
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		Email:                req.Email,
		Password:             req.Password,
		Phone:                req.Phone,
		Specialization:       req.Specialization,
		LicenseNumber:        req.LicenseNumber,
		CertificationURL:     req.CertificationURL,
		Bio:                  req.Bio,
		Department:           req.Department,
		Fees:                 req.Fees,
		JoiningDate:          req.JoiningDate,
		Salary:               req.Salary,
		AttendancePercentage: req.AttendancePercentage,
		LeaveBalance:         req.LeaveBalance,
	}

	resp, err := h.adminService.RegisterDoctor(serviceReq)
	if err != nil {
		// Check what kind of error occurred
		errMsg := err.Error()
		if errMsg == "email already registered" {
			utils.Fail(c, 409, "This email address is already registered. Please use a different email.")
			return
		}
		if errMsg == "phone number already registered" {
			utils.Fail(c, 409, "This phone number is already registered. Please use a different phone number.")
			return
		}
		if errMsg == "license number already registered" {
			utils.Fail(c, 409, "This license number is already registered. Please use a different license number.")
			return
		}
		utils.Fail(c, 400, errMsg)
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"data":    resp,
	})
}

// GetTodayAppointmentsByDoctor retrieves all appointments for today grouped by doctor
// GET /api/v1/admin/appointments/today
// Requires: JWT token + admin role
func (h *AdminHandler) GetTodayAppointmentsByDoctor(c *gin.Context) {
	appointments, err := h.adminService.GetTodayAppointmentsByDoctor()
	if err != nil {
		utils.Fail(c, 500, "Failed to fetch today's appointments: "+err.Error())
		return
	}

	utils.OK(c, gin.H{
		"appointments": appointments,
	})
}

// GetStaffByRole retrieves all users with a specific role
// GET /api/v1/users?role={role}
// Supported roles: pharmacy, lab, nurse
// Requires: JWT token + admin role
func (h *AdminHandler) GetStaffByRole(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		utils.Fail(c, 400, "role query parameter is required")
		return
	}

	// Validate role
	validRoles := map[string]bool{
		"pharmacy": true,
		"lab":      true,
		"nurse":    true,
	}
	if !validRoles[role] {
		utils.Fail(c, 400, "invalid role. supported roles: pharmacy, lab, nurse")
		return
	}

	var staff []model.User
	result := config.DB.Where("role = ?", role).Find(&staff)
	if result.Error != nil {
		utils.Fail(c, 500, "Failed to fetch staff: "+result.Error.Error())
		return
	}

	utils.OK(c, staff)
}

// UpdateStaffUser updates HR fields for a staff user (pharmacy/lab/nurse)
// PATCH /api/v1/users/:id
// Requires: JWT token + admin role
func (h *AdminHandler) UpdateStaffUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid user id")
		return
	}

	type UpdateStaffRequest struct {
		JoiningDate          int64   `json:"joining_date"`
		Salary               float64 `json:"salary"`
		AttendancePercentage float64 `json:"attendance_percentage"`
		LeaveBalance         int     `json:"leave_balance"`
	}

	var req UpdateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	var user model.User
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.Fail(c, 404, "user not found")
		return
	}

	if err := config.DB.Model(&user).
		Select("joining_date", "salary", "attendance_percentage", "leave_balance").
		Updates(model.User{
			JoiningDate:          req.JoiningDate,
			Salary:               req.Salary,
			AttendancePercentage: req.AttendancePercentage,
			LeaveBalance:         req.LeaveBalance,
		}).Error; err != nil {
		utils.Fail(c, 500, "failed to update user: "+err.Error())
		return
	}

	// Return updated user (re-fetch to get fresh data)
	config.DB.Where("id = ?", id).First(&user)
	utils.OK(c, user)
}

// DeleteStaffUser soft-deletes a staff user
// DELETE /api/v1/users/:id
// Requires: JWT token + admin role
func (h *AdminHandler) DeleteStaffUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid user id")
		return
	}

	// Prevent deleting admin accounts
	var user model.User
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.Fail(c, 404, "user not found")
		return
	}
	if user.Role == "admin" {
		utils.Fail(c, 403, "cannot delete admin account")
		return
	}

	if err := config.DB.Delete(&model.User{}, "id = ?", id).Error; err != nil {
		utils.Fail(c, 500, "failed to delete user: "+err.Error())
		return
	}

	utils.OK(c, map[string]string{"message": "user deleted successfully"})
}

// GetBedStats returns bed availability statistics ward-wise
// GET /api/v1/admin/inventory/beds/stats
// Requires: JWT token + admin role
func (h *AdminHandler) GetBedStats(c *gin.Context) {
	stats, err := h.adminService.GetBedStats()
	if err != nil {
		utils.Fail(c, 500, "Failed to fetch bed stats: "+err.Error())
		return
	}
	utils.OK(c, stats)
}

// GetAllBeds returns all beds with pagination
// GET /api/v1/admin/inventory/beds
// Requires: JWT token + admin role
func (h *AdminHandler) GetAllBeds(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	var p, l int
	_, err := fmt.Sscanf(page, "%d", &p)
	if err != nil || p < 1 {
		p = 1
	}
	_, err = fmt.Sscanf(limit, "%d", &l)
	if err != nil || l < 1 {
		l = 20
	}

	beds, total, err := h.adminService.GetAllBeds(p, l)
	if err != nil {
		utils.Fail(c, 500, "Failed to fetch beds: "+err.Error())
		return
	}

	utils.OK(c, map[string]interface{}{
		"data":  beds,
		"total": total,
		"page":  p,
		"limit": l,
	})
}

// GetPatientsByBedType returns patients in a specific bed type
// GET /api/v1/admin/inventory/beds/type/:bedType
// Requires: JWT token + admin role
func (h *AdminHandler) GetPatientsByBedType(c *gin.Context) {
	bedType := c.Param("bedType")
	if bedType == "" {
		utils.Fail(c, 400, "Bed type is required")
		return
	}

	patients, err := h.adminService.GetPatientsByBedType(bedType)
	if err != nil {
		utils.Fail(c, 500, "Failed to fetch patients: "+err.Error())
		return
	}

	utils.OK(c, patients)
}

// GetMedicalEquipmentStats returns equipment statistics
// GET /api/v1/admin/inventory/equipment/stats
// Requires: JWT token + admin role
func (h *AdminHandler) GetMedicalEquipmentStats(c *gin.Context) {
	stats, err := h.adminService.GetMedicalEquipmentStats()
	if err != nil {
		utils.Fail(c, 500, "Failed to fetch equipment stats: "+err.Error())
		return
	}
	utils.OK(c, stats)
}

// GetAllMedicalEquipment returns all medical equipment
// GET /api/v1/admin/inventory/equipment
// Requires: JWT token + admin role
func (h *AdminHandler) GetAllMedicalEquipment(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")
	status := c.Query("status")

	var p, l int
	_, err := fmt.Sscanf(page, "%d", &p)
	if err != nil || p < 1 {
		p = 1
	}
	_, err = fmt.Sscanf(limit, "%d", &l)
	if err != nil || l < 1 {
		l = 20
	}

	equipment, total, err := h.adminService.GetAllMedicalEquipment(p, l, status)
	if err != nil {
		utils.Fail(c, 500, "Failed to fetch equipment: "+err.Error())
		return
	}

	utils.OK(c, map[string]interface{}{
		"data":  equipment,
		"total": total,
		"page":  p,
		"limit": l,
	})
}

// GetOperationTheatreStats returns OT statistics
// GET /api/v1/admin/inventory/ot/stats
// Requires: JWT token + admin role
func (h *AdminHandler) GetOperationTheatreStats(c *gin.Context) {
	stats, err := h.adminService.GetOperationTheatreStats()
	if err != nil {
		utils.Fail(c, 500, "Failed to fetch OT stats: "+err.Error())
		return
	}
	utils.OK(c, stats)
}

// GetOperationSchedules returns operation schedules
// GET /api/v1/admin/inventory/operation-schedules
// Requires: JWT token + admin role
func (h *AdminHandler) GetOperationSchedules(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")
	date := c.Query("date")

	var p, l int
	_, err := fmt.Sscanf(page, "%d", &p)
	if err != nil || p < 1 {
		p = 1
	}
	_, err = fmt.Sscanf(limit, "%d", &l)
	if err != nil || l < 1 {
		l = 20
	}

	operations, total, err := h.adminService.GetOperationSchedules(p, l, date)
	if err != nil {
		utils.Fail(c, 500, "Failed to fetch operation schedules: "+err.Error())
		return
	}

	utils.OK(c, map[string]interface{}{
		"data":  operations,
		"total": total,
		"page":  p,
		"limit": l,
	})
}
