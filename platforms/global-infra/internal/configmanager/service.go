package configmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Publisher interface {
	Publish(ctx context.Context, event ConfigChangeEvent) error
	Close() error
}

type Service struct {
	repo      Repository
	publisher Publisher
	logger    *zap.Logger
}

func NewService(repo Repository, pub Publisher, logger *zap.Logger) *Service {
	return &Service{repo: repo, publisher: pub, logger: logger}
}

func (s *Service) Create(ctx context.Context, entry *ConfigEntry) (*ConfigEntry, error) {
	if entry.Key == "" {
		return nil, fmt.Errorf("key is required")
	}
	if entry.ServiceName == "" {
		return nil, fmt.Errorf("service_name is required")
	}
	if entry.Environment == "" {
		entry.Environment = EnvDev
	}
	if entry.Value == "" {
		return nil, fmt.Errorf("value is required")
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, err
	}
	s.publishEvent(ctx, entry, "created")
	return entry, nil
}

func (s *Service) Get(ctx context.Context, key string, env Environment, serviceName string) (*ConfigEntry, error) {
	return s.repo.Get(ctx, key, env, serviceName)
}

func (s *Service) GetVersion(ctx context.Context, key string, env Environment, serviceName string, version int) (*ConfigEntry, error) {
	return s.repo.GetVersion(ctx, key, env, serviceName, version)
}

func (s *Service) List(ctx context.Context, serviceName string, env Environment) ([]*ConfigEntry, error) {
	return s.repo.List(ctx, serviceName, env)
}

func (s *Service) ListVersions(ctx context.Context, key string, env Environment, serviceName string) ([]*ConfigEntry, error) {
	return s.repo.ListVersions(ctx, key, env, serviceName)
}

func (s *Service) publishEvent(ctx context.Context, entry *ConfigEntry, action string) {
	if s.publisher == nil {
		return
	}
	event := ConfigChangeEvent{
		Key:         entry.Key,
		Environment: entry.Environment,
		ServiceName: entry.ServiceName,
		Version:     entry.Version,
		Action:      action,
		Timestamp:   time.Now(),
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Warn("failed to publish config change event", zap.Error(err))
	}
}

type KafkaPublisher struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func NewKafkaPublisher(brokers []string, topic string, logger *zap.Logger) *KafkaPublisher {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaPublisher{writer: w, logger: logger}
}

func (p *KafkaPublisher) Publish(ctx context.Context, event ConfigChangeEvent) error {
	msg := kafka.Message{
		Key:   []byte(event.Key + ":" + string(event.Environment)),
		Value: []byte(fmt.Sprintf(`{"key":"%s","env":"%s","service":"%s","version":%d,"action":"%s"}`, event.Key, event.Environment, event.ServiceName, event.Version, event.Action)),
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}

type NoOpPublisher struct{}

func NewNoOpPublisher() *NoOpPublisher {
	return &NoOpPublisher{}
}

func (p *NoOpPublisher) Publish(ctx context.Context, event ConfigChangeEvent) error {
	return nil
}

func (p *NoOpPublisher) Close() error {
	return nil
}
