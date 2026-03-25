package service

import (
	"errors"
	"fmt"
	"time"

	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
)

type DoctorService interface {
	CreateDoctor(doctor *model.Doctor) error
	GetDoctor(id uuid.UUID) (*model.Doctor, error)
	ListDoctors(page, limit int, search string) ([]model.Doctor, int64, error)
	UpdateDoctor(doctor *model.Doctor) error
	SoftDeleteDoctor(id uuid.UUID) error
	GetAvailableSlots(doctorID uuid.UUID, date time.Time) ([]time.Time, error)
}

type doctorService struct {
	doctorRepo         repository.DoctorRepository
	doctorScheduleRepo repository.DoctorScheduleRepository
	appointmentRepo    repository.AppointmentRepository
}

func NewDoctorService(
	doctorRepo repository.DoctorRepository,
	doctorScheduleRepo repository.DoctorScheduleRepository,
	appointmentRepo repository.AppointmentRepository,
) DoctorService {
	return &doctorService{
		doctorRepo:         doctorRepo,
		doctorScheduleRepo: doctorScheduleRepo,
		appointmentRepo:    appointmentRepo,
	}
}

func (s *doctorService) CreateDoctor(doctor *model.Doctor) error {
	if doctor == nil {
		return errors.New("doctor cannot be nil")
	}

	if doctor.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}

	if doctor.Specialization == "" {
		return errors.New("specialization is required")
	}

	if doctor.LicenseNumber == "" {
		return errors.New("license_number is required")
	}

	err := s.doctorRepo.Create(doctor)
	if err != nil {
		return err
	}

	return nil
}

func (s *doctorService) GetDoctor(id uuid.UUID) (*model.Doctor, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid doctor id")
	}

	doctor, err := s.doctorRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("doctor not found")
	}

	return doctor, nil
}

func (s *doctorService) ListDoctors(page, limit int, search string) ([]model.Doctor, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100 // Max limit
	}

	doctors, total, err := s.doctorRepo.GetAll(page, limit, search)
	if err != nil {
		return nil, 0, err
	}

	return doctors, total, nil
}

func (s *doctorService) UpdateDoctor(doctor *model.Doctor) error {
	if doctor == nil {
		return errors.New("doctor cannot be nil")
	}

	if doctor.ID == uuid.Nil {
		return errors.New("doctor id is required")
	}

	// Verify doctor exists
	existingDoctor, err := s.doctorRepo.GetByID(doctor.ID)
	if err != nil {
		return errors.New("doctor not found")
	}

	// Update fields
	existingDoctor.Specialization = doctor.Specialization
	existingDoctor.Bio = doctor.Bio
	existingDoctor.Department = doctor.Department
	existingDoctor.ConsultationFee = doctor.ConsultationFee

	return s.doctorRepo.Update(existingDoctor)
}

func (s *doctorService) SoftDeleteDoctor(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid doctor id")
	}

	// Verify doctor exists
	_, err := s.doctorRepo.GetByID(id)
	if err != nil {
		return errors.New("doctor not found")
	}

	return s.doctorRepo.SoftDelete(id)
}

// GetAvailableSlots returns available time slots for a doctor on a given date
func (s *doctorService) GetAvailableSlots(doctorID uuid.UUID, date time.Time) ([]time.Time, error) {
	if doctorID == uuid.Nil {
		return nil, errors.New("invalid doctor id")
	}

	// Verify doctor exists
	_, err := s.doctorRepo.GetByID(doctorID)
	if err != nil {
		return nil, errors.New("doctor not found")
	}

	dayOfWeek := date.Format("Monday")

	// Get doctor's schedule for this day
	schedule, err := s.doctorScheduleRepo.GetByDoctorAndDay(doctorID, dayOfWeek)
	if err != nil || schedule == nil {
		return []time.Time{}, nil // No schedule for this day
	}

	if !schedule.IsAvailable {
		return []time.Time{}, nil
	}

	// Parse start and end times (format: HH:MM)
	startTime, err := time.Parse("15:04", schedule.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %v", err)
	}

	endTime, err := time.Parse("15:04", schedule.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end time format: %v", err)
	}

	dateStr := date.Format("2006-01-02")

	// Get booked appointments for this doctor on this date
	bookedAppointments, err := s.appointmentRepo.GetByDateRange(doctorID, dateStr, dateStr)
	if err != nil {
		return nil, err
	}

	// Create a map of booked times
	bookedTimes := make(map[string]bool)
	for _, apt := range bookedAppointments {
		if apt.Status != model.AppointmentCancelled {
			bookedTimes[apt.AppointmentTime] = true
		}
	}

	// Generate 30-minute slots
	slots := []time.Time{}
	current := startTime

	for current.Before(endTime) {
		slotStr := current.Format("15:04")
		if !bookedTimes[slotStr] {
			// Combine date and time for the result
			fullDateTime := time.Date(date.Year(), date.Month(), date.Day(), current.Hour(), current.Minute(), 0, 0, time.UTC)
			slots = append(slots, fullDateTime)
		}
		current = current.Add(30 * time.Minute)
	}

	return slots, nil
}
