package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type dischargeSummaryRepository struct{}

func NewDischargeSummaryRepository() DischargeSummaryRepository {
	return &dischargeSummaryRepository{}
}

func (r *dischargeSummaryRepository) Create(summary *model.DischargeSummary) error {
	return config.DB.Create(summary).Error
}

func (r *dischargeSummaryRepository) GetByID(id uuid.UUID) (*model.DischargeSummary, error) {
	var summary model.DischargeSummary
	err := config.DB.Where("id = ?", id).
		Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		First(&summary).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *dischargeSummaryRepository) GetAll(page, limit int) ([]model.DischargeSummary, int64, error) {
	var summaries []model.DischargeSummary
	var total int64

	query := config.DB.Model(&model.DischargeSummary{})

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("discharge_date DESC").
		Find(&summaries).Error

	return summaries, total, err
}

func (r *dischargeSummaryRepository) GetByAdmissionID(admissionID uuid.UUID) (*model.DischargeSummary, error) {
	var summary model.DischargeSummary
	err := config.DB.Where("admission_id = ?", admissionID).
		Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		First(&summary).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *dischargeSummaryRepository) GetByPatientID(patientID uuid.UUID, page, limit int) ([]model.DischargeSummary, int64, error) {
	var summaries []model.DischargeSummary
	var total int64

	query := config.DB.Where("patient_id = ?", patientID)

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("discharge_date DESC").
		Find(&summaries).Error

	return summaries, total, err
}

func (r *dischargeSummaryRepository) GetByDoctorID(doctorID uuid.UUID, page, limit int) ([]model.DischargeSummary, int64, error) {
	var summaries []model.DischargeSummary
	var total int64

	query := config.DB.Where("doctor_id = ?", doctorID)

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("discharge_date DESC").
		Find(&summaries).Error

	return summaries, total, err
}

func (r *dischargeSummaryRepository) Update(summary *model.DischargeSummary) error {
	return config.DB.Model(summary).
		Updates(summary).Error
}

func (r *dischargeSummaryRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.DischargeSummary{}).
		Where("id = ?", id).
		Delete(&model.DischargeSummary{}).Error
}
