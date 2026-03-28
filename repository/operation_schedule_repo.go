package repository

import (
	"errors"
	"fmt"
	"meet_sushruta/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OperationTheatreRepositoryImpl struct {
	db *gorm.DB
}

func NewOperationTheatreRepository(db *gorm.DB) OperationTheatreRepository {
	return &OperationTheatreRepositoryImpl{db: db}
}

func (r *OperationTheatreRepositoryImpl) Create(theatre *model.OperationTheatre) error {
	return r.db.Create(theatre).Error
}

func (r *OperationTheatreRepositoryImpl) GetByID(id uuid.UUID) (*model.OperationTheatre, error) {
	var theatre model.OperationTheatre
	err := r.db.Where("id = ?", id).First(&theatre).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &theatre, nil
}

func (r *OperationTheatreRepositoryImpl) GetAll(page, limit int) ([]model.OperationTheatre, int64, error) {
	var theatres []model.OperationTheatre
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	err := r.db.Model(&model.OperationTheatre{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).
		Limit(limit).
		Find(&theatres).Error

	return theatres, total, err
}

func (r *OperationTheatreRepositoryImpl) Update(theatre *model.OperationTheatre) error {
	return r.db.Save(theatre).Error
}

func (r *OperationTheatreRepositoryImpl) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&model.OperationTheatre{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *OperationTheatreRepositoryImpl) GetByStatus(status string) ([]model.OperationTheatre, error) {
	var theatres []model.OperationTheatre
	err := r.db.Where("status = ? AND deleted_at IS NULL", status).Find(&theatres).Error
	return theatres, err
}

// OperationScheduleRepositoryImpl
type OperationScheduleRepositoryImpl struct {
	db *gorm.DB
}

func NewOperationScheduleRepository(db *gorm.DB) OperationScheduleRepository {
	return &OperationScheduleRepositoryImpl{db: db}
}

func (r *OperationScheduleRepositoryImpl) Create(operation *model.OperationSchedule) error {
	return r.db.Create(operation).Error
}

func (r *OperationScheduleRepositoryImpl) GetByID(id uuid.UUID) (*model.OperationSchedule, error) {
	var operation model.OperationSchedule
	err := r.db.
		Preload("Theatre").
		Preload("Patient").
		Preload("SurgeonDoctor").
		Preload("AnesthesiaDoctor").
		Where("id = ?", id).
		First(&operation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &operation, nil
}

func (r *OperationScheduleRepositoryImpl) GetAll(page, limit int) ([]model.OperationSchedule, int64, error) {
	var operations []model.OperationSchedule
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	err := r.db.Model(&model.OperationSchedule{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Preload("Theatre").
		Preload("Patient").
		Preload("SurgeonDoctor").
		Preload("AnesthesiaDoctor").
		Offset(offset).
		Limit(limit).
		Order("operation_date ASC, operation_time ASC").
		Find(&operations).Error

	return operations, total, err
}

func (r *OperationScheduleRepositoryImpl) Update(operation *model.OperationSchedule) error {
	return r.db.Save(operation).Error
}

func (r *OperationScheduleRepositoryImpl) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&model.OperationSchedule{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *OperationScheduleRepositoryImpl) GetByStatus(status string) ([]model.OperationSchedule, error) {
	var operations []model.OperationSchedule
	err := r.db.
		Preload("Theatre").
		Preload("Patient").
		Preload("SurgeonDoctor").
		Preload("AnesthesiaDoctor").
		Where("status = ? AND deleted_at IS NULL", status).
		Order("operation_date ASC, operation_time ASC").
		Find(&operations).Error
	return operations, err
}

func (r *OperationScheduleRepositoryImpl) GetByDate(date string) ([]model.OperationSchedule, error) {
	var operations []model.OperationSchedule
	err := r.db.
		Preload("Theatre").
		Preload("Patient").
		Preload("SurgeonDoctor").
		Preload("AnesthesiaDoctor").
		Where("operation_date = ? AND deleted_at IS NULL", date).
		Order("operation_time ASC").
		Find(&operations).Error
	return operations, err
}

func (r *OperationScheduleRepositoryImpl) GetByTheatreID(theatreID uuid.UUID) ([]model.OperationSchedule, error) {
	var operations []model.OperationSchedule
	err := r.db.
		Preload("Theatre").
		Preload("Patient").
		Preload("SurgeonDoctor").
		Preload("AnesthesiaDoctor").
		Where("theatre_id = ? AND deleted_at IS NULL", theatreID).
		Order("operation_date ASC, operation_time ASC").
		Find(&operations).Error
	return operations, err
}

func (r *OperationScheduleRepositoryImpl) GetUpcomingOperations(daysAhead int) ([]model.OperationSchedule, error) {
	var operations []model.OperationSchedule

	today := time.Now().Format("2006-01-02")
	futureDate := time.Now().AddDate(0, 0, daysAhead).Format("2006-01-02")

	err := r.db.
		Preload("Theatre").
		Preload("Patient").
		Preload("SurgeonDoctor").
		Preload("AnesthesiaDoctor").
		Where("operation_date >= ? AND operation_date <= ? AND status IN ('scheduled', 'in_progress') AND deleted_at IS NULL", today, futureDate).
		Order("operation_date ASC, operation_time ASC").
		Find(&operations).Error

	return operations, err
}

// GetOperationStats returns OT statistics
func (r *OperationScheduleRepositoryImpl) GetOperationStats() (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	today := time.Now().Format("2006-01-02")

	// Today's operations
	var todayOperations int64
	r.db.Model(&model.OperationSchedule{}).
		Where("operation_date = ? AND status IN ('scheduled', 'in_progress')", today).
		Count(&todayOperations)

	// Upcoming operations (next 7 days)
	futureDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	var upcomingOperations int64
	r.db.Model(&model.OperationSchedule{}).
		Where("operation_date > ? AND operation_date <= ? AND status IN ('scheduled')", today, futureDate).
		Count(&upcomingOperations)

	// Status breakdown
	var statusBreakdown []map[string]interface{}
	r.db.Raw(fmt.Sprintf(`
		SELECT 
			status,
			COUNT(*) as count
		FROM operation_schedules
		WHERE deleted_at IS NULL
		GROUP BY status
	`)).Scan(&statusBreakdown)

	// Operation type breakdown
	var typeBreakdown []map[string]interface{}
	r.db.Raw(fmt.Sprintf(`
		SELECT 
			operation_type,
			COUNT(*) as count
		FROM operation_schedules
		WHERE deleted_at IS NULL
		GROUP BY operation_type
	`)).Scan(&typeBreakdown)

	stats["today_operations"] = todayOperations
	stats["upcoming_operations"] = upcomingOperations
	stats["status_breakdown"] = statusBreakdown
	stats["type_breakdown"] = typeBreakdown

	return stats, nil
}
