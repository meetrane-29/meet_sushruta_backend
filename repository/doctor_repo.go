package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type doctorRepository struct{}

func NewDoctorRepository() DoctorRepository {
	return &doctorRepository{}
}

func (r *doctorRepository) Create(doctor *model.Doctor) error {
	return config.DB.Create(doctor).Error
}

func (r *doctorRepository) GetByID(id uuid.UUID) (*model.Doctor, error) {
	var doctor model.Doctor
	err := config.DB.Where("id = ?", id).
		Preload("User").
		First(&doctor).Error
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *doctorRepository) GetAll(page, limit int, search string) ([]model.Doctor, int64, error) {
	var doctors []model.Doctor
	var total int64

	query := config.DB.Model(&model.Doctor{})

	// Apply search filter
	if search != "" {
		query = query.Joins("JOIN users ON users.id = doctors.user_id").
			Where("users.first_name ILIKE ? OR users.last_name ILIKE ? OR doctors.specialization ILIKE ? OR doctors.department ILIKE ?",
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
		Find(&doctors).Error

	return doctors, total, err
}

func (r *doctorRepository) Update(doctor *model.Doctor) error {
	return config.DB.Model(doctor).
		Updates(doctor).Error
}

func (r *doctorRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.Doctor{}).
		Where("id = ?", id).
		Delete(&model.Doctor{}).Error
}

func (r *doctorRepository) GetByUserID(userID uuid.UUID) (*model.Doctor, error) {
	var doctor model.Doctor
	err := config.DB.Where("user_id = ?", userID).
		Preload("User").
		First(&doctor).Error
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}
