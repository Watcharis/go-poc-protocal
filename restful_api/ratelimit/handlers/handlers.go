package handlers

import (
	"context"
	"net/http"
	"watcharis/go-poc-protocal/restful_api/ratelimit/services"
)

type RestFulAPIHandlers interface {
	CreateUserProfile(ctx context.Context) http.HandlerFunc
	CreateOtp(ctx context.Context) http.HandlerFunc
	VerifyOtpRatelimit(ctx context.Context) http.HandlerFunc
}

type restFulAPIHandlers struct {
	services services.Services
}

func NewRestFulAPIHandlers(services services.Services) RestFulAPIHandlers {
	return &restFulAPIHandlers{
		services: services,
	}
}
