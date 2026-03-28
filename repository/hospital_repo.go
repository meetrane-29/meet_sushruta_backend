package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type HospitalRepository interface {
	Create(hospital *model.Hospital) error
	GetByID(id uuid.UUID) (*model.Hospital, error)
	GetAll(page, limit int) ([]model.Hospital, int64, error)
	GetByCity(city string, page, limit int) ([]model.Hospital, int64, error)
	Update(hospital *model.Hospital) error
	SoftDelete(id uuid.UUID) error
	Search(query string, page, limit int) ([]model.Hospital, int64, error)
}

type hospitalRepository struct{}

func NewHospitalRepository() HospitalRepository {
	return &hospitalRepository{}
}

func (r *hospitalRepository) Create(hospital *model.Hospital) error {
	return config.DB.Create(hospital).Error
}

func (r *hospitalRepository) GetByID(id uuid.UUID) (*model.Hospital, error) {
	var hospital model.Hospital
	err := config.DB.Where("id = ? AND is_active = ?", id, true).First(&hospital).Error
	if err != nil {
		return nil, err
	}
	return &hospital, nil
}

func (r *hospitalRepository) GetAll(page, limit int) ([]model.Hospital, int64, error) {
	var hospitals []model.Hospital
	var total int64

	query := config.DB.Model(&model.Hospital{}).Where("is_active = ?", true)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).
		Limit(limit).
		Order("name ASC").
		Find(&hospitals).Error

	return hospitals, total, err
}

func (r *hospitalRepository) GetByCity(city string, page, limit int) ([]model.Hospital, int64, error) {
	var hospitals []model.Hospital
	var total int64

	query := config.DB.Model(&model.Hospital{}).
		Where("city ILIKE ? AND is_active = ?", city, true)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).
		Limit(limit).
		Order("name ASC").
		Find(&hospitals).Error

	return hospitals, total, err
}

func (r *hospitalRepository) Update(hospital *model.Hospital) error {
	return config.DB.Model(hospital).Updates(hospital).Error
}

func (r *hospitalRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.Hospital{}).
		Where("id = ?", id).
		Delete(&model.Hospital{}).Error
}

func (r *hospitalRepository) Search(query string, page, limit int) ([]model.Hospital, int64, error) {
	var hospitals []model.Hospital
	var total int64

	dbQuery := config.DB.Model(&model.Hospital{}).
		Where("(name ILIKE ? OR city ILIKE ? OR address ILIKE ?) AND is_active = ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%", true)
	dbQuery.Count(&total)

	offset := (page - 1) * limit
	err := dbQuery.Offset(offset).
		Limit(limit).
		Order("name ASC").
		Find(&hospitals).Error

	return hospitals, total, err
}
