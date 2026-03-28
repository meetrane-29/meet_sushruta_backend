package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type progressNoteRepository struct{}

func NewProgressNoteRepository() ProgressNoteRepository {
	return &progressNoteRepository{}
}

func (r *progressNoteRepository) Create(note *model.ProgressNote) error {
	return config.DB.Create(note).Error
}

func (r *progressNoteRepository) GetByID(id uuid.UUID) (*model.ProgressNote, error) {
	var note model.ProgressNote
	err := config.DB.Where("id = ?", id).
		Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *progressNoteRepository) GetAll(page, limit int) ([]model.ProgressNote, int64, error) {
	var notes []model.ProgressNote
	var total int64

	query := config.DB.Model(&model.ProgressNote{})

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
		Order("recorded_date DESC").
		Find(&notes).Error

	return notes, total, err
}

func (r *progressNoteRepository) GetByAdmissionID(admissionID uuid.UUID, page, limit int) ([]model.ProgressNote, int64, error) {
	var notes []model.ProgressNote
	var total int64

	query := config.DB.Where("admission_id = ?", admissionID)

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
		Order("recorded_date DESC").
		Find(&notes).Error

	return notes, total, err
}

func (r *progressNoteRepository) GetByPatientID(patientID uuid.UUID, page, limit int) ([]model.ProgressNote, int64, error) {
	var notes []model.ProgressNote
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
		Order("recorded_date DESC").
		Find(&notes).Error

	return notes, total, err
}

func (r *progressNoteRepository) GetByDoctorID(doctorID uuid.UUID, page, limit int) ([]model.ProgressNote, int64, error) {
	var notes []model.ProgressNote
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
		Order("recorded_date DESC").
		Find(&notes).Error

	return notes, total, err
}

func (r *progressNoteRepository) Update(note *model.ProgressNote) error {
	return config.DB.Model(note).
		Updates(note).Error
}

func (r *progressNoteRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.ProgressNote{}).
		Where("id = ?", id).
		Delete(&model.ProgressNote{}).Error
}
