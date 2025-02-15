package domain

import (
	"errors"
	"strings"
	"time"
)

// Notification represents a notification message
type Notification struct {
	ID        string       `json:"id"`
	Recipient string       `json:"recipient"`
	Subject   string       `json:"subject"`
	Content   TemplateData `json:"content"`
	Type      string       `json:"type"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
}

// NotificationType represents the type of notification
type NotificationType string

const (
	Verification  NotificationType = "verification"
	PasswordReset NotificationType = "password-reset"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	Pending NotificationStatus = "pending"
	Sent    NotificationStatus = "sent"
	Failed  NotificationStatus = "failed"
)

// TemplateData represents the data needed for template rendering
type TemplateData map[string]string

// ReplaceString replaces all occurrences of placeholder with value in text
// Returns an error if the placeholder is not found in the text
func ReplaceString(text string, placeholder string, value string) (string, error) {
	replaced := strings.ReplaceAll(text, placeholder, value)
	if replaced == text {
		return replaced, errors.New("placeholder not found in text")
	}
	return replaced, nil
}
