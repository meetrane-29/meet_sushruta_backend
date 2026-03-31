package service

import (
	"context"
	"errors"
	"fmt"

	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
)

type AppointmentService interface {
	BookAppointment(appointment *model.Appointment) error
	GetAppointment(id uuid.UUID) (*model.Appointment, error)
	ListAppointments(page, limit int) ([]model.Appointment, int64, error)
	UpdateAppointmentStatus(id uuid.UUID, newStatus model.AppointmentStatus) error
	SoftDeleteAppointment(id uuid.UUID) error
	GetPatientAppointments(patientID uuid.UUID, page, limit int) ([]model.Appointment, int64, error)
	GetDoctorAppointments(doctorID uuid.UUID, page, limit int) ([]model.Appointment, int64, error)
	GetAppointmentsFiltered(ctx context.Context, filter *AppointmentFilter) ([]*model.Appointment, int64, error)
	GetTodayAppointments(page, limit int, doctorID ...string) ([]model.Appointment, int64, error)
	GetNext7DaysAppointments(page, limit int) ([]model.Appointment, int64, error)
}

type appointmentService struct {
	appointmentRepo repository.AppointmentRepository
	patientRepo     repository.PatientRepository
	doctorRepo      repository.DoctorRepository
}

func NewAppointmentService(
	appointmentRepo repository.AppointmentRepository,
	patientRepo repository.PatientRepository,
	doctorRepo repository.DoctorRepository,
) AppointmentService {
	return &appointmentService{
		appointmentRepo: appointmentRepo,
		patientRepo:     patientRepo,
		doctorRepo:      doctorRepo,
	}
}

// validTransitions defines the state machine for appointment status transitions
var validTransitions = map[model.AppointmentStatus][]model.AppointmentStatus{
	model.AppointmentPending:    {model.AppointmentConfirmed, model.AppointmentCancelled},
	model.AppointmentConfirmed:  {model.AppointmentInProgress, model.AppointmentCancelled},
	model.AppointmentInProgress: {model.AppointmentCompleted},
	model.AppointmentCompleted:  {},
	model.AppointmentCancelled:  {},
}

// BookAppointment creates a new appointment
func (s *appointmentService) BookAppointment(appointment *model.Appointment) error {
	if appointment == nil {
		return errors.New("appointment cannot be nil")
	}

	if appointment.PatientID == uuid.Nil {
		return errors.New("patient_id is required")
	}

	if appointment.DoctorID == uuid.Nil {
		return errors.New("doctor_id is required")
	}

	if appointment.AppointmentDate == "" {
		return errors.New("appointment_date is required")
	}

	if appointment.AppointmentTime == "" {
		return errors.New("appointment_time is required")
	}

	// Verify patient exists
	_, err := s.patientRepo.GetByID(appointment.PatientID)
	if err != nil {
		return errors.New("patient not found")
	}

	// Verify doctor exists
	_, err = s.doctorRepo.GetByID(appointment.DoctorID)
	if err != nil {
		return errors.New("doctor not found")
	}

	// Set default status to pending
	appointment.Status = model.AppointmentPending

	return s.appointmentRepo.Create(appointment)
}

// GetAppointment retrieves an appointment by ID
func (s *appointmentService) GetAppointment(id uuid.UUID) (*model.Appointment, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid appointment id")
	}

	appointment, err := s.appointmentRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	return appointment, nil
}

// ListAppointments retrieves appointments with pagination
func (s *appointmentService) ListAppointments(page, limit int) ([]model.Appointment, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100 // Max limit
	}

	appointments, total, err := s.appointmentRepo.GetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

// UpdateAppointmentStatus updates appointment status with state machine validation
func (s *appointmentService) UpdateAppointmentStatus(id uuid.UUID, newStatus model.AppointmentStatus) error {
	if id == uuid.Nil {
		return errors.New("invalid appointment id")
	}

	if newStatus == "" {
		return errors.New("status is required")
	}

	// Get current appointment
	appointment, err := s.appointmentRepo.GetByID(id)
	if err != nil {
		return errors.New("appointment not found")
	}

	currentStatus := appointment.Status

	// Validate transition
	allowedTransitions, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("invalid current status: %s", currentStatus)
	}

	isValidTransition := false
	for _, validStatus := range allowedTransitions {
		if validStatus == newStatus {
			isValidTransition = true
			break
		}
	}

	if !isValidTransition {
		return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
	}

	// Update status
	appointment.Status = newStatus
	return s.appointmentRepo.Update(appointment)
}

// SoftDeleteAppointment soft deletes an appointment
func (s *appointmentService) SoftDeleteAppointment(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid appointment id")
	}

	// Verify appointment exists
	_, err := s.appointmentRepo.GetByID(id)
	if err != nil {
		return errors.New("appointment not found")
	}

	return s.appointmentRepo.SoftDelete(id)
}

// GetPatientAppointments retrieves appointments for a patient
func (s *appointmentService) GetPatientAppointments(patientID uuid.UUID, page, limit int) ([]model.Appointment, int64, error) {
	if patientID == uuid.Nil {
		return nil, 0, errors.New("invalid patient id")
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	// Get appointments from today onwards only
	appointments, total, err := s.appointmentRepo.GetFutureAppointmentsByPatient(patientID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

// GetDoctorAppointments retrieves appointments for a doctor
func (s *appointmentService) GetDoctorAppointments(doctorID uuid.UUID, page, limit int) ([]model.Appointment, int64, error) {
	if doctorID == uuid.Nil {
		return nil, 0, errors.New("invalid doctor id")
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	// Get appointments for next 7 days only
	appointments, total, err := s.appointmentRepo.GetAppointmentsForNext7Days(doctorID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

// GetTodayAppointments returns all appointments for today (for receptionist)
func (s *appointmentService) GetTodayAppointments(page, limit int, doctorID ...string) ([]model.Appointment, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	// If doctorID is provided, filter by doctor
	if len(doctorID) > 0 && doctorID[0] != "" {
		fmt.Printf("[GetTodayAppointments] Filtering by doctor_id: %s\n", doctorID[0])
		appointments, total, err := s.appointmentRepo.GetTodayAppointmentsByDoctor(page, limit, doctorID[0])
		if err != nil {
			fmt.Printf("[GetTodayAppointments] Error filtering: %v\n", err)
			return nil, 0, err
		}
		fmt.Printf("[GetTodayAppointments] Found %d appointments for doctor\n", len(appointments))
		return appointments, total, nil
	}

	// Otherwise return all appointments for today
	fmt.Println("[GetTodayAppointments] No doctor_id provided, returning all appointments")
	appointments, total, err := s.appointmentRepo.GetTodayAppointments(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

// GetNext7DaysAppointments returns all appointments for next 7 days (for nurse)
func (s *appointmentService) GetNext7DaysAppointments(page, limit int) ([]model.Appointment, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	appointments, total, err := s.appointmentRepo.GetAllAppointmentsForNext7Days(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}
