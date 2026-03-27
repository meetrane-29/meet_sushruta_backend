package handler

import (
	"strconv"

	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BillingHandler struct {
	billingService service.BillingService
}

func NewBillingHandler(billingService service.BillingService) *BillingHandler {
	return &BillingHandler{
		billingService: billingService,
	}
}

// Request/Response structs

type GenerateBillRequest struct {
	PatientID     uuid.UUID                 `json:"patient_id" binding:"required"`
	AppointmentID *uuid.UUID                `json:"appointment_id"`
	Items         []service.BillItemRequest `json:"items" binding:"required,min=1"`
}

type GenerateBillResponse struct {
	Bill       *model.Bill      `json:"bill"`
	Items      []model.BillItem `json:"items"`
	Total      float64          `json:"total"`
	Tax        float64          `json:"tax"`
	GrandTotal float64          `json:"grand_total"`
}

type PayBillRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"` // cash, card, upi, insurance
	TransactionID string `json:"transaction_id"`
}

type ListBillsResponse struct {
	Bills      []model.Bill `json:"bills"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	TotalPages int64        `json:"total_pages"`
}

// GenerateBill creates a bill with line items for a patient
// POST /api/v1/billing/generate
func (h *BillingHandler) GenerateBill(c *gin.Context) {
	var req GenerateBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: patient_id and items are required")
		return
	}

	// Validate items
	if len(req.Items) == 0 {
		utils.Fail(c, 400, "at least one bill item is required")
		return
	}

	bill, err := h.billingService.GenerateBill(req.PatientID, req.AppointmentID, req.Items)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	// Get bill items
	_, billItems, _ := h.billingService.GetBill(bill.ID)

	response := GenerateBillResponse{
		Bill:       bill,
		Items:      billItems,
		Total:      bill.TotalAmount - bill.TaxAmount,
		Tax:        bill.TaxAmount,
		GrandTotal: bill.TotalAmount,
	}

	utils.OK(c, response)
}

// GetBill retrieves a bill by ID with line items
// GET /api/v1/billing/:id
func (h *BillingHandler) GetBill(c *gin.Context) {
	billID := c.Param("id")
	id, err := uuid.Parse(billID)
	if err != nil {
		utils.Fail(c, 400, "invalid bill id")
		return
	}

	bill, billItems, err := h.billingService.GetBill(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	response := GenerateBillResponse{
		Bill:       bill,
		Items:      billItems,
		Total:      bill.TotalAmount - bill.TaxAmount,
		Tax:        bill.TaxAmount,
		GrandTotal: bill.TotalAmount,
	}

	utils.OK(c, response)
}

// PayBill processes payment for a bill
// PATCH /api/v1/billing/:id/pay
func (h *BillingHandler) PayBill(c *gin.Context) {
	billID := c.Param("id")
	id, err := uuid.Parse(billID)
	if err != nil {
		utils.Fail(c, 400, "invalid bill id")
		return
	}

	var req PayBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: payment_method is required")
		return
	}

	bill, err := h.billingService.ProcessPayment(id, req.PaymentMethod, req.TransactionID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, bill)
}

// ListBills retrieves bills with pagination and filtering
// GET /api/v1/billing
func (h *BillingHandler) ListBills(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	// Build filters
	filters := make(map[string]interface{})

	// Filter by status
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// Filter by patient_id
	if patientID := c.Query("patient_id"); patientID != "" {
		if id, err := uuid.Parse(patientID); err == nil {
			filters["patient_id"] = id
		}
	}

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = startDate
		if endDate := c.Query("end_date"); endDate != "" {
			filters["end_date"] = endDate
		}
	}

	bills, total, err := h.billingService.ListBills(page, limit, filters)
	if err != nil {
		utils.Fail(c, 500, "failed to fetch bills")
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	response := ListBillsResponse{
		Bills:      bills,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	utils.OK(c, response)
}
