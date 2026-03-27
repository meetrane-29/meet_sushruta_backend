package repository

import (
	"fmt"
	"time"

	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type medicineRepository struct{}

func NewMedicineRepository() MedicineRepository {
	return &medicineRepository{}
}

func (r *medicineRepository) Create(medicine *model.Medicine) error {
	return config.DB.Create(medicine).Error
}

func (r *medicineRepository) GetByID(id uuid.UUID) (*model.Medicine, error) {
	var medicine model.Medicine
	err := config.DB.Where("id = ?", id).First(&medicine).Error
	if err != nil {
		return nil, err
	}
	return &medicine, nil
}

func (r *medicineRepository) GetAll(page, limit int, search string) ([]model.Medicine, int64, error) {
	var medicines []model.Medicine
	var total int64

	query := config.DB.Model(&model.Medicine{})

	// Apply search filter
	if search != "" {
		query = query.Where("name ILIKE ? OR generic_name ILIKE ? OR manufacturer ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Apply active filter
	query = query.Where("active = ?", true)

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results ordered by name
	err := query.Offset(offset).
		Limit(limit).
		Order("name ASC").
		Find(&medicines).Error

	return medicines, total, err
}

func (r *medicineRepository) Update(medicine *model.Medicine) error {
	return config.DB.Model(medicine).Updates(medicine).Error
}

func (r *medicineRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.Medicine{}).
		Where("id = ?", id).
		Delete(&model.Medicine{}).Error
}

// GetLowStock retrieves medicines where stock is below or equal to reorder level
func (r *medicineRepository) GetLowStock(reorderLevel int, page, limit int) ([]model.Medicine, int64, error) {
	var medicines []model.Medicine
	var total int64

	query := config.DB.Model(&model.Medicine{}).
		Where("stock_quantity <= ? AND active = ?", reorderLevel, true)

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).
		Limit(limit).
		Order("stock_quantity ASC").
		Find(&medicines).Error

	return medicines, total, err
}

// GetExpiring retrieves medicines expiring within daysThreshold days
func (r *medicineRepository) GetExpiring(daysThreshold int, page, limit int) ([]model.Medicine, int64, error) {
	var medicines []model.Medicine
	var total int64

	now := time.Now()
	expiryBefore := now.AddDate(0, 0, daysThreshold).Format("2006-01-02")

	query := config.DB.Model(&model.Medicine{}).
		Where("expiry_date <= ? AND expiry_date >= ? AND active = ?",
			expiryBefore, now.Format("2006-01-02"), true)

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).
		Limit(limit).
		Order("expiry_date ASC").
		Find(&medicines).Error

	return medicines, total, err
}

// UpdateStock updates the medicine stock quantity
func (r *medicineRepository) UpdateStock(medicineID uuid.UUID, quantity int) error {
	if quantity == 0 {
		return fmt.Errorf("quantity cannot be zero")
	}

	result := config.DB.Model(&model.Medicine{}).
		Where("id = ?", medicineID).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("medicine not found")
	}

	return nil
}

// CreateDispenseHistory records a dispense transaction
func (r *medicineRepository) CreateDispenseHistory(history *model.DispenseHistory) error {
	return config.DB.Create(history).Error
}

// GetDispenseHistoryByPrescriptionID retrieves all dispense records for a prescription
func (r *medicineRepository) GetDispenseHistoryByPrescriptionID(prescriptionID uuid.UUID) ([]model.DispenseHistory, error) {
	var histories []model.DispenseHistory
	err := config.DB.
		Preload("Medicine").
		Preload("DispensingUser").
		Where("prescription_id = ?", prescriptionID).
		Order("dispensed_at DESC").
		Find(&histories).Error
	return histories, err
}
