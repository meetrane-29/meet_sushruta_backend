package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type LabOrderInput struct {
	PatientID     uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID      uuid.UUID `json:"doctor_id" binding:"required"`
	TestType      string    `json:"test_type" binding:"required"`
	TestName      string    `json:"test_name" binding:"required"`
	ScheduledDate *string   `json:"scheduled_date"`
	Priority      string    `json:"priority"`
}

type UpdateLabStatusInput struct {
	Status string `json:"status" binding:"required"`
}

type LabService interface {
	CreateOrder(input *LabOrderInput) (*model.LabRequest, error)
	GetOrder(id uuid.UUID) (*model.LabRequest, error)
	ListOrders(page, limit int) ([]model.LabRequest, int64, error)
	UpdateStatus(id uuid.UUID, status string) error
	UploadReport(orderID, staffID uuid.UUID, filePath string, fileName string) (string, error)
	GetPatientLabOrders(patientID uuid.UUID, page, limit int) ([]model.LabRequest, int64, error)
	GetLabOrdersFiltered(ctx context.Context, filter *LabOrderFilter) ([]*model.LabRequest, int64, error)
	BulkUpdateStatus(ctx context.Context, orderIDs []uuid.UUID, status string, labID uuid.UUID) (*BulkOperationResult, error)
}

type labService struct {
	notificationService *NotificationService
}

func NewLabService(notificationService *NotificationService) LabService {
	return &labService{
		notificationService: notificationService,
	}
}

// CreateOrder creates a new lab order
func (s *labService) CreateOrder(input *LabOrderInput) (*model.LabRequest, error) {
	if input.PatientID == uuid.Nil {
		return nil, errors.New("patient_id is required")
	}

	if input.DoctorID == uuid.Nil {
		return nil, errors.New("doctor_id is required")
	}

	if input.TestType == "" {
		return nil, errors.New("test_type is required")
	}

	if input.TestName == "" {
		return nil, errors.New("test_name is required")
	}

	db := config.GetDB()

	// Validate patient exists
	patient := &model.Patient{}
	if err := db.Where("id = ?", input.PatientID).First(patient).Error; err != nil {
		return nil, errors.New("patient not found")
	}

	// Validate doctor exists
	doctor := &model.Doctor{}
	if err := db.Where("id = ?", input.DoctorID).First(doctor).Error; err != nil {
		return nil, errors.New("doctor not found")
	}

	priority := "normal"
	if input.Priority != "" {
		priority = input.Priority
	}

	labOrder := &model.LabRequest{
		ID:            uuid.New(),
		PatientID:     input.PatientID,
		DoctorID:      input.DoctorID,
		TestType:      input.TestType,
		TestName:      input.TestName,
		RequestedAt:   time.Now().Format("2006-01-02 15:04"),
		ScheduledDate: input.ScheduledDate,
		Status:        "pending",
		Priority:      priority,
	}

	if err := db.Create(labOrder).Error; err != nil {
		return nil, fmt.Errorf("failed to create lab order: %w", err)
	}

	return labOrder, nil
}

// GetOrder retrieves a lab order by ID
func (s *labService) GetOrder(id uuid.UUID) (*model.LabRequest, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid lab order id")
	}

	db := config.GetDB()
	labOrder := &model.LabRequest{}

	if err := db.Preload("Patient").Preload("Doctor").
		Where("id = ?", id).
		First(labOrder).Error; err != nil {
		return nil, errors.New("lab order not found")
	}

	return labOrder, nil
}

// ListOrders retrieves lab orders with pagination
func (s *labService) ListOrders(page, limit int) ([]model.LabRequest, int64, error) {
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
	var labOrders []model.LabRequest
	var total int64

	offset := (page - 1) * limit

	if err := db.Preload("Patient").Preload("Doctor").
		Offset(offset).
		Limit(limit).
		Find(&labOrders).
		Offset(-1).
		Limit(-1).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return labOrders, total, nil
}

// UpdateStatus updates the status of a lab order
func (s *labService) UpdateStatus(id uuid.UUID, status string) error {
	if id == uuid.Nil {
		return errors.New("invalid lab order id")
	}

	if status == "" {
		return errors.New("status is required")
	}

	db := config.GetDB()
	labOrder := &model.LabRequest{}

	updateData := map[string]interface{}{
		"status": status,
	}

	// If completing the order, set completed_at
	if status == "completed" {
		updateData["completed_at"] = time.Now().Format("2006-01-02 15:04")
	}

	if err := db.Model(labOrder).Where("id = ?", id).Updates(updateData).Error; err != nil {
		return fmt.Errorf("failed to update lab order status: %w", err)
	}

	return nil
}

// UploadReport uploads lab report and returns the file URL
func (s *labService) UploadReport(orderID, staffID uuid.UUID, filePath string, fileName string) (string, error) {
	if orderID == uuid.Nil {
		return "", errors.New("invalid lab order id")
	}

	if staffID == uuid.Nil {
		return "", errors.New("invalid staff id")
	}

	if filePath == "" {
		return "", errors.New("file path is required")
	}

	db := config.GetDB()

	// Get lab order
	labOrder := &model.LabRequest{}
	if err := db.Where("id = ?", orderID).First(labOrder).Error; err != nil {
		return "", errors.New("lab order not found")
	}

	// Create directory structure: ./uploads/lab/{patientID}/{YYYYMMDD}/
	uploadDir := fmt.Sprintf("./uploads/lab/%s/%s",
		labOrder.PatientID.String(),
		time.Now().Format("20060102"))

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Move or copy file to final location
	finalFilePath := filepath.Join(uploadDir, fileName)
	if err := copyFile(filePath, finalFilePath); err != nil {
		return "", fmt.Errorf("failed to save report file: %w", err)
	}

	// Generate accessible URL
	fileURL := fmt.Sprintf("/uploads/lab/%s/%s/%s",
		labOrder.PatientID.String(),
		time.Now().Format("20060102"),
		fileName)

	// Update lab order with result URL and status
	if err := db.Model(labOrder).Updates(map[string]interface{}{
		"result_url":   fileURL,
		"status":       "completed",
		"completed_at": time.Now().Format("2006-01-02 15:04"),
	}).Error; err != nil {
		return "", fmt.Errorf("failed to update lab order: %w", err)
	}

	// Async notifications
	go s.notifyPatientOfLabResult(labOrder.PatientID, labOrder.DoctorID, orderID, labOrder)
	go s.notifyDoctorOfLabResult(labOrder.DoctorID, labOrder.PatientID, orderID, labOrder)

	return fileURL, nil
}

// GetPatientLabOrders retrieves lab orders for a patient
func (s *labService) GetPatientLabOrders(patientID uuid.UUID, page, limit int) ([]model.LabRequest, int64, error) {
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
	var labOrders []model.LabRequest
	var total int64

	offset := (page - 1) * limit

	if err := db.Preload("Doctor").
		Where("patient_id = ?", patientID).
		Offset(offset).
		Limit(limit).
		Find(&labOrders).
		Offset(-1).
		Limit(-1).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return labOrders, total, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dst, data, 0644); err != nil {
		return err
	}

	return nil
}

// notifyPatientOfLabResult sends async notification to patient when lab results are ready
func (s *labService) notifyPatientOfLabResult(patientID, doctorID, labOrderID uuid.UUID, labOrder *model.LabRequest) {
	if s.notificationService == nil {
		fmt.Printf("[Lab] Notification service not available\n")
		return
	}

	// Get patient details
	db := config.GetDB()
	patient := &model.User{}
	if err := db.Where("id = ?", patientID).First(patient).Error; err != nil {
		fmt.Printf("[Lab] Error fetching patient: %v\n", err)
		return
	}

	message := fmt.Sprintf("Aapki %s ki report ready hai",
		labOrder.TestName)

	job := NotificationJob{
		RecipientID: patientID,
		Type:        "lab_result",
		Title:       "Lab Results Ready",
		Message:     message,
		Channel:     "in_app", // Can be extended to email/sms
		RelatedID:   &labOrderID,
		RelatedType: "LabRequest",
	}
	s.notificationService.Send(job)
	fmt.Printf("[Lab] Patient %s notified about lab result %s\n", patient.Email, labOrderID.String())
}

// notifyDoctorOfLabResult sends async notification to doctor when lab results are ready
func (s *labService) notifyDoctorOfLabResult(doctorID, patientID, labOrderID uuid.UUID, labOrder *model.LabRequest) {
	if s.notificationService == nil {
		fmt.Printf("[Lab] Notification service not available\n")
		return
	}

	// Get doctor details
	db := config.GetDB()
	doctor := &model.User{}
	if err := db.Where("id = ?", doctorID).First(doctor).Error; err != nil {
		fmt.Printf("[Lab] Error fetching doctor: %v\n", err)
		return
	}

	// Get patient name for context
	patient := &model.User{}
	if err := db.Where("id = ?", patientID).First(patient).Error; err != nil {
		fmt.Printf("[Lab] Error fetching patient: %v\n", err)
		return
	}

	message := fmt.Sprintf("Patient %s ki %s test complete ho gayi",
		patient.FirstName+" "+patient.LastName, labOrder.TestName)

	job := NotificationJob{
		RecipientID: doctorID,
		Type:        "lab_result",
		Title:       "Lab Test Completed",
		Message:     message,
		Channel:     "in_app", // Can be extended to email/sms
		RelatedID:   &labOrderID,
		RelatedType: "LabRequest",
	}
	s.notificationService.Send(job)
	fmt.Printf("[Lab] Doctor %s notified about lab result %s\n", doctor.Email, labOrderID.String())
}
