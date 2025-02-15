package sender

import (
	"bytes"
	"encoding/json"
	"github.com/libileh/eegis/notifications/interfaces"
	"io"

	"fmt"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/notifications/internal/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// EmailSenderImpl EmailSender implementation
type EmailSenderImpl struct {
	logger     *zap.SugaredLogger
	emailProps *config.EmailProperties
	client     *http.Client
	baseUrl    string
}

func NewEmailSender(logger *zap.SugaredLogger, props *config.EmailProperties) interfaces.EmailSender {
	return &EmailSenderImpl{
		logger:     logger,
		emailProps: props,
		client:     &http.Client{Timeout: time.Second * 10},
		baseUrl:    props.Mailtrap.Url,
	}
}

func (e *EmailSenderImpl) SendEmail(to string, subject string, content string) error {
	if to == "" {
		return errors.NewBadRequest("recipient email is required")
	}

	payload := map[string]interface{}{
		"from": map[string]interface{}{
			"email": e.emailProps.FromEmail,
			"name":  "Mailtrap Test",
		},
		"to": []map[string]interface{}{
			{
				"email": to,
			},
		},
		"subject":  subject,
		"html":     content,
		"category": "Integration Test",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.NewBadRequest("failed to marshal email payload %s", err)
	}

	req, err := http.NewRequest("POST", e.baseUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return errors.NewInternalServerError("failed to create email request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.emailProps.Mailtrap.ApiKey))

	resp, err := e.client.Do(req)
	if err != nil {
		return errors.NewInternalServerError("failed to send email request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return errors.NewInternalServerError(fmt.Sprintf("failed to send email, status code: %d, cause: %v", resp.StatusCode, bodyString))
	}
	return nil
}
