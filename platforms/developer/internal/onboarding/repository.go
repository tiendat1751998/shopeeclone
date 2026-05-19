package onboarding

import (
	"context"
	"sync"
)

type Repository interface {
	GetTemplate(ctx context.Context, name string) (*Template, error)
	ListTemplates(ctx context.Context) ([]*Template, error)
	StoreTemplate(ctx context.Context, t *Template) error
}

type InMemoryRepository struct {
	mu        sync.RWMutex
	templates map[string]*Template
}

func NewInMemoryRepository() *InMemoryRepository {
	r := &InMemoryRepository{
		templates: make(map[string]*Template),
	}
	r.seed()
	return r
}

func (r *InMemoryRepository) seed() {
	r.templates["microservice"] = &Template{
		Name:        "microservice",
		ServiceType: "microservice",
		Tasks: []OnboardingTask{
			{ID: "ms-1", Title: "Set up repository", Description: "Create a new repository for your service", Category: CategorySetup, Required: true, IsCompleted: false, Order: 1},
			{ID: "ms-2", Title: "Define API contracts", Description: "Define proto or OpenAPI specs", Category: CategoryLearn, Required: true, IsCompleted: false, Order: 2},
			{ID: "ms-3", Title: "Implement business logic", Description: "Write core service logic", Category: CategoryBuild, Required: true, IsCompleted: false, Order: 3},
			{ID: "ms-4", Title: "Deploy to staging", Description: "Deploy service to staging environment", Category: CategoryDeploy, Required: true, IsCompleted: false, Order: 4},
			{ID: "ms-5", Title: "Set up monitoring", Description: "Configure dashboards and alerts", Category: CategoryDeploy, Required: false, IsCompleted: false, Order: 5},
		},
	}
	r.templates["frontend"] = &Template{
		Name:        "frontend",
		ServiceType: "frontend",
		Tasks: []OnboardingTask{
			{ID: "fe-1", Title: "Initialize project", Description: "Scaffold frontend project", Category: CategorySetup, Required: true, IsCompleted: false, Order: 1},
			{ID: "fe-2", Title: "Set up CI/CD", Description: "Configure build and deploy pipelines", Category: CategoryBuild, Required: true, IsCompleted: false, Order: 2},
			{ID: "fe-3", Title: "Implement features", Description: "Build UI components", Category: CategoryBuild, Required: true, IsCompleted: false, Order: 3},
			{ID: "fe-4", Title: "Deploy to production", Description: "Release to production", Category: CategoryDeploy, Required: true, IsCompleted: false, Order: 4},
		},
	}
}

func (r *InMemoryRepository) GetTemplate(ctx context.Context, name string) (*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.templates[name]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (r *InMemoryRepository) ListTemplates(ctx context.Context) ([]*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Template
	for _, t := range r.templates {
		result = append(result, t)
	}
	return result, nil
}

func (r *InMemoryRepository) StoreTemplate(ctx context.Context, t *Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[t.Name] = t
	return nil
}
