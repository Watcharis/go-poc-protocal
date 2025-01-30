package repositories

import (
	"context"
	"time"
)

type RedisRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error)
	Hset(ctx context.Context, key string, values []string) (int64, error)
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
}
