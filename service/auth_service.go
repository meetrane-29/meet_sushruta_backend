package service

import (
	"errors"
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
}

type AuthService interface {
	Login(email, password string) (*AuthResponse, error)
	RefreshToken(refreshToken string) (*AuthResponse, error)
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
