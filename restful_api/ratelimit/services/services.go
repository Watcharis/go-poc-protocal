package services

import (
	"context"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"
	"watcharis/go-poc-protocal/restful_api/ratelimit/repositories"
)

type Services interface {
	CreateUserProfile(ctx context.Context, req models.ProifleRequest) (models.ProifleResponse, error)
	CreateOtp(ctx context.Context, req models.OtpRequest) (models.OtpResponse, error)
	VerifyOtpRatelimit(ctx context.Context, req models.VerifyOtpRatelimitRequest) (models.VerifyOtpRatelimitResponse, error)
}

type services struct {
	redis              repositories.RedisRepository
	profilesRepository repositories.ProfilesRepository
	otpRepository      repositories.OtpRepository
}

func NewServices(redis repositories.RedisRepository,
	profilesRepository repositories.ProfilesRepository,
	otpRepository repositories.OtpRepository) Services {
	return &services{
		redis:              redis,
		profilesRepository: profilesRepository,
		otpRepository:      otpRepository,
	}
}
