package service

import (
	"errors"
	"fmt"

	"meet_sushruta/config"
	"meet_sushruta/model"

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
	Dispense(prescriptionID, staffID uuid.UUID) error
	GetDispenseHistory(prescriptionID uuid.UUID) (map[string]interface{}, error)
}

type pharmacyService struct{}

func NewPharmacyService() PharmacyService {
	return &pharmacyService{}
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

		// Update medicine stock
		for _, item := range items {
			medicine := &model.Medicine{}
			if err := tx.Model(medicine).
				Where("id = ?", item.MedicineID).
				Update("stock_quantity", gorm.Expr("stock_quantity - ?", item.Quantity)).Error; err != nil {
				return fmt.Errorf("failed to update medicine stock: %w", err)
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

// GetDispenseHistory retrieves dispensing information for a prescription
func (s *pharmacyService) GetDispenseHistory(prescriptionID uuid.UUID) (map[string]interface{}, error) {
	if prescriptionID == uuid.Nil {
		return nil, errors.New("invalid prescription_id")
	}

	db := config.GetDB()
	prescription := &model.Prescription{}

	if err := db.Where("id = ?", prescriptionID).
		First(prescription).Error; err != nil {
		return nil, errors.New("prescription not found")
	}

	// Get prescription items
	var items []model.PrescriptionItem
	if err := db.Preload("Medicine").
		Where("prescription_id = ?", prescriptionID).
		Find(&items).Error; err != nil {
		return nil, errors.New("failed to get prescription items")
	}

	result := map[string]interface{}{
		"prescription_id": prescription.ID,
		"patient_id":      prescription.PatientID,
		"doctor_id":       prescription.DoctorID,
		"status":          prescription.Status,
		"issued_at":       prescription.IssuedAt,
		"items":           items,
	}

	return result, nil
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
