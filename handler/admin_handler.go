package handler

import (
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
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
	FirstName        string  `json:"first_name" binding:"required,min=2,max=100"`
	LastName         string  `json:"last_name" binding:"required,min=2,max=100"`
	Email            string  `json:"email" binding:"required,email"`
	Password         string  `json:"password" binding:"required,min=8"`
	Phone            string  `json:"phone" binding:"required"`
	Specialization   string  `json:"specialization" binding:"required"`
	LicenseNumber    string  `json:"license_number" binding:"required"`
	CertificationURL string  `json:"certification_url"`
	Bio              string  `json:"bio"`
	Department       string  `json:"department"`
	Fees             float64 `json:"fees"`
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
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		Password:         req.Password,
		Phone:            req.Phone,
		Specialization:   req.Specialization,
		LicenseNumber:    req.LicenseNumber,
		CertificationURL: req.CertificationURL,
		Bio:              req.Bio,
		Department:       req.Department,
		Fees:             req.Fees,
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
