package handler

import (
	"meet_sushruta/service"
	"meet_sushruta/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OPDReceiptHandler struct {
	receiptService service.OPDReceiptService
}

func NewOPDReceiptHandler(receiptService service.OPDReceiptService) *OPDReceiptHandler {
	return &OPDReceiptHandler{
		receiptService: receiptService,
	}
}

type GenerateOPDReceiptRequest struct {
	BillID           uuid.UUID  `json:"bill_id" binding:"required"`
	AppointmentID    *uuid.UUID `json:"appointment_id"`
	ConsultationFee  float64    `json:"consultation_fee" binding:"required"`
	InsuranceCovered float64    `json:"insurance_covered"`
	PaymentMethod    string     `json:"payment_method" binding:"required"` // cash, card, upi, insurance
	TransactionRef   string     `json:"transaction_ref"`
	Notes            string     `json:"notes"`
}

// GenerateOPDReceipt generates OPD receipt for patient
// POST /api/v1/opd-receipt/generate
func (h *OPDReceiptHandler) GenerateOPDReceipt(c *gin.Context) {
	var req GenerateOPDReceiptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	receipt, err := h.receiptService.GenerateOPDReceipt(
		req.BillID,
		req.AppointmentID,
		req.ConsultationFee,
		req.InsuranceCovered,
		req.PaymentMethod,
	)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	if req.TransactionRef != "" {
		receipt.TransactionRef = req.TransactionRef
	}
	if req.Notes != "" {
		receipt.Notes = req.Notes
	}

	utils.OK(c, receipt)
}

// GetOPDReceipt retrieves OPD receipt by ID
// GET /api/v1/opd-receipt/:id
func (h *OPDReceiptHandler) GetOPDReceipt(c *gin.Context) {
	receiptIDStr := c.Param("id")
	receiptID, err := uuid.Parse(receiptIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid receipt ID")
		return
	}

	receipt, err := h.receiptService.GetOPDReceiptByID(receiptID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	if receipt == nil {
		utils.Fail(c, 404, "receipt not found")
		return
	}

	utils.OK(c, receipt)
}

// GetOPDReceiptByNumber retrieves OPD receipt by receipt number
// GET /api/v1/opd-receipt/number/:receipt_number
func (h *OPDReceiptHandler) GetOPDReceiptByNumber(c *gin.Context) {
	receiptNumber := c.Param("receipt_number")
	if receiptNumber == "" {
		utils.Fail(c, 400, "receipt number is required")
		return
	}

	receipt, err := h.receiptService.GetOPDReceiptByReceiptNumber(receiptNumber)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	if receipt == nil {
		utils.Fail(c, 404, "receipt not found")
		return
	}

	utils.OK(c, receipt)
}

// GetPatientOPDReceipts retrieves all OPD receipts for a patient
// GET /api/v1/opd-receipt/patient/:patient_id
func (h *OPDReceiptHandler) GetPatientOPDReceipts(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient ID")
		return
	}

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	receipts, err := h.receiptService.GetPatientOPDReceipts(patientID, limit, offset)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"receipts": receipts,
		"count":    len(receipts),
	})
}

// PrintOPDReceipt marks receipt as printed (increments print count)
// POST /api/v1/opd-receipt/:id/print
func (h *OPDReceiptHandler) PrintOPDReceipt(c *gin.Context) {
	receiptIDStr := c.Param("id")
	receiptID, err := uuid.Parse(receiptIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid receipt ID")
		return
	}

	if err := h.receiptService.UpdatePrintCount(receiptID); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{"message": "receipt printed successfully"})
}

// GetReceiptSummary retrieves receipt summary for display
// GET /api/v1/opd-receipt/:id/summary
func (h *OPDReceiptHandler) GetReceiptSummary(c *gin.Context) {
	receiptIDStr := c.Param("id")
	receiptID, err := uuid.Parse(receiptIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid receipt ID")
		return
	}

	receipt, err := h.receiptService.GetOPDReceiptByID(receiptID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	if receipt == nil {
		utils.Fail(c, 404, "receipt not found")
		return
	}

	summary := gin.H{
		"receipt_number":    receipt.ReceiptNumber,
		"receipt_date":      receipt.ReceiptDate,
		"consultation_fee":  receipt.ConsultationFee,
		"insurance_covered": receipt.InsuranceCovered,
		"patient_payable":   receipt.PatientPayable,
		"paid_amount":       receipt.PaidAmount,
		"payment_method":    receipt.PaymentMethod,
		"printed_count":     receipt.PrintedCount,
		"transaction_ref":   receipt.TransactionRef,
	}

	utils.OK(c, summary)
}
