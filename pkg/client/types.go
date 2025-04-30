package client

// InputConfig represents the JSON input configuration
type InputConfig struct {
	ConfigURL string `json:"config_url"`
	PlanURL   string `json:"plan_url"`
	Token     string `json:"token"`
	ProductID string `json:"product-id"`
	Name      string `json:"name"`
}

// UploadRequest represents the multipart form data for the final POST request
type UploadRequest struct {
	TFPlan      string `json:"tf-plan"`
	ProductID   string `json:"product-id"`
	Name        string `json:"name"`
	MappingFile string `json:"mapping-file"`
	TFGraphFile string `json:"tfgraph-file"`
}
