package modelregistry

import "time"

type ModelType string

const (
	TypeRecommendation ModelType = "recommendation"
	TypeFraud          ModelType = "fraud"
	TypeRanking        ModelType = "ranking"
	TypeSearch         ModelType = "search"
)

type Framework string

const (
	FrameworkTensorFlow Framework = "tensorflow"
	FrameworkPyTorch    Framework = "pytorch"
	FrameworkONNX       Framework = "onnx"
)

type Stage string

const (
	StageDevelopment Stage = "development"
	StageStaging     Stage = "staging"
	StageProduction  Stage = "production"
	StageArchived    Stage = "archived"
)

type ModelMetrics struct {
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
}

type Model struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Type         ModelType    `json:"type"`
	Framework    Framework    `json:"framework"`
	Status       Stage        `json:"status"`
	Metrics      ModelMetrics `json:"metrics"`
	ArtifactPath string       `json:"artifact_path"`
	CreatedAt    time.Time    `json:"created_at"`
}
