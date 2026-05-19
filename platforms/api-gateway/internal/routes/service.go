package routes

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(req *RegisterRouteRequest) (*Route, error) {
	if req.Path == "" {
		return nil, fmt.Errorf("path is required")
	}
	if len(req.Methods) == 0 {
		return nil, fmt.Errorf("methods are required")
	}
	if req.ServiceName == "" {
		return nil, fmt.Errorf("service_name is required")
	}
	if req.TargetURL == "" {
		return nil, fmt.Errorf("target_url is required")
	}

	id := uuid.New().String()
	now := time.Now()

	route := &Route{
		ID:           id,
		Path:         req.Path,
		Methods:      req.Methods,
		ServiceName:  req.ServiceName,
		TargetURL:    req.TargetURL,
		TimeoutMs:    req.TimeoutMs,
		RateLimit:    req.RateLimit,
		AuthRequired: req.AuthRequired,
		Middleware:   req.Middleware,
		IsActive:     true,
		Version:      req.Version,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if route.TimeoutMs == 0 {
		route.TimeoutMs = 30000
	}
	if route.RateLimit == 0 {
		route.RateLimit = 100
	}
	if route.Version == "" {
		route.Version = "1.0.0"
	}

	if err := s.repo.Store(route); err != nil {
		return nil, err
	}
	return route, nil
}

func (s *Service) Deregister(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	return s.repo.Delete(id)
}

func (s *Service) Get(id string) (*Route, error) {
	return s.repo.Get(id)
}

func (s *Service) List() ([]*Route, error) {
	return s.repo.List()
}

func (s *Service) Match(path, method string) (*Route, error) {
	routes, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	matched := MatchRoute(routes, path, method)
	return matched, nil
}
