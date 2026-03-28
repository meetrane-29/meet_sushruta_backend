package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OperationType string

const (
	OperationScheduled  OperationType = "scheduled"
	OperationInProgress OperationType = "in_progress"
	OperationCompleted  OperationType = "completed"
	OperationCancelled  OperationType = "cancelled"
	OperationPostponed  OperationType = "postponed"
)

type OperationTheatre struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TheatreName          string         `gorm:"uniqueIndex;not null" json:"theatre_name"` // OT 1, OT 2, etc.
	Floor                int            `json:"floor"`
	Capacity             int            `json:"capacity"`                                // Number of patients that can be operated simultaneously
	Status               string         `gorm:"default:'available';index" json:"status"` // available, in_use, maintenance
	Features             string         `gorm:"type:text" json:"features"`               // Comma-separated: Laminar Flow, Advanced Lighting, Video, etc.
	EquipmentList        string         `gorm:"type:text" json:"equipment_list"`         // Comma-separated equipment IDs/names available
	LastSanitizationDate string         `json:"last_sanitization_date"`                  // YYYY-MM-DD
	SanitizationDueDate  string         `json:"sanitization_due_date"`                   // YYYY-MM-DD
	Notes                string         `gorm:"type:text" json:"notes"`
	CreatedAt            int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt            int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (OperationTheatre) TableName() string {
	return "operation_theatres"
}

func (o *OperationTheatre) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// OperationSchedule tracks all surgical operations
type OperationSchedule struct {
	ID                 uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	TheatreID          uuid.UUID         `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"theatre_id"`
	Theatre            *OperationTheatre `gorm:"foreignKey:TheatreID;references:ID" json:"theatre,omitempty"`
	PatientID          uuid.UUID         `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient            *Patient          `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	SurgieDoctorID     uuid.UUID         `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"surgeon_doctor_id"`
	SurgeonDoctor      *Doctor           `gorm:"foreignKey:SurgieDoctorID;references:ID;foreignKeyConstraint:OnDelete:CASCADE" json:"surgeon_doctor,omitempty"`
	AnesthesiaDoctorID *uuid.UUID        `gorm:"type:uuid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"anesthesia_doctor_id"`
	AnesthesiaDoctor   *Doctor           `gorm:"foreignKey:AnesthesiaDoctorID;references:ID" json:"anesthesia_doctor,omitempty"`
	OperationDate      string            `gorm:"not null;index" json:"operation_date"` // YYYY-MM-DD
	OperationTime      string            `gorm:"not null" json:"operation_time"`       // HH:MM
	EstimatedDuration  int               `json:"estimated_duration"`                   // in minutes
	ActualDuration     int               `json:"actual_duration"`                      // in minutes
	OperationType      string            `gorm:"not null" json:"operation_type"`       // e.g., "Appendectomy", "Hernia Repair", etc.
	Status             OperationType     `gorm:"default:'scheduled';index" json:"status"`
	Diagnosis          string            `gorm:"type:text" json:"diagnosis"`
	Notes              string            `gorm:"type:text" json:"notes"`
	PreOperativeNotes  string            `gorm:"type:text" json:"pre_operative_notes"`
	PostOperativeNotes string            `gorm:"type:text" json:"post_operative_notes"`
	Complications      string            `gorm:"type:text" json:"complications"`
	BloodLossEstimate  string            `json:"blood_loss_estimate"` // in ml
	CreatedAt          int64             `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt          int64             `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt          gorm.DeletedAt    `gorm:"index" json:"deleted_at"`
}

func (OperationSchedule) TableName() string {
	return "operation_schedules"
}

func (o *OperationSchedule) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
