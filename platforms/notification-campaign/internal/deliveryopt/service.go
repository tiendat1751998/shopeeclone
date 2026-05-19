package deliveryopt

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrThrottled     = errors.New("delivery: throttled")
	ErrAllChannelsFailed = errors.New("delivery: all channels failed")
)

type Service interface {
	OptimizeSendTime(ctx context.Context, userID string, channel string) (*SendTimeOptimization, error)
	AnalyzePattern(ctx context.Context, userID string, channel string, openHour int, clickHour int) error
	Enqueue(ctx context.Context, msg *QueuedMessage) error
	Dequeue(ctx context.Context) (*QueuedMessage, error)
	SendWithFallback(ctx context.Context, req *SendRequest, cfg *ThrottleConfig) (*SendResult, error)
	CheckThrottle(ctx context.Context, channel string, cfg *ThrottleConfig) (bool, error)
	RecordSend(ctx context.Context, channel string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) OptimizeSendTime(ctx context.Context, userID string, channel string) (*SendTimeOptimization, error) {
	pattern, err := s.repo.GetPattern(ctx, userID, channel)
	if err != nil {
		return nil, err
	}
	if pattern == nil || pattern.SampleSize < 3 {
		return &SendTimeOptimization{
			UserID:     userID,
			Channel:    channel,
			BestHour:   time.Now().Hour(),
			Confidence: "low",
		}, nil
	}

	confidence := "medium"
	if pattern.SampleSize >= 20 {
		confidence = "high"
	}

	bestHour := pattern.PeakOpenHour
	if pattern.ClickRate > pattern.OpenRate {
		bestHour = pattern.PeakClickHour
	}

	return &SendTimeOptimization{
		UserID:     userID,
		Channel:    channel,
		BestHour:   bestHour,
		Confidence: confidence,
	}, nil
}

func (s *service) AnalyzePattern(ctx context.Context, userID string, channel string, openHour int, clickHour int) error {
	pattern, err := s.repo.GetPattern(ctx, userID, channel)
	if err != nil {
		return err
	}

	if pattern == nil {
		pattern = &UserEngagementPattern{
			UserID:        userID,
			Channel:       channel,
			PeakOpenHour:  openHour,
			PeakClickHour: clickHour,
		}
	} else {
		pattern.PeakOpenHour = (pattern.PeakOpenHour*pattern.SampleSize + openHour) / (pattern.SampleSize + 1)
		pattern.PeakClickHour = (pattern.PeakClickHour*pattern.SampleSize + clickHour) / (pattern.SampleSize + 1)
	}
	pattern.SampleSize++
	pattern.UpdatedAt = time.Now()

	return s.repo.SavePattern(ctx, pattern)
}

func (s *service) Enqueue(ctx context.Context, msg *QueuedMessage) error {
	return s.repo.Enqueue(ctx, msg)
}

func (s *service) Dequeue(ctx context.Context) (*QueuedMessage, error) {
	return s.repo.Dequeue(ctx)
}

func (s *service) SendWithFallback(ctx context.Context, req *SendRequest, cfg *ThrottleConfig) (*SendResult, error) {
	plan := ChannelFallbackPlan{Channels: DefaultFallbackOrder}
	if req.Channel != "" {
		plan.Channels = []string{req.Channel}
	}

	for i, ch := range plan.Channels {
		if cfg != nil {
			ok, err := s.CheckThrottle(ctx, ch, cfg)
			if err != nil {
				return nil, err
			}
			if !ok {
				if i < len(plan.Channels)-1 {
					continue
				}
				return nil, ErrThrottled
			}
		}

		if err := s.RecordSend(ctx, ch); err != nil {
			return nil, err
		}

		return &SendResult{
			Success:      true,
			FinalChannel: ch,
		}, nil
	}

	return nil, ErrAllChannelsFailed
}

func (s *service) CheckThrottle(ctx context.Context, channel string, cfg *ThrottleConfig) (bool, error) {
	if cfg == nil {
		return true, nil
	}

	rate, ok := cfg.ChannelMessagesPerHour[channel]
	if !ok {
		return true, nil
	}
	if rate <= 0 {
		return false, nil
	}

	since := time.Now().Add(-1 * time.Hour)
	count, err := s.repo.GetChannelSendCount(ctx, channel, since)
	if err != nil {
		return false, err
	}

	return count < rate, nil
}

func (s *service) RecordSend(ctx context.Context, channel string) error {
	return s.repo.RecordSend(ctx, channel)
}

func (s *StubSender) Send(ctx context.Context, channel string, userID string, subject string, body string) error {
	s.Sent = append(s.Sent, StubSend{Channel: channel, UserID: userID, Subject: subject, Body: body})
	return nil
}

type StubSend struct {
	Channel string
	UserID  string
	Subject string
	Body    string
}

type StubSender struct {
	Sent    []StubSend
	FailOn  map[string]bool
}

func NewStubSender() *StubSender {
	return &StubSender{
		Sent:   make([]StubSend, 0),
		FailOn: make(map[string]bool),
	}
}

type SendFunc func(channel string) error

func CombineErrors(errs []error) error {
	var msgs []string
	for _, e := range errs {
		if e != nil {
			msgs = append(msgs, e.Error())
		}
	}
	if len(msgs) == 0 {
		return nil
	}
	return fmt.Errorf("combined: %s", joinStrings(msgs, "; "))
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
