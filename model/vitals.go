package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vitals struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PatientID       uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient         *Patient       `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	Temperature     float64        `gorm:"not null" json:"temperature"`                                                // in Celsius
	BloodPressure   string         `gorm:"not null" json:"blood_pressure"`                                             // e.g., "120/80"
	HeartRate       int            `gorm:"not null" json:"heart_rate"`                                                 // beats per minute
	RespiratoryRate int            `gorm:"not null" json:"respiratory_rate"`                                           // breaths per minute
	Weight          float64        `gorm:"not null" json:"weight"`                                                     // in kg
	Height          float64        `json:"height"`                                                                     // in cm
	BloodSugar      float64        `json:"blood_sugar"`                                                                // in mg/dL
	Oxygen          int            `json:"oxygen"`                                                                     // SpO2 percentage
	RecordedAt      string         `gorm:"not null;index" json:"recorded_at"`                                          // YYYY-MM-DD HH:MM
	RecordedBy      uuid.UUID      `gorm:"type:uuid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"recorded_by"` // User ID of who recorded
	CreatedAt       int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt       int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Vitals) TableName() string {
	return "vitals"
}

func (v *Vitals) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}
