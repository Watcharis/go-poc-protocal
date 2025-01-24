package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/models"

	"github.com/redis/go-redis/v9"
)

func (s *services) VerifyOtpRatelimit(ctx context.Context, req models.VerifyOtpRatelimitRequest) (models.VerifyOtpRatelimitResponse, error) {
	redisKeyRatelimitOTP := fmt.Sprintf(models.REDIS_RATELIMIT_OTP, req.Uuid)
	countOtp := int64(0)
	countOtpString, err := s.redis.Get(ctx, redisKeyRatelimitOTP)
	if err != nil {
		if errors.Is(err, redis.Nil) || countOtpString == "" {
			countOtp = 0
		} else {
			log.Printf("[error] get otp ratelimit : %v", err)
			return models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
				Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
			}, nil
		}
	} else {
		countOtp, err = strconv.ParseInt(countOtpString, 10, 64)
		if err != nil {
			log.Printf("[error] parse otp ratelimit : %v", err)
			return models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
				Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
			}, nil
		}
	}

	if countOtp == models.RATELIMIT_OTP {
		log.Printf("[error] otp ratelimit exceed count = %d, uuid : %s", countOtp, req.Uuid)
		return models.VerifyOtpRatelimitResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusTooManyRequests),
			Error:          &response.ErrorResponse{ErrorMessage: "otp ratelimit exceed"},
		}, nil
	}

	redisKeyOTP := fmt.Sprintf(models.REDIS_OTP, req.Uuid)
	otpFromRedis, err := s.redis.Get(ctx, redisKeyOTP)
	if err != nil {
		if errors.Is(err, redis.Nil) || otpFromRedis == "" {
			log.Printf("otp not found in redis : %v", err)

			// get otp from db
			otpDB, err := s.otpRepository.GetOtp(ctx, req.Uuid, req.Otp)
			if err != nil {
				log.Printf("[error] get otp from db : %v", err)
				return models.VerifyOtpRatelimitResponse{
					CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusNotFound),
					Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
				}, nil
			}

			// not found otp in db || otp not match
			if otpDB == (models.OtpDB{}) || otpDB.Otp != req.Otp {
				log.Printf("[error] otp not found in DB or otp not match, uuid : %s, otp : %s, otpDB : %s", req.Uuid, req.Otp, otpDB.Otp)
				// start increment otp ratelimit
				countOtp, err := s.redis.Increment(ctx, redisKeyRatelimitOTP)
				if err != nil {
					log.Printf("[error] increment otp ratelimit : %v", err)
					return models.VerifyOtpRatelimitResponse{
						CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
						Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
					}, nil
				}
				log.Printf("increment otp ratelimit : %d\n", countOtp)

				if countOtp == models.RATELIMIT_OTP {
					log.Printf("[error] otp ratelimit exceed count = %d, uuid : %s", countOtp, req.Uuid)
					// set expire otp ratelimit
					_, err := s.redis.Expire(ctx, redisKeyRatelimitOTP, time.Duration(models.RATELIMIT_OTP_EXPIRE))
					if err != nil {
						log.Printf("[error] expire otp ratelimit : %v", err)
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
				log.Printf("[error] set otp to redis : %v", err)
				return models.VerifyOtpRatelimitResponse{
					CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
					Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
				}, nil

			}
			log.Printf("set otp to redis : %s\n", otp)
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
		log.Printf("[error] otp not match, uuid : %s, otp : %s, otpFromRedis : %s", req.Uuid, req.Otp, otpFromRedis)
		// start increment otp ratelimit
		countOtp, err := s.redis.Increment(ctx, redisKeyRatelimitOTP)
		if err != nil {
			log.Printf("[error] increment otp ratelimit : %v", err)
			return models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
				Error:          &response.ErrorResponse{ErrorMessage: err.Error()},
			}, nil
		}
		log.Printf("increment otp ratelimit : %d\n", countOtp)

		if countOtp == models.RATELIMIT_OTP {
			log.Printf("[error] otp ratelimit exceed count = %d, uuid : %s", countOtp, req.Uuid)
			// set expire otp ratelimit
			_, err := s.redis.Expire(ctx, redisKeyRatelimitOTP, time.Duration(models.RATELIMIT_OTP_EXPIRE))
			if err != nil {
				log.Printf("[error] expire otp ratelimit : %v", err)
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

	return models.VerifyOtpRatelimitResponse{
		CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
		Data: &models.VerifyOtpRatelimitDataResponse{
			Otp: req.Otp,
		},
	}, nil
}
