package requests

// swagger:model PostPayload
type PostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required"`
	Tags    []string `json:"tags"`
	Status  string   `json:"status" validate:"required" default:"pending"`
}

type ReviewPostPayload struct {
	Status string `json:"status"`
}
