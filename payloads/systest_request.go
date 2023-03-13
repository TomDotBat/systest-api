package payloads

type SysTestRequest struct {
	App        string            `json:"app" validate:"required"`
	InstanceId string            `json:"instanceId" validate:"required"`
	Method     string            `json:"method" validate:"required"`
	Path       string            `json:"path" validate:"required,startswith=/"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}
