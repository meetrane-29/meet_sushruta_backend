package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rating struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	DoctorID        uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"doctor_id"`
	Doctor          *Doctor        `gorm:"foreignKey:DoctorID;references:ID" json:"doctor,omitempty"`
	PatientID       uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient         *Patient       `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	Rating          int            `gorm:"not null;default:5" json:"rating"`          // 1-5 scale
	Professionalism int            `gorm:"not null;default:5" json:"professionalism"` // 1-5 scale
	Communication   int            `gorm:"not null;default:5" json:"communication"`   // 1-5 scale
	Punctuality     int            `gorm:"not null;default:5" json:"punctuality"`     // 1-5 scale
	Cleanliness     int            `gorm:"not null;default:5" json:"cleanliness"`     // 1-5 scale
	Comment         string         `gorm:"type:text" json:"comment"`
	CreatedAt       int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt       int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Rating) TableName() string {
	return "ratings"
}

func (r *Rating) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
