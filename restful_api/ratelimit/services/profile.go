package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *services) CreateUserProfile(ctx context.Context, req models.ProifleRequest) (models.ProifleResponse, error) {
	// Prepare data
	uuid := pkg.GenerateUUID()

	profile := models.ProfileDB{
		UUID:      uuid,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		CreatedAt: time.Now(),
	}

	// Call repository
	profile, err := s.profilesRepository.CreateUserProfile(ctx, profile)
	if err != nil {
		logger.Error(ctx, "[error] create user profile to db failed", zap.Error(err))
		return models.ProifleResponse{}, err
	}

	profileValue := []string{
		"id", fmt.Sprintf("%d", profile.ID),
		"uuid", profile.UUID,
		"firstname", profile.FirstName,
		"lastname", profile.LastName,
		"email", profile.Email,
		"phone", profile.Phone,
		"created_at", profile.CreatedAt.Format(time.DateTime),
	}

	redisProflekey := fmt.Sprintf(models.REDIS_USER_PROFILE, profile.UUID)
	_, err = s.redis.Hset(ctx, redisProflekey, profileValue)
	if err != nil {
		logger.Error(ctx, "[error] set user profile to redis failed", zap.Error(err))
		return models.ProifleResponse{}, err
	}

	return models.ProifleResponse{
		CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
	}, nil
}

func (s *services) GetUserProfile(ctx context.Context, uuid string) (models.ProifleResponse, error) {

	var profile models.ProfileDB

	redisProflekey := fmt.Sprintf(models.REDIS_USER_PROFILE, uuid)

	profileRedis, err := s.redis.HgetallProfile(ctx, redisProflekey)
	if err != nil && errors.Is(err, redis.Nil) || profileRedis == (models.ProfileDB{}) {
		logger.Info(ctx, "not found profile in redis, query profile from db instead", zap.String("rediskey", redisProflekey), zap.String("uuid", uuid))

		profileDB, err := s.profilesRepository.GetUserProfile(ctx, uuid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Error(ctx, "[error] not found user profile from db", zap.String("uuid", uuid), zap.Error(err))
				return models.ProifleResponse{
					CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusNotFound),
					Data:           nil,
					Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
				}, nil
			}
			logger.Error(ctx, "[error] get user profile from db failed", zap.String("uuid", uuid), zap.Error(err))
			return models.ProifleResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
				Data:           nil,
				Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
			}, nil
		}
		profile = profileDB

	} else if err != nil {
		logger.Error(ctx, "[error] get user profile from redis failed", zap.String("uuid", uuid), zap.Error(err))
		return models.ProifleResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
			Data:           nil,
			Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
		}, nil

	} else {
		logger.Info(ctx, "get user profile from redis success", zap.String("rediskey", redisProflekey), zap.String("uuid", uuid))
		profile = profileRedis
	}

	logger.Info(ctx, "get user profile success", zap.String("uuid", uuid), zap.Any("profile", profile))

	return models.ProifleResponse{
		CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
		Data: &models.ProfileDB{
			ID:        profile.ID,
			UUID:      profile.UUID,
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			Email:     profile.Email,
			Phone:     profile.Phone,
			CreatedAt: profile.CreatedAt,
		},
	}, nil
}
