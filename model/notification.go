package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Title       string         `gorm:"not null" json:"title"`
	Message     string         `gorm:"type:text;not null" json:"message"`
	Type        string         `gorm:"not null;index" json:"type"` // appointment, prescription, lab_result, appointment, bill, general
	RelatedID   *uuid.UUID     `gorm:"type:uuid;index" json:"related_id"`
	RelatedType string         `json:"related_type"` // Appointment, Prescription, LabRequest, Bill, User
	IsRead      bool           `gorm:"default:false;index" json:"is_read"`
	ReadAt      *string        `json:"read_at"` // YYYY-MM-DD HH:MM
	IsPushSent  bool           `gorm:"default:false" json:"is_push_sent"`
	IsEmailSent bool           `gorm:"default:false" json:"is_email_sent"`
	Priority    string         `gorm:"default:'normal'" json:"priority"` // normal, high, critical
	ExpiresAt   *string        `json:"expires_at"`
	CreatedAt   int64          `gorm:"autoCreateTime:milli;index" json:"created_at"`
	UpdatedAt   int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Notification) TableName() string {
	return "notifications"
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
