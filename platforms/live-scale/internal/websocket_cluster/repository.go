package websocket_cluster

import (
	"context"
	"sync"
)

type Repository interface {
	SaveNode(ctx context.Context, node *WSNode) error
	GetNode(ctx context.Context, id string) (*WSNode, error)
	DeleteNode(ctx context.Context, id string) error
	ListNodes(ctx context.Context) ([]*WSNode, error)
	AssignRoom(ctx context.Context, roomID, nodeID string) error
	GetRoomNode(ctx context.Context, roomID string) (string, error)
	UnassignRoom(ctx context.Context, roomID string) error
	ListRooms(ctx context.Context) (map[string]string, error)
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	nodes map[string]*WSNode
	rooms map[string]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		nodes: make(map[string]*WSNode),
		rooms: make(map[string]string),
	}
}

func (r *InMemoryRepository) SaveNode(_ context.Context, node *WSNode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[node.ID] = node
	return nil
}

func (r *InMemoryRepository) GetNode(_ context.Context, id string) (*WSNode, error) {
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
	for room, nodeID := range r.rooms {
		if nodeID == id {
			delete(r.rooms, room)
		}
	}
	return nil
}

func (r *InMemoryRepository) ListNodes(_ context.Context) ([]*WSNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*WSNode, 0, len(r.nodes))
	for _, n := range r.nodes {
		result = append(result, n)
	}
	return result, nil
}

func (r *InMemoryRepository) AssignRoom(_ context.Context, roomID, nodeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rooms[roomID] = nodeID
	return nil
}

func (r *InMemoryRepository) GetRoomNode(_ context.Context, roomID string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	nodeID, ok := r.rooms[roomID]
	if !ok {
		return "", ErrRoomNotFound
	}
	return nodeID, nil
}

func (r *InMemoryRepository) UnassignRoom(_ context.Context, roomID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rooms, roomID)
	return nil
}

func (r *InMemoryRepository) ListRooms(_ context.Context) (map[string]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]string, len(r.rooms))
	for k, v := range r.rooms {
		result[k] = v
	}
	return result, nil
}
