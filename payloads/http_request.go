package payloads

type HttpRequest struct {
	App        string            `json:"app"`
	InstanceId string            `json:"instanceId"`
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Time       int64             `json:"time,omitempty"`
	Parent     *HttpRequest      `json:"-"`
	Children   []*HttpRequest    `json:"children,omitempty"`
	Response   *HttpResponse     `json:"response,omitempty"`
}
