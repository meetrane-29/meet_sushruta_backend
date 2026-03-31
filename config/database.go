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
		&model.UHID{},
		&model.InsurancePolicy{},
		&model.WaitingListEntry{},
		&model.OPDReceipt{},
	}

	// Pre-migration: add columns with defaults to avoid NOT NULL constraint errors on existing rows
	preFixStatements := []string{
		`ALTER TABLE uhids ADD COLUMN IF NOT EXISTS uh_id varchar(20) DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login bigint`,
		// Fix progress_notes constraints for OPD SOAP notes support
		`ALTER TABLE IF EXISTS progress_notes DROP CONSTRAINT IF EXISTS fk_progress_notes_admission`,
		`ALTER TABLE IF EXISTS progress_notes ALTER COLUMN admission_id DROP NOT NULL`,
		`ALTER TABLE IF EXISTS progress_notes ALTER COLUMN recorded_date DROP NOT NULL`,
	}
	for _, stmt := range preFixStatements {
		if err := DB.Exec(stmt).Error; err != nil {
			log.Printf("Pre-migration warning (non-fatal): %v", err)
		}
	}

	// Migrate each model individually so one bad schema doesn't block everything
	for _, m := range models {
		if err := DB.AutoMigrate(m); err != nil {
			log.Printf("AutoMigrate warning (non-fatal) for %T: %v", m, err)
		}
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
