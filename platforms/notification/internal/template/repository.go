package template

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, t *Template) error
	GetByID(ctx context.Context, id string) (*Template, error)
	GetByName(ctx context.Context, name string) (*Template, error)
	List(ctx context.Context) ([]*Template, error)
	Update(ctx context.Context, t *Template) error
	Delete(ctx context.Context, id string) error
	CreateVersion(ctx context.Context, v *TemplateVersion) error
	ListVersions(ctx context.Context, templateID string) ([]*TemplateVersion, error)
}

type InMemoryRepository struct {
	mu        sync.RWMutex
	templates map[string]*Template
	versions  map[string][]*TemplateVersion
	nameIndex map[string]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		templates: make(map[string]*Template),
		versions:  make(map[string][]*TemplateVersion),
		nameIndex: make(map[string]string),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, t *Template) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	t.Version = 1

	r.mu.Lock()
	defer r.mu.Unlock()

	r.templates[t.ID] = t
	r.nameIndex[t.Name] = t.ID

	version := &TemplateVersion{
		ID:         uuid.New().String(),
		TemplateID: t.ID,
		Version:    1,
		Subject:    t.Subject,
		Body:       t.Body,
		CreatedAt:  now,
	}
	r.versions[t.ID] = append(r.versions[t.ID], version)

	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.templates[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (r *InMemoryRepository) GetByName(ctx context.Context, name string) (*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.nameIndex[name]
	if !ok {
		return nil, nil
	}

	t, ok := r.templates[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Template
	for _, t := range r.templates {
		result = append(result, t)
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, t *Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.templates[t.ID]
	if !ok {
		return nil
	}

	t.Version = existing.Version + 1
	t.CreatedAt = existing.CreatedAt
	t.UpdatedAt = time.Now()

	r.templates[t.ID] = t

	version := &TemplateVersion{
		ID:         uuid.New().String(),
		TemplateID: t.ID,
		Version:    t.Version,
		Subject:    t.Subject,
		Body:       t.Body,
		CreatedAt:  t.UpdatedAt,
	}
	r.versions[t.ID] = append(r.versions[t.ID], version)

	return nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.templates[id]
	if ok {
		delete(r.nameIndex, t.Name)
	}
	delete(r.templates, id)
	delete(r.versions, id)
	return nil
}

func (r *InMemoryRepository) CreateVersion(ctx context.Context, v *TemplateVersion) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.versions[v.TemplateID] = append(r.versions[v.TemplateID], v)
	return nil
}

func (r *InMemoryRepository) ListVersions(ctx context.Context, templateID string) ([]*TemplateVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	versions, ok := r.versions[templateID]
	if !ok {
		return []*TemplateVersion{}, nil
	}
	return versions, nil
}
