package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/shopee-clone/shopee/services/payment/internal/config"
	"go.uber.org/zap"
)

type EventHandler func(ctx context.Context, eventType string, payload []byte) error

type Consumer struct {
	reader  *kafka.Reader
	cfg     config.KafkaConfig
	handler EventHandler
}

func NewConsumer(cfg config.KafkaConfig, topics []string, handler EventHandler) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers, GroupID: cfg.ConsumerGroup, GroupTopics: topics,
		MinBytes: 1e3, MaxBytes: 10e6, MaxWait: 500 * time.Millisecond,
	})
	return &Consumer{reader: reader, cfg: cfg, handler: handler}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil { return ctx.Err() }
			zap.L().Error("failed to fetch kafka message", zap.Error(err))
			continue
		}
		var eventType string
		for _, h := range msg.Headers {
			if h.Key == "event_type" { eventType = string(h.Value); break }
		}
		if err := c.handler(ctx, eventType, msg.Value); err != nil {
			zap.L().Error("failed to handle kafka message", zap.String("event_type", eventType), zap.Error(err))
			continue
		}
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			zap.L().Error("failed to commit kafka message", zap.Error(err))
		}
	}
}

func (c *Consumer) Close() error { return c.reader.Close() }
