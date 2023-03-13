package payloads

type HttpResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Time    int64             `json:"time"`
}
