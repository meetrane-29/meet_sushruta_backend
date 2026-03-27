package handler

import (
	"strconv"

	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PharmacyHandler struct {
	pharmacyService service.PharmacyService
}

func NewPharmacyHandler(pharmacyService service.PharmacyService) *PharmacyHandler {
	return &PharmacyHandler{
		pharmacyService: pharmacyService,
	}
}

// ========== Medicine CRUD Handlers ==========

type CreateMedicineRequest struct {
	Name              string  `json:"name" binding:"required"`
	GenericName       string  `json:"generic_name"`
	Dosage            string  `json:"dosage" binding:"required"`
	Composition       string  `json:"composition"`
	Manufacturer      string  `json:"manufacturer"`
	BatchNumber       string  `json:"batch_number"`
	ExpiryDate        string  `json:"expiry_date"`
	Description       string  `json:"description"`
	Instructions      string  `json:"instructions"`
	SideEffects       string  `json:"side_effects"`
	Contraindications string  `json:"contraindications"`
	StockQuantity     int     `json:"stock_quantity" binding:"min=0"`
	ReorderLevel      int     `json:"reorder_level" binding:"min=1"`
	Price             float64 `json:"price" binding:"required,gt=0"`
}

type UpdateMedicineRequest struct {
	Name              string  `json:"name"`
	GenericName       string  `json:"generic_name"`
	Dosage            string  `json:"dosage"`
	Composition       string  `json:"composition"`
	Manufacturer      string  `json:"manufacturer"`
	BatchNumber       string  `json:"batch_number"`
	ExpiryDate        string  `json:"expiry_date"`
	Description       string  `json:"description"`
	Instructions      string  `json:"instructions"`
	SideEffects       string  `json:"side_effects"`
	Contraindications string  `json:"contraindications"`
	ReorderLevel      int     `json:"reorder_level"`
	Price             float64 `json:"price"`
}

// CreateMedicine creates a new medicine
// POST /api/v1/pharmacy/medicines
func (h *PharmacyHandler) CreateMedicine(c *gin.Context) {
	var req CreateMedicineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	medicine := &model.Medicine{
		Name:              req.Name,
		GenericName:       req.GenericName,
		Dosage:            req.Dosage,
		Composition:       req.Composition,
		Manufacturer:      req.Manufacturer,
		BatchNumber:       req.BatchNumber,
		ExpiryDate:        req.ExpiryDate,
		Description:       req.Description,
		Instructions:      req.Instructions,
		SideEffects:       req.SideEffects,
		Contraindications: req.Contraindications,
		StockQuantity:     req.StockQuantity,
		ReorderLevel:      req.ReorderLevel,
		Price:             req.Price,
		Active:            true,
	}

	err := h.pharmacyService.CreateMedicine(medicine)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, medicine)
}

// ListMedicines retrieves all medicines with pagination and search
// GET /api/v1/pharmacy/medicines?page=1&limit=10&search=paracetamol
func (h *PharmacyHandler) ListMedicines(c *gin.Context) {
	page := 1
	limit := 10
	search := c.Query("search")

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	medicines, total, err := h.pharmacyService.ListMedicines(page, limit, search)
	if err != nil {
		utils.Fail(c, 500, "failed to list medicines")
		return
	}

	utils.OKWithMeta(c, medicines, map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GetMedicine retrieves a single medicine by ID
// GET /api/v1/pharmacy/medicines/:id
func (h *PharmacyHandler) GetMedicine(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid medicine id")
		return
	}

	medicine, err := h.pharmacyService.GetMedicine(id)
	if err != nil {
		utils.Fail(c, 404, "medicine not found")
		return
	}

	utils.OK(c, medicine)
}

// UpdateMedicine updates medicine details
// PATCH /api/v1/pharmacy/medicines/:id
func (h *PharmacyHandler) UpdateMedicine(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid medicine id")
		return
	}

	// Get existing medicine
	existingMedicine, err := h.pharmacyService.GetMedicine(id)
	if err != nil {
		utils.Fail(c, 404, "medicine not found")
		return
	}

	var req UpdateMedicineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	// Update only provided fields
	if req.Name != "" {
		existingMedicine.Name = req.Name
	}
	if req.GenericName != "" {
		existingMedicine.GenericName = req.GenericName
	}
	if req.Dosage != "" {
		existingMedicine.Dosage = req.Dosage
	}
	if req.Composition != "" {
		existingMedicine.Composition = req.Composition
	}
	if req.Manufacturer != "" {
		existingMedicine.Manufacturer = req.Manufacturer
	}
	if req.BatchNumber != "" {
		existingMedicine.BatchNumber = req.BatchNumber
	}
	if req.ExpiryDate != "" {
		existingMedicine.ExpiryDate = req.ExpiryDate
	}
	if req.Description != "" {
		existingMedicine.Description = req.Description
	}
	if req.Instructions != "" {
		existingMedicine.Instructions = req.Instructions
	}
	if req.SideEffects != "" {
		existingMedicine.SideEffects = req.SideEffects
	}
	if req.Contraindications != "" {
		existingMedicine.Contraindications = req.Contraindications
	}
	if req.ReorderLevel > 0 {
		existingMedicine.ReorderLevel = req.ReorderLevel
	}
	if req.Price > 0 {
		existingMedicine.Price = req.Price
	}

	err = h.pharmacyService.UpdateMedicine(existingMedicine)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, existingMedicine)
}

// DeleteMedicine soft deletes a medicine
// DELETE /api/v1/pharmacy/medicines/:id
func (h *PharmacyHandler) DeleteMedicine(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid medicine id")
		return
	}

	err = h.pharmacyService.DeleteMedicine(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{"message": "medicine deleted successfully"})
}

// ========== Stock Management Handlers ==========

type StockUpdateRequest struct {
	Action   string `json:"action" binding:"required,oneof=add remove"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	Reason   string `json:"reason"`
}

// UpdateStock adds or removes medicine stock
// PATCH /api/v1/pharmacy/medicines/:id/stock
func (h *PharmacyHandler) UpdateStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid medicine id")
		return
	}

	var req StockUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: "+err.Error())
		return
	}

	var err2 error
	if req.Action == "add" {
		err2 = h.pharmacyService.AddStock(id, req.Quantity, req.Reason)
	} else {
		err2 = h.pharmacyService.RemoveStock(id, req.Quantity, req.Reason)
	}

	if err2 != nil {
		utils.Fail(c, 400, err2.Error())
		return
	}

	medicine, _ := h.pharmacyService.GetMedicine(id)
	utils.OK(c, gin.H{
		"message":         "stock updated successfully",
		"medicine_id":     id,
		"updated_stock":   medicine.StockQuantity,
		"action":          req.Action,
		"quantity_change": req.Quantity,
	})
}

// GetLowStockMedicines retrieves medicines below reorder level
// GET /api/v1/pharmacy/medicines/low-stock?page=1&limit=10
func (h *PharmacyHandler) GetLowStockMedicines(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	medicines, total, err := h.pharmacyService.GetLowStockMedicines(page, limit)
	if err != nil {
		utils.Fail(c, 500, "failed to retrieve low stock medicines")
		return
	}

	utils.OKWithMeta(c, medicines, map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GetExpiringMedicines retrieves medicines expiring within 30 days
// GET /api/v1/pharmacy/medicines/expiring?page=1&limit=10
func (h *PharmacyHandler) GetExpiringMedicines(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	medicines, total, err := h.pharmacyService.GetExpiringMedicines(page, limit)
	if err != nil {
		utils.Fail(c, 500, "failed to retrieve expiring medicines")
		return
	}

	utils.OKWithMeta(c, medicines, map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// ========== Dispensing Handlers ==========

type DispenseRequest struct {
	PrescriptionID uuid.UUID `json:"prescription_id" binding:"required"`
	StaffID        uuid.UUID `json:"staff_id" binding:"required"`
}

// Dispense dispenses medication from a prescription
// POST /api/v1/pharmacy/dispense
func (h *PharmacyHandler) Dispense(c *gin.Context) {
	var req DispenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	err := h.pharmacyService.Dispense(req.PrescriptionID, req.StaffID)
	if err != nil {
		// Check if it's an insufficient stock error
		if _, ok := err.(service.ErrInsufficientStock); ok {
			utils.Fail(c, 400, err.Error())
			return
		}
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message":         "medication dispensed successfully",
		"prescription_id": req.PrescriptionID,
	})
}

// GetDispenseHistory retrieves dispensing history for a prescription
// GET /api/v1/pharmacy/dispense/:prescription_id
func (h *PharmacyHandler) GetDispenseHistory(c *gin.Context) {
	prescriptionIDStr := c.Param("prescription_id")
	prescriptionID, err := uuid.Parse(prescriptionIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid prescription id")
		return
	}

	history, err := h.pharmacyService.GetDispenseHistory(prescriptionID)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, history)
}
