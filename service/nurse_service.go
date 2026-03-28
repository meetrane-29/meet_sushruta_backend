package service

import (
	"errors"
	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NurseService interface {
	CreateNurse(nurse *model.Nurse) error
	GetNurse(id uuid.UUID) (*model.Nurse, error)
	ListNurses(page, limit int, search string) ([]model.Nurse, int64, error)
	UpdateNurse(id uuid.UUID, updates map[string]interface{}) (*model.Nurse, error)
	DeleteNurse(id uuid.UUID) error
	GetNurseByUserID(userID uuid.UUID) (*model.Nurse, error)
}

type nurseService struct {
	nurseRepo repository.NurseRepository
}

func NewNurseService(nurseRepo repository.NurseRepository) NurseService {
	return &nurseService{
		nurseRepo: nurseRepo,
	}
}

func (s *nurseService) CreateNurse(nurse *model.Nurse) error {
	if nurse.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}
	if nurse.LicenseNumber == "" {
		return errors.New("license_number is required")
	}

	return s.nurseRepo.Create(nurse)
}

func (s *nurseService) GetNurse(id uuid.UUID) (*model.Nurse, error) {
	nurse, err := s.nurseRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("nurse not found")
		}
		return nil, err
	}
	return nurse, nil
}

func (s *nurseService) ListNurses(page, limit int, search string) ([]model.Nurse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.nurseRepo.GetAll(page, limit, search)
}

func (s *nurseService) UpdateNurse(id uuid.UUID, updates map[string]interface{}) (*model.Nurse, error) {
	nurse, err := s.nurseRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("nurse not found")
	}

	// Update allowed fields
	if department, ok := updates["department"].(string); ok {
		nurse.Department = department
	}
	if shift, ok := updates["shift"].(string); ok {
		nurse.Shift = shift
	}
	if certURL, ok := updates["certification_url"].(string); ok {
		nurse.CertificationURL = certURL
	}

	err = s.nurseRepo.Update(nurse)
	if err != nil {
		return nil, err
	}

	return nurse, nil
}

func (s *nurseService) DeleteNurse(id uuid.UUID) error {
	return s.nurseRepo.SoftDelete(id)
}

func (s *nurseService) GetNurseByUserID(userID uuid.UUID) (*model.Nurse, error) {
	nurse, err := s.nurseRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("nurse not found")
		}
		return nil, err
	}
	return nurse, nil
}
