package repository

import (
	"database/sql"
	"errors"
	"meet_sushruta/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MedicalEquipmentRepositoryImpl struct {
	db *gorm.DB
}

func NewMedicalEquipmentRepository(db *gorm.DB) MedicalEquipmentRepository {
	return &MedicalEquipmentRepositoryImpl{db: db}
}

func (r *MedicalEquipmentRepositoryImpl) Create(equipment *model.MedicalEquipment) error {
	return r.db.Create(equipment).Error
}

func (r *MedicalEquipmentRepositoryImpl) GetByID(id uuid.UUID) (*model.MedicalEquipment, error) {
	var equipment model.MedicalEquipment
	err := r.db.Where("id = ?", id).First(&equipment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &equipment, nil
}

func (r *MedicalEquipmentRepositoryImpl) GetAll(page, limit int) ([]model.MedicalEquipment, int64, error) {
	var equipment []model.MedicalEquipment
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	err := r.db.Model(&model.MedicalEquipment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).
		Limit(limit).
		Find(&equipment).Error

	return equipment, total, err
}

func (r *MedicalEquipmentRepositoryImpl) Update(equipment *model.MedicalEquipment) error {
	return r.db.Save(equipment).Error
}

func (r *MedicalEquipmentRepositoryImpl) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&model.MedicalEquipment{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *MedicalEquipmentRepositoryImpl) GetByStatus(status string) ([]model.MedicalEquipment, error) {
	var equipment []model.MedicalEquipment
	err := r.db.Where("status = ? AND deleted_at IS NULL", status).Find(&equipment).Error
	return equipment, err
}

func (r *MedicalEquipmentRepositoryImpl) GetCriticalEquipment() ([]model.MedicalEquipment, error) {
	var equipment []model.MedicalEquipment
	err := r.db.Where("critical_equipment = true AND deleted_at IS NULL", true).Find(&equipment).Error
	return equipment, err
}

// GetEquipmentStats returns overall equipment statistics
func (r *MedicalEquipmentRepositoryImpl) GetEquipmentStats() (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	// Equipment status breakdown
	var statusStats []map[string]interface{}
	err := r.db.Raw(`
		SELECT 
			status,
			COUNT(*) as count
		FROM medical_equipments
		WHERE deleted_at IS NULL
		GROUP BY status
	`).Scan(&statusStats).Error

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Equipment type breakdown
	var typeStats []map[string]interface{}
	err = r.db.Raw(`
		SELECT 
			equipment_type,
			COUNT(*) as count,
			SUM(CASE WHEN status = 'working' THEN 1 ELSE 0 END) as working,
			SUM(CASE WHEN status != 'working' THEN 1 ELSE 0 END) as not_working
		FROM medical_equipments
		WHERE deleted_at IS NULL
		GROUP BY equipment_type
	`).Scan(&typeStats).Error

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Overall stats
	var totalEquipment, workingEquipment, underMaintenance, repair, retired int64
	err = r.db.Model(&model.MedicalEquipment{}).
		Select("COUNT(*), SUM(CASE WHEN status = 'working' THEN 1 ELSE 0 END), SUM(CASE WHEN status = 'under_maintenance' THEN 1 ELSE 0 END), SUM(CASE WHEN status = 'repair' THEN 1 ELSE 0 END), SUM(CASE WHEN status = 'retired' THEN 1 ELSE 0 END)").
		Row().
		Scan(&totalEquipment, &workingEquipment, &underMaintenance, &repair, &retired)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	stats["total_equipment"] = totalEquipment
	stats["working"] = workingEquipment
	stats["under_maintenance"] = underMaintenance
	stats["repair"] = repair
	stats["retired"] = retired
	stats["status_breakdown"] = statusStats
	stats["type_breakdown"] = typeStats
	if totalEquipment > 0 {
		stats["working_percentage"] = float64(workingEquipment) / float64(totalEquipment) * 100
	}

	return stats, nil
}
