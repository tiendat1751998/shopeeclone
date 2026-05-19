package multiregion

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Repository interface {
	Create(ctx context.Context, region *Region) error
	Get(ctx context.Context, code string) (*Region, error)
	Update(ctx context.Context, region *Region) error
	List(ctx context.Context) ([]*Region, error)
	GetActive(ctx context.Context) ([]*Region, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	regions map[string]*Region
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		regions: make(map[string]*Region),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, region *Region) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	region.CreatedAt = now
	region.UpdatedAt = now
	if region.Endpoints == nil {
		region.Endpoints = make(map[string]string)
	}
	r.regions[region.Code] = region
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, code string) (*Region, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	region, ok := r.regions[code]
	if !ok {
		return nil, nil
	}
	return region, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, region *Region) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.regions[region.Code]; !ok {
		return fmt.Errorf("region not found: %s", region.Code)
	}
	region.UpdatedAt = time.Now()
	r.regions[region.Code] = region
	return nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*Region, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Region, 0, len(r.regions))
	for _, reg := range r.regions {
		result = append(result, reg)
	}
	return result, nil
}

func (r *InMemoryRepository) GetActive(ctx context.Context) ([]*Region, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Region
	for _, reg := range r.regions {
		if reg.IsActive {
			result = append(result, reg)
		}
	}
	return result, nil
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, region *Region) (*Region, error) {
	if region.Code == "" {
		return nil, fmt.Errorf("code is required")
	}
	if region.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if region.Endpoints == nil {
		region.Endpoints = make(map[string]string)
	}
	if err := s.repo.Create(ctx, region); err != nil {
		return nil, err
	}
	return region, nil
}

func (s *Service) Get(ctx context.Context, code string) (*Region, error) {
	return s.repo.Get(ctx, code)
}

func (s *Service) List(ctx context.Context) ([]*Region, error) {
	return s.repo.List(ctx)
}

func (s *Service) GetActiveRegions(ctx context.Context) ([]*Region, error) {
	return s.repo.GetActive(ctx)
}

func (s *Service) GetRegionEndpoints(ctx context.Context, code string) (map[string]string, error) {
	region, err := s.repo.Get(ctx, code)
	if err != nil {
		return nil, err
	}
	if region == nil {
		return nil, fmt.Errorf("region not found: %s", code)
	}
	return region.Endpoints, nil
}

func (s *Service) GetFailoverStrategy(ctx context.Context, code string, serviceName string) (*FailoverResult, error) {
	region, err := s.repo.Get(ctx, code)
	if err != nil {
		return nil, err
	}
	if region == nil {
		return nil, fmt.Errorf("region not found: %s", code)
	}
	if !region.IsActive {
		if region.FailoverRegion == "" {
			return nil, fmt.Errorf("region %s is inactive and no failover configured", code)
		}
		failoverRegion, err := s.repo.Get(ctx, region.FailoverRegion)
		if err != nil {
			return nil, err
		}
		if failoverRegion == nil {
			return nil, fmt.Errorf("failover region %s not found", region.FailoverRegion)
		}
		endpoint := failoverRegion.Endpoints[serviceName]
		return &FailoverResult{
			PrimaryRegion:  code,
			FailoverRegion: region.FailoverRegion,
			Endpoint:       endpoint,
			IsFailover:     true,
		}, nil
	}
	endpoint := region.Endpoints[serviceName]
	return &FailoverResult{
		PrimaryRegion:   code,
		FailoverRegion:  "",
		Endpoint:        endpoint,
		IsFailover:      false,
	}, nil
}
