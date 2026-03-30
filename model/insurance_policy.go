package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InsurancePolicy represents patient insurance details
type InsurancePolicy struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PatientID         uuid.UUID      `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"patient_id"`
	Patient           *Patient       `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	ProviderName      string         `gorm:"not null" json:"provider_name"` // Insurance company name
	PolicyNumber      string         `gorm:"not null;uniqueIndex" json:"policy_number"`
	MemberID          string         `json:"member_id"`
	CoveragePlanName  string         `json:"coverage_plan_name"`
	CoverageAmount    float64        `gorm:"not null" json:"coverage_amount"`
	CopayPercentage   float64        `gorm:"default:0" json:"copay_percentage"` // Patient pays this percentage
	DeductibleAmount  float64        `gorm:"default:0" json:"deductible_amount"`
	ValidFrom         string         `gorm:"not null" json:"valid_from"` // YYYY-MM-DD
	ValidUpto         string         `gorm:"not null" json:"valid_upto"` // YYYY-MM-DD
	NomineeNames      string         `json:"nominee_names"`
	EmployerName      string         `json:"employer_name"` // For employer-based insurance
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	ProviderContactNo string         `json:"provider_contact_no"`
	ProviderWebsite   string         `json:"provider_website"`
	Notes             string         `gorm:"type:text" json:"notes"`
	CreatedAt         int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt         int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (InsurancePolicy) TableName() string {
	return "insurance_policies"
}

func (ip *InsurancePolicy) BeforeCreate(tx *gorm.DB) error {
	if ip.ID == uuid.Nil {
		ip.ID = uuid.New()
	}
	return nil
}
