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
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	Email            string  `json:"email"`
	Password         string  `json:"password"`
	Phone            string  `json:"phone"`
	Specialization   string  `json:"specialization"`
	LicenseNumber    string  `json:"license_number"`
	CertificationURL string  `json:"certification_url"`
	Bio              string  `json:"bio"`
	Department       string  `json:"department"`
	Fees             float64 `json:"fees"`
}

type AdminService interface {
	GetDashboardStats() (map[string]interface{}, error)
	RegisterDoctor(req RegisterDoctorRequest) (map[string]interface{}, error)
}

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) AdminService {
	return &adminService{
		repo: repo,
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

	return map[string]interface{}{
		"stats": map[string]interface{}{
			"total_patients":     totalPatients,
			"total_doctors":      totalDoctors,
			"today_appointments": todayAppointments,
			"total_users":        totalUsers,
		},
		"recent_patients": recentPatients,
		"top_doctors":     topDoctors,
		"department_load": departmentLoad,
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
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      email,
		Phone:      phone,
		Password:   string(hashedPassword),
		Role:       "doctor",
		Active:     true,
		IsVerified: false,
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
		UserID:           user.ID,
		Specialization:   req.Specialization,
		LicenseNumber:    req.LicenseNumber,
		CertificationURL: req.CertificationURL,
		Bio:              req.Bio,
		Department:       req.Department,
		ConsultationFee:  req.Fees,
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
