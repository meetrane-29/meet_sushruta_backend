package service

import (
	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
)

type VitalsService interface {
	CreateVitals(vitals *model.Vitals) error
	GetVitalsByID(id uuid.UUID) (*model.Vitals, error)
	GetLatestVitalsByPatientID(patientID uuid.UUID) (*model.Vitals, error)
	GetVitalsHistoryByPatientID(patientID uuid.UUID, limit int, offset int) ([]model.Vitals, error)
	GetVitalsByPatientID(patientID uuid.UUID) ([]model.Vitals, error)
	UpdateVitals(vitals *model.Vitals) error
	DeleteVitals(id uuid.UUID) error
}

type vitalsService struct {
	vitalsRepo repository.VitalsRepository
}

func NewVitalsService(vitalsRepo repository.VitalsRepository) VitalsService {
	return &vitalsService{
		vitalsRepo: vitalsRepo,
	}
}

// CreateVitals creates a new vitals record
func (s *vitalsService) CreateVitals(vitals *model.Vitals) error {
	return s.vitalsRepo.CreateVitals(vitals)
}

// GetVitalsByID retrieves vitals by ID
func (s *vitalsService) GetVitalsByID(id uuid.UUID) (*model.Vitals, error) {
	return s.vitalsRepo.GetVitalsByID(id)
}

// GetLatestVitalsByPatientID retrieves the most recent vitals for a patient
func (s *vitalsService) GetLatestVitalsByPatientID(patientID uuid.UUID) (*model.Vitals, error) {
	return s.vitalsRepo.GetLatestVitalsByPatientID(patientID)
}

// GetVitalsHistoryByPatientID retrieves vitals history with pagination
func (s *vitalsService) GetVitalsHistoryByPatientID(patientID uuid.UUID, limit int, offset int) ([]model.Vitals, error) {
	return s.vitalsRepo.GetVitalsHistoryByPatientID(patientID, limit, offset)
}

// GetVitalsByPatientID retrieves all vitals for a patient
func (s *vitalsService) GetVitalsByPatientID(patientID uuid.UUID) ([]model.Vitals, error) {
	// Get a large number of vitals (e.g., last 1000 records)
	// In practice, you might want to limit this or add date range filtering
	return s.vitalsRepo.GetVitalsHistoryByPatientID(patientID, 1000, 0)
}

// UpdateVitals updates an existing vitals record
func (s *vitalsService) UpdateVitals(vitals *model.Vitals) error {
	return s.vitalsRepo.UpdateVitals(vitals)
}

// DeleteVitals deletes a vitals record
func (s *vitalsService) DeleteVitals(id uuid.UUID) error {
	return s.vitalsRepo.DeleteVitals(id)
}
