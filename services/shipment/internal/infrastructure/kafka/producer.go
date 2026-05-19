package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/shopee-clone/shopee/services/shipment/internal/config"
	"github.com/shopee-clone/shopee/services/shipment/internal/domain"
	"github.com/shopee-clone/shopee/services/shipment/internal/metrics"
)

type Producer struct {
	writer *kafka.Writer
	cfg    config.KafkaConfig
}

// [SECURITY] Whitelist of allowed event types to prevent topic injection attacks.
var allowedEventTypes = map[string]bool{
	"shipment.created":    true,
	"shipment.picked_up":  true,
	"shipment.in_transit": true,
	"shipment.delivered":  true,
	"shipment.failed":     true,
	"shipment.returned":   true,
}

func NewProducer(cfg config.KafkaConfig) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
			WriteTimeout: 10 * time.Second,
		},
		cfg: cfg,
	}
}

func (p *Producer) PublishEvent(ctx context.Context, event *domain.ShipmentEvent) error {
	payload, err := json.Marshal(event)
	if err != nil { return err }

	// [SECURITY] Validate event type against whitelist to prevent topic injection
	if !allowedEventTypes[event.EventType] {
		return fmt.Errorf("invalid event type: %s", event.EventType)
	}

	// [SECURITY] Use validated event type only
	topic := fmt.Sprintf("%s.%s", p.cfg.TopicPrefix, event.EventType)

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(event.ShipmentID),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(event.EventType)},
		},
	}

	start := time.Now()
	err = p.writer.WriteMessages(ctx, msg)
	metrics.KafkaPublishLatency.WithLabelValues(string(event.EventType)).Observe(time.Since(start).Seconds())

	if err != nil {
		metrics.KafkaPublishErrors.WithLabelValues(string(event.EventType)).Inc()
		return err
	}
	return nil
}

func (p *Producer) Close() error { return p.writer.Close() }
