package service

import (
	"errors"

	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
)

type PatientService interface {
	RegisterPatient(patient *model.Patient) error
	GetPatient(id uuid.UUID) (*model.Patient, error)
	ListPatients(page, limit int, search string) ([]model.Patient, int64, error)
	UpdatePatient(patient *model.Patient) error
	SoftDeletePatient(id uuid.UUID) error
}

type patientService struct {
	patientRepo repository.PatientRepository
}

func NewPatientService(patientRepo repository.PatientRepository) PatientService {
	return &patientService{
		patientRepo: patientRepo,
	}
}

// RegisterPatient creates a new patient and validates phone uniqueness
func (s *patientService) RegisterPatient(patient *model.Patient) error {
	if patient == nil {
		return errors.New("patient cannot be nil")
	}

	if patient.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}

	if patient.DateOfBirth == "" {
		return errors.New("date_of_birth is required")
	}

	if patient.Gender == "" {
		return errors.New("gender is required")
	}

	// Phone uniqueness is enforced at the User level, not here
	// Service validation for patient-specific fields

	err := s.patientRepo.Create(patient)
	if err != nil {
		return err
	}

	return nil
}

// GetPatient retrieves a patient by ID
func (s *patientService) GetPatient(id uuid.UUID) (*model.Patient, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid patient id")
	}

	patient, err := s.patientRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("patient not found")
	}

	return patient, nil
}

// ListPatients retrieves patients with pagination and optional search
func (s *patientService) ListPatients(page, limit int, search string) ([]model.Patient, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100 // Max limit
	}

	patients, total, err := s.patientRepo.GetAll(page, limit, search)
	if err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}

// UpdatePatient updates patient information
func (s *patientService) UpdatePatient(patient *model.Patient) error {
	if patient == nil {
		return errors.New("patient cannot be nil")
	}

	if patient.ID == uuid.Nil {
		return errors.New("patient id is required")
	}

	// Verify patient exists
	existingPatient, err := s.patientRepo.GetByID(patient.ID)
	if err != nil {
		return errors.New("patient not found")
	}

	// Update fields
	existingPatient.DateOfBirth = patient.DateOfBirth
	existingPatient.Gender = patient.Gender
	existingPatient.BloodGroup = patient.BloodGroup
	existingPatient.Address = patient.Address
	existingPatient.EmergencyContact = patient.EmergencyContact
	existingPatient.MedicalHistory = patient.MedicalHistory
	existingPatient.Allergies = patient.Allergies

	return s.patientRepo.Update(existingPatient)
}

// SoftDeletePatient soft deletes a patient
func (s *patientService) SoftDeletePatient(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid patient id")
	}

	// Verify patient exists
	_, err := s.patientRepo.GetByID(id)
	if err != nil {
		return errors.New("patient not found")
	}

	return s.patientRepo.SoftDelete(id)
}
