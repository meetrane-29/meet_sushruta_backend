package repository

import (
	"errors"
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VitalsRepository interface {
	CreateVitals(vitals *model.Vitals) error
	GetVitalsByID(id uuid.UUID) (*model.Vitals, error)
	GetLatestVitalsByPatientID(patientID uuid.UUID) (*model.Vitals, error)
	GetVitalsHistoryByPatientID(patientID uuid.UUID, limit int, offset int) ([]model.Vitals, error)
	UpdateVitals(vitals *model.Vitals) error
	DeleteVitals(id uuid.UUID) error
}

type vitalsRepository struct{}

func NewVitalsRepository() VitalsRepository {
	return &vitalsRepository{}
}

// CreateVitals creates a new vitals record
func (r *vitalsRepository) CreateVitals(vitals *model.Vitals) error {
	if vitals.ID == uuid.Nil {
		vitals.ID = uuid.New()
	}

	result := config.DB.Create(vitals)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetVitalsByID retrieves vitals by ID
func (r *vitalsRepository) GetVitalsByID(id uuid.UUID) (*model.Vitals, error) {
	var vitals model.Vitals
	result := config.DB.
		Preload("Patient").
		Where("id = ?", id).
		First(&vitals)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("vitals record not found")
		}
		return nil, result.Error
	}

	return &vitals, nil
}

// GetLatestVitalsByPatientID retrieves the most recent vitals record for a patient
func (r *vitalsRepository) GetLatestVitalsByPatientID(patientID uuid.UUID) (*model.Vitals, error) {
	var vitals model.Vitals
	result := config.DB.
		Preload("Patient").
		Where("patient_id = ?", patientID).
		Order("recorded_at DESC").
		First(&vitals)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil without error if no records exist
		}
		return nil, result.Error
	}

	return &vitals, nil
}

// GetVitalsHistoryByPatientID retrieves vitals history for a patient with pagination
func (r *vitalsRepository) GetVitalsHistoryByPatientID(patientID uuid.UUID, limit int, offset int) ([]model.Vitals, error) {
	var vitals []model.Vitals

	// Default pagination limits
	if limit == 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	result := config.DB.
		Preload("Patient").
		Where("patient_id = ?", patientID).
		Order("recorded_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&vitals)

	if result.Error != nil {
		return nil, result.Error
	}

	return vitals, nil
}

// UpdateVitals updates an existing vitals record
func (r *vitalsRepository) UpdateVitals(vitals *model.Vitals) error {
	result := config.DB.Model(&model.Vitals{}).
		Where("id = ?", vitals.ID).
		Updates(vitals)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("vitals record not found")
	}

	return nil
}

// DeleteVitals deletes a vitals record
func (r *vitalsRepository) DeleteVitals(id uuid.UUID) error {
	result := config.DB.Where("id = ?", id).Delete(&model.Vitals{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("vitals record not found")
	}

	return nil
}
