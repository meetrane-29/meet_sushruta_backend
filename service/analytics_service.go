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
// ANALYTICS SERVICE
// ============================================================

type DoctorStatistics struct {
	DoctorID                  uuid.UUID `json:"doctor_id"`
	DoctorName                string    `json:"doctor_name"`
	TotalAppointments         int64     `json:"total_appointments"`
	CompletedAppts            int64     `json:"completed_appointments"`
	CancelledAppts            int64     `json:"cancelled_appointments"`
	UniquePatientsServed      int64     `json:"unique_patients_served"`
	TotalPrescriptions        int64     `json:"total_prescriptions"`
	LabOrdersIssued           int64     `json:"lab_orders_issued"`
	PatientsAdmitted          int64     `json:"patients_admitted"`
	AppointmentCompletionRate float64   `json:"appointment_completion_rate"`
	AverageRating             float64   `json:"average_rating"`
}

type AppointmentTrend struct {
	Date                  string `json:"date"`
	TotalAppointments     int64  `json:"total_appointments"`
	CompletedAppointments int64  `json:"completed_appointments"`
	CancelledAppointments int64  `json:"cancelled_appointments"`
	ScheduledAppointments int64  `json:"scheduled_appointments"`
}

type LabCompletionMetrics struct {
	TotalLabOrders        int64   `json:"total_lab_orders"`
	CompletedLabOrders    int64   `json:"completed_lab_orders"`
	PendingLabOrders      int64   `json:"pending_lab_orders"`
	RejectedLabOrders     int64   `json:"rejected_lab_orders"`
	CompletionRate        float64 `json:"completion_rate"`
	AverageCompletionTime int64   `json:"average_completion_time_minutes"`
}

type PatientOutcomeStatistics struct {
	TotalPatientsAdmitted  int64   `json:"total_patients_admitted"`
	TotalPatientsDischarge int64   `json:"total_patients_discharged"`
	PatientsInWard         int64   `json:"patients_in_ward"`
	DischargeRate          float64 `json:"discharge_rate"`
	PatientSatisfaction    float64 `json:"patient_satisfaction"`
	AverageLengthOfStay    float64 `json:"average_length_of_stay_days"`
}

type RevenueStatistics struct {
	DoctorID           uuid.UUID `json:"doctor_id"`
	DoctorName         string    `json:"doctor_name"`
	TotalRevenue       float64   `json:"total_revenue"`
	AppointmentRevenue float64   `json:"appointment_revenue"`
	LabRevenue         float64   `json:"lab_revenue"`
	ProcedureRevenue   float64   `json:"procedure_revenue"`
	OutstandingBills   float64   `json:"outstanding_bills"`
	CollectionRate     float64   `json:"collection_rate"`
}

type AnalyticsService struct{}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

// GetDoctorStatistics returns comprehensive statistics for a doctor
func (a *AnalyticsService) GetDoctorStatistics(ctx context.Context, doctorID uuid.UUID) (*DoctorStatistics, error) {
	db := config.GetDB().WithContext(ctx)

	stats := &DoctorStatistics{
		DoctorID: doctorID,
	}

	// Get doctor name
	var doctor model.Doctor
	if err := db.Preload("User").Where("id = ?", doctorID).First(&doctor).Error; err == nil {
		if doctor.User != nil {
			stats.DoctorName = doctor.User.FirstName + " " + doctor.User.LastName
		}
	}

	// Total appointments
	db.Model(&model.Appointment{}).Where("doctor_id = ?", doctorID).Count(&stats.TotalAppointments)

	// Completed appointments
	db.Model(&model.Appointment{}).Where("doctor_id = ? AND status = ?", doctorID, "completed").Count(&stats.CompletedAppts)

	// Cancelled appointments
	db.Model(&model.Appointment{}).Where("doctor_id = ? AND status = ?", doctorID, "cancelled").Count(&stats.CancelledAppts)

	// Unique patients
	db.Model(&model.Appointment{}).Where("doctor_id = ?", doctorID).
		Distinct("patient_id").Count(&stats.UniquePatientsServed)

	// Total prescriptions
	db.Model(&model.Prescription{}).Where("doctor_id = ?", doctorID).Count(&stats.TotalPrescriptions)

	// Lab orders issued
	db.Model(&model.LabRequest{}).Where("doctor_id = ?", doctorID).Count(&stats.LabOrdersIssued)

	// Patients admitted
	db.Model(&model.AdmissionRecord{}).Where("doctor_id = ?", doctorID).Count(&stats.PatientsAdmitted)

	// Calculate completion rate
	if stats.TotalAppointments > 0 {
		stats.AppointmentCompletionRate = float64(stats.CompletedAppts) / float64(stats.TotalAppointments) * 100
	}

	return stats, nil
}

// GetAppointmentTrends returns appointment trends for a date range
func (a *AnalyticsService) GetAppointmentTrends(ctx context.Context, doctorID *uuid.UUID, fromDate, toDate time.Time, groupBy string) ([]AppointmentTrend, error) {
	db := config.GetDB().WithContext(ctx)

	// groupBy can be "day", "week", "month"
	query := db.Table("appointments").
		Where("appointment_date BETWEEN ? AND ?", fromDate, toDate)

	if doctorID != nil {
		query = query.Where("doctor_id = ?", doctorID)
	}

	var trends []AppointmentTrend
	if err := query.Group("DATE_TRUNC('" + groupBy + "', appointment_date)").
		Order("DATE_TRUNC('" + groupBy + "', appointment_date)").
		Scan(&trends).Error; err != nil {
		return nil, fmt.Errorf("error fetching appointment trends: %w", err)
	}

	return trends, nil
}

// GetLabCompletionMetrics returns lab test completion statistics
func (a *AnalyticsService) GetLabCompletionMetrics(ctx context.Context, doctorID *uuid.UUID, timeRange string) (*LabCompletionMetrics, error) {
	db := config.GetDB().WithContext(ctx)

	metrics := &LabCompletionMetrics{}

	// Calculate date range
	now := time.Now()
	var startDate time.Time
	switch timeRange {
	case "7days":
		startDate = now.AddDate(0, 0, -7)
	case "30days":
		startDate = now.AddDate(0, -1, 0)
	case "90days":
		startDate = now.AddDate(0, -3, 0)
	default:
		startDate = now.AddDate(-1, 0, 0) // 1 year
	}

	query := db.Model(&model.LabRequest{}).Where("created_at >= ?", startDate)
	if doctorID != nil {
		query = query.Where("doctor_id = ?", doctorID)
	}

	// Total lab orders
	query.Count(&metrics.TotalLabOrders)

	// Completed orders
	query.Where("status = ?", "completed").Count(&metrics.CompletedLabOrders)

	// Pending orders
	query.Where("status = ?", "pending").Count(&metrics.PendingLabOrders)

	// Rejected orders
	query.Where("status = ?", "rejected").Count(&metrics.RejectedLabOrders)

	// Completion rate
	if metrics.TotalLabOrders > 0 {
		metrics.CompletionRate = float64(metrics.CompletedLabOrders) / float64(metrics.TotalLabOrders) * 100
	}

	// Average completion time (placeholder - would need actual time parsing from string)
	metrics.AverageCompletionTime = int64(60) // Default: 60 minutes

	return metrics, nil
}

// GetPatientOutcomeStatistics returns patient outcome metrics
func (a *AnalyticsService) GetPatientOutcomeStatistics(ctx context.Context, doctorID *uuid.UUID) (*PatientOutcomeStatistics, error) {
	db := config.GetDB().WithContext(ctx)

	stats := &PatientOutcomeStatistics{}

	query := db.Model(&model.AdmissionRecord{})
	if doctorID != nil {
		query = query.Where("doctor_id = ?", doctorID)
	}

	// Total admitted
	query.Count(&stats.TotalPatientsAdmitted)

	// Discharged
	query.Where("status = ?", model.AdmissionDischarge).Count(&stats.TotalPatientsDischarge)

	// Currently in ward
	query.Where("status = ?", model.AdmissionActive).Count(&stats.PatientsInWard)

	// Discharge rate
	if stats.TotalPatientsAdmitted > 0 {
		stats.DischargeRate = float64(stats.TotalPatientsDischarge) / float64(stats.TotalPatientsAdmitted) * 100
	}

	// Average length of stay
	var admissions []model.AdmissionRecord
	query.Find(&admissions)

	if len(admissions) > 0 {
		totalDays := 0.0
		for _, adm := range admissions {
			var admissionTime time.Time
			var dischargeTime time.Time

			var err error
			admissionTime, err = time.Parse("2006-01-02 15:04", adm.AdmissionDate)
			if err != nil {
				continue
			}

			if adm.DischargeDate != nil && *adm.DischargeDate != "" {
				dischargeTime, _ = time.Parse("2006-01-02 15:04", *adm.DischargeDate)
			} else {
				dischargeTime = time.Now()
			}

			totalDays += dischargeTime.Sub(admissionTime).Hours() / 24
		}
		stats.AverageLengthOfStay = totalDays / float64(len(admissions))
	}

	return stats, nil
}

// GetRevenueStatistics returns revenue metrics for a doctor
func (a *AnalyticsService) GetRevenueStatistics(ctx context.Context, doctorID uuid.UUID) (*RevenueStatistics, error) {
	db := config.GetDB().WithContext(ctx)

	rev := &RevenueStatistics{
		DoctorID: doctorID,
	}

	// Get doctor name
	var doctor model.Doctor
	if err := db.Preload("User").Where("id = ?", doctorID).First(&doctor).Error; err == nil {
		if doctor.User != nil {
			rev.DoctorName = doctor.User.FirstName + " " + doctor.User.LastName
		}
	}

	// Get all bills for this doctor
	var bills []model.Bill
	db.Where("doctor_id = ?", doctorID).Find(&bills)

	for _, bill := range bills {
		rev.TotalRevenue += bill.TotalAmount
		if bill.PaymentStatus != "paid" {
			rev.OutstandingBills += bill.TotalAmount
		}
	}

	// Collection rate
	if rev.TotalRevenue > 0 {
		rev.CollectionRate = ((rev.TotalRevenue - rev.OutstandingBills) / rev.TotalRevenue) * 100
	}

	// Get appointment revenue (if consultation fees are tracked)
	var appointmentBills int64
	db.Model(&model.Bill{}).Where("doctor_id = ? AND bill_type = ?", doctorID, "consultation").Count(&appointmentBills)
	rev.AppointmentRevenue = float64(appointmentBills) * 500 // Assuming 500 per consultation

	// Get lab revenue
	var labBills int64
	db.Model(&model.Bill{}).Where("doctor_id = ? AND bill_type = ?", doctorID, "lab").Count(&labBills)
	rev.LabRevenue = float64(labBills) * 1000 // Assuming 1000 per lab test

	return rev, nil
}

// GetDoctorRankings returns doctors ranked by specified metric
func (a *AnalyticsService) GetDoctorRankings(ctx context.Context, metric string, limit int) ([]DoctorStatistics, error) {
	db := config.GetDB().WithContext(ctx)

	var doctors []model.Doctor
	if err := db.Limit(limit).Find(&doctors).Error; err != nil {
		return nil, fmt.Errorf("error fetching doctors: %w", err)
	}

	var stats []DoctorStatistics
	for _, doc := range doctors {
		docStats, err := a.GetDoctorStatistics(ctx, doc.ID)
		if err != nil {
			continue
		}
		stats = append(stats, *docStats)
	}

	// Sort by metric (basic sorting - can be improved with custom comparator)
	// Just returning unsorted for now
	if len(stats) > limit {
		stats = stats[:limit]
	}

	return stats, nil
}
