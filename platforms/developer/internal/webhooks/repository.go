package webhooks

import (
	"context"
	"sync"
)

type WebhookRepository interface {
	Store(ctx context.Context, w *Webhook) error
	GetByID(ctx context.Context, id string) (*Webhook, error)
	List(ctx context.Context) ([]*Webhook, error)
	Update(ctx context.Context, w *Webhook) error
	Delete(ctx context.Context, id string) error
	FindByEvent(ctx context.Context, event string) ([]*Webhook, error)
}

type DeliveryRepository interface {
	Store(ctx context.Context, d *Delivery) error
	ListByWebhookID(ctx context.Context, webhookID string) ([]*Delivery, error)
	List(ctx context.Context) ([]*Delivery, error)
}

type InMemoryWebhookRepository struct {
	mu       sync.RWMutex
	webhooks map[string]*Webhook
}

func NewInMemoryWebhookRepository() *InMemoryWebhookRepository {
	return &InMemoryWebhookRepository{
		webhooks: make(map[string]*Webhook),
	}
}

func (r *InMemoryWebhookRepository) Store(ctx context.Context, w *Webhook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.webhooks[w.ID] = w
	return nil
}

func (r *InMemoryWebhookRepository) GetByID(ctx context.Context, id string) (*Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.webhooks[id]
	if !ok {
		return nil, nil
	}
	return w, nil
}

func (r *InMemoryWebhookRepository) List(ctx context.Context) ([]*Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Webhook
	for _, w := range r.webhooks {
		result = append(result, w)
	}
	return result, nil
}

func (r *InMemoryWebhookRepository) Update(ctx context.Context, w *Webhook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.webhooks[w.ID] = w
	return nil
}

func (r *InMemoryWebhookRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.webhooks, id)
	return nil
}

func (r *InMemoryWebhookRepository) FindByEvent(ctx context.Context, event string) ([]*Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Webhook
	for _, w := range r.webhooks {
		if !w.IsActive {
			continue
		}
		for _, e := range w.Events {
			if e == event {
				result = append(result, w)
				break
			}
		}
	}
	return result, nil
}

type InMemoryDeliveryRepository struct {
	mu         sync.RWMutex
	deliveries map[string]*Delivery
}

func NewInMemoryDeliveryRepository() *InMemoryDeliveryRepository {
	return &InMemoryDeliveryRepository{
		deliveries: make(map[string]*Delivery),
	}
}

func (r *InMemoryDeliveryRepository) Store(ctx context.Context, d *Delivery) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deliveries[d.ID] = d
	return nil
}

func (r *InMemoryDeliveryRepository) ListByWebhookID(ctx context.Context, webhookID string) ([]*Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Delivery
	for _, d := range r.deliveries {
		if d.WebhookID == webhookID {
			result = append(result, d)
		}
	}
	return result, nil
}

func (r *InMemoryDeliveryRepository) List(ctx context.Context) ([]*Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Delivery
	for _, d := range r.deliveries {
		result = append(result, d)
	}
	return result, nil
}
