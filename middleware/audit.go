package middleware

import (
	"fmt"
	"log"
	"strings"

	"meet_sushruta/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditMiddleware logs all API actions to audit_logs table asynchronously
func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip health checks and public endpoints
		path := c.Request.URL.Path
		if strings.Contains(path, "/health") || strings.Contains(path, "/public/") {
			c.Next()
			return
		}

		// Extract user info from context (set by AuthMiddleware)
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}
		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			c.Next()
			return
		}

		roleInterface, exists := c.Get("role")
		if !exists {
			c.Next()
			return
		}
		role, ok := roleInterface.(string)
		if !ok {
			c.Next()
			return
		}

		// Capture request details
		method := c.Request.Method
		endpoint := c.Request.URL.Path
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Execute handler
		c.Next()

		// Insert audit log asynchronously (non-blocking)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[Audit] Recovery from panic: %v", r)
				}
			}()

			// Determine action from HTTP method
			action := "READ"
			if method == "POST" {
				action = "CREATE"
			} else if method == "PATCH" || method == "PUT" {
				action = "UPDATE"
			} else if method == "DELETE" {
				action = "DELETE"
			}

			// Extract entity type from endpoint
			parts := strings.Split(strings.Trim(endpoint, "/"), "/")
			entityType := "Unknown"
			if len(parts) > 2 {
				entityType = strings.Title(parts[2]) // e.g., /api/v1/patients -> patients
			}

			auditLog := &model.AuditLog{
				ID:          uuid.New(),
				UserID:      userID,
				Action:      action,
				EntityType:  entityType,
				EntityID:    uuid.Nil, // Could be extracted from route params
				IPAddress:   ipAddress,
				UserAgent:   userAgent,
				Status:      "success",
				Description: fmt.Sprintf("%s %s by %s", method, endpoint, role),
			}

			if err := db.Create(auditLog).Error; err != nil {
				log.Printf("[Audit] Failed to insert audit log: %v", err)
			}
		}()
	}
}
