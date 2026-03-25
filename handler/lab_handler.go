package handler

import (
	"strconv"

	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LabHandler struct {
	labService service.LabService
}

func NewLabHandler(labService service.LabService) *LabHandler {
	return &LabHandler{
		labService: labService,
	}
}

type CreateLabOrderRequest struct {
	PatientID     uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID      uuid.UUID `json:"doctor_id" binding:"required"`
	TestType      string    `json:"test_type" binding:"required"`
	TestName      string    `json:"test_name" binding:"required"`
	ScheduledDate *string   `json:"scheduled_date"`
	Priority      string    `json:"priority"`
}

type UpdateLabStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// CreateOrder creates a new lab order
// POST /api/v1/lab/orders
func (h *LabHandler) CreateOrder(c *gin.Context) {
	var req CreateLabOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	input := &service.LabOrderInput{
		PatientID:     req.PatientID,
		DoctorID:      req.DoctorID,
		TestType:      req.TestType,
		TestName:      req.TestName,
		ScheduledDate: req.ScheduledDate,
		Priority:      req.Priority,
	}

	labOrder, err := h.labService.CreateOrder(input)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, labOrder)
}

// GetOrder retrieves a lab order by ID
// GET /api/v1/lab/orders/:id
func (h *LabHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid lab order id")
		return
	}

	labOrder, err := h.labService.GetOrder(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, labOrder)
}

// ListOrders retrieves all lab orders with pagination
// GET /api/v1/lab/orders
func (h *LabHandler) ListOrders(c *gin.Context) {
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

	labOrders, total, err := h.labService.ListOrders(page, limit)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"lab_orders": labOrders,
		"total":      total,
		"page":       page,
		"limit":      limit,
	})
}

// UpdateStatus updates the status of a lab order
// PATCH /api/v1/lab/orders/:id/status
func (h *LabHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid lab order id")
		return
	}

	var req UpdateLabStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	err = h.labService.UpdateStatus(id, req.Status)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	labOrder, _ := h.labService.GetOrder(id)
	utils.OK(c, labOrder)
}

// UploadReport uploads lab report
// POST /api/v1/lab/orders/:id/report
func (h *LabHandler) UploadReport(c *gin.Context) {
	idStr := c.Param("id")
	labOrderID, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid lab order id")
		return
	}

	staffIDStr := c.PostForm("staff_id")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid staff id")
		return
	}

	// Get uploaded file
	file, err := c.FormFile("report")
	if err != nil {
		utils.Fail(c, 400, "report file is required")
		return
	}

	// Save file temporarily
	tempFilePath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		utils.Fail(c, 500, "failed to save file")
		return
	}

	// Upload report
	fileURL, err := h.labService.UploadReport(labOrderID, staffID, tempFilePath, file.Filename)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"lab_order_id": labOrderID,
		"file_url":     fileURL,
		"message":      "report uploaded successfully",
	})
}

// GetPatientLabOrders retrieves lab orders for a patient
// GET /api/v1/lab/patients/:patient_id/orders
func (h *LabHandler) GetPatientLabOrders(c *gin.Context) {
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

	labOrders, total, err := h.labService.GetPatientLabOrders(patientID, page, limit)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"lab_orders": labOrders,
		"total":      total,
		"page":       page,
		"limit":      limit,
	})
}
