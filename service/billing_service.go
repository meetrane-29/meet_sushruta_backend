package service

import (
	"errors"
	"fmt"
	"time"

	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
)

// BillItemRequest represents a line item for bill generation
type BillItemRequest struct {
	ItemType    string  `json:"item_type" binding:"required"` // consultation, prescription, lab, other
	Description string  `json:"description" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,min=0"`
	Quantity    int     `json:"quantity"`
}

type BillingService interface {
	GenerateBill(patientID uuid.UUID, appointmentID *uuid.UUID, items []BillItemRequest) (*model.Bill, error)
	GetBill(billID uuid.UUID) (*model.Bill, []model.BillItem, error)
	ProcessPayment(billID uuid.UUID, paymentMethod string, transactionID string) (*model.Bill, error)
	ListBills(page, limit int, filters map[string]interface{}) ([]model.Bill, int64, error)
}

type billingService struct {
	billingRepo repository.BillingRepository
	patientRepo repository.PatientRepository
}

func NewBillingService(billingRepo repository.BillingRepository, patientRepo repository.PatientRepository) BillingService {
	return &billingService{
		billingRepo: billingRepo,
		patientRepo: patientRepo,
	}
}

// GenerateBill creates a bill for a patient with specified line items
func (s *billingService) GenerateBill(patientID uuid.UUID, appointmentID *uuid.UUID, items []BillItemRequest) (*model.Bill, error) {
	if patientID == uuid.Nil {
		return nil, errors.New("patient_id is required")
	}

	if len(items) == 0 {
		return nil, errors.New("at least one bill item is required")
	}

	// Verify patient exists
	_, err := s.patientRepo.GetByID(patientID)
	if err != nil {
		return nil, errors.New("patient not found")
	}

	// Calculate total amount from items
	var totalAmount float64
	for _, item := range items {
		if item.Amount < 0 {
			return nil, errors.New("item amount cannot be negative")
		}
		totalAmount += item.Amount
	}

	// Apply taxes (default tax percentage: 0%)
	taxPercentage := 0.0
	taxAmount := (totalAmount * taxPercentage) / 100
	finalTotal := totalAmount + taxAmount

	// Create bill
	bill := &model.Bill{
		ID:            uuid.New(),
		PatientID:     patientID,
		AppointmentID: appointmentID,
		BillNumber:    fmt.Sprintf("BILL-%s-%d", patientID.String()[:8], time.Now().UnixMilli()),
		BillDate:      time.Now().Format("2006-01-02"),
		TaxPercentage: taxPercentage,
		TaxAmount:     taxAmount,
		TotalAmount:   finalTotal,
		PaidAmount:    0,
		DueAmount:     finalTotal,
		PaymentStatus: "pending",
	}

	// Create bill in database
	if err := s.billingRepo.Create(bill); err != nil {
		return nil, fmt.Errorf("failed to create bill: %w", err)
	}

	// Create bill items
	for _, itemReq := range items {
		quantity := itemReq.Quantity
		if quantity < 1 {
			quantity = 1
		}

		billItem := &model.BillItem{
			ID:          uuid.New(),
			BillID:      bill.ID,
			ItemType:    model.BillItemType(itemReq.ItemType),
			Description: itemReq.Description,
			Quantity:    quantity,
			UnitPrice:   itemReq.Amount / float64(quantity),
			Total:       itemReq.Amount,
		}

		if err := s.billingRepo.CreateBillItem(billItem); err != nil {
			return nil, fmt.Errorf("failed to create bill item: %w", err)
		}
	}

	// Return bill with items populated
	billItems, _ := s.billingRepo.GetBillItems(bill.ID)
	bill.ConsultationFee = 0
	bill.Medicines = 0
	bill.LabTests = 0

	// Aggregate items by type for response data structure
	for _, item := range billItems {
		switch item.ItemType {
		case model.BillItemConsultation:
			bill.ConsultationFee += item.Total
		case model.BillItemMedicine:
			bill.Medicines += item.Total
		case model.BillItemLab:
			bill.LabTests += item.Total
		}
	}

	return bill, nil
}

// GetBill retrieves a bill with all its details
func (s *billingService) GetBill(billID uuid.UUID) (*model.Bill, []model.BillItem, error) {
	if billID == uuid.Nil {
		return nil, nil, errors.New("invalid bill id")
	}

	bill, err := s.billingRepo.GetByID(billID)
	if err != nil {
		return nil, nil, errors.New("bill not found")
	}

	billItems, err := s.billingRepo.GetBillItems(billID)
	if err != nil {
		billItems = []model.BillItem{}
	}

	return bill, billItems, nil
}

// ProcessPayment processes payment for a bill
func (s *billingService) ProcessPayment(billID uuid.UUID, paymentMethod string, transactionID string) (*model.Bill, error) {
	if billID == uuid.Nil {
		return nil, errors.New("invalid bill id")
	}

	if paymentMethod == "" {
		return nil, errors.New("payment_method is required")
	}

	// Valid payment methods
	validMethods := map[string]bool{
		"cash":      true,
		"card":      true,
		"upi":       true,
		"insurance": true,
		"check":     true,
		"online":    true,
	}

	if !validMethods[paymentMethod] {
		return nil, errors.New("invalid payment method")
	}

	// Get bill
	bill, err := s.billingRepo.GetByID(billID)
	if err != nil {
		return nil, errors.New("bill not found")
	}

	// Check if bill is already paid
	if bill.PaymentStatus == "paid" {
		return nil, errors.New("bill is already paid")
	}

	// Check if bill is cancelled
	if bill.PaymentStatus == "cancelled" {
		return nil, errors.New("bill is cancelled and cannot be paid")
	}

	// Mark as paid
	now := time.Now()
	paymentDateStr := now.Format("2006-01-02 15:04:05")
	bill.PaymentStatus = "paid"
	bill.PaymentMethod = paymentMethod
	bill.PaymentDate = &paymentDateStr
	bill.PaidAmount = bill.TotalAmount
	bill.DueAmount = 0

	// Update bill
	if err := s.billingRepo.Update(bill); err != nil {
		return nil, fmt.Errorf("failed to update bill: %w", err)
	}

	return bill, nil
}

// ListBills retrieves bills with filtering and pagination
func (s *billingService) ListBills(page, limit int, filters map[string]interface{}) ([]model.Bill, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	// Check for status filter
	if status, ok := filters["status"].(string); ok && status != "" {
		return s.billingRepo.GetByStatus(status, page, limit)
	}

	// Check for patient_id filter
	if patientID, ok := filters["patient_id"].(uuid.UUID); ok && patientID != uuid.Nil {
		return s.billingRepo.GetByPatientID(patientID, page, limit)
	}

	// Check for date range filter
	if startDate, ok := filters["start_date"].(string); ok {
		if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
			return s.billingRepo.GetByDateRange(startDate, endDate, page, limit)
		}
	}

	// Default: return all bills
	return s.billingRepo.GetAll(page, limit)
}
