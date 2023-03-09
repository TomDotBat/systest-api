package payloads

type ErrorResponse struct {
	Status  int    `json:"status" validate:"required,number"`
	Message string `json:"message,omitempty"`
}
