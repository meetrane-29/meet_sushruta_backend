package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Patient struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_id"`
	User             *User          `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	DateOfBirth      string         `gorm:"not null" json:"date_of_birth"`
	Gender           string         `gorm:"not null" json:"gender"`
	BloodGroup       string         `json:"blood_group"`
	Address          string         `json:"address"`
	EmergencyContact string         `json:"emergency_contact"`
	MedicalHistory   string         `gorm:"type:text" json:"medical_history"`
	Allergies        string         `gorm:"type:text" json:"allergies"`
	CreatedAt        int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt        int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Patient) TableName() string {
	return "patients"
}

func (p *Patient) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
