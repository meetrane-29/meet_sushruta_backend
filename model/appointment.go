package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentStatus string

const (
	AppointmentPending    AppointmentStatus = "pending"
	AppointmentConfirmed  AppointmentStatus = "confirmed"
	AppointmentInProgress AppointmentStatus = "in_progress"
	AppointmentCompleted  AppointmentStatus = "completed"
	AppointmentCancelled  AppointmentStatus = "cancelled"
)

type Appointment struct {
	ID              uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	PatientID       uuid.UUID         `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient         *Patient          `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	DoctorID        uuid.UUID         `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"doctor_id"`
	Doctor          *Doctor           `gorm:"foreignKey:DoctorID;references:ID" json:"doctor,omitempty"`
	AppointmentDate string            `gorm:"not null;index" json:"appointment_date"` // YYYY-MM-DD
	AppointmentTime string            `gorm:"not null" json:"appointment_time"`       // HH:MM
	Status          AppointmentStatus `gorm:"not null;default:'pending';index" json:"status"`
	Reason          string            `gorm:"type:text" json:"reason"`
	Notes           string            `gorm:"type:text" json:"notes"`
	CreatedAt       int64             `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt       int64             `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"deleted_at"`
}

func (Appointment) TableName() string {
	return "appointments"
}

func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
