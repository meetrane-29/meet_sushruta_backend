package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"
	"time"

	"github.com/google/uuid"
)

type appointmentRepository struct{}

func NewAppointmentRepository() AppointmentRepository {
	return &appointmentRepository{}
}

func (r *appointmentRepository) Create(appointment *model.Appointment) error {
	return config.DB.Create(appointment).Error
}

func (r *appointmentRepository) GetByID(id uuid.UUID) (*model.Appointment, error) {
	var appointment model.Appointment
	err := config.DB.Where("id = ?", id).
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		First(&appointment).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *appointmentRepository) GetAll(page, limit int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	// Get total count with a separate query
	if err := config.DB.Model(&model.Appointment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with a fresh query
	err := config.DB.Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("appointment_date DESC, appointment_time DESC").
		Find(&appointments).Error

	return appointments, total, err
}

func (r *appointmentRepository) Update(appointment *model.Appointment) error {
	return config.DB.Model(appointment).
		Updates(appointment).Error
}

func (r *appointmentRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.Appointment{}).
		Where("id = ?", id).
		Delete(&model.Appointment{}).Error
}

func (r *appointmentRepository) GetByPatientID(patientID uuid.UUID, page, limit int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	// Get total count with a separate query
	if err := config.DB.Model(&model.Appointment{}).Where("patient_id = ?", patientID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with a fresh query
	err := config.DB.Where("patient_id = ?", patientID).
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("appointment_date DESC, appointment_time DESC").
		Find(&appointments).Error

	return appointments, total, err
}

func (r *appointmentRepository) GetByDoctorID(doctorID uuid.UUID, page, limit int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	// Get total count with a separate query
	if err := config.DB.Model(&model.Appointment{}).Where("doctor_id = ?", doctorID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with a fresh query
	err := config.DB.Where("doctor_id = ?", doctorID).
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("appointment_date DESC, appointment_time DESC").
		Find(&appointments).Error

	return appointments, total, err
}

func (r *appointmentRepository) GetByDateRange(doctorID uuid.UUID, startDate, endDate string) ([]model.Appointment, error) {
	var appointments []model.Appointment
	err := config.DB.Where("doctor_id = ? AND appointment_date >= ? AND appointment_date <= ?",
		doctorID, startDate, endDate).
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Order("appointment_date ASC, appointment_time ASC").
		Find(&appointments).Error

	return appointments, err
}

func (r *appointmentRepository) GetFutureAppointmentsByPatient(patientID uuid.UUID, page, limit int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	today := time.Now().Format("2006-01-02")

	// Get total count with a separate query
	if err := config.DB.Model(&model.Appointment{}).Where("patient_id = ? AND appointment_date >= ?", patientID, today).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with a fresh query
	err := config.DB.Where("patient_id = ? AND appointment_date >= ?", patientID, today).
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("appointment_date ASC, appointment_time ASC").
		Find(&appointments).Error

	return appointments, total, err
}

// GetAppointmentsForNext7Days returns appointments for last 7 days + next 7 days for a doctor
func (r *appointmentRepository) GetAppointmentsForNext7Days(doctorID uuid.UUID, page, limit int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	// Include last 7 days + next 7 days (total 14 days range centered on today)
	lastWeek := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	nextWeek := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	// Get total count with a separate query
	if err := config.DB.Model(&model.Appointment{}).Where("doctor_id = ? AND appointment_date >= ? AND appointment_date <= ?", doctorID, lastWeek, nextWeek).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with a fresh query
	err := config.DB.Where("doctor_id = ? AND appointment_date >= ? AND appointment_date <= ?", doctorID, lastWeek, nextWeek).
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("appointment_date ASC, appointment_time ASC").
		Find(&appointments).Error

	return appointments, total, err
}

// GetTodayAppointments returns all appointments for today (for receptionist)
func (r *appointmentRepository) GetTodayAppointments(page, limit int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	today := time.Now().Format("2006-01-02")

	// Get total count with a separate query
	if err := config.DB.Model(&model.Appointment{}).Where("appointment_date = ?", today).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with a fresh query
	err := config.DB.Where("appointment_date = ?", today).
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("appointment_time ASC").
		Find(&appointments).Error

	return appointments, total, err
}

// GetAllAppointmentsForNext7Days returns all appointments for next 7 days (for nurse)
func (r *appointmentRepository) GetAllAppointmentsForNext7Days(page, limit int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	today := time.Now().Format("2006-01-02")
	nextWeek := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	// Get total count with a separate query
	if err := config.DB.Model(&model.Appointment{}).Where("appointment_date >= ? AND appointment_date <= ?", today, nextWeek).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with a fresh query
	err := config.DB.Where("appointment_date >= ? AND appointment_date <= ?", today, nextWeek).
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("appointment_date ASC, appointment_time ASC").
		Find(&appointments).Error

	return appointments, total, err
}
