package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/libileh/eegis/users/internal/config"
	"github.com/libileh/eegis/users/internal/infra/request"
	"net/http"
)

type HttpNotificationClient struct {
	Props      *config.UserProperties
	HTTPClient *http.Client
}

func NewHttpNotifService(props *config.UserProperties) *HttpNotificationClient {
	return &HttpNotificationClient{
		Props:      props,
		HTTPClient: &http.Client{},
	}
}

// SendConfirmationEmail sends a confirmation email to the specified recipient
func (cl *HttpNotificationClient) SendConfirmationEmail(email, token string) error {
	templateData := request.TemplateDataDTO{
		Name: email,
		ConfirmationLink: fmt.Sprintf("%s/confirm?token=%s",
			cl.Props.CommonProps.FrontendUrl, token),
	}

	notification := request.NotificationDTO{
		Recipient: email,
		Subject:   "Welcome to Qoraal Hub - Confirm Your Email",
		Content:   templateData,
		Type:      "verification",
		Status:    "pending",
	}

	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	// Construct the correct endpoint URL
	endpoint := fmt.Sprintf("%s/v1/notifications/user-confirmation",
		cl.Props.NotificationProps.NotificationBaseUrl)

	resp, err := cl.HTTPClient.Post(endpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create notification request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notification service returned status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	return nil
}
