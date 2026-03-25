package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type patientRepository struct{}

func NewPatientRepository() PatientRepository {
	return &patientRepository{}
}

func (r *patientRepository) Create(patient *model.Patient) error {
	return config.DB.Create(patient).Error
}

func (r *patientRepository) GetByID(id uuid.UUID) (*model.Patient, error) {
	var patient model.Patient
	err := config.DB.Where("id = ?", id).
		Preload("User").
		First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) GetAll(page, limit int, search string) ([]model.Patient, int64, error) {
	var patients []model.Patient
	var total int64

	query := config.DB.Model(&model.Patient{})

	// Apply search filter
	if search != "" {
		query = query.Joins("JOIN users ON users.id = patients.user_id").
			Where("users.first_name ILIKE ? OR users.last_name ILIKE ? OR users.phone ILIKE ? OR users.email ILIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("User").
		Offset(offset).
		Limit(limit).
		Find(&patients).Error

	return patients, total, err
}

func (r *patientRepository) Update(patient *model.Patient) error {
	return config.DB.Model(patient).
		Updates(patient).Error
}

func (r *patientRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.Patient{}).
		Where("id = ?", id).
		Delete(&model.Patient{}).Error
}

func (r *patientRepository) GetByUserID(userID uuid.UUID) (*model.Patient, error) {
	var patient model.Patient
	err := config.DB.Where("user_id = ?", userID).
		Preload("User").
		First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}
