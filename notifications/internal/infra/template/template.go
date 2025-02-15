package template

import (
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/notifications/domain"
	"go.uber.org/zap"
)

// TemplateManager implementation
type TemplateManagerImpl struct {
	logger    *zap.SugaredLogger
	templates map[string]*domain.EmailTemplate
}

func NewTemplateManager(logger *zap.SugaredLogger) *TemplateManagerImpl {
	tm := &TemplateManagerImpl{
		logger:    logger,
		templates: make(map[string]*domain.EmailTemplate),
	}
	// Register the verification template
	err := tm.RegisterTemplate("verification", &domain.EmailTemplate{
		Subject: "Welcome to Qoraal Hub - Confirm Your Email",
		Content: "Hello {{.name}},<br><br>Thank you for registering with Qoraal Hub!<br><br>Please click the following link to confirm your email:<br>{{.confirmationLink}}<br><br>If you didn't create this account, you can safely ignore this email.<br><br>Best regards,<br>The Qoraal Hub Team.",
	})
	if err != nil {
		return nil
	}

	return tm
}

func (t *TemplateManagerImpl) RegisterTemplate(templateType string, template *domain.EmailTemplate) error {
	t.templates[templateType] = template
	return nil
}

func (t *TemplateManagerImpl) GetTemplate(templateType string) (*domain.EmailTemplate, error) {
	template, exists := t.templates[templateType]
	if !exists {
		return nil, errors.NewNotFoundError(templateType)
	}
	return template, nil
}
