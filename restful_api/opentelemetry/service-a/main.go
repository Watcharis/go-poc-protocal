package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/httpclient"
	"watcharis/go-poc-protocal/pkg/logger"

	middlewareTrace "watcharis/go-poc-protocal/pkg/trace"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	ENDPOINT_SERVICE_B = "http://localhost:8779/api/v1/svcb"
	PORT               = ":8778"
)

func handleA(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	if span == nil {
		w.Write([]byte("failed"))
	}
	logger.Info(ctx, "service-a process running")

	if err := callAnotherService(ctx); err != nil {
		logger.Error(ctx, "call service-b failed", zap.Error(err))
		w = pkg.SetHttpStatusCode(w, http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	logger.Info(ctx, "service-a success")

	w = pkg.SetHttpStatusCode(w, http.StatusOK)
	w.Write([]byte("ok"))
}

func InitRouter(ctx context.Context) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", pkg.HealthCheck)
	mux.Handle("GET /api/v1/svca", http.HandlerFunc(handleA))

	handler := middlewareTrace.MiddlewareWarpOtelHttp(ctx, mux)
	return handler
}

func callAnotherService(ctx context.Context) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ENDPOINT_SERVICE_B, nil)
	if err != nil {
		return nil
	}

	httpClient := httpclient.CreateOtelHttpClient()

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Print the response
	fmt.Println("response status_code:", resp.Status)
	fmt.Println("response body:", string(respBody))

	return nil
}

func main() {
	ctx := context.WithValue(context.Background(), dto.APP_NAME, dto.PROJECT_OPENTELEMETRY_SERVICE_A)

	tp, err := middlewareTrace.SetupTracer(ctx, dto.APP_NAME)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Panic(ctx, err.Error())
		}
	}()

	logger.InitOtelZapLogger("develop")
	defer logger.Sync()

	routeHandlers := InitRouter(ctx)

	httpServer := &http.Server{
		Addr:    PORT,
		Handler: routeHandlers,
	}

	go func(httpServer *http.Server) {
		defer httpServer.Close()
		logger.Info(ctx, "Server runnig on http://localhost"+httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Panic(ctx, "cannot start server", zap.Error(err))
		}
	}(httpServer)

	wg := new(sync.WaitGroup)
	signal := make(chan os.Signal, 1)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		s := <-signal
		fmt.Println("signal :", s)
	}()
	wg.Wait()
}
