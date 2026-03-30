package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// ============================================================
// DATA EXPORT SERVICE
// ============================================================

type ExportService struct {
	// dependencies can be added here
}

func NewExportService() *ExportService {
	return &ExportService{}
}

// ExportAppointmentsExcel generates Excel file with appointments
func (e *ExportService) ExportAppointmentsExcel(ctx context.Context, doctorID *uuid.UUID, fromDate, toDate *time.Time) (*bytes.Buffer, error) {
	db := config.GetDB().WithContext(ctx)
	query := db

	// Filter by doctor if specified
	if doctorID != nil {
		query = query.Where("doctor_id = ?", doctorID)
	}

	// Filter by date range if specified
	if fromDate != nil {
		query = query.Where("appointment_date >= ?", fromDate)
	}
	if toDate != nil {
		query = query.Where("appointment_date <= ?", toDate)
	}

	var appointments []model.Appointment
	if err := query.Find(&appointments).Error; err != nil {
		return nil, fmt.Errorf("error fetching appointments: %w", err)
	}

	// Create Excel file
	f := excelize.NewFile()
	defer f.Close()

	// Add headers
	headers := []string{"Appointment ID", "Patient Name", "Doctor Name", "Specialization", "Date", "Time", "Status", "Reason"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, header)
	}

	// Add data rows
	for idx, apt := range appointments {
		row := idx + 2
		cells := []interface{}{
			apt.ID.String(),
			apt.PatientID.String(), // Will be patient name in proper implementation
			apt.DoctorID.String(),
			"", // Specialization
			apt.AppointmentDate,
			apt.AppointmentTime,
			apt.Status,
			apt.Reason,
		}

		for colIdx, cell := range cells {
			cellName, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue("Sheet1", cellName, cell)
		}
	}

	// Auto-fit columns
	f.SetColWidth("Sheet1", "A", "H", 18)

	// Write to buffer
	buf := &bytes.Buffer{}
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("error writing excel file: %w", err)
	}

	return buf, nil
}

// ExportPrescriptionsExcel generates Excel file with prescriptions
func (e *ExportService) ExportPrescriptionsExcel(ctx context.Context, doctorID *uuid.UUID, fromDate, toDate *time.Time) (*bytes.Buffer, error) {
	db := config.GetDB().WithContext(ctx)
	query := db

	if doctorID != nil {
		query = query.Where("doctor_id = ?", doctorID)
	}

	if fromDate != nil {
		query = query.Where("created_at >= ?", fromDate)
	}
	if toDate != nil {
		query = query.Where("created_at <= ?", toDate)
	}

	var prescriptions []model.Prescription
	if err := query.Find(&prescriptions).Error; err != nil {
		return nil, fmt.Errorf("error fetching prescriptions: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	headers := []string{"Prescription ID", "Patient ID", "Doctor ID", "Appointment ID", "Status", "Notes", "Created At", "Updated At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, header)
	}

	for idx, presc := range prescriptions {
		row := idx + 2
		cells := []interface{}{
			presc.ID.String(),
			presc.PatientID.String(),
			presc.DoctorID.String(),
			presc.AppointmentID.String(),
			presc.Status,
			presc.Notes,
			time.UnixMilli(presc.CreatedAt).Format("2006-01-02 15:04"),
			time.UnixMilli(presc.UpdatedAt).Format("2006-01-02 15:04"),
		}

		for colIdx, cell := range cells {
			cellName, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue("Sheet1", cellName, cell)
		}
	}

	f.SetColWidth("Sheet1", "A", "H", 18)

	buf := &bytes.Buffer{}
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("error writing excel file: %w", err)
	}

	return buf, nil
}

// ExportAdmissionsExcel generates Excel file with admission records
func (e *ExportService) ExportAdmissionsExcel(ctx context.Context, doctorID *uuid.UUID, fromDate, toDate *time.Time) (*bytes.Buffer, error) {
	db := config.GetDB().WithContext(ctx)
	query := db

	if doctorID != nil {
		query = query.Where("doctor_id = ?", doctorID)
	}

	if fromDate != nil {
		query = query.Where("admission_date >= ?", fromDate)
	}
	if toDate != nil {
		query = query.Where("admission_date <= ?", toDate)
	}

	var admissions []model.AdmissionRecord
	if err := query.Find(&admissions).Error; err != nil {
		return nil, fmt.Errorf("error fetching admissions: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	headers := []string{"Admission ID", "Patient ID", "Doctor ID", "Bed ID", "Admission Date", "Discharge Date", "Reason", "Diagnosis", "Status", "Ward", "Room"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, header)
	}

	for idx, adm := range admissions {
		row := idx + 2
		dischargeDate := ""
		if adm.DischargeDate != nil {
			dischargeDate = *adm.DischargeDate
		}

		cells := []interface{}{
			adm.ID.String(),
			adm.PatientID.String(),
			adm.DoctorID.String(),
			adm.BedID.String(),
			adm.AdmissionDate,
			dischargeDate,
			adm.Reason,
			adm.Diagnosis,
			adm.Status,
			adm.Ward,
			adm.RoomNumber,
		}

		for colIdx, cell := range cells {
			cellName, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue("Sheet1", cellName, cell)
		}
	}

	f.SetColWidth("Sheet1", "A", "K", 18)

	buf := &bytes.Buffer{}
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("error writing excel file: %w", err)
	}

	return buf, nil
}

// ExportLabOrdersExcel generates Excel file with lab orders
func (e *ExportService) ExportLabOrdersExcel(ctx context.Context, doctorID *uuid.UUID, fromDate, toDate *time.Time) (*bytes.Buffer, error) {
	db := config.GetDB().WithContext(ctx)
	query := db

	if doctorID != nil {
		query = query.Where("doctor_id = ?", doctorID)
	}

	if fromDate != nil {
		query = query.Where("created_at >= ?", fromDate)
	}
	if toDate != nil {
		query = query.Where("created_at <= ?", toDate)
	}

	var labOrders []model.LabRequest
	if err := query.Find(&labOrders).Error; err != nil {
		return nil, fmt.Errorf("error fetching lab orders: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	headers := []string{"Order ID", "Patient ID", "Doctor ID", "Test Type", "Priority", "Status", "Created At", "Completed At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, header)
	}

	for idx, order := range labOrders {
		row := idx + 2
		completedAt := ""
		if order.CompletedAt != nil && *order.CompletedAt != "" {
			completedAt = *order.CompletedAt
		}

		cells := []interface{}{
			order.ID.String(),
			order.PatientID.String(),
			order.DoctorID.String(),
			order.TestType,
			order.Priority,
			order.Status,
			time.UnixMilli(order.CreatedAt).Format("2006-01-02 15:04"),
			completedAt,
		}

		for colIdx, cell := range cells {
			cellName, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue("Sheet1", cellName, cell)
		}
	}

	f.SetColWidth("Sheet1", "A", "H", 18)

	buf := &bytes.Buffer{}
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("error writing excel file: %w", err)
	}

	return buf, nil
}

// DischargeSummaryPDF represents data for discharge summary PDF
type DischargeSummaryData struct {
	PatientName     string
	UHID            string
	AdmissionDate   string
	DischargeDate   string
	Diagnosis       string
	Procedures      string
	Medications     string
	FollowUp        string
	DietRecommends  string
	ActivityLevel   string
	WarningSymptoms string
}

// GenerateDischargeSummary creates a simple text representation (basic implementation)
// For production, use a PDF library like gofpdf
func (e *ExportService) GenerateDischargeSummary(ctx context.Context, admissionID uuid.UUID) ([]byte, error) {
	db := config.GetDB().WithContext(ctx)

	var admission model.AdmissionRecord
	if err := db.Where("id = ?", admissionID).First(&admission).Error; err != nil {
		return nil, fmt.Errorf("admission not found: %w", err)
	}

	// Fetch discharge summary if exists
	var dischargeSummary model.DischargeSummary
	db.Where("admission_id = ?", admissionID).First(&dischargeSummary)

	// Create summary text
	summary := fmt.Sprintf(`
================================================================================
                          DISCHARGE SUMMARY
================================================================================

ADMISSION DETAILS:
  Admission ID: %s
  Admission Date: %s
  Discharge Date: %s
  
PATIENT INFORMATION:
  Patient ID: %s
  Doctor ID: %s
  Ward: %s
  
CLINICAL INFORMATION:
  Reason for Admission: %s
  Diagnosis: %s
  
DISCHARGE INFORMATION:
  Final Diagnosis: %s
  Procedures Performed: %s
  Discharge Medications: %s
  
PATIENT INSTRUCTIONS:
  Follow-up: %s
  Diet Recommendations: %s
  Activity Level: %s
  Warning Symptoms: %s
  
Patient Outcome: %s
Generated: %s

================================================================================
                    Please follow all instructions carefully
================================================================================
`,
		admission.ID.String(),
		admission.AdmissionDate,
		func() string {
			if admission.DischargeDate != nil {
				return *admission.DischargeDate
			}
			return "Not Yet Discharged"
		}(),
		admission.PatientID.String(),
		admission.DoctorID.String(),
		admission.Ward,
		admission.Reason,
		admission.Diagnosis,
		dischargeSummary.FinalDiagnosis,
		dischargeSummary.ProceduresPerformed,
		dischargeSummary.DischargeMedications,
		dischargeSummary.FollowUpInstructions,
		dischargeSummary.DietRecommendation,
		dischargeSummary.ActivityRecommendation,
		dischargeSummary.WarningSymptoms,
		dischargeSummary.Outcome,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return []byte(summary), nil
}
