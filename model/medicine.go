package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Medicine struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name              string         `gorm:"uniqueIndex;not null" json:"name"`
	GenericName       string         `gorm:"index" json:"generic_name"`
	Dosage            string         `gorm:"not null" json:"dosage"` // e.g., "500mg"
	Manufacturer      string         `json:"manufacturer"`
	Description       string         `gorm:"type:text" json:"description"`
	Instructions      string         `gorm:"type:text" json:"instructions"`
	SideEffects       string         `gorm:"type:text" json:"side_effects"`
	Contraindications string         `gorm:"type:text" json:"contraindications"`
	StockQuantity     int            `gorm:"default:0;index" json:"stock_quantity"`
	ReorderLevel      int            `json:"reorder_level"`
	Price             float64        `gorm:"not null" json:"price"`
	Active            bool           `gorm:"default:true;index" json:"active"`
	CreatedAt         int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt         int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Medicine) TableName() string {
	return "medicines"
}

func (m *Medicine) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
