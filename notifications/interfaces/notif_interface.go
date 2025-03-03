package interfaces

import (
	"github.com/libileh/eegis/notifications/domain"
)

// EmailSender interface defines the contract for email sending
type EmailSender interface {
	SendEmail(to string, subject string, content string) error
}

// TemplateManager interface defines the contract for template management
type TemplateManager interface {
	RegisterTemplate(templateType string, template *domain.EmailTemplate) error
	GetTemplate(templateType domain.NotificationType) (*domain.EmailTemplate, error)
}

// ContentRenderer interface defines the contract for content rendering
type ContentRenderer interface {
	Render(template *domain.EmailTemplate, data domain.TemplateData) (*string, error)
}
