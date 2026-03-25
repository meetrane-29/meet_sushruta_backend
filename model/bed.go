package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bed struct {
	ID                    uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BedNumber             string         `gorm:"uniqueIndex;not null" json:"bed_number"`
	Ward                  string         `gorm:"not null;index" json:"ward"`
	Room                  string         `gorm:"index" json:"room"`
	Floor                 int            `json:"floor"`
	BedType               string         `gorm:"not null" json:"bed_type"`                // ICU, General, Private, Semi-Private
	Status                string         `gorm:"default:'available';index" json:"status"` // available, occupied, maintenance, reserved
	PatientID             *uuid.UUID     `gorm:"type:uuid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"patient_id"`
	Patient               *Patient       `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	AdmittedAt            *string        `json:"admitted_at"`   // YYYY-MM-DD HH:MM
	DischargedAt          *string        `json:"discharged_at"` // YYYY-MM-DD HH:MM
	Features              string         `gorm:"type:text" json:"features"`
	DailyRate             float64        `json:"daily_rate"`
	IsMaintenanceRequired bool           `gorm:"default:false" json:"is_maintenance_required"`
	MaintenanceNotes      string         `gorm:"type:text" json:"maintenance_notes"`
	CreatedAt             int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt             int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Bed) TableName() string {
	return "beds"
}

func (b *Bed) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
