package cache

import (
	"context"
	"time"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"

	"github.com/stretchr/testify/mock"
)

type RedisRepositoryMock struct {
	mock.Mock
}

// mockgen -source=cache/redis.go -destination=cache/mocks/redis_mock.go -package=mocks

func (m *RedisRepositoryMock) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *RedisRepositoryMock) Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error) {
	args := m.Called(ctx, key, value, expiration)
	return args.String(0), args.Error(1)
}

func (m *RedisRepositoryMock) Hset(ctx context.Context, key string, values []string) (int64, error) {
	args := m.Called(ctx, key, values)
	return args.Get(0).(int64), args.Error(1)
}

func (m *RedisRepositoryMock) Increment(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *RedisRepositoryMock) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	args := m.Called(ctx, key, expiration)
	return args.Bool(0), args.Error(1)
}

func (m *RedisRepositoryMock) Hgetall(ctx context.Context, key string) (map[string]string, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *RedisRepositoryMock) HgetallProfile(ctx context.Context, key string) (models.ProfileDB, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(models.ProfileDB), args.Error(1)
}
