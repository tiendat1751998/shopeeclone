package featurestore

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterFeature(ctx context.Context, feature *Feature) error {
	return s.repo.StoreFeature(ctx, feature)
}

func (s *Service) GetFeature(ctx context.Context, name string) (*Feature, error) {
	return s.repo.GetFeature(ctx, name)
}

func (s *Service) ListFeatures(ctx context.Context) ([]*Feature, error) {
	return s.repo.ListFeatures(ctx)
}

func (s *Service) SetFeatureValue(ctx context.Context, value *FeatureValue) error {
	return s.repo.StoreFeatureValue(ctx, value)
}

func (s *Service) GetFeatureValue(ctx context.Context, featureName, entityID string) (*FeatureValue, error) {
	return s.repo.GetFeatureValue(ctx, featureName, entityID)
}

func (s *Service) BatchGet(ctx context.Context, featureNames []string, entityID string) (map[string]*FeatureValue, error) {
	return s.repo.BatchGetFeatureValues(ctx, featureNames, entityID)
}

func (s *Service) GetFeaturesForEntity(ctx context.Context, entityID string, entityType EntityType) (map[string]*FeatureValue, error) {
	return s.repo.GetFeaturesForEntity(ctx, entityID, entityType)
}
