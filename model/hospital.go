package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hospital struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name          string         `gorm:"uniqueIndex;not null" json:"name"`
	Address       string         `gorm:"type:text;not null" json:"address"`
	City          string         `gorm:"index;not null" json:"city"`
	State         string         `gorm:"not null" json:"state"`
	PostalCode    string         `json:"postal_code"`
	Phone         string         `gorm:"not null" json:"phone"`
	Email         string         `gorm:"not null" json:"email"`
	Website       string         `json:"website"`
	Description   string         `gorm:"type:text" json:"description"`
	TotalBeds     int            `json:"total_beds"`
	AvailableBeds int            `json:"available_beds"`
	Latitude      float64        `json:"latitude"`
	Longitude     float64        `json:"longitude"`
	IsVerified    bool           `gorm:"default:false" json:"is_verified"`
	IsActive      bool           `gorm:"default:true;index" json:"is_active"`
	CreatedAt     int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt     int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Hospital) TableName() string {
	return "hospitals"
}

func (h *Hospital) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}
