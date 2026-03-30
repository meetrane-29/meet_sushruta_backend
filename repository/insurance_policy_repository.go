package repository

import (
	"errors"
	"fmt"
	"meet_sushruta/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InsurancePolicyRepository interface {
	CreateInsurancePolicy(policy *model.InsurancePolicy) error
	GetInsurancePolicyByID(id uuid.UUID) (*model.InsurancePolicy, error)
	GetInsurancePoliciesByPatientID(patientID uuid.UUID) ([]model.InsurancePolicy, error)
	GetActiveInsurancePolicyByPatientID(patientID uuid.UUID) (*model.InsurancePolicy, error)
	UpdateInsurancePolicy(policy *model.InsurancePolicy) error
	DeactivateInsurancePolicy(id uuid.UUID) error
}

type insurancePolicyRepository struct {
	db *gorm.DB
}

func NewInsurancePolicyRepository(db *gorm.DB) InsurancePolicyRepository {
	return &insurancePolicyRepository{db: db}
}

// CreateInsurancePolicy creates a new insurance policy
func (r *insurancePolicyRepository) CreateInsurancePolicy(policy *model.InsurancePolicy) error {
	if err := r.db.Create(policy).Error; err != nil {
		return fmt.Errorf("error creating insurance policy: %w", err)
	}
	return nil
}

// GetInsurancePolicyByID retrieves insurance policy by ID
func (r *insurancePolicyRepository) GetInsurancePolicyByID(id uuid.UUID) (*model.InsurancePolicy, error) {
	var policy model.InsurancePolicy
	if err := r.db.First(&policy, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching insurance policy: %w", err)
	}
	return &policy, nil
}

// GetInsurancePoliciesByPatientID retrieves all insurance policies for a patient
func (r *insurancePolicyRepository) GetInsurancePoliciesByPatientID(patientID uuid.UUID) ([]model.InsurancePolicy, error) {
	var policies []model.InsurancePolicy
	if err := r.db.Where("patient_id = ?", patientID).
		Order("created_at DESC").
		Find(&policies).Error; err != nil {
		return nil, fmt.Errorf("error fetching insurance policies: %w", err)
	}
	return policies, nil
}

// GetActiveInsurancePolicyByPatientID retrieves active insurance policy for a patient
func (r *insurancePolicyRepository) GetActiveInsurancePolicyByPatientID(patientID uuid.UUID) (*model.InsurancePolicy, error) {
	var policy model.InsurancePolicy
	if err := r.db.Where("patient_id = ? AND is_active = true", patientID).
		First(&policy).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching active insurance policy: %w", err)
	}
	return &policy, nil
}

// UpdateInsurancePolicy updates existing insurance policy
func (r *insurancePolicyRepository) UpdateInsurancePolicy(policy *model.InsurancePolicy) error {
	if err := r.db.Save(policy).Error; err != nil {
		return fmt.Errorf("error updating insurance policy: %w", err)
	}
	return nil
}

// DeactivateInsurancePolicy deactivates an insurance policy
func (r *insurancePolicyRepository) DeactivateInsurancePolicy(id uuid.UUID) error {
	if err := r.db.Model(&model.InsurancePolicy{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return fmt.Errorf("error deactivating insurance policy: %w", err)
	}
	return nil
}
