package models

import "watcharis/go-poc-protocal/pkg/response"

type VerifyOtpRatelimitRequest struct {
	Uuid string `json:"uuid"`
	Otp  string `json:"otp"`
}

type VerifyOtpRatelimitResponse struct {
	response.CommonResponse
	Data  *VerifyOtpRatelimitDataResponse `json:"data,omitempty"`
	Error *response.ErrorResponse         `json:"error,omitempty"`
}

type VerifyOtpRatelimitDataResponse struct {
	Otp string `json:"otp"`
}
