package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemInput struct {
	MedicineID   uuid.UUID `json:"medicine_id"`
	MedicineName string    `json:"medicine_name"`
	Dosage       string    `json:"dosage" binding:"required"`
	Frequency    string    `json:"frequency" binding:"required"`
	Duration     string    `json:"duration" binding:"required"`
	TimeOfDay    string    `json:"time_of_day"`
	Quantity     int       `json:"quantity" binding:"required"`
	Instructions string    `json:"instructions"`
}

type PrescriptionService interface {
	CreatePrescription(appointmentID, doctorID uuid.UUID, items []ItemInput, notes string) (*model.Prescription, error)
	GetPrescription(id uuid.UUID) (*model.Prescription, error)
	ListPrescriptions(page, limit int) ([]model.Prescription, int64, error)
	UpdatePrescriptionStatus(id uuid.UUID, status string) error
	GetPatientPrescriptions(patientID uuid.UUID, page, limit int) ([]model.Prescription, int64, error)
	GetPrescriptionMedicines(prescriptionID uuid.UUID) ([]model.PrescriptionItem, error)
	BulkUpdateStatus(ctx context.Context, prescriptionIDs []uuid.UUID, status string, doctorID uuid.UUID) (*BulkOperationResult, error)
}

type prescriptionService struct {
	notificationService *NotificationService
}

func NewPrescriptionService(notificationService *NotificationService) PrescriptionService {
	return &prescriptionService{
		notificationService: notificationService,
	}
}

// CreatePrescription creates a new prescription with items in a transaction
func (s *prescriptionService) CreatePrescription(appointmentID, doctorID uuid.UUID, items []ItemInput, notes string) (*model.Prescription, error) {
	if appointmentID == uuid.Nil {
		return nil, errors.New("appointment_id is required")
	}

	if doctorID == uuid.Nil {
		return nil, errors.New("doctor_id is required")
	}

	if len(items) == 0 {
		return nil, errors.New("at least one prescription item is required")
	}

	db := config.GetDB()

	// Validate appointment and get its details
	appointment := &model.Appointment{}
	if err := db.Where("id = ?", appointmentID).First(appointment).Error; err != nil {
		return nil, errors.New("appointment not found")
	}

	if appointment.DoctorID != doctorID {
		return nil, errors.New("appointment does not belong to this doctor")
	}

	if appointment.Status != model.AppointmentInProgress {
		return nil, fmt.Errorf("appointment status must be 'in_progress', current status: %s", appointment.Status)
	}

	var prescription *model.Prescription

	// Use transaction for atomic operation
	err := db.Transaction(func(tx *gorm.DB) error {
		// Create prescription
		prescription = &model.Prescription{
			ID:               uuid.New(),
			PatientID:        appointment.PatientID,
			DoctorID:         doctorID,
			AppointmentID:    &appointmentID,
			PrescriptionDate: time.Now().Format("2006-01-02"),
			IssuedAt:         time.Now().Format("2006-01-02 15:04"),
			Notes:            notes,
			Status:           "active",
		}

		if err := tx.Create(prescription).Error; err != nil {
			return fmt.Errorf("failed to create prescription: %w", err)
		}

		// Create prescription items
		for _, item := range items {
			// Validate medicine exists - lookup by ID or name
			medicine := &model.Medicine{}
			var lookupErr error

			// Try lookup by medicine_id first if provided
			if item.MedicineID != uuid.Nil {
				lookupErr = tx.Where("id = ?", item.MedicineID).First(medicine).Error
			} else if item.MedicineName != "" {
				// Lookup by medicine name
				lookupErr = tx.Where("LOWER(name) = LOWER(?)", item.MedicineName).First(medicine).Error
			} else {
				return fmt.Errorf("either medicine_id or medicine_name is required for item")
			}

			if lookupErr != nil {
				medIdentifier := ""
				if item.MedicineID != uuid.Nil {
					medIdentifier = item.MedicineID.String()
				} else {
					medIdentifier = item.MedicineName
				}
				return fmt.Errorf("medicine not found: %s", medIdentifier)
			}

			prescriptionItem := &model.PrescriptionItem{
				ID:             uuid.New(),
				PrescriptionID: prescription.ID,
				MedicineID:     medicine.ID,
				Dosage:         item.Dosage,
				Frequency:      item.Frequency,
				Duration:       item.Duration,
				TimeOfDay:      item.TimeOfDay,
				Quantity:       item.Quantity,
				Instructions:   item.Instructions,
			}

			if err := tx.Create(prescriptionItem).Error; err != nil {
				return fmt.Errorf("failed to create prescription item: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Async notification to pharmacy
	go s.notifyPharmacy(prescription.ID, prescription.DoctorID, prescription.PatientID)

	return prescription, nil
}

// GetPrescription retrieves a prescription by ID
func (s *prescriptionService) GetPrescription(id uuid.UUID) (*model.Prescription, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid prescription id")
	}

	db := config.GetDB()
	prescription := &model.Prescription{}

	if err := db.Preload("Patient").Preload("Doctor").Preload("Appointment").First(prescription, "id = ?", id).Error; err != nil {
		return nil, errors.New("prescription not found")
	}

	return prescription, nil
}

// ListPrescriptions retrieves prescriptions with pagination
func (s *prescriptionService) ListPrescriptions(page, limit int) ([]model.Prescription, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	db := config.GetDB()
	var prescriptions []model.Prescription
	var total int64

	offset := (page - 1) * limit

	if err := db.Preload("Patient").Preload("Doctor").Preload("PrescriptionItems").Preload("PrescriptionItems.Medicine").
		Offset(offset).
		Limit(limit).
		Find(&prescriptions).
		Offset(-1).
		Limit(-1).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return prescriptions, total, nil
}

// UpdatePrescriptionStatus updates prescription status
func (s *prescriptionService) UpdatePrescriptionStatus(id uuid.UUID, status string) error {
	if id == uuid.Nil {
		return errors.New("invalid prescription id")
	}

	if status == "" {
		return errors.New("status is required")
	}

	db := config.GetDB()
	prescription := &model.Prescription{}

	if err := db.Model(prescription).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update prescription status: %w", err)
	}

	return nil
}

// GetPatientPrescriptions retrieves prescriptions for a patient
func (s *prescriptionService) GetPatientPrescriptions(patientID uuid.UUID, page, limit int) ([]model.Prescription, int64, error) {
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

	db := config.GetDB()
	var prescriptions []model.Prescription
	var total int64

	offset := (page - 1) * limit

	if err := db.Preload("Doctor").
		Where("patient_id = ?", patientID).
		Offset(offset).
		Limit(limit).
		Find(&prescriptions).
		Offset(-1).
		Limit(-1).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return prescriptions, total, nil
}

// GetPrescriptionMedicines retrieves all medicines in a prescription with details
func (s *prescriptionService) GetPrescriptionMedicines(prescriptionID uuid.UUID) ([]model.PrescriptionItem, error) {
	if prescriptionID == uuid.Nil {
		return nil, errors.New("invalid prescription id")
	}

	db := config.GetDB()
	var items []model.PrescriptionItem

	// Fetch prescription items with medicine details
	if err := db.Preload("Medicine").
		Where("prescription_id = ?", prescriptionID).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get prescription medicines: %w", err)
	}

	return items, nil
}

// notifyPharmacy sends async notification to pharmacy when prescription is created
func (s *prescriptionService) notifyPharmacy(prescriptionID, doctorID, patientID uuid.UUID) {
	if s.notificationService == nil {
		fmt.Printf("[Prescription] Notification service not available\n")
		return
	}

	// Get prescription details
	db := config.GetDB()
	prescription := &model.Prescription{}
	if err := db.Where("id = ?", prescriptionID).
		Preload("PrescriptionItems").
		Preload("Patient").
		First(prescription).Error; err != nil {
		fmt.Printf("[Prescription] Error fetching prescription: %v\n", err)
		return
	}

	// Get doctor name
	doctor := &model.User{}
	if err := db.Where("id = ?", doctorID).First(doctor).Error; err != nil {
		fmt.Printf("[Prescription] Error fetching doctor: %v\n", err)
		return
	}

	// Get patient name
	patient := &model.User{}
	if err := db.Where("id = ?", patientID).First(patient).Error; err != nil {
		fmt.Printf("[Prescription] Error fetching patient: %v\n", err)
		return
	}

	// Get all pharmacist users to notify
	var pharmacists []model.User
	if err := db.Where("role = ?", "pharmacist").Find(&pharmacists).Error; err != nil {
		fmt.Printf("[Prescription] Error fetching pharmacists: %v\n", err)
		return
	}

	itemCount := len(prescription.PrescriptionItems)
	patientName := patient.FirstName + " " + patient.LastName
	message := fmt.Sprintf("Patient %s ki prescription aayi hai, %d medicines hain", patientName, itemCount)

	// Notify all pharmacists
	for _, pharmacist := range pharmacists {
		job := NotificationJob{
			RecipientID: pharmacist.ID,
			Type:        "prescription",
			Title:       "New Prescription Received",
			Message:     message,
			Channel:     "in_app", // Can be extended to email/sms
			RelatedID:   &prescriptionID,
			RelatedType: "Prescription",
		}
		s.notificationService.Send(job)
		fmt.Printf("[Prescription] Pharmacist %s notified about prescription %s\n", pharmacist.Email, prescriptionID.String())
	}
}
