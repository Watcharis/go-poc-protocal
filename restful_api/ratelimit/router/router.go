package router

import (
	"context"
	"net/http"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/trace"
	"watcharis/go-poc-protocal/restful_api/ratelimit/handlers"
)

func InitRouter(ctx context.Context, handlers handlers.RestFulAPIHandlers) http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", pkg.HealthCheck)

	mux.HandleFunc("POST /api/v1/create-user-profile", handlers.CreateUserProfile(ctx))
	mux.HandleFunc("POST /api/v1/create-otp", handlers.CreateOtp(ctx))
	mux.HandleFunc("POST /api/v1/verify-otp-ratelimit", handlers.VerifyOtpRatelimit(ctx))

	return trace.MiddlewareAddTrace(ctx, mux)
	// return mux
}
