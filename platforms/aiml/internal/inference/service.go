package inference

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Predictor interface {
	Predict(ctx context.Context, req *InferenceRequest) (*InferenceResult, error)
	BatchPredict(ctx context.Context, reqs []*InferenceRequest) ([]*InferenceResult, error)
	GetModelInfo(ctx context.Context, name, version string) (*ModelInfo, error)
}

type MockPredictor struct {
	mu     sync.RWMutex
	models map[string]*ModelInfo
}

func NewMockPredictor() *MockPredictor {
	return &MockPredictor{
		models: make(map[string]*ModelInfo),
	}
}

func (p *MockPredictor) RegisterModel(info *ModelInfo) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.models[info.Name+":"+info.Version] = info
}

func (p *MockPredictor) Predict(ctx context.Context, req *InferenceRequest) (*InferenceResult, error) {
	if req.Input == nil {
		return nil, ErrInvalidInput
	}
	output := make(map[string]interface{})
	for k, v := range req.Input {
		switch val := v.(type) {
		case float64:
			output[k+"_prediction"] = val * 1.5
		case float32:
			output[k+"_prediction"] = float64(val) * 1.5
		case int:
			output[k+"_prediction"] = float64(val) * 1.5
		case int64:
			output[k+"_prediction"] = float64(val) * 1.5
		default:
			output[k+"_prediction"] = v
		}
	}
	version := req.ModelVersion
	if version == "" {
		version = "latest"
	}
	latency := math.Round((rand.Float64()*50+5)*100) / 100
	return &InferenceResult{
		Output:           output,
		Confidence:       math.Round((rand.Float64()*0.3+0.65)*100) / 100,
		LatencyMs:        latency,
		ModelVersionUsed: version,
	}, nil
}

func (p *MockPredictor) BatchPredict(ctx context.Context, reqs []*InferenceRequest) ([]*InferenceResult, error) {
	results := make([]*InferenceResult, 0, len(reqs))
	for _, req := range reqs {
		result, err := p.Predict(ctx, req)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (p *MockPredictor) GetModelInfo(ctx context.Context, name, version string) (*ModelInfo, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	info, ok := p.models[name+":"+version]
	if !ok {
		info, ok = p.models[name+":latest"]
		if !ok {
			return nil, ErrModelNotFound
		}
	}
	return info, nil
}

type Service struct {
	predictor Predictor
}

func NewService(predictor Predictor) *Service {
	return &Service{predictor: predictor}
}

func (s *Service) Predict(ctx context.Context, req *InferenceRequest) (*InferenceResult, error) {
	start := time.Now()
	result, err := s.predictor.Predict(ctx, req)
	if err != nil {
		return nil, err
	}
	result.LatencyMs = float64(time.Since(start).Microseconds()) / 1000.0
	return result, nil
}

func (s *Service) BatchPredict(ctx context.Context, reqs []*InferenceRequest) ([]*InferenceResult, error) {
	return s.predictor.BatchPredict(ctx, reqs)
}

func (s *Service) GetModelInfo(ctx context.Context, name, version string) (*ModelInfo, error) {
	return s.predictor.GetModelInfo(ctx, name, version)
}
