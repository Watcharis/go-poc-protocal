package cache

import (
	"context"
	"time"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"

	"github.com/redis/go-redis/v9"
)

// mockgen -source=cache/redis.go -destination=cache/mocks/redis_mock.go -package=mocks

type RedisRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error)
	Hset(ctx context.Context, key string, values []string) (int64, error)
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	Hgetall(ctx context.Context, key string) (map[string]string, error)
	HgetallProfile(ctx context.Context, key string) (models.ProfileDB, error)
}

type redisRepository struct {
	redisClient *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) RedisRepository {
	return &redisRepository{
		redisClient: redisClient,
	}
}

func (r *redisRepository) Increment(ctx context.Context, key string) (int64, error) {
	result, err := r.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *redisRepository) Hset(ctx context.Context, key string, values []string) (int64, error) {
	result, err := r.redisClient.HSet(ctx, key, values).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *redisRepository) Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error) {
	result, err := r.redisClient.Set(ctx, key, value, expiration).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	result, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *redisRepository) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	result, err := r.redisClient.Expire(ctx, key, expiration).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

func (r *redisRepository) Hgetall(ctx context.Context, key string) (map[string]string, error) {
	result, err := r.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *redisRepository) HgetallProfile(ctx context.Context, key string) (models.ProfileDB, error) {
	var result models.ProfileDB
	if err := r.redisClient.HGetAll(ctx, key).Scan(&result); err != nil {
		return models.ProfileDB{}, err
	}
	return result, nil
}
