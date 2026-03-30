package service

import (
	"context"
	"fmt"
	"time"

	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

// ============================================================
// FILTERING MODELS
// ============================================================

type AppointmentFilter struct {
	DoctorID       string
	PatientID      string
	Status         string
	Specialization string
	SearchTerm     string
	FromDate       *time.Time
	ToDate         *time.Time
	Page           int
	Limit          int
}

type LabOrderFilter struct {
	DoctorID  string
	PatientID string
	Status    string
	TestType  string
	Priority  string
	FromDate  *time.Time
	ToDate    *time.Time
	Page      int
	Limit     int
}

// ============================================================
// BULK OPERATION RESULTS
// ============================================================

type BulkOperationResult struct {
	UpdatedCount int                   `json:"updated_count"`
	FailedCount  int                   `json:"failed_count"`
	Details      []BulkOperationDetail `json:"details"`
}

type BulkOperationDetail struct {
	ID    string `json:"id"`
	Error string `json:"error,omitempty"`
	OK    bool   `json:"ok"`
}

// ============================================================
// APPOINTMENT SERVICE - FILTERING
// ============================================================

// GetAppointmentsFiltered retrieves appointments with advanced filtering
func (s *appointmentService) GetAppointmentsFiltered(ctx context.Context, filter *AppointmentFilter) ([]*model.Appointment, int64, error) {
	db := config.GetDB().WithContext(ctx)
	query := db

	// Apply filters
	if filter.DoctorID != "" {
		query = query.Where("doctor_id = ?", filter.DoctorID)
	}

	if filter.PatientID != "" {
		query = query.Where("patient_id = ?", filter.PatientID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Date range filtering
	if filter.FromDate != nil {
		fromDateStr := filter.FromDate.Format("2006-01-02")
		query = query.Where("DATE(appointment_date) >= ?", fromDateStr)
	}

	if filter.ToDate != nil {
		toDateStr := filter.ToDate.Format("2006-01-02")
		query = query.Where("DATE(appointment_date) <= ?", toDateStr)
	}

	// Specialization filtering (join with doctor table)
	if filter.Specialization != "" {
		query = query.Joins("JOIN doctors ON appointments.doctor_id = doctors.id").
			Where("doctors.specialization = ?", filter.Specialization)
	}

	// Patient search (by name or UHID)
	if filter.SearchTerm != "" {
		query = query.Joins("JOIN patients ON appointments.patient_id = patients.id").
			Where("CONCAT(patients.first_name, ' ', patients.last_name) LIKE ? OR patients.uhid LIKE ?",
				"%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%")
	}

	// Get total count
	var total int64
	query.Model(&model.Appointment{}).Count(&total)

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit)

	var appointments []*model.Appointment
	if err := query.Find(&appointments).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching filtered appointments: %w", err)
	}

	return appointments, total, nil
}

// ============================================================
// PRESCRIPTION SERVICE - BULK OPERATIONS
// ============================================================

// BulkUpdateStatus updates multiple prescriptions in a transaction
func (s *prescriptionService) BulkUpdateStatus(ctx context.Context, prescriptionIDs []uuid.UUID, status string, doctorID uuid.UUID) (*BulkOperationResult, error) {
	result := &BulkOperationResult{
		Details: make([]BulkOperationDetail, 0),
	}

	// Validate status
	validStatuses := map[string]bool{"approved": true, "rejected": true, "completed": true}
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	db := config.GetDB()
	tx := db.WithContext(ctx).Begin()

	for _, prescriptionID := range prescriptionIDs {
		detail := BulkOperationDetail{ID: prescriptionID.String(), OK: true}

		// Fetch prescription
		var prescription model.Prescription
		if err := tx.Where("id = ?", prescriptionID).First(&prescription).Error; err != nil {
			detail.Error = fmt.Sprintf("prescription not found: %v", err)
			detail.OK = false
			result.FailedCount++
			result.Details = append(result.Details, detail)
			continue
		}

		// Verify ownership (doctor can only update own prescriptions)
		if prescription.DoctorID != doctorID {
			detail.Error = "unauthorized - doctor can only update own prescriptions"
			detail.OK = false
			result.FailedCount++
			result.Details = append(result.Details, detail)
			continue
		}

		// Update status
		prescription.Status = status
		prescription.UpdatedAt = time.Now().UnixMilli()

		if err := tx.Save(&prescription).Error; err != nil {
			detail.Error = fmt.Sprintf("failed to update: %v", err)
			detail.OK = false
			result.FailedCount++
			result.Details = append(result.Details, detail)
			continue
		}

		result.UpdatedCount++
		result.Details = append(result.Details, detail)
	}

	// Commit transaction
	if result.FailedCount == 0 {
		if err := tx.Commit().Error; err != nil {
			return nil, fmt.Errorf("transaction commit failed: %w", err)
		}
	} else {
		tx.Rollback()
	}

	return result, nil
}

// ============================================================
// LAB SERVICE - FILTERING & BULK OPERATIONS
// ============================================================

// GetLabOrdersFiltered retrieves lab orders with advanced filtering
func (s *labService) GetLabOrdersFiltered(ctx context.Context, filter *LabOrderFilter) ([]*model.LabRequest, int64, error) {
	db := config.GetDB()
	query := db.WithContext(ctx)

	if filter.DoctorID != "" {
		query = query.Where("doctor_id = ?", filter.DoctorID)
	}

	if filter.PatientID != "" {
		query = query.Where("patient_id = ?", filter.PatientID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.TestType != "" {
		query = query.Where("test_type = ?", filter.TestType)
	}

	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}

	// Date range filtering
	if filter.FromDate != nil {
		fromDateStr := filter.FromDate.Format("2006-01-02")
		query = query.Where("DATE(created_at) >= ?", fromDateStr)
	}

	if filter.ToDate != nil {
		toDateStr := filter.ToDate.Format("2006-01-02")
		query = query.Where("DATE(created_at) <= ?", toDateStr)
	}

	// Get total count
	var total int64
	query.Model(&model.LabRequest{}).Count(&total)

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Preload("Patient").Preload("Doctor").Offset(offset).Limit(filter.Limit).Order("created_at DESC")

	var orders []*model.LabRequest
	if err := query.Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching filtered lab orders: %w", err)
	}

	return orders, total, nil
}

// BulkUpdateStatus marks multiple lab orders with same status
func (s *labService) BulkUpdateStatus(ctx context.Context, orderIDs []uuid.UUID, status string, labID uuid.UUID) (*BulkOperationResult, error) {
	result := &BulkOperationResult{
		Details: make([]BulkOperationDetail, 0),
	}

	validStatuses := map[string]bool{"completed": true, "rejected": true, "pending": true}
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	db := config.GetDB()
	tx := db.WithContext(ctx).Begin()

	for _, orderID := range orderIDs {
		detail := BulkOperationDetail{ID: orderID.String(), OK: true}

		var order model.LabRequest
		if err := tx.Where("id = ?", orderID).First(&order).Error; err != nil {
			detail.Error = fmt.Sprintf("order not found: %v", err)
			detail.OK = false
			result.FailedCount++
			result.Details = append(result.Details, detail)
			continue
		}

		order.Status = status
		order.UpdatedAt = time.Now().UnixMilli()

		if err := tx.Save(&order).Error; err != nil {
			detail.Error = fmt.Sprintf("failed to update: %v", err)
			detail.OK = false
			result.FailedCount++
			result.Details = append(result.Details, detail)
			continue
		}

		result.UpdatedCount++
		result.Details = append(result.Details, detail)
	}

	if result.FailedCount == 0 {
		if err := tx.Commit().Error; err != nil {
			return nil, fmt.Errorf("transaction commit failed: %w", err)
		}
	} else {
		tx.Rollback()
	}

	return result, nil
}

// ============================================================
// ADMISSION SERVICE - BULK DISCHARGE
// ============================================================
