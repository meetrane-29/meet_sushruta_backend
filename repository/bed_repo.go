package repository

import (
	"database/sql"
	"errors"
	"meet_sushruta/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BedRepositoryImpl struct {
	db *gorm.DB
}

func NewBedRepository(db *gorm.DB) BedRepository {
	return &BedRepositoryImpl{db: db}
}

func (r *BedRepositoryImpl) Create(bed *model.Bed) error {
	return r.db.Create(bed).Error
}

func (r *BedRepositoryImpl) GetByID(id uuid.UUID) (*model.Bed, error) {
	var bed model.Bed
	err := r.db.Preload("Patient").Where("id = ?", id).First(&bed).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bed, nil
}

func (r *BedRepositoryImpl) GetAll(page, limit int) ([]model.Bed, int64, error) {
	var beds []model.Bed
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	err := r.db.Model(&model.Bed{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Patient").
		Offset(offset).
		Limit(limit).
		Find(&beds).Error

	return beds, total, err
}

func (r *BedRepositoryImpl) Update(bed *model.Bed) error {
	return r.db.Save(bed).Error
}

func (r *BedRepositoryImpl) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&model.Bed{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *BedRepositoryImpl) GetByWard(ward string) ([]model.Bed, error) {
	var beds []model.Bed
	err := r.db.Preload("Patient").Where("ward = ? AND deleted_at IS NULL", ward).Find(&beds).Error
	return beds, err
}

func (r *BedRepositoryImpl) GetByStatus(status string) ([]model.Bed, error) {
	var beds []model.Bed
	err := r.db.Preload("Patient").Where("status = ? AND deleted_at IS NULL", status).Find(&beds).Error
	return beds, err
}

// GetByBedType returns beds by bed type (General, ICU, Private, Semi-Private) with patient details
func (r *BedRepositoryImpl) GetByBedType(bedType string) ([]model.Bed, error) {
	var beds []model.Bed
	err := r.db.Preload("Patient").Preload("Patient.User").
		Where("bed_type = ? AND deleted_at IS NULL", bedType).
		Find(&beds).Error
	return beds, err
}

// GetWardStats returns ward-wise bed statistics
func (r *BedRepositoryImpl) GetWardStats() (map[string]interface{}, error) {
	type WardStat struct {
		Ward        string
		Total       int64
		Available   int64
		Occupied    int64
		Maintenance int64
	}

	var wardStats []WardStat
	err := r.db.Raw(`
		SELECT 
			ward,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'available' THEN 1 ELSE 0 END) as available,
			SUM(CASE WHEN status = 'occupied' THEN 1 ELSE 0 END) as occupied,
			SUM(CASE WHEN status IN ('maintenance', 'reserved') THEN 1 ELSE 0 END) as maintenance
		FROM beds
		WHERE deleted_at IS NULL
		GROUP BY ward
	`).Scan(&wardStats).Error

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"wards": wardStats,
	}

	return result, nil
}

// GetBedStats returns overall bed statistics
func (r *BedRepositoryImpl) GetBedStats() (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	// Total beds by type
	var bedTypeStats []map[string]interface{}
	err := r.db.Raw(`
		SELECT 
			bed_type,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'available' THEN 1 ELSE 0 END) as available,
			SUM(CASE WHEN status = 'occupied' THEN 1 ELSE 0 END) as occupied
		FROM beds
		WHERE deleted_at IS NULL
		GROUP BY bed_type
	`).Scan(&bedTypeStats).Error

	if err != nil {
		return nil, err
	}

	// Overall stats
	var totalBeds, availableBeds, occupiedBeds, maintenanceBeds int64
	err = r.db.Model(&model.Bed{}).
		Select("COUNT(*), SUM(CASE WHEN status = 'available' THEN 1 ELSE 0 END), SUM(CASE WHEN status = 'occupied' THEN 1 ELSE 0 END), SUM(CASE WHEN status IN ('maintenance', 'reserved') THEN 1 ELSE 0 END)").
		Row().
		Scan(&totalBeds, &availableBeds, &occupiedBeds, &maintenanceBeds)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	stats["total_beds"] = totalBeds
	stats["available_beds"] = availableBeds
	stats["occupied_beds"] = occupiedBeds
	stats["maintenance_beds"] = maintenanceBeds
	stats["bed_types"] = bedTypeStats
	stats["occupancy_rate"] = float64(occupiedBeds) / float64(totalBeds) * 100

	return stats, nil
}
