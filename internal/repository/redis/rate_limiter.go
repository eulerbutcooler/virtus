package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{client: client}
}

func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, int, error) {
	redisKey := fmt.Sprintf("rl:%s", key)
	pipe := r.client.Pipeline()
	incr := pipe.Incr(ctx, redisKey)
	pipe.Expire(ctx, redisKey, window)
	if _, err := pipe.Exec(ctx); err != nil {
		return false, 0, fmt.Errorf("ratelimiter.Allow: %w", err)
	}
	count := int(incr.Val())
	remaining := limit - count
	if remaining < 0 {
		remaining = 0
	}
	return count <= limit, remaining, nil
}
