package sfu

import "errors"

var (
	ErrNodeNotFound        = errors.New("sfu node not found")
	ErrNodeAlreadyExists   = errors.New("sfu node already exists")
	ErrNoAvailableNodes    = errors.New("no available sfu nodes")
	ErrStreamSessionNotFound = errors.New("stream session not found")
	ErrNodeAtCapacity      = errors.New("sfu node at capacity")
	ErrInvalidNodeData     = errors.New("invalid node data")
)
