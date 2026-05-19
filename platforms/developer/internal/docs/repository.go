package docs

import (
	"context"
	"strings"
	"sync"
)

type Repository interface {
	Store(ctx context.Context, doc *DocPage) error
	GetByID(ctx context.Context, id string) (*DocPage, error)
	ListByService(ctx context.Context, service string) ([]*DocPage, error)
	ListByCategory(ctx context.Context, category string) ([]*DocPage, error)
	Search(ctx context.Context, query string) ([]*DocPage, error)
	List(ctx context.Context) ([]*DocPage, error)
	Update(ctx context.Context, doc *DocPage) error
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	docs map[string]*DocPage
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		docs: make(map[string]*DocPage),
	}
}

func (r *InMemoryRepository) Store(ctx context.Context, doc *DocPage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.docs[doc.ID] = doc
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*DocPage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	doc, ok := r.docs[id]
	if !ok {
		return nil, nil
	}
	return doc, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*DocPage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*DocPage
	for _, doc := range r.docs {
		result = append(result, doc)
	}
	return result, nil
}

func (r *InMemoryRepository) ListByService(ctx context.Context, service string) ([]*DocPage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*DocPage
	for _, doc := range r.docs {
		if doc.Service == service {
			result = append(result, doc)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) ListByCategory(ctx context.Context, category string) ([]*DocPage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*DocPage
	for _, doc := range r.docs {
		if doc.Category == category {
			result = append(result, doc)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) Search(ctx context.Context, query string) ([]*DocPage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	q := strings.ToLower(query)
	var result []*DocPage
	for _, doc := range r.docs {
		if strings.Contains(strings.ToLower(doc.Title), q) ||
			strings.Contains(strings.ToLower(doc.Content), q) ||
			strings.Contains(strings.ToLower(doc.Category), q) {
			result = append(result, doc)
		}
		for _, tag := range doc.Tags {
			if strings.Contains(strings.ToLower(tag), q) {
				result = append(result, doc)
				break
			}
		}
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, doc *DocPage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.docs[doc.ID] = doc
	return nil
}
