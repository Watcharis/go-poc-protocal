package router

import (
	"context"
	"net/http"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/restful_api/ratelimit/handlers"

	"go.opentelemetry.io/otel"
)

func InitRouter(ctx context.Context, handlers handlers.RestFulAPIHandlers) http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", pkg.HealthCheck)

	mux.HandleFunc("POST /api/v1/create-user-profile", handlers.CreateUserProfile(ctx))
	mux.HandleFunc("POST /api/v1/create-otp", handlers.CreateOtp(ctx))
	mux.HandleFunc("POST /api/v1/verify-otp-ratelimit", handlers.VerifyOtpRatelimit(ctx))

	return MiddlewareAddTrace(ctx, mux)
	// return mux
}

func MiddlewareAddTrace(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tracer := otel.Tracer(logger.APP_NAME)
		ctx, span := tracer.Start(ctx, logger.PROJECT_RATELIMIT)
		defer span.End()

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
