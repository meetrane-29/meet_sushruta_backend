package handler

import (
	"fmt"
	"strconv"

	"meet_sushruta/model"
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RatingHandler struct {
	ratingService service.RatingService
}

func NewRatingHandler(ratingService service.RatingService) *RatingHandler {
	return &RatingHandler{
		ratingService: ratingService,
	}
}

type CreateRatingRequest struct {
	DoctorID        string `json:"doctor_id" binding:"required,uuid"`
	PatientID       string `json:"patient_id" binding:"required,uuid"`
	Rating          int    `json:"rating" binding:"required,min=1,max=5"`
	Professionalism int    `json:"professionalism" binding:"required,min=1,max=5"`
	Communication   int    `json:"communication" binding:"required,min=1,max=5"`
	Punctuality     int    `json:"punctuality" binding:"required,min=1,max=5"`
	Cleanliness     int    `json:"cleanliness" binding:"required,min=1,max=5"`
	Comment         string `json:"comment"`
}

type UpdateRatingRequest struct {
	Rating          int    `json:"rating" binding:"min=1,max=5"`
	Professionalism int    `json:"professionalism" binding:"min=1,max=5"`
	Communication   int    `json:"communication" binding:"min=1,max=5"`
	Punctuality     int    `json:"punctuality" binding:"min=1,max=5"`
	Cleanliness     int    `json:"cleanliness" binding:"min=1,max=5"`
	Comment         string `json:"comment"`
}

// CreateRating creates a new rating
// POST /api/v1/ratings
func (h *RatingHandler) CreateRating(c *gin.Context) {
	var req CreateRatingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[CreateRating] Binding error: %v\n", err.Error())
		utils.Fail(c, 400, fmt.Sprintf("invalid request: %v", err.Error()))
		return
	}

	// Convert string UUIDs to uuid.UUID
	doctorID, err := uuid.Parse(req.DoctorID)
	if err != nil {
		utils.Fail(c, 400, fmt.Sprintf("invalid doctor_id: %v", err.Error()))
		return
	}

	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		utils.Fail(c, 400, fmt.Sprintf("invalid patient_id: %v", err.Error()))
		return
	}

	rating := &model.Rating{
		DoctorID:        doctorID,
		PatientID:       patientID,
		Rating:          req.Rating,
		Professionalism: req.Professionalism,
		Communication:   req.Communication,
		Punctuality:     req.Punctuality,
		Cleanliness:     req.Cleanliness,
		Comment:         req.Comment,
	}

	fmt.Printf("[CreateRating] Creating rating: %+v\n", rating)

	err = h.ratingService.CreateRating(rating)
	if err != nil {
		fmt.Printf("[CreateRating] Error creating rating: %v\n", err.Error())
		utils.Fail(c, 400, err.Error())
		return
	}

	fmt.Printf("[CreateRating] Rating created successfully with ID: %s\n", rating.ID)

	utils.OK(c, rating)
}

// GetRating retrieves a rating by ID
// GET /api/v1/ratings/:id
func (h *RatingHandler) GetRating(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid rating id")
		return
	}

	rating, err := h.ratingService.GetRating(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, rating)
}

// ListRatings retrieves all ratings with pagination
// GET /api/v1/ratings
func (h *RatingHandler) ListRatings(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	ratings, total, err := h.ratingService.ListRatings(page, limit)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"data":  ratings,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GetDoctorRatings retrieves all ratings for a specific doctor
// GET /api/v1/ratings/doctor/:doctor_id
func (h *RatingHandler) GetDoctorRatings(c *gin.Context) {
	doctorIDStr := c.Param("doctor_id")
	doctorID, err := uuid.Parse(doctorIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor_id")
		return
	}

	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	ratings, total, err := h.ratingService.GetDoctorRatings(doctorID, page, limit)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"data":  ratings,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GetPatientRatings retrieves all ratings by a specific patient
// GET /api/v1/ratings/patient/:patient_id
func (h *RatingHandler) GetPatientRatings(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient_id")
		return
	}

	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	ratings, total, err := h.ratingService.GetPatientRatings(patientID, page, limit)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"data":  ratings,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// UpdateRating updates an existing rating
// PUT /api/v1/ratings/:id
func (h *RatingHandler) UpdateRating(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid rating id")
		return
	}

	var req UpdateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, fmt.Sprintf("invalid request: %v", err.Error()))
		return
	}

	// Get existing rating
	rating, err := h.ratingService.GetRating(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	// Update fields if provided
	if req.Rating > 0 {
		rating.Rating = req.Rating
	}
	if req.Professionalism > 0 {
		rating.Professionalism = req.Professionalism
	}
	if req.Communication > 0 {
		rating.Communication = req.Communication
	}
	if req.Punctuality > 0 {
		rating.Punctuality = req.Punctuality
	}
	if req.Cleanliness > 0 {
		rating.Cleanliness = req.Cleanliness
	}
	if req.Comment != "" {
		rating.Comment = req.Comment
	}

	err = h.ratingService.UpdateRating(rating)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, rating)
}

// DeleteRating deletes a rating
// DELETE /api/v1/ratings/:id
func (h *RatingHandler) DeleteRating(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, 400, "invalid rating id")
		return
	}

	err = h.ratingService.DeleteRating(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"message": "rating deleted successfully",
	})
}
