package service

import (
	"fmt"
	"meet_sushruta/model"
	"meet_sushruta/repository"
	"time"

	"github.com/google/uuid"
)

type WaitingListService interface {
	AddToWaitingList(appointmentID, patientID, doctorID uuid.UUID) (*model.WaitingListEntry, error)
	GetDoctorWaitingList(doctorID uuid.UUID) ([]model.WaitingListEntry, error)
	CallNextPatient(doctorID uuid.UUID) (*model.WaitingListEntry, error)
	MarkPatientSeen(waitingListID uuid.UUID) error
	CompleteConsultation(waitingListID uuid.UUID) error
	CancelWaitingListEntry(waitingListID uuid.UUID) error
	GetPatientQueuePosition(patientID uuid.UUID) (int64, *model.WaitingListEntry, error)
	GetEstimatedWaitTime(doctorID uuid.UUID) int64
}

type waitingListService struct {
	waitingListRepo repository.WaitingListRepository
	appointmentRepo repository.AppointmentRepository
}

func NewWaitingListService(
	waitingListRepo repository.WaitingListRepository,
	appointmentRepo repository.AppointmentRepository,
) WaitingListService {
	return &waitingListService{
		waitingListRepo: waitingListRepo,
		appointmentRepo: appointmentRepo,
	}
}

// AddToWaitingList adds a patient to waiting list for doctor
func (s *waitingListService) AddToWaitingList(appointmentID, patientID, doctorID uuid.UUID) (*model.WaitingListEntry, error) {
	// Get next token number
	tokenNum, err := s.waitingListRepo.GetNextTokenNumber(doctorID)
	if err != nil {
		return nil, fmt.Errorf("error getting token number: %w", err)
	}

	now := time.Now().UnixMilli()
	entry := &model.WaitingListEntry{
		ID:            uuid.New(),
		AppointmentID: appointmentID,
		PatientID:     patientID,
		DoctorID:      doctorID,
		TokenNumber:   tokenNum,
		Status:        model.WaitingListWaiting,
		ArrivalTime:   &now,
	}

	if err := s.waitingListRepo.CreateWaitingListEntry(entry); err != nil {
		return nil, fmt.Errorf("error adding to waiting list: %w", err)
	}

	return entry, nil
}

// GetDoctorWaitingList retrieves waiting list for a doctor for today
func (s *waitingListService) GetDoctorWaitingList(doctorID uuid.UUID) ([]model.WaitingListEntry, error) {
	entries, err := s.waitingListRepo.GetWaitingListByDoctorIDToday(doctorID)
	if err != nil {
		return nil, fmt.Errorf("error fetching waiting list: %w", err)
	}
	return entries, nil
}

// CallNextPatient marks next patient as called
func (s *waitingListService) CallNextPatient(doctorID uuid.UUID) (*model.WaitingListEntry, error) {
	// Get waiting patients
	entries, err := s.waitingListRepo.GetWaitingListByDoctorIDToday(doctorID)
	if err != nil {
		return nil, fmt.Errorf("error fetching waiting list: %w", err)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no patients in waiting list")
	}

	// Call the first patient
	nextPatient := entries[0]
	now := time.Now().UnixMilli()
	nextPatient.CalledTime = &now

	err = s.waitingListRepo.UpdateWaitingListStatus(nextPatient.ID, model.WaitingListCalled)
	if err != nil {
		return nil, fmt.Errorf("error updating waiting list: %w", err)
	}

	return &nextPatient, nil
}

// MarkPatientSeen marks patient as seen/in consultation
func (s *waitingListService) MarkPatientSeen(waitingListID uuid.UUID) error {
	if err := s.waitingListRepo.UpdateWaitingListStatus(waitingListID, model.WaitingListSeen); err != nil {
		return fmt.Errorf("error marking patient as seen: %w", err)
	}
	return nil
}

// CompleteConsultation marks consultation as completed
func (s *waitingListService) CompleteConsultation(waitingListID uuid.UUID) error {
	if err := s.waitingListRepo.UpdateWaitingListStatus(waitingListID, model.WaitingListCompleted); err != nil {
		return fmt.Errorf("error completing consultation: %w", err)
	}
	return nil
}

// CancelWaitingListEntry cancels waiting list entry
func (s *waitingListService) CancelWaitingListEntry(waitingListID uuid.UUID) error {
	if err := s.waitingListRepo.UpdateWaitingListStatus(waitingListID, model.WaitingListCancelled); err != nil {
		return fmt.Errorf("error cancelling waiting list entry: %w", err)
	}
	return nil
}

// GetPatientQueuePosition gets patient's current queue position and details
func (s *waitingListService) GetPatientQueuePosition(patientID uuid.UUID) (int64, *model.WaitingListEntry, error) {
	entry, err := s.waitingListRepo.GetActiveWaitingListByPatientID(patientID)
	if err != nil {
		return 0, nil, fmt.Errorf("error fetching queue position: %w", err)
	}

	if entry == nil {
		return 0, nil, fmt.Errorf("patient not in queue")
	}

	return entry.TokenNumber, entry, nil
}

// GetEstimatedWaitTime calculates estimated wait time for a doctor
func (s *waitingListService) GetEstimatedWaitTime(doctorID uuid.UUID) int64 {
	// Average consultation time in minutes (5 minutes per consultation)
	avgConsultationTime := int64(5)

	entries, err := s.waitingListRepo.GetWaitingListByDoctorIDToday(doctorID)
	if err != nil {
		return 0
	}

	// Count waiting patients
	waitingCount := 0
	for _, entry := range entries {
		if entry.Status == model.WaitingListWaiting || entry.Status == model.WaitingListCalled {
			waitingCount++
		}
	}

	return int64(waitingCount) * avgConsultationTime
}
