package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type NurseRepository interface {
	Create(nurse *model.Nurse) error
	GetByID(id uuid.UUID) (*model.Nurse, error)
	GetAll(page, limit int, search string) ([]model.Nurse, int64, error)
	Update(nurse *model.Nurse) error
	SoftDelete(id uuid.UUID) error
	GetByUserID(userID uuid.UUID) (*model.Nurse, error)
}

type nurseRepository struct{}

func NewNurseRepository() NurseRepository {
	return &nurseRepository{}
}

func (r *nurseRepository) Create(nurse *model.Nurse) error {
	return config.DB.Create(nurse).Error
}

func (r *nurseRepository) GetByID(id uuid.UUID) (*model.Nurse, error) {
	var nurse model.Nurse
	err := config.DB.Where("id = ?", id).
		Preload("User").
		First(&nurse).Error
	if err != nil {
		return nil, err
	}
	return &nurse, nil
}

func (r *nurseRepository) GetAll(page, limit int, search string) ([]model.Nurse, int64, error) {
	var nurses []model.Nurse
	var total int64

	query := config.DB.Model(&model.Nurse{})

	// Apply search filter
	if search != "" {
		query = query.Joins("JOIN users ON users.id = nurses.user_id").
			Where("users.first_name ILIKE ? OR users.last_name ILIKE ? OR nurses.department ILIKE ? OR nurses.license_number ILIKE ?",
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
		Find(&nurses).Error

	return nurses, total, err
}

func (r *nurseRepository) Update(nurse *model.Nurse) error {
	return config.DB.Model(nurse).
		Updates(nurse).Error
}

func (r *nurseRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.Nurse{}).
		Where("id = ?", id).
		Delete(&model.Nurse{}).Error
}

func (r *nurseRepository) GetByUserID(userID uuid.UUID) (*model.Nurse, error) {
	var nurse model.Nurse
	err := config.DB.Where("user_id = ?", userID).
		Preload("User").
		First(&nurse).Error
	if err != nil {
		return nil, err
	}
	return &nurse, nil
}
