package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email                string         `gorm:"uniqueIndex;not null" json:"email"`
	Phone                string         `gorm:"uniqueIndex;not null" json:"phone"`
	FirstName            string         `gorm:"not null" json:"first_name"`
	LastName             string         `gorm:"not null" json:"last_name"`
	Role                 string         `gorm:"not null;index" json:"role"` // admin, patient, doctor, nurse, pharmacy, lab
	Password             string         `gorm:"not null" json:"password"`
	Active               bool           `gorm:"default:true;index" json:"active"`
	IsVerified           bool           `gorm:"default:false;index" json:"is_verified"`
	JoiningDate          int64          `json:"joining_date"`
	Salary               float64        `json:"salary"`
	AttendancePercentage float64        `json:"attendance_percentage"`
	LeaveBalance         int            `json:"leave_balance"`
	LastLogin            *int64         `gorm:"index" json:"last_login"`
	CreatedAt            int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt            int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
