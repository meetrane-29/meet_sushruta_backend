package handler

import (
	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InsurancePolicyHandler struct {
	insuranceService service.InsurancePolicyService
}

func NewInsurancePolicyHandler(insuranceService service.InsurancePolicyService) *InsurancePolicyHandler {
	return &InsurancePolicyHandler{
		insuranceService: insuranceService,
	}
}

type CreateInsurancePolicyRequest struct {
	PatientID         uuid.UUID `json:"patient_id" binding:"required"`
	ProviderName      string    `json:"provider_name" binding:"required"`
	PolicyNumber      string    `json:"policy_number" binding:"required"`
	MemberID          string    `json:"member_id"`
	CoveragePlanName  string    `json:"coverage_plan_name"`
	CoverageAmount    float64   `json:"coverage_amount" binding:"required"`
	CopayPercentage   float64   `json:"copay_percentage"`
	DeductibleAmount  float64   `json:"deductible_amount"`
	ValidFrom         string    `json:"valid_from" binding:"required"` // YYYY-MM-DD
	ValidUpto         string    `json:"valid_upto" binding:"required"` // YYYY-MM-DD
	NomineeNames      string    `json:"nominee_names"`
	EmployerName      string    `json:"employer_name"`
	ProviderContactNo string    `json:"provider_contact_no"`
	ProviderWebsite   string    `json:"provider_website"`
	Notes             string    `json:"notes"`
}

type UpdateInsurancePolicyRequest struct {
	ProviderName      string  `json:"provider_name"`
	CoveragePlanName  string  `json:"coverage_plan_name"`
	CoverageAmount    float64 `json:"coverage_amount"`
	CopayPercentage   float64 `json:"copay_percentage"`
	DeductibleAmount  float64 `json:"deductible_amount"`
	ValidUpto         string  `json:"valid_upto"` // YYYY-MM-DD
	IsActive          bool    `json:"is_active"`
	ProviderContactNo string  `json:"provider_contact_no"`
	ProviderWebsite   string  `json:"provider_website"`
	Notes             string  `json:"notes"`
}

// CreateInsurancePolicy creates a new insurance policy for patient
// POST /api/v1/insurance/policies
func (h *InsurancePolicyHandler) CreateInsurancePolicy(c *gin.Context) {
	var req CreateInsurancePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	policy := &model.InsurancePolicy{
		ID:                uuid.New(),
		PatientID:         req.PatientID,
		ProviderName:      req.ProviderName,
		PolicyNumber:      req.PolicyNumber,
		MemberID:          req.MemberID,
		CoveragePlanName:  req.CoveragePlanName,
		CoverageAmount:    req.CoverageAmount,
		CopayPercentage:   req.CopayPercentage,
		DeductibleAmount:  req.DeductibleAmount,
		ValidFrom:         req.ValidFrom,
		ValidUpto:         req.ValidUpto,
		NomineeNames:      req.NomineeNames,
		EmployerName:      req.EmployerName,
		IsActive:          true,
		ProviderContactNo: req.ProviderContactNo,
		ProviderWebsite:   req.ProviderWebsite,
		Notes:             req.Notes,
	}

	if err := h.insuranceService.CreateInsurancePolicy(policy); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, policy)
}

// GetInsurancePoliciesByPatient retrieves all insurance policies for a patient
// GET /api/v1/insurance/patient/:patient_id
func (h *InsurancePolicyHandler) GetInsurancePoliciesByPatient(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient ID")
		return
	}

	policies, err := h.insuranceService.GetPatientInsurancePolicies(patientID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, policies)
}

// GetActiveInsurancePolicy retrieves active insurance policy for patient
// GET /api/v1/insurance/active/:patient_id
func (h *InsurancePolicyHandler) GetActiveInsurancePolicy(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient ID")
		return
	}

	policy, err := h.insuranceService.GetActiveInsurancePolicy(patientID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	if policy == nil {
		utils.Fail(c, 404, "no active insurance policy found")
		return
	}

	utils.OK(c, policy)
}

// UpdateInsurancePolicy updates insurance policy details
// PUT /api/v1/insurance/policies/:policy_id
func (h *InsurancePolicyHandler) UpdateInsurancePolicy(c *gin.Context) {
	policyIDStr := c.Param("policy_id")
	policyID, err := uuid.Parse(policyIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid policy ID")
		return
	}

	var req UpdateInsurancePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	policy, err := h.insuranceService.GetInsurancePolicyByID(policyID)
	if err != nil || policy == nil {
		utils.Fail(c, 404, "policy not found")
		return
	}

	// Update provided fields
	if req.ProviderName != "" {
		policy.ProviderName = req.ProviderName
	}
	if req.CoveragePlanName != "" {
		policy.CoveragePlanName = req.CoveragePlanName
	}
	if req.CoverageAmount > 0 {
		policy.CoverageAmount = req.CoverageAmount
	}
	if req.CopayPercentage >= 0 {
		policy.CopayPercentage = req.CopayPercentage
	}
	if req.ValidUpto != "" {
		policy.ValidUpto = req.ValidUpto
	}
	if req.ProviderContactNo != "" {
		policy.ProviderContactNo = req.ProviderContactNo
	}
	if req.ProviderWebsite != "" {
		policy.ProviderWebsite = req.ProviderWebsite
	}
	if req.Notes != "" {
		policy.Notes = req.Notes
	}

	if err := h.insuranceService.UpdateInsurancePolicy(policy); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, policy)
}

// DeactivateInsurancePolicy deactivates an insurance policy
// DELETE /api/v1/insurance/policies/:policy_id
func (h *InsurancePolicyHandler) DeactivateInsurancePolicy(c *gin.Context) {
	policyIDStr := c.Param("policy_id")
	policyID, err := uuid.Parse(policyIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid policy ID")
		return
	}

	if err := h.insuranceService.DeactivateInsurancePolicy(policyID); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{"message": "insurance policy deactivated successfully"})
}
