package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog is immutable - no DeletedAt field
type AuditLog struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Action        string    `gorm:"not null;index" json:"action"`      // CREATE, UPDATE, DELETE
	EntityType    string    `gorm:"not null;index" json:"entity_type"` // e.g., User, Patient, Appointment
	EntityID      uuid.UUID `gorm:"type:uuid;not null;index" json:"entity_id"`
	OldValues     string    `gorm:"type:text" json:"old_values"`     // JSON string
	NewValues     string    `gorm:"type:text" json:"new_values"`     // JSON string
	ChangedFields string    `gorm:"type:text" json:"changed_fields"` // Comma-separated list
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	Status        string    `gorm:"default:'success'" json:"status"` // success, failure
	Description   string    `gorm:"type:text" json:"description"`
	CreatedAt     int64     `gorm:"autoCreateTime:milli;index" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

func (al *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == uuid.Nil {
		al.ID = uuid.New()
	}
	return nil
}
