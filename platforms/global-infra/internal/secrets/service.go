package secrets

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, secret *Secret) error
	Get(ctx context.Context, id string) (*Secret, error)
	GetByName(ctx context.Context, name, serviceName string) (*Secret, error)
	List(ctx context.Context, serviceName string) ([]*Secret, error)
	Update(ctx context.Context, secret *Secret) error
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	secrets map[string]*Secret
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		secrets: make(map[string]*Secret),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, secret *Secret) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	secret.CreatedAt = now
	secret.UpdatedAt = now
	secret.Version = 1
	r.secrets[secret.ID] = secret
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (*Secret, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.secrets[id]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *InMemoryRepository) GetByName(ctx context.Context, name, serviceName string) (*Secret, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.secrets {
		if s.Name == name && s.ServiceName == serviceName {
			return s, nil
		}
	}
	return nil, nil
}

func (r *InMemoryRepository) List(ctx context.Context, serviceName string) ([]*Secret, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Secret
	for _, s := range r.secrets {
		if serviceName != "" && s.ServiceName != serviceName {
			continue
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, secret *Secret) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.secrets[secret.ID]; !ok {
		return fmt.Errorf("secret not found: %s", secret.ID)
	}
	secret.UpdatedAt = time.Now()
	r.secrets[secret.ID] = secret
	return nil
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, secret *Secret) (*Secret, error) {
	if secret.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if secret.Value == "" {
		return nil, fmt.Errorf("value is required")
	}
	if secret.ServiceName == "" {
		return nil, fmt.Errorf("service_name is required")
	}
	if secret.RotationPeriod <= 0 {
		secret.RotationPeriod = 30
	}
	secret.ID = uuid.New().String()
	secret.Version = 1
	secret.LastRotated = time.Now()
	encoded := base64.StdEncoding.EncodeToString([]byte(secret.Value))
	secret.Value = encoded
	if err := s.repo.Create(ctx, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Secret, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) GetByName(ctx context.Context, name, serviceName string) (*Secret, error) {
	return s.repo.GetByName(ctx, name, serviceName)
}

func (s *Service) List(ctx context.Context, serviceName string) ([]*Secret, error) {
	return s.repo.List(ctx, serviceName)
}

func (s *Service) Rotate(ctx context.Context, id string) (*Secret, error) {
	secret, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found: %s", id)
	}
	decoded, err := base64.StdEncoding.DecodeString(secret.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret: %w", err)
	}
	newValue := deriveNewValue(string(decoded), secret.Name, secret.ServiceName, secret.Version+1)
	secret.Value = base64.StdEncoding.EncodeToString(newValue)
	secret.Version++
	secret.LastRotated = time.Now()
	secret.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func deriveNewValue(current, name, serviceName string, version int) []byte {
	data := fmt.Sprintf("%s:%s:%s:%d", current, name, serviceName, version)
	hash := sha256.Sum256([]byte(data))
	return hash[:16]
}
