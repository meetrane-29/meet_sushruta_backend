package service

import (
	"fmt"
	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
)

type InsurancePolicyService interface {
	CreateInsurancePolicy(policy *model.InsurancePolicy) error
	GetInsurancePolicyByID(id uuid.UUID) (*model.InsurancePolicy, error)
	GetPatientInsurancePolicies(patientID uuid.UUID) ([]model.InsurancePolicy, error)
	GetActiveInsurancePolicy(patientID uuid.UUID) (*model.InsurancePolicy, error)
	UpdateInsurancePolicy(policy *model.InsurancePolicy) error
	DeactivateInsurancePolicy(id uuid.UUID) error
	ValidateInsurancePolicy(policyID uuid.UUID) error
}

type insurancePolicyService struct {
	insuranceRepo repository.InsurancePolicyRepository
}

func NewInsurancePolicyService(insuranceRepo repository.InsurancePolicyRepository) InsurancePolicyService {
	return &insurancePolicyService{
		insuranceRepo: insuranceRepo,
	}
}

// CreateInsurancePolicy creates a new insurance policy for a patient
func (s *insurancePolicyService) CreateInsurancePolicy(policy *model.InsurancePolicy) error {
	if err := s.insuranceRepo.CreateInsurancePolicy(policy); err != nil {
		return fmt.Errorf("error creating insurance policy: %w", err)
	}
	return nil
}

// GetInsurancePolicyByID retrieves insurance policy by ID
func (s *insurancePolicyService) GetInsurancePolicyByID(id uuid.UUID) (*model.InsurancePolicy, error) {
	policy, err := s.insuranceRepo.GetInsurancePolicyByID(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching insurance policy: %w", err)
	}
	return policy, nil
}

// GetPatientInsurancePolicies retrieves all insurance policies for a patient
func (s *insurancePolicyService) GetPatientInsurancePolicies(patientID uuid.UUID) ([]model.InsurancePolicy, error) {
	policies, err := s.insuranceRepo.GetInsurancePoliciesByPatientID(patientID)
	if err != nil {
		return nil, fmt.Errorf("error fetching insurance policies: %w", err)
	}
	return policies, nil
}

// GetActiveInsurancePolicy retrieves the active insurance policy for a patient
func (s *insurancePolicyService) GetActiveInsurancePolicy(patientID uuid.UUID) (*model.InsurancePolicy, error) {
	policy, err := s.insuranceRepo.GetActiveInsurancePolicyByPatientID(patientID)
	if err != nil {
		return nil, fmt.Errorf("error fetching active insurance policy: %w", err)
	}
	return policy, nil
}

// UpdateInsurancePolicy updates existing insurance policy
func (s *insurancePolicyService) UpdateInsurancePolicy(policy *model.InsurancePolicy) error {
	if err := s.insuranceRepo.UpdateInsurancePolicy(policy); err != nil {
		return fmt.Errorf("error updating insurance policy: %w", err)
	}
	return nil
}

// DeactivateInsurancePolicy deactivates an insurance policy
func (s *insurancePolicyService) DeactivateInsurancePolicy(id uuid.UUID) error {
	if err := s.insuranceRepo.DeactivateInsurancePolicy(id); err != nil {
		return fmt.Errorf("error deactivating insurance policy: %w", err)
	}
	return nil
}

// ValidateInsurancePolicy checks if insurance policy is valid and active
func (s *insurancePolicyService) ValidateInsurancePolicy(policyID uuid.UUID) error {
	policy, err := s.insuranceRepo.GetInsurancePolicyByID(policyID)
	if err != nil {
		return fmt.Errorf("error validating insurance policy: %w", err)
	}

	if policy == nil {
		return fmt.Errorf("insurance policy not found")
	}

	if !policy.IsActive {
		return fmt.Errorf("insurance policy is not active")
	}

	return nil
}
