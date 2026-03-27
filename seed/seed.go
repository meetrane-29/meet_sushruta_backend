package seed

import (
	"fmt"
	"log"
	"time"

	"meet_sushruta/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ComprehensiveSeed performs all database seeding with proper error handling and idempotency
func ComprehensiveSeed(db *gorm.DB) error {
	var (
		userCount         int64
		doctorCount       int64
		patientCount      int64
		appointmentCount  int64
		prescriptionCount int64
		medicineCount     int64
		labCount          int64
		billCount         int64
	)

	// Seed users (14 total: 1 admin, 3 doctors, 2 nurses, 1 pharmacist, 1 lab tech, 6 patients)
	if err := SeedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}
	db.Model(&model.User{}).Count(&userCount)

	// Seed doctors (3 with specializations linked to doctor users)
	if err := SeedDoctors(db); err != nil {
		return fmt.Errorf("failed to seed doctors: %w", err)
	}
	db.Model(&model.Doctor{}).Count(&doctorCount)

	// Seed patients (6 with medical history linked to patient users)
	if err := SeedPatients(db); err != nil {
		return fmt.Errorf("failed to seed patients: %w", err)
	}
	db.Model(&model.Patient{}).Count(&patientCount)

	// Seed appointments (9 total spread across different dates/statuses)
	if err := SeedAppointments(db); err != nil {
		return fmt.Errorf("failed to seed appointments: %w", err)
	}

	// Seed vitals for in_progress and completed appointments
	if err := SeedVitals(db); err != nil {
		return fmt.Errorf("failed to seed vitals: %w", err)
	}

	// Seed nurses
	if err := SeedNurses(db); err != nil {
		return fmt.Errorf("failed to seed nurses: %w", err)
	}

	// Seed medicines
	if err := SeedMedicines(db); err != nil {
		return fmt.Errorf("failed to seed medicines: %w", err)
	}
	db.Model(&model.Medicine{}).Count(&medicineCount)

	// Seed prescriptions
	if err := SeedPrescriptions(db); err != nil {
		return fmt.Errorf("failed to seed prescriptions: %w", err)
	}
	db.Model(&model.Prescription{}).Count(&prescriptionCount)

	// Seed lab orders (6 total with varied statuses)
	if err := SeedLabOrders(db); err != nil {
		return fmt.Errorf("failed to seed lab orders: %w", err)
	}
	db.Model(&model.LabRequest{}).Count(&labCount)

	// Seed bills
	if err := SeedBills(db); err != nil {
		return fmt.Errorf("failed to seed bills: %w", err)
	}
	db.Model(&model.Bill{}).Count(&billCount)

	// Count appointments after seeding
	db.Model(&model.Appointment{}).Count(&appointmentCount)

	// Print confirmation summary
	log.Printf("\n✅ Database seeding completed successfully!")
	log.Printf("✅ Summary: %d users, %d doctors, %d patients, %d appointments, %d prescriptions, %d medicines, %d lab requests, %d bills\n",
		userCount, doctorCount, patientCount, appointmentCount, prescriptionCount, medicineCount, labCount, billCount)

	return nil
}

// SeedUsers creates all user roles: admin, doctors, nurses, pharmacist, lab tech, and patients
// Uses FirstOrCreate for idempotency - safe to run multiple times
func SeedUsers(db *gorm.DB) error {
	users := []struct {
		email     string
		phone     string
		firstName string
		lastName  string
		role      string
		password  string
	}{
		// Admin (1)
		{"admin@test.com", "9900000001", "Admin", "User", "admin", "password123"},
		// Doctors (3)
		{"doctor@test.com", "9900000002", "Priya", "Sharma", "doctor", "password123"},
		{"doctor2@test.com", "9900000003", "Arjun", "Mehta", "doctor", "password123"},
		{"doctor3@test.com", "9900000004", "Sunita", "Patel", "doctor", "password123"},
		// Nurses (2)
		{"nurse@test.com", "9900000005", "Anjali", "Singh", "nurse", "password123"},
		{"nurse2@test.com", "9900000006", "Deepa", "Kumar", "nurse", "password123"},
		// Pharmacist (1)
		{"pharmacy@test.com", "9900000007", "Rajesh", "Desai", "pharmacy", "password123"},
		// Lab Tech (1)
		{"lab@test.com", "9900000008", "Vikram", "Verma", "lab", "password123"},
		// Patients (6)
		{"patient1@test.com", "9900000009", "Rohan", "Kapoor", "patient", "password123"},
		{"patient2@test.com", "9900000010", "Priyanka", "Gupta", "patient", "password123"},
		{"patient3@test.com", "9900000011", "Aditya", "Mishra", "patient", "password123"},
		{"patient4@test.com", "9900000012", "Neha", "Nair", "patient", "password123"},
		{"patient5@test.com", "9900000013", "Sanjay", "Reddy", "patient", "password123"},
		{"patient6@test.com", "9900000014", "Kavya", "Das", "patient", "password123"},
	}

	for _, u := range users {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", u.email, err)
			continue
		}

		// Use FirstOrCreate with email as unique identifier for idempotency
		user := model.User{
			Email:     u.email,
			Phone:     u.phone,
			FirstName: u.firstName,
			LastName:  u.lastName,
			Role:      u.role,
			Password:  string(hashedPassword),
			Active:    true,
		}

		if err := db.FirstOrCreate(&user, model.User{Email: u.email}).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", u.email, err)
			continue
		}

		log.Printf("✓ User ensured: %s (%s)", u.email, u.role)
	}

	// Print user count
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	log.Printf("📊 Total users in database: %d\n", userCount)

	return nil
}

// SeedDoctors creates doctor records linked to doctor users
// Uses FirstOrCreate with license number as unique identifier for idempotency
func SeedDoctors(db *gorm.DB) error {
	doctors := []struct {
		email           string
		specialization  string
		licenseNumber   string
		consultationFee float64
	}{
		{"doctor@test.com", "Cardiology", "MC-1001", 500},
		{"doctor2@test.com", "Orthopedics", "MC-1002", 400},
		{"doctor3@test.com", "Neurology", "MC-1003", 600},
	}

	for _, d := range doctors {
		// Find user by email
		var user model.User
		if err := db.Where("email = ?", d.email).First(&user).Error; err != nil {
			log.Printf("User %s not found for doctor: %v", d.email, err)
			continue
		}

		// Use FirstOrCreate with license number as unique identifier for idempotency
		doctor := model.Doctor{
			UserID:          user.ID,
			Specialization:  d.specialization,
			LicenseNumber:   d.licenseNumber,
			ConsultationFee: d.consultationFee,
			Department:      d.specialization,
			Bio:             fmt.Sprintf("Experienced %s specialist with expertise in patient care and diagnosis.", d.specialization),
		}

		if err := db.FirstOrCreate(&doctor, model.Doctor{LicenseNumber: d.licenseNumber}).Error; err != nil {
			log.Printf("Failed to seed doctor %s: %v", d.email, err)
			continue
		}

		log.Printf("✓ Doctor ensured: Dr. %s %s - %s", user.FirstName, user.LastName, d.specialization)
	}

	// Print doctor count
	var doctorCount int64
	db.Model(&model.Doctor{}).Count(&doctorCount)
	log.Printf("📊 Total doctors in database: %d\n", doctorCount)

	return nil
}

// SeedPatients creates patient records with realistic Indian data (6 patients)
// Uses FirstOrCreate with email lookup for idempotency
func SeedPatients(db *gorm.DB) error {
	patients := []struct {
		email            string
		dateOfBirth      string
		gender           string
		bloodGroup       string
		address          string
		emergencyContact string
		medicalHistory   string
		allergies        string
	}{
		{
			"patient1@test.com",
			"1985-03-15",
			"Male",
			"A+",
			"123 MG Road, Bangalore",
			"9999999999",
			"Hypertension since 2015, controlled with medication. Seasonal allergies.",
			"Penicillin (rash)",
		},
		{
			"patient2@test.com",
			"1992-07-22",
			"Female",
			"B+",
			"456 Bandra Street, Mumbai",
			"9999999998",
			"Diabetes Type 2 diagnosed in 2018. Family history of heart disease.",
			"None",
		},
		{
			"patient3@test.com",
			"1988-11-08",
			"Male",
			"O+",
			"789 Connaught Place, Delhi",
			"9999999997",
			"Previous fracture of left arm in 2019. Recovering well. Regular physiotherapy.",
			"Aspirin (GI upset)",
		},
		{
			"patient4@test.com",
			"1995-05-30",
			"Female",
			"AB+",
			"321 Cathedral Road, Chennai",
			"9999999996",
			"Thyroid disorder (hypothyroidism) - on levothyroxine. No other chronic conditions.",
			"None",
		},
		{
			"patient5@test.com",
			"1980-09-12",
			"Male",
			"O-",
			"654 Park Street, Kolkata",
			"9999999995",
			"Asthma since childhood, well controlled. Occasional seasonal exacerbations.",
			"Sulfonamides (rash), Shellfish (anaphylaxis risk)",
		},
		{
			"patient6@test.com",
			"1990-01-28",
			"Female",
			"AB-",
			"987 Salt Lake, Kolkata",
			"9999999994",
			"No significant medical history. Regular wellness checkups. Minor seasonal cough.",
			"Lactose intolerance",
		},
	}

	for _, p := range patients {
		// Find user by email
		var user model.User
		if err := db.Where("email = ?", p.email).First(&user).Error; err != nil {
			log.Printf("User %s not found for patient: %v", p.email, err)
			continue
		}

		// Use FirstOrCreate with UserID as unique identifier for idempotency
		patient := model.Patient{
			UserID:           user.ID,
			DateOfBirth:      p.dateOfBirth,
			Gender:           p.gender,
			BloodGroup:       p.bloodGroup,
			Address:          p.address,
			EmergencyContact: p.emergencyContact,
			MedicalHistory:   p.medicalHistory,
			Allergies:        p.allergies,
		}

		if err := db.FirstOrCreate(&patient, model.Patient{UserID: user.ID}).Error; err != nil {
			log.Printf("Failed to seed patient %s: %v", p.email, err)
			continue
		}

		log.Printf("✓ Patient ensured: %s %s (%s, %s)", user.FirstName, user.LastName, p.bloodGroup, p.dateOfBirth)
	}

	// Print patient count
	var patientCount int64
	db.Model(&model.Patient{}).Count(&patientCount)
	log.Printf("📊 Total patients in database: %d\n", patientCount)

	return nil
}

// SeedAppointments creates 9 appointments spread across multiple days with different statuses
// Uses FirstOrCreate with composite keys for idempotency
func SeedAppointments(db *gorm.DB) error {
	// Fetch all doctors and patients
	var doctors []model.Doctor
	var patients []model.Patient

	if err := db.Find(&doctors).Error; err != nil {
		return fmt.Errorf("failed to fetch doctors: %w", err)
	}
	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to fetch patients: %w", err)
	}

	if len(doctors) == 0 || len(patients) == 0 {
		return fmt.Errorf("no doctors or patients found to create appointments")
	}

	now := time.Now()

	// Appointments slice: each entry specifies days offset and status
	appointments := []struct {
		patientIdx int
		doctorIdx  int
		daysAgo    int
		status     string
		time       string
		reason     string
	}{
		// TODAY (3 appointments)
		{0, 0, 0, "pending", "09:00", "General checkup"},
		{1, 1, 0, "confirmed", "10:30", "Follow-up consultation"},
		{2, 2, 0, "in_progress", "14:00", "Neurological assessment"},
		// YESTERDAY (3 appointments)
		{3, 0, 1, "completed", "11:00", "Cardiac evaluation"},
		{4, 1, 1, "completed", "15:00", "Orthopedic consultation"},
		{5, 2, 1, "completed", "09:30", "Allergy screening"},
		// 3 DAYS AGO (3 appointments)
		{0, 1, 3, "completed", "10:00", "Annual health checkup"},
		{1, 2, 3, "completed", "13:00", "Neurological follow-up"},
		{2, 0, 3, "completed", "16:00", "Gastrointestinal review"},
		// 7 DAYS AGO (3 appointments)
		{3, 2, 7, "completed", "11:30", "Hypertension management"},
		{4, 0, 7, "completed", "14:30", "Orthopedic review"},
		{5, 1, 7, "completed", "10:15", "Diabetes management"},
	}

	for _, a := range appointments {
		appointmentDate := now.AddDate(0, 0, -a.daysAgo).Format("2006-01-02")
		doctorID := doctors[a.doctorIdx%len(doctors)].ID
		patientID := patients[a.patientIdx%len(patients)].ID

		// Use FirstOrCreate with composite key for idempotency
		appointment := model.Appointment{
			PatientID:       patientID,
			DoctorID:        doctorID,
			AppointmentDate: appointmentDate,
			AppointmentTime: a.time,
			Status:          model.AppointmentStatus(a.status),
			Reason:          a.reason,
			Notes:           fmt.Sprintf("Patient status: %s. Appointment recorded.", a.status),
		}

		if err := db.FirstOrCreate(&appointment, model.Appointment{
			PatientID:       patientID,
			DoctorID:        doctorID,
			AppointmentDate: appointmentDate,
			AppointmentTime: a.time,
		}).Error; err != nil {
			log.Printf("Failed to seed appointment: %v", err)
			continue
		}

		log.Printf("✓ Appointment ensured: Patient %d -> Doctor %d on %s at %s [%s]",
			a.patientIdx+1, a.doctorIdx+1, appointmentDate, a.time, a.status)
	}

	// Print appointment count
	var appointmentCount int64
	db.Model(&model.Appointment{}).Count(&appointmentCount)
	log.Printf("📊 Total appointments in database: %d\n", appointmentCount)

	return nil
}

// SeedVitals creates vital signs records for in_progress and completed appointments
// Uses FirstOrCreate with composite keys for idempotency
func SeedVitals(db *gorm.DB) error {
	// Fetch all in_progress and completed appointments
	var appointments []model.Appointment
	if err := db.Where("status IN ?", []string{"in_progress", "completed"}).
		Find(&appointments).Error; err != nil {
		return fmt.Errorf("failed to fetch appointments: %w", err)
	}

	if len(appointments) == 0 {
		log.Println("No in_progress or completed appointments found for vitals seeding")
		return nil
	}

	// Vitals data pool for variety
	vitalsData := []struct {
		temperature     float64
		bloodPressure   string
		heartRate       int
		respiratoryRate int
		weight          float64
		height          float64
		bloodSugar      float64
		oxygen          int
	}{
		{36.5, "120/80", 72, 16, 70.0, 175.0, 100.0, 98},
		{36.8, "130/85", 78, 18, 65.0, 160.0, 110.0, 97},
		{37.2, "140/90", 82, 20, 75.0, 180.0, 120.0, 96},
		{36.7, "125/82", 75, 17, 68.0, 170.0, 105.0, 98},
		{36.9, "118/78", 70, 16, 72.0, 172.0, 95.0, 99},
	}

	for idx, appointment := range appointments {
		// Assign vitals round-robin from the data pool
		vitalData := vitalsData[idx%len(vitalsData)]

		// Create recorded_at timestamp based on appointment date and time
		recordedAtStr := appointment.AppointmentDate + " " + appointment.AppointmentTime

		// Use FirstOrCreate with composite key for idempotency
		vital := model.Vitals{
			PatientID:       appointment.PatientID,
			Temperature:     vitalData.temperature,
			BloodPressure:   vitalData.bloodPressure,
			HeartRate:       vitalData.heartRate,
			RespiratoryRate: vitalData.respiratoryRate,
			Weight:          vitalData.weight,
			Height:          vitalData.height,
			BloodSugar:      vitalData.bloodSugar,
			Oxygen:          vitalData.oxygen,
			RecordedAt:      recordedAtStr,
		}

		if err := db.FirstOrCreate(&vital, model.Vitals{
			PatientID:  appointment.PatientID,
			RecordedAt: recordedAtStr,
		}).Error; err != nil {
			log.Printf("Failed to seed vital signs: %v", err)
			continue
		}

		log.Printf("✓ Vitals ensured: Patient %s - Temp: %.1f°C, BP: %s, HR: %d bpm, O2: %d%%",
			appointment.PatientID.String()[:8], vitalData.temperature, vitalData.bloodPressure, vitalData.heartRate, vitalData.oxygen)
	}

	// Print vitals count
	var vitalsCount int64
	db.Model(&model.Vitals{}).Count(&vitalsCount)
	log.Printf("📊 Total vitals records in database: %d\n", vitalsCount)

	return nil
}

// SeedNurses creates nurse records linked to nurse users
func SeedNurses(db *gorm.DB) error {
	nurses := []struct {
		email         string
		licenseNumber string
		department    string
		shift         string
	}{
		{"nurse@test.com", "NL-5001", "General Ward", "morning"},
		{"nurse2@test.com", "NL-5002", "ICU", "evening"},
	}

	for _, n := range nurses {
		// Check if nurse already exists with this email
		var existingNurse model.Nurse
		var user model.User

		// Find user by email
		if err := db.Where("email = ?", n.email).First(&user).Error; err != nil {
			log.Printf("User %s not found for nurse: %v", n.email, err)
			continue
		}

		// Check if nurse already exists
		if err := db.Where("user_id = ?", user.ID).First(&existingNurse).Error; err == nil {
			continue // Nurse already exists
		}

		// Create nurse
		nurse := model.Nurse{
			ID:            uuid.New(),
			UserID:        user.ID,
			LicenseNumber: n.licenseNumber,
			Department:    n.department,
			Shift:         n.shift,
		}

		if err := db.Create(&nurse).Error; err != nil {
			log.Printf("Failed to create nurse %s: %v", n.email, err)
			continue
		}

		log.Printf("✓ Created nurse: %s %s - %s (%s shift)", user.FirstName, user.LastName, n.department, n.shift)
	}

	return nil
}

// SeedMedicines creates 10 medicines with realistic Indian pharmacy data
// Uses FirstOrCreate with name as unique identifier for idempotency
func SeedMedicines(db *gorm.DB) error {
	medicines := []struct {
		name          string
		genericName   string
		dosage        string
		stockQuantity int
		reorderLevel  int
		price         float64
		description   string
	}{
		{"Paracetamol", "Acetaminophen", "500mg", 200, 50, 2, "Pain reliever and fever reducer"},
		{"Amoxicillin", "Amoxicillin", "250mg", 8, 50, 15, "Antibiotic for bacterial infections (LOW STOCK)"},
		{"Omeprazole", "Omeprazole", "20mg", 150, 40, 5, "Proton pump inhibitor for acid reflux"},
		{"Metformin", "Metformin", "500mg", 300, 100, 3, "Antidiabetic medication for Type 2 diabetes"},
		{"Atorvastatin", "Atorvastatin", "10mg", 50, 30, 8, "Statin for cholesterol management"},
		{"Azithromycin", "Azithromycin", "500mg", 5, 30, 25, "Macrolide antibiotic (LOW STOCK)"},
		{"Cetirizine", "Cetirizine", "10mg", 180, 60, 2, "Antihistamine for allergies"},
		{"Pantoprazole", "Pantoprazole", "40mg", 120, 40, 6, "Proton pump inhibitor for gastric ulcers"},
		{"Amlodipine", "Amlodipine", "5mg", 90, 30, 4, "Calcium channel blocker for hypertension"},
		{"Ibuprofen", "Ibuprofen", "400mg", 160, 50, 3, "NSAID for pain and inflammation"},
	}

	for _, m := range medicines {
		// Use FirstOrCreate with name as unique identifier for idempotency
		medicine := model.Medicine{
			Name:          m.name,
			GenericName:   m.genericName,
			Dosage:        m.dosage,
			StockQuantity: m.stockQuantity,
			ReorderLevel:  m.reorderLevel,
			Price:         m.price,
			Description:   m.description,
			Manufacturer:  "Indian Pharma Pvt Ltd",
			ExpiryDate:    "2026-12-31",
			Active:        true,
		}

		if err := db.FirstOrCreate(&medicine, model.Medicine{Name: m.name}).Error; err != nil {
			log.Printf("Failed to seed medicine %s: %v", m.name, err)
			continue
		}

		log.Printf("✓ Medicine ensured: %s %s - Stock: %d, Price: ₹%.2f", m.name, m.dosage, m.stockQuantity, m.price)
	}

	// Print medicine count
	var medicineCount int64
	db.Model(&model.Medicine{}).Count(&medicineCount)
	log.Printf("📊 Total medicines in database: %d\n", medicineCount)

	return nil
}

// SeedPrescriptions creates 6 prescriptions linked to completed appointments
// Each prescription has 2-3 items with dosage, frequency and duration
// Uses FirstOrCreate with composite keys for idempotency
func SeedPrescriptions(db *gorm.DB) error {
	// Get completed appointments
	var completedAppointments []model.Appointment
	if err := db.Where("status = ?", "completed").Find(&completedAppointments).Error; err != nil {
		return fmt.Errorf("failed to fetch completed appointments: %w", err)
	}

	if len(completedAppointments) < 6 {
		log.Printf("Warning: Only %d completed appointments found, need 6 for prescriptions", len(completedAppointments))
	}

	// Get medicines
	var medicines []model.Medicine
	if err := db.Find(&medicines).Error; err != nil {
		return fmt.Errorf("failed to fetch medicines: %w", err)
	}

	if len(medicines) < 10 {
		return fmt.Errorf("not enough medicines for prescription items (need 10, have %d)", len(medicines))
	}

	// Define 6 prescriptions with specific medicines and dosages
	prescriptions := []struct {
		appointmentIdx int
		items          []struct {
			medicineIdx  int
			dosage       string
			frequency    string
			duration     string
			instructions string
		}
		status string
	}{
		{
			appointmentIdx: 0,
			items: []struct {
				medicineIdx  int
				dosage       string
				frequency    string
				duration     string
				instructions string
			}{
				{0, "1 tablet", "3 times daily", "5 days", "Take after meals"},
				{1, "1 capsule", "2 times daily", "7 days", "With full glass of water"},
			},
			status: "active",
		},
		{
			appointmentIdx: 1,
			items: []struct {
				medicineIdx  int
				dosage       string
				frequency    string
				duration     string
				instructions string
			}{
				{2, "1 tablet", "1 time daily", "14 days", "Take before breakfast"},
				{3, "1 tablet", "2 times daily", "30 days", "With food"},
			},
			status: "active",
		},
		{
			appointmentIdx: 2,
			items: []struct {
				medicineIdx  int
				dosage       string
				frequency    string
				duration     string
				instructions string
			}{
				{4, "1 tablet", "1 time daily", "30 days", "Take in evening"},
				{5, "1 capsule", "2 times daily", "10 days", "Complete full course"},
				{6, "1 tablet", "As needed", "PRN", "Maximum 3 tablets per day"},
			},
			status: "dispensed",
		},
		{
			appointmentIdx: 3,
			items: []struct {
				medicineIdx  int
				dosage       string
				frequency    string
				duration     string
				instructions string
			}{
				{7, "1 tablet", "2 times daily", "20 days", "Before meals"},
				{8, "1 tablet", "1 time daily", "30 days", "In morning"},
			},
			status: "dispensed",
		},
		{
			appointmentIdx: 4,
			items: []struct {
				medicineIdx  int
				dosage       string
				frequency    string
				duration     string
				instructions string
			}{
				{9, "1 tablet", "2 times daily", "7 days", "With food if needed"},
				{0, "2 tablets", "1 time daily", "10 days", "At bedtime"},
				{3, "1 tablet", "1 time daily", "30 days", "With breakfast"},
			},
			status: "active",
		},
		{
			appointmentIdx: 5,
			items: []struct {
				medicineIdx  int
				dosage       string
				frequency    string
				duration     string
				instructions string
			}{
				{1, "1 capsule", "1 time daily", "7 days", "Complete course"},
				{6, "1 tablet", "3 times daily", "5 days", "After meals"},
			},
			status: "active",
		},
	}

	prescriptionCount := 0
	for pIdx, p := range prescriptions {
		if pIdx >= len(completedAppointments) {
			log.Printf("Warning: Not enough appointments for prescription %d", pIdx+1)
			continue
		}

		app := completedAppointments[p.appointmentIdx]
		prescriptionDate := app.AppointmentDate

		// Use FirstOrCreate with appointmentID as composite key
		prescription := model.Prescription{
			PatientID:        app.PatientID,
			DoctorID:         app.DoctorID,
			AppointmentID:    &app.ID,
			PrescriptionDate: prescriptionDate,
			IssuedAt:         prescriptionDate + " 10:00",
			Status:           p.status,
			Notes:            fmt.Sprintf("Prescription for appointment on %s. Follow-up in 2 weeks.", prescriptionDate),
		}

		if err := db.FirstOrCreate(&prescription, model.Prescription{AppointmentID: &app.ID}).Error; err != nil {
			log.Printf("Failed to seed prescription: %v", err)
			continue
		}

		// Create prescription items
		for _, item := range p.items {
			if item.medicineIdx >= len(medicines) {
				log.Printf("Warning: Medicine index %d out of range", item.medicineIdx)
				continue
			}

			prescriptionItem := model.PrescriptionItem{
				PrescriptionID: prescription.ID,
				MedicineID:     medicines[item.medicineIdx].ID,
				Dosage:         item.dosage,
				Frequency:      item.frequency,
				Duration:       item.duration,
				TimeOfDay:      "morning/evening",
				Quantity:       int(item.frequency[0]-'0') * 10, // Rough calculation based on frequency
				Instructions:   item.instructions,
			}

			if err := db.FirstOrCreate(&prescriptionItem, model.PrescriptionItem{
				PrescriptionID: prescription.ID,
				MedicineID:     medicines[item.medicineIdx].ID,
			}).Error; err != nil {
				log.Printf("Failed to seed prescription item: %v", err)
				continue
			}
		}

		prescriptionCount++
		log.Printf("✓ Prescription ensured %d: %d items [%s]", pIdx+1, len(p.items), p.status)
	}

	// Print prescription count
	var finalPrescriptionCount int64
	db.Model(&model.Prescription{}).Count(&finalPrescriptionCount)
	log.Printf("📊 Total prescriptions in database: %d\n", finalPrescriptionCount)

	return nil
}

// SeedLabRequests creates 5 lab orders with varied statuses
// SeedLabOrders creates 6 lab orders with varied statuses
// Uses FirstOrCreate with composite keys for idempotency
func SeedLabOrders(db *gorm.DB) error {
	// Get patients and doctors
	var patients []model.Patient
	var doctors []model.Doctor

	if err := db.Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to fetch patients: %w", err)
	}
	if err := db.Find(&doctors).Error; err != nil {
		return fmt.Errorf("failed to fetch doctors: %w", err)
	}

	if len(patients) == 0 || len(doctors) == 0 {
		return fmt.Errorf("no patients or doctors found for lab orders")
	}

	now := time.Now()
	labOrders := []struct {
		patientIdx int
		doctorIdx  int
		testType   string
		testName   string
		status     string
		daysAgo    int
		resultURL  string
	}{
		// 2 completed with result URLs
		{0, 0, "BloodWork", "CBC", "completed", 5, "https://reports.hospital.com/lab/cbc_001_20260320.pdf"},
		{1, 1, "BloodWork", "LFT", "completed", 3, "https://reports.hospital.com/lab/lft_001_20260322.pdf"},
		// 1 in_progress
		{2, 2, "Metabolic", "Blood Sugar", "in_progress", 1, ""},
		// 3 pending
		{3, 0, "BloodWork", "Lipid Profile", "pending", 0, ""},
		{4, 1, "Urinalysis", "Urine Routine", "pending", 0, ""},
		{5, 2, "Metabolic", "TFT", "pending", 0, ""},
	}

	for _, lr := range labOrders {
		requestDate := now.AddDate(0, 0, -lr.daysAgo).Format("2006-01-02")
		patientID := patients[lr.patientIdx%len(patients)].ID
		doctorID := doctors[lr.doctorIdx%len(doctors)].ID

		// Use FirstOrCreate with composite key for idempotency
		labOrder := model.LabRequest{
			PatientID:   patientID,
			DoctorID:    doctorID,
			TestType:    lr.testType,
			TestName:    lr.testName,
			RequestedAt: requestDate + " 09:00",
			Status:      lr.status,
			Priority:    "normal",
			ResultURL:   lr.resultURL,
		}

		if err := db.FirstOrCreate(&labOrder, model.LabRequest{
			PatientID: patientID,
			DoctorID:  doctorID,
			TestName:  lr.testName,
		}).Error; err != nil {
			log.Printf("Failed to seed lab order: %v", err)
			continue
		}

		// Add completion details if completed
		if lr.status == "completed" && lr.resultURL != "" {
			completedDate := requestDate
			db.Model(&labOrder).Updates(map[string]interface{}{
				"completed_at":   completedDate,
				"result_url":     lr.resultURL,
				"result_summary": fmt.Sprintf("Normal %s results. All parameters within reference range.", lr.testName),
			})
		}

		log.Printf("✓ Lab order ensured: %s (%s) - Patient %d [%s]",
			lr.testName, lr.testType, lr.patientIdx+1, lr.status)
	}

	// Print lab order count
	var labOrderCount int64
	db.Model(&model.LabRequest{}).Count(&labOrderCount)
	log.Printf("📊 Total lab orders in database: %d\n", labOrderCount)

	return nil
}

// SeedBills creates 4 bills with payment details and line items
// 2 paid bills (upi, cash) and 2 pending bills
// Uses FirstOrCreate with composite keys for idempotency
func SeedBills(db *gorm.DB) error {
	// Get appointments and patients
	var completedAppointments []model.Appointment
	if err := db.Where("status = ?", "completed").Find(&completedAppointments).Error; err != nil {
		return fmt.Errorf("failed to fetch completed appointments: %w", err)
	}

	if len(completedAppointments) < 4 {
		log.Printf("Warning: Only %d completed appointments found, need 4 for bills", len(completedAppointments))
	}

	now := time.Now()
	bills := []struct {
		appointmentIdx  int
		consultationFee float64
		medicines       float64
		labTests        float64
		paymentStatus   string
		paymentMethod   string
		daysAgo         int
	}{
		// 2 paid bills
		{0, 500, 150, 300, "paid", "upi", 3},
		{1, 400, 200, 250, "paid", "cash", 2},
		// 2 pending bills
		{2, 600, 100, 400, "pending", "card", 0},
		{3, 450, 175, 200, "pending", "online", 0},
	}

	for bIdx, b := range bills {
		if b.appointmentIdx >= len(completedAppointments) {
			log.Printf("Warning: Not enough appointments for bill %d", bIdx+1)
			continue
		}

		app := completedAppointments[b.appointmentIdx]

		// Use FirstOrCreate with appointmentID and bill number for idempotency
		billNumber := fmt.Sprintf("BILL-%06d", bIdx+1001)
		billDate := now.AddDate(0, 0, -b.daysAgo).Format("2006-01-02")

		taxPercentage := 5.0
		subtotal := b.consultationFee + b.medicines + b.labTests
		taxAmount := (subtotal * taxPercentage) / 100
		totalAmount := subtotal + taxAmount

		var paidAmount, dueAmount float64
		if b.paymentStatus == "paid" {
			paidAmount = totalAmount
			dueAmount = 0
		} else {
			paidAmount = 0
			dueAmount = totalAmount
		}

		bill := model.Bill{
			PatientID:       app.PatientID,
			AppointmentID:   &app.ID,
			BillNumber:      billNumber,
			BillDate:        billDate,
			ConsultationFee: b.consultationFee,
			LabTests:        b.labTests,
			Medicines:       b.medicines,
			OtherCharges:    0,
			Discount:        0,
			TaxPercentage:   taxPercentage,
			TaxAmount:       taxAmount,
			TotalAmount:     totalAmount,
			PaidAmount:      paidAmount,
			DueAmount:       dueAmount,
			PaymentStatus:   b.paymentStatus,
			PaymentMethod:   b.paymentMethod,
			Notes:           fmt.Sprintf("Bill %d: Consultation, lab tests and medicines. Tax included.", bIdx+1),
		}

		if b.paymentStatus == "paid" {
			bill.PaymentDate = &billDate
		}

		if err := db.FirstOrCreate(&bill, model.Bill{BillNumber: billNumber}).Error; err != nil {
			log.Printf("Failed to seed bill: %v", err)
			continue
		}

		// Create bill items
		billItems := []struct {
			itemType    string
			description string
			quantity    int
			unitPrice   float64
		}{
			{string(model.BillItemConsultation), "Consultation Fee", 1, b.consultationFee},
			{string(model.BillItemLab), "Lab Tests", 1, b.labTests},
			{string(model.BillItemMedicine), "Medicines", 1, b.medicines},
		}

		for _, bi := range billItems {
			if err := db.FirstOrCreate(&model.BillItem{}, model.BillItem{
				BillID:      bill.ID,
				Description: bi.description,
			}).Error; err != nil {
				log.Printf("Failed to seed bill item: %v", err)
				continue
			}

			// Fetch the created bill item to get its ID for logging
			var billItem model.BillItem
			db.Where("bill_id = ? AND description = ?", bill.ID, bi.description).First(&billItem)
			db.Model(&billItem).Updates(map[string]interface{}{
				"item_type":  model.BillItemType(bi.itemType),
				"quantity":   bi.quantity,
				"unit_price": bi.unitPrice,
			})
		}

		log.Printf("✓ Bill ensured: %s - Total ₹%.2f [%s via %s]", billNumber, totalAmount, b.paymentStatus, b.paymentMethod)
	}

	// Print bill count
	var billCount int64
	db.Model(&model.Bill{}).Count(&billCount)
	log.Printf("📊 Total bills in database: %d\n", billCount)

	return nil
}
