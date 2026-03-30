package handler

import (
	"meet_sushruta/service"
	"meet_sushruta/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WaitingListHandler struct {
	waitingListService service.WaitingListService
}

func NewWaitingListHandler(waitingListService service.WaitingListService) *WaitingListHandler {
	return &WaitingListHandler{
		waitingListService: waitingListService,
	}
}

type AddToWaitingListRequest struct {
	AppointmentID uuid.UUID `json:"appointment_id" binding:"required"`
	PatientID     uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID      uuid.UUID `json:"doctor_id" binding:"required"`
}

type UpdateWaitingListStatusRequest struct {
	Status string `json:"status" binding:"required"` // waiting, called, seen, completed, cancelled
}

// AddToWaitingList adds patient to OPD waiting list
// POST /api/v1/waiting-list/add
func (h *WaitingListHandler) AddToWaitingList(c *gin.Context) {
	var req AddToWaitingListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "invalid request: appointment_id, patient_id, doctor_id are required")
		return
	}

	entry, err := h.waitingListService.AddToWaitingList(req.AppointmentID, req.PatientID, req.DoctorID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, entry)
}

// GetDoctorWaitingList retrieves waiting list for a doctor
// GET /api/v1/waiting-list/doctor/:doctor_id
func (h *WaitingListHandler) GetDoctorWaitingList(c *gin.Context) {
	doctorIDStr := c.Param("doctor_id")
	doctorID, err := uuid.Parse(doctorIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor ID")
		return
	}

	entries, err := h.waitingListService.GetDoctorWaitingList(doctorID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"waiting_list": entries,
		"total":        len(entries),
	})
}

// CallNextPatient calls next patient in queue
// POST /api/v1/waiting-list/call-next/:doctor_id
func (h *WaitingListHandler) CallNextPatient(c *gin.Context) {
	doctorIDStr := c.Param("doctor_id")
	doctorID, err := uuid.Parse(doctorIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor ID")
		return
	}

	nextPatient, err := h.waitingListService.CallNextPatient(doctorID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, nextPatient)
}

// MarkPatientSeen marks patient as seen/in consultation
// PUT /api/v1/waiting-list/:id/mark-seen
func (h *WaitingListHandler) MarkPatientSeen(c *gin.Context) {
	entryIDStr := c.Param("id")
	entryID, err := uuid.Parse(entryIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid entry ID")
		return
	}

	if err := h.waitingListService.MarkPatientSeen(entryID); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{"message": "patient marked as seen"})
}

// CompleteConsultation marks consultation as completed
// PUT /api/v1/waiting-list/:id/complete
func (h *WaitingListHandler) CompleteConsultation(c *gin.Context) {
	entryIDStr := c.Param("id")
	entryID, err := uuid.Parse(entryIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid entry ID")
		return
	}

	if err := h.waitingListService.CompleteConsultation(entryID); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{"message": "consultation completed"})
}

// CancelWaitingListEntry cancels patient's waiting list entry
// DELETE /api/v1/waiting-list/:id/cancel
func (h *WaitingListHandler) CancelWaitingListEntry(c *gin.Context) {
	entryIDStr := c.Param("id")
	entryID, err := uuid.Parse(entryIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid entry ID")
		return
	}

	if err := h.waitingListService.CancelWaitingListEntry(entryID); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.OK(c, gin.H{"message": "waiting list entry cancelled"})
}

// GetPatientQueuePosition gets patient's current queue position
// GET /api/v1/waiting-list/patient/:patient_id/position
func (h *WaitingListHandler) GetPatientQueuePosition(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid patient ID")
		return
	}

	position, entry, err := h.waitingListService.GetPatientQueuePosition(patientID)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"token_number":        position,
		"status":              entry.Status,
		"estimated_wait_time": entry.EstimatedWaitTime,
	})
}

// GetEstimatedWaitTime gets estimated wait time for a doctor
// GET /api/v1/waiting-list/doctor/:doctor_id/estimated-time
func (h *WaitingListHandler) GetEstimatedWaitTime(c *gin.Context) {
	doctorIDStr := c.Param("doctor_id")
	doctorID, err := uuid.Parse(doctorIDStr)
	if err != nil {
		utils.Fail(c, 400, "invalid doctor ID")
		return
	}

	estimatedTime := h.waitingListService.GetEstimatedWaitTime(doctorID)
	utils.OK(c, gin.H{
		"estimated_wait_time_minutes": estimatedTime,
	})
}
