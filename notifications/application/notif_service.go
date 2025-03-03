package application

import (
	"fmt"
	"github.com/libileh/eegis/common/eventbus"
	"github.com/libileh/eegis/notifications/domain"
	"github.com/libileh/eegis/notifications/interfaces"
	"go.uber.org/zap"
	"time"
)

type ServiceManager struct {
	NotificationService *NotificationService
}

func NewServiceManager(notificationService *NotificationService) *ServiceManager {
	return &ServiceManager{
		NotificationService: notificationService,
	}
}

// NotificationService implements the core notification functionality
type NotificationService struct {
	emailSender     interfaces.EmailSender
	templateManager interfaces.TemplateManager
	contentRenderer interfaces.ContentRenderer
	logger          *zap.SugaredLogger
	Event           *eventbus.EventBus
}

// NewNotificationService creates a new instance of NotificationService
func NewNotificationService(
	emailSender interfaces.EmailSender,
	templateManager interfaces.TemplateManager,
	contentRenderer interfaces.ContentRenderer,
	logger *zap.SugaredLogger,
) *NotificationService {
	return &NotificationService{
		emailSender:     emailSender,
		templateManager: templateManager,
		contentRenderer: contentRenderer,
		logger:          logger,
	}
}

// SendUserVerificationEmail is a helper function for sending user verification emails.
// It wraps SendEmail using the 'user-verification' notification type.
func (s *NotificationService) SendUserVerificationEmail(recipient string, data domain.TemplateData) error {
	return s.SendEmail(domain.UserVerification, recipient, data)
}

// StartReviewPostNotificationTrigger is a convenience method that encapsulates
// the binding of the "review-post" event type with its associated handler.
func (s *NotificationService) StartReviewPostNotificationTrigger() error {
	return s.StartNotificationTrigger("review-post", s.HandlePostStatusChange)
}

// StartNotificationTrigger subscribes to events of a specified type
// and processes them using the provided handler function.
func (s *NotificationService) StartNotificationTrigger(eventType string, handler func(eventbus.Event) error) error {

	// channel to receive domain events
	eventsChan := make(chan eventbus.Event, 10) // using a buffered channel for resilience

	// Subscribe to events of the specified type
	if err := s.Event.Subscribe(eventType, eventsChan); err != nil {
		s.logger.Errorf("Failed to subscribe to %s events: %v", eventType, err)
		return err
	}

	s.logger.Infof("Notification trigger subscribed to events of type: %s", eventType)

	// Process events asynchronously
	go func() {
		for event := range eventsChan {
			if err := handler(event); err != nil {
				s.logger.Errorf("Error handling event: %v", err)
			}
		}
	}()
	return nil
}

// HandlePostStatusChange processes a post status change event and sends appropriate notifications.
// Updated to return error to be compatible with StartNotificationTrigger.
func (s *NotificationService) HandlePostStatusChange(event eventbus.Event) error {
	// Assume event.Data is of type *domain.PostReviewEvent.
	postStatusEvent, ok := event.Data.(*domain.PostReviewEvent)
	if !ok {
		err := fmt.Errorf("received event with unexpected data format: %+v", event.Data)
		s.logger.Errorf(err.Error())
		return err
	}

	// Filter only approved or rejected posts.
	if postStatusEvent.NewStatus != domain.Approved && postStatusEvent.NewStatus != domain.Rejected {
		// No notification required. Not an error.
		return nil
	}

	// Determine the recipient's email.
	recipientEmail := postStatusEvent.AuthorEmail

	// Prepare the notification message content.
	notificationData := domain.TemplateData{
		"post_id":    postStatusEvent.PostID.String(),
		"new_status": string(postStatusEvent.NewStatus),
		"occurred":   event.OccurredAt.Format(time.RFC1123),
	}

	// Select a notification type based on the event.
	notificationType, err := domain.NewNotificationType(event.Type)
	if err != nil {
		s.logger.Errorf("failed to determine notification type: %v", err)
		return err
	}

	// Send the notification (email, in-app, or push based on your interfaces).
	if err := s.SendEmail(notificationType, recipientEmail, notificationData); err != nil {
		s.logger.Errorf("failed to send notification for post %v: %v", postStatusEvent.PostID, err)
		return err
	}

	s.logger.Infof("Notification sent for post %v with new status: %s", postStatusEvent.PostID, postStatusEvent.NewStatus)
	return nil
}

// SendEmail sends a notification email based on the specified notification type.
func (s *NotificationService) SendEmail(notificationType domain.NotificationType, recipient string, data domain.TemplateData) error {
	template, err := s.templateManager.GetTemplate(notificationType)
	if err != nil {
		return err
	}

	content, err := s.contentRenderer.Render(template, data)
	if err != nil {
		return err
	}

	if err := s.emailSender.SendEmail(recipient, template.Subject, *content); err != nil {
		return err
	}

	s.logger.Infow("Email sent successfully", "type", notificationType, "recipient", recipient)
	return nil
}
