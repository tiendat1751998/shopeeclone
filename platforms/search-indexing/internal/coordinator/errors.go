package coordinator

import "errors"

var (
	ErrNodeNotFound       = errors.New("index node not found")
	ErrShardNotFound      = errors.New("shard not found")
	ErrNodeInactive       = errors.New("node is not active")
	ErrNoAvailableNodes   = errors.New("no available nodes")
	ErrShardAlreadyExists = errors.New("shard already exists")
)
