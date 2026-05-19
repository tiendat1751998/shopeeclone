package inference

type InferenceRequest struct {
	ModelName    string                 `json:"model_name"`
	ModelVersion string                 `json:"model_version"`
	Input        map[string]interface{} `json:"input"`
	Features     map[string]interface{} `json:"features,omitempty"`
	Context      map[string]interface{} `json:"context,omitempty"`
}

type InferenceResult struct {
	Output            map[string]interface{} `json:"output"`
	Confidence        float64                `json:"confidence"`
	LatencyMs         float64                `json:"latency_ms"`
	ModelVersionUsed  string                 `json:"model_version_used"`
}

type ModelInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	InputSchema string `json:"input_schema"`
	OutputSchema string `json:"output_schema"`
}
