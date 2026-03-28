package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type nurseInstructionRepository struct{}

func NewNurseInstructionRepository() NurseInstructionRepository {
	return &nurseInstructionRepository{}
}

func (r *nurseInstructionRepository) Create(instruction *model.NurseInstruction) error {
	return config.DB.Create(instruction).Error
}

func (r *nurseInstructionRepository) GetByID(id uuid.UUID) (*model.NurseInstruction, error) {
	var instruction model.NurseInstruction
	err := config.DB.Where("id = ?", id).
		Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		First(&instruction).Error
	if err != nil {
		return nil, err
	}
	return &instruction, nil
}

func (r *nurseInstructionRepository) GetAll(page, limit int) ([]model.NurseInstruction, int64, error) {
	var instructions []model.NurseInstruction
	var total int64

	query := config.DB.Model(&model.NurseInstruction{})

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("issued_date DESC").
		Find(&instructions).Error

	return instructions, total, err
}

func (r *nurseInstructionRepository) GetByAdmissionID(admissionID uuid.UUID, page, limit int) ([]model.NurseInstruction, int64, error) {
	var instructions []model.NurseInstruction
	var total int64

	query := config.DB.Where("admission_id = ?", admissionID)

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("issued_date DESC").
		Find(&instructions).Error

	return instructions, total, err
}

func (r *nurseInstructionRepository) GetByPatientID(patientID uuid.UUID, page, limit int) ([]model.NurseInstruction, int64, error) {
	var instructions []model.NurseInstruction
	var total int64

	query := config.DB.Where("patient_id = ?", patientID)

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("issued_date DESC").
		Find(&instructions).Error

	return instructions, total, err
}

func (r *nurseInstructionRepository) GetByDoctorID(doctorID uuid.UUID, page, limit int) ([]model.NurseInstruction, int64, error) {
	var instructions []model.NurseInstruction
	var total int64

	query := config.DB.Where("doctor_id = ?", doctorID)

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	err := query.Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Offset(offset).
		Limit(limit).
		Order("issued_date DESC").
		Find(&instructions).Error

	return instructions, total, err
}

func (r *nurseInstructionRepository) GetLatestByAdmissionID(admissionID uuid.UUID) (*model.NurseInstruction, error) {
	var instruction model.NurseInstruction
	err := config.DB.Where("admission_id = ?", admissionID).
		Preload("Admission").
		Preload("Patient").
		Preload("Patient.User").
		Preload("Doctor").
		Preload("Doctor.User").
		Order("issued_date DESC").
		First(&instruction).Error
	if err != nil {
		return nil, err
	}
	return &instruction, nil
}

func (r *nurseInstructionRepository) Update(instruction *model.NurseInstruction) error {
	return config.DB.Model(instruction).
		Updates(instruction).Error
}

func (r *nurseInstructionRepository) SoftDelete(id uuid.UUID) error {
	return config.DB.Model(&model.NurseInstruction{}).
		Where("id = ?", id).
		Delete(&model.NurseInstruction{}).Error
}
