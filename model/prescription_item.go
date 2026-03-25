package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrescriptionItem struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PrescriptionID uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"prescription_id"`
	Prescription   *Prescription  `gorm:"foreignKey:PrescriptionID;references:ID" json:"prescription,omitempty"`
	MedicineID     uuid.UUID      `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"medicine_id"`
	Medicine       *Medicine      `gorm:"foreignKey:MedicineID;references:ID" json:"medicine,omitempty"`
	Dosage         string         `gorm:"not null" json:"dosage"`    // e.g., "1 tablet"
	Frequency      string         `gorm:"not null" json:"frequency"` // e.g., "twice daily"
	Duration       string         `gorm:"not null" json:"duration"`  // e.g., "7 days"
	TimeOfDay      string         `json:"time_of_day"`               // e.g., "morning, evening"
	Quantity       int            `gorm:"not null" json:"quantity"`  // number of units
	Instructions   string         `gorm:"type:text" json:"instructions"`
	CreatedAt      int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt      int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (PrescriptionItem) TableName() string {
	return "prescription_items"
}

func (pi *PrescriptionItem) BeforeCreate(tx *gorm.DB) error {
	if pi.ID == uuid.Nil {
		pi.ID = uuid.New()
	}
	return nil
}
