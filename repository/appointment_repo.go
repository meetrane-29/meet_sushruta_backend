package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

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

	query := config.DB.Model(&model.Appointment{})

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Patient").
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

	query := config.DB.Where("patient_id = ?", patientID)

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Patient").
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

	query := config.DB.Where("doctor_id = ?", doctorID)

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Patient").
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
