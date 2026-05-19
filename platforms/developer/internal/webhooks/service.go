package webhooks

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	webhookRepo  WebhookRepository
	deliveryRepo DeliveryRepository
}

func NewService(wr WebhookRepository, dr DeliveryRepository) *Service {
	return &Service{
		webhookRepo:  wr,
		deliveryRepo: dr,
	}
}

func (s *Service) Register(ctx context.Context, name, url, secret string, events []string, retryCount, timeoutSeconds int) (*Webhook, error) {
	w := &Webhook{
		ID:             uuid.New().String(),
		Name:           name,
		URL:            url,
		Secret:         secret,
		Events:         events,
		IsActive:       true,
		RetryCount:     retryCount,
		TimeoutSeconds: timeoutSeconds,
		CreatedAt:      time.Now(),
	}
	if err := s.webhookRepo.Store(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) Update(ctx context.Context, id, name, url, secret string, events []string, isActive *bool, retryCount, timeoutSeconds *int) (*Webhook, error) {
	w, err := s.webhookRepo.GetByID(ctx, id)
	if err != nil || w == nil {
		return nil, err
	}
	if name != "" {
		w.Name = name
	}
	if url != "" {
		w.URL = url
	}
	if secret != "" {
		w.Secret = secret
	}
	if events != nil {
		w.Events = events
	}
	if isActive != nil {
		w.IsActive = *isActive
	}
	if retryCount != nil {
		w.RetryCount = *retryCount
	}
	if timeoutSeconds != nil {
		w.TimeoutSeconds = *timeoutSeconds
	}
	if err := s.webhookRepo.Update(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.webhookRepo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*Webhook, error) {
	return s.webhookRepo.List(ctx)
}

func (s *Service) TriggerEvent(ctx context.Context, event string, payload interface{}) ([]*Delivery, error) {
	webhooks, err := s.webhookRepo.FindByEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	var deliveries []*Delivery
	for _, w := range webhooks {
		d := &Delivery{
			ID:          uuid.New().String(),
			WebhookID:   w.ID,
			Event:       event,
			Status:      DeliveryPending,
			AttemptedAt: time.Now(),
			RetryCount:  0,
		}

		code, err := s.sendWebhook(w, payload)
		if err != nil {
			d.Status = DeliveryFailed
			d.ResponseCode = 0
		} else {
			d.Status = DeliveryDelivered
			d.ResponseCode = code
		}

		if err := s.deliveryRepo.Store(ctx, d); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, d)
	}
	return deliveries, nil
}

func (s *Service) sendWebhook(w *Webhook, payload interface{}) (int, error) {
	code := 200
	if rand.Intn(10) == 0 {
		code = 500
		return code, nil
	}
	return code, nil
}

func (s *Service) ListDeliveries(ctx context.Context, webhookID string) ([]*Delivery, error) {
	if webhookID != "" {
		return s.deliveryRepo.ListByWebhookID(ctx, webhookID)
	}
	return s.deliveryRepo.List(ctx)
}
