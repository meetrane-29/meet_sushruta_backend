package handler

import (
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
