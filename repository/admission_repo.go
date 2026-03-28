package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type admissionRepository struct{}

func NewAdmissionRepository() AdmissionRepository {
	return &admissionRepository{}
}

func (r *admissionRepository) Create(admission *model.AdmissionRecord) error {
	return config.DB.Create(admission).Error
}

func (r *admissionRepository) GetByID(id uuid.UUID) (*model.AdmissionRecord, error) {
	var admission model.AdmissionRecord
	err := config.DB.Where("id = ?", id).
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Preload("Bed").
		First(&admission).Error
	if err != nil {
		return nil, err
	}
	return &admission, nil
}

func (r *admissionRepository) GetAll(page, limit int) ([]model.AdmissionRecord, int64, error) {
	var admissions []model.AdmissionRecord
	var total int64

	query := config.DB.Model(&model.AdmissionRecord{})

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Preload("Bed").
		Offset(offset).
		Limit(limit).
		Order("admission_date DESC").
		Find(&admissions).Error

	return admissions, total, err
}

func (r *admissionRepository) GetByPatientID(patientID uuid.UUID, page, limit int) ([]model.AdmissionRecord, int64, error) {
	var admissions []model.AdmissionRecord
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
		Preload("Bed").
		Offset(offset).
		Limit(limit).
		Order("admission_date DESC").
		Find(&admissions).Error

	return admissions, total, err
}

func (r *admissionRepository) GetByDoctorID(doctorID uuid.UUID, page, limit int) ([]model.AdmissionRecord, int64, error) {
	var admissions []model.AdmissionRecord
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
		Preload("Bed").
		Offset(offset).
		Limit(limit).
		Order("admission_date DESC").
		Find(&admissions).Error

	return admissions, total, err
}

func (r *admissionRepository) GetActiveByPatientID(patientID uuid.UUID) (*model.AdmissionRecord, error) {
	var admission model.AdmissionRecord
	err := config.DB.Where("patient_id = ? AND status = ?", patientID, model.AdmissionActive).
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Preload("Bed").
		First(&admission).Error
	if err != nil {
		return nil, err
	}
	return &admission, nil
}

func (r *admissionRepository) GetActiveByDoctorID(doctorID uuid.UUID, page, limit int) ([]model.AdmissionRecord, int64, error) {
	var admissions []model.AdmissionRecord
	var total int64

	query := config.DB.Where("doctor_id = ? AND status = ?", doctorID, model.AdmissionActive)

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Preload("Bed").
		Offset(offset).
		Limit(limit).
		Order("admission_date DESC").
		Find(&admissions).Error

	return admissions, total, err
}

func (r *admissionRepository) Update(admission *model.AdmissionRecord) error {
	return config.DB.Model(admission).
		Updates(admission).Error
}

func (r *admissionRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.AdmissionRecord{}).
		Where("id = ?", id).
		Delete(&model.AdmissionRecord{}).Error
}
