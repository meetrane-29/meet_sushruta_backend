package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProgressNote struct {
	ID           uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	AdmissionID  uuid.UUID        `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"admission_id"`
	Admission    *AdmissionRecord `gorm:"foreignKey:AdmissionID;references:ID" json:"admission,omitempty"`
	PatientID    uuid.UUID        `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient      *Patient         `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	DoctorID     uuid.UUID        `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"doctor_id"`
	Doctor       *Doctor          `gorm:"foreignKey:DoctorID;references:ID" json:"doctor,omitempty"`
	RecordedDate string           `gorm:"not null;index" json:"recorded_date"` // YYYY-MM-DD

	// SOAP Format
	Subjective string `gorm:"type:text" json:"subjective"` // Patient's complaints and symptoms
	Objective  string `gorm:"type:text" json:"objective"`  // Vital signs, findings, test results
	Assessment string `gorm:"type:text" json:"assessment"` // Doctor's diagnosis and impression
	Plan       string `gorm:"type:text" json:"plan"`       // Treatment plan, medications, follow-up

	Vitals         string         `gorm:"type:text" json:"vitals"`       // JSON string of vitals recorded (temperature, BP, etc.)
	Medications    string         `gorm:"type:text" json:"medications"`  // Current medications
	NextReviewDate *string        `gorm:"index" json:"next_review_date"` // Next check-up date
	Notes          string         `gorm:"type:text" json:"notes"`        // Additional notes
	CreatedAt      int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt      int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (ProgressNote) TableName() string {
	return "progress_notes"
}

func (p *ProgressNote) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
