package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{client: client}
}

func (s *Store) Client() *redis.Client { return s.client }

// Pre-computed key builders using string concatenation instead of fmt.Sprintf
// These are the hot paths called on every WebSocket event

func (s *Store) AddViewer(ctx context.Context, roomID, userID string) (int64, error) {
	key := "room:" + roomID + ":viewers"
	count, err := s.client.SAdd(ctx, key, userID).Result()
	if err != nil {
		return 0, fmt.Errorf("add viewer: %w", err)
	}
	s.client.Expire(ctx, key, 2*time.Hour)
	return count, nil
}

func (s *Store) RemoveViewer(ctx context.Context, roomID, userID string) (int64, error) {
	key := "room:" + roomID + ":viewers"
	count, err := s.client.SRem(ctx, key, userID).Result()
	if err != nil {
		return 0, fmt.Errorf("remove viewer: %w", err)
	}
	return count, nil
}

func (s *Store) GetViewerCount(ctx context.Context, roomID string) (int64, error) {
	key := "room:" + roomID + ":viewers"
	count, err := s.client.SCard(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("get viewer count: %w", err)
	}
	return count, nil
}

func (s *Store) GetViewers(ctx context.Context, roomID string) ([]string, error) {
	key := "room:" + roomID + ":viewers"
	return s.client.SMembers(ctx, key).Result()
}

func (s *Store) IncrementReaction(ctx context.Context, roomID, reactionType string) (int64, error) {
	key := "room:" + roomID + ":reactions:" + reactionType
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("incr reaction: %w", err)
	}
	s.client.Expire(ctx, key, 2*time.Hour)
	return count, nil
}

func (s *Store) GetReactionCount(ctx context.Context, roomID, reactionType string) (int64, error) {
	key := "room:" + roomID + ":reactions:" + reactionType
	return s.client.Get(ctx, key).Int64()
}

func (s *Store) GetAllReactionCounts(ctx context.Context, roomID string) (map[string]int64, error) {
	pattern := "room:" + roomID + ":reactions:*"
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("keys: %w", err)
	}
	result := make(map[string]int64, len(keys))
	prefix := "room:" + roomID + ":reactions:"
	for _, key := range keys {
		val, err := s.client.Get(ctx, key).Int64()
		if err != nil {
			continue
		}
		reactionType := key[len(prefix):]
		result[reactionType] = val
	}
	return result, nil
}

func (s *Store) AddGiftAmount(ctx context.Context, roomID string, amount int64) (int64, error) {
	key := "room:" + roomID + ":gifts:total"
	total, err := s.client.IncrBy(ctx, key, amount).Result()
	if err != nil {
		return 0, fmt.Errorf("add gift: %w", err)
	}
	s.client.Expire(ctx, key, 2*time.Hour)
	return total, nil
}

func (s *Store) GetGiftTotal(ctx context.Context, roomID string) (int64, error) {
	key := "room:" + roomID + ":gifts:total"
	return s.client.Get(ctx, key).Int64()
}

func (s *Store) SetRoomStatus(ctx context.Context, roomID, status string) error {
	key := "room:" + roomID + ":status"
	return s.client.Set(ctx, key, status, 2*time.Hour).Err()
}

func (s *Store) GetRoomStatus(ctx context.Context, roomID string) (string, error) {
	key := "room:" + roomID + ":status"
	return s.client.Get(ctx, key).Result()
}

func (s *Store) SetUserMuted(ctx context.Context, roomID, userID string, duration time.Duration) error {
	key := "room:" + roomID + ":muted:" + userID
	return s.client.Set(ctx, key, "1", duration).Err()
}

func (s *Store) IsUserMuted(ctx context.Context, roomID, userID string) (bool, error) {
	key := "room:" + roomID + ":muted:" + userID
	_, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) SetUserBanned(ctx context.Context, roomID, userID string) error {
	key := "room:" + roomID + ":banned:" + userID
	return s.client.Set(ctx, key, "1", 24*time.Hour).Err()
}

func (s *Store) IsUserBanned(ctx context.Context, roomID, userID string) (bool, error) {
	key := "room:" + roomID + ":banned:" + userID
	_, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) AddToGiftLeaderboard(ctx context.Context, roomID, userID, username string, amount int64) error {
	key := "room:" + roomID + ":giftlb"
	member := userID + ":" + username
	return s.client.ZIncrBy(ctx, key, float64(amount), member).Err()
}

func (s *Store) GetGiftLeaderboard(ctx context.Context, roomID string, limit int) ([]redis.Z, error) {
	key := "room:" + roomID + ":giftlb"
	return s.client.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
}

func (s *Store) StoreReplayEvent(ctx context.Context, roomID string, event interface{}) error {
	key := "room:" + roomID + ":replay"
	data, err := sonic.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal replay: %w", err)
	}
	_, err = s.client.RPush(ctx, key, data).Result()
	if err != nil {
		return fmt.Errorf("push replay: %w", err)
	}
	s.client.LTrim(ctx, key, -200, -1)
	s.client.Expire(ctx, key, 1*time.Hour)
	return nil
}

func (s *Store) GetReplayEvents(ctx context.Context, roomID string, sinceSeq int64) ([]string, error) {
	key := "room:" + roomID + ":replay"
	return s.client.LRange(ctx, key, sinceSeq, -1).Result()
}

func (s *Store) SetConnectionMapping(ctx context.Context, connID, roomID, userID string) error {
	key := "conn:" + connID
	mapping := map[string]string{"room_id": roomID, "user_id": userID}
	data, _ := sonic.Marshal(mapping)
	return s.client.Set(ctx, key, data, 1*time.Hour).Err()
}

func (s *Store) GetConnectionMapping(ctx context.Context, connID string) (roomID, userID string, err error) {
	key := "conn:" + connID
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return "", "", err
	}
	var mapping map[string]string
	if err := sonic.Unmarshal(data, &mapping); err != nil {
		return "", "", err
	}
	return mapping["room_id"], mapping["user_id"], nil
}

func (s *Store) RemoveConnectionMapping(ctx context.Context, connID string) error {
	key := "conn:" + connID
	return s.client.Del(ctx, key).Err()
}

func (s *Store) RateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	current, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if current == 1 {
		s.client.Expire(ctx, key, window)
	}
	return current <= int64(limit), nil
}

// FormatSeq formats a sequence number for replay key lookups (replaces any remaining Sprintf).
func FormatSeq(seq int64) string {
	return strconv.FormatInt(seq, 10)
}
