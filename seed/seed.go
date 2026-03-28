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

	// Seed users (21 total: 1 admin, 10 doctors, 5 nurses, 1 pharmacist, 1 lab tech, 6 patients)
	if err := SeedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}
	db.Model(&model.User{}).Count(&userCount)

	// Seed specializations
	if err := SeedSpecializations(db); err != nil {
		return fmt.Errorf("failed to seed specializations: %w", err)
	}

	// Seed hospitals
	if err := SeedHospitals(db); err != nil {
		return fmt.Errorf("failed to seed hospitals: %w", err)
	}

	// Seed beds (ward-wise allocation)
	if err := SeedBeds(db); err != nil {
		return fmt.Errorf("failed to seed beds: %w", err)
	}

	// Seed medical equipment
	if err := SeedMedicalEquipment(db); err != nil {
		return fmt.Errorf("failed to seed medical equipment: %w", err)
	}

	// Seed operation theatres
	if err := SeedOperationTheatres(db); err != nil {
		return fmt.Errorf("failed to seed operation theatres: %w", err)
	}

	// Seed operation schedules
	if err := SeedOperationSchedules(db); err != nil {
		return fmt.Errorf("failed to seed operation schedules: %w", err)
	}

	// Seed doctors (10 with specializations linked to doctor users)
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
		// Doctors (10)
		{"doctor@test.com", "9900000002", "Priya", "Sharma", "doctor", "password123"},
		{"doctor2@test.com", "9900000003", "Arjun", "Mehta", "doctor", "password123"},
		{"doctor3@test.com", "9900000004", "Sunita", "Patel", "doctor", "password123"},
		{"doctor4@test.com", "9900000018", "Rajesh", "Verma", "doctor", "password123"},
		{"doctor5@test.com", "9900000019", "Aisha", "Khan", "doctor", "password123"},
		{"doctor6@test.com", "9900000020", "Vikram", "Sharma", "doctor", "password123"},
		{"doctor7@test.com", "9900000021", "Anjali", "Menon", "doctor", "password123"},
		{"doctor8@test.com", "9900000022", "Nikhil", "Gupta", "doctor", "password123"},
		{"doctor9@test.com", "9900000023", "Neha", "Singh", "doctor", "password123"},
		{"doctor10@test.com", "9900000024", "Amit", "Reddy", "doctor", "password123"},
		// Nurses (5) with different roles
		{"nurse@test.com", "9900000005", "Anjali", "Singh", "nurse", "password123"},
		{"nurse2@test.com", "9900000006", "Deepa", "Kumar", "nurse", "password123"},
		{"nurse3@test.com", "9900000015", "Priya", "Desai", "nurse", "password123"},
		{"nurse4@test.com", "9900000016", "Neelima", "Rao", "nurse", "password123"},
		{"nurse5@test.com", "9900000017", "Sneha", "Iyer", "nurse", "password123"},
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
		{"doctor4@test.com", "Pediatrics", "MC-1004", 350},
		{"doctor5@test.com", "Dermatology", "MC-1005", 450},
		{"doctor6@test.com", "Ophthalmology", "MC-1006", 550},
		{"doctor7@test.com", "ENT", "MC-1007", 400},
		{"doctor8@test.com", "General Medicine", "MC-1008", 300},
		{"doctor9@test.com", "Psychiatry", "MC-1009", 500},
		{"doctor10@test.com", "Gynecology", "MC-1010", 550},
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
		// TODAY (18 appointments - showcasing all 10 doctors)
		{0, 0, 0, "pending", "09:00", "General checkup"},
		{1, 1, 0, "confirmed", "10:30", "Follow-up consultation"},
		{2, 2, 0, "in_progress", "14:00", "Neurological assessment"},
		{3, 3, 0, "pending", "11:00", "Pediatric consultation"},
		{4, 4, 0, "confirmed", "12:30", "Skin treatment"},
		{5, 5, 0, "pending", "15:30", "Eye examination"},
		{0, 6, 0, "confirmed", "13:00", "Ear, nose and throat check"},
		{1, 7, 0, "in_progress", "16:00", "General medical exam"},
		{2, 8, 0, "confirmed", "10:00", "Mental health consultation"},
		{3, 9, 0, "pending", "14:30", "Women's health checkup"},
		{4, 0, 0, "completed", "08:00", "Cardiac evaluation"},
		{5, 1, 0, "completed", "11:00", "Orthopedic assessment"},
		{0, 2, 0, "completed", "09:30", "Neurological review"},
		{1, 3, 0, "confirmed", "15:00", "Child vaccination"},
		{2, 4, 0, "pending", "13:00", "Skin allergy test"},
		{3, 5, 0, "confirmed", "10:00", "Cataract consultation"},
		{4, 6, 0, "pending", "11:30", "Hearing test"},
		{5, 7, 0, "confirmed", "12:00", "Blood pressure checkup"},
		// YESTERDAY (10 appointments)
		{0, 8, 1, "completed", "09:00", "Psychiatric evaluation"},
		{1, 9, 1, "completed", "10:30", "Pregnancy checkup"},
		{2, 0, 1, "completed", "14:00", "Heart disease consultation"},
		{3, 1, 1, "completed", "11:00", "Bone injury treatment"},
		{4, 2, 1, "completed", "15:30", "Brain disorder assessment"},
		{5, 3, 1, "completed", "09:30", "Child health checkup"},
		{0, 4, 1, "completed", "12:00", "Dermatology treatment"},
		{1, 5, 1, "completed", "13:00", "Vision correction"},
		{2, 6, 1, "completed", "10:15", "ENT surgery follow-up"},
		{3, 7, 1, "completed", "16:00", "Routine checkup"},
		// 3 DAYS AGO (6 appointments)
		{4, 8, 3, "completed", "11:00", "Psychology session"},
		{5, 9, 3, "completed", "14:00", "Gynecology consultation"},
		{0, 0, 3, "completed", "10:00", "Cardiology follow-up"},
		{1, 1, 3, "completed", "13:00", "Physical therapy"},
		{2, 2, 3, "completed", "16:00", "Neurology session"},
		{3, 3, 3, "completed", "09:00", "Pediatric follow-up"},
		// 7 DAYS AGO (6 appointments)
		{4, 4, 7, "completed", "11:30", "Dermatology review"},
		{5, 5, 7, "completed", "14:30", "Ophthalmology checkup"},
		{0, 6, 7, "completed", "10:15", "ENT consultation"},
		{1, 7, 7, "completed", "13:15", "General medicine review"},
		{2, 8, 7, "completed", "15:00", "Psychiatric follow-up"},
		{3, 9, 7, "completed", "10:45", "Women's health"},
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

// SeedNurses creates nurse records linked to nurse users with various roles
func SeedNurses(db *gorm.DB) error {
	nurses := []struct {
		email         string
		licenseNumber string
		role          string
		department    string
		shift         string
	}{
		{"nurse@test.com", "NL-5001", "Staff Nurse", "General Ward", "morning"},
		{"nurse2@test.com", "NL-5002", "ICU Nurse", "ICU", "evening"},
		{"nurse3@test.com", "NL-5003", "Charge Nurse", "General Ward", "night"},
		{"nurse4@test.com", "NL-5004", "Operation Theatre Nurse", "OT", "morning"},
		{"nurse5@test.com", "NL-5005", "Pediatric Nurse", "Pediatrics", "evening"},
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
			log.Printf("✓ Nurse already exists: %s %s - %s (%s, %s shift)", user.FirstName, user.LastName, n.role, n.department, n.shift)
			continue // Nurse already exists
		}

		// Create nurse
		nurse := model.Nurse{
			ID:            uuid.New(),
			UserID:        user.ID,
			LicenseNumber: n.licenseNumber,
			Role:          n.role,
			Department:    n.department,
			Shift:         n.shift,
		}

		if err := db.Create(&nurse).Error; err != nil {
			log.Printf("Failed to create nurse %s: %v", n.email, err)
			continue
		}

		log.Printf("✓ Created nurse: %s %s - %s (%s, %s shift)", user.FirstName, user.LastName, n.role, n.department, n.shift)
	}

	// Print nurses count
	var nurseCount int64
	db.Model(&model.Nurse{}).Count(&nurseCount)
	log.Printf("📊 Total nurses in database: %d\n", nurseCount)

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

// SeedSpecializations creates common medical specializations
func SeedSpecializations(db *gorm.DB) error {
	specializations := []struct {
		name     string
		category string
	}{
		{"Cardiology", "Medicine"},
		{"Neurology", "Medicine"},
		{"Orthopedics", "Surgery"},
		{"General Surgery", "Surgery"},
		{"Pediatrics", "Medicine"},
		{"Dermatology", "Medicine"},
		{"ENT", "Medicine"},
		{"Ophthalmology", "Medicine"},
		{"Psychiatry", "Medicine"},
		{"Oncology", "Medicine"},
		{"Gastroenterology", "Medicine"},
		{"Urology", "Surgery"},
		{"Nephrology", "Medicine"},
		{"Rheumatology", "Medicine"},
		{"Endocrinology", "Medicine"},
		{"Pulmonology", "Medicine"},
		{"Emergency Medicine", "Medicine"},
		{"Anesthesiology", "Medicine"},
		{"Radiology", "Diagnostics"},
		{"Pathology", "Diagnostics"},
	}

	for _, s := range specializations {
		spec := model.Specialization{
			Name:     s.name,
			Category: s.category,
			IsActive: true,
		}

		if err := db.FirstOrCreate(&spec, model.Specialization{Name: s.name}).Error; err != nil {
			log.Printf("Failed to seed specialization %s: %v", s.name, err)
			continue
		}
		log.Printf("✓ Specialization ensured: %s", s.name)
	}

	var specCount int64
	db.Model(&model.Specialization{}).Count(&specCount)
	log.Printf("📊 Total specializations in database: %d\n", specCount)

	return nil
}

// SeedHospitals creates sample hospitals
func SeedHospitals(db *gorm.DB) error {
	hospitals := []struct {
		name          string
		address       string
		city          string
		state         string
		phone         string
		email         string
		totalBeds     int
		availableBeds int
	}{
		{
			"City Central Hospital",
			"123 Medical Lane, Healthcare District",
			"New Delhi",
			"Delhi",
			"011-40123456",
			"info@cityhospital.com",
			500,
			150,
		},
		{
			"Apollo Medical Center",
			"456 Health Avenue, Tech Park",
			"Bangalore",
			"Karnataka",
			"080-40000123",
			"contact@apollo.com",
			400,
			120,
		},
		{
			"Fortis Healthcare",
			"789 Wellness Blvd, Medical Zone",
			"Mumbai",
			"Maharashtra",
			"022-67890123",
			"support@fortis.com",
			350,
			100,
		},
		{
			"Max Super Specialty Hospital",
			"321 Care Street, Health Campus",
			"Gurgaon",
			"Haryana",
			"124-40000456",
			"info@maxhospital.com",
			450,
			130,
		},
		{
			"Fortis No. 1 Hospital",
			"654 Cure Drive, Hospital Zone",
			"Hyderabad",
			"Telangana",
			"040-67890789",
			"contact@fortis-hyd.com",
			380,
			110,
		},
	}

	for _, h := range hospitals {
		hospital := model.Hospital{
			Name:          h.name,
			Address:       h.address,
			City:          h.city,
			State:         h.state,
			Phone:         h.phone,
			Email:         h.email,
			TotalBeds:     h.totalBeds,
			AvailableBeds: h.availableBeds,
			IsActive:      true,
			IsVerified:    true,
		}

		if err := db.FirstOrCreate(&hospital, model.Hospital{Name: h.name}).Error; err != nil {
			log.Printf("Failed to seed hospital %s: %v", h.name, err)
			continue
		}
		log.Printf("✓ Hospital ensured: %s (%s)", h.name, h.city)
	}

	var hospitalCount int64
	db.Model(&model.Hospital{}).Count(&hospitalCount)
	log.Printf("📊 Total hospitals in database: %d\n", hospitalCount)

	return nil
}

// SeedBeds creates ward-wise bed allocation
func SeedBeds(db *gorm.DB) error {
	bedData := []struct {
		bedNumber string
		ward      string
		room      string
		floor     int
		bedType   string
		features  string
		dailyRate float64
	}{
		// ICU (10 beds)
		{"ICU-001", "ICU", "ICU-101", 1, "ICU", "Cardiac monitor, Ventilator, Advanced lighting", 5000},
		{"ICU-002", "ICU", "ICU-101", 1, "ICU", "Cardiac monitor, Ventilator", 5000},
		{"ICU-003", "ICU", "ICU-102", 1, "ICU", "Ventilator, Laminar flow", 5000},
		{"ICU-004", "ICU", "ICU-102", 1, "ICU", "Cardiac monitor, Oxygen concentrator", 5000},
		{"ICU-005", "ICU", "ICU-103", 1, "ICU", "Ventilator, Infusion pump", 5000},
		{"ICU-006", "ICU", "ICU-103", 1, "ICU", "Monitors available", 5000},
		{"ICU-007", "ICU", "ICU-104", 2, "ICU", "Full ICU setup", 5000},
		{"ICU-008", "ICU", "ICU-104", 2, "ICU", "Ventilator, Monitor", 5000},
		{"ICU-009", "ICU", "ICU-105", 2, "ICU", "Cardiac monitoring", 5000},
		{"ICU-010", "ICU", "ICU-105", 2, "ICU", "Standard ICU", 5000},

		// General Ward (15 beds)
		{"GW-001", "General", "Ward-201", 2, "General", "Basic amenities", 1500},
		{"GW-002", "General", "Ward-201", 2, "General", "Basic amenities", 1500},
		{"GW-003", "General", "Ward-202", 2, "General", "TV, AC", 1500},
		{"GW-004", "General", "Ward-202", 2, "General", "Basic amenities", 1500},
		{"GW-005", "General", "Ward-203", 3, "General", "Basic amenities", 1500},
		{"GW-006", "General", "Ward-203", 3, "General", "TV, AC", 1500},
		{"GW-007", "General", "Ward-204", 3, "General", "Basic amenities", 1500},
		{"GW-008", "General", "Ward-204", 3, "General", "Basic amenities", 1500},
		{"GW-009", "General", "Ward-205", 3, "General", "TV, AC", 1500},
		{"GW-010", "General", "Ward-205", 3, "General", "Basic amenities", 1500},
		{"GW-011", "General", "Ward-206", 4, "General", "Basic amenities", 1500},
		{"GW-012", "General", "Ward-206", 4, "General", "TV", 1500},
		{"GW-013", "General", "Ward-207", 4, "General", "Basic amenities", 1500},
		{"GW-014", "General", "Ward-207", 4, "General", "AC", 1500},
		{"GW-015", "General", "Ward-208", 4, "General", "Basic amenities", 1500},

		// Private Rooms (8 beds)
		{"PR-001", "Private", "Room-301", 3, "Private", "Premium: AC, TV, Attached bath, Sofa", 3500},
		{"PR-002", "Private", "Room-302", 3, "Private", "Premium: AC, TV, Attached bath", 3500},
		{"PR-003", "Private", "Room-303", 3, "Private", "Premium: AC, TV, Sofa bed", 3500},
		{"PR-004", "Private", "Room-304", 4, "Private", "De-luxe: AC, TV, Attached bath, Fridge", 3500},
		{"PR-005", "Private", "Room-305", 4, "Private", "Premium: AC, TV", 3500},
		{"PR-006", "Private", "Room-306", 4, "Private", "Premium: AC, TV, Sofa", 3500},
		{"PR-007", "Private", "Room-307", 4, "Private", "Standard: AC", 2500},
		{"PR-008", "Private", "Room-308", 5, "Private", "Deluxe: Full amenities", 4000},

		// Semi-Private Rooms (6 beds)
		{"SP-001", "Semi-Private", "Room-401", 4, "Semi-Private", "AC, TV, Shared bathroom", 2500},
		{"SP-002", "Semi-Private", "Room-401", 4, "Semi-Private", "AC, TV", 2500},
		{"SP-003", "Semi-Private", "Room-402", 4, "Semi-Private", "AC, Basic amenities", 2500},
		{"SP-004", "Semi-Private", "Room-402", 4, "Semi-Private", "AC, TV", 2500},
		{"SP-005", "Semi-Private", "Room-403", 5, "Semi-Private", "AC, TV", 2500},
		{"SP-006", "Semi-Private", "Room-403", 5, "Semi-Private", "AC, Basic amenities", 2500},
	}

	for idx, bd := range bedData {
		var existingBed model.Bed
		err := db.Where("bed_number = ?", bd.bedNumber).First(&existingBed).Error

		var status string
		// Make some beds occupied
		if idx%4 == 0 {
			status = "occupied"
		} else if idx%7 == 0 {
			status = "maintenance"
		} else {
			status = "available"
		}

		bed := model.Bed{
			BedNumber: bd.bedNumber,
			Ward:      bd.ward,
			Room:      bd.room,
			Floor:     bd.floor,
			BedType:   bd.bedType,
			Status:    status,
			Features:  bd.features,
			DailyRate: bd.dailyRate,
		}

		if err != nil {
			// Bed doesn't exist, create it
			if err := db.Create(&bed).Error; err != nil {
				log.Printf("Failed to seed bed %s: %v", bd.bedNumber, err)
				continue
			}
		}
		log.Printf("✓ Bed ensured: %s (%s - %s) [%s]", bd.bedNumber, bd.ward, bd.bedType, status)
	}

	var bedCount int64
	db.Model(&model.Bed{}).Count(&bedCount)
	log.Printf("📊 Total beds in database: %d\n", bedCount)

	return nil
}

// SeedMedicalEquipment creates medical equipment inventory
func SeedMedicalEquipment(db *gorm.DB) error {
	// Drop the old unique constraint on serial_number if it exists
	if db.Migrator().HasIndex("medical_equipments", "idx_medical_equipments_serial_number") {
		db.Migrator().DropIndex("medical_equipments", "idx_medical_equipments_serial_number")
		log.Println("Dropped old unique constraint on serial_number")
	}

	equipment := []struct {
		equipmentName     string
		equipmentType     string
		model             string
		manufacturer      string
		location          string
		status            string
		criticalEquipment bool
		dailyRentalRate   float64
	}{
		// Critical ICU Equipment
		{"Ventilator-ICU-01", "Ventilator", "Siemens SERVO-i", "Siemens", "ICU Ward-101", "working", true, 2000},
		{"Ventilator-ICU-02", "Ventilator", "Hamilton-G5", "Hamilton", "ICU Ward-102", "working", true, 2000},
		{"Cardiac-Monitor-01", "Cardiac Monitor", "GE Case", "GE Healthcare", "ICU Ward-101", "working", true, 1500},
		{"Cardiac-Monitor-02", "Cardiac Monitor", "Philips Intellivue", "Philips", "ICU Ward-102", "working", true, 1500},
		{"Defibrillator-01", "Defibrillator", "Philips Heartstart", "Philips", "ICU Ward-103", "working", true, 1200},
		{"Infusion-Pump-01", "Infusion Pump", "Baxter Colleague", "Baxter", "ICU Ward-101", "working", true, 800},
		{"Infusion-Pump-02", "Infusion Pump", "B. Braun Infusomat", "B. Braun", "ICU Ward-102", "working", true, 800},
		{"Oxygen-Concentrator-01", "Oxygen Concentrator", "Invacare Perfecto2", "Invacare", "ICU Ward-104", "working", false, 500},

		// Operating Theatre Equipment
		{"Surgical-Light-OT-01", "Surgical Light", "Draeger Symbia", "Draeger", "OT-1", "working", true, 1000},
		{"Surgical-Light-OT-02", "Surgical Light", "Stryker TPS", "Stryker", "OT-2", "working", true, 1000},
		{"Electrosurgical-Unit-01", "Electrosurgical Unit", "Conmed Sabre", "Conmed", "OT-1", "working", true, 800},
		{"Anesthesia-Machine-01", "Anesthesia Machine", "Draeger Primus", "Draeger", "OT-1", "working", true, 1500},
		{"Anesthesia-Machine-02", "Anesthesia Machine", "GE Avance", "GE", "OT-2", "repair", true, 1500},
		{"Autoclave-01", "Autoclave Sterilizer", "Tuttnauer", "Tuttnauer", "Sterilization", "working", false, 0},

		// Diagnostic Equipment
		{"CT-Scanner-01", "CT Scanner", "Siemens SOMATOM", "Siemens", "Radiology", "working", true, 5000},
		{"Ultrasound-01", "Ultrasound", "GE Vivid E95", "GE Healthcare", "Radiology", "working", false, 800},
		{"Ultrasound-02", "Ultrasound", "Philips EPIQ", "Philips", "Cardiology", "working", false, 1000},
		{"X-Ray-Machine-01", "X-Ray Machine", "Siemens Mobilett", "Siemens", "Radiology", "working", false, 2000},
		{"ECG-Machine-01", "ECG Machine", "Philips PageWriter", "Philips", "Cardiology", "working", false, 500},
		{"ABG-Analyzer-01", "ABG Analyzer", "Siemens RapidLab", "Siemens", "Lab", "working", false, 0},

		// Laboratory Equipment
		{"Hematology-Analyzer-01", "Hematology Analyzer", "Sysmex XN1000", "Sysmex", "Lab", "working", false, 0},
		{"Biochemistry-Analyzer-01", "Biochemistry Analyzer", "Roche Cobas", "Roche", "Lab", "working", false, 0},
		{"Centrifuge-01", "Centrifuge", "Eppendorf 5810R", "Eppendorf", "Lab", "working", false, 0},

		// Patient Monitoring
		{"BP-Monitor-01", "BP Monitor", "Omron Digital", "Omron", "General Ward", "working", false, 0},
		{"Pulse-Oximeter-01", "Pulse Oximeter", "Nellcor N600", "Nellcor", "ICU", "working", false, 200},

		// Equipment in Maintenance/Repair
		{"Ventilator-Old", "Ventilator", "Helpman 2000", "Helpman", "Maintenance", "under_maintenance", true, 0},
		{"Cardiac-Monitor-Old", "Cardiac Monitor", "Spacelabs 91370", "Spacelabs", "Maintenance", "repair", true, 0},
	}

	for _, eq := range equipment {
		var existingEquip model.MedicalEquipment
		err := db.Where("equipment_name = ?", eq.equipmentName).First(&existingEquip).Error

		medicalEquip := model.MedicalEquipment{
			EquipmentName:     eq.equipmentName,
			EquipmentType:     eq.equipmentType,
			Model:             eq.model,
			Manufacturer:      eq.manufacturer,
			Location:          eq.location,
			Status:            eq.status,
			CriticalEquipment: eq.criticalEquipment,
			DailyRentalRate:   eq.dailyRentalRate,
		}

		if err != nil && err == gorm.ErrRecordNotFound {
			if err := db.Create(&medicalEquip).Error; err != nil {
				log.Printf("Failed to seed equipment %s: %v", eq.equipmentName, err)
				continue
			}
		} else if err != nil {
			log.Printf("Error querying equipment %s: %v", eq.equipmentName, err)
			continue
		}
		log.Printf("✓ Equipment ensured: %s (%s) [%s]", eq.equipmentName, eq.equipmentType, eq.status)
	}

	var equipmentCount int64
	db.Model(&model.MedicalEquipment{}).Count(&equipmentCount)
	log.Printf("📊 Total medical equipment in database: %d\n", equipmentCount)

	return nil
}

// SeedOperationTheatres creates operation theatres
func SeedOperationTheatres(db *gorm.DB) error {
	theatres := []struct {
		theatreName   string
		floor         int
		capacity      int
		features      string
		equipmentList string
	}{
		{"OT-1 (General Surgery)", 2, 1, "Laminar flow, Advanced lighting, Video recording, Telemedicine", "Surgical Light, Electrosurgical Unit, Anesthesia Machine"},
		{"OT-2 (Cardiac Surgery)", 2, 1, "Laminar flow, Advanced lighting, ECMO ready", "Surgical Light, Cardiac bypass machine, Electrosurgical Unit"},
		{"OT-3 (Orthopedic)", 3, 1, "Laminar flow, C-arm compatible, Bone imaging", "Surgical Light, C-arm machine, Orthopedic instruments"},
		{"OT-4 (Emergency)", 1, 1, "Rapid setup, Full equipped", "Surgical Light, Electrosurgical Unit, Anesthesia Machine"},
		{"OT-5 (Laparoscopic)", 3, 1, "Laparoscopic tower, HD monitors", "Surgical Light, Laparoscopic equipment, Video monitors"},
		{"Minor OT (Procedures)", 2, 2, "Basic setup, Multiple beds", "Procedure lights, Basic instruments"},
	}

	for _, ot := range theatres {
		var existingOT model.OperationTheatre
		err := db.Where("theatre_name = ?", ot.theatreName).First(&existingOT).Error

		theatre := model.OperationTheatre{
			TheatreName:   ot.theatreName,
			Floor:         ot.floor,
			Capacity:      ot.capacity,
			Status:        "available",
			Features:      ot.features,
			EquipmentList: ot.equipmentList,
		}

		if err != nil {
			if err := db.Create(&theatre).Error; err != nil {
				log.Printf("Failed to seed OT %s: %v", ot.theatreName, err)
				continue
			}
		}
		log.Printf("✓ Operation Theatre ensured: %s (Floor %d)", ot.theatreName, ot.floor)
	}

	var otCount int64
	db.Model(&model.OperationTheatre{}).Count(&otCount)
	log.Printf("📊 Total operation theatres in database: %d\n", otCount)

	return nil
}

// SeedOperationSchedules creates operation schedules
func SeedOperationSchedules(db *gorm.DB) error {
	// Get theatres, patients, and doctors
	var theatres []model.OperationTheatre
	if err := db.Find(&theatres).Error; err != nil {
		return fmt.Errorf("failed to fetch theatres: %w", err)
	}

	var patients []model.Patient
	if err := db.Limit(10).Find(&patients).Error; err != nil {
		return fmt.Errorf("failed to fetch patients: %w", err)
	}

	var doctors []model.Doctor
	if err := db.Limit(10).Find(&doctors).Error; err != nil {
		return fmt.Errorf("failed to fetch doctors: %w", err)
	}

	if len(theatres) == 0 || len(patients) == 0 || len(doctors) == 0 {
		log.Printf("Warning: Not enough data to seed operations (theatres: %d, patients: %d, doctors: %d)", len(theatres), len(patients), len(doctors))
		return nil
	}

	now := time.Now()
	operations := []struct {
		theatreIdx    int
		patientIdx    int
		surgeonIdx    int
		operationType string
		diagnosis     string
		daysOffset    int
		status        string
		timeStr       string
	}{
		{0, 0, 0, "Appendectomy", "Acute appendicitis", 1, "scheduled", "09:00"},
		{1, 1, 1, "Coronary Artery Bypass", "Triple vessel CAD", 1, "scheduled", "10:30"},
		{2, 2, 2, "Knee Replacement", "Severe osteoarthritis knee", 2, "scheduled", "14:00"},
		{3, 3, 3, "Emergency Laparotomy", "Acute abdomen", 0, "in_progress", "11:00"},
		{4, 4, 4, "Gallbladder Removal (Lap)", "Cholelithiasis", -1, "completed", "09:30"},
		{0, 5, 0, "Hernia Repair", "Inguinal hernia", -2, "completed", "15:00"},
		{1, 6, 1, "Cardiac Valve Replacement", "Aortic stenosis", 3, "scheduled", "08:00"},
		{2, 7, 2, "Hip Arthroplasty", "Hip fracture repair", 2, "scheduled", "13:00"},
	}

	for idx, op := range operations {
		theatreIdx := op.theatreIdx % len(theatres)
		patientIdx := op.patientIdx % len(patients)
		surgeonIdx := op.surgeonIdx % len(doctors)

		theatre := theatres[theatreIdx]
		patient := patients[patientIdx]
		surgeon := doctors[surgeonIdx]

		operationDate := now.AddDate(0, 0, op.daysOffset).Format("2006-01-02")

		var existingOp model.OperationSchedule
		err := db.Where("patient_id = ? AND operation_date = ? AND operation_time = ?", patient.ID, operationDate, op.timeStr).First(&existingOp).Error

		operationSchedule := model.OperationSchedule{
			TheatreID:         theatre.ID,
			PatientID:         patient.ID,
			SurgieDoctorID:    surgeon.UserID,
			OperationDate:     operationDate,
			OperationTime:     op.timeStr,
			EstimatedDuration: 120,
			OperationType:     op.operationType,
			Status:            model.OperationType(op.status),
			Diagnosis:         op.diagnosis,
			PreOperativeNotes: "Patient pre-op assessment done. All investigations normal.",
		}

		if op.status == "completed" {
			operationSchedule.ActualDuration = 135
			operationSchedule.PostOperativeNotes = "Surgery completed successfully. Patient stable. Sent to recovery."
		}

		if err != nil {
			if err := db.Create(&operationSchedule).Error; err != nil {
				log.Printf("Failed to seed operation %d: %v", idx+1, err)
				continue
			}
		}
		log.Printf("✓ Operation ensured: %s on %s at %s [%s]", op.operationType, operationDate, op.timeStr, op.status)
	}

	var operationCount int64
	db.Model(&model.OperationSchedule{}).Count(&operationCount)
	log.Printf("📊 Total operations in database: %d\n", operationCount)

	return nil
}
