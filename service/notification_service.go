package service

import (
	"log"
	"os"

	"meet_sushruta/model"

	"github.com/google/uuid"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

	// Handle external integrations based on channel
	switch job.Channel {
	case "email":
		s.sendEmailNotification(notification)
	case "sms":
		s.sendSMSNotification(notification)
	case "push":
		s.sendPushNotification(notification)
	case "in_app":
		// Already saved to database above
		log.Printf("[Notification] In-app notification ready for user %s", notification.UserID.String())
	default:
		log.Printf("[Notification] Unknown channel: %s", job.Channel)
	}
}

// sendEmailNotification sends email via SendGrid
func (s *NotificationService) sendEmailNotification(notification *model.Notification) {
	user := &model.User{}
	if err := s.db.Where("id = ?", notification.UserID).First(user).Error; err != nil {
		log.Printf("[Email] Failed to fetch user: %v", err)
		return
	}

	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Printf("[Email] SENDGRID_API_KEY not set, skipping email to %s", user.Email)
		return
	}

	from := mail.NewEmail("Meet Sushruta", "noreply@meetsushruta.com")
	to := mail.NewEmail(user.FirstName+" "+user.LastName, user.Email)
	message := mail.NewSingleEmail(from, notification.Title, to, notification.Message, notification.Message)

	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Printf("[Email] SendGrid error for user %s: %v", user.Email, err)
		return
	}
	if response.StatusCode >= 400 {
		log.Printf("[Email] SendGrid returned status %d for user %s", response.StatusCode, user.Email)
		return
	}

	s.db.Model(notification).Update("is_email_sent", true)
	log.Printf("[Email] Sent to %s: [%s]", user.Email, notification.Title)
}

// sendSMSNotification sends SMS via configured provider
// TODO: Integrate with real SMS service (Twilio, AWS SNS, etc.)
func (s *NotificationService) sendSMSNotification(notification *model.Notification) {
	// Placeholder: Get user phone from database
	user := &model.User{}
	if err := s.db.Where("id = ?", notification.UserID).First(user).Error; err != nil {
		log.Printf("[SMS] Failed to fetch user: %v", err)
		return
	}

	log.Printf("[SMS] Would send SMS to %s: %s", user.Phone, notification.Message)

	// TODO: Uncomment when SMS service is configured
	// Example with Twilio:
	// err := twilioClient.SendSMS(user.Phone, notification.Message)
	// if err == nil {
	//     s.db.Model(notification).Update("is_push_sent", true)
	// }
}

// sendPushNotification sends push notification via configured provider
// TODO: Integrate with real push service (Firebase Cloud Messaging, etc.)
func (s *NotificationService) sendPushNotification(notification *model.Notification) {
	log.Printf("[Push] Would send push notification to user %s: [%s] %s",
		notification.UserID.String(), notification.Title, notification.Message)

	// TODO: Uncomment when push service is configured
	// Example with Firebase:
	// topic := fmt.Sprintf("user_%s", notification.UserID.String())
	// msg := &messaging.Message{
	//     Topic: topic,
	//     Data: map[string]string{"type": notification.Type},
	// }
	// if _, err := firebaseApp.SendMessage(ctx, msg); err == nil {
	//     s.db.Model(notification).Update("is_push_sent", true)
	// }
}
