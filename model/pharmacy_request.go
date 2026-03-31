package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PharmacyRequest struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PatientID      uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient        *Patient       `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	DoctorID       uuid.UUID      `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"doctor_id"`
	Doctor         *Doctor        `gorm:"foreignKey:DoctorID;references:ID" json:"doctor,omitempty"`
	PrescriptionID uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"prescription_id"`
	Prescription   *Prescription  `gorm:"foreignKey:PrescriptionID;references:ID" json:"prescription,omitempty"`
	RequestedAt    string         `gorm:"not null;index" json:"requested_at"` // YYYY-MM-DD HH:MM
	CompletedAt    *string        `json:"completed_at"`
	Status         string         `gorm:"default:'pending';index" json:"status"` // pending, acknowledged, dispensed, cancelled
	Priority       string         `gorm:"default:'normal'" json:"priority"`      // normal, urgent
	Notes          string         `gorm:"type:text" json:"notes"`
	CreatedAt      int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt      int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (PharmacyRequest) TableName() string {
	return "pharmacy_requests"
}

func (pr *PharmacyRequest) BeforeCreate(tx *gorm.DB) error {
	if pr.ID == uuid.Nil {
		pr.ID = uuid.New()
	}
	return nil
}
