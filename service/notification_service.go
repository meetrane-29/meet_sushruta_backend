package service

import (
	"log"

	"meet_sushruta/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationJob represents a notification task in the queue
type NotificationJob struct {
	RecipientID uuid.UUID
	Type        string // "appointment", "prescription", "lab_result", "bill", "general"
	Message     string
	Channel     string // "in_app", "email", "sms"
	Title       string
	RelatedID   *uuid.UUID
	RelatedType string
}

// NotificationService handles async notification delivery
type NotificationService struct {
	queue chan NotificationJob
	db    *gorm.DB
}

// NewNotificationService creates a new notification service with a buffered queue
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{
		queue: make(chan NotificationJob, 100), // Buffer 100 jobs
		db:    db,
	}
}

// StartWorker starts the background worker goroutine that processes notifications
func (s *NotificationService) StartWorker() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Notification] Worker recovered from panic: %v", r)
				// Restart worker
				s.StartWorker()
			}
		}()

		log.Println("[Notification] Worker started")

		for job := range s.queue {
			s.processNotification(job)
		}
	}()
}

// Send adds a notification job to the queue (non-blocking)
func (s *NotificationService) Send(job NotificationJob) {
	select {
	case s.queue <- job:
		log.Printf("[Notification] Job queued for user %s", job.RecipientID.String())
	default:
		log.Printf("[Notification] Queue full, dropping job for user %s", job.RecipientID.String())
	}
}

// processNotification handles a single notification by inserting to DB and any external API calls
func (s *NotificationService) processNotification(job NotificationJob) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Notification] Recovered from panic while processing job: %v", r)
		}
	}()

	// Create notification record
	notification := &model.Notification{
		ID:          uuid.New(),
		UserID:      job.RecipientID,
		Title:       job.Title,
		Message:     job.Message,
		Type:        job.Type,
		RelatedID:   job.RelatedID,
		RelatedType: job.RelatedType,
		IsRead:      false,
		IsPushSent:  false,
		IsEmailSent: false,
		Priority:    "normal",
	}

	// Insert into database
	if err := s.db.Create(notification).Error; err != nil {
		log.Printf("[Notification] Failed to save notification: %v", err)
		return
	}

	// Log to console (easy to swap with SMS/Email API later)
	log.Printf("[Notification] Sent %s notification to user %s: %s", job.Type, job.RecipientID.String(), job.Message)

	// TODO: Add external integrations here
	// - SMS API call for Channel="sms"
	// - Email API call for Channel="email"
	// - Push notification for Channel="push"
}
