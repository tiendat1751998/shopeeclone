package content

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateTemplate(ctx context.Context, t *ContentTemplate) error
	GetTemplateByID(ctx context.Context, id string) (*ContentTemplate, error)
	ListTemplates(ctx context.Context) ([]*ContentTemplate, error)
	UpdateTemplate(ctx context.Context, t *ContentTemplate) error
	DeleteTemplate(ctx context.Context, id string) error
	CreateVariant(ctx context.Context, v *Variant) error
	ListVariants(ctx context.Context, templateID string) ([]*Variant, error)
	GetVariantByID(ctx context.Context, id string) (*Variant, error)
}

type InMemoryRepository struct {
	mu        sync.RWMutex
	templates map[string]*ContentTemplate
	variants  map[string][]*Variant
	variantByID map[string]*Variant
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		templates:   make(map[string]*ContentTemplate),
		variants:    make(map[string][]*Variant),
		variantByID: make(map[string]*Variant),
	}
}

func (r *InMemoryRepository) CreateTemplate(ctx context.Context, t *ContentTemplate) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[t.ID] = t
	return nil
}

func (r *InMemoryRepository) GetTemplateByID(ctx context.Context, id string) (*ContentTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.templates[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (r *InMemoryRepository) ListTemplates(ctx context.Context) ([]*ContentTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*ContentTemplate
	for _, t := range r.templates {
		result = append(result, t)
	}
	return result, nil
}

func (r *InMemoryRepository) UpdateTemplate(ctx context.Context, t *ContentTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.templates[t.ID]
	if !ok {
		return nil
	}
	t.CreatedAt = existing.CreatedAt
	t.UpdatedAt = time.Now()
	r.templates[t.ID] = t
	return nil
}

func (r *InMemoryRepository) DeleteTemplate(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.templates, id)
	return nil
}

func (r *InMemoryRepository) CreateVariant(ctx context.Context, v *Variant) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	v.CreatedAt = time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.variants[v.TemplateID] = append(r.variants[v.TemplateID], v)
	r.variantByID[v.ID] = v
	return nil
}

func (r *InMemoryRepository) ListVariants(ctx context.Context, templateID string) ([]*Variant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	variants, ok := r.variants[templateID]
	if !ok {
		return []*Variant{}, nil
	}
	return variants, nil
}

func (r *InMemoryRepository) GetVariantByID(ctx context.Context, id string) (*Variant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.variantByID[id]
	if !ok {
		return nil, nil
	}
	return v, nil
}
