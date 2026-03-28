package service

import (
	"errors"
	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
)

type HospitalService interface {
	CreateHospital(hospital *model.Hospital) error
	GetHospital(id uuid.UUID) (*model.Hospital, error)
	ListHospitals(page, limit int) ([]model.Hospital, int64, error)
	ListHospitalsByCity(city string, page, limit int) ([]model.Hospital, int64, error)
	UpdateHospital(hospital *model.Hospital) error
	DeleteHospital(id uuid.UUID) error
	SearchHospitals(query string, page, limit int) ([]model.Hospital, int64, error)
}

type hospitalService struct {
	repo repository.HospitalRepository
}

func NewHospitalService(repo repository.HospitalRepository) HospitalService {
	return &hospitalService{repo: repo}
}

func (s *hospitalService) CreateHospital(hospital *model.Hospital) error {
	if hospital.Name == "" {
		return errors.New("hospital name is required")
	}
	if hospital.Phone == "" {
		return errors.New("hospital phone is required")
	}
	if hospital.Email == "" {
		return errors.New("hospital email is required")
	}
	return s.repo.Create(hospital)
}

func (s *hospitalService) GetHospital(id uuid.UUID) (*model.Hospital, error) {
	return s.repo.GetByID(id)
}

func (s *hospitalService) ListHospitals(page, limit int) ([]model.Hospital, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.GetAll(page, limit)
}

func (s *hospitalService) ListHospitalsByCity(city string, page, limit int) ([]model.Hospital, int64, error) {
	if city == "" {
		return []model.Hospital{}, 0, errors.New("city is required")
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.GetByCity(city, page, limit)
}

func (s *hospitalService) UpdateHospital(hospital *model.Hospital) error {
	if hospital.ID == uuid.Nil {
		return errors.New("hospital id is required")
	}
	return s.repo.Update(hospital)
}

func (s *hospitalService) DeleteHospital(id uuid.UUID) error {
	return s.repo.SoftDelete(id)
}

func (s *hospitalService) SearchHospitals(query string, page, limit int) ([]model.Hospital, int64, error) {
	if query == "" {
		return []model.Hospital{}, 0, errors.New("search query is required")
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.Search(query, page, limit)
}
