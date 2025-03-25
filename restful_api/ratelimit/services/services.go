package services

import (
	"context"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"
	"watcharis/go-poc-protocal/restful_api/ratelimit/repositories/cache"
	"watcharis/go-poc-protocal/restful_api/ratelimit/repositories/db"
)

// mockgen -source=services/services.go -destination=services/mocks/services_mock.go -package=mocks
type Services interface {
	CreateUserProfile(ctx context.Context, req models.ProifleRequest) (models.ProifleResponse, error)
	CreateOtp(ctx context.Context, req models.OtpRequest) (models.OtpResponse, error)
	VerifyOtpRatelimit(ctx context.Context, req models.VerifyOtpRatelimitRequest) (models.VerifyOtpRatelimitResponse, error)
	GetUserProfile(ctx context.Context, uuid string) (models.ProifleResponse, error)
}

type services struct {
	redis              cache.RedisRepository
	profilesRepository db.ProfilesRepository
	otpRepository      db.OtpRepository
}

func NewServices(redis cache.RedisRepository,
	profilesRepository db.ProfilesRepository,
	otpRepository db.OtpRepository) Services {
	return &services{
		redis:              redis,
		profilesRepository: profilesRepository,
		otpRepository:      otpRepository,
	}
}
