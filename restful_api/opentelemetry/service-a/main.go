package main

import (
	"context"
	"fmt"
	"io"
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

var otelHttpTransport *otelhttp.Transport
var httpTransport *http.Transport

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

		// การใช้ otelhttp.NewTransport สำหรับ HTTP Client เพื่อส่ง Trace Context ไปยัง Service ปลายทาง
		// Trace Context จะถูกส่งจาก Service A ไปยัง Service B
		otelHttpTransport = otelhttp.NewTransport(transport)
	}

	httpClient := &http.Client{
		Transport: otelHttpTransport,
		Timeout:   30 * time.Second,
	}

	return httpClient
}

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

	handler := middlewareTrace.MiddlewareAddTrace(ctx, mux)
	return handler
}

func callAnotherService(ctx context.Context) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8779/api/v1/svcb", nil)
	if err != nil {
		return nil
	}

	httpClient := CreateOtelHttpClient()

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
		Addr:    ":8778",
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
