package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DispenseHistory struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PrescriptionID uuid.UUID      `gorm:"type:uuid;index;not null" json:"prescription_id"`
	MedicineID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"medicine_id"`
	Medicine       *Medicine      `gorm:"foreignKey:MedicineID" json:"medicine,omitempty"`
	Quantity       int            `gorm:"not null" json:"quantity"`
	UnitPrice      float64        `gorm:"not null" json:"unit_price"`
	TotalPrice     float64        `gorm:"not null" json:"total_price"`
	DispensedBy    uuid.UUID      `gorm:"type:uuid;not null" json:"dispensed_by"`
	DispensingUser *User          `gorm:"foreignKey:DispensedBy" json:"dispensing_user,omitempty"`
	DispensedAt    int64          `gorm:"not null" json:"dispensed_at"`
	BatchNumber    string         `json:"batch_number"`
	ExpiryDate     string         `json:"expiry_date"`
	CreatedAt      int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt      int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (DispenseHistory) TableName() string {
	return "dispense_histories"
}

func (d *DispenseHistory) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
