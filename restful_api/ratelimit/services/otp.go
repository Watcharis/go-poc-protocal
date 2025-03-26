package services

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"

	"go.uber.org/zap"
)

func (s *services) CreateOtp(ctx context.Context, req models.OtpRequest) (models.OtpResponse, error) {
	if req.Otp == "" {
		err := fmt.Errorf("otp is required")
		logger.Error(ctx, "request missing otp", zap.String("uuid", req.UUID), zap.Error(err))
		return models.OtpResponse{
			CommonResponse: response.CommonResponse{
				Status: response.STATUS_ERROR,
				Code:   http.StatusBadRequest,
			},
			Error: &response.ErrorResponse{ErrorMessage: err.Error()},
		}, nil
	}

	dataOTP := models.OtpDB{
		UUID:      req.UUID,
		Otp:       req.Otp,
		CreatedAt: TimeNow(),
	}

	otp, err := s.otpRepository.CreateOtp(ctx, dataOTP)
	if err != nil {
		logger.Error(ctx, "[error] set otp to db failed", zap.String("uuid", req.UUID), zap.Error(err))
		return models.OtpResponse{
			CommonResponse: response.CommonResponse{
				Status: response.STATUS_ERROR,
				Code:   http.StatusInternalServerError,
			},
			Error: &response.ErrorResponse{ErrorMessage: err.Error()},
		}, nil
	}

	result, err := s.redis.Set(ctx, fmt.Sprintf(models.REDIS_OTP, otp.UUID), otp.Otp, time.Duration(models.OTP_EXPIRE))
	if err != nil {
		logger.Error(ctx, "[error] set otp to redis failed", zap.String("uuid", req.UUID), zap.Error(err))
		return models.OtpResponse{
			CommonResponse: response.CommonResponse{
				Status: response.STATUS_ERROR,
				Code:   http.StatusInternalServerError,
			},
			Error: &response.ErrorResponse{ErrorMessage: err.Error()},
		}, nil
	}

	logger.Info(ctx, "set redis otp success", zap.String("result", result))
	return models.OtpResponse{
		CommonResponse: response.CommonResponse{
			Status: response.STATUS_SUCCESS,
			Code:   http.StatusOK,
		},
	}, nil
}
