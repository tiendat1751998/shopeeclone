package experiments

import (
	"context"
	"hash/fnv"
	"math"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateExperiment(ctx context.Context, exp *Experiment) error {
	exp.Variants = []Variant{
		{ModelName: exp.ModelA, TrafficAllocationPct: exp.TrafficPct},
		{ModelName: exp.ModelB, TrafficAllocationPct: 100 - exp.TrafficPct},
	}
	exp.Status = StatusRunning
	return s.repo.StoreExperiment(ctx, exp)
}

func (s *Service) GetExperiment(ctx context.Context, id string) (*Experiment, error) {
	return s.repo.GetExperiment(ctx, id)
}

func (s *Service) ListExperiments(ctx context.Context) ([]*Experiment, error) {
	return s.repo.ListExperiments(ctx)
}

func consistentHash(userID string, modulus uint32) uint32 {
	h := fnv.New32a()
	h.Write([]byte(userID))
	return h.Sum32() % modulus
}

func (s *Service) AssignVariant(ctx context.Context, experimentID, userID string) (*Assignment, error) {
	exp, err := s.repo.GetExperiment(ctx, experimentID)
	if err != nil {
		return nil, err
	}
	if exp.Status != StatusRunning {
		return nil, ErrExperimentClosed
	}
	hash := consistentHash(userID, 100)
	var variant string
	if float64(hash) < exp.TrafficPct {
		variant = exp.ModelA
	} else {
		variant = exp.ModelB
	}
	return &Assignment{
		ExperimentID: experimentID,
		UserID:       userID,
		Variant:      variant,
	}, nil
}

func (s *Service) RecordResult(ctx context.Context, experimentID, variant, userID string, value float64) error {
	exp, err := s.repo.GetExperiment(ctx, experimentID)
	if err != nil {
		return err
	}
	if exp.Status != StatusRunning {
		return ErrExperimentClosed
	}
	return s.repo.StoreResult(ctx, &Result{
		ExperimentID: experimentID,
		Variant:      variant,
		UserID:       userID,
		Value:        value,
	})
}

func (s *Service) GetExperimentResults(ctx context.Context, experimentID string) (*ExperimentResults, error) {
	exp, err := s.repo.GetExperiment(ctx, experimentID)
	if err != nil {
		return nil, err
	}
	results, err := s.repo.GetResults(ctx, experimentID)
	if err != nil {
		return nil, err
	}
	var aValues, bValues []float64
	for _, r := range results {
		if r.Variant == exp.ModelA {
			aValues = append(aValues, r.Value)
		} else {
			bValues = append(bValues, r.Value)
		}
	}
	aMean, bMean := 0.0, 0.0
	if len(aValues) > 0 {
		var sum float64
		for _, v := range aValues {
			sum += v
		}
		aMean = sum / float64(len(aValues))
	}
	if len(bValues) > 0 {
		var sum float64
		for _, v := range bValues {
			sum += v
		}
		bMean = sum / float64(len(bValues))
	}
	improvement := 0.0
	if aMean > 0 {
		improvement = math.Round(((bMean-aMean)/aMean)*10000) / 100
	}
	return &ExperimentResults{
		ExperimentID: experimentID,
		Name:         exp.Name,
		Metric:       exp.Metric,
		ModelA:       exp.ModelA,
		ModelB:       exp.ModelB,
		AResults: VariantResults{
			ModelName:  exp.ModelA,
			SampleSize: len(aValues),
			Sum:        sumSlice(aValues),
			Mean:       aMean,
		},
		BResults: VariantResults{
			ModelName:  exp.ModelB,
			SampleSize: len(bValues),
			Sum:        sumSlice(bValues),
			Mean:       bMean,
		},
		Improvement: improvement,
	}, nil
}

func sumSlice(vals []float64) float64 {
	var s float64
	for _, v := range vals {
		s += v
	}
	return math.Round(s*100) / 100
}
