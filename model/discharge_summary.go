package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PatientOutcome string

const (
	OutcomeRecovered            PatientOutcome = "recovered"
	OutcomeImproved             PatientOutcome = "improved"
	OutcomeStableCondition      PatientOutcome = "stable"
	OutcomeReferredToSpecialist PatientOutcome = "referred"
	OutcomeAMADischarge         PatientOutcome = "ama_discharge" // Against Medical Advice
	OutcomeDeceased             PatientOutcome = "deceased"
)

type DischargeSummary struct {
	ID            uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	AdmissionID   uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"admission_id"`
	Admission     *AdmissionRecord `gorm:"foreignKey:AdmissionID;references:ID" json:"admission,omitempty"`
	PatientID     uuid.UUID        `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient       *Patient         `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	DoctorID      uuid.UUID        `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"doctor_id"`
	Doctor        *Doctor          `gorm:"foreignKey:DoctorID;references:ID" json:"doctor,omitempty"`
	DischargeDate string           `gorm:"not null;index" json:"discharge_date"` // YYYY-MM-DD HH:MM

	// Clinical Information
	FinalDiagnosis      string `gorm:"type:text;not null" json:"final_diagnosis"`
	ProceduresPerformed string `gorm:"type:text" json:"procedures_performed"`
	ComplicationsIfAny  string `gorm:"type:text" json:"complications_if_any"`

	// Medication at Discharge
	DischargeMedications   string `gorm:"type:text" json:"discharge_medications"` // JSON formatted list
	MedicationInstructions string `gorm:"type:text" json:"medication_instructions"`

	// Diet & Activity
	DietRecommendation     string `gorm:"type:text" json:"diet_recommendation"`
	ActivityRecommendation string `gorm:"type:text" json:"activity_recommendation"`

	// Follow-up
	FollowUpInstructions string  `gorm:"type:text" json:"follow_up_instructions"`
	FollowUpDate         *string `json:"follow_up_date"`      // YYYY-MM-DD
	FollowUpDoctor       *string `json:"follow_up_doctor"`    // Doctor name to see
	FollowUpSpecialty    *string `json:"follow_up_specialty"` // If referred to specialist

	// Warning Symptoms
	WarningSymptoms      string `gorm:"type:text" json:"warning_symptoms"` // Red flags patient should watch for
	WhenToReturnHospital string `gorm:"type:text" json:"when_to_return_hospital"`

	// Outcomes
	Outcome          PatientOutcome `gorm:"not null" json:"outcome"`
	PatientEducation string         `gorm:"type:text" json:"patient_education"`

	// Additional Notes
	Notes     string         `gorm:"type:text" json:"notes"`
	CreatedAt int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (DischargeSummary) TableName() string {
	return "discharge_summaries"
}

func (d *DischargeSummary) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
