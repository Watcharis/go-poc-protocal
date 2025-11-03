package router

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/pprof"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/trace"
	"watcharis/go-poc-protocal/restful_api/ratelimit/handlers"

	"github.com/rs/cors"
)

func InitRouter(ctx context.Context, handlers handlers.RestFulAPIHandlers) http.Handler {

	mux := http.NewServeMux()

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Restrict to specific origin
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "*"},
		// AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	InitProfiling(mux)

	mux.HandleFunc("GET /health", pkg.HealthCheck)

	mux.HandleFunc("POST /api/v1/create-user-profile", handlers.CreateUserProfile(ctx))
	mux.HandleFunc("GET /api/v1/get-user-profile", handlers.GetUserProfile(ctx))
	mux.HandleFunc("POST /api/v1/create-otp", handlers.CreateOtp(ctx))
	mux.HandleFunc("POST /api/v1/verify-otp-ratelimit", handlers.VerifyOtpRatelimit(ctx))
	mux.HandleFunc("GET /api/v1/cal", func(w http.ResponseWriter, r *http.Request) {
		// run heavy calculation and return result to client
		result := heavyCalculation(20)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		// write the numeric result as plain text
		_, _ = fmt.Fprintf(w, "%f", result)
	})

	return c.Handler(trace.MiddlewareAddTrace(ctx, mux))
}

func InitProfiling(mux *http.ServeMux) {
	// init route profiling CPU
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// init route profiling Memory
	mux.HandleFunc("/debug/pprof/heap", pprof.Index)
}

func heavyCalculation(n int) float64 {
	var result float64
	for i := 1; i < n; i++ {
		result += math.Sqrt(float64(i)) * math.Sin(float64(i))
	}
	return result
}
