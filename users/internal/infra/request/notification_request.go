package request

type NotificationDTO struct {
	Recipient string          `json:"recipient"`
	Subject   string          `json:"subject"`
	Content   TemplateDataDTO `json:"content"`
	Type      string          `json:"type"`
	Status    string          `json:"status"`
}

type TemplateDataDTO struct {
	Name             string `json:"name"`
	ConfirmationLink string `json:"confirmationLink"`
}
