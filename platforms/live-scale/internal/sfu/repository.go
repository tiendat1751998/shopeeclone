package sfu

import (
	"context"
	"sync"
)

type Repository interface {
	SaveNode(ctx context.Context, node *SFUNode) error
	GetNode(ctx context.Context, id string) (*SFUNode, error)
	DeleteNode(ctx context.Context, id string) error
	ListNodes(ctx context.Context) ([]*SFUNode, error)
	SaveStreamSession(ctx context.Context, session *StreamSession) error
	GetStreamSession(ctx context.Context, id string) (*StreamSession, error)
	DeleteStreamSession(ctx context.Context, id string) error
	ListStreamSessionsByNode(ctx context.Context, nodeID string) ([]*StreamSession, error)
	ListStreamSessions(ctx context.Context) ([]*StreamSession, error)
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	nodes    map[string]*SFUNode
	sessions map[string]*StreamSession
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		nodes:    make(map[string]*SFUNode),
		sessions: make(map[string]*StreamSession),
	}
}

func (r *InMemoryRepository) SaveNode(_ context.Context, node *SFUNode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[node.ID] = node
	return nil
}

func (r *InMemoryRepository) GetNode(_ context.Context, id string) (*SFUNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.nodes[id]
	if !ok {
		return nil, ErrNodeNotFound
	}
	return n, nil
}

func (r *InMemoryRepository) DeleteNode(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.nodes, id)
	return nil
}

func (r *InMemoryRepository) ListNodes(_ context.Context) ([]*SFUNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*SFUNode, 0, len(r.nodes))
	for _, n := range r.nodes {
		result = append(result, n)
	}
	return result, nil
}

func (r *InMemoryRepository) SaveStreamSession(_ context.Context, session *StreamSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}

func (r *InMemoryRepository) GetStreamSession(_ context.Context, id string) (*StreamSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sessions[id]
	if !ok {
		return nil, ErrStreamSessionNotFound
	}
	return s, nil
}

func (r *InMemoryRepository) DeleteStreamSession(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, id)
	return nil
}

func (r *InMemoryRepository) ListStreamSessionsByNode(_ context.Context, nodeID string) ([]*StreamSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*StreamSession
	for _, s := range r.sessions {
		if s.NodeID == nodeID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) ListStreamSessions(_ context.Context) ([]*StreamSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*StreamSession, 0, len(r.sessions))
	for _, s := range r.sessions {
		result = append(result, s)
	}
	return result, nil
}
