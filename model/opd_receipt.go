package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OPDReceipt represents OPD consultation receipt
type OPDReceipt struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BillID           uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bill_id"`
	Bill             *Bill          `gorm:"foreignKey:BillID;references:ID" json:"bill,omitempty"`
	AppointmentID    uuid.UUID      `gorm:"type:uuid;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"appointment_id"`
	Appointment      *Appointment   `gorm:"foreignKey:AppointmentID;references:ID" json:"appointment,omitempty"`
	ReceiptNumber    string         `gorm:"uniqueIndex;not null;size:20" json:"receipt_number"` // Format: RCP-YYYY-XXXXX
	ReceiptDate      string         `gorm:"not null;index" json:"receipt_date"`                 // YYYY-MM-DD
	ConsultationFee  float64        `gorm:"not null" json:"consultation_fee"`
	InsuranceCovered float64        `gorm:"default:0" json:"insurance_covered"` // Amount covered by insurance
	PatientPayable   float64        `gorm:"not null" json:"patient_payable"`
	PaidAmount       float64        `gorm:"not null" json:"paid_amount"`
	PaymentMethod    string         `json:"payment_method"`  // cash, card, upi, insurance
	TransactionRef   string         `json:"transaction_ref"` // For digital payments
	Notes            string         `gorm:"type:text" json:"notes"`
	PrintedCount     int64          `gorm:"default:0" json:"printed_count"` // Number of times printed
	LastPrintedTime  *int64         `json:"last_printed_time"`
	CreatedAt        int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt        int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (OPDReceipt) TableName() string {
	return "opd_receipts"
}

func (o *OPDReceipt) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
