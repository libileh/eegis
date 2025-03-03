package domain

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

// Notification represents a notification message
type Notification struct {
	ID        string             `json:"id"`
	Recipient string             `json:"recipient"`
	Subject   string             `json:"subject"`
	Content   TemplateData       `json:"content"`
	Type      NotificationType   `json:"type"`
	Status    NotificationStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
}

// NotificationType represents the type of notification
type NotificationType string

const (
	UserVerification NotificationType = "user-verification"
	PasswordReset    NotificationType = "password-reset"
	ReviewPost       NotificationType = "review-post"
)

func NewNotificationType(notificationType string) (NotificationType, error) {
	validateTypes := map[string]NotificationType{
		string(UserVerification): UserVerification,
		string(PasswordReset):    PasswordReset,
		string(ReviewPost):       ReviewPost,
	}
	if validType, status := validateTypes[notificationType]; status {
		return validType, nil
	}
	return "", errors.New("invalid notification type")

}

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	Sent   NotificationStatus = "sent"
	Failed NotificationStatus = "failed"
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

type PostReviewEvent struct {
	PostID      uuid.UUID        `json:"post_id"`
	NewStatus   PostReviewStatus `json:"new_status"`
	AuthorEmail string           `json:"author_email"`
}

// PostReviewStatus represents the valid statuses for a Post.
type PostReviewStatus string

const (
	Pending  PostReviewStatus = "pending"
	Approved PostReviewStatus = "approved"
	Rejected PostReviewStatus = "rejected"
)
