package services

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"

	"go.uber.org/zap"
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
