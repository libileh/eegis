package template

import (
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/notifications/domain"
	"go.uber.org/zap"
)

type TemplateManagerImpl struct {
	logger    *zap.SugaredLogger
	templates map[string]*domain.EmailTemplate
}

// NewTemplateManagerImpl creates a new instance of TemplateManagerImpl
func NewTemplateManager(logger *zap.SugaredLogger) *TemplateManagerImpl {
	tm := &TemplateManagerImpl{
		logger:    logger,
		templates: make(map[string]*domain.EmailTemplate),
	}

	// User Register the verification template
	if err := tm.RegisterTemplate("user-verification", &domain.EmailTemplate{
		Subject: "Welcome to Qoraal Hub - Confirm Your Email",
		Content: "Hello {{.name}},<br><br>Thank you for registering with Qoraal Hub!<br><br>Please click the following link to confirm your email:<br>{{.confirmationLink}}<br><br>If you didn't create this account, you can safely ignore this email.<br><br>Best regards,<br>The Qoraal Hub Team.",
	}); err != nil {
		return nil
	}

	// Register the review post notification template.
	if err := tm.RegisterTemplate("review-post", &domain.EmailTemplate{
		Subject: "New Review Posted - Check It Out",
		Content: "Hello {{.name}},<br><br>A new review has been posted.<br>Please visit the following link to read it:<br>{{.reviewLink}}<br><br>Thank you,<br>The Qoraal Hub Team.",
	}); err != nil {
		return nil
	}

	return tm
}

func (t *TemplateManagerImpl) RegisterTemplate(templateType string, template *domain.EmailTemplate) error {
	t.templates[templateType] = template
	return nil
}

func (t *TemplateManagerImpl) GetTemplate(templateType domain.NotificationType) (*domain.EmailTemplate, error) {
	template, exists := t.templates[string(templateType)]
	if !exists {
		return nil, errors.NewNotFoundError(string(templateType))
	}
	return template, nil
}
