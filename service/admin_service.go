package service

import (
	"errors"
	"strings"

	"meet_sushruta/config"
	"meet_sushruta/model"
	"meet_sushruta/repository"

	"golang.org/x/crypto/bcrypt"
)

type RegisterDoctorRequest struct {
	FirstName            string  `json:"first_name"`
	LastName             string  `json:"last_name"`
	Email                string  `json:"email"`
	Password             string  `json:"password"`
	Phone                string  `json:"phone"`
	Specialization       string  `json:"specialization"`
	LicenseNumber        string  `json:"license_number"`
	CertificationURL     string  `json:"certification_url"`
	Bio                  string  `json:"bio"`
	Department           string  `json:"department"`
	Fees                 float64 `json:"fees"`
	JoiningDate          int64   `json:"joining_date"`
	Salary               float64 `json:"salary"`
	AttendancePercentage float64 `json:"attendance_percentage"`
	LeaveBalance         int     `json:"leave_balance"`
}

type AdminService interface {
	GetDashboardStats() (map[string]interface{}, error)
	RegisterDoctor(req RegisterDoctorRequest) (map[string]interface{}, error)
	GetTodayAppointmentsByDoctor() ([]map[string]interface{}, error)
	GetBedStats() (map[string]interface{}, error)
	GetAllBeds(page, limit int) ([]model.Bed, int64, error)
	GetPatientsByBedType(bedType string) ([]map[string]interface{}, error)
	GetMedicalEquipmentStats() (map[string]interface{}, error)
	GetAllMedicalEquipment(page, limit int, status string) ([]model.MedicalEquipment, int64, error)
	GetOperationTheatreStats() (map[string]interface{}, error)
	GetOperationSchedules(page, limit int, date string) ([]model.OperationSchedule, int64, error)
}

type adminService struct {
	repo                  repository.AdminRepository
	bedRepo               repository.BedRepository
	medicalEquipmentRepo  repository.MedicalEquipmentRepository
	operationTheatreRepo  repository.OperationTheatreRepository
	operationScheduleRepo repository.OperationScheduleRepository
}

func NewAdminService(repo repository.AdminRepository, bedRepo repository.BedRepository, medicalEquipmentRepo repository.MedicalEquipmentRepository, operationTheatreRepo repository.OperationTheatreRepository, operationScheduleRepo repository.OperationScheduleRepository) AdminService {
	return &adminService{
		repo:                  repo,
		bedRepo:               bedRepo,
		medicalEquipmentRepo:  medicalEquipmentRepo,
		operationTheatreRepo:  operationTheatreRepo,
		operationScheduleRepo: operationScheduleRepo,
	}
}

// GetDashboardStats returns all dashboard statistics
func (s *adminService) GetDashboardStats() (map[string]interface{}, error) {
	// Get all stats in parallel-like structure
	totalPatients, err := s.repo.GetTotalPatients()
	if err != nil {
		return nil, err
	}

	totalDoctors, err := s.repo.GetTotalDoctors()
	if err != nil {
		return nil, err
	}

	totalNurses, err := s.repo.GetTotalNurses()
	if err != nil {
		return nil, err
	}

	totalUsers, err := s.repo.GetTotalUsers()
	if err != nil {
		return nil, err
	}

	todayAppointments, err := s.repo.GetTodayAppointments()
	if err != nil {
		return nil, err
	}

	recentPatients, err := s.repo.GetRecentPatients()
	if err != nil {
		return nil, err
	}
	if recentPatients == nil {
		recentPatients = []map[string]interface{}{}
	}

	topDoctors, err := s.repo.GetTopDoctorsByAppointmentsToday()
	if err != nil {
		return nil, err
	}
	if topDoctors == nil {
		topDoctors = []map[string]interface{}{}
	}

	departmentLoad, err := s.repo.GetDepartmentLoad()
	if err != nil {
		return nil, err
	}
	if departmentLoad == nil {
		departmentLoad = []map[string]interface{}{}
	}

	nurseRoleBreakdown, err := s.repo.GetNurseRoleBreakdown()
	if err != nil {
		return nil, err
	}
	if nurseRoleBreakdown == nil {
		nurseRoleBreakdown = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"stats": map[string]interface{}{
			"total_patients":     totalPatients,
			"total_doctors":      totalDoctors,
			"total_nurses":       totalNurses,
			"today_appointments": todayAppointments,
			"total_users":        totalUsers,
		},
		"recent_patients":      recentPatients,
		"top_doctors":          topDoctors,
		"department_load":      departmentLoad,
		"nurse_role_breakdown": nurseRoleBreakdown,
	}, nil
}

// RegisterDoctor creates a new doctor user and doctor record
func (s *adminService) RegisterDoctor(req RegisterDoctorRequest) (map[string]interface{}, error) {
	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" || req.Phone == "" || req.Specialization == "" || req.LicenseNumber == "" {
		return nil, errors.New("missing required fields: first_name, last_name, email, password, phone, specialization, license_number are required")
	}

	// Normalize email and password
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := strings.TrimSpace(req.Password)
	phone := strings.TrimSpace(req.Phone)

	// Check if email already exists
	var existingUser model.User
	if err := config.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	// Check if phone already exists
	var existingPhone model.User
	if err := config.DB.Where("phone = ?", phone).First(&existingPhone).Error; err == nil {
		return nil, errors.New("phone number already registered")
	}

	// Check if license number already exists
	var existingLicense model.Doctor
	if err := config.DB.Where("license_number = ?", req.LicenseNumber).First(&existingLicense).Error; err == nil {
		return nil, errors.New("license number already registered")
	}

	// Hash password with bcrypt (cost=12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user first
	user := &model.User{
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		Email:                email,
		Phone:                phone,
		Password:             string(hashedPassword),
		Role:                 "doctor",
		Active:               true,
		IsVerified:           false,
		JoiningDate:          req.JoiningDate,
		Salary:               req.Salary,
		AttendancePercentage: req.AttendancePercentage,
		LeaveBalance:         req.LeaveBalance,
	}

	// Use transaction to ensure both user and doctor are created together
	tx := config.DB.Begin()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "email") {
			return nil, errors.New("email already registered")
		}
		if strings.Contains(err.Error(), "phone") {
			return nil, errors.New("phone number already registered")
		}
		return nil, errors.New("failed to create user: " + err.Error())
	}

	// Create doctor record
	doctor := &model.Doctor{
		UserID:               user.ID,
		Specialization:       req.Specialization,
		LicenseNumber:        req.LicenseNumber,
		CertificationURL:     req.CertificationURL,
		Bio:                  req.Bio,
		Department:           req.Department,
		ConsultationFee:      req.Fees,
		JoiningDate:          req.JoiningDate,
		Salary:               req.Salary,
		AttendancePercentage: req.AttendancePercentage,
		LeaveBalance:         req.LeaveBalance,
	}

	if err := tx.Create(doctor).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create doctor record: " + err.Error())
	}

	// Commit transaction (user and doctor)
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit transaction: " + err.Error())
	}

	// Log audit entry AFTER commit (non-critical)
	auditLog := &model.AuditLog{
		UserID:     user.ID,
		Action:     "DOCTOR_REGISTERED",
		EntityType: "Doctor",
		EntityID:   doctor.ID,
		NewValues:  email + " (" + req.Specialization + ")",
		Status:     "success",
	}

	// Don't fail if audit logging fails - it's non-critical
	if err := config.DB.Create(auditLog).Error; err != nil {
		// Just ignore audit error - user and doctor are already saved
	}

	return map[string]interface{}{
		"doctor_id":      doctor.ID.String(),
		"user_id":        user.ID.String(),
		"email":          user.Email,
		"name":           req.FirstName + " " + req.LastName,
		"role":           user.Role,
		"specialization": req.Specialization,
		"message":        "Doctor registered successfully",
	}, nil
}

// GetTodayAppointmentsByDoctor returns all appointments for today grouped by doctor
func (s *adminService) GetTodayAppointmentsByDoctor() ([]map[string]interface{}, error) {
	appointments, err := s.repo.GetTodayAppointmentsByDoctor()
	if err != nil {
		return nil, err
	}
	if appointments == nil {
		appointments = []map[string]interface{}{}
	}
	return appointments, nil
}

// GetBedStats returns bed statistics
func (s *adminService) GetBedStats() (map[string]interface{}, error) {
	return s.bedRepo.GetBedStats()
}

// GetAllBeds returns all beds with pagination
func (s *adminService) GetAllBeds(page, limit int) ([]model.Bed, int64, error) {
	return s.bedRepo.GetAll(page, limit)
}

// GetPatientsByBedType returns patients in a specific bed type with details
func (s *adminService) GetPatientsByBedType(bedType string) ([]map[string]interface{}, error) {
	beds, err := s.bedRepo.GetByBedType(bedType)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, bed := range beds {
		bedData := map[string]interface{}{
			"id":            bed.ID,
			"bed_number":    bed.BedNumber,
			"bed_type":      bed.BedType,
			"status":        bed.Status,
			"admitted_at":   bed.AdmittedAt,
			"discharged_at": bed.DischargedAt,
			"daily_rate":    bed.DailyRate,
		}

		// Add patient details if occupied
		if bed.Patient != nil {
			bedData["patient"] = map[string]interface{}{
				"id":            bed.Patient.ID,
				"gender":        bed.Patient.Gender,
				"blood_group":   bed.Patient.BloodGroup,
				"date_of_birth": bed.Patient.DateOfBirth,
				"address":       bed.Patient.Address,
			}
			if bed.Patient.User != nil {
				bedData["patient_name"] = bed.Patient.User.FirstName + " " + bed.Patient.User.LastName
				bedData["patient_email"] = bed.Patient.User.Email
				bedData["patient_phone"] = bed.Patient.User.Phone
			}
		}

		result = append(result, bedData)
	}

	return result, nil
}

// GetMedicalEquipmentStats returns medical equipment statistics
func (s *adminService) GetMedicalEquipmentStats() (map[string]interface{}, error) {
	return s.medicalEquipmentRepo.GetEquipmentStats()
}

// GetAllMedicalEquipment returns all medical equipment with optional filtering
func (s *adminService) GetAllMedicalEquipment(page, limit int, status string) ([]model.MedicalEquipment, int64, error) {
	equipment, total, err := s.medicalEquipmentRepo.GetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// If status filter is provided, filter the results in-memory
	if status != "" {
		var filtered []model.MedicalEquipment
		for _, e := range equipment {
			if e.Status == status {
				filtered = append(filtered, e)
			}
		}
		return filtered, total, nil
	}

	return equipment, total, nil
}

// GetOperationTheatreStats returns operation theatre statistics
func (s *adminService) GetOperationTheatreStats() (map[string]interface{}, error) {
	// Get theatre stats
	theatres, _, err := s.operationTheatreRepo.GetAll(1, 100)
	if err != nil {
		return nil, err
	}

	// Get operation schedule stats
	scheduleStats, err := s.operationScheduleRepo.GetOperationStats()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_theatres":   len(theatres),
		"operations_stats": scheduleStats,
	}

	// Add theatre status breakdown
	availableCount := 0
	inUseCount := 0
	maintenanceCount := 0
	for _, theatre := range theatres {
		switch theatre.Status {
		case "available":
			availableCount++
		case "in_use":
			inUseCount++
		case "maintenance":
			maintenanceCount++
		}
	}

	stats["theatre_status"] = map[string]interface{}{
		"available":   availableCount,
		"in_use":      inUseCount,
		"maintenance": maintenanceCount,
	}

	return stats, nil
}

// GetOperationSchedules returns operation schedules with optional date filter
func (s *adminService) GetOperationSchedules(page, limit int, date string) ([]model.OperationSchedule, int64, error) {
	if date != "" {
		operations, err := s.operationScheduleRepo.GetByDate(date)
		if err != nil {
			return nil, 0, err
		}
		return operations, int64(len(operations)), nil
	}

	return s.operationScheduleRepo.GetAll(page, limit)
}
