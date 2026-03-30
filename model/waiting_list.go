package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WaitingListStatus string

const (
	WaitingListWaiting   WaitingListStatus = "waiting"
	WaitingListCalled    WaitingListStatus = "called"
	WaitingListSeen      WaitingListStatus = "seen"
	WaitingListCompleted WaitingListStatus = "completed"
	WaitingListCancelled WaitingListStatus = "cancelled"
)

// WaitingListEntry represents OPD waiting list/queue management
type WaitingListEntry struct {
	ID                uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	AppointmentID     uuid.UUID         `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"appointment_id"`
	Appointment       *Appointment      `gorm:"foreignKey:AppointmentID;references:ID" json:"appointment,omitempty"`
	PatientID         uuid.UUID         `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient           *Patient          `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	DoctorID          uuid.UUID         `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"doctor_id"`
	Doctor            *Doctor           `gorm:"foreignKey:DoctorID;references:ID" json:"doctor,omitempty"`
	TokenNumber       int64             `gorm:"not null;index" json:"token_number"` // Queue position
	Status            WaitingListStatus `gorm:"not null;default:'waiting';index" json:"status"`
	ArrivalTime       *int64            `gorm:"index" json:"arrival_time"` // When patient arrived
	CalledTime        *int64            `json:"called_time"`               // When called by receptionist
	CompletionTime    *int64            `json:"completion_time"`           // When consultation completed
	EstimatedWaitTime int64             `json:"estimated_wait_time"`       // Minutes
	Notes             string            `gorm:"type:text" json:"notes"`
	CreatedAt         int64             `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt         int64             `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt         gorm.DeletedAt    `gorm:"index" json:"deleted_at"`
}

func (WaitingListEntry) TableName() string {
	return "waiting_list_entries"
}

func (w *WaitingListEntry) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
