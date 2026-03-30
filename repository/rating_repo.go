package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type ratingRepository struct{}

func NewRatingRepository() RatingRepository {
	return &ratingRepository{}
}

func (r *ratingRepository) Create(rating *model.Rating) error {
	return config.DB.Create(rating).Error
}

func (r *ratingRepository) GetByID(id uuid.UUID) (*model.Rating, error) {
	var rating model.Rating
	err := config.DB.Where("id = ?", id).
		Preload("Doctor").
		Preload("Doctor.User").
		Preload("Patient").
		Preload("Patient.User").
		First(&rating).Error
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *ratingRepository) GetAll(page, limit int) ([]model.Rating, int64, error) {
	var ratings []model.Rating
	var total int64

	query := config.DB.Model(&model.Rating{})

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Doctor").
		Preload("Doctor.User").
		Preload("Patient").
		Preload("Patient.User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&ratings).Error

	return ratings, total, err
}

func (r *ratingRepository) GetByDoctorID(doctorID uuid.UUID, page, limit int) ([]model.Rating, int64, error) {
	var ratings []model.Rating
	var total int64

	query := config.DB.Where("doctor_id = ?", doctorID)

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Doctor").
		Preload("Doctor.User").
		Preload("Patient").
		Preload("Patient.User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&ratings).Error

	return ratings, total, err
}

func (r *ratingRepository) GetByPatientID(patientID uuid.UUID, page, limit int) ([]model.Rating, int64, error) {
	var ratings []model.Rating
	var total int64

	query := config.DB.Where("patient_id = ?", patientID)

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Doctor").
		Preload("Doctor.User").
		Preload("Patient").
		Preload("Patient.User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&ratings).Error

	return ratings, total, err
}

func (r *ratingRepository) Update(rating *model.Rating) error {
	return config.DB.Model(rating).
		Updates(rating).Error
}

func (r *ratingRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.Rating{}).
		Where("id = ?", id).
		Delete(&model.Rating{}).Error
}
