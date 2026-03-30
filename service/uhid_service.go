package service

import (
	"fmt"
	"meet_sushruta/model"
	"meet_sushruta/repository"
	"time"

	"github.com/google/uuid"
)

type UHIDService interface {
	GenerateUHID(patientID uuid.UUID) (string, error)
	GetUHIDByPatientID(patientID uuid.UUID) (*model.UHID, error)
	ValidateUHID(uhidString string) (*model.UHID, error)
}

type uhidService struct {
	uhidRepo repository.UHIDRepository
}

func NewUHIDService(uhidRepo repository.UHIDRepository) UHIDService {
	return &uhidService{
		uhidRepo: uhidRepo,
	}
}

// GenerateUHID generates a unique UHID for a patient
// Format: MS-YYYY-XXXXX (MS = hospital code, YYYY = year, XXXXX = sequence)
func (s *uhidService) GenerateUHID(patientID uuid.UUID) (string, error) {
	// Check if UHID already exists for this patient
	existingUHID, err := s.uhidRepo.GetUHIDByPatientID(patientID)
	if err != nil {
		return "", fmt.Errorf("error checking existing UHID: %w", err)
	}

	if existingUHID != nil {
		return existingUHID.UHID, nil // Return existing UHID
	}

	// Hospital code (can be configured based on hospital)
	hospitalCode := "MS" // Meet Sushruta
	currentYear := time.Now().Year()

	// Get next sequence number
	seqNum, err := s.uhidRepo.GetNextSequenceNumber(hospitalCode)
	if err != nil {
		return "", fmt.Errorf("error getting sequence number: %w", err)
	}

	// Generate UHID string
	uhidString := fmt.Sprintf("%s-%d-%05d", hospitalCode, currentYear, seqNum)

	// Create UHID record
	uhid := &model.UHID{
		ID:             uuid.New(),
		PatientID:      patientID,
		UHID:           uhidString,
		HospitalCode:   hospitalCode,
		SequenceNumber: seqNum,
		IsActive:       true,
	}

	if err := s.uhidRepo.CreateUHID(uhid); err != nil {
		return "", fmt.Errorf("error creating UHID: %w", err)
	}

	return uhidString, nil
}

// GetUHIDByPatientID retrieves UHID for a patient
func (s *uhidService) GetUHIDByPatientID(patientID uuid.UUID) (*model.UHID, error) {
	uhid, err := s.uhidRepo.GetUHIDByPatientID(patientID)
	if err != nil {
		return nil, fmt.Errorf("error fetching UHID: %w", err)
	}
	return uhid, nil
}

// ValidateUHID validates if a UHID exists
func (s *uhidService) ValidateUHID(uhidString string) (*model.UHID, error) {
	uhid, err := s.uhidRepo.GetUHIDByUHIDString(uhidString)
	if err != nil {
		return nil, fmt.Errorf("error validating UHID: %w", err)
	}

	if uhid == nil {
		return nil, fmt.Errorf("UHID not found or inactive")
	}

	return uhid, nil
}
