package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UHID represents Unique Hospital ID for patients
type UHID struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PatientID      uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient        *Patient       `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	UHID           string         `gorm:"uniqueIndex;not null;size:20" json:"uhid"` // Format: MS-YYYY-XXXXX
	HospitalCode   string         `gorm:"not null;size:5" json:"hospital_code"`     // e.g., "MS" for hospital
	SequenceNumber int64          `gorm:"not null" json:"sequence_number"`
	IssuedDate     int64          `gorm:"autoCreateTime:milli;not null" json:"issued_date"`
	IsActive       bool           `gorm:"default:true;not null" json:"is_active"`
	CreatedAt      int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt      int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (UHID) TableName() string {
	return "uhids"
}

func (u *UHID) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
