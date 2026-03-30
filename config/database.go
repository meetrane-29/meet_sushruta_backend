package config

import (
	"fmt"
	"log"

	"meet_sushruta/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(cfg *Config) error {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("Database connected successfully")

	// Run auto migrations
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func AutoMigrate() error {
	models := []interface{}{
		&model.User{},
		&model.Patient{},
		&model.Doctor{},
		&model.Nurse{},
		&model.DoctorSchedule{},
		&model.Appointment{},
		&model.Vitals{},
		&model.Medicine{},
		&model.Prescription{},
		&model.PrescriptionItem{},
		&model.LabRequest{},
		&model.Bill{},
		&model.BillItem{},
		&model.Bed{},
		&model.AuditLog{},
		&model.Notification{},
		&model.Specialization{},
		&model.Hospital{},
		&model.MedicalEquipment{},
		&model.OperationTheatre{},
		&model.OperationSchedule{},
		&model.DispenseHistory{},
		&model.AdmissionRecord{},
		&model.ProgressNote{},
		&model.NurseInstruction{},
		&model.DischargeSummary{},
		&model.Rating{},
	}

	if err := DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database not initialized. Call InitDatabase() first")
	}
	return DB
}

func CloseDatabase() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
