package service

import (
	"errors"
	"strings"
	"time"

	"meet_sushruta/config"
	"meet_sushruta/middleware"
	"meet_sushruta/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"user_id"`
	Role         string `json:"role"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type RegisterRequest struct {
	FirstName             string `json:"first_name" binding:"required,min=2,max=100"`
	LastName              string `json:"last_name" binding:"required,min=2,max=100"`
	Email                 string `json:"email" binding:"required,email"`
	Password              string `json:"password" binding:"required,min=8"`
	Phone                 string `json:"phone" binding:"required"`
	Role                  string `json:"role"`
	DateOfBirth           string `json:"date_of_birth"`
	Gender                string `json:"gender"`
	BloodGroup            string `json:"blood_group"`
	Address               string `json:"address"`
	EmergencyContactName  string `json:"emergency_contact_name"`
	EmergencyContactPhone string `json:"emergency_contact_phone"`
	MedicalHistory        string `json:"medical_history"`
	Allergies             string `json:"allergies"`
}

type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Message string `json:"message"`
}

type AuthService interface {
	Login(email, password string) (*AuthResponse, error)
	RefreshToken(refreshToken string) (*AuthResponse, error)
	RegisterUser(req RegisterRequest) (*RegisterResponse, error)
}

type authService struct {
	jwtSecret string
}

func NewAuthService(jwtSecret string) AuthService {
	return &authService{
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Login(email, password string) (*AuthResponse, error) {
	// Normalize email
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)

	// Find user by email
	var user model.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.Active {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Record last login time
	now := time.Now().UnixMilli()
	config.DB.Model(&user).UpdateColumn("last_login", now)

	// Generate access token (15 minutes)
	accessToken, err := s.generateAccessToken(&user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token (7 days)
	refreshToken, err := s.generateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64((time.Duration(15) * time.Minute).Seconds()),
		TokenType:    "Bearer",
		UserID:       user.ID.String(),
		Role:         user.Role,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}, nil
}

func (s *authService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	// Parse refresh token
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Get user from database
	var user model.User
	if err := config.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if !user.Active {
		return nil, errors.New("user account is inactive")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(&user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64((time.Duration(15) * time.Minute).Seconds()),
		TokenType:    "Bearer",
		UserID:       user.ID.String(),
		Role:         user.Role,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}, nil
}

func (s *authService) generateAccessToken(user *model.User) (string, error) {
	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(15) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "meet_sushruta",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	return tokenString, err
}

func (s *authService) generateRefreshToken(user *model.User) (string, error) {
	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "meet_sushruta",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	return tokenString, err
}

func (s *authService) RegisterUser(req RegisterRequest) (*RegisterResponse, error) {
	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" || req.Phone == "" {
		return nil, errors.New("missing required fields: first_name, last_name, email, password, phone")
	}

	// Normalize email and password
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	// Set default role to patient if not provided
	if req.Role == "" {
		req.Role = "patient"
	}

	// Validate role
	validRoles := map[string]bool{
		"admin":        true,
		"doctor":       true,
		"nurse":        true,
		"pharmacy":     true,
		"lab":          true,
		"patient":      true,
		"receptionist": true,
	}
	if !validRoles[req.Role] {
		return nil, errors.New("invalid role - must be admin, doctor, nurse, pharmacy, lab, patient, or receptionist")
	}

	// Check if email already exists
	var existingUser model.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	// Hash password with bcrypt (cost=12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &model.User{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   string(hashedPassword),
		Role:       req.Role,
		Active:     true,
		IsVerified: false,
	}

	if err := config.DB.Create(user).Error; err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return nil, errors.New("email already registered")
		}
		return nil, err
	}

	// If role is patient, create patient record
	if req.Role == "patient" {
		patient := &model.Patient{
			UserID:           user.ID,
			DateOfBirth:      req.DateOfBirth,
			Gender:           req.Gender,
			BloodGroup:       req.BloodGroup,
			Address:          req.Address,
			EmergencyContact: req.EmergencyContactName, // Combining name and phone into contact field for simplicity
			MedicalHistory:   req.MedicalHistory,
			Allergies:        req.Allergies,
		}

		if err := config.DB.Create(patient).Error; err != nil {
			// If patient creation fails, delete the created user
			config.DB.Delete(user)
			return nil, errors.New("failed to create patient record")
		}
	}

	// Log audit entry (optional - can be enhanced later)
	// For now, skip audit logging if it fails
	auditLog := &model.AuditLog{
		UserID:     user.ID,
		Action:     "USER_REGISTERED",
		EntityType: "User",
		EntityID:   user.ID,
		NewValues:  user.Email + " (" + user.Role + ")",
		Status:     "success",
	}
	// Don't fail registration if audit logging fails
	if err := config.DB.Create(auditLog).Error; err != nil {
		// Log the error but continue (audit is non-critical)
		// In production, you might want to log this error to a separate logging system
	}

	return &RegisterResponse{
		UserID:  user.ID.String(),
		Email:   user.Email,
		Role:    user.Role,
		Message: "Account created successfully",
	}, nil
}
