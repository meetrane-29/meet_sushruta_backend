package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabRequest struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PatientID     uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient       *Patient       `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	DoctorID      uuid.UUID      `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"doctor_id"`
	Doctor        *Doctor        `gorm:"foreignKey:DoctorID;references:ID" json:"doctor,omitempty"`
	TestType      string         `gorm:"not null;index" json:"test_type"`    // e.g., "BloodWork", "Xray"
	TestName      string         `gorm:"not null" json:"test_name"`          // e.g., "CBC", "Chest X-ray"
	RequestedAt   string         `gorm:"not null;index" json:"requested_at"` // YYYY-MM-DD HH:MM
	ScheduledDate *string        `json:"scheduled_date"`                     // YYYY-MM-DD
	CompletedAt   *string        `json:"completed_at"`
	ResultURL     string         `json:"result_url"`
	ResultSummary string         `gorm:"type:text" json:"result_summary"`
	Status        string         `gorm:"default:'pending';index" json:"status"` // pending, completed, cancelled
	Priority      string         `gorm:"default:'normal'" json:"priority"`      // normal, urgent
	CreatedAt     int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt     int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (LabRequest) TableName() string {
	return "lab_requests"
}

func (lr *LabRequest) BeforeCreate(tx *gorm.DB) error {
	if lr.ID == uuid.Nil {
		lr.ID = uuid.New()
	}
	return nil
}
