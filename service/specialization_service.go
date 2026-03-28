package service

import (
	"errors"
	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
)

type SpecializationService interface {
	CreateSpecialization(specialization *model.Specialization) error
	GetSpecialization(id uuid.UUID) (*model.Specialization, error)
	ListSpecializations(page, limit int) ([]model.Specialization, int64, error)
	ListMainSpecializations(page, limit int) ([]model.Specialization, int64, error)
	GetSubSpecializations(parentID uuid.UUID) ([]model.Specialization, error)
	UpdateSpecialization(specialization *model.Specialization) error
	DeleteSpecialization(id uuid.UUID) error
}

type specializationService struct {
	repo repository.SpecializationRepository
}

func NewSpecializationService(repo repository.SpecializationRepository) SpecializationService {
	return &specializationService{repo: repo}
}

func (s *specializationService) CreateSpecialization(specialization *model.Specialization) error {
	if specialization.Name == "" {
		return errors.New("specialization name is required")
	}
	return s.repo.Create(specialization)
}

func (s *specializationService) GetSpecialization(id uuid.UUID) (*model.Specialization, error) {
	return s.repo.GetByID(id)
}

func (s *specializationService) ListSpecializations(page, limit int) ([]model.Specialization, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.GetAll(page, limit)
}

func (s *specializationService) ListMainSpecializations(page, limit int) ([]model.Specialization, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.GetMainSpecializations(page, limit)
}

func (s *specializationService) GetSubSpecializations(parentID uuid.UUID) ([]model.Specialization, error) {
	return s.repo.GetSubSpecializations(parentID)
}

func (s *specializationService) UpdateSpecialization(specialization *model.Specialization) error {
	if specialization.ID == uuid.Nil {
		return errors.New("specialization id is required")
	}
	return s.repo.Update(specialization)
}

func (s *specializationService) DeleteSpecialization(id uuid.UUID) error {
	return s.repo.SoftDelete(id)
}
