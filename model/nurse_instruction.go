package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NurseInstruction struct {
	ID          uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	AdmissionID uuid.UUID        `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"admission_id"`
	Admission   *AdmissionRecord `gorm:"foreignKey:AdmissionID;references:ID" json:"admission,omitempty"`
	PatientID   uuid.UUID        `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient     *Patient         `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	DoctorID    uuid.UUID        `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"doctor_id"`
	Doctor      *Doctor          `gorm:"foreignKey:DoctorID;references:ID" json:"doctor,omitempty"`
	IssuedDate  string           `gorm:"not null;index" json:"issued_date"` // YYYY-MM-DD

	// Vital Monitoring
	VitalsFrequency  string `json:"vitals_frequency"`                   // Every 2 hours, 4 hours, 6 hours, 8 hours, 12 hours
	VitalsParameters string `gorm:"type:text" json:"vitals_parameters"` // JSON: which vitals to monitor

	// Dietary Instructions
	DietType            string `json:"diet_type"` // NPO, Soft, Liquid, Regular, Diabetic, etc.
	DietaryRestrictions string `gorm:"type:text" json:"dietary_restrictions"`
	FluidsRestriction   string `json:"fluids_restriction"` // e.g., "1000 ml/day"

	// Medication Instructions
	MedicationNotes string `gorm:"type:text" json:"medication_notes"`
	DrugAllergies   string `gorm:"type:text" json:"drug_allergies"`

	// Activity Restrictions
	ActivityLevel string `json:"activity_level"` // Bed rest, Limited mobility, Ambulatory
	ActivityNotes string `gorm:"type:text" json:"activity_notes"`

	// Monitoring Parameters
	SpecialMonitoring string `gorm:"type:text" json:"special_monitoring"` // Strict I/O chart, Catheter care, Wound dressing, etc.
	CautionPoints     string `gorm:"type:text" json:"caution_points"`     // Things to watch for

	// General Instructions
	HygieneInstructions string `gorm:"type:text" json:"hygiene_instructions"`
	OtherInstructions   string `gorm:"type:text" json:"other_instructions"`

	Status    string         `json:"status"` // Active, Updated, Cancelled
	CreatedAt int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (NurseInstruction) TableName() string {
	return "nurse_instructions"
}

func (n *NurseInstruction) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
