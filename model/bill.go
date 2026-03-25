package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bill struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PatientID       uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient         *Patient       `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	AppointmentID   *uuid.UUID     `gorm:"type:uuid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"appointment_id"`
	Appointment     *Appointment   `gorm:"foreignKey:AppointmentID;references:ID" json:"appointment,omitempty"`
	BillNumber      string         `gorm:"uniqueIndex;not null" json:"bill_number"`
	BillDate        string         `gorm:"not null;index" json:"bill_date"` // YYYY-MM-DD
	ConsultationFee float64        `gorm:"not null" json:"consultation_fee"`
	LabTests        float64        `json:"lab_tests"`
	Medicines       float64        `json:"medicines"`
	OtherCharges    float64        `json:"other_charges"`
	Discount        float64        `gorm:"default:0" json:"discount"`
	TaxPercentage   float64        `json:"tax_percentage"`
	TaxAmount       float64        `json:"tax_amount"`
	TotalAmount     float64        `gorm:"not null" json:"total_amount"`
	PaidAmount      float64        `gorm:"default:0" json:"paid_amount"`
	DueAmount       float64        `gorm:"not null" json:"due_amount"`
	PaymentStatus   string         `gorm:"default:'pending';index" json:"payment_status"` // pending, paid, partial, cancelled
	PaymentDate     *string        `json:"payment_date"`
	PaymentMethod   string         `json:"payment_method"` // cash, card, check, online
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt       int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Bill) TableName() string {
	return "bills"
}

func (b *Bill) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
