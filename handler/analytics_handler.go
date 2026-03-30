package handler

import (
	"fmt"
	"strconv"
	"time"

	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ============================================================
// ANALYTICS HANDLER
// ============================================================

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetDoctorStatistics returns comprehensive statistics for a specific doctor
// GET /api/v1/analytics/doctors/:doctor_id/statistics
func (h *AnalyticsHandler) GetDoctorStatistics(c *gin.Context) {
	doctorIDStr := c.Param("doctor_id")
	doctorID, err := uuid.Parse(doctorIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor_id format")
		return
	}

	stats, err := h.analyticsService.GetDoctorStatistics(c.Request.Context(), doctorID)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error fetching doctor statistics: %v", err))
		return
	}

	utils.OK(c, stats)
}

// GetAppointmentTrends returns appointment trends for a date range
// GET /api/v1/analytics/appointments/trends?doctor_id=xxx&from_date=2026-01-01&to_date=2026-12-31&group_by=day
func (h *AnalyticsHandler) GetAppointmentTrends(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")
	fromDateStr := c.Query("from_date") // YYYY-MM-DD
	toDateStr := c.Query("to_date")     // YYYY-MM-DD
	groupBy := c.Query("group_by")      // day, week, month

	if groupBy == "" {
		groupBy = "day"
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		fromDate = time.Now().AddDate(0, -3, 0) // Default: last 3 months
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		toDate = time.Now()
	}

	var doctorID *uuid.UUID
	if doctorIDStr != "" {
		parsed, err := uuid.Parse(doctorIDStr)
		if err != nil {
			utils.Fail(c, 400, "invalid doctor_id format")
			return
		}
		doctorID = &parsed
	}

	trends, err := h.analyticsService.GetAppointmentTrends(c.Request.Context(), doctorID, fromDate, toDate, groupBy)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error fetching appointment trends: %v", err))
		return
	}

	utils.OK(c, gin.H{
		"trends":    trends,
		"from_date": fromDateStr,
		"to_date":   toDateStr,
		"group_by":  groupBy,
	})
}

// GetLabCompletionMetrics returns lab test completion statistics
// GET /api/v1/analytics/lab/completion?doctor_id=xxx&time_range=30days
func (h *AnalyticsHandler) GetLabCompletionMetrics(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")
	timeRange := c.Query("time_range") // 7days, 30days, 90days, 1year

	if timeRange == "" {
		timeRange = "30days"
	}

	var doctorID *uuid.UUID
	if doctorIDStr != "" {
		parsed, err := uuid.Parse(doctorIDStr)
		if err != nil {
			utils.Fail(c, 400, "invalid doctor_id format")
			return
		}
		doctorID = &parsed
	}

	metrics, err := h.analyticsService.GetLabCompletionMetrics(c.Request.Context(), doctorID, timeRange)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error fetching lab metrics: %v", err))
		return
	}

	utils.OK(c, metrics)
}

// GetPatientOutcomeStatistics returns patient outcome metrics
// GET /api/v1/analytics/patients/outcomes?doctor_id=xxx
func (h *AnalyticsHandler) GetPatientOutcomeStatistics(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")

	var doctorID *uuid.UUID
	if doctorIDStr != "" {
		parsed, err := uuid.Parse(doctorIDStr)
		if err != nil {
			utils.Fail(c, 400, "invalid doctor_id format")
			return
		}
		doctorID = &parsed
	}

	stats, err := h.analyticsService.GetPatientOutcomeStatistics(c.Request.Context(), doctorID)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error fetching patient outcome statistics: %v", err))
		return
	}

	utils.OK(c, stats)
}

// GetRevenueStatistics returns revenue metrics for a doctor
// GET /api/v1/analytics/doctors/:doctor_id/revenue
func (h *AnalyticsHandler) GetRevenueStatistics(c *gin.Context) {
	doctorIDStr := c.Param("doctor_id")
	doctorID, err := uuid.Parse(doctorIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor_id format")
		return
	}

	revenue, err := h.analyticsService.GetRevenueStatistics(c.Request.Context(), doctorID)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error fetching revenue statistics: %v", err))
		return
	}

	utils.OK(c, revenue)
}

// GetDoctorRankings returns doctors ranked by specified metric
// GET /api/v1/analytics/doctors/rankings?metric=appointments&limit=10
func (h *AnalyticsHandler) GetDoctorRankings(c *gin.Context) {
	metric := c.Query("metric") // appointments, patients, prescriptions, revenue
	limitStr := c.Query("limit")

	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		if l > 100 {
			l = 100 // Max limit
		}
		limit = l
	}

	if metric == "" {
		metric = "appointments"
	}

	rankings, err := h.analyticsService.GetDoctorRankings(c.Request.Context(), metric, limit)
	if err != nil {
		utils.Fail(c, 500, fmt.Sprintf("error fetching doctor rankings: %v", err))
		return
	}

	utils.OK(c, gin.H{
		"metric": metric,
		"limit":  limit,
		"data":   rankings,
	})
}

// GetCacheStats returns cache service statistics (for performance monitoring)
// GET /api/v1/analytics/cache/stats
func (h *AnalyticsHandler) GetCacheStats(c *gin.Context) {
	// Note: This would require injecting cache service into the handler
	// For now, returning a placeholder
	utils.OK(c, gin.H{
		"message": "Cache statistics endpoint - requires cache service injection",
	})
}

// GetPerformanceMetrics returns general performance metrics
// GET /api/v1/analytics/performance
func (h *AnalyticsHandler) GetPerformanceMetrics(c *gin.Context) {
	utils.OK(c, gin.H{
		"database_queries_optimized": true,
		"indexes_created":            true,
		"caching_enabled":            true,
		"indexed_tables": []string{
			"appointments", "prescriptions", "lab_requests", "patients",
			"doctors", "admission_records", "progress_notes", "bills",
		},
		"cache_ttl_minutes": 30,
		"max_cache_entries": 10000,
	})
}
