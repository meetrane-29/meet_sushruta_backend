package handler

import (
	"fmt"
	"time"

	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ============================================================
// EXPORT HANDLER
// ============================================================

type ExportHandler struct {
	exportService *service.ExportService
}

func NewExportHandler(exportService *service.ExportService) *ExportHandler {
	return &ExportHandler{
		exportService: exportService,
	}
}

// ExportAppointmentsExcel exports appointments to Excel file
// GET /api/v1/exports/appointments/excel?doctor_id=xxx&from_date=2026-01-01&to_date=2026-12-31
func (h *ExportHandler) ExportAppointmentsExcel(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	var doctorID *uuid.UUID
	if doctorIDStr != "" {
		parsed, err := uuid.Parse(doctorIDStr)
		if err != nil {
			utils.Fail(c, 400, "invalid doctor_id format")
			return
		}
		doctorID = &parsed
	}

	var fromDate, toDate *time.Time
	if fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = &parsed
		}
	}
	if toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = &parsed
		}
	}

	buf, err := h.exportService.ExportAppointmentsExcel(c.Request.Context(), doctorID, fromDate, toDate)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error exporting appointments: %v", err))
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=appointments_%s.xlsx", time.Now().Format("20060102_150405")))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ExportPrescriptionsExcel exports prescriptions to Excel file
// GET /api/v1/exports/prescriptions/excel?doctor_id=xxx&from_date=2026-01-01&to_date=2026-12-31
func (h *ExportHandler) ExportPrescriptionsExcel(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	var doctorID *uuid.UUID
	if doctorIDStr != "" {
		parsed, err := uuid.Parse(doctorIDStr)
		if err != nil {
			utils.Fail(c, 400, "invalid doctor_id format")
			return
		}
		doctorID = &parsed
	}

	var fromDate, toDate *time.Time
	if fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = &parsed
		}
	}
	if toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = &parsed
		}
	}

	buf, err := h.exportService.ExportPrescriptionsExcel(c.Request.Context(), doctorID, fromDate, toDate)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error exporting prescriptions: %v", err))
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=prescriptions_%s.xlsx", time.Now().Format("20060102_150405")))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ExportAdmissionsExcel exports admission records to Excel file
// GET /api/v1/exports/admissions/excel?doctor_id=xxx&from_date=2026-01-01&to_date=2026-12-31
func (h *ExportHandler) ExportAdmissionsExcel(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	var doctorID *uuid.UUID
	if doctorIDStr != "" {
		parsed, err := uuid.Parse(doctorIDStr)
		if err != nil {
			utils.Fail(c, 400, "invalid doctor_id format")
			return
		}
		doctorID = &parsed
	}

	var fromDate, toDate *time.Time
	if fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = &parsed
		}
	}
	if toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = &parsed
		}
	}

	buf, err := h.exportService.ExportAdmissionsExcel(c.Request.Context(), doctorID, fromDate, toDate)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error exporting admissions: %v", err))
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=admissions_%s.xlsx", time.Now().Format("20060102_150405")))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ExportLabOrdersExcel exports lab orders to Excel file
// GET /api/v1/exports/lab-orders/excel?doctor_id=xxx&from_date=2026-01-01&to_date=2026-12-31
func (h *ExportHandler) ExportLabOrdersExcel(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	var doctorID *uuid.UUID
	if doctorIDStr != "" {
		parsed, err := uuid.Parse(doctorIDStr)
		if err != nil {
			utils.Fail(c, 400, "invalid doctor_id format")
			return
		}
		doctorID = &parsed
	}

	var fromDate, toDate *time.Time
	if fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = &parsed
		}
	}
	if toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = &parsed
		}
	}

	buf, err := h.exportService.ExportLabOrdersExcel(c.Request.Context(), doctorID, fromDate, toDate)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error exporting lab orders: %v", err))
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=lab_orders_%s.xlsx", time.Now().Format("20060102_150405")))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ExportDischargeSummaryText exports discharge summary as text file
// GET /api/v1/exports/discharge-summary/:admission_id
func (h *ExportHandler) ExportDischargeSummaryText(c *gin.Context) {
	admissionIDStr := c.Param("admission_id")
	admissionID, err := uuid.Parse(admissionIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid admission_id format")
		return
	}

	data, err := h.exportService.GenerateDischargeSummary(c.Request.Context(), admissionID)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error generating discharge summary: %v", err))
		return
	}

	c.Header("Content-Type", "text/plain")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=discharge_summary_%s.txt", admissionID.String()[:8]))
	c.Data(200, "text/plain", data)
}
