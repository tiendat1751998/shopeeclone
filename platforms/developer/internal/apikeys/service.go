package apikeys

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func generateRawKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *Service) Generate(ctx context.Context, name string, permissions []string, serviceName string, expiresAt time.Time) (*APIKey, string, error) {
	rawKey := generateRawKey()
	hash := HashKey(rawKey)

	key := &APIKey{
		ID:          uuid.New().String(),
		Name:        name,
		KeyHash:     hash,
		Permissions: permissions,
		ServiceName: serviceName,
		ExpiresAt:   expiresAt,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Store(ctx, key); err != nil {
		return nil, "", err
	}
	return key, rawKey, nil
}

func (s *Service) Validate(ctx context.Context, rawKey string) (*APIKey, bool) {
	hash := HashKey(rawKey)
	key, err := s.repo.GetByKeyHash(ctx, hash)
	if err != nil || key == nil {
		return nil, false
	}
	if !key.IsActive {
		return nil, false
	}
	if !key.ExpiresAt.IsZero() && time.Now().After(key.ExpiresAt) {
		return nil, false
	}
	return key, true
}

func (s *Service) Revoke(ctx context.Context, id string) error {
	key, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if key == nil {
		return nil
	}
	key.IsActive = false
	return s.repo.Update(ctx, key)
}

func (s *Service) List(ctx context.Context) ([]*APIKey, error) {
	return s.repo.List(ctx)
}

func (s *Service) Rotate(ctx context.Context, id string) (*APIKey, string, error) {
	key, err := s.repo.GetByID(ctx, id)
	if err != nil || key == nil {
		return nil, "", err
	}
	oldHash := key.KeyHash
	rawKey := generateRawKey()
	hash := HashKey(rawKey)
	key.KeyHash = hash
	key.CreatedAt = time.Now()

	if err := s.repo.Update(ctx, key); err != nil {
		return nil, "", err
	}
	if err := s.repo.DeleteByHash(ctx, oldHash); err != nil {
		return nil, "", err
	}
	return key, rawKey, nil
}
