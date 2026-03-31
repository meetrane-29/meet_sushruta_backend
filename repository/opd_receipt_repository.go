package repository

import (
	"errors"
	"fmt"
	"meet_sushruta/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OPDReceiptRepository interface {
	CreateOPDReceipt(receipt *model.OPDReceipt) error
	GetOPDReceiptByID(id uuid.UUID) (*model.OPDReceipt, error)
	GetOPDReceiptByReceiptNumber(receiptNumber string) (*model.OPDReceipt, error)
	GetOPDReceiptsByPatientID(patientID uuid.UUID, limit, offset int) ([]model.OPDReceipt, error)
	GetOPDReceiptsByBillID(billID uuid.UUID) (*model.OPDReceipt, error)
	UpdateOPDReceipt(receipt *model.OPDReceipt) error
	GetNextReceiptNumber() (string, error)
}

type opdReceiptRepository struct {
	db *gorm.DB
}

func NewOPDReceiptRepository(db *gorm.DB) OPDReceiptRepository {
	return &opdReceiptRepository{db: db}
}

// CreateOPDReceipt creates a new OPD receipt
func (r *opdReceiptRepository) CreateOPDReceipt(receipt *model.OPDReceipt) error {
	if err := r.db.Create(receipt).Error; err != nil {
		return fmt.Errorf("error creating OPD receipt: %w", err)
	}
	return nil
}

// GetOPDReceiptByID retrieves OPD receipt by ID
func (r *opdReceiptRepository) GetOPDReceiptByID(id uuid.UUID) (*model.OPDReceipt, error) {
	var receipt model.OPDReceipt
	if err := r.db.Preload("Bill").Preload("Appointment").
		First(&receipt, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching OPD receipt: %w", err)
	}
	return &receipt, nil
}

// GetOPDReceiptByReceiptNumber retrieves OPD receipt by receipt number
func (r *opdReceiptRepository) GetOPDReceiptByReceiptNumber(receiptNumber string) (*model.OPDReceipt, error) {
	var receipt model.OPDReceipt
	if err := r.db.Preload("Bill").Preload("Appointment").
		Where("receipt_number = ?", receiptNumber).
		First(&receipt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching OPD receipt: %w", err)
	}
	return &receipt, nil
}

// GetOPDReceiptsByPatientID retrieves OPD receipts for a patient with pagination
func (r *opdReceiptRepository) GetOPDReceiptsByPatientID(patientID uuid.UUID, limit, offset int) ([]model.OPDReceipt, error) {
	var receipts []model.OPDReceipt
	if err := r.db.Preload("Bill").Preload("Appointment").
		Where("bill_id IN (SELECT id FROM bills WHERE patient_id = ?)", patientID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&receipts).Error; err != nil {
		return nil, fmt.Errorf("error fetching OPD receipts: %w", err)
	}
	return receipts, nil
}

// GetOPDReceiptsByBillID retrieves OPD receipt by bill ID
func (r *opdReceiptRepository) GetOPDReceiptsByBillID(billID uuid.UUID) (*model.OPDReceipt, error) {
	var receipt model.OPDReceipt
	if err := r.db.Preload("Bill").Preload("Appointment").
		Where("bill_id = ?", billID).
		First(&receipt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching OPD receipt: %w", err)
	}
	return &receipt, nil
}

// UpdateOPDReceipt updates existing OPD receipt
func (r *opdReceiptRepository) UpdateOPDReceipt(receipt *model.OPDReceipt) error {
	if err := r.db.Save(receipt).Error; err != nil {
		return fmt.Errorf("error updating OPD receipt: %w", err)
	}
	return nil
}

// GetNextReceiptNumber generates next receipt number in format RCP-YYYY-XXXXX
func (r *opdReceiptRepository) GetNextReceiptNumber() (string, error) {
	var maxNumber int64
	currentYear := time.Now().Year()
	yearPrefix := fmt.Sprintf("%d-", currentYear)

	if err := r.db.Model(&model.OPDReceipt{}).
		Where("receipt_number LIKE ?", "RCP-"+yearPrefix+"%").
		Count(&maxNumber).Error; err != nil {
		return "", fmt.Errorf("error generating receipt number: %w", err)
	}

	receiptNumber := fmt.Sprintf("RCP-%d-%05d", currentYear, maxNumber+1)
	return receiptNumber, nil
}
