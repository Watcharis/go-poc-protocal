package router

import (
	"context"
	"net/http"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/trace"
	"watcharis/go-poc-protocal/restful_api/ratelimit/handlers"

	"github.com/rs/cors"
)

func InitRouter(ctx context.Context, handlers handlers.RestFulAPIHandlers) http.Handler {

	mux := http.NewServeMux()

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Restrict to specific origin
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	})

	mux.HandleFunc("GET /health", pkg.HealthCheck)

	mux.HandleFunc("POST /api/v1/create-user-profile", handlers.CreateUserProfile(ctx))
	mux.HandleFunc("GET /api/v1/get-user-profile", handlers.GetUserProfile(ctx))
	mux.HandleFunc("POST /api/v1/create-otp", handlers.CreateOtp(ctx))
	mux.HandleFunc("POST /api/v1/verify-otp-ratelimit", handlers.VerifyOtpRatelimit(ctx))

	return c.Handler(trace.MiddlewareAddTrace(ctx, mux))
}
