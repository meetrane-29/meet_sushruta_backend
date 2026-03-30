package service

import (
	"fmt"
	"meet_sushruta/model"
	"meet_sushruta/repository"
	"time"

	"github.com/google/uuid"
)

type OPDReceiptService interface {
	GenerateOPDReceipt(billID uuid.UUID, appointmentID *uuid.UUID, consultationFee, insuranceCovered float64, paymentMethod string) (*model.OPDReceipt, error)
	GetOPDReceiptByID(id uuid.UUID) (*model.OPDReceipt, error)
	GetOPDReceiptByReceiptNumber(receiptNumber string) (*model.OPDReceipt, error)
	GetPatientOPDReceipts(patientID uuid.UUID, limit, offset int) ([]model.OPDReceipt, error)
	UpdatePrintCount(receiptID uuid.UUID) error
	CalculatePatientPayable(consultationFee, insuranceCovered float64) float64
}

type opdReceiptService struct {
	receiptRepo repository.OPDReceiptRepository
}

func NewOPDReceiptService(receiptRepo repository.OPDReceiptRepository) OPDReceiptService {
	return &opdReceiptService{
		receiptRepo: receiptRepo,
	}
}

// GenerateOPDReceipt creates a new OPD receipt for a patient
func (s *opdReceiptService) GenerateOPDReceipt(
	billID uuid.UUID,
	appointmentID *uuid.UUID,
	consultationFee,
	insuranceCovered float64,
	paymentMethod string,
) (*model.OPDReceipt, error) {

	// Generate receipt number
	receiptNum, err := s.receiptRepo.GetNextReceiptNumber()
	if err != nil {
		return nil, fmt.Errorf("error generating receipt number: %w", err)
	}

	// Calculate patient payable amount
	patientPayable := s.CalculatePatientPayable(consultationFee, insuranceCovered)

	today := time.Now().Format("2006-01-02")

	// Handle appointmentID which may be nil
	var appID uuid.UUID
	if appointmentID != nil {
		appID = *appointmentID
	} else {
		appID = uuid.Nil
	}

	receipt := &model.OPDReceipt{
		ID:               uuid.New(),
		BillID:           billID,
		AppointmentID:    appID,
		ReceiptNumber:    receiptNum,
		ReceiptDate:      today,
		ConsultationFee:  consultationFee,
		InsuranceCovered: insuranceCovered,
		PatientPayable:   patientPayable,
		PaidAmount:       patientPayable,
		PaymentMethod:    paymentMethod,
		PrintedCount:     0,
	}

	if err := s.receiptRepo.CreateOPDReceipt(receipt); err != nil {
		return nil, fmt.Errorf("error creating OPD receipt: %w", err)
	}

	return receipt, nil
}

// GetOPDReceiptByID retrieves OPD receipt by ID
func (s *opdReceiptService) GetOPDReceiptByID(id uuid.UUID) (*model.OPDReceipt, error) {
	receipt, err := s.receiptRepo.GetOPDReceiptByID(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching OPD receipt: %w", err)
	}
	return receipt, nil
}

// GetOPDReceiptByReceiptNumber retrieves OPD receipt by receipt number
func (s *opdReceiptService) GetOPDReceiptByReceiptNumber(receiptNumber string) (*model.OPDReceipt, error) {
	receipt, err := s.receiptRepo.GetOPDReceiptByReceiptNumber(receiptNumber)
	if err != nil {
		return nil, fmt.Errorf("error fetching OPD receipt: %w", err)
	}
	return receipt, nil
}

// GetPatientOPDReceipts retrieves OPD receipts for a patient
func (s *opdReceiptService) GetPatientOPDReceipts(patientID uuid.UUID, limit, offset int) ([]model.OPDReceipt, error) {
	receipts, err := s.receiptRepo.GetOPDReceiptsByPatientID(patientID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error fetching patient OPD receipts: %w", err)
	}
	return receipts, nil
}

// UpdatePrintCount increments print count for a receipt
func (s *opdReceiptService) UpdatePrintCount(receiptID uuid.UUID) error {
	receipt, err := s.receiptRepo.GetOPDReceiptByID(receiptID)
	if err != nil {
		return fmt.Errorf("error fetching receipt: %w", err)
	}

	if receipt == nil {
		return fmt.Errorf("receipt not found")
	}

	receipt.PrintedCount++
	now := time.Now().UnixMilli()
	receipt.LastPrintedTime = &now

	if err := s.receiptRepo.UpdateOPDReceipt(receipt); err != nil {
		return fmt.Errorf("error updating print count: %w", err)
	}

	return nil
}

// CalculatePatientPayable calculates final patient payable amount
// Patient pays: Consultation Fee - Insurance Covered Amount
func (s *opdReceiptService) CalculatePatientPayable(consultationFee, insuranceCovered float64) float64 {
	if insuranceCovered > consultationFee {
		insuranceCovered = consultationFee
	}
	return consultationFee - insuranceCovered
}
