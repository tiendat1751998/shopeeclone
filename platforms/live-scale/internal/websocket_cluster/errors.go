package websocket_cluster

import "errors"

var (
	ErrNodeNotFound    = errors.New("ws cluster node not found")
	ErrNodeAlreadyExists = errors.New("ws cluster node already exists")
	ErrRoomNotFound    = errors.New("room not found in cluster")
	ErrNoAvailableNode = errors.New("no available ws node for assignment")
	ErrRoomAlreadyAssigned = errors.New("room already assigned to a node")
)
