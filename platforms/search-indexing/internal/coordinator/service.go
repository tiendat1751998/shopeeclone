package coordinator

import (
	"context"
	"math"

	"github.com/google/uuid"
)

type Service interface {
	RegisterNode(ctx context.Context, node *IndexNode) (*IndexNode, error)
	GetNode(ctx context.Context, id string) (*IndexNode, error)
	ListNodes(ctx context.Context) ([]*IndexNode, error)
	RemoveNode(ctx context.Context, id string) error
	AssignShard(ctx context.Context, shard *IndexShard) (*IndexShard, error)
	RebalanceShards(ctx context.Context) error
	GetShardDistribution(ctx context.Context) ([]*ShardDistribution, error)
	HandleNodeFailure(ctx context.Context, nodeID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RegisterNode(_ context.Context, node *IndexNode) (*IndexNode, error) {
	if node.ID == "" {
		node.ID = uuid.New().String()
	}
	node.IsActive = true
	return node, s.repo.CreateNode(nil, node)
}

func (s *service) GetNode(ctx context.Context, id string) (*IndexNode, error) {
	return s.repo.GetNode(ctx, id)
}

func (s *service) ListNodes(ctx context.Context) ([]*IndexNode, error) {
	return s.repo.ListNodes(ctx)
}

func (s *service) RemoveNode(ctx context.Context, id string) error {
	return s.repo.DeleteNode(ctx, id)
}

func (s *service) AssignShard(ctx context.Context, shard *IndexShard) (*IndexShard, error) {
	nodes, err := s.repo.ListNodes(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, ErrNoAvailableNodes
	}

	var best *IndexNode
	minLoad := math.MaxFloat64
	for _, n := range nodes {
		if !n.IsActive {
			continue
		}
		if n.LoadPercentage < minLoad {
			best = n
			minLoad = n.LoadPercentage
		}
	}
	if best == nil {
		return nil, ErrNoAvailableNodes
	}

	if shard.ID == "" {
		shard.ID = uuid.New().String()
	}
	shard.NodeID = best.ID
	shard.Status = ShardActive

	if err := s.repo.CreateShard(ctx, shard); err != nil {
		return nil, err
	}

	best.LoadPercentage = math.Min(100, best.LoadPercentage+10)
	best.AvailableShards--
	if err := s.repo.UpdateNode(ctx, best); err != nil {
		return nil, err
	}

	return shard, nil
}

func (s *service) RebalanceShards(ctx context.Context) error {
	shards, err := s.repo.ListShards(ctx)
	if err != nil {
		return err
	}
	if len(shards) == 0 {
		return nil
	}

	nodes, err := s.repo.ListNodes(ctx)
	if err != nil {
		return err
	}
	activeNodes := make([]*IndexNode, 0, len(nodes))
	for _, n := range nodes {
		if n.IsActive {
			activeNodes = append(activeNodes, n)
		}
	}
	if len(activeNodes) == 0 {
		return ErrNoAvailableNodes
	}

	activeShards := make([]*IndexShard, 0, len(shards))
	for _, sh := range shards {
		if sh.Status == ShardActive || sh.Status == ShardRelocating {
			activeShards = append(activeShards, sh)
		}
	}

	targetPerNode := len(activeShards) / len(activeNodes)
	if targetPerNode == 0 {
		targetPerNode = 1
	}

	shardIndex := 0
	for _, node := range activeNodes {
		nodeShards := make([]*IndexShard, 0)
		for _, sh := range activeShards {
			if sh.NodeID == node.ID {
				nodeShards = append(nodeShards, sh)
			}
		}

		if len(nodeShards) > targetPerNode {
			excess := len(nodeShards) - targetPerNode
			for i := 0; i < excess && shardIndex < len(activeShards); i++ {
				sh := nodeShards[i]
				sh.Status = ShardRelocating
				s.repo.UpdateShard(ctx, sh)
			}
		}
	}

	for _, sh := range activeShards {
		if sh.Status != ShardRelocating {
			continue
		}
		for _, node := range activeNodes {
			nodeShards := make([]*IndexShard, 0)
			for _, ns := range activeShards {
				if ns.NodeID == node.ID && ns.Status == ShardActive {
					nodeShards = append(nodeShards, ns)
				}
			}
			if len(nodeShards) < targetPerNode {
				sh.NodeID = node.ID
				sh.Status = ShardActive
				s.repo.UpdateShard(ctx, sh)
				break
			}
		}
	}

	return nil
}

func (s *service) GetShardDistribution(ctx context.Context) ([]*ShardDistribution, error) {
	nodes, err := s.repo.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	shards, err := s.repo.ListShards(ctx)
	if err != nil {
		return nil, err
	}

	shardsByNode := make(map[string][]*IndexShard)
	for _, sh := range shards {
		shardsByNode[sh.NodeID] = append(shardsByNode[sh.NodeID], sh)
	}

	dist := make([]*ShardDistribution, 0, len(nodes))
	for _, node := range nodes {
		nodeShards := shardsByNode[node.ID]
		if nodeShards == nil {
			nodeShards = []*IndexShard{}
		}
		dist = append(dist, &ShardDistribution{
			NodeID: node.ID,
			Node:   node,
			Shards: nodeShards,
			Count:  len(nodeShards),
		})
	}
	return dist, nil
}

func (s *service) HandleNodeFailure(ctx context.Context, nodeID string) error {
	node, err := s.repo.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}

	node.IsActive = false
	node.LoadPercentage = 0
	if err := s.repo.UpdateNode(ctx, node); err != nil {
		return err
	}

	shards, err := s.repo.ListShardsByNode(ctx, nodeID)
	if err != nil {
		return err
	}

	for _, sh := range shards {
		sh.Status = ShardDead
		if err := s.repo.UpdateShard(ctx, sh); err != nil {
			return err
		}
	}

	activeNodes, err := s.repo.ListNodes(ctx)
	if err != nil {
		return err
	}
	availableCount := 0
	for _, n := range activeNodes {
		if n.IsActive {
			availableCount++
		}
	}

	if availableCount > 0 {
		return s.RebalanceShards(ctx)
	}

	return nil
}
