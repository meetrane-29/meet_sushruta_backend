package service

import (
	"errors"
	"fmt"
	"time"

	"meet_sushruta/config"
	"meet_sushruta/model"
	"meet_sushruta/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrInsufficientStock represents an error when medicine stock is insufficient
type ErrInsufficientStock struct {
	MedicineName string
}

func (e ErrInsufficientStock) Error() string {
	return fmt.Sprintf("insufficient stock for medicine: %s", e.MedicineName)
}

type PharmacyService interface {
	// Medicine CRUD operations
	CreateMedicine(medicine *model.Medicine) error
	GetMedicine(id uuid.UUID) (*model.Medicine, error)
	ListMedicines(page, limit int, search string) ([]model.Medicine, int64, error)
	UpdateMedicine(medicine *model.Medicine) error
	DeleteMedicine(id uuid.UUID) error

	// Stock management
	AddStock(medicineID uuid.UUID, quantity int, reason string) error
	RemoveStock(medicineID uuid.UUID, quantity int, reason string) error
	GetLowStockMedicines(page, limit int) ([]model.Medicine, int64, error)
	GetExpiringMedicines(page, limit int) ([]model.Medicine, int64, error)

	// Dispensing
	Dispense(prescriptionID, staffID uuid.UUID) error
	GetDispenseHistory(prescriptionID uuid.UUID) ([]model.DispenseHistory, error)
}

type pharmacyService struct {
	medicineRepo repository.MedicineRepository
}

func NewPharmacyService(medicineRepo repository.MedicineRepository) PharmacyService {
	return &pharmacyService{
		medicineRepo: medicineRepo,
	}
}

// CreateMedicine creates a new medicine entry
func (s *pharmacyService) CreateMedicine(medicine *model.Medicine) error {
	if medicine == nil {
		return errors.New("medicine cannot be nil")
	}

	if medicine.Name == "" {
		return errors.New("medicine name is required")
	}

	if medicine.Price <= 0 {
		return errors.New("medicine price must be greater than 0")
	}

	if medicine.ReorderLevel <= 0 {
		medicine.ReorderLevel = 10 // Default reorder level
	}

	return s.medicineRepo.Create(medicine)
}

// GetMedicine retrieves a medicine by ID
func (s *pharmacyService) GetMedicine(id uuid.UUID) (*model.Medicine, error) {
	if id == uuid.Nil {
		return nil, errors.New("medicine id is required")
	}

	return s.medicineRepo.GetByID(id)
}

// ListMedicines retrieves all active medicines with pagination
func (s *pharmacyService) ListMedicines(page, limit int, search string) ([]model.Medicine, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.medicineRepo.GetAll(page, limit, search)
}

// UpdateMedicine updates medicine details
func (s *pharmacyService) UpdateMedicine(medicine *model.Medicine) error {
	if medicine == nil {
		return errors.New("medicine cannot be nil")
	}

	if medicine.ID == uuid.Nil {
		return errors.New("medicine id is required")
	}

	if medicine.Name == "" {
		return errors.New("medicine name is required")
	}

	if medicine.Price <= 0 {
		return errors.New("medicine price must be greater than 0")
	}

	return s.medicineRepo.Update(medicine)
}

// DeleteMedicine soft deletes a medicine
func (s *pharmacyService) DeleteMedicine(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("medicine id is required")
	}

	return s.medicineRepo.SoftDelete(id)
}

// AddStock adds quantity to medicine stock
func (s *pharmacyService) AddStock(medicineID uuid.UUID, quantity int, reason string) error {
	if medicineID == uuid.Nil {
		return errors.New("medicine id is required")
	}

	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	if reason == "" {
		reason = "stock_addition"
	}

	return s.medicineRepo.UpdateStock(medicineID, quantity)
}

// RemoveStock removes quantity from medicine stock
func (s *pharmacyService) RemoveStock(medicineID uuid.UUID, quantity int, reason string) error {
	if medicineID == uuid.Nil {
		return errors.New("medicine id is required")
	}

	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	if reason == "" {
		reason = "stock_removal"
	}

	// Check if sufficient stock exists
	medicine, err := s.medicineRepo.GetByID(medicineID)
	if err != nil {
		return errors.New("medicine not found")
	}

	if medicine.StockQuantity < quantity {
		return ErrInsufficientStock{MedicineName: medicine.Name}
	}

	return s.medicineRepo.UpdateStock(medicineID, -quantity)
}

// GetLowStockMedicines retrieves medicines below reorder level
func (s *pharmacyService) GetLowStockMedicines(page, limit int) ([]model.Medicine, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.medicineRepo.GetLowStock(10, page, limit)
}

// GetExpiringMedicines retrieves medicines expiring within 30 days
func (s *pharmacyService) GetExpiringMedicines(page, limit int) ([]model.Medicine, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.medicineRepo.GetExpiring(30, page, limit)
}

// Dispense handles the medication dispensing process with stock management
func (s *pharmacyService) Dispense(prescriptionID, staffID uuid.UUID) error {
	if prescriptionID == uuid.Nil {
		return errors.New("prescription_id is required")
	}

	if staffID == uuid.Nil {
		return errors.New("staff_id is required")
	}

	db := config.GetDB()

	// Use transaction for atomic operation with proper locking
	err := db.Transaction(func(tx *gorm.DB) error {
		// SELECT prescription FOR UPDATE (lock for update)
		prescription := &model.Prescription{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", prescriptionID).
			Preload("PrescriptionItems").
			First(prescription).Error; err != nil {
			return fmt.Errorf("prescription not found: %w", err)
		}

		// Get prescription items
		var items []model.PrescriptionItem
		if err := tx.Where("prescription_id = ?", prescriptionID).Find(&items).Error; err != nil {
			return fmt.Errorf("failed to get prescription items: %w", err)
		}

		// Check stock for all items first (fail fast)
		for _, item := range items {
			medicine := &model.Medicine{}
			if err := tx.Where("id = ?", item.MedicineID).First(medicine).Error; err != nil {
				return fmt.Errorf("medicine not found: %w", err)
			}

			if medicine.StockQuantity < item.Quantity {
				return ErrInsufficientStock{MedicineName: medicine.Name}
			}
		}

		// Update medicine stock and create dispense history for each item
		now := time.Now().UnixMilli()
		for _, item := range items {
			medicine := &model.Medicine{}
			if err := tx.Where("id = ?", item.MedicineID).First(medicine).Error; err != nil {
				return fmt.Errorf("medicine not found: %w", err)
			}

			// Update stock
			if err := tx.Model(medicine).
				Where("id = ?", item.MedicineID).
				Update("stock_quantity", gorm.Expr("stock_quantity - ?", item.Quantity)).Error; err != nil {
				return fmt.Errorf("failed to update medicine stock: %w", err)
			}

			// Create dispense history
			history := &model.DispenseHistory{
				PrescriptionID: prescriptionID,
				MedicineID:     item.MedicineID,
				Quantity:       item.Quantity,
				UnitPrice:      medicine.Price,
				TotalPrice:     medicine.Price * float64(item.Quantity),
				DispensedBy:    staffID,
				DispensedAt:    now,
				BatchNumber:    medicine.BatchNumber,
				ExpiryDate:     medicine.ExpiryDate,
			}

			if err := tx.Create(history).Error; err != nil {
				return fmt.Errorf("failed to create dispense history: %w", err)
			}
		}

		// Update prescription status to dispensed
		if err := tx.Model(prescription).Update("status", "dispensed").Error; err != nil {
			return fmt.Errorf("failed to update prescription status: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// After transaction commit, check reorder levels and notify admin
	go checkAndNotifyReorderLevels()

	return nil
}

// GetDispenseHistory retrieves dispense records for a prescription
func (s *pharmacyService) GetDispenseHistory(prescriptionID uuid.UUID) ([]model.DispenseHistory, error) {
	if prescriptionID == uuid.Nil {
		return nil, errors.New("prescription_id is required")
	}

	return s.medicineRepo.GetDispenseHistoryByPrescriptionID(prescriptionID)
}

// checkAndNotifyReorderLevels checks medicine stock levels and notifies if below threshold
func checkAndNotifyReorderLevels() {
	// TODO: Implement reorder level notification logic
	// This could be done via email, SMS, or in-app notification to admin
	fmt.Println("Checking medicine reorder levels...")

	db := config.GetDB()
	var medicines []model.Medicine

	// Find medicines below reorder level
	if err := db.Where("stock_quantity <= reorder_level AND active = ?", true).Find(&medicines).Error; err != nil {
		fmt.Printf("Error checking reorder levels: %v\n", err)
		return
	}

	if len(medicines) > 0 {
		fmt.Printf("Medicines below reorder level: %d\n", len(medicines))
		for _, med := range medicines {
			fmt.Printf("Medicine %s: stock=%d, reorder_level=%d\n", med.Name, med.StockQuantity, med.ReorderLevel)
		}
	}
}
