package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type SpecializationRepository interface {
	Create(specialization *model.Specialization) error
	GetByID(id uuid.UUID) (*model.Specialization, error)
	GetAll(page, limit int) ([]model.Specialization, int64, error)
	GetByName(name string) (*model.Specialization, error)
	GetMainSpecializations(page, limit int) ([]model.Specialization, int64, error)
	GetSubSpecializations(parentID uuid.UUID) ([]model.Specialization, error)
	Update(specialization *model.Specialization) error
	SoftDelete(id uuid.UUID) error
}

type specializationRepository struct{}

func NewSpecializationRepository() SpecializationRepository {
	return &specializationRepository{}
}

func (r *specializationRepository) Create(specialization *model.Specialization) error {
	return config.DB.Create(specialization).Error
}

func (r *specializationRepository) GetByID(id uuid.UUID) (*model.Specialization, error) {
	var specialization model.Specialization
	err := config.DB.Where("id = ? AND is_active = ?", id, true).First(&specialization).Error
	if err != nil {
		return nil, err
	}
	return &specialization, nil
}

func (r *specializationRepository) GetAll(page, limit int) ([]model.Specialization, int64, error) {
	var specializations []model.Specialization
	var total int64

	query := config.DB.Model(&model.Specialization{}).Where("is_active = ?", true)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).
		Limit(limit).
		Order("name ASC").
		Find(&specializations).Error

	return specializations, total, err
}

func (r *specializationRepository) GetByName(name string) (*model.Specialization, error) {
	var specialization model.Specialization
	err := config.DB.Where("name ILIKE ? AND is_active = ?", name, true).First(&specialization).Error
	if err != nil {
		return nil, err
	}
	return &specialization, nil
}

// GetMainSpecializations retrieves specializations without parent (top-level)
func (r *specializationRepository) GetMainSpecializations(page, limit int) ([]model.Specialization, int64, error) {
	var specializations []model.Specialization
	var total int64

	query := config.DB.Model(&model.Specialization{}).
		Where("parent_id IS NULL AND is_active = ?", true)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).
		Limit(limit).
		Order("name ASC").
		Find(&specializations).Error

	return specializations, total, err
}

// GetSubSpecializations retrieves specializations with given parent
func (r *specializationRepository) GetSubSpecializations(parentID uuid.UUID) ([]model.Specialization, error) {
	var specializations []model.Specialization
	err := config.DB.Where("parent_id = ? AND is_active = ?", parentID, true).
		Order("name ASC").
		Find(&specializations).Error
	return specializations, err
}

func (r *specializationRepository) Update(specialization *model.Specialization) error {
	return config.DB.Model(specialization).Updates(specialization).Error
}

func (r *specializationRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.Specialization{}).
		Where("id = ?", id).
		Delete(&model.Specialization{}).Error
}
