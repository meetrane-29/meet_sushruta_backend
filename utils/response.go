package utils

import (
	"github.com/gin-gonic/gin"
)

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
		"error":   msg,
		"code":    code,
	})
}
