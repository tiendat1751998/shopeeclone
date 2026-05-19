package events

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		Async:        true,
		RequiredAcks: kafka.RequireOne,
	}
	return &KafkaProducer{writer: w}
}

func (p *KafkaProducer) Publish(ctx context.Context, event Event) error {
	if event.ID == "" {
		event.ID = newID()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ID),
		Value: data,
	})
}

func (p *KafkaProducer) PublishBatch(ctx context.Context, events []Event) error {
	msgs := make([]kafka.Message, len(events))
	for i, e := range events {
		if e.ID == "" {
			e.ID = newID()
		}
		if e.Timestamp.IsZero() {
			e.Timestamp = time.Now().UTC()
		}
		data, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("marshal event %d: %w", i, err)
		}
		msgs[i] = kafka.Message{Key: []byte(e.ID), Value: data}
	}
	return p.writer.WriteMessages(ctx, msgs...)
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

func newID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type NoopProducer struct{}

func (n *NoopProducer) Publish(_ context.Context, _ Event) error         { return nil }
func (n *NoopProducer) PublishBatch(_ context.Context, _ []Event) error  { return nil }
func (n *NoopProducer) Close() error                                     { return nil }
