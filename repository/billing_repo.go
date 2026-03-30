package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type billingRepository struct{}

func NewBillingRepository() BillingRepository {
	return &billingRepository{}
}

func (r *billingRepository) Create(bill *model.Bill) error {
	return config.DB.Create(bill).Error
}

func (r *billingRepository) GetByID(id uuid.UUID) (*model.Bill, error) {
	var bill model.Bill
	err := config.DB.Where("id = ?", id).
		Preload("Patient").
		Preload("Appointment").
		Preload("Appointment.Doctor").
		First(&bill).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

func (r *billingRepository) GetAll(page, limit int) ([]model.Bill, int64, error) {
	var bills []model.Bill
	var total int64

	// Get total count with a separate query
	if err := config.DB.Model(&model.Bill{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with a fresh query
	err := config.DB.Preload("Patient").
		Preload("Appointment").
		Preload("Appointment.Doctor").
		Offset(offset).
		Limit(limit).
		Order("bill_date DESC, created_at DESC").
		Find(&bills).Error

	return bills, total, err
}

func (r *billingRepository) GetByPatientID(patientID uuid.UUID, page, limit int) ([]model.Bill, int64, error) {
	var bills []model.Bill
	var total int64

	// Get total count with a separate query
	if err := config.DB.Model(&model.Bill{}).Where("patient_id = ?", patientID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with a fresh query
	err := config.DB.Where("patient_id = ?", patientID).
		Preload("Patient").
		Preload("Appointment").
		Preload("Appointment.Doctor").
		Offset(offset).
		Limit(limit).
		Order("bill_date DESC, created_at DESC").
		Find(&bills).Error

	return bills, total, err
}

func (r *billingRepository) GetByStatus(status string, page, limit int) ([]model.Bill, int64, error) {
	var bills []model.Bill
	var total int64

	// Get total count with a separate query
	if err := config.DB.Model(&model.Bill{}).Where("payment_status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with a fresh query
	err := config.DB.Where("payment_status = ?", status).
		Preload("Patient").
		Preload("Appointment").
		Preload("Appointment.Doctor").
		Offset(offset).
		Limit(limit).
		Order("bill_date DESC, created_at DESC").
		Find(&bills).Error

	return bills, total, err
}

func (r *billingRepository) GetByDateRange(startDate, endDate string, page, limit int) ([]model.Bill, int64, error) {
	var bills []model.Bill
	var total int64

	// Get total count with a separate query
	if err := config.DB.Model(&model.Bill{}).Where("bill_date >= ? AND bill_date <= ?", startDate, endDate).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with a fresh query
	err := config.DB.Where("bill_date >= ? AND bill_date <= ?", startDate, endDate).
		Preload("Patient").
		Preload("Appointment").
		Preload("Appointment.Doctor").
		Offset(offset).
		Limit(limit).
		Order("bill_date DESC, created_at DESC").
		Find(&bills).Error

	return bills, total, err
}

func (r *billingRepository) Update(bill *model.Bill) error {
	return config.DB.Model(bill).Updates(bill).Error
}

func (r *billingRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.Bill{}).
		Where("id = ?", id).
		Delete(&model.Bill{}).Error
}

func (r *billingRepository) CreateBillItem(billItem *model.BillItem) error {
	return config.DB.Create(billItem).Error
}

func (r *billingRepository) GetBillItems(billID uuid.UUID) ([]model.BillItem, error) {
	var items []model.BillItem
	err := config.DB.Where("bill_id = ?", billID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}
