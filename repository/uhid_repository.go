package repository

import (
	"errors"
	"fmt"
	"meet_sushruta/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UHIDRepository interface {
	CreateUHID(uhid *model.UHID) error
	GetUHIDByPatientID(patientID uuid.UUID) (*model.UHID, error)
	GetUHIDByUHIDString(uhidString string) (*model.UHID, error)
	GetNextSequenceNumber(hospitalCode string) (int64, error)
}

type uhidRepository struct {
	db *gorm.DB
}

func NewUHIDRepository(db *gorm.DB) UHIDRepository {
	return &uhidRepository{db: db}
}

// CreateUHID creates a new UHID entry
func (r *uhidRepository) CreateUHID(uhid *model.UHID) error {
	if err := r.db.Create(uhid).Error; err != nil {
		return fmt.Errorf("error creating UHID: %w", err)
	}
	return nil
}

// GetUHIDByPatientID retrieves UHID by patient ID
func (r *uhidRepository) GetUHIDByPatientID(patientID uuid.UUID) (*model.UHID, error) {
	var uhid model.UHID
	if err := r.db.Where("patient_id = ? AND is_active = true", patientID).
		First(&uhid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching UHID: %w", err)
	}
	return &uhid, nil
}

// GetUHIDByUHIDString retrieves UHID by UHID string
func (r *uhidRepository) GetUHIDByUHIDString(uhidString string) (*model.UHID, error) {
	var uhid model.UHID
	if err := r.db.Where("uhid = ? AND is_active = true", uhidString).
		First(&uhid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching UHID: %w", err)
	}
	return &uhid, nil
}

// GetNextSequenceNumber gets the next sequence number for UHID generation
func (r *uhidRepository) GetNextSequenceNumber(hospitalCode string) (int64, error) {
	var maxSeqNum int64
	now := time.Now()
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).UnixMilli()
	endOfYear := time.Date(now.Year(), 12, 31, 23, 59, 59, 999999999, now.Location()).UnixMilli()

	if err := r.db.Model(&model.UHID{}).
		Where("hospital_code = ? AND issued_date >= ? AND issued_date <= ?", hospitalCode, startOfYear, endOfYear).
		Select("COALESCE(MAX(sequence_number), 0)").
		Scan(&maxSeqNum).Error; err != nil {
		return 0, fmt.Errorf("error getting sequence number: %w", err)
	}

	return maxSeqNum + 1, nil
}
