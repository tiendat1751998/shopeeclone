package featurestore

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	StoreFeature(ctx context.Context, feature *Feature) error
	GetFeature(ctx context.Context, name string) (*Feature, error)
	ListFeatures(ctx context.Context) ([]*Feature, error)
	StoreFeatureValue(ctx context.Context, value *FeatureValue) error
	GetFeatureValue(ctx context.Context, featureName, entityID string) (*FeatureValue, error)
	BatchGetFeatureValues(ctx context.Context, featureNames []string, entityID string) (map[string]*FeatureValue, error)
	GetFeaturesForEntity(ctx context.Context, entityID string, entityType EntityType) (map[string]*FeatureValue, error)
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	features map[string]*Feature
	values   map[string]*FeatureValue
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		features: make(map[string]*Feature),
		values:   make(map[string]*FeatureValue),
	}
}

func featureValueKey(featureName, entityID string) string {
	return featureName + ":" + entityID
}

func (r *InMemoryRepository) StoreFeature(ctx context.Context, feature *Feature) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.features[feature.Name]; ok {
		return ErrFeatureAlreadyExists
	}
	if feature.CreatedAt.IsZero() {
		feature.CreatedAt = time.Now()
	}
	r.features[feature.Name] = feature
	return nil
}

func (r *InMemoryRepository) GetFeature(ctx context.Context, name string) (*Feature, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.features[name]
	if !ok {
		return nil, ErrFeatureNotFound
	}
	return f, nil
}

func (r *InMemoryRepository) ListFeatures(ctx context.Context) ([]*Feature, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Feature, 0, len(r.features))
	for _, f := range r.features {
		result = append(result, f)
	}
	return result, nil
}

func (r *InMemoryRepository) StoreFeatureValue(ctx context.Context, value *FeatureValue) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := featureValueKey(value.FeatureName, value.EntityID)
	if value.Timestamp.IsZero() {
		value.Timestamp = time.Now()
	}
	existing, ok := r.values[key]
	if !ok {
		value.Version = 1
	} else {
		value.Version = existing.Version + 1
	}
	r.values[key] = value
	return nil
}

func (r *InMemoryRepository) GetFeatureValue(ctx context.Context, featureName, entityID string) (*FeatureValue, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.values[featureValueKey(featureName, entityID)]
	if !ok {
		return nil, ErrFeatureValueNotFound
	}
	return v, nil
}

func (r *InMemoryRepository) BatchGetFeatureValues(ctx context.Context, featureNames []string, entityID string) (map[string]*FeatureValue, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]*FeatureValue)
	for _, name := range featureNames {
		if v, ok := r.values[featureValueKey(name, entityID)]; ok {
			result[name] = v
		}
	}
	return result, nil
}

func (r *InMemoryRepository) GetFeaturesForEntity(ctx context.Context, entityID string, entityType EntityType) (map[string]*FeatureValue, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]*FeatureValue)
	for _, f := range r.features {
		if f.Entity == entityType {
			if v, ok := r.values[featureValueKey(f.Name, entityID)]; ok {
				result[f.Name] = v
			}
		}
	}
	return result, nil
}
