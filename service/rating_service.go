package service

import (
	"errors"

	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
)

type RatingService interface {
	CreateRating(rating *model.Rating) error
	GetRating(id uuid.UUID) (*model.Rating, error)
	ListRatings(page, limit int) ([]model.Rating, int64, error)
	GetDoctorRatings(doctorID uuid.UUID, page, limit int) ([]model.Rating, int64, error)
	GetPatientRatings(patientID uuid.UUID, page, limit int) ([]model.Rating, int64, error)
	UpdateRating(rating *model.Rating) error
	DeleteRating(id uuid.UUID) error
}

type ratingService struct {
	ratingRepo  repository.RatingRepository
	patientRepo repository.PatientRepository
	doctorRepo  repository.DoctorRepository
}

func NewRatingService(
	ratingRepo repository.RatingRepository,
	patientRepo repository.PatientRepository,
	doctorRepo repository.DoctorRepository,
) RatingService {
	return &ratingService{
		ratingRepo:  ratingRepo,
		patientRepo: patientRepo,
		doctorRepo:  doctorRepo,
	}
}

// CreateRating creates a new rating
func (s *ratingService) CreateRating(rating *model.Rating) error {
	if rating == nil {
		return errors.New("rating cannot be nil")
	}

	if rating.DoctorID == uuid.Nil {
		return errors.New("doctor_id is required")
	}

	if rating.PatientID == uuid.Nil {
		return errors.New("patient_id is required")
	}

	// Validate rating values (1-5 scale)
	if rating.Rating < 1 || rating.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	if rating.Professionalism < 1 || rating.Professionalism > 5 {
		return errors.New("professionalism rating must be between 1 and 5")
	}

	if rating.Communication < 1 || rating.Communication > 5 {
		return errors.New("communication rating must be between 1 and 5")
	}

	if rating.Punctuality < 1 || rating.Punctuality > 5 {
		return errors.New("punctuality rating must be between 1 and 5")
	}

	if rating.Cleanliness < 1 || rating.Cleanliness > 5 {
		return errors.New("cleanliness rating must be between 1 and 5")
	}

	// Verify doctor exists
	_, err := s.doctorRepo.GetByID(rating.DoctorID)
	if err != nil {
		return errors.New("doctor not found")
	}

	// Verify patient exists
	_, err = s.patientRepo.GetByID(rating.PatientID)
	if err != nil {
		return errors.New("patient not found")
	}

	return s.ratingRepo.Create(rating)
}

// GetRating retrieves a rating by ID
func (s *ratingService) GetRating(id uuid.UUID) (*model.Rating, error) {
	if id == uuid.Nil {
		return nil, errors.New("rating id is required")
	}

	return s.ratingRepo.GetByID(id)
}

// ListRatings retrieves all ratings with pagination
func (s *ratingService) ListRatings(page, limit int) ([]model.Rating, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.ratingRepo.GetAll(page, limit)
}

// GetDoctorRatings retrieves all ratings for a specific doctor
func (s *ratingService) GetDoctorRatings(doctorID uuid.UUID, page, limit int) ([]model.Rating, int64, error) {
	if doctorID == uuid.Nil {
		return nil, 0, errors.New("doctor_id is required")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Verify doctor exists
	_, err := s.doctorRepo.GetByID(doctorID)
	if err != nil {
		return nil, 0, errors.New("doctor not found")
	}

	return s.ratingRepo.GetByDoctorID(doctorID, page, limit)
}

// GetPatientRatings retrieves all ratings by a specific patient
func (s *ratingService) GetPatientRatings(patientID uuid.UUID, page, limit int) ([]model.Rating, int64, error) {
	if patientID == uuid.Nil {
		return nil, 0, errors.New("patient_id is required")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Verify patient exists
	_, err := s.patientRepo.GetByID(patientID)
	if err != nil {
		return nil, 0, errors.New("patient not found")
	}

	return s.ratingRepo.GetByPatientID(patientID, page, limit)
}

// UpdateRating updates an existing rating
func (s *ratingService) UpdateRating(rating *model.Rating) error {
	if rating == nil {
		return errors.New("rating cannot be nil")
	}

	if rating.ID == uuid.Nil {
		return errors.New("rating id is required")
	}

	// Validate rating values (1-5 scale)
	if rating.Rating < 1 || rating.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	if rating.Professionalism < 1 || rating.Professionalism > 5 {
		return errors.New("professionalism rating must be between 1 and 5")
	}

	if rating.Communication < 1 || rating.Communication > 5 {
		return errors.New("communication rating must be between 1 and 5")
	}

	if rating.Punctuality < 1 || rating.Punctuality > 5 {
		return errors.New("punctuality rating must be between 1 and 5")
	}

	if rating.Cleanliness < 1 || rating.Cleanliness > 5 {
		return errors.New("cleanliness rating must be between 1 and 5")
	}

	return s.ratingRepo.Update(rating)
}

// DeleteRating soft deletes a rating
func (s *ratingService) DeleteRating(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("rating id is required")
	}

	return s.ratingRepo.SoftDelete(id)
}
