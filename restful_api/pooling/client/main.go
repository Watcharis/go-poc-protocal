package main

import (
	"context"
	"log"
	"net/http"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/pkg/trace"
	"watcharis/go-poc-protocal/restful_api/pooling/client/handlers"
	"watcharis/go-poc-protocal/restful_api/pooling/client/repositories"
	"watcharis/go-poc-protocal/restful_api/pooling/client/router"
	"watcharis/go-poc-protocal/restful_api/pooling/client/services"
)

const (
	PORT = ":8780"
)

func main() {

	ctx := context.WithValue(context.Background(), dto.APP_NAME, dto.PROJECT_POOLING_CLIENT)

	tp, err := trace.SetupTracer(ctx, dto.APP_NAME)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Panic(ctx, err.Error())
		}
	}()

	// Create logger with TraceID and SpanID automatically included
	logger.InitOtelZapLogger("develop")
	defer logger.Sync()

	poolingRepositories := repositories.NewPoolingRepository()
	poolingService := services.NewPoolinClientService(poolingRepositories)
	poolingHandlers := handlers.NewPoolingClientHandler(poolingService)

	handler := router.InitRouter(ctx, poolingHandlers)

	s := &http.Server{
		Addr:    PORT,
		Handler: handler,
		//ReadTimeout: 30 * time.Second, // customize http.Server timeouts
	}
	logger.Info(ctx, "Server runnig on http://localhost"+s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
