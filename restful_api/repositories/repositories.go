package repositories

import (
	"context"
	"time"
	"watcharis/go-poc-protocal/restful_api/models"
)

type RedisRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error)
	Hset(ctx context.Context, key string, values []string) (int64, error)
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
}

type ProfilesRepository interface {
	CreateUserProfile(ctx context.Context, data models.ProfileDB) (models.ProfileDB, error)
}

type OtpRepository interface {
	CreateOtp(ctx context.Context, data models.OtpDB) (models.OtpDB, error)
	GetOtp(ctx context.Context, uuid string, otp string) (models.OtpDB, error)
}
