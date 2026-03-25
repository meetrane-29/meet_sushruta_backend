package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DoctorSchedule struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	DoctorID    uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"doctor_id"`
	Doctor      *Doctor        `gorm:"foreignKey:DoctorID;references:ID" json:"doctor,omitempty"`
	DayOfWeek   string         `gorm:"not null;index" json:"day_of_week"` // Monday, Tuesday, etc.
	StartTime   string         `gorm:"not null" json:"start_time"`        // HH:MM format
	EndTime     string         `gorm:"not null" json:"end_time"`          // HH:MM format
	IsAvailable bool           `gorm:"default:true" json:"is_available"`
	MaxPatients int            `json:"max_patients"`
	CreatedAt   int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (DoctorSchedule) TableName() string {
	return "doctor_schedules"
}

func (ds *DoctorSchedule) BeforeCreate(tx *gorm.DB) error {
	if ds.ID == uuid.Nil {
		ds.ID = uuid.New()
	}
	return nil
}
