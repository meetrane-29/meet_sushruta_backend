package handler

import (
	"fmt"
	"time"

	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ============================================================
// ADVANCED FILTERING ENHANCEMENTS
// ============================================================

// AppointmentFilterRequest handles advanced filtering
type AppointmentFilterRequest struct {
	DoctorID       string `form:"doctor_id"`
	PatientID      string `form:"patient_id"`
	Specialization string `form:"specialization"`
	Status         string `form:"status"`    // scheduled, in-progress, completed, cancelled
	SearchTerm     string `form:"search"`    // Patient name search
	FromDate       string `form:"from_date"` // YYYY-MM-DD
	ToDate         string `form:"to_date"`   // YYYY-MM-DD
	Page           int    `form:"page,default=1"`
	Limit          int    `form:"limit,default=20"`
}

// GetAppointmentsAdvanced retrieves appointments with advanced filtering
// GET /api/v1/appointments/advanced?doctor_id=xxx&status=scheduled&from_date=2026-01-01&to_date=2026-12-31&search=patient_name
func (h *AppointmentHandler) GetAppointmentsAdvanced(c *gin.Context) {
	var filter AppointmentFilterRequest

	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Fail(c, 400, "invalid query parameters")
		return
	}

	// Validate pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	// Validate date range
	var fromDate, toDate *time.Time
	if filter.FromDate != "" {
		if parsed, err := time.Parse("2006-01-02", filter.FromDate); err == nil {
			fromDate = &parsed
		}
	}
	if filter.ToDate != "" {
		if parsed, err := time.Parse("2006-01-02", filter.ToDate); err == nil {
			toDate = &parsed
		}
	}

	// Build query with filters
	query := c.Request.Context()
	appointments, total, err := h.appointmentService.GetAppointmentsFiltered(query, &service.AppointmentFilter{
		DoctorID:       filter.DoctorID,
		PatientID:      filter.PatientID,
		Status:         filter.Status,
		Specialization: filter.Specialization,
		SearchTerm:     filter.SearchTerm,
		FromDate:       fromDate,
		ToDate:         toDate,
		Page:           filter.Page,
		Limit:          filter.Limit,
	})
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error fetching appointments: %v", err))
		return
	}

	utils.OK(c, gin.H{
		"data":  appointments,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	})
}

// ============================================================
// BULK OPERATIONS - PRESCRIPTIONS
// ============================================================

type BulkPrescriptionStatusRequest struct {
	PrescriptionIDs []string `json:"prescription_ids" binding:"required,min=1"`
	Status          string   `json:"status" binding:"required"` // approved, rejected, completed
	Reason          string   `json:"reason"`
	DoctorID        string   `json:"doctor_id" binding:"required,uuid"`
}

// BulkUpdatePrescriptionStatus updates multiple prescriptions at once
// PUT /api/v1/prescriptions/bulk/status
func (h *PrescriptionHandler) BulkUpdatePrescriptionStatus(c *gin.Context) {
	var req BulkPrescriptionStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	// Validate status
	validStatuses := map[string]bool{"approved": true, "rejected": true, "completed": true}
	if !validStatuses[req.Status] {
		utils.Fail(c, 400, "invalid status - must be approved, rejected, or completed")
		return
	}

	// Parse UUIDs
	doctorID, err := uuid.Parse(req.DoctorID)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor_id")
		return
	}

	prescriptionUUIDs := make([]uuid.UUID, 0, len(req.PrescriptionIDs))
	for _, id := range req.PrescriptionIDs {
		parsed, err := uuid.Parse(id)
		if err != nil {
			utils.Fail(c, 400, fmt.Sprintf("invalid prescription id: %s", id))
			return
		}
		prescriptionUUIDs = append(prescriptionUUIDs, parsed)
	}

	// Update all prescriptions
	result, err := h.prescriptionService.BulkUpdateStatus(c.Request.Context(), prescriptionUUIDs, req.Status, doctorID)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error updating prescriptions: %v", err))
		return
	}

	utils.OK(c, gin.H{
		"updated": result.UpdatedCount,
		"failed":  result.FailedCount,
		"details": result.Details,
	})
}

// ============================================================
// BULK OPERATIONS - LAB TESTS
// ============================================================

type BulkLabStatusRequest struct {
	OrderIDs []string `json:"order_ids" binding:"required,min=1"`
	Status   string   `json:"status" binding:"required"` // completed, rejected, pending
	Results  string   `json:"results"`
	LabID    string   `json:"lab_id" binding:"required,uuid"`
}

// BulkUpdateLabStatus marks multiple lab tests as completed/rejected
// PUT /api/v1/lab/orders/bulk/status
func (h *LabHandler) BulkUpdateLabStatus(c *gin.Context) {
	var req BulkLabStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	// Validate status
	validStatuses := map[string]bool{"completed": true, "rejected": true, "pending": true}
	if !validStatuses[req.Status] {
		utils.Fail(c, 400, "invalid status - must be completed, rejected, or pending")
		return
	}

	// Parse UUIDs
	labID, err := uuid.Parse(req.LabID)
	if err != nil {
		utils.Fail(c, 400, "invalid lab_id")
		return
	}

	orderUUIDs := make([]uuid.UUID, 0, len(req.OrderIDs))
	for _, id := range req.OrderIDs {
		parsed, err := uuid.Parse(id)
		if err != nil {
			utils.Fail(c, 400, fmt.Sprintf("invalid order id: %s", id))
			return
		}
		orderUUIDs = append(orderUUIDs, parsed)
	}

	// Update all lab orders
	result, err := h.labService.BulkUpdateStatus(c.Request.Context(), orderUUIDs, req.Status, labID)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error updating lab orders: %v", err))
		return
	}

	utils.OK(c, gin.H{
		"updated": result.UpdatedCount,
		"failed":  result.FailedCount,
		"details": result.Details,
	})
}

// ============================================================
// BULK OPERATIONS - ADMISSIONS
// ============================================================

type BulkDischargeRequest struct {
	AdmissionIDs    []string `json:"admission_ids" binding:"required,min=1"`
	DoctorID        string   `json:"doctor_id" binding:"required,uuid"`
	DischargeReason string   `json:"discharge_reason"`
}

// BulkDischargePatients discharged multiple IPD patients in batch
// POST /api/v1/admissions/bulk/discharge
func (h *AdmissionHandler) BulkDischargePatients(c *gin.Context) {
	var req BulkDischargeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	// Parse UUIDs
	doctorID, err := uuid.Parse(req.DoctorID)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor_id")
		return
	}

	admissionUUIDs := make([]uuid.UUID, 0, len(req.AdmissionIDs))
	for _, id := range req.AdmissionIDs {
		parsed, err := uuid.Parse(id)
		if err != nil {
			utils.Fail(c, 400, fmt.Sprintf("invalid admission id: %s", id))
			return
		}
		admissionUUIDs = append(admissionUUIDs, parsed)
	}

	// Process bulk discharge
	dischargedCount := 0
	failedCount := 0
	details := make([]service.BulkOperationDetail, 0)

	for _, admissionID := range admissionUUIDs {
		detail := service.BulkOperationDetail{ID: admissionID.String(), OK: true}

		// Get admission
		admission, err := h.admissionRepo.GetByID(admissionID)
		if err != nil {
			detail.Error = fmt.Sprintf("admission not found: %v", err)
			detail.OK = false
			failedCount++
			details = append(details, detail)
			continue
		}

		// Verify doctor owns this admission
		if admission.DoctorID != doctorID {
			detail.Error = "unauthorized - doctor can only discharge own patients"
			detail.OK = false
			failedCount++
			details = append(details, detail)
			continue
		}

		// Check if already discharged
		if admission.Status == model.AdmissionDischarge {
			detail.Error = "admission already discharged"
			detail.OK = false
			failedCount++
			details = append(details, detail)
			continue
		}

		// Discharge patient
		now := time.Now()
		admissionToUpdate := admission
		admissionToUpdate.Status = model.AdmissionDischarge
		dischargeDate := now.Format("2006-01-02 15:04")
		admissionToUpdate.DischargeDate = &dischargeDate
		admissionToUpdate.UpdatedAt = now.UnixMilli()

		if err := h.admissionRepo.Update(admissionToUpdate); err != nil {
			detail.Error = fmt.Sprintf("failed to discharge: %v", err)
			detail.OK = false
			failedCount++
			details = append(details, detail)
			continue
		}

		dischargedCount++
		details = append(details, detail)
	}

	utils.OK(c, gin.H{
		"discharged": dischargedCount,
		"failed":     failedCount,
		"details":    details,
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
	})
}

// ============================================================
// ADVANCED FILTERING - LAB ORDERS
// ============================================================

type LabFilterRequest struct {
	DoctorID  string `form:"doctor_id"`
	PatientID string `form:"patient_id"`
	Status    string `form:"status"` // ordered, pending, completed, rejected
	TestType  string `form:"test_type"`
	FromDate  string `form:"from_date"` // YYYY-MM-DD
	ToDate    string `form:"to_date"`
	Priority  string `form:"priority"` // urgent, high, normal, low
	Page      int    `form:"page,default=1"`
	Limit     int    `form:"limit,default=20"`
}

// GetLabOrdersAdvanced retrieves lab orders with advanced filtering
// GET /api/v1/lab/orders/advanced?doctor_id=xxx&status=completed&from_date=2026-01-01
func (h *LabHandler) GetLabOrdersAdvanced(c *gin.Context) {
	var filter LabFilterRequest

	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Fail(c, 400, "invalid query parameters")
		return
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	// Validate date range
	var fromDate, toDate *time.Time
	if filter.FromDate != "" {
		if parsed, err := time.Parse("2006-01-02", filter.FromDate); err == nil {
			fromDate = &parsed
		}
	}
	if filter.ToDate != "" {
		if parsed, err := time.Parse("2006-01-02", filter.ToDate); err == nil {
			toDate = &parsed
		}
	}

	orders, total, err := h.labService.GetLabOrdersFiltered(c.Request.Context(), &service.LabOrderFilter{
		DoctorID:  filter.DoctorID,
		PatientID: filter.PatientID,
		Status:    filter.Status,
		TestType:  filter.TestType,
		Priority:  filter.Priority,
		FromDate:  fromDate,
		ToDate:    toDate,
		Page:      filter.Page,
		Limit:     filter.Limit,
	})
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error fetching lab orders: %v", err))
		return
	}

	utils.OK(c, gin.H{
		"data":  orders,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	})
}
