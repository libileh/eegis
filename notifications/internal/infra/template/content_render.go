package template

import (
	"github.com/libileh/eegis/notifications/domain"
	"github.com/libileh/eegis/notifications/interfaces"
	"go.uber.org/zap"
)

// ContentRenderer implementation
type ContentRendererImpl struct {
	logger *zap.SugaredLogger
}

func NewContentRenderer(logger *zap.SugaredLogger) interfaces.ContentRenderer {
	return &ContentRendererImpl{
		logger: logger,
	}
}

func (c *ContentRendererImpl) Render(template *domain.EmailTemplate, data domain.TemplateData) (*string, error) {
	// Simple template rendering implementation

	content := template.Content
	var err error

	// Replace placeholders with actual data
	for key, value := range data {
		content, err = replacePlaceholders(content, key, value)
		if err != nil {
			return nil, err
		}
	}
	return &content, nil
}

func replacePlaceholders(text string, key string, value string) (string, error) {
	placeholder := "{{." + key + "}}"
	return domain.ReplaceString(text, placeholder, value)
}
