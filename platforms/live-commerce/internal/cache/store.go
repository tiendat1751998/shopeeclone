package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/shopee-clone/shopee/platforms/live-commerce/internal/domain"
)

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{client: client}
}

func (s *Store) GetLivestream(ctx context.Context, id string) (*domain.Livestream, error) {
	key := fmt.Sprintf("cache:livestream:%s", id)
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	ls := &domain.Livestream{}
	if err := json.Unmarshal(data, ls); err != nil {
		return nil, err
	}
	return ls, nil
}

func (s *Store) SetLivestream(ctx context.Context, ls *domain.Livestream) error {
	key := fmt.Sprintf("cache:livestream:%s", ls.ID)
	data, err := json.Marshal(ls)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, data, 5*time.Minute).Err()
}

func (s *Store) InvalidateLivestream(ctx context.Context, id string) error {
	return s.client.Del(ctx, fmt.Sprintf("cache:livestream:%s", id)).Err()
}

func (s *Store) GetActiveLivestreams(ctx context.Context, offset, limit int) ([]*domain.Livestream, error) {
	key := "cache:livestreams:active"
	data, err := s.client.LRange(ctx, key, int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		return nil, err
	}
	var result []*domain.Livestream
	for _, d := range data {
		ls := &domain.Livestream{}
		if err := json.Unmarshal([]byte(d), ls); err != nil {
			continue
		}
		result = append(result, ls)
	}
	return result, nil
}

func (s *Store) SetActiveLivestreams(ctx context.Context, livestreams []*domain.Livestream) error {
	key := "cache:livestreams:active"
	s.client.Del(ctx, key)
	pipe := s.client.Pipeline()
	for _, ls := range livestreams {
		data, _ := json.Marshal(ls)
		pipe.RPush(ctx, key, data)
	}
	pipe.Expire(ctx, key, 1*time.Minute)
	_, err := pipe.Exec(ctx)
	return err
}
