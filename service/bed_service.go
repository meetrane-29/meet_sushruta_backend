package service

import (
	"fmt"
	"time"

	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
)

// BedService handles bed management operations
type BedService interface {
	// Admission and Discharge
	AdmitPatient(bedID, patientID uuid.UUID, notes string) error
	DischargePatient(bedID uuid.UUID, dischargeNotes string) error
	GetPatientBed(patientID uuid.UUID) (*model.Bed, error)

	// Bed queries
	GetAvailableBeds(bedType string) ([]model.Bed, error)
	GetBedsByWard(ward string) ([]model.Bed, error)
	GetBedsByStatus(status string) ([]model.Bed, error)

	// Bed status management
	UpdateBedStatus(bedID uuid.UUID, status string) error
	GetWardStats(ward string) (map[string]interface{}, error)
	GetBedStats() (map[string]interface{}, error)
}

type bedService struct {
	bedRepository repository.BedRepository
}

// NewBedService creates a new bed service instance
func NewBedService(bedRepository repository.BedRepository) BedService {
	return &bedService{
		bedRepository: bedRepository,
	}
}

// AdmitPatient admits a patient to a bed
// Sets bed status to occupied, assigns patient, and records admission time
func (s *bedService) AdmitPatient(bedID, patientID uuid.UUID, notes string) error {
	// Validate bedID and patientID
	if bedID == uuid.Nil || patientID == uuid.Nil {
		return fmt.Errorf("invalid bed or patient id")
	}

	// Get the bed to check if it's available
	bed, err := s.bedRepository.GetByID(bedID)
	if err != nil {
		return fmt.Errorf("bed not found: %w", err)
	}

	if bed == nil {
		return fmt.Errorf("bed does not exist")
	}

	// Check if bed is available
	if bed.Status != "available" {
		return fmt.Errorf("bed is not available (current status: %s)", bed.Status)
	}

	// Update bed with patient info
	bed.Status = "occupied"
	bed.PatientID = &patientID
	now := time.Now().Format("2006-01-02 15:04")
	bed.AdmittedAt = &now
	bed.DischargedAt = nil // Clear any previous discharge time
	if notes != "" {
		bed.MaintenanceNotes = notes
	}

	// Save changes
	if err := s.bedRepository.Update(bed); err != nil {
		return fmt.Errorf("failed to admit patient: %w", err)
	}

	return nil
}

// DischargePatient discharges a patient from a bed
// Sets bed status back to available and records discharge time
func (s *bedService) DischargePatient(bedID uuid.UUID, dischargeNotes string) error {
	// Get the bed
	bed, err := s.bedRepository.GetByID(bedID)
	if err != nil {
		return fmt.Errorf("bed not found: %w", err)
	}

	if bed == nil {
		return fmt.Errorf("bed does not exist")
	}

	// Check if bed is occupied
	if bed.Status != "occupied" {
		return fmt.Errorf("bed is not occupied (current status: %s)", bed.Status)
	}

	if bed.PatientID == nil {
		return fmt.Errorf("no patient assigned to this bed")
	}

	// Clear patient and update bed status
	bed.Status = "available"
	bed.PatientID = nil
	now := time.Now().Format("2006-01-02 15:04")
	bed.DischargedAt = &now
	if dischargeNotes != "" {
		bed.MaintenanceNotes = dischargeNotes
	}

	// Save changes
	if err := s.bedRepository.Update(bed); err != nil {
		return fmt.Errorf("failed to discharge patient: %w", err)
	}

	return nil
}

// GetPatientBed retrieves the current bed for a patient
func (s *bedService) GetPatientBed(patientID uuid.UUID) (*model.Bed, error) {
	if patientID == uuid.Nil {
		return nil, fmt.Errorf("invalid patient id")
	}

	// Query beds where patient is currently admitted (status = occupied)
	beds, err := s.bedRepository.GetByStatus("occupied")
	if err != nil {
		return nil, fmt.Errorf("failed to query beds: %w", err)
	}

	for _, bed := range beds {
		if bed.PatientID != nil && *bed.PatientID == patientID {
			return &bed, nil
		}
	}

	// Patient not in any bed
	return nil, nil
}

// GetAvailableBeds retrieves all available beds, optionally filtered by bed type
func (s *bedService) GetAvailableBeds(bedType string) ([]model.Bed, error) {
	if bedType == "" {
		// Get all available beds
		return s.bedRepository.GetByStatus("available")
	}

	// Get all available beds and filter by type
	beds, err := s.bedRepository.GetByStatus("available")
	if err != nil {
		return nil, err
	}

	var filteredBeds []model.Bed
	for _, bed := range beds {
		if bed.BedType == bedType {
			filteredBeds = append(filteredBeds, bed)
		}
	}

	return filteredBeds, nil
}

// GetBedsByWard retrieves all beds in a specific ward
func (s *bedService) GetBedsByWard(ward string) ([]model.Bed, error) {
	if ward == "" {
		return nil, fmt.Errorf("ward name cannot be empty")
	}

	return s.bedRepository.GetByWard(ward)
}

// GetBedsByStatus retrieves all beds with a specific status
func (s *bedService) GetBedsByStatus(status string) ([]model.Bed, error) {
	if status == "" {
		return nil, fmt.Errorf("status cannot be empty")
	}

	validStatuses := map[string]bool{
		"available":   true,
		"occupied":    true,
		"maintenance": true,
		"reserved":    true,
	}

	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid bed status: %s", status)
	}

	return s.bedRepository.GetByStatus(status)
}

// UpdateBedStatus updates the status of a bed
func (s *bedService) UpdateBedStatus(bedID uuid.UUID, status string) error {
	if bedID == uuid.Nil {
		return fmt.Errorf("invalid bed id")
	}

	// Validate status
	validStatuses := map[string]bool{
		"available":   true,
		"occupied":    true,
		"maintenance": true,
		"reserved":    true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid bed status: %s", status)
	}

	// Get the bed
	bed, err := s.bedRepository.GetByID(bedID)
	if err != nil {
		return fmt.Errorf("bed not found: %w", err)
	}

	if bed == nil {
		return fmt.Errorf("bed does not exist")
	}

	// Update status
	bed.Status = status

	// If changing to maintenance or reserved, clear patient if present
	if status == "maintenance" || status == "reserved" {
		bed.PatientID = nil
		now := time.Now().Format("2006-01-02 15:04")
		bed.DischargedAt = &now
	}

	// Save changes
	if err := s.bedRepository.Update(bed); err != nil {
		return fmt.Errorf("failed to update bed status: %w", err)
	}

	return nil
}

// GetWardStats retrieves statistics for a specific ward
func (s *bedService) GetWardStats(ward string) (map[string]interface{}, error) {
	if ward == "" {
		return nil, fmt.Errorf("ward name cannot be empty")
	}

	return s.bedRepository.GetWardStats()
}

// GetBedStats retrieves overall bed statistics
func (s *bedService) GetBedStats() (map[string]interface{}, error) {
	return s.bedRepository.GetBedStats()
}
