package pipeline

type ProcessorType string

const (
	ProcessorTokenize   ProcessorType = "tokenize"
	ProcessorAnalyze    ProcessorType = "analyze"
	ProcessorTransform  ProcessorType = "transform"
	ProcessorEnrich     ProcessorType = "enrich"
)

type PipelineStage struct {
	Name      string         `json:"name"`
	Processor ProcessorType  `json:"processor"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

type Pipeline struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	IndexName string           `json:"index_name"`
	Stages    []PipelineStage  `json:"stages"`
	IsActive  bool             `json:"is_active"`
}

type Document struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}
