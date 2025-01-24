package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/models"
)

func (s *services) CreateOtp(ctx context.Context, req models.OtpRequest) (models.OtpResponse, error) {

	if req.Otp == "" {
		err := fmt.Errorf("otp is required")
		log.Printf("[error] request missing otp : %v", err)
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
		CreatedAt: time.Now(),
	}

	otp, err := s.otpRepository.CreateOtp(ctx, dataOTP)
	if err != nil {
		log.Printf("[error] set otp to db : %v", err)
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
		log.Printf("[error] set otp to redis : %v", err)
		return models.OtpResponse{
			CommonResponse: response.CommonResponse{
				Status: response.STATUS_ERROR,
				Code:   http.StatusInternalServerError,
			},
			Error: &response.ErrorResponse{ErrorMessage: err.Error()},
		}, nil
	}

	fmt.Printf("result : %v\n", result)
	return models.OtpResponse{
		CommonResponse: response.CommonResponse{
			Status: response.STATUS_SUCCESS,
			Code:   http.StatusOK,
		},
	}, nil
}
