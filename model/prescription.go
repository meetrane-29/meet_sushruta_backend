package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Prescription struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PatientID        uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient          *Patient       `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	DoctorID         uuid.UUID      `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"doctor_id"`
	Doctor           *Doctor        `gorm:"foreignKey:DoctorID;references:ID" json:"doctor,omitempty"`
	AppointmentID    *uuid.UUID     `gorm:"type:uuid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"appointment_id"`
	Appointment      *Appointment   `gorm:"foreignKey:AppointmentID;references:ID" json:"appointment,omitempty"`
	PrescriptionDate string         `gorm:"not null;index" json:"prescription_date"` // YYYY-MM-DD
	Notes            string         `gorm:"type:text" json:"notes"`
	IssuedAt         string         `gorm:"not null" json:"issued_at"` // YYYY-MM-DD HH:MM
	ExpiresAt        *string        `json:"expires_at"`
	Status           string         `gorm:"default:'active';index" json:"status"` // active, expired, completed
	CreatedAt        int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt        int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Prescription) TableName() string {
	return "prescriptions"
}

func (p *Prescription) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
