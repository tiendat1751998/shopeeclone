package coordinator

import (
	"context"
	"sync"
)

type Repository interface {
	CreateNode(ctx context.Context, node *IndexNode) error
	GetNode(ctx context.Context, id string) (*IndexNode, error)
	ListNodes(ctx context.Context) ([]*IndexNode, error)
	UpdateNode(ctx context.Context, node *IndexNode) error
	DeleteNode(ctx context.Context, id string) error
	CreateShard(ctx context.Context, shard *IndexShard) error
	GetShard(ctx context.Context, id string) (*IndexShard, error)
	ListShards(ctx context.Context) ([]*IndexShard, error)
	ListShardsByNode(ctx context.Context, nodeID string) ([]*IndexShard, error)
	UpdateShard(ctx context.Context, shard *IndexShard) error
	DeleteShard(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu     sync.RWMutex
	nodes  map[string]*IndexNode
	shards map[string]*IndexShard
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		nodes:  make(map[string]*IndexNode),
		shards: make(map[string]*IndexShard),
	}
}

func (r *InMemoryRepository) CreateNode(_ context.Context, node *IndexNode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[node.ID] = node
	return nil
}

func (r *InMemoryRepository) GetNode(_ context.Context, id string) (*IndexNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	node, ok := r.nodes[id]
	if !ok {
		return nil, ErrNodeNotFound
	}
	return node, nil
}

func (r *InMemoryRepository) ListNodes(_ context.Context) ([]*IndexNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*IndexNode, 0, len(r.nodes))
	for _, n := range r.nodes {
		list = append(list, n)
	}
	return list, nil
}

func (r *InMemoryRepository) UpdateNode(_ context.Context, node *IndexNode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.nodes[node.ID]; !ok {
		return ErrNodeNotFound
	}
	r.nodes[node.ID] = node
	return nil
}

func (r *InMemoryRepository) DeleteNode(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.nodes[id]; !ok {
		return ErrNodeNotFound
	}
	delete(r.nodes, id)
	return nil
}

func (r *InMemoryRepository) CreateShard(_ context.Context, shard *IndexShard) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.shards[shard.ID]; ok {
		return ErrShardAlreadyExists
	}
	r.shards[shard.ID] = shard
	return nil
}

func (r *InMemoryRepository) GetShard(_ context.Context, id string) (*IndexShard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	shard, ok := r.shards[id]
	if !ok {
		return nil, ErrShardNotFound
	}
	return shard, nil
}

func (r *InMemoryRepository) ListShards(_ context.Context) ([]*IndexShard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*IndexShard, 0, len(r.shards))
	for _, s := range r.shards {
		list = append(list, s)
	}
	return list, nil
}

func (r *InMemoryRepository) ListShardsByNode(_ context.Context, nodeID string) ([]*IndexShard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*IndexShard, 0)
	for _, s := range r.shards {
		if s.NodeID == nodeID {
			list = append(list, s)
		}
	}
	return list, nil
}

func (r *InMemoryRepository) UpdateShard(_ context.Context, shard *IndexShard) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.shards[shard.ID]; !ok {
		return ErrShardNotFound
	}
	r.shards[shard.ID] = shard
	return nil
}

func (r *InMemoryRepository) DeleteShard(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.shards[id]; !ok {
		return ErrShardNotFound
	}
	delete(r.shards, id)
	return nil
}
