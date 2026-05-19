package routes

import (
	"strings"
	"sync"
)

type Repository interface {
	Store(route *Route) error
	Get(id string) (*Route, error)
	List() ([]*Route, error)
	Delete(id string) error
}

type InMemoryRepository struct {
	mu     sync.RWMutex
	routes map[string]*Route
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		routes: make(map[string]*Route),
	}
}

func (r *InMemoryRepository) Store(route *Route) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes[route.ID] = route
	return nil
}

func (r *InMemoryRepository) Get(id string) (*Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	route, ok := r.routes[id]
	if !ok {
		return nil, nil
	}
	return route, nil
}

func (r *InMemoryRepository) List() ([]*Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Route, 0, len(r.routes))
	for _, route := range r.routes {
		result = append(result, route)
	}
	return result, nil
}

func (r *InMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.routes, id)
	return nil
}

func MatchRoute(routes []*Route, path, method string) *Route {
	var best *Route
	bestPrefixLen := -1

	for _, route := range routes {
		if !route.IsActive {
			continue
		}
		if !methodMatches(route.Methods, method) {
			continue
		}
		if exactMatch(route.Path, path) {
			if bestPrefixLen < len(route.Path) {
				best = route
				bestPrefixLen = len(route.Path)
			}
			continue
		}
		if prefixLen := prefixMatch(route.Path, path); prefixLen > 0 {
			if prefixLen > bestPrefixLen {
				best = route
				bestPrefixLen = prefixLen
			}
		}
	}
	return best
}

func methodMatches(methods []string, method string) bool {
	for _, m := range methods {
		if strings.EqualFold(m, method) {
			return true
		}
	}
	return false
}

func exactMatch(pattern, path string) bool {
	return pattern == path
}

func prefixMatch(pattern, path string) int {
	if !strings.HasSuffix(pattern, "*") {
		return 0
	}
	prefix := strings.TrimSuffix(pattern, "*")
	if strings.HasPrefix(path, prefix) {
		return len(prefix)
	}
	return 0
}
