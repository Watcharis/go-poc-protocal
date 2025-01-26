package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/logger"
	middlewareTrace "watcharis/go-poc-protocal/pkg/trace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var httpTransport *http.Transport
var otelHttpTransport *otelhttp.Transport

func CreateHttpClient() *http.Client {
	if httpTransport == nil {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxConnsPerHost = 10
		t.MaxIdleConns = 10
		t.MaxIdleConnsPerHost = 10
		t.IdleConnTimeout = 30 * time.Second
		t.ResponseHeaderTimeout = 30 * time.Second
		t.DisableKeepAlives = false
		// t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		httpTransport = t
	}

	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   30 * time.Second,
	}

	return httpClient
}

func CreateOtelHttpClient() *http.Client {
	// client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	if otelHttpTransport == nil {
		transport := &http.Transport{
			MaxConnsPerHost:       10,
			MaxIdleConns:          10,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       15 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
			DisableKeepAlives:     false,
		}
		otelHttpTransport = otelhttp.NewTransport(transport)
	}

	httpClient := &http.Client{
		Transport: otelHttpTransport,
		Timeout:   30 * time.Second,
	}

	return httpClient
}

func handleB(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	span := trace.SpanFromContext(ctx)
	if span == nil {
		w.Write([]byte("failed"))
	}

	logger.Info(ctx, "service-b success")
	w = pkg.SetHttpStatusCode(w, http.StatusOK)
	w.Write([]byte("ok"))
}

func InitRouter(ctx context.Context) http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", pkg.HealthCheck)

	mux.Handle("GET /api/v1/svcb", http.HandlerFunc(handleB))

	handler := middlewareTrace.MiddlewareAddTrace(ctx, mux)
	return handler
}

func main() {
	ctx := context.WithValue(context.Background(), dto.APP_NAME, dto.PROJECT_RATELIMIT)

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

	httpServer := http.Server{
		Addr:    ":8779",
		Handler: routeHandlers,
	}

	go func(port string) {
		defer httpServer.Close()
		logger.Info(ctx, "Server runnig on http://localhost"+port)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Panic(ctx, "cannot start server", zap.Error(err))
		}
	}(httpServer.Addr)

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
