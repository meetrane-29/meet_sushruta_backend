package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdmissionStatus string

const (
	AdmissionActive    AdmissionStatus = "active"
	AdmissionDischarge AdmissionStatus = "discharged"
	AdmissionCancelled AdmissionStatus = "cancelled"
)

type AdmissionRecord struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	PatientID     uuid.UUID       `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient       *Patient        `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	DoctorID      uuid.UUID       `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"doctor_id"`
	Doctor        *Doctor         `gorm:"foreignKey:DoctorID;references:ID" json:"doctor,omitempty"`
	BedID         uuid.UUID       `gorm:"type:uuid;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"bed_id"`
	Bed           *Bed            `gorm:"foreignKey:BedID;references:ID" json:"bed,omitempty"`
	AdmissionDate string          `gorm:"not null;index" json:"admission_date"` // YYYY-MM-DD HH:MM
	DischargeDate *string         `gorm:"index" json:"discharge_date"`          // YYYY-MM-DD HH:MM (null if active)
	Reason        string          `gorm:"type:text" json:"reason"`              // Reason for admission
	Diagnosis     string          `gorm:"type:text" json:"diagnosis"`           // Initial diagnosis
	Status        AdmissionStatus `gorm:"not null;default:'active';index" json:"status"`
	Ward          string          `gorm:"index" json:"ward"` // ICU, General, etc.
	RoomNumber    string          `json:"room_number"`
	IsEmergency   bool            `gorm:"default:false" json:"is_emergency"`
	Notes         string          `gorm:"type:text" json:"notes"`
	CreatedAt     int64           `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt     int64           `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
}

func (AdmissionRecord) TableName() string {
	return "admission_records"
}

func (a *AdmissionRecord) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
