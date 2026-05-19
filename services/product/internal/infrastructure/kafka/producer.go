package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"go.uber.org/zap"
)

// Producer wraps kafka.Writer for publishing events
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:          kafka.TCP(brokers...),
			Balancer:      &kafka.LeastBytes{},
			BatchTimeout:  10 * time.Millisecond,
			WriteTimeout:  10 * time.Second,
			Async:         false,
		},
	}
}

// Publish publishes a message to a Kafka topic
func (p *Producer) Publish(ctx context.Context, topic string, key string, payload []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		observability.GetLogger().Error("failed to publish kafka message",
			zap.String("topic", topic),
			zap.String("key", key),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.writer.Close()
}

// EventPublisher interface for domain events
type EventPublisher interface {
	Publish(ctx context.Context, topic string, key string, payload []byte) error
}

// Ensure interface is satisfied
var _ EventPublisher = (*Producer)(nil)

// MarshalEvent is a helper to marshal events
func MarshalEvent(event interface{}) ([]byte, error) {
	return json.Marshal(event)
}
