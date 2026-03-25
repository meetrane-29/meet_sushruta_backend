package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Nurse struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_id"`
	User             *User          `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	LicenseNumber    string         `gorm:"uniqueIndex;not null" json:"license_number"`
	Department       string         `gorm:"index" json:"department"`
	Shift            string         `json:"shift"` // morning, evening, night
	AssignedBedID    *uuid.UUID     `gorm:"type:uuid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"assigned_bed_id"`
	AssignedBed      *Bed           `gorm:"foreignKey:AssignedBedID;references:ID" json:"assigned_bed,omitempty"`
	CertificationURL string         `json:"certification_url"`
	CreatedAt        int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt        int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Nurse) TableName() string {
	return "nurses"
}

func (n *Nurse) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
