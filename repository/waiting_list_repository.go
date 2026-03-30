package repository

import (
	"errors"
	"fmt"
	"meet_sushruta/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WaitingListRepository interface {
	CreateWaitingListEntry(entry *model.WaitingListEntry) error
	GetWaitingListByAppointmentID(appointmentID uuid.UUID) (*model.WaitingListEntry, error)
	GetWaitingListByDoctorIDToday(doctorID uuid.UUID) ([]model.WaitingListEntry, error)
	GetWaitingListByStatus(status string) ([]model.WaitingListEntry, error)
	UpdateWaitingListStatus(id uuid.UUID, status model.WaitingListStatus) error
	GetNextTokenNumber(doctorID uuid.UUID) (int64, error)
	DeleteWaitingListEntry(id uuid.UUID) error
	GetActiveWaitingListByPatientID(patientID uuid.UUID) (*model.WaitingListEntry, error)
}

type waitingListRepository struct {
	db *gorm.DB
}

func NewWaitingListRepository(db *gorm.DB) WaitingListRepository {
	return &waitingListRepository{db: db}
}

// CreateWaitingListEntry creates a new waiting list entry
func (r *waitingListRepository) CreateWaitingListEntry(entry *model.WaitingListEntry) error {
	if err := r.db.Create(entry).Error; err != nil {
		return fmt.Errorf("error creating waiting list entry: %w", err)
	}
	return nil
}

// GetWaitingListByAppointmentID retrieves waiting list entry by appointment ID
func (r *waitingListRepository) GetWaitingListByAppointmentID(appointmentID uuid.UUID) (*model.WaitingListEntry, error) {
	var entry model.WaitingListEntry
	if err := r.db.Where("appointment_id = ?", appointmentID).
		First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching waiting list entry: %w", err)
	}
	return &entry, nil
}

// GetWaitingListByDoctorIDToday retrieves today's waiting list for a doctor
func (r *waitingListRepository) GetWaitingListByDoctorIDToday(doctorID uuid.UUID) ([]model.WaitingListEntry, error) {
	var entries []model.WaitingListEntry
	today := getCurrentDateString()

	if err := r.db.Joins("JOIN appointments ON waiting_list_entries.appointment_id = appointments.id").
		Where("appointments.doctor_id = ? AND DATE(appointments.appointment_date) = ?", doctorID, today).
		Order("waiting_list_entries.token_number ASC").
		Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("error fetching waiting list: %w", err)
	}
	return entries, nil
}

// GetWaitingListByStatus retrieves waiting list entries by status
func (r *waitingListRepository) GetWaitingListByStatus(status string) ([]model.WaitingListEntry, error) {
	var entries []model.WaitingListEntry
	if err := r.db.Where("status = ?", status).
		Order("token_number ASC").
		Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("error fetching waiting list by status: %w", err)
	}
	return entries, nil
}

// UpdateWaitingListStatus updates the status of a waiting list entry
func (r *waitingListRepository) UpdateWaitingListStatus(id uuid.UUID, status model.WaitingListStatus) error {
	if err := r.db.Model(&model.WaitingListEntry{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("error updating waiting list status: %w", err)
	}
	return nil
}

// GetNextTokenNumber gets the next token number for a doctor's queue
func (r *waitingListRepository) GetNextTokenNumber(doctorID uuid.UUID) (int64, error) {
	var maxToken int64
	today := getCurrentDateString()

	if err := r.db.Model(&model.WaitingListEntry{}).
		Joins("JOIN appointments ON waiting_list_entries.appointment_id = appointments.id").
		Where("appointments.doctor_id = ? AND DATE(appointments.appointment_date) = ?", doctorID, today).
		Select("COALESCE(MAX(token_number), 0)").
		Scan(&maxToken).Error; err != nil {
		return 0, fmt.Errorf("error getting next token number: %w", err)
	}

	return maxToken + 1, nil
}

// DeleteWaitingListEntry deletes a waiting list entry
func (r *waitingListRepository) DeleteWaitingListEntry(id uuid.UUID) error {
	if err := r.db.Delete(&model.WaitingListEntry{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("error deleting waiting list entry: %w", err)
	}
	return nil
}

// GetActiveWaitingListByPatientID retrieves active waiting list entry for a patient
func (r *waitingListRepository) GetActiveWaitingListByPatientID(patientID uuid.UUID) (*model.WaitingListEntry, error) {
	var entry model.WaitingListEntry
	if err := r.db.Where("patient_id = ? AND status IN (?)", patientID, []string{"waiting", "called"}).
		First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching active waiting list entry: %w", err)
	}
	return &entry, nil
}

// Helper function to get current date string in format YYYY-MM-DD
func getCurrentDateString() string {
	// This will be implemented based on your timezone requirements
	return ""
}
