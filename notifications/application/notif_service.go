package application

import (
	"github.com/libileh/eegis/notifications/domain"
	"github.com/libileh/eegis/notifications/interfaces"
	"go.uber.org/zap"
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

// SendEmail sends a notification email
func (s *NotificationService) SendEmail(notificationType string, recipient string, data domain.TemplateData) error {
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
