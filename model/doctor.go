package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Doctor struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID               uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_id"`
	User                 *User          `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Specialization       string         `gorm:"not null;index" json:"specialization"`
	LicenseNumber        string         `gorm:"uniqueIndex;not null" json:"license_number"`
	CertificationURL     string         `json:"certification_url"`
	Bio                  string         `gorm:"type:text" json:"bio"`
	Department           string         `json:"department"`
	ConsultationFee      float64        `json:"consultation_fee"`
	JoiningDate          int64          `json:"joining_date"`
	Salary               float64        `json:"salary"`
	AttendancePercentage float64        `json:"attendance_percentage"`
	LeaveBalance         int            `json:"leave_balance"`
	CreatedAt            int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt            int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Doctor) TableName() string {
	return "doctors"
}

func (d *Doctor) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
