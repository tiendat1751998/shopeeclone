package idempotency

import (
	"context"
	"errors"
	"time"
)

var (
	ErrKeyExists    = errors.New("idempotency key already exists")
	ErrKeyNotFound  = errors.New("idempotency key not found")
)

type Store interface {
	Exists(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

type KeyGenerator func(prefix, resourceID, operation string) string

func DefaultKeyGenerator(prefix, resourceID, operation string) string {
	return prefix + ":" + resourceID + ":" + operation
}

type Middleware struct {
	store Store
	ttl   time.Duration
	gen   KeyGenerator
}

func NewMiddleware(store Store, ttl time.Duration) *Middleware {
	return &Middleware{store: store, ttl: ttl, gen: DefaultKeyGenerator}
}

func (m *Middleware) CheckAndSet(ctx context.Context, key string) (bool, error) {
	exists, err := m.store.Exists(ctx, key)
	if err != nil {
		return false, err
	}
	if exists {
		return true, ErrKeyExists
	}
	if err := m.store.Set(ctx, key, nil, m.ttl); err != nil {
		return false, err
	}
	return false, nil
}

func (m *Middleware) Complete(ctx context.Context, key string, result []byte) error {
	return m.store.Set(ctx, key, result, m.ttl)
}
