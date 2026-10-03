package migration

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const legacySessionPrefix = "migration:registration:legacy:"

type redisSessionAffinity struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisSessionAffinity constructs the Redis-backed registration ownership
// store used to keep multi-step legacy fallback sessions sticky.
func NewRedisSessionAffinity(client *redis.Client, ttl time.Duration) (SessionAffinity, error) {
	if client == nil {
		return nil, ErrNilRedis
	}
	if ttl <= 0 {
		ttl = 20 * time.Minute
	}
	return &redisSessionAffinity{client: client, ttl: ttl}, nil
}

func (s *redisSessionAffinity) SetLegacy(ctx context.Context, sessionID string) error {
	if strings.TrimSpace(sessionID) == "" {
		return ErrEmptySessionID
	}
	return s.client.Set(ctx, legacySessionPrefix+sessionID, "legacy", s.ttl).Err()
}

func (s *redisSessionAffinity) IsLegacy(ctx context.Context, sessionID string) (bool, error) {
	if strings.TrimSpace(sessionID) == "" {
		return false, ErrEmptySessionID
	}
	value, err := s.client.Exists(ctx, legacySessionPrefix+sessionID).Result()
	return value == 1, err
}
