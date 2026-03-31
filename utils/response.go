package utils

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

// ValidatePhone validates that phone number is exactly 10 digits
func ValidatePhone(phone string) bool {
	// Check if phone has exactly 10 digits
	pattern := `^\d{10}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"success": true,
		"data":    data,
	})
}

func OKWithMeta(c *gin.Context, data interface{}, meta interface{}) {
	c.JSON(200, gin.H{
		"success": true,
		"data":    data,
		"meta":    meta,
	})
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"success": false,
		"message": msg,
		"error":   msg,
		"code":    code,
	})
}
