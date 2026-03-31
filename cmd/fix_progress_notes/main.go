package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection - UPDATE THESE VALUES
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "meet_sushruta"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}

	fmt.Println("✓ Connected to database")
	fmt.Printf("  Host: %s | DB: %s\n\n", host, dbname)

	// Fix statements
	fixes := []string{
		// Drop old foreign key if exists
		`ALTER TABLE progress_notes DROP CONSTRAINT IF EXISTS fk_progress_notes_admission CASCADE`,

		// Make admission_id nullable
		`ALTER TABLE progress_notes ALTER COLUMN admission_id DROP NOT NULL`,

		// Make recorded_date nullable
		`ALTER TABLE progress_notes ALTER COLUMN recorded_date DROP NOT NULL`,

		// Recreate foreign key with proper cascade
		`ALTER TABLE progress_notes
		 ADD CONSTRAINT fk_progress_notes_admission
		 FOREIGN KEY(admission_id) REFERENCES admission_records(id)
		 ON DELETE SET NULL ON UPDATE CASCADE`,
	}

	fmt.Println("Running migrations...\n")

	for i, fix := range fixes {
		fmt.Printf("[%d/%d] Executing: %s...\n", i+1, len(fixes), fix[:60]+"...")
		if err := db.Exec(fix).Error; err != nil {
			log.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Println("✅ Success\n")
		}
	}

	// Verify
	fmt.Println("Verifying schema...")
	type ColumnInfo struct {
		ColumnName    string
		DataType      string
		IsNullable    string
		ColumnDefault interface{}
	}

	var columns []ColumnInfo
	db.Raw(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_name = 'progress_notes'
		AND column_name IN ('admission_id', 'recorded_date', 'patient_id', 'doctor_id')
		ORDER BY ordinal_position
	`).Scan(&columns)

	fmt.Println("\nColumn Schema:")
	fmt.Println("─────────────────────────────────────────────────────")
	for _, col := range columns {
		nullable := "✓ NULL"
		if col.IsNullable == "NO" {
			nullable = "❌ NOT NULL"
		}
		fmt.Printf("%-15s | %-10s | %s\n", col.ColumnName, col.DataType, nullable)
	}
	fmt.Println("─────────────────────────────────────────────────────\n")

	fmt.Println("✅ Migration complete!")
	fmt.Println("\nNow try saving SOAP notes again from the frontend.")
}
