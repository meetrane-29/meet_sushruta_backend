package repository

import (
	"meet_sushruta/config"
	"meet_sushruta/model"

	"github.com/google/uuid"
)

type doctorScheduleRepository struct{}

func NewDoctorScheduleRepository() DoctorScheduleRepository {
	return &doctorScheduleRepository{}
}

func (r *doctorScheduleRepository) GetByDoctorID(doctorID uuid.UUID) ([]model.DoctorSchedule, error) {
	var schedules []model.DoctorSchedule
	err := config.DB.Where("doctor_id = ?", doctorID).
		Preload("Doctor").
		Find(&schedules).Error

	return schedules, err
}

func (r *doctorScheduleRepository) GetByDoctorAndDay(doctorID uuid.UUID, dayOfWeek string) (*model.DoctorSchedule, error) {
	var schedule model.DoctorSchedule
	err := config.DB.Where("doctor_id = ? AND day_of_week = ?", doctorID, dayOfWeek).
		Preload("Doctor").
		First(&schedule).Error

	if err != nil {
		return nil, err
	}

	return &schedule, nil
}
