package handlers

import (
	"context"
	"log"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/pkg/trace"
)

func init() {

	ctx := context.WithValue(context.Background(), "", "")

	tp, err := trace.SetupTracer(ctx, dto.APP_NAME)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Panic(ctx, err.Error())
		}
	}()

	// tracer := otel.Tracer(logger.APP_NAME)
	// ctx, span := tracer.Start(ctx, logger.PROJECT_RATELIMIT)
	// defer span.End()

	// Create logger with TraceID and SpanID automatically included
	logger.InitOtelZapLogger("develop")
	defer logger.Sync()
}
