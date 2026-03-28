package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Specialization struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name           string         `gorm:"uniqueIndex;not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	Category       string         `gorm:"index" json:"category"` // e.g., "Surgery", "Medicine", "Pediatrics"
	Icon           string         `json:"icon"`                   // icon name or URL
	ParentID       *uuid.UUID     `gorm:"type:uuid;index" json:"parent_id"` // for sub-specializations
	IsActive       bool           `gorm:"default:true;index" json:"is_active"`
	CreatedAt      int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt      int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Specialization) TableName() string {
	return "specializations"
}

func (s *Specialization) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
