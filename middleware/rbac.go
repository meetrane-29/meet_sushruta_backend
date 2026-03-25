package middleware

import (
	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(403, gin.H{
				"success": false,
				"error":   "forbidden",
				"code":    403,
			})
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			c.JSON(403, gin.H{
				"success": false,
				"error":   "forbidden",
				"code":    403,
			})
			c.Abort()
			return
		}

		// Check if user role is in allowed list
		allowed := false
		for _, r := range roles {
			if userRole == r {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(403, gin.H{
				"success": false,
				"error":   "forbidden",
				"code":    403,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
