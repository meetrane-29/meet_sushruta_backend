package handler

import (
	"strings"

	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	resp, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.Fail(c, 401, err.Error())
		return
	}

	utils.OK(c, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	resp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.Fail(c, 401, err.Error())
		return
	}

	utils.OK(c, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT is stateless, logout is just a client-side action
	// Return success response
	utils.OK(c, gin.H{
		"message": "logged out successfully",
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request")
		return
	}

	// Trim and validate email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Basic field validation
	if len(req.FirstName) < 2 || len(req.FirstName) > 100 {
		utils.Fail(c, 400, "first name must be between 2 and 100 characters")
		return
	}

	if len(req.LastName) < 2 || len(req.LastName) > 100 {
		utils.Fail(c, 400, "last name must be between 2 and 100 characters")
		return
	}

	if len(req.Password) < 8 {
		utils.Fail(c, 400, "password must be at least 8 characters")
		return
	}

	if len(req.Phone) == 0 || len(req.Phone) < 10 {
		utils.Fail(c, 400, "phone must be at least 10 characters")
		return
	}

	// Blood group validation if provided
	if req.BloodGroup != "" {
		validBloodGroups := map[string]bool{
			"A+": true, "A-": true, "B+": true, "B-": true,
			"O+": true, "O-": true, "AB+": true, "AB-": true,
		}
		if !validBloodGroups[req.BloodGroup] {
			utils.Fail(c, 400, "invalid blood group")
			return
		}
	}

	// Gender validation if provided
	if req.Gender != "" {
		validGenders := map[string]bool{"male": true, "female": true, "other": true}
		if !validGenders[strings.ToLower(req.Gender)] {
			utils.Fail(c, 400, "invalid gender")
			return
		}
	}

	// Call service to register user
	resp, err := h.authService.RegisterUser(req)
	if err != nil {
		// Check what kind of error occurred
		errorMsg := err.Error()

		// Email already registered
		if strings.Contains(errorMsg, "email already registered") {
			utils.Fail(c, 409, "This email is already registered. Please login or use a different email address.")
			return
		}

		// Invalid role or other validation errors
		if strings.Contains(errorMsg, "invalid") {
			utils.Fail(c, 400, errorMsg)
			return
		}

		// Server errors
		utils.Fail(c, 500, "Registration failed. Please try again later.")
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"data":    resp,
	})
}
