package handler

import (
	"strconv"

	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NurseHandler struct {
	nurseService service.NurseService
}

func NewNurseHandler(nurseService service.NurseService) *NurseHandler {
	return &NurseHandler{
		nurseService: nurseService,
	}
}

type CreateNurseRequest struct {
	UserID               uuid.UUID `json:"user_id" binding:"required"`
	LicenseNumber        string    `json:"license_number" binding:"required"`
	Department           string    `json:"department"`
	Shift                string    `json:"shift"` // morning, evening, night
	CertificationURL     string    `json:"certification_url"`
	JoiningDate          int64     `json:"joining_date"`
	Salary               float64   `json:"salary"`
	AttendancePercentage float64   `json:"attendance_percentage"`
	LeaveBalance         int       `json:"leave_balance"`
}

type UpdateNurseRequest struct {
	Department           string  `json:"department"`
	Shift                string  `json:"shift"`
	CertificationURL     string  `json:"certification_url"`
	JoiningDate          int64   `json:"joining_date"`
	Salary               float64 `json:"salary"`
	AttendancePercentage float64 `json:"attendance_percentage"`
	LeaveBalance         int     `json:"leave_balance"`
}

// CreateNurse creates a new nurse
// POST /api/v1/nurses
func (h *NurseHandler) CreateNurse(c *gin.Context) {
	var req CreateNurseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	nurse := &model.Nurse{
		UserID:               req.UserID,
		LicenseNumber:        req.LicenseNumber,
		Department:           req.Department,
		Shift:                req.Shift,
		CertificationURL:     req.CertificationURL,
		JoiningDate:          req.JoiningDate,
		Salary:               req.Salary,
		AttendancePercentage: req.AttendancePercentage,
		LeaveBalance:         req.LeaveBalance,
	}

	err := h.nurseService.CreateNurse(nurse)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, nurse)
}

// GetAllNurses retrieves all nurses with pagination
// GET /api/v1/nurses
func (h *NurseHandler) GetAllNurses(c *gin.Context) {
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

	nurses, total, err := h.nurseService.ListNurses(page, limit, search)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, map[string]interface{}{
		"data":  nurses,
		"total": total,
	})
}

// GetNurse retrieves a nurse by ID
// GET /api/v1/nurses/:id
func (h *NurseHandler) GetNurse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid nurse id")
		return
	}

	nurse, err := h.nurseService.GetNurse(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, nurse)
}

// UpdateNurse updates a nurse
// PATCH /api/v1/nurses/:id
func (h *NurseHandler) UpdateNurse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid nurse id")
		return
	}

	var req UpdateNurseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	updates := map[string]interface{}{
		"department":            req.Department,
		"shift":                 req.Shift,
		"certification_url":     req.CertificationURL,
		"joining_date":          req.JoiningDate,
		"salary":                req.Salary,
		"attendance_percentage": req.AttendancePercentage,
		"leave_balance":         req.LeaveBalance,
	}

	nurse, err := h.nurseService.UpdateNurse(id, updates)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, nurse)
}

// DeleteNurse deletes a nurse
// DELETE /api/v1/nurses/:id
func (h *NurseHandler) DeleteNurse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid nurse id")
		return
	}

	err = h.nurseService.DeleteNurse(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, map[string]string{"message": "nurse deleted successfully"})
}

// GetNurseByUserID retrieves a nurse by user ID
// GET /api/v1/nurses/user/:user_id
func (h *NurseHandler) GetNurseByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid user id")
		return
	}

	nurse, err := h.nurseService.GetNurseByUserID(userID)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, nurse)
}
