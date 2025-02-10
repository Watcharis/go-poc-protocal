package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *services) VerifyOtpRatelimit(ctx context.Context, req models.VerifyOtpRatelimitRequest) (models.VerifyOtpRatelimitResponse, error) {
	logger.Info(ctx, "service - VerifyOtpRatelimit", zap.Any("req", req))

	redisKeyRatelimitOTP := fmt.Sprintf(models.REDIS_RATELIMIT_OTP, req.Uuid)
	countOtp := int64(0)
	countOtpString, err := s.redis.Get(ctx, redisKeyRatelimitOTP)
	if err != nil {
		if errors.Is(err, redis.Nil) || countOtpString == "" {
			countOtp = 0
		} else {
			logger.Error(ctx, "[error] get otp ratelimit", zap.String("uuid", req.Uuid), zap.Error(err))
			return models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
				Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
			}, nil
		}
	} else {
		countOtp, err = strconv.ParseInt(countOtpString, 10, 64)
		if err != nil {
			logger.Error(ctx, "[error] parse otp ratelimit ", zap.String("uuid", req.Uuid), zap.Error(err))
			return models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
				Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
			}, nil
		}
	}

	if countOtp == models.RATELIMIT_OTP {
		logger.Error(ctx, "[error] otp ratelimit exceed count", zap.Int64("count", countOtp), zap.String("uuid", req.Uuid))
		return models.VerifyOtpRatelimitResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusTooManyRequests),
			Error:          &response.ErrorResponse{ErrorMessage: "otp ratelimit exceed"},
		}, nil
	}

	redisKeyOTP := fmt.Sprintf(models.REDIS_OTP, req.Uuid)
	otpFromRedis, err := s.redis.Get(ctx, redisKeyOTP)
	if err != nil {
		if errors.Is(err, redis.Nil) || otpFromRedis == "" {
			logger.Info(ctx, "otp not found in redis", zap.Error(err))

			// get otp from db
			otpDB, err := s.otpRepository.GetOtp(ctx, req.Uuid, req.Otp)
			if err != nil || errors.Is(err, gorm.ErrRecordNotFound) || otpDB == (models.OtpDB{}) || otpDB.Otp != req.Otp {
				// not found otp in db || otp not match
				logger.Info(ctx, "otp not found in DB or otp not match", zap.String("uuid", req.Uuid),
					zap.String("otp", req.Otp),
					zap.String("otpDB", otpDB.Otp))

				// start increment otp ratelimit
				countOtp, err := s.redis.Increment(ctx, redisKeyRatelimitOTP)
				if err != nil {
					logger.Error(ctx, "[error] increment otp ratelimit", zap.String("uuid", req.Uuid), zap.Error(err))
					return models.VerifyOtpRatelimitResponse{
						CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
						Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
					}, nil
				}
				logger.Info(ctx, "increment otp ratelimit", zap.Int64("countOtp", countOtp), zap.String("uuid", req.Uuid))

				if countOtp == models.RATELIMIT_OTP {
					logger.Info(ctx, "otp ratelimit exceed ", zap.Int64("countOtp", countOtp), zap.String("uuid", req.Uuid))
					// set expire otp ratelimit
					_, err := s.redis.Expire(ctx, redisKeyRatelimitOTP, time.Duration(models.RATELIMIT_OTP_EXPIRE))
					if err != nil {
						logger.Error(ctx, "[error] expire otp ratelimit", zap.String("uuid", req.Uuid), zap.Error(err))
						return models.VerifyOtpRatelimitResponse{
							CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
							Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
						}, nil
					}
					return models.VerifyOtpRatelimitResponse{
						CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusTooManyRequests),
						Error:          &response.ErrorResponse{ErrorMessage: "otp ratelimit exceed"},
					}, nil
				}

				return models.VerifyOtpRatelimitResponse{
					CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusNotFound),
					Error:          &response.ErrorResponse{ErrorMessage: "otp not match"},
				}, nil
			}

			// otp in db match
			otp, err := s.redis.Set(ctx, redisKeyOTP, req.Otp, time.Duration(models.OTP_EXPIRE))
			if err != nil {
				logger.Error(ctx, "[error] set otp to redis", zap.String("uuid", req.Uuid), zap.Error(err))
				return models.VerifyOtpRatelimitResponse{
					CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
					Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
				}, nil

			}

			logger.Info(ctx, "[error] set otp to redis", zap.String("otp", otp), zap.String("uuid", req.Uuid), zap.Error(err))
			return models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
				Data: &models.VerifyOtpRatelimitDataResponse{
					Otp: req.Otp,
				},
			}, nil
		}

		return models.VerifyOtpRatelimitResponse{
			CommonResponse: response.CommonResponse{
				Status: response.STATUS_ERROR,
				Code:   http.StatusInternalServerError,
			},
			Error: &response.ErrorResponse{ErrorMessage: err.Error()},
		}, nil
	}

	if otpFromRedis != req.Otp {
		logger.Info(ctx, "otp in redis not match with request otp", zap.String("uuid", req.Uuid),
			zap.String("redis_otp", otpFromRedis),
			zap.String("request_otp", req.Otp))

		// start increment otp ratelimit
		countOtp, err := s.redis.Increment(ctx, redisKeyRatelimitOTP)
		if err != nil {
			logger.Error(ctx, "[error] increment otp ratelimit", zap.String("uuid", req.Uuid), zap.Error(err))
			return models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
				Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
			}, nil
		}
		logger.Info(ctx, "increment otp ratelimit", zap.Int64("countOtp", countOtp), zap.String("uuid", req.Uuid))

		if countOtp == models.RATELIMIT_OTP {
			logger.Info(ctx, "otp ratelimit exceed ", zap.Int64("countOtp", countOtp), zap.String("uuid", req.Uuid))
			// set expire otp ratelimit
			_, err := s.redis.Expire(ctx, redisKeyRatelimitOTP, time.Duration(models.RATELIMIT_OTP_EXPIRE))
			if err != nil {
				logger.Error(ctx, "[error] expire otp ratelimit", zap.String("uuid", req.Uuid), zap.Error(err))
				return models.VerifyOtpRatelimitResponse{
					CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
					Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
				}, nil
			}
			return models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusTooManyRequests),
				Error:          &response.ErrorResponse{ErrorMessage: "otp ratelimit exceed"},
			}, nil
		}

		return models.VerifyOtpRatelimitResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusNotFound),
			Error:          &response.ErrorResponse{ErrorMessage: "otp not match"},
		}, nil
	}

	logger.Info(ctx, "verify otp success", zap.String("otp", req.Otp), zap.String("uuid", req.Uuid))
	return models.VerifyOtpRatelimitResponse{
		CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
		Data: &models.VerifyOtpRatelimitDataResponse{
			Otp: req.Otp,
		},
	}, nil
}
